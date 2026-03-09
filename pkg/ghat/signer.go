package ghat

import (
	"context"

	"github.com/yagihash/ghat/v2/internal/kms"
)

// Signer signs data using a Google Cloud KMS asymmetric key.
// It wraps the internal KMS implementation and is the entry point for
// requirement 1 (KMS access) and requirement 2 (JWT signing).
type Signer struct {
	inner *kms.Signer
}

// NewSigner creates a Signer backed by Google Cloud KMS.
// projectID, location, keyRingID, keyID, and version identify the CryptoKeyVersion.
func NewSigner(ctx context.Context, projectID, location, keyRingID, keyID, version string) (*Signer, error) {
	s, err := kms.NewSigner(ctx, projectID, location, keyRingID, keyID, version)
	if err != nil {
		return nil, err
	}
	return &Signer{inner: s}, nil
}

// Close releases the underlying KMS client connection.
func (s *Signer) Close() error {
	return s.inner.Close()
}
