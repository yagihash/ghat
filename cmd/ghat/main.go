package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/yagihash/ghat/v2/actions"
	"github.com/yagihash/ghat/v2/client"
	"github.com/yagihash/ghat/v2/input"
	"github.com/yagihash/ghat/v2/kms"
)

const (
	exitOK = iota
	exitErr
)

const (
	// jwtIssuedAtSkew adjusts the iat claim into the past to account for
	// clock skew between the local machine and GitHub's servers.
	jwtIssuedAtSkew = -60 * time.Second
	// jwtExpiry is the duration for which the JWT is valid.
	// GitHub Apps require JWTs to expire within 10 minutes.
	jwtExpiry = 600 * time.Second
)

var isActions = os.Getenv("GITHUB_ACTIONS") == "true"

func main() {
	os.Exit(realMain())
}

func realMain() int {
	ctx := context.Background()

	args, err := input.Load()
	if err != nil {
		actions.LogError("failed to load inputs: " + err.Error())
		return exitErr
	}

	signer, err := kms.NewSigner(ctx, args.ProjectID, args.Location, args.KeyRingID, args.KeyID, args.KeyVersion)
	if err != nil {
		actions.LogError("failed to create signer: " + err.Error())
		return exitErr
	}
	defer func(signer *kms.Signer) {
		if err := signer.Close(); err != nil {
			actions.LogWarning("failed to close KMS signer: " + err.Error())
		}
	}(signer)

	now := time.Now()
	iat := now.Add(jwtIssuedAtSkew).Unix()
	exp := now.Add(jwtExpiry).Unix()

	header := map[string]any{
		"typ": "token",
		"alg": "RS256",
	}

	payload := map[string]any{
		"iat": iat,
		"exp": exp,
		"iss": args.AppID,
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		actions.LogError("failed to marshal jwt header: " + err.Error())
		return exitErr
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		actions.LogError("failed to marshal jwt payload: " + err.Error())
		return exitErr
	}

	headerBase64url := base64.RawURLEncoding.EncodeToString(headerJSON)
	payloadBase64url := base64.RawURLEncoding.EncodeToString(payloadJSON)

	unsignedJWT := fmt.Sprintf("%s.%s", headerBase64url, payloadBase64url)

	sig, err := signer.Sign(ctx, []byte(unsignedJWT))
	if err != nil {
		actions.LogError("failed to sign jwt: " + err.Error())
		return exitErr
	}

	signatureBase64url := base64.RawURLEncoding.EncodeToString(sig)

	signedJWT := fmt.Sprintf("%s.%s", unsignedJWT, signatureBase64url)

	c := client.New(args.BaseURL, signedJWT)

	installation, err := c.GetInstallationByOwner(args.Owner)
	if err != nil {
		actions.LogError("failed to get installation: " + err.Error())
		return exitErr
	}

	accessToken, err := c.GetInstallationAccessToken(installation.ID, args.Permissions, args.Repositories)
	if err != nil {
		actions.LogError("failed to get access token: " + err.Error())
		return exitErr
	}
	if isActions {
		actions.AddMask(accessToken.Token)

		if err := actions.SetState("token", accessToken.Token); err != nil {
			actions.LogError(err.Error())
			return exitErr
		}

		if err := actions.SetOutput("token", accessToken.Token); err != nil {
			actions.LogError(err.Error())
			return exitErr
		}
	} else {
		fmt.Print(accessToken.Token)
	}

	return exitOK
}
