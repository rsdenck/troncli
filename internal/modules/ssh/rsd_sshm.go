package ssh

// Package ssh provides SSH key management capabilities.

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
	tunnels    map[string]*exec.Cmd
}

type sshProfile struct {
	Name string   `json:"name"`
	Host string   `json:"host"`
	User string   `json:"user"`
	Tags []string `json:"tags"`
}

// NewRSDSSHMClient creates a new client
func NewRSDSSHMClient() (ports.SSHClient, error) {
	path, err := exec.LookPath("rsd-sshm")
	if err != nil {
		// Return client with empty path to allow app startup
		// Methods will check for empty path and return error
		return &RSDSSHMClient{
			binaryPath: "",
			tunnels:    make(map[string]*exec.Cmd),
		}, nil
	}
	return &RSDSSHMClient{
		binaryPath: path,
		tunnels:    make(map[string]*exec.Cmd),
	}, nil
}

// ListProfiles returns a list of available SSH profiles
func (c *RSDSSHMClient) ListProfiles() ([]string, error) {
	if c.binaryPath == "" {
		return nil, fmt.Errorf("rsd-sshm binary not found. Please install it to use SSH features")
	}
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
	if c.binaryPath == "" {
		return fmt.Errorf("rsd-sshm binary not found")
	}
	cmd := exec.Command(c.binaryPath, "connect", profile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// Execute runs a command on the remote host and returns the output
func (c *RSDSSHMClient) Execute(profile string, command string) (string, error) {
	if c.binaryPath == "" {
		return "", fmt.Errorf("rsd-sshm binary not found")
	}
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

// CreateTunnel establishes an SSH tunnel
func (c *RSDSSHMClient) CreateTunnel(profile string, localPort, remoteHost, remotePort string, reverse bool) error {
	// Construct tunnel argument
	// Local: -L localPort:remoteHost:remotePort
	// Remote: -R remotePort:localHost:localPort (simplified mapping)

	tunnelArg := fmt.Sprintf("%s:%s:%s", localPort, remoteHost, remotePort)
	flag := "-L"
	if reverse {
		flag = "-R"
	}

	// We use "exec" to pass arguments to the underlying ssh command
	// assuming rsd-sshm exec profile -- [args...] works and passes args to ssh
	// We use -N (no command) and -f (background)? No, exec.Command runs in background if we don't wait.
	// We just need -N to not execute a shell.

	cmd := exec.Command(c.binaryPath, "exec", profile, "--", flag, tunnelArg, "-N")

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start tunnel: %w", err)
	}

	// Store reference to kill later
	key := fmt.Sprintf("%s-%s", profile, localPort)
	c.tunnels[key] = cmd

	return nil
}

// CloseTunnel closes an active tunnel
func (c *RSDSSHMClient) CloseTunnel(profile string, localPort string) error {
	key := fmt.Sprintf("%s-%s", profile, localPort)
	cmd, exists := c.tunnels[key]
	if !exists {
		return fmt.Errorf("tunnel not found: %s", key)
	}

	if err := cmd.Process.Kill(); err != nil {
		return fmt.Errorf("failed to kill tunnel process: %w", err)
	}

	delete(c.tunnels, key)
	return nil
}

// CopyFile copies a file to the remote host
func (c *RSDSSHMClient) CopyFile(profile string, src string, dest string) error {
	if c.binaryPath == "" {
		return fmt.Errorf("rsd-sshm binary not found")
	}
	cmd := exec.Command(c.binaryPath, "scp", src, fmt.Sprintf("%s:%s", profile, dest))
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("scp failed: %w, stderr: %s", err, stderr.String())
	}
	return nil
}
