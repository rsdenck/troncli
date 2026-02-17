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

func (m *LinuxAuditManager) queryJournal(filters []string, limit int) ([]ports.AuditEntry, error) {
	// Add JSON output format and reverse order (newest first)
	args := append([]string{"journalctl", "-o", "json", "-r", "-n", strconv.Itoa(limit)}, filters...)
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
			Message:   entry.Message,
			Service:   entry.SyslogIdentifier,
			Host:      entry.Hostname,
			Severity:  "Info",
		}

		// Basic parsing logic based on message content
		lowerMsg := strings.ToLower(entry.Message)
		if strings.Contains(lowerMsg, "failed") || 
		   strings.Contains(lowerMsg, "invalid") || 
		   strings.Contains(lowerMsg, "error") || 
		   strings.Contains(lowerMsg, "denied") {
			auditEntry.Severity = "High"
			auditEntry.Result = "Fail"
		} else {
			auditEntry.Result = "Success"
		}

		// Extract user if possible (simple heuristic)
		if strings.Contains(entry.Message, "user=") {
			// Extract user=name
			start := strings.Index(entry.Message, "user=") + 5
			end := strings.Index(entry.Message[start:], " ")
			if end == -1 {
				end = len(entry.Message[start:])
			}
			auditEntry.User = entry.Message[start : start+end]
		} else if strings.Contains(entry.Message, "user ") {
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

func (m *LinuxAuditManager) GetAuthLogs(limit int) ([]ports.AuditEntry, error) {
	// auth.log is usually covered by syslog identifier 'sshd', 'sudo', 'login', 'su'
	// Using journalctl facility 'auth' and 'authpriv'
	return m.queryJournal([]string{"SYSLOG_FACILITY=4", "SYSLOG_FACILITY=10"}, limit)
}

func (m *LinuxAuditManager) GetJournalLogs(service string, limit int) ([]ports.AuditEntry, error) {
	return m.queryJournal([]string{"-u", service}, limit)
}

func (m *LinuxAuditManager) GetPAMTrace(limit int) ([]ports.AuditEntry, error) {
	// PAM logs often come from various services, but we can grep for "pam" in message
	// or look for specific PAM libs.
	// journalctl -g "pam_" might be slow, so we rely on auth facility and filter in code?
	// Or use _TRANSPORT=syslog
	return m.queryJournal([]string{"SYSLOG_FACILITY=4", "SYSLOG_FACILITY=10"}, limit)
}

func (m *LinuxAuditManager) CheckCriticalFiles() ([]ports.AuditEntry, error) {
	// Check /etc/passwd, /etc/shadow permissions
	files := []string{"/etc/passwd", "/etc/shadow", "/etc/sudoers", "/etc/ssh/sshd_config"}
	var entries []ports.AuditEntry

	for _, f := range files {
		cmd := exec.Command("stat", "-c", "%a %U %G", f)
		out, err := cmd.Output()
		if err != nil {
			continue
		}

		perms := strings.TrimSpace(string(out))
		// Example check: 644 root root
		// Shadow must be 640 or 600
		if f == "/etc/shadow" && !strings.HasPrefix(perms, "640") && !strings.HasPrefix(perms, "600") {
			entries = append(entries, ports.AuditEntry{
				Timestamp: time.Now(),
				Severity:  "High",
				Message:   fmt.Sprintf("File %s has unsafe permissions: %s", f, perms),
				Result:    "Fail",
			})
		}
		// Sudoers must be 440
		if f == "/etc/sudoers" && !strings.HasPrefix(perms, "440") {
			entries = append(entries, ports.AuditEntry{
				Timestamp: time.Now(),
				Severity:  "High",
				Message:   fmt.Sprintf("File %s has unsafe permissions: %s", f, perms),
				Result:    "Fail",
			})
		}
	}
	return entries, nil
}

func (m *LinuxAuditManager) CheckSUIDBinaries() ([]string, error) {
	// find / -perm -4000
	// Limit to /bin /usr/bin /sbin /usr/sbin to avoid full disk scan in interactive mode
	cmd := exec.Command("find", "/bin", "/usr/bin", "/sbin", "/usr/sbin", "-type", "f", "-perm", "-4000")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimSpace(string(out)), "\n"), nil
}

func (m *LinuxAuditManager) CheckListeningPorts() ([]string, error) {
	// ss -tulpn
	cmd := exec.Command("ss", "-tulpn")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimSpace(string(out)), "\n"), nil
}
