package ports

// SSHClient defines the interface for SSH operations
type SSHClient interface {
	// ListProfiles returns a list of available SSH profiles
	ListProfiles() ([]string, error)

	// Connect establishes a connection to a profile (interactive)
	Connect(profile string) error

	// Execute runs a command on the remote host and returns the output
	Execute(profile string, command string) (string, error)

	// CopyFile copies a file to the remote host
	CopyFile(profile string, src string, dest string) error
}
