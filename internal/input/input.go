package input

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	AppID            string            `envconfig:"APP_ID" required:"true"`
	Owner            string            `envconfig:"OWNER"`
	Repositories     Repositories      `envconfig:"REPOSITORIES"`
	PermissionInputs Permissions       `envconfig:"PERMISSION"`
	Permissions      map[string]string `ignored:"true"`
	BaseURL          string            `envconfig:"BASE_URL" default:"https://api.github.com"`

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

	c.Permissions = c.PermissionInputs.toMap()

	return &c, nil
}

// toMap flattens the envconfig-bound per-resource fields into the
// map[string]string shape the GitHub API expects, keyed by each field's perm
// tag (its exact API resource name). Fields envconfig left unset stay "" and
// are excluded.
func (p Permissions) toMap() map[string]string {
	v := reflect.ValueOf(p)

	permissions := make(map[string]string, v.NumField())
	for f := range reflect.TypeFor[Permissions]().Fields() {
		value := v.FieldByName(f.Name).String()
		if value == "" {
			continue
		}
		permissions[f.Tag.Get("perm")] = value
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
