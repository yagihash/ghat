package main

import (
	"github.com/yagihash/ghat/actions"
	"github.com/yagihash/ghat/client"
	"github.com/yagihash/ghat/input"
)

func main() {
	token, err := actions.GetState("token")
	if err != nil {
		actions.LogError("token is not found in state")
		return
	}

	args, err := input.Load()
	if err != nil {
		actions.LogError(err.Error())
		return
	}

	c := client.New(args.BaseURL, token)
	if err := c.DeleteInstallationAccessToken(); err != nil {
		actions.LogError(err.Error())
		return
	}

	actions.LogNotice("Successfully deleted installation access token")
}
