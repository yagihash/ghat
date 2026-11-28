# Changelog

## [v2.5.2](https://github.com/yagihash/ghat/compare/v2.5.1...v2.5.2) - 2026-11-28

## [v2.5.1](https://github.com/yagihash/ghat/compare/v2.5.0...v2.5.1) - 2026-11-16

## [v2.5.0](https://github.com/yagihash/ghat/compare/v2.4.0...v2.5.0) - 2026-10-29

## [v2.4.0](https://github.com/yagihash/ghat/compare/v2.3.1...v2.4.0) - 2026-09-20

## [v2.3.1](https://github.com/yagihash/ghat/compare/v2.3.0...v2.3.1) - 2026-08-21

## [v2.3.0](https://github.com/yagihash/ghat/compare/v2.2.8...v2.3.0) - 2026-06-23

## [v2.2.8](https://github.com/yagihash/ghat/compare/v2.2.7...v2.2.8) - 2026-06-12

## [v2.2.7](https://github.com/yagihash/ghat/compare/v2.2.6...v2.2.7) - 2026-05-17
### Dependencies
- chore(deps): update dependency suzuki-shunsuke/ghalint to v1.5.6 by @renovate[bot] in https://github.com/yagihash/ghat/pull/149
- fix(deps): update module google.golang.org/api to v0.278.0 by @renovate[bot] in https://github.com/yagihash/ghat/pull/151
- fix(deps): update module cloud.google.com/go/kms to v1.31.0 by @renovate[bot] in https://github.com/yagihash/ghat/pull/150
- chore(deps): update songmu/tagpr action to v1.19.0 by @renovate[bot] in https://github.com/yagihash/ghat/pull/155
- chore(deps): update yagihash/ghmint-action action to v1 by @renovate[bot] in https://github.com/yagihash/ghat/pull/156
- chore(deps): update golang docker tag to v1.26.3 by @renovate[bot] in https://github.com/yagihash/ghat/pull/154
- chore(deps): update dependency go to v1.26.3 by @renovate[bot] in https://github.com/yagihash/ghat/pull/153

## [v2.2.6](https://github.com/yagihash/ghat/compare/v2.2.5...v2.2.6) - 2026-05-07
### Dependencies
- chore(deps): update yagihash/ghmint-action digest to 4b6f18d by @renovate[bot] in https://github.com/yagihash/ghat/pull/141

## [v2.2.5](https://github.com/yagihash/ghat/compare/v2.2.4...v2.2.5) - 2026-04-28
### Features
- feat: log token hash to Actions execution log for incident investigation by @yagihash in https://github.com/yagihash/ghat/pull/144

## [v2.2.4](https://github.com/yagihash/ghat/compare/v2.2.3...v2.2.4) - 2026-04-19
### Dependencies
- chore(deps): update dependency go to v1.26.2 by @renovate[bot] in https://github.com/yagihash/ghat/pull/118
- chore(deps): update songmu/tagpr action to v1.18.2 by @renovate[bot] in https://github.com/yagihash/ghat/pull/119
- fix(deps): update module cloud.google.com/go/kms to v1.29.0 by @renovate[bot] in https://github.com/yagihash/ghat/pull/120
- chore(deps): update alpine docker tag to v3.23.4 by @renovate[bot] in https://github.com/yagihash/ghat/pull/125
- chore(deps): update songmu/tagpr action to v1.18.3 by @renovate[bot] in https://github.com/yagihash/ghat/pull/126
- fix(deps): update module github.com/googleapis/gax-go/v2 to v2.22.0 by @renovate[bot] in https://github.com/yagihash/ghat/pull/127
- fix(deps): update module google.golang.org/api to v0.276.0 by @renovate[bot] in https://github.com/yagihash/ghat/pull/128

## [v2.2.3](https://github.com/yagihash/ghat/compare/v2.2.2...v2.2.3) - 2026-04-11
### Bug Fixes
- fix: ensure attestation succeeds before creating GitHub Release by @yagihash in https://github.com/yagihash/ghat/pull/114

## [v2.2.2](https://github.com/yagihash/ghat/compare/v2.2.1...v2.2.2) - 2026-04-11
### Features
- feat: show token scope as collapsible log groups in GitHub Actions by @yagihash in https://github.com/yagihash/ghat/pull/113

## [v2.2.1](https://github.com/yagihash/ghat/compare/v2.2.0...v2.2.1) - 2026-04-11
### Dependencies
- fix(deps): update module github.com/googleapis/gax-go/v2 to v2.21.0 by @renovate[bot] in https://github.com/yagihash/ghat/pull/104
- chore(deps): update github actions by @renovate[bot] in https://github.com/yagihash/ghat/pull/105
- fix(deps): update module google.golang.org/api to v0.275.0 by @renovate[bot] in https://github.com/yagihash/ghat/pull/107
- fix(deps): update module cloud.google.com/go/kms to v1.27.0 by @renovate[bot] in https://github.com/yagihash/ghat/pull/106

## [v2.2.0](https://github.com/yagihash/ghat/compare/v2.1.2...v2.2.0) - 2026-04-02
- fix: pass GH_TOKEN to git push step in build-release.yml by @yagihash in https://github.com/yagihash/ghat/pull/76
- remove executable which is commited by mistake by @yagihash in https://github.com/yagihash/ghat/pull/77
- ci: remove push-edge workflow by @yagihash in https://github.com/yagihash/ghat/pull/78
- ci: retag Docker image via docker pull/tag/push instead of imagetools create by @yagihash in https://github.com/yagihash/ghat/pull/79
- fix(ci): use GitHub API to create verified commits in build-release workflow by @yagihash in https://github.com/yagihash/ghat/pull/82
- fix(ci): trigger build-release on push to tagpr-*, add verify check by @yagihash in https://github.com/yagihash/ghat/pull/83
- fix(ci): trigger build-release on push to tagpr-*, add verify-release-image check by @yagihash in https://github.com/yagihash/ghat/pull/84
- fix(ci): run verify-release-image after build via workflow_run and Checks API by @yagihash in https://github.com/yagihash/ghat/pull/85
- fix(ci): run verify-release-image after build via workflow_run and Checks API by @yagihash in https://github.com/yagihash/ghat/pull/86
- fix(ci): specify branch when fetching action.yml SHA to avoid 409 conflict by @yagihash in https://github.com/yagihash/ghat/pull/87
- separate github actions lint into two jobs by @yagihash in https://github.com/yagihash/ghat/pull/89
- ci: remove [skip ci], prevent loops with idempotency check by @yagihash in https://github.com/yagihash/ghat/pull/88
- ci: add job-level condition to skip bot digest commits by @yagihash in https://github.com/yagihash/ghat/pull/90
- chore: disable CHANGELOG generation by @yagihash in https://github.com/yagihash/ghat/pull/91
- ci: fix bot author check in verify-release-image by @yagihash in https://github.com/yagihash/ghat/pull/92
- fix(ci): use branch name as Docker image tag to fix verify check by @yagihash in https://github.com/yagihash/ghat/pull/93
- omit tagpr-temp to run build wf by @yagihash in https://github.com/yagihash/ghat/pull/95
- feat(kms): auto-detect latest enabled key version when kms_key_version is unset by @yagihash in https://github.com/yagihash/ghat/pull/80
- chore: add tab indentation to .tagpr to match tagpr's auto-generated style by @yagihash in https://github.com/yagihash/ghat/pull/98
- feat(kms): expose KeyVersionInfo and log resolved key version in main by @yagihash in https://github.com/yagihash/ghat/pull/99
- fix(ci): run build-release on tagpr force-push even when HEAD is digest commit by @yagihash in https://github.com/yagihash/ghat/pull/100

## [v2.1.2](https://github.com/yagihash/ghat/compare/v1.0.0...v2.1.2) - 2026-04-01
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
- fix: configure tagpr to use version.txt for version detection by @yagihash in https://github.com/yagihash/ghat/pull/60
- feat: verify ghalint and actionlint binaries with GitHub Attestation by @yagihash in https://github.com/yagihash/ghat/pull/62
- fix: set version.txt to v2.1.1 to align with latest released tag by @yagihash in https://github.com/yagihash/ghat/pull/64
- fix: bump version.txt to v2.1.2 to unblock tagpr by @yagihash in https://github.com/yagihash/ghat/pull/65
- feat: log permissions when creating GitHub App installation token by @yagihash in https://github.com/yagihash/ghat/pull/66
- chore(deps): update actions/setup-go action to v6.4.0 by @renovate[bot] in https://github.com/yagihash/ghat/pull/67
- chore(config): migrate Renovate config by @renovate[bot] in https://github.com/yagihash/ghat/pull/68
- remove version.txt by @yagihash in https://github.com/yagihash/ghat/pull/69
- fix release.yml by @yagihash in https://github.com/yagihash/ghat/pull/70
- refactor: build Docker image on tagpr PR, retag on release by @yagihash in https://github.com/yagihash/ghat/pull/71
- ci: add workflow_dispatch trigger to tagpr by @yagihash in https://github.com/yagihash/ghat/pull/73

## [v2.4.0](https://github.com/yagihash/ghat/compare/v2.4.0...v2) - 2026-10-29

## [v1.0.0](https://github.com/yagihash/ghat/commits/v1.0.0) - 2026-01-22
- First commit by @yagihash in https://github.com/yagihash/ghat/pull/1
- add sync permissions workflow by @yagihash in https://github.com/yagihash/ghat/pull/2
- fix sync permissions workflow by @yagihash in https://github.com/yagihash/ghat/pull/3
- 🤖 Sync GitHub App Permissions by @yagihash-bot[bot] in https://github.com/yagihash/ghat/pull/4
- 🤖 Sync GitHub App Permissions by @yagihash-bot[bot] in https://github.com/yagihash/ghat/pull/5
- 🤖 Sync GitHub App Permissions by @yagihash-bot[bot] in https://github.com/yagihash/ghat/pull/6
- add ghalint wf by @yagihash in https://github.com/yagihash/ghat/pull/7
