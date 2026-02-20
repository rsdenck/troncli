package ports

// ServiceUnit represents a system service (universal)
type ServiceUnit struct {
	Name        string
	Status      string // active, inactive, failed
	Enabled     bool
	PID         int
	Description string
	LoadState   string // systemd specific, optional
	ActiveState string // systemd specific, optional
	SubState    string // systemd specific, optional
}

// ServiceManager defines the interface for universal service operations
type ServiceManager interface {
	// ListServices returns a list of system services
	ListServices() ([]ServiceUnit, error)

	// StartService starts a service
	StartService(name string) error

	// StopService stops a service
	StopService(name string) error

	// RestartService restarts a service
	RestartService(name string) error

	// EnableService enables a service to start on boot
	EnableService(name string) error

	// DisableService disables a service from starting on boot
	DisableService(name string) error

	// GetServiceStatus returns the status output of a service
	GetServiceStatus(name string) (string, error)

	// GetServiceJournal returns the logs for a service (journald or log file)
	GetServiceLogs(name string, lines int) (string, error)
}
