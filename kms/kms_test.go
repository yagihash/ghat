package kms

import (
	"context"
	"crypto/sha256"
	"errors"
	"testing"

	"cloud.google.com/go/kms/apiv1/kmspb"
	"github.com/googleapis/gax-go/v2"
)

// mockKMSClient is a mock implementation of KMSClient for testing
type mockKMSClient struct {
	asymmetricSignFunc func(ctx context.Context, req *kmspb.AsymmetricSignRequest, opts ...gax.CallOption) (*kmspb.AsymmetricSignResponse, error)
	closeFunc          func() error
}

func (m *mockKMSClient) AsymmetricSign(ctx context.Context, req *kmspb.AsymmetricSignRequest, opts ...gax.CallOption) (*kmspb.AsymmetricSignResponse, error) {
	if m.asymmetricSignFunc != nil {
		return m.asymmetricSignFunc(ctx, req, opts...)
	}
	return &kmspb.AsymmetricSignResponse{}, nil
}

func (m *mockKMSClient) Close() error {
	if m.closeFunc != nil {
		return m.closeFunc()
	}
	return nil
}

func TestNewSigner(t *testing.T) {
	tests := []struct {
		name      string
		projectID string
		location  string
		keyRingID string
		keyID     string
		version   string
		wantPath  string
	}{
		{
			name:      "creates signer with correct path",
			projectID: "test-project",
			location:  "us-central1",
			keyRingID: "test-keyring",
			keyID:     "test-key",
			version:   "1",
			wantPath:  "projects/test-project/locations/us-central1/keyRings/test-keyring/cryptoKeys/test-key/cryptoKeyVersions/1",
		},
		{
			name:      "handles different parameters",
			projectID: "another-project",
			location:  "asia-northeast1",
			keyRingID: "prod-keyring",
			keyID:     "signing-key",
			version:   "2",
			wantPath:  "projects/another-project/locations/asia-northeast1/keyRings/prod-keyring/cryptoKeys/signing-key/cryptoKeyVersions/2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockKMSClient{}
			signer := newSigner(mockClient, tt.projectID, tt.location, tt.keyRingID, tt.keyID, tt.version)

			if signer == nil {
				t.Fatal("newSigner returned nil")
			}

			if signer.keyPath != tt.wantPath {
				t.Errorf("keyPath = %v, want %v", signer.keyPath, tt.wantPath)
			}

			if signer.client != mockClient {
				t.Error("client was not set correctly")
			}
		})
	}
}

func TestSigner_Sign(t *testing.T) {
	tests := []struct {
		name           string
		data           []byte
		mockResponse   *kmspb.AsymmetricSignResponse
		mockError      error
		wantSignature  []byte
		wantError      bool
		validateDigest bool
	}{
		{
			name: "successfully signs data",
			data: []byte("test data"),
			mockResponse: &kmspb.AsymmetricSignResponse{
				Signature: []byte("mock-signature"),
			},
			mockError:      nil,
			wantSignature:  []byte("mock-signature"),
			wantError:      false,
			validateDigest: true,
		},
		{
			name:          "returns error when signing fails",
			data:          []byte("test data"),
			mockResponse:  nil,
			mockError:     errors.New("signing failed"),
			wantSignature: nil,
			wantError:     true,
		},
		{
			name: "handles empty data",
			data: []byte(""),
			mockResponse: &kmspb.AsymmetricSignResponse{
				Signature: []byte("signature-for-empty"),
			},
			mockError:      nil,
			wantSignature:  []byte("signature-for-empty"),
			wantError:      false,
			validateDigest: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedDigest []byte
			mockClient := &mockKMSClient{
				asymmetricSignFunc: func(ctx context.Context, req *kmspb.AsymmetricSignRequest, opts ...gax.CallOption) (*kmspb.AsymmetricSignResponse, error) {
					if tt.validateDigest && req.Digest != nil {
						capturedDigest = req.Digest.GetSha256()
					}
					if tt.mockError != nil {
						return nil, tt.mockError
					}
					return tt.mockResponse, nil
				},
			}

			signer := newSigner(mockClient, "project", "location", "keyring", "key", "1")
			ctx := context.Background()

			signature, err := signer.Sign(ctx, tt.data)

			if tt.wantError {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if string(signature) != string(tt.wantSignature) {
				t.Errorf("signature = %v, want %v", signature, tt.wantSignature)
			}

			if tt.validateDigest {
				expectedDigest := sha256.Sum256(tt.data)
				if string(capturedDigest) != string(expectedDigest[:]) {
					t.Errorf("digest = %x, want %x", capturedDigest, expectedDigest)
				}
			}
		})
	}
}

func TestSigner_Close(t *testing.T) {
	tests := []struct {
		name      string
		mockError error
		wantError bool
	}{
		{
			name:      "successfully closes client",
			mockError: nil,
			wantError: false,
		},
		{
			name:      "returns error when close fails",
			mockError: errors.New("close failed"),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockKMSClient{
				closeFunc: func() error {
					return tt.mockError
				},
			}

			signer := newSigner(mockClient, "project", "location", "keyring", "key", "1")
			err := signer.Close()

			if tt.wantError && err == nil {
				t.Error("expected error but got nil")
			}

			if !tt.wantError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
