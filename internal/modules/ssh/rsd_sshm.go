package ssh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/mascli/troncli/internal/core/ports"
)

// RSDSSHMClient implements ports.SSHClient using rsd-sshm binary
type RSDSSHMClient struct {
	binaryPath string
}

type sshProfile struct {
	Name string `json:"name"`
	Host string `json:"host"`
	User string `json:"user"`
	Tags []string `json:"tags"`
}

// NewRSDSSHMClient creates a new client
func NewRSDSSHMClient() (ports.SSHClient, error) {
	path, err := exec.LookPath("rsd-sshm")
	if err != nil {
		return nil, fmt.Errorf("rsd-sshm binary not found in PATH: %w", err)
	}
	return &RSDSSHMClient{
		binaryPath: path,
	}, nil
}

// ListProfiles returns a list of available SSH profiles
func (c *RSDSSHMClient) ListProfiles() ([]string, error) {
	// Execute real command to get profiles in JSON format
	cmd := exec.Command(c.binaryPath, "--list", "--json")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to list profiles: %w, stderr: %s", err, stderr.String())
	}
	
	var profiles []sshProfile
	if err := json.Unmarshal(out.Bytes(), &profiles); err != nil {
		return nil, fmt.Errorf("failed to parse rsd-sshm output: %w", err)
	}
	
	result := make([]string, len(profiles))
	for i, p := range profiles {
		result[i] = p.Name
	}
	
	return result, nil
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
	cmd := exec.Command(c.binaryPath, "scp", src, fmt.Sprintf("%s:%s", profile, dest))
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("scp failed: %w, stderr: %s", err, stderr.String())
	}
	return nil
}
