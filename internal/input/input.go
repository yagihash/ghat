package input

import (
	"fmt"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	AppID        string            `envconfig:"APP_ID" required:"true"`
	Owner        string            `envconfig:"OWNER"`
	Repositories Repositories      `envconfig:"REPOSITORIES"`
	Permissions  map[string]string `envconfig:"PERMISSION"`
	BaseURL      string            `envconfig:"BASE_URL" default:"https://api.github.com"`

	ProjectID string `envconfig:"KMS_PROJECT_ID" required:"true"`
	KeyRingID string `envconfig:"KMS_KEYRING_ID" required:"true"`
	KeyID     string `envconfig:"KMS_KEY_ID" required:"true"`
	// KeyVersion defaults to "" to signal automatic detection of the latest enabled
	// KMS key version. GitHub Actions passes an explicit empty string for optional
	// inputs the user leaves unset, so the zero value carries the same semantic.
	KeyVersion string `envconfig:"KMS_KEY_VERSION" default:""`
	Location   string `envconfig:"KMS_LOCATION" required:"true"`
}

func Load() (*Config, error) {
	var c Config
	if err := envconfig.Process("INPUT", &c); err != nil {
		return nil, err
	}

	// The JWT and the issued token are sent in the Authorization header, so
	// plaintext transport must never be allowed.
	if !strings.HasPrefix(c.BaseURL, "https://") {
		return nil, fmt.Errorf("base_url must start with https://: %s", c.BaseURL)
	}

	if c.Owner == "" {
		c.Owner = os.Getenv("GITHUB_REPOSITORY_OWNER")
	}

	c.Permissions = loadPermissions(c.Permissions)

	return &c, nil
}

const permissionEnvPrefix = "INPUT_PERMISSION_"

// loadPermissions merges the aggregate INPUT_PERMISSION map (if any) with the
// ~60 individual INPUT_PERMISSION_<NAME> variables that action.yml's
// permission_<resource> inputs actually produce. envconfig only binds the
// single aggregate variable, so the per-resource inputs must be collected
// manually here.
func loadPermissions(aggregate map[string]string) map[string]string {
	permissions := make(map[string]string, len(aggregate))
	for k, v := range aggregate {
		permissions[strings.ToLower(k)] = v
	}

	for _, kv := range os.Environ() {
		key, value, ok := strings.Cut(kv, "=")
		if !ok || !strings.HasPrefix(key, permissionEnvPrefix) || value == "" {
			continue
		}
		name := strings.ToLower(strings.TrimPrefix(key, permissionEnvPrefix))
		permissions[name] = value
	}

	if len(permissions) == 0 {
		return nil
	}

	return permissions
}

type Repositories []string

func (r *Repositories) Decode(value string) error {
	if value == "" {
		return nil
	}

	res := make(Repositories, 0)
	normalized := strings.ReplaceAll(value, "\n", ",")
	repos := strings.SplitSeq(normalized, ",")
	for repo := range repos {
		trimmed := strings.TrimSpace(repo)
		if trimmed != "" {
			res = append(res, trimmed)
		}
	}

	*r = res

	return nil
}
