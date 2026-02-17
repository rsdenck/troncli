package ports

// ProcessManager defines operations for managing system processes
type ProcessManager interface {
	// KillProcess sends a signal to a process
	// signal can be "SIGTERM", "SIGKILL", etc.
	KillProcess(pid int, signal string) error

	// ReniceProcess changes the priority of a process
	// priority range is usually -20 (highest) to 19 (lowest)
	ReniceProcess(pid int, priority int) error

	// KillZombies identifies zombie processes and attempts to eliminate them
	// Returns the number of zombies found and eliminated (or attempted)
	KillZombies() (int, error)

	// Phase 4: Universal Process Features
	GetProcessTree() ([]ProcessNode, error)
	GetOpenFiles(pid int) ([]string, error)
	GetProcessPorts(pid int) ([]string, error)
}

// ProcessNode represents a process in a tree structure
type ProcessNode struct {
	PID      int
	PPID     int
	Name     string
	User     string
	State    string
	Children []ProcessNode
}
