package ports

// BashManager defines the interface for bash operations
type BashManager interface {
	// RunCommand executes a single bash command
	RunCommand(cmd string) (string, error)
	// RunScript executes a bash script from a file
	RunScript(path string) (string, error)
}
