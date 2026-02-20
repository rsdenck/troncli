package remote

// Package remote provides remote execution capabilities.

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/domain"
)

type UniversalRemoteManager struct {
	executor adapter.Executor
	profile  *domain.SystemProfile
}

func NewUniversalRemoteManager(executor adapter.Executor, profile *domain.SystemProfile) *UniversalRemoteManager {
	return &UniversalRemoteManager{
		executor: executor,
		profile:  profile,
	}
}

func (m *UniversalRemoteManager) ListProfiles() ([]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	configPath := filepath.Join(home, ".ssh", "config")

	file, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var profiles []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "Host ") {
			// Handle "Host alias1 alias2"
			hosts := strings.Fields(line)[1:]
			for _, h := range hosts {
				if h != "*" && !strings.Contains(h, "?") {
					profiles = append(profiles, h)
				}
			}
		}
	}
	return profiles, nil
}

func (m *UniversalRemoteManager) Connect(profile string) error {
	// Connect is interactive, so we cannot use m.executor.Exec which captures output.
	// We need to attach Stdin/Stdout/Stderr directly.
	cmd := exec.Command("ssh", profile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (m *UniversalRemoteManager) Execute(profile string, command string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	res, err := m.executor.Exec(ctx, "ssh", profile, command)
	if err != nil {
		return "", fmt.Errorf("ssh execution failed: %w\nOutput: %s", err, res.Stderr)
	}
	return res.Stdout, nil
}

func (m *UniversalRemoteManager) CopyFile(profile string, src string, dest string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	// Using scp
	// If dest is absolute, fine. If relative, it's relative to home on remote.
	// Format: scp src profile:dest
	remoteDest := fmt.Sprintf("%s:%s", profile, dest)
	res, err := m.executor.Exec(ctx, "scp", src, remoteDest)
	if err != nil {
		return fmt.Errorf("scp failed: %w\nOutput: %s", err, res.Stderr)
	}
	return nil
}

func (m *UniversalRemoteManager) CreateTunnel(profile string, localPort, remoteHost, remotePort string, reverse bool) error {
	// Not implemented for now or requires background process management
	return fmt.Errorf("tunnel creation not fully implemented in CLI mode yet")
}

func (m *UniversalRemoteManager) CloseTunnel(profile string, localPort string) error {
	return fmt.Errorf("tunnel closing not implemented")
}
