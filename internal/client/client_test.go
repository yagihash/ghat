package client

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

type mockTransport struct {
	roundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTripFunc(req)
}

func newResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}

func newClientWithMock(baseURL, jwt string, transport http.RoundTripper) *Client {
	c := New(baseURL, jwt)
	c.HTTPClient.Transport = transport
	return c
}

func TestNew(t *testing.T) {
	c := New("https://api.github.com", "test-jwt")

	if c.BaseURL != "https://api.github.com" {
		t.Errorf("BaseURL = %q, want %q", c.BaseURL, "https://api.github.com")
	}
	if c.token != "test-jwt" {
		t.Errorf("token = %q, want %q", c.token, "test-jwt")
	}
	if c.HTTPClient == nil {
		t.Fatal("HTTPClient is nil")
	}
	if c.HTTPClient.Timeout != 10*time.Second {
		t.Errorf("Timeout = %v, want %v", c.HTTPClient.Timeout, 10*time.Second)
	}
}

func TestGetInstallationByOwner(t *testing.T) {
	tests := []struct {
		name          string
		owner         string
		roundTripFunc func(req *http.Request) (*http.Response, error)
		wantID        int64
		wantErr       bool
		checkReq      func(t *testing.T, req *http.Request)
	}{
		{
			name:  "正常系",
			owner: "myorg",
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return newResponse(http.StatusOK, `{"id": 42}`), nil
			},
			wantID:  42,
			wantErr: false,
		},
		{
			name:  "非200応答",
			owner: "myorg",
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return newResponse(http.StatusNotFound, `{"message": "Not Found"}`), nil
			},
			wantErr: true,
		},
		{
			name:  "不正なJSON",
			owner: "myorg",
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return newResponse(http.StatusOK, `not-json`), nil
			},
			wantErr: true,
		},
		{
			name:  "トランスポートエラー",
			owner: "myorg",
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErr: true,
		},
		{
			name:  "リクエストヘッダー検証",
			owner: "testowner",
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return newResponse(http.StatusOK, `{"id": 1}`), nil
			},
			wantID:  1,
			wantErr: false,
			checkReq: func(t *testing.T, req *http.Request) {
				t.Helper()
				if got := req.Header.Get("Authorization"); got != "Bearer test-jwt" {
					t.Errorf("Authorization = %q, want %q", got, "Bearer test-jwt")
				}
				if got := req.Header.Get("Accept"); got != "application/vnd.github+json" {
					t.Errorf("Accept = %q, want %q", got, "application/vnd.github+json")
				}
				if got := req.Header.Get("X-GitHub-Api-Version"); got != "2022-11-28" {
					t.Errorf("X-GitHub-Api-Version = %q, want %q", got, "2022-11-28")
				}
				if got := req.Method; got != http.MethodGet {
					t.Errorf("Method = %q, want %q", got, http.MethodGet)
				}
				wantPath := "/users/testowner/installation"
				if got := req.URL.Path; got != wantPath {
					t.Errorf("Path = %q, want %q", got, wantPath)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedReq *http.Request
			transport := &mockTransport{
				roundTripFunc: func(req *http.Request) (*http.Response, error) {
					capturedReq = req
					return tt.roundTripFunc(req)
				},
			}
			c := newClientWithMock("https://api.github.com", "test-jwt", transport)

			got, err := c.GetInstallationByOwner(tt.owner)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.ID != tt.wantID {
				t.Errorf("ID = %d, want %d", got.ID, tt.wantID)
			}
			if tt.checkReq != nil {
				tt.checkReq(t, capturedReq)
			}
		})
	}
}

func TestGetInstallationAccessToken(t *testing.T) {
	tests := []struct {
		name           string
		installationID int64
		permissions    map[string]string
		repos          []string
		roundTripFunc  func(req *http.Request) (*http.Response, error)
		wantToken      string
		wantErr        bool
		checkReq       func(t *testing.T, req *http.Request)
	}{
		{
			name:           "正常系",
			installationID: 123,
			permissions:    map[string]string{"contents": "read"},
			repos:          []string{"repo1", "repo2"},
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return newResponse(http.StatusCreated, `{"token": "ghs_xxx"}`), nil
			},
			wantToken: "ghs_xxx",
			wantErr:   false,
		},
		{
			name:           "非201応答",
			installationID: 123,
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return newResponse(http.StatusUnprocessableEntity, `{"message": "Validation Failed"}`), nil
			},
			wantErr: true,
		},
		{
			name:           "不正なJSON",
			installationID: 123,
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return newResponse(http.StatusCreated, `not-json`), nil
			},
			wantErr: true,
		},
		{
			name:           "トランスポートエラー",
			installationID: 123,
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErr: true,
		},
		{
			name:           "リクエスト検証",
			installationID: 456,
			permissions:    map[string]string{"issues": "write"},
			repos:          []string{"myrepo"},
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return newResponse(http.StatusCreated, `{"token": "ghs_yyy"}`), nil
			},
			wantToken: "ghs_yyy",
			wantErr:   false,
			checkReq: func(t *testing.T, req *http.Request) {
				t.Helper()
				if got := req.Method; got != http.MethodPost {
					t.Errorf("Method = %q, want %q", got, http.MethodPost)
				}
				wantPath := "/app/installations/456/access_tokens"
				if got := req.URL.Path; got != wantPath {
					t.Errorf("Path = %q, want %q", got, wantPath)
				}
				if got := req.Header.Get("Content-Type"); got != "application/json" {
					t.Errorf("Content-Type = %q, want %q", got, "application/json")
				}

				body, err := io.ReadAll(req.Body)
				if err != nil {
					t.Fatalf("failed to read body: %v", err)
				}
				var payload AccessTokenRequest
				if err := json.Unmarshal(body, &payload); err != nil {
					t.Fatalf("failed to unmarshal body: %v", err)
				}
				if len(payload.Repositories) != 1 || payload.Repositories[0] != "myrepo" {
					t.Errorf("Repositories = %v, want [myrepo]", payload.Repositories)
				}
				if payload.Permissions["issues"] != "write" {
					t.Errorf("Permissions[issues] = %q, want %q", payload.Permissions["issues"], "write")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedReq *http.Request
			transport := &mockTransport{
				roundTripFunc: func(req *http.Request) (*http.Response, error) {
					capturedReq = req
					return tt.roundTripFunc(req)
				},
			}
			c := newClientWithMock("https://api.github.com", "test-jwt", transport)

			got, err := c.GetInstallationAccessToken(tt.installationID, tt.permissions, tt.repos)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Token != tt.wantToken {
				t.Errorf("Token = %q, want %q", got.Token, tt.wantToken)
			}
			if tt.checkReq != nil {
				tt.checkReq(t, capturedReq)
			}
		})
	}
}

func TestDeleteInstallationAccessToken(t *testing.T) {
	tests := []struct {
		name          string
		roundTripFunc func(req *http.Request) (*http.Response, error)
		wantErr       bool
		checkReq      func(t *testing.T, req *http.Request)
	}{
		{
			name: "正常系",
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return newResponse(http.StatusNoContent, ""), nil
			},
			wantErr: false,
		},
		{
			name: "非204応答",
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return newResponse(http.StatusInternalServerError, "Internal Server Error"), nil
			},
			wantErr: true,
		},
		{
			name: "トランスポートエラー",
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErr: true,
		},
		{
			name: "リクエスト検証",
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return newResponse(http.StatusNoContent, ""), nil
			},
			wantErr: false,
			checkReq: func(t *testing.T, req *http.Request) {
				t.Helper()
				if got := req.Method; got != http.MethodDelete {
					t.Errorf("Method = %q, want %q", got, http.MethodDelete)
				}
				wantPath := "/installation/token"
				if got := req.URL.Path; got != wantPath {
					t.Errorf("Path = %q, want %q", got, wantPath)
				}
				if got := req.Header.Get("Authorization"); got != "Bearer test-jwt" {
					t.Errorf("Authorization = %q, want %q", got, "Bearer test-jwt")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedReq *http.Request
			transport := &mockTransport{
				roundTripFunc: func(req *http.Request) (*http.Response, error) {
					capturedReq = req
					return tt.roundTripFunc(req)
				},
			}
			c := newClientWithMock("https://api.github.com", "test-jwt", transport)

			err := c.DeleteInstallationAccessToken()

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.checkReq != nil {
				tt.checkReq(t, capturedReq)
			}
		})
	}
}
