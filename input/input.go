package input

import (
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

	ProjectID  string `envconfig:"KMS_PROJECT_ID" required:"true"`
	KeyRingID  string `envconfig:"KMS_KEYRING_ID" required:"true"`
	KeyID      string `envconfig:"KMS_KEY_ID" required:"true"`
	KeyVersion string `envconfig:"KMS_KEY_VERSION"`
	Location   string `envconfig:"KMS_LOCATION" required:"true"`
}

func Load() (*Config, error) {
	var c Config
	if err := envconfig.Process("INPUT", &c); err != nil {
		return nil, err
	}

	if c.Owner == "" {
		c.Owner = os.Getenv("GITHUB_REPOSITORY_OWNER")
	}

	if c.KeyVersion == "" {
		c.KeyVersion = "1"
	}

	if len(c.Permissions) > 0 {
		lowered := make(map[string]string, len(c.Permissions))
		for k, v := range c.Permissions {
			lowered[strings.ToLower(k)] = v
		}
		c.Permissions = lowered
	}

	return &c, nil
}

type Repositories []string

func (r *Repositories) Decode(value string) error {
	if value == "" {
		return nil
	}

	res := make(Repositories, 0)
	normalized := strings.ReplaceAll(value, "\n", ",")
	repos := strings.Split(normalized, ",")
	for _, repo := range repos {
		trimmed := strings.TrimSpace(repo)
		if trimmed != "" {
			res = append(res, trimmed)
		}
	}

	*r = res

	return nil
}
