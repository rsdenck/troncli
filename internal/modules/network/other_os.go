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
			Index:        i.Index,
			MTU:          i.MTU,
			Name:         i.Name,
			HardwareAddr: i.HardwareAddr.String(),
			Flags:        i.Flags,
			IPAddresses:  ips,
		})
	}
	return result, nil
}

func (m *OtherOSNetworkManager) GetHostname() (string, error) {
	return os.Hostname()
}

func (m *OtherOSNetworkManager) GetDNSConfig() ([]string, error) {
	return nil, errors.New("DNS config parsing not supported on this OS")
}
