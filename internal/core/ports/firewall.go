package ports

// FirewallRule represents a firewall rule
type FirewallRule struct {
	ID       string
	Action   string // ALLOW, BLOCK
	Protocol string // tcp, udp
	Port     string // 80, 443, etc.
	Source   string // IP or ANY
	Comment  string
}

// FirewallManager defines the interface for universal firewall operations
type FirewallManager interface {
	// DetectFirewall returns the detected firewall manager (nftables, ufw, etc.)
	DetectFirewall() (string, error)

	// AllowPort allows incoming traffic on a port
	AllowPort(port string, protocol string) error

	// BlockPort blocks incoming traffic on a port
	BlockPort(port string, protocol string) error

	// ListRules returns the current rules
	ListRules() ([]FirewallRule, error)

	// Enable enables the firewall
	Enable() error

	// Disable disables the firewall
	Disable() error
}
