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
}

func (m *LinuxAuditManager) queryJournal(filters []string, limit int) ([]ports.AuditEntry, error) {
	args := append([]string{"journalctl", "-o", "json", "-n", strconv.Itoa(limit)}, filters...)
	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("journalctl failed: %w", err)
	}

	var entries []ports.AuditEntry
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var entry JournalEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}

		ts, _ := strconv.ParseInt(entry.Realtime, 10, 64)
		// Realtime is in microseconds
		timestamp := time.Unix(0, ts*1000)

		auditEntry := ports.AuditEntry{
			Timestamp: timestamp,
			Details:   entry.Message,
			Severity:  "Info",
		}

		// Basic parsing logic based on message content
		if strings.Contains(strings.ToLower(entry.Message), "failed") || strings.Contains(strings.ToLower(entry.Message), "invalid") {
			auditEntry.Severity = "High"
			auditEntry.Result = "Fail"
		} else {
			auditEntry.Result = "Success"
		}

		// Extract user if possible (simple heuristic)
		if strings.Contains(entry.Message, "user ") {
			parts := strings.Fields(entry.Message)
			for i, p := range parts {
				if p == "user" && i+1 < len(parts) {
					auditEntry.User = parts[i+1]
					break
				}
			}
		}

		entries = append(entries, auditEntry)
	}
	return entries, nil
}

func (m *LinuxAuditManager) GetSSHAudit(limit int) ([]ports.AuditEntry, error) {
	return m.queryJournal([]string{"-u", "ssh", "-u", "sshd"}, limit)
}

func (m *LinuxAuditManager) GetSudoAudit(limit int) ([]ports.AuditEntry, error) {
	return m.queryJournal([]string{"-t", "sudo"}, limit)
}

func (m *LinuxAuditManager) CheckCriticalFiles() ([]ports.AuditEntry, error) {
	// Check /etc/passwd, /etc/shadow permissions
	files := []string{"/etc/passwd", "/etc/shadow", "/etc/sudoers"}
	var entries []ports.AuditEntry

	for _, f := range files {
		cmd := exec.Command("stat", "-c", "%a %U %G", f)
		out, err := cmd.Output()
		if err != nil {
			continue
		}

		perms := strings.TrimSpace(string(out))
		// Example check: 644 root root
		if f == "/etc/shadow" && !strings.HasPrefix(perms, "640") && !strings.HasPrefix(perms, "600") {
			entries = append(entries, ports.AuditEntry{
				Timestamp: time.Now(),
				Severity:  "High",
				Details:   fmt.Sprintf("File %s has unsafe permissions: %s", f, perms),
				Result:    "Fail",
			})
		}
	}
	return entries, nil
}

func (m *LinuxAuditManager) CheckSUIDBinaries() ([]string, error) {
	// find / -perm -4000
	cmd := exec.Command("find", "/bin", "/usr/bin", "-type", "f", "-perm", "-4000")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimSpace(string(out)), "\n"), nil
}
