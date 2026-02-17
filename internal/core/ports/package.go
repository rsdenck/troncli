package ports

// PackageInfo represents a software package
type PackageInfo struct {
	Name        string
	Version     string
	Description string
	Installed   bool
	Manager     string // apt, dnf, etc.
}

// PackageManager defines the interface for universal package operations
type PackageManager interface {
	// DetectManager returns the detected package manager name (apt, dnf, etc.)
	DetectManager() (string, error)

	// Install installs a package
	Install(packageName string) error

	// Remove removes a package
	Remove(packageName string) error

	// Update updates the package list and upgrades the system
	Update() error

	// Search searches for a package
	Search(query string) ([]PackageInfo, error)
}
