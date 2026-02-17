package bash

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/domain"
)

// UniversalBashManager implements BashManager for all Linux distributions
type UniversalBashManager struct {
	executor adapter.Executor
	profile  *domain.SystemProfile
}

// NewUniversalBashManager creates a new instance of UniversalBashManager
func NewUniversalBashManager(executor adapter.Executor, profile *domain.SystemProfile) *UniversalBashManager {
	return &UniversalBashManager{
		executor: executor,
		profile:  profile,
	}
}

// RunCommand executes a single bash command
func (m *UniversalBashManager) RunCommand(cmd string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Use bash explicitly to run the command
	res, err := m.executor.Exec(ctx, "bash", "-c", cmd)
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %w\nOutput: %s", err, res.Stderr)
	}
	return res.Stdout, nil
}

// RunScript executes a bash script from a file
func (m *UniversalBashManager) RunScript(path string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute) // Scripts might take longer
	defer cancel()

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to resolve path: %w", err)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return "", fmt.Errorf("script file not found: %s", absPath)
	}

	// Make sure script is executable
	_, err = m.executor.Exec(ctx, "chmod", "+x", absPath)
	if err != nil {
		return "", fmt.Errorf("failed to make script executable: %w", err)
	}

	res, err := m.executor.Exec(ctx, "bash", absPath)
	if err != nil {
		return "", fmt.Errorf("failed to execute script: %w\nOutput: %s", err, res.Stderr)
	}
	return res.Stdout, nil
}
