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
	State        string // UP/DOWN
}

// SocketStat represents socket statistics (ss)
type SocketStat struct {
	Protocol string
	State    string
	Local    string
	Remote   string
	Process  string
}

// PortScanResult represents a port scan result
type PortScanResult struct {
	Port     int
	Protocol string
	State    string
	Service  string
}

// NetworkConfig represents a network configuration to apply
type NetworkConfig struct {
	Interface string
	DHCP      bool
	IP        string // CIDR format
	Gateway   string
	DNS       []string
}

// NetworkManager defines operations for network configuration and status
type NetworkManager interface {
	// Interface Management
	GetInterfaces() ([]NetworkInterface, error)
	SetInterfaceState(name string, up bool) error

	// Configuration (Universal Phase 4)
	ApplyConfig(config NetworkConfig) error
	ValidateConfig(config NetworkConfig) error
	BackupConfig() error
	GetActiveStack() (string, error)

	// Legacy/Read-only
	GetHostname() (string, error)
	GetDNSConfig() ([]string, error)

	// Diagnostics Tools
	GetSocketStats() ([]SocketStat, error)
	GetNftablesRules() ([]string, error)
	RunTraceRoute(target string) (string, error)
	RunDig(target string) (string, error)
	RunNmap(target string, options string) ([]PortScanResult, error) // Changed signature
	RunTcpdump(interfaceName string, filter string, durationSeconds int) (string, error)
}
