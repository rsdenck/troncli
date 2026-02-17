package service

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/mascli/troncli/internal/core/ports"
)

// LinuxServiceManager implements ports.ServiceManager using systemd
type LinuxServiceManager struct{}

// NewLinuxServiceManager creates a new service manager
func NewLinuxServiceManager() *LinuxServiceManager {
	return &LinuxServiceManager{}
}

// ListServices returns a list of system services
func (m *LinuxServiceManager) ListServices() ([]ports.ServiceUnit, error) {
	// systemctl list-units --type=service --all --no-pager --no-legend
	cmd := exec.Command("systemctl", "list-units", "--type=service", "--all", "--no-pager", "--no-legend")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	var services []ports.ServiceUnit
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		
		// systemctl output: UNIT LOAD ACTIVE SUB DESCRIPTION...
		// But description can contain spaces, so we need to be careful.
		// Fields[0] = name, [1] = load, [2] = active, [3] = sub, [4:] = desc
		
		desc := ""
		if len(fields) > 4 {
			desc = strings.Join(fields[4:], " ")
		}

		services = append(services, ports.ServiceUnit{
			Name:        fields[0],
			LoadState:   fields[1],
			ActiveState: fields[2],
			SubState:    fields[3],
			Description: desc,
		})
	}
	return services, nil
}

// StartService starts a service
func (m *LinuxServiceManager) StartService(name string) error {
	return exec.Command("systemctl", "start", name).Run()
}

// StopService stops a service
func (m *LinuxServiceManager) StopService(name string) error {
	return exec.Command("systemctl", "stop", name).Run()
}

// RestartService restarts a service
func (m *LinuxServiceManager) RestartService(name string) error {
	return exec.Command("systemctl", "restart", name).Run()
}

// EnableService enables a service to start on boot
func (m *LinuxServiceManager) EnableService(name string) error {
	return exec.Command("systemctl", "enable", name).Run()
}

// DisableService disables a service from starting on boot
func (m *LinuxServiceManager) DisableService(name string) error {
	return exec.Command("systemctl", "disable", name).Run()
}

// GetServiceStatus returns the status output of a service
func (m *LinuxServiceManager) GetServiceStatus(name string) (string, error) {
	cmd := exec.Command("systemctl", "status", name, "--no-pager", "-l")
	// Status returns non-zero exit code if service is stopped/failed, but we still want the output
	out, _ := cmd.CombinedOutput()
	return string(out), nil
}

// GetServiceJournal returns the journal logs for a service
func (m *LinuxServiceManager) GetServiceJournal(name string, lines int) (string, error) {
	// journalctl -u name -n lines --no-pager
	cmd := exec.Command("journalctl", "-u", name, "-n", fmt.Sprintf("%d", lines), "--no-pager")
	out, err := cmd.CombinedOutput()
	return string(out), err
}
