package actions

import (
	"fmt"
	"os"
	"strings"
)

const (
	EnvGitHubOutput = "GITHUB_OUTPUT"
	EnvGitHubState  = "GITHUB_STATE"
)

func SetOutput(key, value string) error {
	outputFilePath := os.Getenv("GITHUB_OUTPUT")
	if outputFilePath == "" {
		return fmt.Errorf("GITHUB_OUTPUT environment variable is not set")
	}

	if err := writeKeyValue(outputFilePath, key, value); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	return nil
}

func SetState(key, value string) error {
	stateFilePath := os.Getenv("GITHUB_STATE")
	if stateFilePath == "" {
		return fmt.Errorf("GITHUB_STATE environment variable is not set")
	}

	if key == "" {
		return fmt.Errorf("key is empty")
	}

	if value == "" {
		return fmt.Errorf("value is empty")
	}

	if strings.Contains(value, "\n") {
		return fmt.Errorf("SetState does not support new-line characters in the value: %s", value)
	}

	if err := writeKeyValue(stateFilePath, strings.ToUpper(key), value); err != nil {
		return fmt.Errorf("failed to write state: %w", err)
	}

	return nil
}

func GetState(key string) (string, error) {
	formatted := "STATE_" + normalizeToEnvKey(key)
	value := os.Getenv(formatted)
	if value == "" {
		return "", fmt.Errorf("state %s is not found", formatted)
	}

	return value, nil
}

func LogGroup(title string, messages ...string) {
	fmt.Print("::group::" + title)
	for _, v := range messages {
		fmt.Println(v)
	}
	fmt.Println("::endgroup::")
}

func LogDebug(value string) {
	workflowCommand("debug", value, nil)
}

func LogNotice(value string) {
	workflowCommand("notice", value, nil)
}

func LogWarning(value string) {
	workflowCommand("warning", value, nil)
}

func LogError(value string) {
	workflowCommand("error", value, nil)
}

func AddMask(value string) {
	workflowCommand("add-mask", value, nil)
}

func workflowCommand(command string, value string, params map[string]string) {
	if params == nil || len(params) == 0 {
		fmt.Printf("::%s::%s\n", command, value)
	} else {
		p := make([]string, 0, len(params))
		for k, v := range params {
			p = append(p, fmt.Sprintf("%s=%s", k, v))
		}
		fmt.Printf("::%s %s::%s\n", command, strings.Join(p, ","), value)
	}
}

func normalizeToEnvKey(key string) string {
	return strings.ToUpper(strings.ReplaceAll(key, "-", "_"))
}

func writeKeyValue(path, key, value string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", path, err)
	}

	if _, err := fmt.Fprintf(f, "%s=%s\n", key, value); err != nil {
		_ = f.Close()
		return fmt.Errorf("failed to write value: %w", err)
	}

	if err := f.Sync(); err != nil {
		_ = f.Close()
		return fmt.Errorf("failed to sync: %w", err)
	}

	return f.Close()
}
