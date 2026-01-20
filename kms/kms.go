package kms

import (
	"context"
	"crypto/sha256"
	"fmt"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
)

type Signer struct {
	client  *kms.KeyManagementClient
	keyPath string
}

func NewSigner(ctx context.Context, projectID, location, keyRingID, keyID, version string) (*Signer, error) {
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create kms client: %w", err)
	}

	path := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s/cryptoKeyVersions/%s",
		projectID, location, keyRingID, keyID, version)

	return &Signer{
		client:  client,
		keyPath: path,
	}, nil
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
