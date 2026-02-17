package network

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/domain"
	"github.com/mascli/troncli/internal/core/ports"
)

// UniversalNetworkManager implements ports.NetworkManager
type UniversalNetworkManager struct {
	executor adapter.Executor
	profile  *domain.SystemProfile
}

// NewUniversalNetworkManager creates a new network manager
func NewUniversalNetworkManager(executor adapter.Executor, profile *domain.SystemProfile) *UniversalNetworkManager {
	return &UniversalNetworkManager{
		executor: executor,
		profile:  profile,
	}
}

// GetInterfaces returns detailed interface info
func (m *UniversalNetworkManager) GetInterfaces() ([]ports.NetworkInterface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get interfaces: %w", err)
	}

	var result []ports.NetworkInterface
	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		var ipAddrs []string
		for _, addr := range addrs {
			ipAddrs = append(ipAddrs, addr.String())
		}

		state := "DOWN"
		if iface.Flags&net.FlagUp != 0 {
			state = "UP"
		}

		result = append(result, ports.NetworkInterface{
			Index:        iface.Index,
			MTU:          iface.MTU,
			Name:         iface.Name,
			HardwareAddr: iface.HardwareAddr.String(),
			Flags:        iface.Flags,
			IPAddresses:  ipAddrs,
			State:        state,
		})
	}
	return result, nil
}

// SetInterfaceState sets the state of an interface (UP/DOWN)
func (m *UniversalNetworkManager) SetInterfaceState(name string, up bool) error {
	ctx := context.Background()
	state := "down"
	if up {
		state = "up"
	}
	// ip link set dev <name> <state>
	_, err := m.executor.Exec(ctx, "ip", "link", "set", "dev", name, state)
	return err
}

// GetActiveStack returns the detected network stack
func (m *UniversalNetworkManager) GetActiveStack() (string, error) {
	if m.profile.NetworkStack == "" {
		return "unknown", nil
	}
	return m.profile.NetworkStack, nil
}

// ApplyConfig applies a network configuration
func (m *UniversalNetworkManager) ApplyConfig(config ports.NetworkConfig) error {
	// Simplified implementation for demo purposes
	// Real implementation requires handling multiple network managers complex configs
	return fmt.Errorf("apply config not fully implemented for %s", m.profile.NetworkStack)
}

// ValidateConfig validates the configuration syntax
func (m *UniversalNetworkManager) ValidateConfig(config ports.NetworkConfig) error {
	if config.Interface == "" {
		return fmt.Errorf("interface name required")
	}
	if !config.DHCP && config.IP == "" {
		return fmt.Errorf("IP address required for static config")
	}
	return nil
}

// BackupConfig creates a backup of current network configuration
func (m *UniversalNetworkManager) BackupConfig() error {
	return nil // Placeholder
}

// GetHostname returns the system hostname
func (m *UniversalNetworkManager) GetHostname() (string, error) {
	ctx := context.Background()
	res, err := m.executor.Exec(ctx, "hostname")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(res.Stdout), nil
}

// GetDNSConfig returns the current DNS configuration
func (m *UniversalNetworkManager) GetDNSConfig() ([]string, error) {
	ctx := context.Background()
	// cat /etc/resolv.conf
	res, err := m.executor.Exec(ctx, "cat", "/etc/resolv.conf")
	if err != nil {
		return nil, err
	}

	var nameservers []string
	scanner := bufio.NewScanner(strings.NewReader(res.Stdout))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "nameserver") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				nameservers = append(nameservers, parts[1])
			}
		}
	}
	return nameservers, nil
}

// GetSocketStats returns socket statistics (ss)
func (m *UniversalNetworkManager) GetSocketStats() ([]ports.SocketStat, error) {
	ctx := context.Background()
	// ss -tuln
	res, err := m.executor.Exec(ctx, "ss", "-tulnH") // H for no header
	if err != nil {
		return nil, err
	}

	var stats []ports.SocketStat
	scanner := bufio.NewScanner(strings.NewReader(res.Stdout))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) >= 5 {
			stats = append(stats, ports.SocketStat{
				Protocol: parts[0],
				State:    parts[1],
				Local:    parts[3],
				Remote:   parts[4],
				Process:  "", // parsing process info requires -p and sudo
			})
		}
	}
	return stats, nil
}

// GetNftablesRules returns current nftables rules
func (m *UniversalNetworkManager) GetNftablesRules() ([]string, error) {
	ctx := context.Background()
	res, err := m.executor.Exec(ctx, "nft", "list", "ruleset")
	if err != nil {
		return nil, err
	}
	return strings.Split(res.Stdout, "\n"), nil
}

// RunTraceRoute runs traceroute to a target
func (m *UniversalNetworkManager) RunTraceRoute(target string) (string, error) {
	ctx := context.Background()
	res, err := m.executor.Exec(ctx, "traceroute", target)
	if err != nil {
		// Try tracepath if traceroute fails (common on some distros)
		res, err = m.executor.Exec(ctx, "tracepath", target)
		if err != nil {
			return "", err
		}
	}
	return res.Stdout, nil
}

// RunDig runs dig on a target
func (m *UniversalNetworkManager) RunDig(target string) (string, error) {
	ctx := context.Background()
	res, err := m.executor.Exec(ctx, "dig", target)
	if err != nil {
		return "", err
	}
	return res.Stdout, nil
}

// RunNmap runs nmap on a target
func (m *UniversalNetworkManager) RunNmap(target string, options string) (string, error) {
	ctx := context.Background()
	args := []string{}
	if options != "" {
		args = append(args, strings.Fields(options)...)
	}
	args = append(args, target)
	res, err := m.executor.Exec(ctx, "nmap", args...)
	if err != nil {
		return "", err
	}
	return res.Stdout, nil
}

// RunTcpdump runs tcpdump on an interface
func (m *UniversalNetworkManager) RunTcpdump(interfaceName string, filter string, durationSeconds int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(durationSeconds)*time.Second)
	defer cancel()

	args := []string{"-i", interfaceName, "-n", "-c", "10"} // Limit to 10 packets for safety in CLI
	if filter != "" {
		args = append(args, filter)
	}

	res, err := m.executor.Exec(ctx, "tcpdump", args...)
	if err != nil {
		// Timeout is expected
		if ctx.Err() == context.DeadlineExceeded {
			return res.Stdout, nil // Return what we captured
		}
		return "", err
	}
	return res.Stdout, nil
}
