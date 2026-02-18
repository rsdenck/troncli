package audit

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/domain"
	"github.com/mascli/troncli/internal/core/ports"
)

// UniversalAuditManager implements ports.AuditManager
type UniversalAuditManager struct {
	executor adapter.Executor
	profile  *domain.SystemProfile
}

// NewUniversalAuditManager creates a new audit manager
func NewUniversalAuditManager(executor adapter.Executor, profile *domain.SystemProfile) *UniversalAuditManager {
	return &UniversalAuditManager{
		executor: executor,
		profile:  profile,
	}
}

// AnalyzeSSH analyzes SSH logs
func (m *UniversalAuditManager) AnalyzeSSH(since time.Duration) ([]ports.AuditEvent, error) {
	// Decide log source
	logFile := "/var/log/auth.log"
	if _, err := os.Stat("/var/log/secure"); err == nil {
		logFile = "/var/log/secure"
	}

	// If systemd, try journalctl first
	if m.profile.InitSystem == "systemd" {
		return m.analyzeJournal("sshd", since)
	}

	return m.analyzeLogFile(logFile, "sshd", since)
}

func (m *UniversalAuditManager) AnalyzeSudo(since time.Duration) ([]ports.AuditEvent, error) {
	if m.profile.InitSystem == "systemd" {
		return m.analyzeJournal("sudo", since)
	}
	return m.analyzeLogFile("/var/log/auth.log", "sudo", since)
}

func (m *UniversalAuditManager) AnalyzeLogins(since time.Duration) ([]ports.AuditEvent, error) {
	ctx := context.Background()
	// last -n 50 -F -i (Show IP always, Full date)
	res, err := m.executor.Exec(ctx, "last", "-n", "50", "-F", "-i")
	if err != nil {
		// Fallback without -i if it fails (some busybox/older versions might not support -i)
		res, err = m.executor.Exec(ctx, "last", "-n", "50", "-F")
		if err != nil {
			return nil, err
		}
	}

	// Regex for last -F -i
	// root     pts/0        0.0.0.0          Fri Feb 17 10:00:00 2026   still logged in
	// user     pts/0        192.168.1.50     Fri Feb 17 10:00:00 2026 - Fri Feb 17 ...
	// Group 1: User
	// Group 2: TTY
	// Group 3: IP
	// Group 4: Date start
	reLast := regexp.MustCompile(`^(\S+)\s+(\S+)\s+(\S+)\s+(\w+\s+\w+\s+\d+\s+\d+:\d+:\d+\s+\d+)`)

	var events []ports.AuditEvent
	lines := strings.Split(res.Stdout, "\n")
	for _, line := range lines {
		if line == "" || strings.HasPrefix(line, "wtmp begins") {
			continue
		}
		
		matches := reLast.FindStringSubmatch(line)
		if len(matches) >= 5 {
			user := matches[1]
			tty := matches[2]
			ip := matches[3]
			dateStr := matches[4]

			// Parse date
			// Mon Jan 02 15:04:05 2006
			timestamp, _ := time.Parse(time.UnixDate, dateStr) // approximate, actually `Fri Feb 17 10:00:00 2026` is close to UnixDate/RubyDate
			// Standard date format in `last -F`: "Fri Feb 17 10:00:00 2026"
			// Go layout: "Mon Jan _2 15:04:05 2006"
			if t, err := time.Parse("Mon Jan _2 15:04:05 2006", dateStr); err == nil {
				timestamp = t
			}

			if time.Since(timestamp) > since {
				continue
			}

			events = append(events, ports.AuditEvent{
				Type:      "LOGIN_SUCCESS",
				User:      user,
				IP:        ip,
				Message:   fmt.Sprintf("Login on %s from %s", tty, ip),
				Severity:  "INFO",
				Timestamp: timestamp,
			})
		}
	}
	return events, nil
}

func (m *UniversalAuditManager) parseLogs(output string, logType string, since time.Duration) []ports.AuditEvent {
	var events []ports.AuditEvent
	scanner := bufio.NewScanner(strings.NewReader(output))

	// Regex for SSH failures
	reSSHFail := regexp.MustCompile(`Failed password for (?:invalid user )?(\S+) from (\S+)`)
	
	// Regex for Sudo
	reSudo := regexp.MustCompile(`sudo:\s+(\S+)\s*:.*COMMAND=(.*)`)

	currentYear := time.Now().Year()
	cutoff := time.Now().Add(-since)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 16 {
			continue
		}

		// Parse syslog date: "Feb 17 10:00:00"
		// Go layout: "Jan _2 15:04:05"
		dateStr := line[:15]
		timestamp, err := time.Parse("Jan _2 15:04:05", dateStr)
		if err != nil {
			continue
		}
		
		// Add year (syslog usually doesn't have year)
		timestamp = timestamp.AddDate(currentYear, 0, 0)
		
		// Handle year wraparound (e.g. reading Dec logs in Jan)
		if timestamp.After(time.Now().Add(24 * time.Hour)) {
			timestamp = timestamp.AddDate(-1, 0, 0)
		}

		if timestamp.Before(cutoff) {
			continue
		}

		if logType == "sshd" {
			matches := reSSHFail.FindStringSubmatch(line)
			if len(matches) >= 3 {
				events = append(events, ports.AuditEvent{
					Type:      "SSH_FAILURE",
					User:      matches[1],
					IP:        matches[2],
					Message:   line,
					Severity:  "WARNING",
					Timestamp: timestamp,
				})
			}
		} else if logType == "sudo" {
			matches := reSudo.FindStringSubmatch(line)
			if len(matches) >= 3 {
				events = append(events, ports.AuditEvent{
					Type:      "SUDO_ATTEMPT",
					User:      matches[1],
					Message:   fmt.Sprintf("Command: %s", matches[2]),
					Severity:  "INFO",
					Timestamp: timestamp,
				})
			}
		}
	}
	return events
}

// JournalEntry represents a JSON entry from journalctl
type JournalEntry struct {
	Message          string `json:"MESSAGE"`
	SyslogIdentifier string `json:"SYSLOG_IDENTIFIER"`
	Realtime         string `json:"__REALTIME_TIMESTAMP"` // Microseconds since epoch
	Hostname         string `json:"_HOSTNAME"`
	PID              string `json:"_PID"`
	UID              string `json:"_UID"`
	GID              string `json:"_GID"`
	Comm             string `json:"_COMM"`
	Priority         string `json:"PRIORITY"` // 0-7 string
	SystemdUnit      string `json:"_SYSTEMD_UNIT"`
}

// AnalyzeFileChanges analyzes file changes in the given paths within the specified duration
func (m *UniversalAuditManager) AnalyzeFileChanges(paths []string, since time.Duration) ([]ports.AuditEvent, error) {
	// TODO: Implement file change analysis using auditd or fsnotify
	// For now, return empty list to satisfy interface
	return []ports.AuditEvent{}, nil
}

// AuditUsers checks user accounts for security issues
func (m *UniversalAuditManager) AuditUsers() ([]ports.UserAudit, error) {
	// TODO: Implement user auditing (check for empty passwords, uid 0, etc)
	return []ports.UserAudit{}, nil
}

// CheckPrivilegedGroups checks for users in privileged groups
func (m *UniversalAuditManager) CheckPrivilegedGroups() ([]string, error) {
	// TODO: Implement group auditing
	return []string{}, nil
}

// CheckBashCompatibility checks for bash configuration issues
func (m *UniversalAuditManager) CheckBashCompatibility() ([]string, error) {
	// TODO: Implement bash checks
	return []string{}, nil
}

// analyzeJournal queries journalctl for logs
func (m *UniversalAuditManager) analyzeJournal(service string, since time.Duration) ([]ports.AuditEvent, error) {
	ctx := context.Background()
	sinceStr := time.Now().Add(-since).Format("2006-01-02 15:04:05")

	// Try JSON output first
	// journalctl -u service --since "X" --output=json --no-pager
	res, err := m.executor.Exec(ctx, "journalctl", "-u", service, "--since", sinceStr, "--output=json", "--no-pager")
	if err == nil {
		var events []ports.AuditEvent
		// Parse JSON stream (one JSON object per line)
		scanner := bufio.NewScanner(strings.NewReader(res.Stdout))
		for scanner.Scan() {
			line := scanner.Text()
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
			prio, _ := strconv.Atoi(entry.Priority)
			if prio <= 3 { // 0=Emerg, 1=Alert, 2=Crit, 3=Err
				severity = "CRITICAL"
			} else if prio == 4 { // Warning
				severity = "WARNING"
			}

			events = append(events, ports.AuditEvent{
				Type:      entry.SyslogIdentifier,
				User:      entry.UID, // Use UID if available
				Timestamp: ts,
				Message:   entry.Message,
				Severity:  severity,
			})
		}
		return events, nil
	}
	
	return nil, fmt.Errorf("journalctl failed or not available")
}

// analyzeLogFile reads log file and parses it
func (m *UniversalAuditManager) analyzeLogFile(path string, logType string, since time.Duration) ([]ports.AuditEvent, error) {
	ctx := context.Background()
	// tail -n 500 path
	res, err := m.executor.Exec(ctx, "tail", "-n", "500", path)
	if err != nil {
		return nil, err
	}
	
	return m.parseLogs(res.Stdout, logType, since), nil
}
