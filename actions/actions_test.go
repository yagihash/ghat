package actions

import (
	"os"
	"path/filepath"
	"testing"
)

const ()

func TestSetOutput(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "/github_output")

	if err := os.Setenv(EnvGitHubOutput, path); err != nil {
		t.Fatal(err)
	}

	if err := SetOutput("key", "value"); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("file %s does not exist", path)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != "key=value\n" {
		t.Fatalf("unexpected content: %s", content)
	}

	_ = os.Unsetenv(EnvGitHubOutput)
}

func TestSetState(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "/github_state")

	if err := os.Setenv(EnvGitHubState, path); err != nil {
		t.Fatal(err)
	}

	if err := SetState("key", "value"); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("file %s does not exist", path)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != "KEY=value\n" {
		t.Fatalf("unexpected content: %s", content)
	}

	_ = os.Unsetenv(EnvGitHubState)
}

func TestGetState(t *testing.T) {
	if err := os.Setenv("STATE_KEY", "value"); err != nil {
		t.Fatal(err)
	}

	val, err := GetState("key")
	if err != nil {
		t.Fatal(err)
	}

	if val != "value" {
		t.Fatalf("unexpected value: %s", val)
	}

	_ = os.Unsetenv("STATE_KEY")
}
