package audit

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/mascli/troncli/internal/core/ports"
)

// LinuxAuditManager implements ports.AuditManager using journalctl and file checks
type LinuxAuditManager struct{}

func NewLinuxAuditManager() ports.AuditManager {
	return &LinuxAuditManager{}
}

// JournalEntry represents a JSON entry from journalctl
type JournalEntry struct {
	Message          string `json:"MESSAGE"`
	SyslogIdentifier string `json:"SYSLOG_IDENTIFIER"`
	Realtime         string `json:"__REALTIME_TIMESTAMP"`
	Hostname         string `json:"_HOSTNAME"`
	PID              string `json:"_PID"`
	SystemdUnit      string `json:"_SYSTEMD_UNIT"`
}

func (m *LinuxAuditManager) queryJournal(filters []string, limit int) ([]ports.AuditEvent, error) {
	// Add JSON output format and reverse order (newest first)
	args := append([]string{"journalctl", "-o", "json", "-r", "-n", strconv.Itoa(limit)}, filters...)
	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("journalctl failed: %w", err)
	}

	var entries []ports.AuditEvent
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var entry JournalEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}

		// Convert to AuditEvent
		tsVal, _ := strconv.ParseInt(entry.Realtime, 10, 64)
		ts := time.Unix(0, tsVal*1000) // Microseconds to Nanoseconds

		severity := "INFO"
		// Simple heuristic
		if strings.Contains(strings.ToLower(entry.Message), "failed") || strings.Contains(strings.ToLower(entry.Message), "failure") {
			severity = "WARNING"
		}
		if strings.Contains(strings.ToLower(entry.Message), "root") && severity == "WARNING" {
			severity = "CRITICAL"
		}

		entries = append(entries, ports.AuditEvent{
			Type:      entry.SyslogIdentifier,
			User:      "unknown", // Need parsing from message
			Timestamp: ts,
			Message:   entry.Message,
			Severity:  severity,
		})
	}
	return entries, nil
}

func (m *LinuxAuditManager) AnalyzeSSH(since time.Duration) ([]ports.AuditEvent, error) {
	return m.queryJournal([]string{"_COMM=sshd"}, 100)
}

func (m *LinuxAuditManager) AnalyzeSudo(since time.Duration) ([]ports.AuditEvent, error) {
	return m.queryJournal([]string{"_COMM=sudo"}, 100)
}

func (m *LinuxAuditManager) AnalyzeLogins(since time.Duration) ([]ports.AuditEvent, error) {
	// Combine SSH and Sudo + local login
	ssh, _ := m.AnalyzeSSH(since)
	sudo, _ := m.AnalyzeSudo(since)
	return append(ssh, sudo...), nil
}

func (m *LinuxAuditManager) AuditUsers() ([]ports.UserAudit, error) {
	// Parse /etc/shadow or use chage
	return []ports.UserAudit{}, nil // TODO: Implement
}

func (m *LinuxAuditManager) CheckPrivilegedGroups() ([]string, error) {
	// grep sudo /etc/group
	return []string{}, nil // TODO: Implement
}

func (m *LinuxAuditManager) CheckBashCompatibility() ([]string, error) {
	return []string{}, nil // TODO: Implement
}

func (m *LinuxAuditManager) AnalyzeFileChanges(paths []string, since time.Duration) ([]ports.AuditEvent, error) {
	return []ports.AuditEvent{}, nil // TODO: Implement
}
