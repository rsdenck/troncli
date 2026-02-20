package service

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/mascli/troncli/internal/core/ports"
)

type SystemdManager struct{}

func NewSystemdManager() ports.ServiceManager {
	return &SystemdManager{}
}

func (m *SystemdManager) ListServices() ([]ports.ServiceUnit, error) {
	// List all service units
	cmd := exec.Command("systemctl", "list-units", "--type=service", "--all", "--no-pager", "--no-legend")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	var services []ports.ServiceUnit
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		// Fields are separated by spaces, but description can contain spaces.
		// systemctl output is fixed width-ish but better to parse carefully or use JSON output if available.
		// Standard output: UNIT LOAD ACTIVE SUB DESCRIPTION
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		unit := ports.ServiceUnit{
			Name:        fields[0],
			LoadState:   fields[1],
			ActiveState: fields[2],
			SubState:    fields[3],
			Status:      fields[2], // Map ActiveState to Status
		}
		
		if len(fields) > 4 {
			unit.Description = strings.Join(fields[4:], " ")
		}

		services = append(services, unit)
	}
	return services, nil
}

func (m *SystemdManager) StartService(name string) error {
	return exec.Command("systemctl", "start", name).Run()
}

func (m *SystemdManager) StopService(name string) error {
	return exec.Command("systemctl", "stop", name).Run()
}

func (m *SystemdManager) RestartService(name string) error {
	return exec.Command("systemctl", "restart", name).Run()
}

func (m *SystemdManager) EnableService(name string) error {
	return exec.Command("systemctl", "enable", name).Run()
}

func (m *SystemdManager) DisableService(name string) error {
	return exec.Command("systemctl", "disable", name).Run()
}

func (m *SystemdManager) GetServiceStatus(name string) (string, error) {
	output, err := exec.Command("systemctl", "status", name, "--no-pager").CombinedOutput()
	// systemctl status returns non-zero if service is not running, which is not necessarily an error for us
	// but we return the output anyway.
	if err != nil && len(output) == 0 {
		return "", err
	}
	return string(output), nil
}

func (m *SystemdManager) GetServiceLogs(name string, lines int) (string, error) {
	cmd := exec.Command("journalctl", "-u", name, "-n", fmt.Sprintf("%d", lines), "--no-pager")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
