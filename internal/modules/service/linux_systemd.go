//go:build linux

package service

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/mascli/troncli/internal/core/ports"
)

type LinuxServiceManager struct{}

func NewLinuxServiceManager() ports.ServiceManager {
	return &LinuxServiceManager{}
}

func (m *LinuxServiceManager) ListServices() ([]ports.ServiceStatus, error) {
	// systemctl list-units --type=service --all
	cmd := exec.Command("systemctl", "list-units", "--type=service", "--all", "--no-pager", "--plain")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var services []ports.ServiceStatus
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 4 {
			services = append(services, ports.ServiceStatus{
				Name:   fields[0],
				Status: fields[2], // active/inactive
				State:  fields[3], // running/exited/dead
			})
		}
	}
	return services, nil
}

func (m *LinuxServiceManager) StartService(name string) error {
	return exec.Command("systemctl", "start", name).Run()
}

func (m *LinuxServiceManager) StopService(name string) error {
	return exec.Command("systemctl", "stop", name).Run()
}

func (m *LinuxServiceManager) RestartService(name string) error {
	return exec.Command("systemctl", "restart", name).Run()
}

func (m *LinuxServiceManager) GetServiceStatus(name string) (*ports.ServiceStatus, error) {
	// systemctl status name
	cmd := exec.Command("systemctl", "status", name, "--no-pager")
	output, err := cmd.Output()
	// Exit code 3 means inactive, which is fine, but Run() might return error.
	// We should parse output.
	status := "unknown"
	state := "unknown"

	if err != nil {
		// If service not found or error
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() == 4 {
				return nil, fmt.Errorf("service not found")
			}
		}
	}

	outStr := string(output)
	if strings.Contains(outStr, "Active: active (running)") {
		status = "active"
		state = "running"
	} else if strings.Contains(outStr, "Active: inactive") {
		status = "inactive"
		state = "dead"
	}

	return &ports.ServiceStatus{
		Name:   name,
		Status: status,
		State:  state,
	}, nil
}

func (m *LinuxServiceManager) GetServiceLogs(name string, lines int) (string, error) {
	cmd := exec.Command("journalctl", "-u", name, "-n", fmt.Sprintf("%d", lines), "--no-pager")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func (m *LinuxServiceManager) EnableService(name string) error {
	return exec.Command("systemctl", "enable", name).Run()
}

func (m *LinuxServiceManager) DisableService(name string) error {
	return exec.Command("systemctl", "disable", name).Run()
}
