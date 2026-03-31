# Changelog

## [v2.1.0](https://github.com/yagihash/ghat/compare/v2...v2.1.0) - 2026-03-31

## [v1.0.0](https://github.com/yagihash/ghat/compare/v1.0.0...v2) - 2026-03-31
- add tests for actions by @yagihash in https://github.com/yagihash/ghat/pull/8
- add tests for kms by @yagihash in https://github.com/yagihash/ghat/pull/9
- configure renovate by @yagihash in https://github.com/yagihash/ghat/pull/10
- update renovate config by @yagihash in https://github.com/yagihash/ghat/pull/11
- add pinact by @yagihash in https://github.com/yagihash/ghat/pull/12
- fix typo kms_key_version by @yagihash in https://github.com/yagihash/ghat/pull/13
- explicitly input default key version by @yagihash in https://github.com/yagihash/ghat/pull/14
- add simple support for local execution by @yagihash in https://github.com/yagihash/ghat/pull/15
- v2 modules by @yagihash in https://github.com/yagihash/ghat/pull/16
- go mod tidy by @yagihash in https://github.com/yagihash/ghat/pull/17
- update actions by @yagihash in https://github.com/yagihash/ghat/pull/19
- chore(deps): update alpine docker tag to v3.23.3 by @renovate[bot] in https://github.com/yagihash/ghat/pull/20
- chore(deps): update golang docker tag to v1.26 by @renovate[bot] in https://github.com/yagihash/ghat/pull/21
- fix(deps): update module cloud.google.com/go/kms to v1.26.0 by @renovate[bot] in https://github.com/yagihash/ghat/pull/22
- fix(deps): update module github.com/yagihash/ghat to v2 by @renovate[bot] in https://github.com/yagihash/ghat/pull/24
- chore(config): migrate Renovate config by @renovate[bot] in https://github.com/yagihash/ghat/pull/25
- Add CLAUDE.md: AI assistant guide for ghat project by @yagihash in https://github.com/yagihash/ghat/pull/26
- refactor: use t.Setenv() instead of os.Setenv/os.Unsetenv in tests by @yagihash in https://github.com/yagihash/ghat/pull/27
- Add test for empty KeyVersion environment variable handling by @yagihash in https://github.com/yagihash/ghat/pull/35
- refactor: define named constants for JWT expiry durations by @yagihash in https://github.com/yagihash/ghat/pull/30
- fix: use LogWarning instead of LogError for Close failures by @yagihash in https://github.com/yagihash/ghat/pull/36
- fix: replace panic with error log on resource close failure by @yagihash in https://github.com/yagihash/ghat/pull/29
- fix: restrict file permissions to 0600 for state and output files by @yagihash in https://github.com/yagihash/ghat/pull/32
- fix: add newline validation to SetOutput by @yagihash in https://github.com/yagihash/ghat/pull/31
- fix: handle io.ReadAll error in GetInstallationAccessToken by @yagihash in https://github.com/yagihash/ghat/pull/34
- fix: validate HTTP status code in GetInstallationByOwner by @yagihash in https://github.com/yagihash/ghat/pull/33
- Add comprehensive unit tests for GitHub client by @yagihash in https://github.com/yagihash/ghat/pull/37
- Add Claude settings configuration for Go toolchain by @yagihash in https://github.com/yagihash/ghat/pull/39
- chore: upgrade Go version to 1.26 by @yagihash in https://github.com/yagihash/ghat/pull/41
- fix: SessionStart hookでgo.modのバージョンのGoを自動インストールするよう修正 by @yagihash in https://github.com/yagihash/ghat/pull/43
- chore(deps): pin dependency go to 1.26.1 by @renovate[bot] in https://github.com/yagihash/ghat/pull/42
- build(deps): strengthen Renovate supply chain security and enable aut… by @yagihash in https://github.com/yagihash/ghat/pull/40
- test(actions): add missing error cases and output tests by @yagihash in https://github.com/yagihash/ghat/pull/38
- chore(deps): update docker/login-action action to v4 by @renovate[bot] in https://github.com/yagihash/ghat/pull/44
- Extract public API and refactor internal packages by @yagihash in https://github.com/yagihash/ghat/pull/45
- Fix CI workflow to trigger on main branch pushes by @yagihash in https://github.com/yagihash/ghat/pull/46
- chore: migrate to mise for Go version and task management by @yagihash in https://github.com/yagihash/ghat/pull/47
- ci: add job summary to sync-permissions workflow by @yagihash in https://github.com/yagihash/ghat/pull/48
- add GHTKN_APP env var by @yagihash in https://github.com/yagihash/ghat/pull/49
- feat: log token repository scope on creation by @babarot in https://github.com/yagihash/ghat/pull/50
- feat: introduce tagpr for automated release management by @yagihash in https://github.com/yagihash/ghat/pull/51
- feat: automate GitHub Release creation on tag push by @yagihash in https://github.com/yagihash/ghat/pull/52
- feat: set deployment: false on all job environments by @yagihash in https://github.com/yagihash/ghat/pull/53
- feat: use GitHub App Token in tagpr for PR creation by @yagihash in https://github.com/yagihash/ghat/pull/54
- fix: use ubuntu-latest in tagpr to support Docker-based actions by @yagihash in https://github.com/yagihash/ghat/pull/55
- fix: push release commit to main to fix tagpr version detection by @yagihash in https://github.com/yagihash/ghat/pull/57
- Fix/retag v2 series by @yagihash in https://github.com/yagihash/ghat/pull/59

## [v1.0.0](https://github.com/yagihash/ghat/commits/v1.0.0) - 2026-01-22
- First commit by @yagihash in https://github.com/yagihash/ghat/pull/1
- add sync permissions workflow by @yagihash in https://github.com/yagihash/ghat/pull/2
- fix sync permissions workflow by @yagihash in https://github.com/yagihash/ghat/pull/3
- 🤖 Sync GitHub App Permissions by @yagihash-bot[bot] in https://github.com/yagihash/ghat/pull/4
- 🤖 Sync GitHub App Permissions by @yagihash-bot[bot] in https://github.com/yagihash/ghat/pull/5
- 🤖 Sync GitHub App Permissions by @yagihash-bot[bot] in https://github.com/yagihash/ghat/pull/6
- add ghalint wf by @yagihash in https://github.com/yagihash/ghat/pull/7
