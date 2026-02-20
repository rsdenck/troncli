package ports

import "time"

// AuditEvent represents a security event
type AuditEvent struct {
	Type      string // SSH_FAILURE, SUDO_ATTEMPT, LOGIN_SUCCESS
	User      string
	IP        string
	Timestamp time.Time
	Message   string
	Severity  string // INFO, WARNING, CRITICAL
}

// UserAudit represents user security audit
type UserAudit struct {
	Username       string
	UID            string
	GID            string
	Home           string
	Shell          string
	PasswordStatus string // Empty, Expired, Valid
	LastChange     string
	SSHPermissions string // Safe, Unsafe
	SSHKeysInvalid bool
	Message        string // Security findings
}

// AuditManager defines interface for security auditing
type AuditManager interface {
	// Log Analysis
	AnalyzeSSH(since time.Duration) ([]AuditEvent, error)
	AnalyzeSudo(since time.Duration) ([]AuditEvent, error)
	AnalyzeLogins(since time.Duration) ([]AuditEvent, error)
	AnalyzeFileChanges(paths []string, since time.Duration) ([]AuditEvent, error)

	// User & Permissions
	AuditUsers() ([]UserAudit, error)
	CheckPrivilegedGroups() ([]string, error)

	// Environment
	CheckBashCompatibility() ([]string, error) // .bashrc checks etc.
}
