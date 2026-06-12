package input

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLoad(t *testing.T) {
	t.Setenv("INPUT_APP_ID", "12345")
	t.Setenv("INPUT_OWNER", "owner")
	t.Setenv("INPUT_REPOSITORIES", "owner/repo1,owner/repo2")
	t.Setenv("INPUT_PERMISSION_CONTENTS", "write")
	t.Setenv("INPUT_PERMISSION_ISSUES", "read")
	t.Setenv("INPUT_BASE_URL", "https://api.github.com")
	t.Setenv("INPUT_KMS_PROJECT_ID", "project-id")
	t.Setenv("INPUT_KMS_KEYRING_ID", "keyring-id")
	t.Setenv("INPUT_KMS_KEY_ID", "key-id")
	t.Setenv("INPUT_KMS_KEY_VERSION", "1")
	t.Setenv("INPUT_KMS_LOCATION", "us-central1")

	_, err := Load()
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoad_DefaultKeyVersion(t *testing.T) {
	t.Setenv("INPUT_APP_ID", "12345")
	t.Setenv("INPUT_OWNER", "owner")
	t.Setenv("INPUT_REPOSITORIES", "owner/repo1,owner/repo2")
	t.Setenv("INPUT_PERMISSION_CONTENTS", "write")
	t.Setenv("INPUT_PERMISSION_ISSUES", "read")
	t.Setenv("INPUT_BASE_URL", "https://api.github.com")
	t.Setenv("INPUT_KMS_PROJECT_ID", "project-id")
	t.Setenv("INPUT_KMS_KEYRING_ID", "keyring-id")
	t.Setenv("INPUT_KMS_KEY_ID", "key-id")
	t.Setenv("INPUT_KMS_LOCATION", "us-central1")

	i, err := Load()
	if err != nil {
		t.Fatal(err)
	}

	if i.KeyVersion != "" {
		t.Errorf("KeyVersion should be empty string when not set (auto-detect): %s", i.KeyVersion)
	}
}

// TestLoad_EmptyKeyVersion verifies that KeyVersion remains empty when
// INPUT_KMS_KEY_VERSION is explicitly set to an empty string. This is a real
// scenario in GitHub Actions: optional inputs with no default in action.yml
// are passed as empty strings, not as unset variables. An empty KeyVersion
// signals to the KMS layer that the latest enabled version should be detected
// automatically.
func TestLoad_EmptyKeyVersion(t *testing.T) {
	t.Setenv("INPUT_APP_ID", "12345")
	t.Setenv("INPUT_OWNER", "owner")
	t.Setenv("INPUT_REPOSITORIES", "owner/repo1,owner/repo2")
	t.Setenv("INPUT_PERMISSION_CONTENTS", "write")
	t.Setenv("INPUT_PERMISSION_ISSUES", "read")
	t.Setenv("INPUT_BASE_URL", "https://api.github.com")
	t.Setenv("INPUT_KMS_PROJECT_ID", "project-id")
	t.Setenv("INPUT_KMS_KEYRING_ID", "keyring-id")
	t.Setenv("INPUT_KMS_KEY_ID", "key-id")
	t.Setenv("INPUT_KMS_KEY_VERSION", "")
	t.Setenv("INPUT_KMS_LOCATION", "us-central1")

	i, err := Load()
	if err != nil {
		t.Fatal(err)
	}

	if i.KeyVersion != "" {
		t.Errorf("KeyVersion should remain empty string when INPUT_KMS_KEY_VERSION is empty string, got: %q", i.KeyVersion)
	}
}

func TestLoad_NonHTTPSBaseURL(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
		wantErr bool
	}{
		{
			name:    "https URL is accepted",
			baseURL: "https://github.example.com/api/v3",
			wantErr: false,
		},
		{
			name:    "http URL is rejected",
			baseURL: "http://api.github.com",
			wantErr: true,
		},
		{
			name:    "scheme-less URL is rejected",
			baseURL: "api.github.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("INPUT_APP_ID", "12345")
			t.Setenv("INPUT_BASE_URL", tt.baseURL)
			t.Setenv("INPUT_KMS_PROJECT_ID", "project-id")
			t.Setenv("INPUT_KMS_KEYRING_ID", "keyring-id")
			t.Setenv("INPUT_KMS_KEY_ID", "key-id")
			t.Setenv("INPUT_KMS_LOCATION", "us-central1")

			_, err := Load()
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
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
