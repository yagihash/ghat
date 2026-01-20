package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/yagihash/ghat/actions"
	"github.com/yagihash/ghat/client"
	"github.com/yagihash/ghat/input"
	"github.com/yagihash/ghat/kms"
)

func main() {
	ctx := context.Background()

	args, err := input.Load()
	if err != nil {
		actions.LogError("failed to load inputs: " + err.Error())
		return
	}

	signer, err := kms.NewSigner(ctx, args.ProjectID, args.Location, args.KeyRingID, args.KeyID, "1")
	if err != nil {
		actions.LogError("failed to create signer: " + err.Error())
		return
	}
	defer func(signer *kms.Signer) {
		err := signer.Close()
		if err != nil {
			panic(err)
		}
	}(signer)

	now := time.Now()
	iat := now.Add(-60 * time.Second).Unix()
	exp := now.Add(600 * time.Second).Unix()

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
		return
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		actions.LogError("failed to marshal jwt payload: " + err.Error())
		return
	}

	headerBase64url := base64.RawURLEncoding.EncodeToString(headerJSON)
	payloadBase64url := base64.RawURLEncoding.EncodeToString(payloadJSON)

	unsignedJWT := fmt.Sprintf("%s.%s", headerBase64url, payloadBase64url)

	sig, err := signer.Sign(ctx, []byte(unsignedJWT))
	if err != nil {
		actions.LogError("failed to sign jwt: " + err.Error())
		return
	}

	signatureBase64url := base64.RawURLEncoding.EncodeToString(sig)

	signedJWT := fmt.Sprintf("%s.%s", unsignedJWT, signatureBase64url)

	c := client.New(args.BaseURL, signedJWT)

	res, err := c.GetInstallationByOwner(args.Owner)
	if err != nil {
		actions.LogError("failed to get installation: " + err.Error())
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		actions.LogError("failed to read body of installation response: " + err.Error())
		return
	}

	var installation client.InstallationResponse
	if err := json.Unmarshal(body, &installation); err != nil {
		actions.LogError("failed to unmarshal installation response: " + err.Error())
		return
	}

	accessToken, err := c.GetInstallationAccessToken(installation.ID, map[string]string{}, args.Repositories)
	if err != nil {
		actions.LogError("failed to get access token: " + err.Error())
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(res.Body)

	actions.AddMask(accessToken.Token)

	if err := actions.SetState("token", accessToken.Token); err != nil {
		actions.LogError(err.Error())
	}

	if err := actions.SetOutput("token", accessToken.Token); err != nil {
		actions.LogError(err.Error())
	}
}
