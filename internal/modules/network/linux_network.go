//go:build linux

package network

import (
	"bufio"
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/mascli/troncli/internal/core/ports"
)

type LinuxNetworkManager struct{}

func NewLinuxNetworkManager() ports.NetworkManager {
	return &LinuxNetworkManager{}
}

func (m *LinuxNetworkManager) GetInterfaces() ([]ports.NetworkInterface, error) {
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

func (m *LinuxNetworkManager) GetHostname() (string, error) {
	// Use hostnamectl as requested
	cmd := exec.Command("hostnamectl", "--static")
	out, err := cmd.Output()
	if err == nil {
		return strings.TrimSpace(string(out)), nil
	}
	// Fallback to os.Hostname
	return os.Hostname()
}

func (m *LinuxNetworkManager) GetDNSConfig() ([]string, error) {
	file, err := os.Open("/etc/resolv.conf")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var servers []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) >= 2 && parts[0] == "nameserver" {
			servers = append(servers, parts[1])
		}
	}
	return servers, nil
}
