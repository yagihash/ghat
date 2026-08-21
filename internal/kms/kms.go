package kms

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strconv"
	"strings"

	kmsv1 "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
	"github.com/googleapis/gax-go/v2"
	"google.golang.org/api/iterator"
)

// KMSClient defines the interface for KMS operations
type KMSClient interface {
	AsymmetricSign(ctx context.Context, req *kmspb.AsymmetricSignRequest, opts ...gax.CallOption) (*kmspb.AsymmetricSignResponse, error)
	LatestEnabledKeyVersion(ctx context.Context, parent string) (string, error)
	Close() error
}

type Signer struct {
	client     KMSClient
	keyPath    string
	resolution ResolveResult
}

// kmsClientWrapper wraps the real KMS client to implement KMSClient.
type kmsClientWrapper struct {
	inner *kmsv1.KeyManagementClient
}

func (w *kmsClientWrapper) AsymmetricSign(ctx context.Context, req *kmspb.AsymmetricSignRequest, opts ...gax.CallOption) (*kmspb.AsymmetricSignResponse, error) {
	return w.inner.AsymmetricSign(ctx, req, opts...)
}

func (w *kmsClientWrapper) Close() error {
	return w.inner.Close()
}

// LatestEnabledKeyVersion lists all enabled versions of the given crypto key and
// returns the one with the highest version number.
func (w *kmsClientWrapper) LatestEnabledKeyVersion(ctx context.Context, parent string) (string, error) {
	it := w.inner.ListCryptoKeyVersions(ctx, &kmspb.ListCryptoKeyVersionsRequest{
		Parent: parent,
		Filter: "state=ENABLED",
	})

	var latest int64
	var found bool
	for {
		v, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return "", fmt.Errorf("failed to list crypto key versions: %w", err)
		}
		_, last, _ := strings.CutLast(v.Name, "/")
		n, err := strconv.ParseInt(last, 10, 64)
		if err != nil {
			continue
		}
		if n > latest {
			latest = n
			found = true
		}
	}

	if !found {
		return "", fmt.Errorf("no enabled key versions found for %s", parent)
	}
	return strconv.FormatInt(latest, 10), nil
}

// NewKMSClient creates a real KMS client
func NewKMSClient(ctx context.Context) (KMSClient, error) {
	client, err := kmsv1.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create kms client: %w", err)
	}
	return &kmsClientWrapper{inner: client}, nil
}

// ResolveResult holds the outcome of key version resolution.
type ResolveResult struct {
	Version      string // the key version that will be used
	AutoDetected bool   // true when version was empty and auto-detection was attempted
	DetectionErr error  // non-nil when auto-detection failed and "1" was used as fallback
}

// resolveVersion returns version if non-empty. Otherwise it attempts to detect
// the latest enabled key version via the KMS API. If detection fails (e.g. the SA
// lacks cloudkms.cryptoKeyVersions.list), it falls back to "1".
func resolveVersion(ctx context.Context, client KMSClient, projectID, location, keyRingID, keyID, version string) ResolveResult {
	if version != "" {
		return ResolveResult{Version: version}
	}
	parent := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s",
		projectID, location, keyRingID, keyID)
	resolved, err := client.LatestEnabledKeyVersion(ctx, parent)
	if err != nil {
		return ResolveResult{Version: "1", AutoDetected: true, DetectionErr: err}
	}
	return ResolveResult{Version: resolved, AutoDetected: true}
}

// NewSigner creates a new Signer with a KMS client.
// If version is empty, the latest enabled key version is detected automatically.
// Falls back to version "1" when the SA lacks the required list permission.
// The resolution outcome is available via Signer.Resolution.
func NewSigner(ctx context.Context, projectID, location, keyRingID, keyID, version string) (*Signer, error) {
	client, err := NewKMSClient(ctx)
	if err != nil {
		return nil, err
	}

	result := resolveVersion(ctx, client, projectID, location, keyRingID, keyID, version)
	s := newSigner(client, projectID, location, keyRingID, keyID, result.Version)
	s.resolution = result
	return s, nil
}

// KeyVersionInfo returns the key version that will be used for signing and whether
// it was detected automatically (i.e. kms_key_version was not explicitly specified).
// When auto-detection failed, autoDetected is true and version is "1" (fallback).
func (s *Signer) KeyVersionInfo() (version string, autoDetected bool) {
	return s.resolution.Version, s.resolution.AutoDetected
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
