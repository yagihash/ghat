package jwt

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"
)

type mockSigner struct {
	signFn func(ctx context.Context, data []byte) ([]byte, error)
}

func (m *mockSigner) Sign(ctx context.Context, data []byte) ([]byte, error) {
	return m.signFn(ctx, data)
}

func TestBuild(t *testing.T) {
	fixedNow := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	fakeSig := []byte("fakesignature")

	tests := []struct {
		name        string
		appID       string
		signFn      func(ctx context.Context, data []byte) ([]byte, error)
		wantErr     bool
		validateJWT func(t *testing.T, jwt string)
	}{
		{
			name:  "happy path produces valid JWT structure",
			appID: "12345",
			signFn: func(ctx context.Context, data []byte) ([]byte, error) {
				return fakeSig, nil
			},
			wantErr: false,
			validateJWT: func(t *testing.T, jwt string) {
				t.Helper()
				parts := strings.Split(jwt, ".")
				if len(parts) != 3 {
					t.Fatalf("expected 3 dot-separated parts, got %d: %q", len(parts), jwt)
				}

				// Decode and validate header
				headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
				if err != nil {
					t.Fatalf("failed to decode header: %v", err)
				}
				var header map[string]any
				if err := json.Unmarshal(headerBytes, &header); err != nil {
					t.Fatalf("failed to unmarshal header: %v", err)
				}
				if header["typ"] != "token" {
					t.Errorf("header[typ] = %q, want %q", header["typ"], "token")
				}
				if header["alg"] != "RS256" {
					t.Errorf("header[alg] = %q, want %q", header["alg"], "RS256")
				}

				// Decode and validate payload
				payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
				if err != nil {
					t.Fatalf("failed to decode payload: %v", err)
				}
				var payload map[string]any
				if err := json.Unmarshal(payloadBytes, &payload); err != nil {
					t.Fatalf("failed to unmarshal payload: %v", err)
				}
				if payload["iss"] != "12345" {
					t.Errorf("payload[iss] = %q, want %q", payload["iss"], "12345")
				}
				wantIat := float64(fixedNow.Add(IssuedAtSkew).Unix())
				if payload["iat"] != wantIat {
					t.Errorf("payload[iat] = %v, want %v", payload["iat"], wantIat)
				}
				wantExp := float64(fixedNow.Add(Expiry).Unix())
				if payload["exp"] != wantExp {
					t.Errorf("payload[exp] = %v, want %v", payload["exp"], wantExp)
				}

				// Validate signature
				wantSig := base64.RawURLEncoding.EncodeToString(fakeSig)
				if parts[2] != wantSig {
					t.Errorf("signature = %q, want %q", parts[2], wantSig)
				}
			},
		},
		{
			name:  "sign error is propagated",
			appID: "12345",
			signFn: func(ctx context.Context, data []byte) ([]byte, error) {
				return nil, errors.New("KMS unavailable")
			},
			wantErr: true,
		},
		{
			name:  "iat skew is -60s",
			appID: "99",
			signFn: func(ctx context.Context, data []byte) ([]byte, error) {
				return fakeSig, nil
			},
			wantErr: false,
			validateJWT: func(t *testing.T, jwt string) {
				t.Helper()
				parts := strings.Split(jwt, ".")
				payloadBytes, _ := base64.RawURLEncoding.DecodeString(parts[1])
				var payload map[string]any
				_ = json.Unmarshal(payloadBytes, &payload)
				wantIat := float64(fixedNow.Unix() - 60)
				if payload["iat"] != wantIat {
					t.Errorf("iat = %v, want %v (fixedNow - 60s)", payload["iat"], wantIat)
				}
			},
		},
		{
			name:  "exp is +600s",
			appID: "99",
			signFn: func(ctx context.Context, data []byte) ([]byte, error) {
				return fakeSig, nil
			},
			wantErr: false,
			validateJWT: func(t *testing.T, jwt string) {
				t.Helper()
				parts := strings.Split(jwt, ".")
				payloadBytes, _ := base64.RawURLEncoding.DecodeString(parts[1])
				var payload map[string]any
				_ = json.Unmarshal(payloadBytes, &payload)
				wantExp := float64(fixedNow.Unix() + 600)
				if payload["exp"] != wantExp {
					t.Errorf("exp = %v, want %v (fixedNow + 600s)", payload["exp"], wantExp)
				}
			},
		},
		{
			name:  "deterministic with same inputs",
			appID: "42",
			signFn: func(ctx context.Context, data []byte) ([]byte, error) {
				return fakeSig, nil
			},
			wantErr: false,
			validateJWT: func(t *testing.T, jwt string) {
				t.Helper()
				signer := &mockSigner{signFn: func(ctx context.Context, data []byte) ([]byte, error) {
					return fakeSig, nil
				}}
				jwt2, err := Build(context.Background(), signer, "42", fixedNow)
				if err != nil {
					t.Fatalf("second Build failed: %v", err)
				}
				if jwt != jwt2 {
					t.Errorf("Build is not deterministic: %q != %q", jwt, jwt2)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signer := &mockSigner{signFn: tt.signFn}
			got, err := Build(context.Background(), signer, tt.appID, fixedNow)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.validateJWT != nil {
				tt.validateJWT(t, got)
			}
		})
	}
}
