package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/yagihash/ghat/v2/internal/actions"
	"github.com/yagihash/ghat/v2/internal/client"
	"github.com/yagihash/ghat/v2/internal/input"
	"github.com/yagihash/ghat/v2/internal/jwt"
	"github.com/yagihash/ghat/v2/internal/kms"
)

const (
	exitOK = iota
	exitErr
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

	signedJWT, err := jwt.Build(ctx, signer, args.AppID, time.Now())
	if err != nil {
		actions.LogError("failed to build jwt: " + err.Error())
		return exitErr
	}

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
