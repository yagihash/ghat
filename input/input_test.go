package input

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLoad(t *testing.T) {
	_, err := Load()
	if err != nil {
		t.Fatal(err)
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
