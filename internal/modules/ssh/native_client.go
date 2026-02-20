package ssh

// Package ssh provides SSH configuration parsing and connection management.

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kevinburke/ssh_config"
	"github.com/mascli/troncli/internal/core/ports"
)

// NativeSSHClient implements ports.SSHClient by parsing ~/.ssh/config directly
type NativeSSHClient struct {
	configPath string
}

// NewNativeSSHClient creates a new client that reads the user's SSH config
func NewNativeSSHClient() (ports.SSHClient, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	configPath := filepath.Join(home, ".ssh", "config")
	return &NativeSSHClient{
		configPath: configPath,
	}, nil
}

// ListProfiles returns a list of available SSH hosts from ~/.ssh/config
func (c *NativeSSHClient) ListProfiles() ([]string, error) {
	f, err := os.Open(c.configPath)
	if os.IsNotExist(err) {
		// Return empty list if no config exists, do not return error as per user request (no fake data, just empty)
		return []string{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to open ssh config: %w", err)
	}
	defer f.Close()

	cfg, err := ssh_config.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ssh config: %w", err)
	}

	var profiles []string
	seen := make(map[string]bool)

	for _, host := range cfg.Hosts {
		for _, pattern := range host.Patterns {
			name := pattern.String()
			// Filter out wildcards and duplicates
			if name != "*" && !strings.Contains(name, "*") && !strings.Contains(name, "?") && !seen[name] {
				profiles = append(profiles, name)
				seen[name] = true
			}
		}
	}

	sort.Strings(profiles)
	return profiles, nil
}

// Connect establishes an interactive SSH connection using the system's ssh client
func (c *NativeSSHClient) Connect(profile string) error {
	// We use the system ssh client for the best interactive experience (terminal emulation, keys, etc.)
	// This replaces the "rsd-sshm connect" wrapper.
	sshPath, err := exec.LookPath("ssh")
	if err != nil {
		return fmt.Errorf("ssh binary not found in PATH")
	}

	cmd := exec.Command(sshPath, profile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// Execute runs a command on the remote host and returns the output
func (c *NativeSSHClient) Execute(profile string, command string) (string, error) {
	sshPath, err := exec.LookPath("ssh")
	if err != nil {
		return "", fmt.Errorf("ssh binary not found in PATH")
	}

	cmd := exec.Command(sshPath, profile, command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ssh execution failed: %w, output: %s", err, string(output))
	}

	return string(output), nil
}

// CreateTunnel establishes an SSH tunnel
func (c *NativeSSHClient) CreateTunnel(profile string, localPort, remoteHost, remotePort string, reverse bool) error {
	sshPath, err := exec.LookPath("ssh")
	if err != nil {
		return fmt.Errorf("ssh binary not found in PATH")
	}

	// Construct tunnel argument
	// Local: -L localPort:remoteHost:remotePort
	// Remote: -R remotePort:localHost:localPort (simplified mapping)

	tunnelArg := fmt.Sprintf("%s:%s:%s", localPort, remoteHost, remotePort)
	flag := "-L"
	if reverse {
		flag = "-R"
	}

	// -N: Do not execute a remote command.
	// -f: Requests ssh to go to background just before command execution.
	cmd := exec.Command(sshPath, "-N", flag, tunnelArg, profile)
	
	// Since we want to manage the process, we don't use -f which forks.
	// We start it and keep the process handle.
	// Note: This simple implementation doesn't track tunnels persistently across app restarts yet,
	// but it replaces the rsd-sshm wrapper logic.
	
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start tunnel: %w", err)
	}

	// In a real implementation we would store cmd.Process to kill it later.
	// For now, we follow the interface.
	return nil
}

// CloseTunnel closes an active tunnel
func (c *NativeSSHClient) CloseTunnel(profile string, localPort string) error {
	// TODO: Implement tunnel tracking in NativeSSHClient struct if needed
	return fmt.Errorf("tunnel management not fully implemented in native client yet")
}

// CopyFile copies a file to the remote host using scp
func (c *NativeSSHClient) CopyFile(profile string, src string, dest string) error {
	scpPath, err := exec.LookPath("scp")
	if err != nil {
		return fmt.Errorf("scp binary not found in PATH")
	}

	cmd := exec.Command(scpPath, src, fmt.Sprintf("%s:%s", profile, dest))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("scp failed: %w, output: %s", err, string(output))
	}
	return nil
}
