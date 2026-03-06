package actions

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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

func TestSetOutput_EnvNotSet(t *testing.T) {
	t.Setenv(EnvGitHubOutput, "")

	err := SetOutput("key", "value")
	if err == nil {
		t.Fatal("expected error when GITHUB_OUTPUT is not set, got nil")
	}
}

func TestSetOutput_Append(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "github_output")
	t.Setenv(EnvGitHubOutput, path)

	if err := SetOutput("key1", "value1"); err != nil {
		t.Fatal(err)
	}
	if err := SetOutput("key2", "value2"); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	want := "key1=value1\nkey2=value2\n"
	if string(content) != want {
		t.Errorf("content = %q, want %q", string(content), want)
	}
}

func TestSetState_EnvNotSet(t *testing.T) {
	t.Setenv(EnvGitHubState, "")

	err := SetState("key", "value")
	if err == nil {
		t.Fatal("expected error when GITHUB_STATE is not set, got nil")
	}
}

func TestSetState_EmptyKey(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv(EnvGitHubState, filepath.Join(tempDir, "github_state"))

	err := SetState("", "value")
	if err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
}

func TestSetState_EmptyValue(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv(EnvGitHubState, filepath.Join(tempDir, "github_state"))

	err := SetState("key", "")
	if err == nil {
		t.Fatal("expected error for empty value, got nil")
	}
}

func TestSetState_NewlineInValue(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv(EnvGitHubState, filepath.Join(tempDir, "github_state"))

	err := SetState("key", "line1\nline2")
	if err == nil {
		t.Fatal("expected error for newline in value, got nil")
	}
}

func TestSetState_KeyNormalization(t *testing.T) {
	tests := []struct {
		key     string
		wantKey string
	}{
		{"my-key", "MY_KEY"},
		{"MY-KEY", "MY_KEY"},
		{"token", "TOKEN"},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			tempDir := t.TempDir()
			path := filepath.Join(tempDir, "github_state")
			t.Setenv(EnvGitHubState, path)

			if err := SetState(tt.key, "value"); err != nil {
				t.Fatal(err)
			}

			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			want := fmt.Sprintf("%s=value\n", tt.wantKey)
			if string(content) != want {
				t.Errorf("content = %q, want %q", string(content), want)
			}
		})
	}
}

func TestGetState_NotFound(t *testing.T) {
	_, err := GetState("nonexistent-key")
	if err == nil {
		t.Fatal("expected error for missing state key, got nil")
	}
}

func TestGetState_KeyNormalization(t *testing.T) {
	t.Setenv("STATE_MY_KEY", "hello")

	val, err := GetState("my-key")
	if err != nil {
		t.Fatal(err)
	}
	if val != "hello" {
		t.Errorf("val = %q, want %q", val, "hello")
	}
}

func captureStdout(t *testing.T, f func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	orig := os.Stdout
	os.Stdout = w
	f()
	w.Close()
	os.Stdout = orig
	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	return string(out)
}

func TestLogDebug(t *testing.T) {
	out := captureStdout(t, func() { LogDebug("debug message") })
	if want := "::debug::debug message\n"; out != want {
		t.Errorf("output = %q, want %q", out, want)
	}
}

func TestLogNotice(t *testing.T) {
	out := captureStdout(t, func() { LogNotice("notice message") })
	if want := "::notice::notice message\n"; out != want {
		t.Errorf("output = %q, want %q", out, want)
	}
}

func TestLogWarning(t *testing.T) {
	out := captureStdout(t, func() { LogWarning("warning message") })
	if want := "::warning::warning message\n"; out != want {
		t.Errorf("output = %q, want %q", out, want)
	}
}

func TestLogError(t *testing.T) {
	out := captureStdout(t, func() { LogError("error message") })
	if want := "::error::error message\n"; out != want {
		t.Errorf("output = %q, want %q", out, want)
	}
}

func TestAddMask(t *testing.T) {
	out := captureStdout(t, func() { AddMask("secret-token") })
	if want := "::add-mask::secret-token\n"; out != want {
		t.Errorf("output = %q, want %q", out, want)
	}
}

func TestLogGroup(t *testing.T) {
	out := captureStdout(t, func() { LogGroup("my group", "line1", "line2") })
	if !strings.HasPrefix(out, "::group::my group") {
		t.Errorf("output should start with ::group::my group, got %q", out)
	}
	if !strings.Contains(out, "line1") || !strings.Contains(out, "line2") {
		t.Errorf("output should contain messages, got %q", out)
	}
	if !strings.Contains(out, "::endgroup::") {
		t.Errorf("output should contain ::endgroup::, got %q", out)
	}
}

func TestWorkflowCommand_WithParams(t *testing.T) {
	out := captureStdout(t, func() {
		workflowCommand("error", "something went wrong", map[string]string{"file": "main.go", "line": "42"})
	})
	if !strings.HasPrefix(out, "::error ") {
		t.Errorf("output should start with ::error , got %q", out)
	}
	if !strings.Contains(out, "something went wrong") {
		t.Errorf("output should contain value, got %q", out)
	}
	if !strings.Contains(out, "file=main.go") || !strings.Contains(out, "line=42") {
		t.Errorf("output should contain params, got %q", out)
	}
}
