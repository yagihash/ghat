package jwt

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

const (
	// IssuedAtSkew adjusts the iat claim into the past to account for
	// clock skew between the local machine and GitHub's servers.
	IssuedAtSkew = -60 * time.Second
	// Expiry is the duration for which the JWT is valid.
	// GitHub Apps require JWTs to expire within 10 minutes.
	Expiry = 600 * time.Second
)

// Signer signs arbitrary byte slices. internal/kms.Signer satisfies this interface.
type Signer interface {
	Sign(ctx context.Context, data []byte) ([]byte, error)
}

// Build constructs and returns a signed GitHub App JWT.
// appID is the GitHub App's numeric ID (as a string).
// now is the reference time; callers should pass time.Now().
func Build(ctx context.Context, signer Signer, appID string, now time.Time) (string, error) {
	header := map[string]any{
		"typ": "token",
		"alg": "RS256",
	}

	payload := map[string]any{
		"iat": now.Add(IssuedAtSkew).Unix(),
		"exp": now.Add(Expiry).Unix(),
		"iss": appID,
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("jwt: marshal header: %w", err)
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("jwt: marshal payload: %w", err)
	}

	headerBase64url := base64.RawURLEncoding.EncodeToString(headerJSON)
	payloadBase64url := base64.RawURLEncoding.EncodeToString(payloadJSON)

	unsigned := fmt.Sprintf("%s.%s", headerBase64url, payloadBase64url)

	sig, err := signer.Sign(ctx, []byte(unsigned))
	if err != nil {
		return "", fmt.Errorf("jwt: sign: %w", err)
	}

	return fmt.Sprintf("%s.%s", unsigned, base64.RawURLEncoding.EncodeToString(sig)), nil
}
