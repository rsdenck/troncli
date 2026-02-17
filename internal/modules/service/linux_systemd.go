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

func (m *LinuxServiceManager) ListServices() ([]ports.ServiceUnit, error) {
	// systemctl list-units --type=service --all
	cmd := exec.Command("systemctl", "list-units", "--type=service", "--all", "--no-pager", "--plain")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var services []ports.ServiceUnit
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 4 {
			services = append(services, ports.ServiceUnit{
				Name:        fields[0],
				Status:      fields[2], // active/inactive
				ActiveState: fields[3], // running/exited/dead
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

func (m *LinuxServiceManager) GetServiceStatus(name string) (string, error) {
	// systemctl status name
	cmd := exec.Command("systemctl", "is-active", name)
	output, err := cmd.Output()
	// Exit code != 0 means inactive or failed usually

	status := strings.TrimSpace(string(output))
	if err != nil {
		// Try to see if it exists
		if status == "" {
			return "unknown", err
		}
	}

	return status, nil
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
