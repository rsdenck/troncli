package audit

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"
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
	// last -n 50 -F
	res, err := m.executor.Exec(ctx, "last", "-n", "50", "-F")
	if err != nil {
		return nil, err
	}

	var events []ports.AuditEvent
	lines := strings.Split(res.Stdout, "\n")
	for _, line := range lines {
		if line == "" || strings.HasPrefix(line, "wtmp begins") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 10 {
			// user pts/0 192.168.1.10 Fri Feb 17 10:00:00 2026
			user := parts[0]
			ip := parts[2]
			
			// Simple check if within duration (approximate for now without complex date parsing)
			// For a real implementation, we would parse the date.
			// Given the complexity of `last` date format and "since", we rely on `last` output mostly.
			
			events = append(events, ports.AuditEvent{
				Type:     "LOGIN_SUCCESS",
				User:     user,
				IP:       ip,
				Message:  line,
				Severity: "INFO",
				Timestamp: time.Now(), // Placeholder
			})
		}
	}
	return events, nil
}

func (m *UniversalAuditManager) parseLogs(output string, logType string) []ports.AuditEvent {
	var events []ports.AuditEvent
	scanner := bufio.NewScanner(strings.NewReader(output))
	
	// Regex for SSH failures
	// Feb 17 10:00:00 host sshd[123]: Failed password for invalid user admin from 192.168.1.50 port 55555 ssh2
	reSSHFail := regexp.MustCompile(`Failed password for (?:invalid user )?(\S+) from (\S+)`)
	
	// Regex for Sudo
	// Feb 17 10:00:00 host sudo:    user : TTY=pts/0 ; PWD=/home/user ; USER=root ; COMMAND=/bin/ls
	reSudo := regexp.MustCompile(`sudo:\s+(\S+)\s*:.*COMMAND=(.*)`)

	for scanner.Scan() {
		line := scanner.Text()
		
		if logType == "sshd" {
			matches := reSSHFail.FindStringSubmatch(line)
			if len(matches) >= 3 {
				events = append(events, ports.AuditEvent{
					Type:      "SSH_FAILURE",
					User:      matches[1],
					IP:        matches[2],
					Message:   line,
					Severity:  "WARNING",
					Timestamp: time.Now(), // Placeholder
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
					Timestamp: time.Now(), // Placeholder
				})
			}
		}
	}
	return events
}

// AuditUsers checks for user issues
func (m *UniversalAuditManager) AuditUsers() ([]ports.UserAudit, error) {
	// Parse /etc/passwd and /etc/shadow (requires root)
	// Using generic "cat" via executor
	ctx := context.Background()
	
	// Passwd
	res, err := m.executor.Exec(ctx, "cat", "/etc/passwd")
	if err != nil {
		return nil, err
	}
	passwdLines := strings.Split(res.Stdout, "\n")
	
	// Shadow (might fail if not root)
	resShadow, _ := m.executor.Exec(ctx, "cat", "/etc/shadow")
	shadowMap := make(map[string]string) // user -> hash
	if resShadow != nil {
		for _, line := range strings.Split(resShadow.Stdout, "\n") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				shadowMap[parts[0]] = parts[1]
			}
		}
	}

	var users []ports.UserAudit
	uidMap := make(map[string]bool)

	for _, line := range passwdLines {
		parts := strings.Split(line, ":")
		if len(parts) < 7 { continue }
		
		username := parts[0]
		uid := parts[2]
		
		audit := ports.UserAudit{
			Username: username,
			UID:      uid,
			GID:      parts[3],
			Home:     parts[5],
			Shell:    parts[6],
		}

		// Check duplicate UID
		if uidMap[uid] {
			audit.Message = "Duplicate UID"
		}
		uidMap[uid] = true

		// Check Shadow
		if hash, ok := shadowMap[username]; ok {
			if hash == "" || hash == "*" || hash == "!" {
				audit.PasswordStatus = "Locked/Empty"
			} else {
				audit.PasswordStatus = "Valid"
			}
			// Check expiration fields in shadow would require more parsing of parts[2] etc.
		}

		// Check SSH keys
		sshDir := fmt.Sprintf("%s/.ssh", audit.Home)
		if _, err := os.Stat(sshDir); err == nil {
			// Check permissions (should be 700)
			// stat -c %a
			res, _ := m.executor.Exec(ctx, "stat", "-c", "%a", sshDir)
			if res != nil && strings.TrimSpace(res.Stdout) != "700" {
				audit.SSHPermissions = "Unsafe (" + strings.TrimSpace(res.Stdout) + ")"
			} else {
				audit.SSHPermissions = "Safe"
			}
			
			// Check authorized_keys for invalid keys (basic check)
			authKeys := fmt.Sprintf("%s/authorized_keys", sshDir)
			if _, err := os.Stat(authKeys); err == nil {
				resKeys, _ := m.executor.Exec(ctx, "cat", authKeys)
				if resKeys != nil {
					if strings.Contains(resKeys.Stdout, "ssh-dss") { // Deprecated
						audit.SSHKeysInvalid = true
						audit.Message += " Deprecated DSA key found;"
					}
				}
			}
		}

		users = append(users, audit)
	}
	return users, nil
}

func (m *UniversalAuditManager) CheckPrivilegedGroups() ([]string, error) {
	ctx := context.Background()
	res, err := m.executor.Exec(ctx, "cat", "/etc/group")
	if err != nil {
		return nil, err
	}

	var privilegedUsers []string
	targetGroups := map[string]bool{"sudo": true, "wheel": true, "adm": true, "root": true}

	lines := strings.Split(res.Stdout, "\n")
	for _, line := range lines {
		// sudo:x:27:user1,user2
		parts := strings.Split(line, ":")
		if len(parts) >= 4 {
			groupName := parts[0]
			if targetGroups[groupName] {
				users := parts[3]
				if users != "" {
					userList := strings.Split(users, ",")
					for _, u := range userList {
						privilegedUsers = append(privilegedUsers, fmt.Sprintf("%s (%s)", u, groupName))
					}
				}
			}
		}
	}
	return privilegedUsers, nil
}

func (m *UniversalAuditManager) CheckBashCompatibility() ([]string, error) {
	var issues []string
	ctx := context.Background()

	// Get home directories from /etc/passwd
	res, err := m.executor.Exec(ctx, "cat", "/etc/passwd")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(res.Stdout, "\n")
	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) < 6 {
			continue
		}
		user := parts[0]
		home := parts[5]
		shell := parts[6]

		// Only check users with valid shells (bash/sh/zsh)
		if !strings.Contains(shell, "sh") || strings.Contains(shell, "nologin") || strings.Contains(shell, "false") {
			continue
		}

		filesToCheck := []string{".bashrc", ".profile", ".bash_profile", ".zshrc"}
		for _, f := range filesToCheck {
			path := fmt.Sprintf("%s/%s", home, f)
			// Read file (using cat via executor)
			// Note: This might fail if permission denied, but we try.
			contentRes, err := m.executor.Exec(ctx, "cat", path)
			if err == nil {
				issues = append(issues, m.analyzeShellConfig(user, path, contentRes.Stdout)...)
			}
		}
	}
	return issues, nil
}

func (m *UniversalAuditManager) analyzeShellConfig(user, path, content string) []string {
	var issues []string
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			continue
		}

		// Check for dangerous aliases
		if strings.Contains(line, "alias sudo") {
			issues = append(issues, fmt.Sprintf("[%s] %s: Dangerous alias 'sudo' found on line %d", user, path, i+1))
		}
		
		// Check for unsafe PATH
		if strings.Contains(line, "export PATH=.:") || strings.Contains(line, "export PATH=:") || strings.Contains(line, "::") {
			issues = append(issues, fmt.Sprintf("[%s] %s: Unsafe PATH (current directory) found on line %d", user, path, i+1))
		}

		// Check for suspicious permissions changes
		if strings.Contains(line, "chmod 777") || strings.Contains(line, "chmod -R 777") {
			issues = append(issues, fmt.Sprintf("[%s] %s: Suspicious 'chmod 777' found on line %d", user, path, i+1))
		}
		
		// Check for rm -rf / (unlikely but fatal)
		if strings.Contains(line, "rm -rf /") {
			issues = append(issues, fmt.Sprintf("[%s] %s: CRITICAL 'rm -rf /' found on line %d", user, path, i+1))
		}
	}
	return issues
}
