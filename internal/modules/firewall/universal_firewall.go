package firewall

// Package firewall provides firewall management capabilities.

import (
	"context"
	"fmt"
	"strings"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/domain"
	"github.com/mascli/troncli/internal/core/ports"
)

const (
	firewallCmd  = "firewall-cmd"
	iptablesCmd  = "iptables"
	nftCmd       = "nft"
	firewalldSvc = "firewalld"
	ufwCmd       = "ufw"
	nftablesSvc  = "nftables"
)

// UniversalFirewallManager implements ports.FirewallManager
type UniversalFirewallManager struct {
	executor adapter.Executor
	profile  *domain.SystemProfile
}

// NewUniversalFirewallManager creates a new firewall manager
func NewUniversalFirewallManager(executor adapter.Executor, profile *domain.SystemProfile) *UniversalFirewallManager {
	return &UniversalFirewallManager{
		executor: executor,
		profile:  profile,
	}
}

// DetectFirewall returns the detected firewall manager
func (m *UniversalFirewallManager) DetectFirewall() (string, error) {
	if m.profile.Firewall == "" {
		return "", fmt.Errorf("no supported firewall manager detected in profile")
	}
	return m.profile.Firewall, nil
}

// AllowPort allows incoming traffic on a port
func (m *UniversalFirewallManager) AllowPort(port string, protocol string) error {
	ctx := context.Background()
	var args []string
	cmd := m.profile.Firewall

	// Adjust command name for firewalld
	execCmd := cmd
	if cmd == firewalldSvc {
		execCmd = firewallCmd
	}

	switch cmd {
	case ufwCmd:
		args = []string{"allow", fmt.Sprintf("%s/%s", port, protocol)}
	case firewalldSvc:
		args = []string{"--permanent", "--add-port", fmt.Sprintf("%s/%s", port, protocol)}
	case iptablesCmd:
		args = []string{"-A", "INPUT", "-p", protocol, "--dport", port, "-j", "ACCEPT"}
	case nftablesSvc:
		// nft add rule inet filter input tcp dport 80 accept
		// Simplified assumption: table inet filter, chain input exists
		args = []string{"add", "rule", "inet", "filter", "input", protocol, "dport", port, "accept"}
		execCmd = nftCmd
	default:
		return fmt.Errorf("unsupported firewall manager: %s", cmd)
	}

	_, err := m.executor.Exec(ctx, execCmd, args...)

	// Reload firewalld if needed
	if err == nil && cmd == firewalldSvc {
		_, _ = m.executor.Exec(ctx, firewallCmd, "--reload")
	}
	return err
}

// BlockPort blocks incoming traffic on a port
func (m *UniversalFirewallManager) BlockPort(port string, protocol string) error {
	ctx := context.Background()
	var args []string
	cmd := m.profile.Firewall

	execCmd := cmd
	if cmd == firewalldSvc {
		execCmd = firewallCmd
	}

	switch cmd {
	case ufwCmd:
		args = []string{"deny", fmt.Sprintf("%s/%s", port, protocol)}
	case firewalldSvc:
		args = []string{"--permanent", "--remove-port", fmt.Sprintf("%s/%s", port, protocol)}
	case iptablesCmd:
		args = []string{"-D", "INPUT", "-p", protocol, "--dport", port, "-j", "ACCEPT"}
	case nftablesSvc:
		// nft add rule inet filter input tcp dport 80 drop
		args = []string{"add", "rule", "inet", "filter", "input", protocol, "dport", port, "drop"}
		execCmd = nftCmd
	default:
		return fmt.Errorf("unsupported firewall manager: %s", cmd)
	}

	_, err := m.executor.Exec(ctx, execCmd, args...)
	if err == nil && cmd == firewalldSvc {
		_, _ = m.executor.Exec(ctx, firewallCmd, "--reload")
	}
	return err
}

// ListRules returns the current rules
func (m *UniversalFirewallManager) ListRules() ([]ports.FirewallRule, error) {
	ctx := context.Background()
	var args []string
	cmd := m.profile.Firewall
	execCmd := cmd

	if cmd == firewalldSvc {
		execCmd = firewallCmd
		args = []string{"--list-all"}
	} else if cmd == ufwCmd {
		args = []string{"status", "numbered"}
	} else if cmd == iptablesCmd {
		args = []string{"-L", "-n", "--line-numbers"}
	} else if cmd == nftablesSvc {
		execCmd = nftCmd
		args = []string{"list", "ruleset"}
	}

	res, err := m.executor.Exec(ctx, execCmd, args...)
	if err != nil {
		return nil, err
	}

	// Parsing logic based on detected firewall
	switch cmd {
	case ufwCmd:
		return m.parseUfwRules(res.Stdout), nil
	case iptablesCmd:
		return m.parseIptablesRules(res.Stdout), nil
	default:
		// For others, return raw output as comment in a single rule for now
		// or implement specific parsers
		return []ports.FirewallRule{{
			ID:      "RAW",
			Comment: res.Stdout,
		}}, nil
	}
}

func (m *UniversalFirewallManager) parseUfwRules(output string) []ports.FirewallRule {
	lines := strings.Split(output, "\n")
	rules := make([]ports.FirewallRule, 0, len(lines))

	// Skip header lines until we find "--" separator or just look for patterns
	// UFW status numbered:
	// Status: active
	//
	//      To                         Action      From
	//      --                         ------      ----
	// [ 1] 22/tcp                     ALLOW IN    Anywhere

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "[") && strings.Contains(line, "]") {
			// [ 1] 22/tcp                     ALLOW IN    Anywhere
			parts := strings.Fields(line)
			if len(parts) >= 4 {
				id := strings.Trim(parts[0], "[]")
				portProto := parts[1]
				action := parts[2]

				pp := strings.Split(portProto, "/")
				port := pp[0]
				proto := "tcp"
				if len(pp) > 1 {
					proto = pp[1]
				}

				rules = append(rules, ports.FirewallRule{
					ID:       id,
					Action:   action,
					Protocol: proto,
					Port:     port,
					Comment:  line,
				})
			}
		}
	}
	return rules
}

func (m *UniversalFirewallManager) parseIptablesRules(output string) []ports.FirewallRule {
	lines := strings.Split(output, "\n")
	rules := make([]ports.FirewallRule, 0, len(lines))

	// Chain INPUT (policy ACCEPT)
	// num  target     prot opt source               destination
	// 1    ACCEPT     tcp  --  0.0.0.0/0            0.0.0.0/0            tcp dpt:22

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 5 || fields[0] == "num" || strings.HasPrefix(fields[0], "Chain") {
			continue
		}

		// 1 ACCEPT tcp -- 0.0.0.0/0 0.0.0.0/0 tcp dpt:22
		id := fields[0]
		action := fields[1]
		proto := fields[2]

		// Extract port if possible (e.g., dpt:22)
		port := "any"
		for _, f := range fields {
			if strings.HasPrefix(f, "dpt:") {
				port = strings.TrimPrefix(f, "dpt:")
			}
		}

		rules = append(rules, ports.FirewallRule{
			ID:       id,
			Action:   action,
			Protocol: proto,
			Port:     port,
			Comment:  line,
		})
	}
	return rules
}

// Enable enables the firewall
func (m *UniversalFirewallManager) Enable() error {
	ctx := context.Background()
	cmd := m.profile.Firewall

	switch cmd {
	case ufwCmd:
		_, err := m.executor.Exec(ctx, ufwCmd, "enable")
		return err
	case firewalldSvc:
		_, err := m.executor.Exec(ctx, "systemctl", "enable", "--now", firewalldSvc)
		return err
	case iptablesCmd:
		// No direct enable, assume service
		_, err := m.executor.Exec(ctx, "systemctl", "enable", "--now", iptablesCmd)
		return err
	case nftablesSvc:
		_, err := m.executor.Exec(ctx, "systemctl", "enable", "--now", nftablesSvc)
		return err
	}
	return fmt.Errorf("unsupported firewall manager: %s", cmd)
}

// Disable disables the firewall
func (m *UniversalFirewallManager) Disable() error {
	ctx := context.Background()
	cmd := m.profile.Firewall

	switch cmd {
	case ufwCmd:
		_, err := m.executor.Exec(ctx, ufwCmd, "disable")
		return err
	case firewalldSvc:
		_, err := m.executor.Exec(ctx, "systemctl", "disable", "--now", firewalldSvc)
		return err
	case iptablesCmd:
		_, err := m.executor.Exec(ctx, "systemctl", "disable", "--now", iptablesCmd)
		return err
	case nftablesSvc:
		_, err := m.executor.Exec(ctx, "systemctl", "disable", "--now", nftablesSvc)
		return err
	}
	return fmt.Errorf("unsupported firewall manager: %s", cmd)
}
