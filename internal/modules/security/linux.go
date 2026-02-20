package security

import (
	"github.com/mascli/troncli/internal/core/ports"
)

type LinuxSecurityManager struct {
	scanner *NativeScanner
}

func NewLinuxSecurityManager() ports.SecurityManager {
	return &LinuxSecurityManager{
		scanner: NewNativeScanner(),
	}
}

func (m *LinuxSecurityManager) IsToolInstalled() bool {
	// Native scanner is always "installed" as it's part of the binary
	return true
}

func (m *LinuxSecurityManager) InstallTool() error {
	// No installation needed for native scanner
	return nil
}

func (m *LinuxSecurityManager) ScanDirectory(path string) ([]ports.CVEVulnerability, error) {
	return m.scanner.ScanDirectory(path)
}

func (m *LinuxSecurityManager) ScanBinary(path string) ([]ports.CVEVulnerability, error) {
	return m.scanner.ScanFile(path)
}
