package ports

// Container defines the structure for a container
type Container struct {
	ID      string
	Names   []string
	Image   string
	State   string
	Status  string
	Runtime string // "docker" or "podman"
}

// ContainerManager defines the interface for container operations
type ContainerManager interface {
	// ListContainers returns a list of containers from all available runtimes
	ListContainers(all bool) ([]Container, error)

	// StartContainer starts a container
	StartContainer(id string) error

	// StopContainer stops a container
	StopContainer(id string) error

	// RestartContainer restarts a container
	RestartContainer(id string) error

	// RemoveContainer removes a container
	RemoveContainer(id string, force bool) error

	// GetContainerLogs returns the logs of a container
	GetContainerLogs(id string, tail int) (string, error)

	// PruneSystem removes unused data
	PruneSystem() (string, error)
}
