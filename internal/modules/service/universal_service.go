package service

// Package service provides service management capabilities.

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/domain"
	"github.com/mascli/troncli/internal/core/ports"
	"github.com/mascli/troncli/internal/policy"
)

// UniversalServiceManager implements ports.ServiceManager
type UniversalServiceManager struct {
	executor     adapter.Executor
	profile      *domain.SystemProfile
	policyEngine *policy.PolicyEngine
}

// NewUniversalServiceManager creates a new service manager
func NewUniversalServiceManager(executor adapter.Executor, profile *domain.SystemProfile, policyEngine *policy.PolicyEngine) *UniversalServiceManager {
	return &UniversalServiceManager{
		executor:     executor,
		profile:      profile,
		policyEngine: policyEngine,
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
	case "runit":
		return m.listRunitServices(ctx)
	default:
		return nil, fmt.Errorf("unsupported init system: %s", init)
	}
}

type systemdUnit struct {
	Unit        string `json:"unit"`
	Load        string `json:"load"`
	Active      string `json:"active"`
	Sub         string `json:"sub"`
	Description string `json:"description"`
}

func (m *UniversalServiceManager) listSystemdServices(ctx context.Context) ([]ports.ServiceUnit, error) {
	// Try JSON output first (systemd v218+)
	res, err := m.executor.Exec(ctx, "systemctl", "list-units", "--type=service", "--all", "--output=json")
	if err == nil {
		var units []systemdUnit
		if jsonErr := json.Unmarshal([]byte(res.Stdout), &units); jsonErr == nil {
			var services []ports.ServiceUnit
			for _, u := range units {
				services = append(services, ports.ServiceUnit{
					Name:        u.Unit,
					LoadState:   u.Load,
					ActiveState: u.Active,
					SubState:    u.Sub,
					Status:      u.Active,
					Description: u.Description,
					Enabled:     u.Load == "loaded",
				})
			}
			return services, nil
		}
	}

	// Fallback to text parsing if JSON fails or not supported
	res, err = m.executor.Exec(ctx, "systemctl", "list-units", "--type=service", "--all", "--no-pager", "--no-legend")
	if err != nil {
		return nil, err
	}

	// Group 1: Unit
	// Group 2: Load
	// Group 3: Active
	// Group 4: Sub
	// Group 5: Description (optional)
	// Example: cron.service loaded active running Regular background program processing daemon
	// Regex breakdown:
	// ^(\S+)      -> Unit name (non-whitespace)
	// \s+         -> Separator
	// (\S+)       -> Load state
	// \s+         -> Separator
	// (\S+)       -> Active state
	// \s+         -> Separator
	// (\S+)       -> Sub state
	// \s*         -> Optional separator
	// (.*)$       -> Description (rest of line)
	reService := regexp.MustCompile(`^(\S+)\s+(\S+)\s+(\S+)\s+(\S+)\s*(.*)$`)

	var services []ports.ServiceUnit
	lines := strings.Split(res.Stdout, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		matches := reService.FindStringSubmatch(line)
		if len(matches) >= 5 {
			// Skip header lines if they match accidentally
			if matches[1] == "UNIT" {
				continue
			}

			name := matches[1]
			load := matches[2]
			active := matches[3]
			sub := matches[4]
			desc := matches[5]

			services = append(services, ports.ServiceUnit{
				Name:        name,
				LoadState:   load,
				ActiveState: active,
				SubState:    sub,
				Status:      active,
				Description: desc,
				Enabled:     load == "loaded",
			})
		}
	}
	return services, nil
}

func (m *UniversalServiceManager) listSysvServices(ctx context.Context) ([]ports.ServiceUnit, error) {
	// service --status-all
	// Output format:
	// [ + ]  nginx
	// [ - ]  apache2
	// [ ? ]  unknown
	res, err := m.executor.Exec(ctx, "service", "--status-all")
	if err != nil {
		return nil, err
	}

	// Regex for sysvinit status
	// Group 1: Status (+, -, ?)
	// Group 2: Name
	reSysv := regexp.MustCompile(`\[\s*([+\-?])\s*\]\s+(\S+)`)

	var services []ports.ServiceUnit
	lines := strings.Split(res.Stdout, "\n")
	for _, line := range lines {
		matches := reSysv.FindStringSubmatch(line)
		if len(matches) >= 3 {
			statusSymbol := matches[1]
			name := matches[2]

			state := "inactive"
			if statusSymbol == "+" {
				state = "active"
			} else if statusSymbol == "?" {
				state = "unknown"
			}

			services = append(services, ports.ServiceUnit{
				Name:        name,
				Status:      state,
				ActiveState: state,
				LoadState:   "loaded", // Assumed loaded if listed
			})
		}
	}
	return services, nil
}

func (m *UniversalServiceManager) listOpenRCServices(ctx context.Context) ([]ports.ServiceUnit, error) {
	// rc-status --all
	_, err := m.executor.Exec(ctx, "rc-status", "--all")
	if err != nil {
		return nil, err
	}
	// parsing rc-status output...
	return []ports.ServiceUnit{}, nil // TODO: Implement OpenRC parsing
}

func (m *UniversalServiceManager) listRunitServices(ctx context.Context) ([]ports.ServiceUnit, error) {
	// sv status /var/service/*
	res, err := m.executor.Exec(ctx, "sh", "-c", "sv status /var/service/*")
	if err != nil {
		return nil, err
	}

	// Runit sv status output format:
	// run: /var/service/sshd: (pid 1234) 123s
	// down: /var/service/nginx: 0s, normally up
	reRunit := regexp.MustCompile(`^(run|down):\s+([^:]+):\s+(.*)$`)

	var services []ports.ServiceUnit
	lines := strings.Split(res.Stdout, "\n")
	for _, line := range lines {
		matches := reRunit.FindStringSubmatch(line)
		if len(matches) >= 3 {
			state := matches[1]
			servicePath := matches[2]
			
			// Extract service name from path (e.g., /var/service/sshd -> sshd)
			parts := strings.Split(strings.TrimSpace(servicePath), "/")
			name := parts[len(parts)-1]

			activeState := "inactive"
			if state == "run" {
				activeState = "active"
			}

			services = append(services, ports.ServiceUnit{
				Name:        name,
				Status:      activeState,
				ActiveState: activeState,
				LoadState:   "loaded",
			})
		}
	}
	return services, nil
}

// StartService starts a service
func (m *UniversalServiceManager) StartService(name string) error {
	ctx := context.Background()
	init := m.profile.InitSystem

	// Build command string for policy checking
	var cmd string
	switch init {
	case "systemd":
		cmd = fmt.Sprintf("systemctl start %s", name)
	case "sysvinit":
		cmd = fmt.Sprintf("service %s start", name)
	case "openrc":
		cmd = fmt.Sprintf("rc-service %s start", name)
	case "runit":
		cmd = fmt.Sprintf("sv start %s", name)
	default:
		return fmt.Errorf("unsupported init system: %s", init)
	}

	// Check command against policy engine
	if err := m.policyEngine.CheckCommand(cmd); err != nil {
		return fmt.Errorf("policy engine blocked execution: %w", err)
	}

	switch init {
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
	default:
		return fmt.Errorf("unsupported init system: %s", init)
	}
}

// StopService stops a service
func (m *UniversalServiceManager) StopService(name string) error {
	ctx := context.Background()
	init := m.profile.InitSystem

	// Build command string for policy checking
	var cmd string
	switch init {
	case "systemd":
		cmd = fmt.Sprintf("systemctl stop %s", name)
	case "sysvinit":
		cmd = fmt.Sprintf("service %s stop", name)
	case "openrc":
		cmd = fmt.Sprintf("rc-service %s stop", name)
	case "runit":
		cmd = fmt.Sprintf("sv stop %s", name)
	default:
		return fmt.Errorf("unsupported init system: %s", init)
	}

	// Check command against policy engine
	if err := m.policyEngine.CheckCommand(cmd); err != nil {
		return fmt.Errorf("policy engine blocked execution: %w", err)
	}

	switch init {
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
	init := m.profile.InitSystem

	// Build command string for policy checking
	var cmd string
	switch init {
	case "systemd":
		cmd = fmt.Sprintf("systemctl restart %s", name)
	case "sysvinit":
		cmd = fmt.Sprintf("service %s restart", name)
	case "openrc":
		cmd = fmt.Sprintf("rc-service %s restart", name)
	case "runit":
		cmd = fmt.Sprintf("sv restart %s", name)
	default:
		return fmt.Errorf("unsupported init system: %s", init)
	}

	// Check command against policy engine
	if err := m.policyEngine.CheckCommand(cmd); err != nil {
		return fmt.Errorf("policy engine blocked execution: %w", err)
	}

	switch init {
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

// EnableService enables a service to start on boot
func (m *UniversalServiceManager) EnableService(name string) error {
	ctx := context.Background()
	init := m.profile.InitSystem

	// Build command string for policy checking
	var cmd string
	switch init {
	case "systemd":
		cmd = fmt.Sprintf("systemctl enable %s", name)
	case "sysvinit":
		cmd = fmt.Sprintf("update-rc.d %s enable", name)
	case "openrc":
		cmd = fmt.Sprintf("rc-update add %s default", name)
	case "runit":
		cmd = fmt.Sprintf("ln -s /etc/sv/%s /var/service/%s", name, name)
	default:
		return fmt.Errorf("unsupported init system for enable: %s", init)
	}

	// Check command against policy engine
	if err := m.policyEngine.CheckCommand(cmd); err != nil {
		return fmt.Errorf("policy engine blocked execution: %w", err)
	}

	switch init {
	case "systemd":
		_, err := m.executor.Exec(ctx, "systemctl", "enable", name)
		return err
	case "sysvinit":
		_, err := m.executor.Exec(ctx, "update-rc.d", name, "enable")
		return err
	case "openrc":
		_, err := m.executor.Exec(ctx, "rc-update", "add", name, "default")
		return err
	case "runit":
		// In runit, enabling a service means creating a symlink from /etc/sv/[name] to /var/service/[name]
		_, err := m.executor.Exec(ctx, "ln", "-s", fmt.Sprintf("/etc/sv/%s", name), fmt.Sprintf("/var/service/%s", name))
		return err
	}
	return fmt.Errorf("unsupported init system for enable: %s", m.profile.InitSystem)
}

// DisableService disables a service from starting on boot
func (m *UniversalServiceManager) DisableService(name string) error {
	ctx := context.Background()
	init := m.profile.InitSystem

	// Build command string for policy checking
	var cmd string
	switch init {
	case "systemd":
		cmd = fmt.Sprintf("systemctl disable %s", name)
	case "sysvinit":
		cmd = fmt.Sprintf("update-rc.d %s disable", name)
	case "openrc":
		cmd = fmt.Sprintf("rc-update del %s default", name)
	case "runit":
		cmd = fmt.Sprintf("rm /var/service/%s", name)
	default:
		return fmt.Errorf("unsupported init system for disable: %s", init)
	}

	// Check command against policy engine
	if err := m.policyEngine.CheckCommand(cmd); err != nil {
		return fmt.Errorf("policy engine blocked execution: %w", err)
	}

	switch init {
	case "systemd":
		_, err := m.executor.Exec(ctx, "systemctl", "disable", name)
		return err
	case "sysvinit":
		_, err := m.executor.Exec(ctx, "update-rc.d", name, "disable")
		return err
	case "openrc":
		_, err := m.executor.Exec(ctx, "rc-update", "del", name, "default")
		return err
	case "runit":
		// In runit, disabling a service means removing the symlink from /var/service/[name]
		_, err := m.executor.Exec(ctx, "rm", fmt.Sprintf("/var/service/%s", name))
		return err
	}
	return fmt.Errorf("unsupported init system for disable: %s", m.profile.InitSystem)
}

// GetServiceStatus returns the status output of a service
func (m *UniversalServiceManager) GetServiceStatus(name string) (string, error) {
	ctx := context.Background()
	init := m.profile.InitSystem

	// Build command string for policy checking
	var cmd string
	switch init {
	case "systemd":
		cmd = fmt.Sprintf("systemctl status %s", name)
	case "sysvinit":
		cmd = fmt.Sprintf("service %s status", name)
	case "openrc":
		cmd = fmt.Sprintf("rc-service %s status", name)
	case "runit":
		cmd = fmt.Sprintf("sv status %s", name)
	default:
		return "", fmt.Errorf("unsupported init system: %s", init)
	}

	// Check command against policy engine
	if err := m.policyEngine.CheckCommand(cmd); err != nil {
		return "", fmt.Errorf("policy engine blocked execution: %w", err)
	}

	switch init {
	case "systemd":
		res, err := m.executor.Exec(ctx, "systemctl", "status", name)
		return res.Stdout, err
	case "sysvinit":
		res, err := m.executor.Exec(ctx, "service", name, "status")
		return res.Stdout, err
	case "openrc":
		res, err := m.executor.Exec(ctx, "rc-service", name, "status")
		return res.Stdout, err
	case "runit":
		res, err := m.executor.Exec(ctx, "sv", "status", name)
		return res.Stdout, err
	}
	return "", fmt.Errorf("unsupported init system")
}

// GetServiceLogs returns the logs for a service (journald or log file)
func (m *UniversalServiceManager) GetServiceLogs(name string, lines int) (string, error) {
	ctx := context.Background()

	if m.profile.InitSystem == "systemd" {
		// Build command string for policy checking
		cmd := fmt.Sprintf("journalctl -u %s -n %d --no-pager", name, lines)

		// Check command against policy engine
		if err := m.policyEngine.CheckCommand(cmd); err != nil {
			return "", fmt.Errorf("policy engine blocked execution: %w", err)
		}

		res, err := m.executor.Exec(ctx, "journalctl", "-u", name, "-n", fmt.Sprintf("%d", lines), "--no-pager")
		return res.Stdout, err
	}
	return "", fmt.Errorf("log retrieval only supported for systemd")
}
