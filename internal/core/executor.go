package core

import (
	"fmt"
	"os/exec"
	"strings"
)

// Executor interface allows mocking for tests
type Executor interface {
	Run(name string, args ...string) (string, error)
	RunSilent(name string, args ...string) error
	CombinedOutput(name string, args ...string) (string, error)
}

// RealExecutor executes real system commands
type RealExecutor struct{}

func (r *RealExecutor) Run(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("command %s failed: %w", name, err)
	}
	return strings.TrimSpace(string(output)), nil
}

func (r *RealExecutor) RunSilent(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	return cmd.Run()
}

func (r *RealExecutor) CombinedOutput(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command %s failed: %w - %s", name, err, string(output))
	}
	return strings.TrimSpace(string(output)), nil
}

// MockExecutor for testing
type MockExecutor struct {
	Output string
	Err    error
}

func (m *MockExecutor) Run(name string, args ...string) (string, error) {
	return m.Output, m.Err
}

func (m *MockExecutor) RunSilent(name string, args ...string) error {
	return m.Err
}

func (m *MockExecutor) CombinedOutput(name string, args ...string) (string, error) {
	return m.Output, m.Err
}

// SanitizeInput prevents shell injection attacks
func SanitizeInput(input string) string {
	dangerous := []string{";", "&&", "||", "`", "$(", "${", "|", ">", "<", "\n", "\r"}
	result := input
	for _, d := range dangerous {
		result = strings.ReplaceAll(result, d, "")
	}
	return result
}

// ValidatePath ensures path is safe
func ValidatePath(path string) bool {
	// Must be absolute path
	if !strings.HasPrefix(path, "/") {
		return false
	}
	// Check for path traversal
	if strings.Contains(path, "..") {
		return false
	}
	return true
}
