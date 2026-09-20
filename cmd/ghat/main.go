package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
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
	if version, autoDetected := signer.KeyVersionInfo(); autoDetected {
		if version == "1" {
			actions.LogNotice("failed to detect latest KMS key version; using version \"1\" as fallback")
		} else {
			actions.LogNotice(fmt.Sprintf("auto-detected latest KMS key version: %s", version))
		}
	}

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

	if len(args.Repositories) == 0 {
		actions.LogNotice("Token will be scoped to all repositories accessible by the GitHub App installation")
	} else {
		actions.LogGroup(
			fmt.Sprintf("Token scope: %d repositories", len(args.Repositories)),
			args.Repositories...,
		)
	}

	if len(args.Permissions) == 0 {
		actions.LogNotice("Token will be granted all permissions of the GitHub App installation")
	} else {
		perms := make([]string, 0, len(args.Permissions))
		for k, v := range args.Permissions {
			perms = append(perms, k+":"+v)
		}
		actions.LogGroup(
			fmt.Sprintf("Token scope: %d permissions", len(args.Permissions)),
			perms...,
		)
	}

	accessToken, err := c.GetInstallationAccessToken(installation.ID, args.Permissions, args.Repositories)
	if err != nil {
		actions.LogError("failed to get access token: " + err.Error())
		return exitErr
	}
	if isActions {
		actions.AddMask(accessToken.Token)

		hash := sha256.Sum256([]byte(accessToken.Token))
		actions.LogNotice("Token hash: " + base64.StdEncoding.EncodeToString(hash[:]))

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
