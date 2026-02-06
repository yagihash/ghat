package input

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLoad(t *testing.T) {
	os.Setenv("INPUT_APP_ID", "12345")
	os.Setenv("INPUT_OWNER", "owner")
	os.Setenv("INPUT_REPOSITORIES", "owner/repo1,owner/repo2")
	os.Setenv("INPUT_PERMISSION_CONTENTS", "write")
	os.Setenv("INPUT_PERMISSION_ISSUES", "read")
	os.Setenv("INPUT_BASE_URL", "https://api.github.com")
	os.Setenv("INPUT_KMS_PROJECT_ID", "project-id")
	os.Setenv("INPUT_KMS_KEYRING_ID", "keyring-id")
	os.Setenv("INPUT_KMS_KEY_ID", "key-id")
	os.Setenv("INPUT_KMS_KEY_VERSION", "1")
	os.Setenv("INPUT_KMS_LOCATION", "us-central1")

	_, err := Load()
	if err != nil {
		t.Fatal(err)
	}

	os.Unsetenv("INPUT_APP_ID")
	os.Unsetenv("INPUT_OWNER")
	os.Unsetenv("INPUT_REPOSITORIES")
	os.Unsetenv("INPUT_PERMISSION_CONTENTS")
	os.Unsetenv("INPUT_PERMISSION_ISSUES")
	os.Unsetenv("INPUT_BASE_URL")
	os.Unsetenv("INPUT_KMS_PROJECT_ID")
	os.Unsetenv("INPUT_KMS_KEYRING_ID")
	os.Unsetenv("INPUT_KMS_KEY_ID")
	os.Unsetenv("INPUT_KMS_KEY_VERSION")
	os.Unsetenv("INPUT_KMS_LOCATION")
}

func TestLoad_DefaultKeyVersion(t *testing.T) {
	os.Setenv("INPUT_APP_ID", "12345")
	os.Setenv("INPUT_OWNER", "owner")
	os.Setenv("INPUT_REPOSITORIES", "owner/repo1,owner/repo2")
	os.Setenv("INPUT_PERMISSION_CONTENTS", "write")
	os.Setenv("INPUT_PERMISSION_ISSUES", "read")
	os.Setenv("INPUT_BASE_URL", "https://api.github.com")
	os.Setenv("INPUT_KMS_PROJECT_ID", "project-id")
	os.Setenv("INPUT_KMS_KEYRING_ID", "keyring-id")
	os.Setenv("INPUT_KMS_KEY_ID", "key-id")
	os.Setenv("INPUT_KMS_LOCATION", "us-central1")

	i, err := Load()
	if err != nil {
		t.Fatal(err)
	}

	if i.KeyVersion != "1" {
		t.Errorf("KeyVersion should be set to default value: %s", i.KeyVersion)
	}

	os.Unsetenv("INPUT_APP_ID")
	os.Unsetenv("INPUT_OWNER")
	os.Unsetenv("INPUT_REPOSITORIES")
	os.Unsetenv("INPUT_PERMISSION_CONTENTS")
	os.Unsetenv("INPUT_PERMISSION_ISSUES")
	os.Unsetenv("INPUT_BASE_URL")
	os.Unsetenv("INPUT_KMS_PROJECT_ID")
	os.Unsetenv("INPUT_KMS_KEYRING_ID")
	os.Unsetenv("INPUT_KMS_KEY_ID")
	os.Unsetenv("INPUT_KMS_LOCATION")
}

func TestRepositories_Decode(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Repositories
		wantErr bool
	}{
		{
			name:    "Comma-separated string",
			input:   "owner/repo1,owner/repo2",
			want:    Repositories{"owner/repo1", "owner/repo2"},
			wantErr: false,
		},
		{
			name:    "New-line separated string",
			input:   "owner/repo1\nowner/repo2",
			want:    Repositories{"owner/repo1", "owner/repo2"},
			wantErr: false,
		},
		{
			name:    "Mixed string",
			input:   "owner/repo1, \n owner/repo2 \n,owner/repo3",
			want:    Repositories{"owner/repo1", "owner/repo2", "owner/repo3"},
			wantErr: false,
		},
		{
			name:    "Blank string",
			input:   "",
			want:    nil,
			wantErr: false,
		},
		{
			name:    "Extra commas and new lines",
			input:   ",,,repo1\n\nrepo2,",
			want:    Repositories{"repo1", "repo2"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r Repositories
			err := r.Decode(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.want, r); diff != "" {
				t.Errorf("Decode() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
