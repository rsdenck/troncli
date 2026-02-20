package ports

// CVE Vulnerability represents a found vulnerability
type CVEVulnerability struct {
	CVEID       string
	Severity    string
	Product     string
	Version     string
	Description string
	Path        string
}

// SecurityManager defines the interface for security and CVE operations
type SecurityManager interface {
	// ScanDirectory runs cve-bin-tool on a directory
	ScanDirectory(path string) ([]CVEVulnerability, error)

	// ScanBinary runs cve-bin-tool on a specific binary
	ScanBinary(path string) ([]CVEVulnerability, error)

	// IsToolInstalled checks if cve-bin-tool is available
	IsToolInstalled() bool

	// InstallTool attempts to install cve-bin-tool via pip
	InstallTool() error
}
