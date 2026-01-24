package kms

import (
	"context"
	"crypto/sha256"
	"fmt"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
	"github.com/googleapis/gax-go/v2"
)

// KMSClient defines the interface for KMS operations
type KMSClient interface {
	AsymmetricSign(ctx context.Context, req *kmspb.AsymmetricSignRequest, opts ...gax.CallOption) (*kmspb.AsymmetricSignResponse, error)
	Close() error
}

type Signer struct {
	client  KMSClient
	keyPath string
}

// NewKMSClient creates a real KMS client
func NewKMSClient(ctx context.Context) (KMSClient, error) {
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create kms client: %w", err)
	}
	return client, nil
}

// NewSigner creates a new Signer with a KMS client
func NewSigner(ctx context.Context, projectID, location, keyRingID, keyID, version string) (*Signer, error) {
	client, err := NewKMSClient(ctx)
	if err != nil {
		return nil, err
	}

	return newSigner(client, projectID, location, keyRingID, keyID, version), nil
}

// newSigner creates a new Signer with the given KMS client (for testing)
func newSigner(client KMSClient, projectID, location, keyRingID, keyID, version string) *Signer {
	path := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s/cryptoKeyVersions/%s",
		projectID, location, keyRingID, keyID, version)

	return &Signer{
		client:  client,
		keyPath: path,
	}
}

func (s *Signer) Sign(ctx context.Context, data []byte) ([]byte, error) {
	digest := sha256.Sum256(data)

	req := &kmspb.AsymmetricSignRequest{
		Name: s.keyPath,
		Digest: &kmspb.Digest{
			Digest: &kmspb.Digest_Sha256{
				Sha256: digest[:],
			},
		},
	}

	result, err := s.client.AsymmetricSign(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to asymmetric sign: %w", err)
	}

	return result.Signature, nil
}

func (s *Signer) Close() error {
	return s.client.Close()
}
