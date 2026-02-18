package audit

import (
"context"
"fmt"
"testing"
"time"

"github.com/mascli/troncli/internal/core/adapter"
"github.com/mascli/troncli/internal/core/domain"
)

// MockExecutor implements adapter.Executor for testing
type MockExecutor struct {
ExecFunc func(ctx context.Context, command string, args ...string) (*adapter.CommandResult, error)
}

func (m *MockExecutor) Exec(ctx context.Context, command string, args ...string) (*adapter.CommandResult, error) {
if m.ExecFunc != nil {
return m.ExecFunc(ctx, command, args...)
}
return nil, fmt.Errorf("unexpected call to Exec")
}

func (m *MockExecutor) ExecWithInput(ctx context.Context, input string, command string, args ...string) (*adapter.CommandResult, error) {
return m.Exec(ctx, command, args...)
}

func TestAnalyzeLogins(t *testing.T) {
mockExec := &MockExecutor{}
profile := &domain.SystemProfile{Distro: "ubuntu", InitSystem: "systemd"}
manager := NewUniversalAuditManager(mockExec, profile)

// Mock last output
// Format: User TTY IP Date...
lastOutput := `root     pts/0        192.168.1.10     Fri Feb 17 10:00:00 2026   still logged in
user     tty1         127.0.0.1        Thu Feb 16 09:00:00 2026 - down`

mockExec.ExecFunc = func(ctx context.Context, command string, args ...string) (*adapter.CommandResult, error) {
if command == "last" {
return &adapter.CommandResult{Stdout: lastOutput, ExitCode: 0}, nil
}
return nil, fmt.Errorf("unexpected command: %s", command)
}

// We set a long duration to include the mocked dates
events, err := manager.AnalyzeLogins(365 * 24 * time.Hour)
if err != nil {
t.Fatalf("AnalyzeLogins failed: %v", err)
}

if len(events) != 2 {
t.Errorf("Expected 2 events, got %d", len(events))
}

if events[0].User != "root" {
t.Errorf("Expected user root, got %s", events[0].User)
}
if events[0].IP != "192.168.1.10" {
t.Errorf("Expected IP 192.168.1.10, got %s", events[0].IP)
}
}

func TestAnalyzeSSH_Journald(t *testing.T) {
mockExec := &MockExecutor{}
profile := &domain.SystemProfile{Distro: "ubuntu", InitSystem: "systemd"}
manager := NewUniversalAuditManager(mockExec, profile)

// Mock journalctl JSON output
jsonOutput := `{"MESSAGE": "Accepted password for root from 192.168.1.20 port 54321 ssh2", "SYSLOG_IDENTIFIER": "sshd", "__REALTIME_TIMESTAMP": "1708164000000000", "_HOSTNAME": "server1", "_PID": "1234", "PRIORITY": "6"}
{"MESSAGE": "Failed password for invalid user admin from 192.168.1.50 port 55555 ssh2", "SYSLOG_IDENTIFIER": "sshd", "__REALTIME_TIMESTAMP": "1708164060000000", "_HOSTNAME": "server1", "_PID": "1235", "PRIORITY": "3"}`

mockExec.ExecFunc = func(ctx context.Context, command string, args ...string) (*adapter.CommandResult, error) {
if command == "journalctl" {
// Check args
hasOutputJson := false
for _, arg := range args {
if arg == "--output=json" {
hasOutputJson = true
}
}
if !hasOutputJson {
return nil, fmt.Errorf("expected --output=json arg")
}
return &adapter.CommandResult{Stdout: jsonOutput, ExitCode: 0}, nil
}
return nil, fmt.Errorf("unexpected command: %s", command)
}

events, err := manager.AnalyzeSSH(24 * time.Hour)
if err != nil {
t.Fatalf("AnalyzeSSH failed: %v", err)
}

if len(events) != 2 {
t.Errorf("Expected 2 events, got %d", len(events))
}

// Priority 3 is Critical
if events[1].Severity != "CRITICAL" {
t.Errorf("Expected severity CRITICAL for event 1, got %s", events[1].Severity)
}
}

func TestAnalyzeSudo_Journald(t *testing.T) {
mockExec := &MockExecutor{}
profile := &domain.SystemProfile{Distro: "ubuntu", InitSystem: "systemd"}
manager := NewUniversalAuditManager(mockExec, profile)

jsonOutput := `{"MESSAGE": "sudo: user : TTY=pts/0 ; PWD=/home/user ; USER=root ; COMMAND=/bin/ls", "SYSLOG_IDENTIFIER": "sudo", "__REALTIME_TIMESTAMP": "1708164000000000", "PRIORITY": "6"}`

mockExec.ExecFunc = func(ctx context.Context, command string, args ...string) (*adapter.CommandResult, error) {
if command == "journalctl" {
return &adapter.CommandResult{Stdout: jsonOutput, ExitCode: 0}, nil
}
return nil, fmt.Errorf("unexpected command: %s", command)
}

events, err := manager.AnalyzeSudo(24 * time.Hour)
if err != nil {
t.Fatalf("AnalyzeSudo failed: %v", err)
}

if len(events) != 1 {
t.Errorf("Expected 1 event, got %d", len(events))
}
}

func TestAnalyzeSSH_File(t *testing.T) {
mockExec := &MockExecutor{}
// InitSystem is empty, so it falls back to file
profile := &domain.SystemProfile{Distro: "ubuntu", InitSystem: ""}
manager := NewUniversalAuditManager(mockExec, profile)

// Mock file output from tail
// Use recent timestamps to pass 'since' check
ts1 := time.Now().Add(-1 * time.Minute).Format("Jan _2 15:04:05")
ts2 := time.Now().Add(-2 * time.Minute).Format("Jan _2 15:04:05")

fileOutput := fmt.Sprintf(`%s host sshd[123]: Failed password for invalid user admin from 192.168.1.50 port 55555 ssh2
%s host sshd[124]: Accepted password for root from 192.168.1.20 port 54321 ssh2`, ts1, ts2)

mockExec.ExecFunc = func(ctx context.Context, command string, args ...string) (*adapter.CommandResult, error) {
if command == "tail" {
return &adapter.CommandResult{Stdout: fileOutput, ExitCode: 0}, nil
}
return nil, fmt.Errorf("unexpected command: %s", command)
}

events, err := manager.AnalyzeSSH(24 * time.Hour)
if err != nil {
t.Fatalf("AnalyzeSSH failed: %v", err)
}

// Only 1 failure event is parsed by regex
if len(events) != 1 {
t.Errorf("Expected 1 event, got %d", len(events))
}

if events[0].Type != "SSH_FAILURE" {
t.Errorf("Expected SSH_FAILURE, got %s", events[0].Type)
}
if events[0].User != "admin" {
t.Errorf("Expected user admin, got %s", events[0].User)
}
}

func TestAnalyzeSudo_File(t *testing.T) {
mockExec := &MockExecutor{}
profile := &domain.SystemProfile{Distro: "ubuntu", InitSystem: ""}
manager := NewUniversalAuditManager(mockExec, profile)

ts := time.Now().Add(-1 * time.Minute).Format("Jan _2 15:04:05")
fileOutput := fmt.Sprintf(`%s host sudo:    user : TTY=pts/0 ; PWD=/home/user ; USER=root ; COMMAND=/bin/ls`, ts)

mockExec.ExecFunc = func(ctx context.Context, command string, args ...string) (*adapter.CommandResult, error) {
if command == "tail" {
return &adapter.CommandResult{Stdout: fileOutput, ExitCode: 0}, nil
}
return nil, fmt.Errorf("unexpected command: %s", command)
}

events, err := manager.AnalyzeSudo(24 * time.Hour)
if err != nil {
t.Fatalf("AnalyzeSudo failed: %v", err)
}

if len(events) != 1 {
t.Errorf("Expected 1 event, got %d", len(events))
}

if events[0].Type != "SUDO_ATTEMPT" {
t.Errorf("Expected SUDO_ATTEMPT, got %s", events[0].Type)
}
}

func TestAuditPlaceholders(t *testing.T) {
mockExec := &MockExecutor{}
profile := &domain.SystemProfile{}
manager := NewUniversalAuditManager(mockExec, profile)

if _, err := manager.AnalyzeFileChanges(nil, 0); err != nil {
t.Error("AnalyzeFileChanges failed")
}
if _, err := manager.AuditUsers(); err != nil {
t.Error("AuditUsers failed")
}
if _, err := manager.CheckPrivilegedGroups(); err != nil {
t.Error("CheckPrivilegedGroups failed")
}
if _, err := manager.CheckBashCompatibility(); err != nil {
t.Error("CheckBashCompatibility failed")
}
}

func TestAnalyzeLogins_Fallback(t *testing.T) {
	mockExec := &MockExecutor{}
	manager := NewUniversalAuditManager(mockExec, &domain.SystemProfile{})

	// Dynamic date within 24 hours
	ts := time.Now().Add(-1 * time.Hour).Format("Mon Jan _2 15:04:05 2006")
	// Ensure format matches "last" output exactly (it uses spaces padding)
	// But Time.Format handles padding for _2.
	// "Fri Feb 17 10:00:00 2026"

	mockExec.ExecFunc = func(ctx context.Context, command string, args ...string) (*adapter.CommandResult, error) {
		if command == "last" {
			// Check if -i is present
			hasI := false
			for _, arg := range args {
				if arg == "-i" {
					hasI = true
				}
			}
			if hasI {
				return nil, fmt.Errorf("last -i failed")
			}
			// Return output that matches regex to cover parsing lines
			// User TTY IP Date
			output := fmt.Sprintf("root     pts/0        192.168.1.10     %s   still logged in", ts)
			return &adapter.CommandResult{Stdout: output, ExitCode: 0}, nil
		}
		return nil, fmt.Errorf("unexpected command: %s", command)
	}

	events, err := manager.AnalyzeLogins(24 * time.Hour)
if err != nil {
t.Fatalf("AnalyzeLogins failed: %v", err)
}

if len(events) != 1 {
t.Errorf("Expected 1 event, got %d", len(events))
}
}

func TestAnalyzeJournal_Error(t *testing.T) {
mockExec := &MockExecutor{}
profile := &domain.SystemProfile{Distro: "ubuntu", InitSystem: "systemd"}
manager := NewUniversalAuditManager(mockExec, profile)

mockExec.ExecFunc = func(ctx context.Context, command string, args ...string) (*adapter.CommandResult, error) {
return nil, fmt.Errorf("journalctl error")
}

_, err := manager.AnalyzeSSH(24 * time.Hour)
if err == nil {
t.Error("Expected error from AnalyzeSSH when journalctl fails")
}
}

func TestAnalyzeLogFile_Error(t *testing.T) {
mockExec := &MockExecutor{}
profile := &domain.SystemProfile{Distro: "ubuntu", InitSystem: ""}
manager := NewUniversalAuditManager(mockExec, profile)

mockExec.ExecFunc = func(ctx context.Context, command string, args ...string) (*adapter.CommandResult, error) {
return nil, fmt.Errorf("tail error")
}

_, err := manager.AnalyzeSSH(24 * time.Hour)
if err == nil {
t.Error("Expected error from AnalyzeSSH when tail fails")
}
}

func TestAnalyzeJournal_EdgeCases(t *testing.T) {
mockExec := &MockExecutor{}
profile := &domain.SystemProfile{Distro: "ubuntu", InitSystem: "systemd"}
manager := NewUniversalAuditManager(mockExec, profile)

// 1. Empty line
// 2. Invalid JSON
// 3. Warning priority (4)
// 4. Info priority (>4)
jsonOutput := `
{"invalid": json}
{"MESSAGE": "Warning event", "SYSLOG_IDENTIFIER": "test", "__REALTIME_TIMESTAMP": "1708164000000000", "PRIORITY": "4"}
{"MESSAGE": "Info event", "SYSLOG_IDENTIFIER": "test", "__REALTIME_TIMESTAMP": "1708164000000000", "PRIORITY": "6"}`

mockExec.ExecFunc = func(ctx context.Context, command string, args ...string) (*adapter.CommandResult, error) {
return &adapter.CommandResult{Stdout: jsonOutput, ExitCode: 0}, nil
}

events, err := manager.AnalyzeSSH(24 * time.Hour) // AnalyzeSSH calls analyzeJournal
if err != nil {
t.Fatalf("AnalyzeSSH failed: %v", err)
}

if len(events) != 2 {
t.Errorf("Expected 2 events, got %d", len(events))
}

if events[0].Severity != "WARNING" {
t.Errorf("Expected WARNING severity, got %s", events[0].Severity)
}
if events[1].Severity != "INFO" {
t.Errorf("Expected INFO severity, got %s", events[1].Severity)
}
}

func TestAnalyzeLogins_OldEvent(t *testing.T) {
mockExec := &MockExecutor{}
manager := NewUniversalAuditManager(mockExec, &domain.SystemProfile{})

// Old date: 2020
output := `root     pts/0        192.168.1.10     Mon Jan 01 10:00:00 2020   still logged in`

mockExec.ExecFunc = func(ctx context.Context, command string, args ...string) (*adapter.CommandResult, error) {
return &adapter.CommandResult{Stdout: output, ExitCode: 0}, nil
}

events, err := manager.AnalyzeLogins(24 * time.Hour)
if err != nil {
t.Fatalf("AnalyzeLogins failed: %v", err)
}

if len(events) != 0 {
t.Errorf("Expected 0 events, got %d", len(events))
}
}

func TestAnalyzeLogins_AllFail(t *testing.T) {
mockExec := &MockExecutor{}
manager := NewUniversalAuditManager(mockExec, &domain.SystemProfile{})

mockExec.ExecFunc = func(ctx context.Context, command string, args ...string) (*adapter.CommandResult, error) {
return nil, fmt.Errorf("command failed")
}

_, err := manager.AnalyzeLogins(24 * time.Hour)
if err == nil {
t.Error("Expected error when all commands fail")
}
}

func TestAnalyzeSSH_File_NoMatch(t *testing.T) {
mockExec := &MockExecutor{}
profile := &domain.SystemProfile{Distro: "ubuntu", InitSystem: ""}
manager := NewUniversalAuditManager(mockExec, profile)

output := `Feb 17 10:00:00 host sshd[123]: Something unrelated`

mockExec.ExecFunc = func(ctx context.Context, command string, args ...string) (*adapter.CommandResult, error) {
return &adapter.CommandResult{Stdout: output, ExitCode: 0}, nil
}

events, err := manager.AnalyzeSSH(24 * time.Hour)
if err != nil {
t.Fatalf("AnalyzeSSH failed: %v", err)
}

if len(events) != 0 {
t.Errorf("Expected 0 events, got %d", len(events))
}
}

func TestAnalyzeLogins_MalformedOutput(t *testing.T) {
mockExec := &MockExecutor{}
manager := NewUniversalAuditManager(mockExec, &domain.SystemProfile{})

output := `garbage output line`

mockExec.ExecFunc = func(ctx context.Context, command string, args ...string) (*adapter.CommandResult, error) {
return &adapter.CommandResult{Stdout: output, ExitCode: 0}, nil
}

events, err := manager.AnalyzeLogins(24 * time.Hour)
if err != nil {
t.Fatalf("AnalyzeLogins failed: %v", err)
}

if len(events) != 0 {
t.Errorf("Expected 0 events, got %d", len(events))
}
}

func TestAnalyzeSSH_File_InvalidDate(t *testing.T) {
mockExec := &MockExecutor{}
profile := &domain.SystemProfile{Distro: "ubuntu", InitSystem: ""}
manager := NewUniversalAuditManager(mockExec, profile)

// Invalid date format
output := `InvalidDate host sshd[123]: Failed password for invalid user admin from 1.2.3.4 port 55555 ssh2`

mockExec.ExecFunc = func(ctx context.Context, command string, args ...string) (*adapter.CommandResult, error) {
return &adapter.CommandResult{Stdout: output, ExitCode: 0}, nil
}

events, err := manager.AnalyzeSSH(24 * time.Hour)
if err != nil {
t.Fatalf("AnalyzeSSH failed: %v", err)
}

if len(events) != 0 {
t.Errorf("Expected 0 events, got %d", len(events))
}
}
