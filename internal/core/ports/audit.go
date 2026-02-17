package ports

import "time"

// AuditEntry represents an audit log entry
type AuditEntry struct {
	Timestamp time.Time
	User      string
	Command   string
	Result    string // Success/Fail
	Severity  string // High, Medium, Low
	Details   string
}

// AuditManager defines operations for system auditing
type AuditManager interface {
	GetSSHAudit(limit int) ([]AuditEntry, error)
	GetSudoAudit(limit int) ([]AuditEntry, error)
	CheckCriticalFiles() ([]AuditEntry, error)
	CheckSUIDBinaries() ([]string, error)
}
