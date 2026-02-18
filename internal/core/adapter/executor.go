package adapter

// Package adapter defines interfaces for external system adapters.

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// CommandResult represents the result of an execution
type CommandResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Duration time.Duration
}

// Executor defines the interface for executing commands
type Executor interface {
	Exec(ctx context.Context, command string, args ...string) (*CommandResult, error)
	ExecWithInput(ctx context.Context, input string, command string, args ...string) (*CommandResult, error)
}

// SystemExecutor implements Executor with safe execution patterns
type SystemExecutor struct {
	// mu sync.Mutex // Optional: use for concurrency limiting if needed globally
}

// NewExecutor creates a new executor
func NewExecutor() Executor {
	return &SystemExecutor{}
}

// Exec executes a command with context
func (e *SystemExecutor) Exec(ctx context.Context, command string, args ...string) (*CommandResult, error) {
	return e.ExecWithInput(ctx, "", command, args...)
}

// ExecWithInput executes a command with input and context
func (e *SystemExecutor) ExecWithInput(ctx context.Context, input string, command string, args ...string) (*CommandResult, error) {
	// Sanitization could happen here, but usually args are safe with exec.Command
	// We might want to log execution here for audit

	start := time.Now()

	cmd := exec.CommandContext(ctx, command, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if input != "" {
		cmd.Stdin = strings.NewReader(input)
	}

	err := cmd.Run()
	duration := time.Since(start)

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = -1 // Internal error or signal
		}
	}

	result := &CommandResult{
		Stdout:   strings.TrimSpace(stdout.String()),
		Stderr:   strings.TrimSpace(stderr.String()),
		ExitCode: exitCode,
		Duration: duration,
	}

	// Wrap error if exit code is non-zero and err is nil (shouldn't happen with Run)
	// Or if err is not nil, return it but also return result for inspection
	if err != nil {
		return result, fmt.Errorf("command execution failed: %w (stderr: %s)", err, result.Stderr)
	}

	return result, nil
}
