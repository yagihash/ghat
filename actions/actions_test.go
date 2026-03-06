package actions

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSetOutput(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "/github_output")

	t.Setenv(EnvGitHubOutput, path)

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

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Fatalf("unexpected file permission: %o", perm)
	}

	_ = os.Unsetenv(EnvGitHubOutput)
}

func TestSetOutput_NewlineInValue(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "/github_output")

	if err := os.Setenv(EnvGitHubOutput, path); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Unsetenv(EnvGitHubOutput) }()

	err := SetOutput("key", "line1\nline2")
	if err == nil {
		t.Fatal("expected error for newline in value, got nil")
	}
}

func TestSetState(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "/github_state")

	t.Setenv(EnvGitHubState, path)

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
}

func TestGetState(t *testing.T) {
	t.Setenv("STATE_KEY", "value")

	val, err := GetState("key")
	if err != nil {
		t.Fatal(err)
	}

	if val != "value" {
		t.Fatalf("unexpected value: %s", val)
	}
}
