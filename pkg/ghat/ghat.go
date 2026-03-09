// Package ghat provides a public API for generating and revoking GitHub App
// installation access tokens using Google Cloud KMS for JWT signing.
//
// Typical usage:
//
//	signer, err := ghat.NewSigner(ctx, projectID, location, keyRingID, keyID, version)
//	if err != nil { ... }
//	defer signer.Close()
//
//	app := ghat.New(appID, signer, "")
//	token, err := app.GetGitHubAppToken(ctx, owner, nil, nil)
//	if err != nil { ... }
//	// use token ...
//	if err := app.RevokeGitHubAppToken(ctx, token); err != nil { ... }
package ghat

import (
	"context"
	"time"

	"github.com/yagihash/ghat/v2/internal/client"
	"github.com/yagihash/ghat/v2/internal/jwt"
)

// signerIface is the minimal interface required by App for JWT signing.
// *Signer.inner (*kms.Signer) satisfies this interface.
type signerIface interface {
	Sign(ctx context.Context, data []byte) ([]byte, error)
}

// App orchestrates GitHub App JWT signing, token issuance, and token revocation.
type App struct {
	appID   string
	baseURL string
	signer  signerIface
}

// New constructs an App.
// signer must be obtained from NewSigner.
// baseURL is the GitHub API base URL; pass "" to use "https://api.github.com".
func New(appID string, signer *Signer, baseURL string) *App {
	if baseURL == "" {
		baseURL = "https://api.github.com"
	}
	return &App{
		appID:   appID,
		baseURL: baseURL,
		signer:  signer.inner,
	}
}

// newApp is the internal constructor used by tests to inject a mock signer.
func newApp(appID string, signer signerIface, baseURL string) *App {
	if baseURL == "" {
		baseURL = "https://api.github.com"
	}
	return &App{
		appID:   appID,
		baseURL: baseURL,
		signer:  signer,
	}
}

// GetGitHubAppToken generates a signed JWT, resolves the GitHub App installation for
// the given owner, and returns an installation access token.
// This satisfies requirement 3: GitHub App Token issuance.
//
// permissions is a map of permission name to access level (e.g. "contents": "read").
// Pass nil to request the default permissions of the GitHub App installation.
//
// repositories is an optional list of repository names to scope the token to.
// Pass nil to grant access to all repositories the installation can access.
func (a *App) CreateGitHubAppToken(ctx context.Context, owner string, permissions map[string]string, repositories []string) (string, error) {
	signedJWT, err := jwt.Build(ctx, a.signer, a.appID, time.Now())
	if err != nil {
		return "", err
	}

	c := client.New(a.baseURL, signedJWT)

	installation, err := c.GetInstallationByOwner(owner)
	if err != nil {
		return "", err
	}

	accessToken, err := c.GetInstallationAccessToken(installation.ID, permissions, repositories)
	if err != nil {
		return "", err
	}

	return accessToken.Token, nil
}

// RevokeGitHubAppToken revokes an installation access token.
// This satisfies requirement 4: GitHub App Token revocation.
//
// token is the value previously returned by GetGitHubAppToken.
func (a *App) RevokeGitHubAppToken(ctx context.Context, token string) error {
	c := client.New(a.baseURL, token)
	return c.DeleteInstallationAccessToken()
}
