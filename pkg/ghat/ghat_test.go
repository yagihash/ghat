package ghat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// mockSigner satisfies signerIface for tests without requiring real KMS.
type mockSigner struct {
	signFn func(ctx context.Context, data []byte) ([]byte, error)
}

func (m *mockSigner) Sign(ctx context.Context, data []byte) ([]byte, error) {
	return m.signFn(ctx, data)
}

// fakeSig is a minimal signature that produces a valid base64url encoding.
var fakeSig = []byte("fakesig")

func successfulSigner() *mockSigner {
	return &mockSigner{signFn: func(ctx context.Context, data []byte) ([]byte, error) {
		return fakeSig, nil
	}}
}

func failingSigner(errMsg string) *mockSigner {
	return &mockSigner{signFn: func(ctx context.Context, data []byte) ([]byte, error) {
		return nil, errors.New(errMsg)
	}}
}

// newResponse builds an *http.Response suitable for use in httptest handlers.
func jsonResponse(w http.ResponseWriter, statusCode int, body string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = io.WriteString(w, body)
}

func TestNew_DefaultBaseURL(t *testing.T) {
	s := &mockSigner{signFn: func(ctx context.Context, data []byte) ([]byte, error) { return fakeSig, nil }}
	app := newApp("123", s, "")
	if app.baseURL != "https://api.github.com" {
		t.Errorf("baseURL = %q, want %q", app.baseURL, "https://api.github.com")
	}
}

func TestNew_CustomBaseURL(t *testing.T) {
	s := &mockSigner{signFn: func(ctx context.Context, data []byte) ([]byte, error) { return fakeSig, nil }}
	app := newApp("123", s, "https://github.example.com/api/v3")
	if app.baseURL != "https://github.example.com/api/v3" {
		t.Errorf("baseURL = %q, want %q", app.baseURL, "https://github.example.com/api/v3")
	}
}

func TestApp_CreateGitHubAppToken(t *testing.T) {
	tests := []struct {
		name        string
		signer      *mockSigner
		handler     http.HandlerFunc
		owner       string
		permissions map[string]string
		repos       []string
		wantToken   string
		wantErr     bool
	}{
		{
			name:   "happy path returns token",
			signer: successfulSigner(),
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch {
				case strings.Contains(r.URL.Path, "/installation") && r.Method == http.MethodGet:
					jsonResponse(w, http.StatusOK, `{"id": 42}`)
				case strings.Contains(r.URL.Path, "/access_tokens") && r.Method == http.MethodPost:
					jsonResponse(w, http.StatusCreated, `{"token": "ghs_testtoken"}`)
				default:
					http.Error(w, "unexpected path: "+r.URL.Path, http.StatusNotFound)
				}
			}),
			owner:     "myorg",
			wantToken: "ghs_testtoken",
			wantErr:   false,
		},
		{
			name:    "KMS sign error",
			signer:  failingSigner("sign failed"),
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
			owner:   "myorg",
			wantErr: true,
		},
		{
			name:   "installation not found",
			signer: successfulSigner(),
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				jsonResponse(w, http.StatusNotFound, `{"message": "Not Found"}`)
			}),
			owner:   "myorg",
			wantErr: true,
		},
		{
			name:   "access token creation fails",
			signer: successfulSigner(),
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch {
				case strings.Contains(r.URL.Path, "/installation") && r.Method == http.MethodGet:
					jsonResponse(w, http.StatusOK, `{"id": 99}`)
				case strings.Contains(r.URL.Path, "/access_tokens") && r.Method == http.MethodPost:
					jsonResponse(w, http.StatusUnprocessableEntity, `{"message": "Validation Failed"}`)
				}
			}),
			owner:   "myorg",
			wantErr: true,
		},
		{
			name:   "permissions and repos are forwarded",
			signer: successfulSigner(),
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch {
				case strings.Contains(r.URL.Path, "/installation") && r.Method == http.MethodGet:
					jsonResponse(w, http.StatusOK, `{"id": 7}`)
				case strings.Contains(r.URL.Path, "/access_tokens") && r.Method == http.MethodPost:
					body, _ := io.ReadAll(r.Body)
					var req struct {
						Repositories []string          `json:"repositories"`
						Permissions  map[string]string `json:"permissions"`
					}
					if err := json.Unmarshal(body, &req); err != nil {
						http.Error(w, "bad body", http.StatusBadRequest)
						return
					}
					if len(req.Repositories) != 1 || req.Repositories[0] != "myrepo" {
						http.Error(w, fmt.Sprintf("unexpected repos: %v", req.Repositories), http.StatusBadRequest)
						return
					}
					if req.Permissions["contents"] != "read" {
						http.Error(w, fmt.Sprintf("unexpected perms: %v", req.Permissions), http.StatusBadRequest)
						return
					}
					jsonResponse(w, http.StatusCreated, `{"token": "ghs_scoped"}`)
				}
			}),
			owner:       "myorg",
			permissions: map[string]string{"contents": "read"},
			repos:       []string{"myrepo"},
			wantToken:   "ghs_scoped",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(tt.handler)
			defer srv.Close()

			app := newApp("12345", tt.signer, srv.URL)
			got, err := app.CreateGitHubAppToken(context.Background(), tt.owner, tt.permissions, tt.repos)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.wantToken {
				t.Errorf("token = %q, want %q", got, tt.wantToken)
			}
		})
	}
}

func TestApp_RevokeGitHubAppToken(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		token   string
		wantErr bool
	}{
		{
			name: "happy path",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodDelete {
					http.Error(w, "expected DELETE", http.StatusMethodNotAllowed)
					return
				}
				w.WriteHeader(http.StatusNoContent)
			}),
			token:   "ghs_sometoken",
			wantErr: false,
		},
		{
			name: "server error",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}),
			token:   "ghs_sometoken",
			wantErr: true,
		},
		{
			name: "authorization header contains token",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if got := r.Header.Get("Authorization"); got != "Bearer ghs_checktoken" {
					http.Error(w, "wrong auth: "+got, http.StatusUnauthorized)
					return
				}
				w.WriteHeader(http.StatusNoContent)
			}),
			token:   "ghs_checktoken",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(tt.handler)
			defer srv.Close()

			// RevokeGitHubAppToken only needs baseURL and token; signer is irrelevant.
			app := newApp("12345", successfulSigner(), srv.URL)
			err := app.RevokeGitHubAppToken(context.Background(), tt.token)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
