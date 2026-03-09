# CLAUDE.md — AI Assistant Guide for ghat

## Project Overview

**ghat** is a GitHub Action (and standalone CLI tool) that generates GitHub App installation access tokens by signing JWTs with Google Cloud KMS rather than storing private keys locally. This improves security for CI/CD pipelines by leveraging Workload Identity Federation and hardware-backed key storage.

It also exposes a public Go library (`pkg/ghat`) that external Go modules can import to perform the same operations programmatically.

**Module path:** `github.com/yagihash/ghat/v2`
**Language:** Go 1.26+
**License:** MIT

---

## Repository Structure

```
ghat/
├── cmd/
│   ├── ghat/main.go           # Main action: generates and outputs GitHub App token
│   └── post/main.go           # Post-action cleanup: revokes the token after use
├── internal/
│   ├── actions/
│   │   ├── actions.go         # GitHub Actions workflow commands (output, state, masking, logging)
│   │   └── actions_test.go
│   ├── client/
│   │   ├── client.go          # GitHub API client (installation, token endpoints)
│   │   └── client_test.go
│   ├── input/
│   │   ├── input.go           # Config struct + envconfig-based loader
│   │   └── input_test.go
│   ├── jwt/
│   │   ├── jwt.go             # GitHub App JWT construction and signing
│   │   └── jwt_test.go
│   └── kms/
│       ├── kms.go             # Google Cloud KMS signer (interface + implementation)
│       └── kms_test.go
├── pkg/
│   └── ghat/
│       ├── ghat.go            # Public API: App struct (IssueToken, RevokeToken, SignJWT)
│       ├── signer.go          # Public API: Signer type (NewSigner, Close)
│       └── ghat_test.go
├── scripts/
│   └── update-permissions.sh  # Syncs permission list in action.yml from GitHub API
├── .github/
│   └── workflows/             # CI/CD pipelines (test, edge build, release, linting)
├── action.yml                 # GitHub Action metadata (inputs, outputs, Docker entrypoint)
├── Dockerfile                 # Multi-stage build: golang:1.26-alpine → alpine:3.23.3
├── Taskfile.yml               # Task runner (build, test, coverage)
├── go.mod / go.sum
├── renovate.json
└── README.md
```

---

## Development Workflows

### Prerequisites

- Go 1.25+
- [Taskfile](https://taskfile.dev) (`task` CLI)
- Docker (for building the container image)

### Common Tasks

```bash
# Run unit tests (verbose, with coverage and race detection)
task test
# equivalent: go test -v -cover -race ./...

# Build both binaries to dist/
task

# Generate HTML coverage report
task coverage
```

### Building the Docker Image

```bash
docker build -t ghat .
```

The Dockerfile uses a two-stage build:
1. `golang:1.26-alpine` — compiles both binaries with `CGO_ENABLED=0` and compresses with `upx`
2. `alpine:3.23.3` — minimal runtime image

---

## Key Conventions

### Package Structure

The repository uses two package hierarchies with different access rules:

**`internal/`** — implementation packages, not importable by external modules:

| Package | Responsibility |
|---------|---------------|
| `internal/actions` | GitHub Actions workflow command helpers (masking, output, state, logging) |
| `internal/client` | GitHub API HTTP client |
| `internal/input` | Environment variable parsing into `Config` struct |
| `internal/jwt` | GitHub App JWT construction and signing |
| `internal/kms` | KMS signing abstraction (interface + real implementation) |

**`pkg/`** — public API, importable by external Go modules:

| Package | Responsibility |
|---------|---------------|
| `pkg/ghat` | Public API: KMS signer, JWT signing, token issuance, token revocation |

**`cmd/`** — entry points:

| Package | Responsibility |
|---------|---------------|
| `cmd/ghat` | Orchestrates config, KMS, JWT, GitHub API; outputs token |
| `cmd/post` | Revokes token post-workflow |

`cmd/` packages import from `internal/` only. `pkg/ghat` imports from `internal/` and wraps them in a clean public API.

### Public API (`pkg/ghat`)

External Go consumers import `github.com/yagihash/ghat/v2/pkg/ghat`:

```go
// 1. Create a KMS-backed signer
signer, err := ghat.NewSigner(ctx, projectID, location, keyRingID, keyID, version)
defer signer.Close()

// 2. Create an App
app := ghat.NewApp(appID, signer, "") // "" defaults to https://api.github.com

// 3. Sign a JWT (optional; IssueToken does this internally)
jwt, err := app.SignJWT(ctx)

// 4. Issue a GitHub App installation access token
token, err := app.IssueToken(ctx, owner, permissions, repositories)

// 5. Revoke the token when done
err = app.RevokeToken(ctx, token)
```

### Naming Conventions

- Package names: short, lowercase (`input`, `kms`, `client`, `actions`, `jwt`, `ghat`)
- Exported identifiers: `PascalCase` (`NewSigner`, `SetOutput`, `Config`)
- Unexported identifiers: `camelCase` (`newRequest`, `writeKeyValue`, `newApp`)
- Environment variable constants: `UPPER_SNAKE_CASE` (`EnvGitHubOutput`, `EnvGitHubState`)

### Error Handling

- Always return explicit `error` values; never swallow errors silently
- Wrap errors with context using `fmt.Errorf("context: %w", err)`
- In `main.go`, errors are printed to stderr via `actions.LogError` then `os.Exit(1)`

### Testing Patterns

- **Table-driven tests** for all non-trivial logic (see `internal/input/input_test.go`, `internal/kms/kms_test.go`)
- **Mock implementations** of interfaces for external dependencies (see `mockKMSClient` in `internal/kms/kms_test.go`)
- **Temporary directories** for any tests involving file I/O (see `internal/actions/actions_test.go`)
- **Environment variable setup/teardown** using `t.Setenv()` so tests don't leak state
- **`httptest.NewServer`** for GitHub API calls in `pkg/ghat/ghat_test.go`
- Tests live alongside their implementation files (`*_test.go`)
- `pkg/ghat` tests use `package ghat` (whitebox) to inject mock signers via the unexported `signerIface`

### Dependency Injection via Interfaces

External dependencies (KMS client, JWT signer) are abstracted behind interfaces to allow mocking in tests:

```go
// internal/kms/kms.go
type KMSClient interface {
    AsymmetricSign(ctx context.Context, ...) (*kmspb.AsymmetricSignResponse, error)
    Close() error
}

// internal/jwt/jwt.go
type Signer interface {
    Sign(ctx context.Context, data []byte) ([]byte, error)
}
```

Always define interface boundaries for any new external service integrations.

### Security Patterns

- **Token masking:** Always call `actions.AddMask(token)` before any token is logged or output
- **Token revocation:** The `post` binary deletes the token after use; preserve this behaviour
- **No private keys stored locally:** All signing is done via KMS
- **Minimal permissions:** Tokens are scoped to the minimal required repositories and permissions

---

## Configuration and Environment Variables

All configuration is read from environment variables with the `INPUT_` prefix using [`kelseyhightower/envconfig`](https://github.com/kelseyhightower/envconfig).

### Required Inputs

| Variable | Description |
|----------|-------------|
| `INPUT_APP_ID` | GitHub App ID |
| `INPUT_KMS_PROJECT_ID` | Google Cloud project ID |
| `INPUT_KMS_KEYRING_ID` | KMS key ring ID |
| `INPUT_KMS_KEY_ID` | KMS key ID |
| `INPUT_KMS_LOCATION` | KMS location (e.g., `asia-northeast1`) |

### Optional Inputs

| Variable | Default | Description |
|----------|---------|-------------|
| `INPUT_OWNER` | `GITHUB_REPOSITORY_OWNER` | GitHub org/user for app installation |
| `INPUT_REPOSITORIES` | (all) | Comma- or newline-separated repo list |
| `INPUT_KMS_KEY_VERSION` | `"1"` | KMS key version |
| `INPUT_BASE_URL` | `https://api.github.com` | GitHub API base URL (for GHES) |
| `INPUT_PERMISSION_*` | (unset) | Fine-grained permissions (60+ options) |

### GitHub Actions Runtime Variables

| Variable | Purpose |
|----------|---------|
| `GITHUB_ACTIONS` | Detected to switch between Actions mode and CLI mode |
| `GITHUB_REPOSITORY_OWNER` | Default owner for token installation |
| `GITHUB_OUTPUT` | File path used by `actions.SetOutput()` |
| `GITHUB_STATE` | File path used by `actions.SetState()` / `GetState()` |
| `STATE_TOKEN` | Token passed from main action to post-action for revocation |

---

## CI/CD Pipelines

All workflows live in `.github/workflows/`.

| Workflow | Trigger | Purpose |
|----------|---------|---------|
| `test.yml` | PR, push to `main` | Runs unit tests (`go test -v -cover -race ./...`) |
| `push-edge.yml` | Push to `main`, workflow dispatch | Builds and pushes `edge` Docker image; runs integration test |
| `release.yml` | Git tag `v*.*.*` | Builds versioned Docker image, updates `action.yml` digest, force-pushes tags |
| `ghalint.yml` | PR, push | Lints GitHub Actions workflow files |
| `pinact.yml` | PR, push | Validates action references are pinned to commit SHAs |
| `sync-permissions.yml` | PR, push | Ensures `action.yml` permission inputs stay in sync |

**Required secrets/variables for workflows:**

- Secrets: `BOT_GITHUB_APP_ID`, `KMS_KEY_ID`, `KMS_KEYRING_ID`, `WORKLOAD_IDENTITY_PROVIDER`, `SERVICE_ACCOUNT`
- Variables: `KMS_LOCATION`, `KMS_PROJECT_ID`

---

## action.yml Conventions

- All GitHub App permission inputs follow the pattern `permission-<resource>` (e.g., `permission-contents`, `permission-pull-requests`)
- Permissions are kept in sync with the GitHub API via `scripts/update-permissions.sh`
- The `token` output is the only output
- Dual entrypoints: `main: /ghat`, `post: /post`

---

## Adding New Features

### Adding a New Permission

Run `scripts/update-permissions.sh` to regenerate the permissions section in `action.yml` from the GitHub API rather than editing it by hand.

### Adding a New Input

1. Add the field to `input.Config` in `internal/input/input.go` with appropriate `envconfig` tags
2. Add a corresponding input definition to `action.yml`
3. Write table-driven tests in `internal/input/input_test.go`

### Adding a New GitHub API Call

1. Add the method to `internal/client/client.go`
2. Keep the HTTP client reusable (10s timeout, standard headers already configured)
3. Follow the existing pattern: build request → set headers → execute → decode JSON response
4. If the call should also be available in the public API, expose it via `pkg/ghat/ghat.go`

### Extending the Public API (`pkg/ghat`)

1. Add the new operation to `pkg/ghat/ghat.go` (or a new file in `pkg/ghat/`)
2. Use `signerIface` for any signer dependency injection (keeps `internal/kms` out of public signatures)
3. Write whitebox tests in `pkg/ghat/ghat_test.go` using `mockSigner` and `httptest.NewServer`
4. Do not expose `internal/` types in any exported function signature

---

## Dependencies

Key production dependencies (see `go.mod` for pinned versions):

| Dependency | Purpose |
|-----------|---------|
| `cloud.google.com/go/kms` | Google Cloud KMS API client |
| `github.com/googleapis/gax-go/v2` | Google API extensions (retry, call options) |
| `github.com/kelseyhightower/envconfig` | Environment variable parsing |
| `github.com/google/go-cmp` | Deep equality comparisons in tests |

Dependencies are managed via Renovate with semantic commits, scheduled weekly (Friday evenings JST), and GitHub Actions references are pinned to commit SHA digests for supply chain security.
