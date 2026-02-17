package ports

import "net"

// NetworkInterface represents detailed interface info
type NetworkInterface struct {
	Index        int
	MTU          int
	Name         string
	HardwareAddr string
	Flags        net.Flags
	IPAddresses  []string
}

// NetworkManager defines operations for network configuration and status
type NetworkManager interface {
	GetInterfaces() ([]NetworkInterface, error)
	GetHostname() (string, error)
	GetDNSConfig() ([]string, error)
}
