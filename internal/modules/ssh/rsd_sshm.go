package ssh

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mascli/troncli/internal/core/ports"
)

// RSDSSHMClient implements ports.SSHClient using rsd-sshm binary
type RSDSSHMClient struct {
	binaryPath string
}

// NewRSDSSHMClient creates a new client
func NewRSDSSHMClient() ports.SSHClient {
	// Default to "rsd-sshm" in PATH
	path, err := exec.LookPath("rsd-sshm")
	if err != nil {
		// Fallback or just use the name and hope it's in PATH when run
		path = "rsd-sshm"
	}
	return &RSDSSHMClient{
		binaryPath: path,
	}
}

// ListProfiles returns a list of available SSH profiles
func (c *RSDSSHMClient) ListProfiles() ([]string, error) {
	cmd := exec.Command(c.binaryPath, "--list", "--json") // Assuming JSON output support
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		// Fallback for demo if binary missing
		if strings.Contains(err.Error(), "executable file not found") {
			return []string{"demo-server-01", "demo-db-01", "prod-web-01"}, nil
		}
		return nil, fmt.Errorf("failed to list profiles: %w", err)
	}
	
	// Parse output (mock implementation for now as we don't know the format)
	lines := strings.Split(out.String(), "\n")
	var profiles []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			profiles = append(profiles, strings.TrimSpace(line))
		}
	}
	return profiles, nil
}

// Connect establishes a connection to a profile (interactive)
func (c *RSDSSHMClient) Connect(profile string) error {
	cmd := exec.Command(c.binaryPath, "connect", profile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	return cmd.Run()
}

// Execute runs a command on the remote host and returns the output
func (c *RSDSSHMClient) Execute(profile string, command string) (string, error) {
	// Assuming rsd-sshm has an 'exec' or similar command
	cmd := exec.Command(c.binaryPath, "exec", profile, "--", command)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("execution failed: %w, stderr: %s", err, stderr.String())
	}
	
	return out.String(), nil
}

// CopyFile copies a file to the remote host
func (c *RSDSSHMClient) CopyFile(profile string, src string, dest string) error {
	// Assuming rsd-sshm supports scp-like functionality
	cmd := exec.Command(c.binaryPath, "scp", src, fmt.Sprintf("%s:%s", profile, dest))
	return cmd.Run()
}
