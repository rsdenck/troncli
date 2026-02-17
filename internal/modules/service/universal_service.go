package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/domain"
	"github.com/mascli/troncli/internal/core/ports"
)

// UniversalServiceManager implements ports.ServiceManager
type UniversalServiceManager struct {
	executor adapter.Executor
	profile  *domain.SystemProfile
}

// NewUniversalServiceManager creates a new service manager
func NewUniversalServiceManager(executor adapter.Executor, profile *domain.SystemProfile) *UniversalServiceManager {
	return &UniversalServiceManager{
		executor: executor,
		profile:  profile,
	}
}

// ListServices returns a list of system services
func (m *UniversalServiceManager) ListServices() ([]ports.ServiceUnit, error) {
	ctx := context.Background()
	init := m.profile.InitSystem

	switch init {
	case "systemd":
		return m.listSystemdServices(ctx)
	case "sysvinit":
		return m.listSysvServices(ctx)
	case "openrc":
		return m.listOpenRCServices(ctx)
	default:
		return nil, fmt.Errorf("unsupported init system: %s", init)
	}
}

func (m *UniversalServiceManager) listSystemdServices(ctx context.Context) ([]ports.ServiceUnit, error) {
	// systemctl list-units --type=service --all --no-pager --no-legend
	res, err := m.executor.Exec(ctx, "systemctl", "list-units", "--type=service", "--all", "--no-pager", "--no-legend")
	if err != nil {
		return nil, err
	}

	var services []ports.ServiceUnit
	lines := strings.Split(res.Stdout, "\n")
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) >= 4 {
			// unit load active sub description...
			name := parts[0]
			load := parts[1]
			active := parts[2]
			sub := parts[3]
			desc := strings.Join(parts[4:], " ")

			services = append(services, ports.ServiceUnit{
				Name:        name,
				LoadState:   load,
				ActiveState: active,
				SubState:    sub,
				Status:      active, // active/inactive
				Description: desc,
				Enabled:     load == "loaded", // simplification
			})
		}
	}
	return services, nil
}

func (m *UniversalServiceManager) listSysvServices(ctx context.Context) ([]ports.ServiceUnit, error) {
	// service --status-all
	res, err := m.executor.Exec(ctx, "service", "--status-all")
	if err != nil {
		return nil, err
	}

	var services []ports.ServiceUnit
	lines := strings.Split(res.Stdout, "\n")
	for _, line := range lines {
		// [ + ]  nginx
		// [ - ]  apache2
		parts := strings.Fields(line)
		if len(parts) >= 4 {
			status := parts[1]
			name := parts[3]
			
			state := "inactive"
			if status == "+" {
				state = "active"
			}

			services = append(services, ports.ServiceUnit{
				Name:   name,
				Status: state,
			})
		}
	}
	return services, nil
}

func (m *UniversalServiceManager) listOpenRCServices(ctx context.Context) ([]ports.ServiceUnit, error) {
	// rc-status --all
	res, err := m.executor.Exec(ctx, "rc-status", "--all")
	if err != nil {
		return nil, err
	}
	// parsing rc-status output...
	return []ports.ServiceUnit{}, nil // TODO: Implement OpenRC parsing
}

// StartService starts a service
func (m *UniversalServiceManager) StartService(name string) error {
	ctx := context.Background()
	switch m.profile.InitSystem {
	case "systemd":
		_, err := m.executor.Exec(ctx, "systemctl", "start", name)
		return err
	case "sysvinit":
		_, err := m.executor.Exec(ctx, "service", name, "start")
		return err
	case "openrc":
		_, err := m.executor.Exec(ctx, "rc-service", name, "start")
		return err
	case "runit":
		_, err := m.executor.Exec(ctx, "sv", "start", name)
		return err
	}
	return fmt.Errorf("unsupported init system")
}

// StopService stops a service
func (m *UniversalServiceManager) StopService(name string) error {
	ctx := context.Background()
	switch m.profile.InitSystem {
	case "systemd":
		_, err := m.executor.Exec(ctx, "systemctl", "stop", name)
		return err
	case "sysvinit":
		_, err := m.executor.Exec(ctx, "service", name, "stop")
		return err
	case "openrc":
		_, err := m.executor.Exec(ctx, "rc-service", name, "stop")
		return err
	case "runit":
		_, err := m.executor.Exec(ctx, "sv", "stop", name)
		return err
	}
	return fmt.Errorf("unsupported init system")
}

// RestartService restarts a service
func (m *UniversalServiceManager) RestartService(name string) error {
	ctx := context.Background()
	switch m.profile.InitSystem {
	case "systemd":
		_, err := m.executor.Exec(ctx, "systemctl", "restart", name)
		return err
	case "sysvinit":
		_, err := m.executor.Exec(ctx, "service", name, "restart")
		return err
	case "openrc":
		_, err := m.executor.Exec(ctx, "rc-service", name, "restart")
		return err
	case "runit":
		_, err := m.executor.Exec(ctx, "sv", "restart", name)
		return err
	}
	return fmt.Errorf("unsupported init system")
}

// EnableService enables a service
func (m *UniversalServiceManager) EnableService(name string) error {
	ctx := context.Background()
	switch m.profile.InitSystem {
	case "systemd":
		_, err := m.executor.Exec(ctx, "systemctl", "enable", name)
		return err
	case "sysvinit":
		// update-rc.d name defaults (Debian) or chkconfig name on (RHEL)
		// Detecting distro for sysvinit enabling is complex
		return fmt.Errorf("enable not fully implemented for sysvinit")
	case "openrc":
		_, err := m.executor.Exec(ctx, "rc-update", "add", name, "default")
		return err
	}
	return fmt.Errorf("unsupported init system")
}

// DisableService disables a service
func (m *UniversalServiceManager) DisableService(name string) error {
	ctx := context.Background()
	switch m.profile.InitSystem {
	case "systemd":
		_, err := m.executor.Exec(ctx, "systemctl", "disable", name)
		return err
	case "openrc":
		_, err := m.executor.Exec(ctx, "rc-update", "del", name, "default")
		return err
	}
	return fmt.Errorf("unsupported init system")
}

// GetServiceStatus returns the status output
func (m *UniversalServiceManager) GetServiceStatus(name string) (string, error) {
	ctx := context.Background()
	var args []string
	cmd := ""

	switch m.profile.InitSystem {
	case "systemd":
		cmd = "systemctl"
		args = []string{"status", name}
	case "sysvinit":
		cmd = "service"
		args = []string{name, "status"}
	case "openrc":
		cmd = "rc-service"
		args = []string{name, "status"}
	default:
		return "", fmt.Errorf("unsupported init system")
	}

	res, err := m.executor.Exec(ctx, cmd, args...)
	// status command often returns non-zero if service is stopped, but we want the output
	if res != nil {
		return res.Stdout + "\n" + res.Stderr, nil
	}
	return "", err
}

// GetServiceLogs returns the logs for a service
func (m *UniversalServiceManager) GetServiceLogs(name string, lines int) (string, error) {
	ctx := context.Background()
	if m.profile.InitSystem == "systemd" {
		res, err := m.executor.Exec(ctx, "journalctl", "-u", name, "-n", strconv.Itoa(lines), "--no-pager")
		if err != nil {
			return "", err
		}
		return res.Stdout, nil
	}
	return "", fmt.Errorf("logs only supported for systemd currently")
}
