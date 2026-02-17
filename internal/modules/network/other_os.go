//go:build !linux

package network

import (
	"errors"
	"net"
	"os"

	"github.com/mascli/troncli/internal/core/ports"
)

type OtherOSNetworkManager struct{}

func NewLinuxNetworkManager() ports.NetworkManager {
	return &OtherOSNetworkManager{}
}

func (m *OtherOSNetworkManager) GetInterfaces() ([]ports.NetworkInterface, error) {
	// Generic implementation works
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var result []ports.NetworkInterface
	for _, i := range ifaces {
		var ips []string
		addrs, err := i.Addrs()
		if err == nil {
			for _, addr := range addrs {
				ips = append(ips, addr.String())
			}
		}

		result = append(result, ports.NetworkInterface{
			Name:         i.Name,
			State:        i.Flags.String(),
			MTU:          i.MTU,
			IPAddresses:  ips,
			HardwareAddr: i.HardwareAddr.String(),
		})
	}
	return result, nil
}

func (m *OtherOSNetworkManager) ApplyConfig(config ports.NetworkConfig) error {
	return errors.New("network configuration not supported on this OS")
}

func (m *OtherOSNetworkManager) GetHostname() (string, error) {
	return os.Hostname()
}

func (m *OtherOSNetworkManager) SetHostname(hostname string) error {
	return errors.New("hostname setting not supported on this OS")
}

func (m *OtherOSNetworkManager) GetDNSConfig() ([]string, error) {
	return nil, errors.New("DNS config not supported on this OS")
}

func (m *OtherOSNetworkManager) SetDNSConfig(servers []string) error {
	return errors.New("DNS setting not supported on this OS")
}

func (m *OtherOSNetworkManager) BackupConfig() error {
	return errors.New("backup config not supported on this OS")
}

func (m *OtherOSNetworkManager) SetInterfaceState(name string, up bool) error {
	return errors.New("interface state not supported on this OS")
}

func (m *OtherOSNetworkManager) ValidateConfig(config ports.NetworkConfig) error {
	return errors.New("config validation not supported on this OS")
}

func (m *OtherOSNetworkManager) GetActiveStack() (string, error) {
	return "unknown", nil
}

func (m *OtherOSNetworkManager) GetSocketStats() ([]ports.SocketStat, error) {
	return nil, errors.New("socket stats not supported on this OS")
}

func (m *OtherOSNetworkManager) GetNftablesRules() ([]string, error) {
	return nil, errors.New("nftables not supported on this OS")
}

func (m *OtherOSNetworkManager) RunTraceRoute(target string) (string, error) {
	return "", errors.New("traceroute not supported on this OS")
}

func (m *OtherOSNetworkManager) RunDig(target string) (string, error) {
	return "", errors.New("dig not supported on this OS")
}

func (m *OtherOSNetworkManager) RunNmap(target string, options string) (string, error) {
	return "", errors.New("nmap not supported on this OS")
}

func (m *OtherOSNetworkManager) RunTcpdump(interfaceName string, filter string, durationSeconds int) (string, error) {
	return "", errors.New("tcpdump not supported on this OS")
}
