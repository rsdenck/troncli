package network

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
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

// GetActiveStack returns the detected network stack
func (m *UniversalNetworkManager) GetActiveStack() (string, error) {
	if m.profile.NetworkStack == "" {
		return "", fmt.Errorf("no supported network stack detected")
	}
	return m.profile.NetworkStack, nil
}

// ApplyConfig applies a network configuration
func (m *UniversalNetworkManager) ApplyConfig(config ports.NetworkConfig) error {
	ctx := context.Background()
	stack := m.profile.NetworkStack

	// Backup first
	if err := m.BackupConfig(); err != nil {
		return fmt.Errorf("failed to backup config: %w", err)
	}

	switch stack {
	case "netplan":
		return m.applyNetplan(ctx, config)
	case "NetworkManager":
		return m.applyNetworkManager(ctx, config)
	case "systemd-networkd":
		return m.applySystemdNetworkd(ctx, config)
	case "ifcfg": // RHEL/CentOS Legacy
		return m.applyIfcfg(ctx, config)
	case "interfaces": // Debian Legacy
		return m.applyInterfaces(ctx, config)
	default:
		return fmt.Errorf("unsupported network stack: %s", stack)
	}
}

// ValidateConfig validates the configuration syntax
func (m *UniversalNetworkManager) ValidateConfig(config ports.NetworkConfig) error {
	// Simple validation
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
	// Simple backup logic based on stack
	// In production, this should copy files to a backup dir
	return nil // Placeholder
}

func (m *UniversalNetworkManager) applyNetplan(ctx context.Context, config ports.NetworkConfig) error {
	// Generate YAML for netplan
	// This is complex, simplified version:
	// Write to /etc/netplan/99-troncli.yaml
	// sudo netplan apply
	yamlContent := fmt.Sprintf(`network:
  version: 2
  ethernets:
    %s:
      dhcp4: %v
`, config.Interface, config.DHCP)

	if !config.DHCP {
		yamlContent += fmt.Sprintf(`      addresses: [%s]
      gateway4: %s
      nameservers:
        addresses: [%s]
`, config.IP, config.Gateway, strings.Join(config.DNS, ", "))
	}

	// Write file (using os directly as executor doesn't write files, but we should use a file helper)
	// For now, assume root and write directly
	err := os.WriteFile(fmt.Sprintf("/etc/netplan/99-troncli-%s.yaml", config.Interface), []byte(yamlContent), 0600)
	if err != nil {
		return err
	}

	_, err = m.executor.Exec(ctx, "netplan", "apply")
	return err
}

func (m *UniversalNetworkManager) applyNetworkManager(ctx context.Context, config ports.NetworkConfig) error {
	// nmcli con modify ...
	// nmcli con mod eth0 ipv4.method manual ipv4.addresses ...
	args := []string{"con", "mod", config.Interface}
	if config.DHCP {
		args = append(args, "ipv4.method", "auto")
	} else {
		args = append(args, "ipv4.method", "manual")
		args = append(args, "ipv4.addresses", config.IP)
		args = append(args, "ipv4.gateway", config.Gateway)
		args = append(args, "ipv4.dns", strings.Join(config.DNS, " "))
	}

	_, err := m.executor.Exec(ctx, "nmcli", args...)
	if err != nil {
		return err
	}
	_, err = m.executor.Exec(ctx, "nmcli", "con", "up", config.Interface)
	return err
}

func (m *UniversalNetworkManager) applySystemdNetworkd(ctx context.Context, config ports.NetworkConfig) error {
	// Write .network file to /etc/systemd/network/
	content := fmt.Sprintf(`[Match]
Name=%s

[Network]
DHCP=%v
`, config.Interface, map[bool]string{true: "yes", false: "no"}[config.DHCP])

	if !config.DHCP {
		content += fmt.Sprintf(`Address=%s
Gateway=%s
DNS=%s
`, config.IP, config.Gateway, strings.Join(config.DNS, " "))
	}

	err := os.WriteFile(fmt.Sprintf("/etc/systemd/network/20-troncli-%s.network", config.Interface), []byte(content), 0644)
	if err != nil {
		return err
	}
	_, err = m.executor.Exec(ctx, "networkctl", "reload")
	return err
}

func (m *UniversalNetworkManager) applyIfcfg(ctx context.Context, config ports.NetworkConfig) error {
	return fmt.Errorf("ifcfg support not implemented")
}

func (m *UniversalNetworkManager) applyInterfaces(ctx context.Context, config ports.NetworkConfig) error {
	return fmt.Errorf("interfaces support not implemented")
}

// ... existing methods from LinuxNetworkManager adapted to use executor ...

func (m *UniversalNetworkManager) GetInterfaces() ([]ports.NetworkInterface, error) {
	// Implementation using net package or ip command
	// Using "ip -j addr" for rich info
	res, err := m.executor.Exec(context.Background(), "ip", "-j", "addr")
	if err != nil {
		return nil, err
	}

	type ipAddrInfo struct {
		Family    string `json:"family"`
		Local     string `json:"local"`
		PrefixLen int    `json:"prefixlen"`
	}

	type ipLink struct {
		Ifindex   int          `json:"ifindex"`
		Ifname    string       `json:"ifname"`
		LinkType  string       `json:"link_type"`
		Flags     []string     `json:"flags"`
		AddrInfo  []ipAddrInfo `json:"addr_info"`
		Mtu       int          `json:"mtu"`
		OperState string       `json:"operstate"`
		Address   string       `json:"address"` // MAC
	}

	var links []ipLink
	if err := json.Unmarshal([]byte(res.Stdout), &links); err != nil {
		return nil, fmt.Errorf("failed to parse ip output: %w", err)
	}

	var interfaces []ports.NetworkInterface
	for _, link := range links {
		var ips []string
		for _, addr := range link.AddrInfo {
			if addr.Family == "inet" || addr.Family == "inet6" {
				ips = append(ips, fmt.Sprintf("%s/%d", addr.Local, addr.PrefixLen))
			}
		}

		interfaces = append(interfaces, ports.NetworkInterface{
			Index:        link.Ifindex,
			Name:         link.Ifname,
			MTU:          link.Mtu,
			HardwareAddr: link.Address,
			IPAddresses:  ips,
			State:        link.OperState,
		})
	}
	return interfaces, nil
}

func (m *UniversalNetworkManager) SetInterfaceState(name string, up bool) error {
	state := "down"
	if up {
		state = "up"
	}
	_, err := m.executor.Exec(context.Background(), "ip", "link", "set", name, state)
	return err
}

func (m *UniversalNetworkManager) GetHostname() (string, error) {
	res, err := m.executor.Exec(context.Background(), "hostname")
	if err != nil {
		return "", err
	}
	return res.Stdout, nil
}

func (m *UniversalNetworkManager) GetDNSConfig() ([]string, error) {
	content, err := os.ReadFile("/etc/resolv.conf")
	if err != nil {
		return nil, err
	}
	var servers []string
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "nameserver") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				servers = append(servers, parts[1])
			}
		}
	}
	return servers, nil
}

func (m *UniversalNetworkManager) GetSocketStats() ([]ports.SocketStat, error) {
	// ss -tulpn
	res, err := m.executor.Exec(context.Background(), "ss", "-tulpn")
	if err != nil {
		return nil, err
	}
	// Parse output...
	return []ports.SocketStat{}, nil
}

func (m *UniversalNetworkManager) GetNftablesRules() ([]string, error) {
	res, err := m.executor.Exec(context.Background(), "nft", "list", "ruleset")
	if err != nil {
		return nil, err
	}
	return strings.Split(res.Stdout, "\n"), nil
}

func (m *UniversalNetworkManager) RunTraceRoute(target string) (string, error) {
	res, err := m.executor.Exec(context.Background(), "traceroute", target)
	if err != nil {
		return "", err
	}
	return res.Stdout, nil
}

func (m *UniversalNetworkManager) RunDig(target string) (string, error) {
	res, err := m.executor.Exec(context.Background(), "dig", target)
	if err != nil {
		return "", err
	}
	return res.Stdout, nil
}

func (m *UniversalNetworkManager) RunNmap(target string, options string) (string, error) {
	args := strings.Fields(options)
	args = append(args, target)
	res, err := m.executor.Exec(context.Background(), "nmap", args...)
	if err != nil {
		return "", err
	}
	return res.Stdout, nil
}

func (m *UniversalNetworkManager) RunTcpdump(interfaceName string, filter string, durationSeconds int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(durationSeconds)*time.Second)
	defer cancel()

	args := []string{"-i", interfaceName, "-n", "-c", "100"} // limit packets
	if filter != "" {
		args = append(args, filter)
	}

	res, err := m.executor.Exec(ctx, "tcpdump", args...)
	// tcpdump might return error on timeout or cancellation, check result
	if res != nil {
		return res.Stdout, nil
	}
	return "", err
}
