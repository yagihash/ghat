package mini_gh_sts

issuer := "https://token.actions.githubusercontent.com"

permissions := {
	"contents": "write",
	"pull_requests": "write",
}

default allow := false

allow if {
	input.sub == "repo:yagihash/ghat:environment:main"
}

allow if {
	input.sub == "repo:yagihash/ghat:refs/heads/main"
}
