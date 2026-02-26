package bash

// Package bash provides bash command execution.

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/domain"
	"github.com/mascli/troncli/internal/policy"
)

// UniversalBashManager implements BashManager for all Linux distributions
type UniversalBashManager struct {
	executor     adapter.Executor
	profile      *domain.SystemProfile
	policyEngine *policy.PolicyEngine
}

// NewUniversalBashManager creates a new instance of UniversalBashManager
func NewUniversalBashManager(executor adapter.Executor, profile *domain.SystemProfile, policyEngine *policy.PolicyEngine) *UniversalBashManager {
	return &UniversalBashManager{
		executor:     executor,
		profile:      profile,
		policyEngine: policyEngine,
	}
}

// RunCommand executes a single bash command
func (m *UniversalBashManager) RunCommand(cmd string) (string, error) {
	// Check command against policy engine
	if err := m.policyEngine.CheckCommand(cmd); err != nil {
		return "", fmt.Errorf("policy engine blocked execution: %w", err)
	}

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
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to resolve path: %w", err)
	}

	if _, errStat := os.Stat(absPath); os.IsNotExist(errStat) {
		return "", fmt.Errorf("script file not found: %s", absPath)
	}

	// Read script content for policy checking
	scriptContent, err := os.ReadFile(absPath)
	if err != nil {
		return "", fmt.Errorf("failed to read script file: %w", err)
	}

	// Check script content against policy engine
	scriptStr := string(scriptContent)
	if err := m.policyEngine.CheckCommand(scriptStr); err != nil {
		return "", fmt.Errorf("policy engine blocked script execution: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute) // Scripts might take longer
	defer cancel()

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
