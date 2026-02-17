//go:build !linux

package firewall

import (
	"errors"
	"github.com/mascli/troncli/internal/core/ports"
)

type OtherOSFirewallManager struct{}

func NewLinuxFirewallManager() ports.FirewallManager {
	return &OtherOSFirewallManager{}
}

func (m *OtherOSFirewallManager) DetectFirewall() (string, error) {
	return "", errors.New("firewall not supported on this OS")
}

func (m *OtherOSFirewallManager) AllowPort(port string, protocol string) error {
	return errors.New("firewall not supported on this OS")
}

func (m *OtherOSFirewallManager) BlockPort(port string, protocol string) error {
	return errors.New("firewall not supported on this OS")
}

func (m *OtherOSFirewallManager) ListRules() ([]ports.FirewallRule, error) {
	return nil, errors.New("firewall not supported on this OS")
}

func (m *OtherOSFirewallManager) Enable() error {
	return errors.New("firewall not supported on this OS")
}

func (m *OtherOSFirewallManager) Disable() error {
	return errors.New("firewall not supported on this OS")
}
