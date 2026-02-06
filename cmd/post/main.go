package main

import (
	"os"

	"github.com/yagihash/ghat/actions"
	"github.com/yagihash/ghat/client"
	"github.com/yagihash/ghat/input"
)

const (
	exitOK = iota
	exitErr
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	token, err := actions.GetState("token")
	if err != nil {
		actions.LogError("token is not found in state")
		return exitErr
	}

	args, err := input.Load()
	if err != nil {
		actions.LogError(err.Error())
		return exitErr
	}

	c := client.New(args.BaseURL, token)
	if err := c.DeleteInstallationAccessToken(); err != nil {
		actions.LogError(err.Error())
		return exitErr
	}

	actions.LogNotice("Successfully deleted installation access token")

	return exitOK
}
