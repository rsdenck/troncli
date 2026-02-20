package service

import (
	"context"
	"fmt"
	"testing"

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

func TestListSystemdServices_JSON(t *testing.T) {
	mockExec := &MockExecutor{}
	profile := &domain.SystemProfile{Distro: "ubuntu", InitSystem: "systemd"}
	manager := NewUniversalServiceManager(mockExec, profile)

	// Mock JSON output from systemctl
	jsonOutput := `[
		{"unit": "ssh.service", "load": "loaded", "active": "active", "sub": "running", "description": "OpenBSD Secure Shell server"},
		{"unit": "nginx.service", "load": "loaded", "active": "inactive", "sub": "dead", "description": "A high performance web server and a reverse proxy server"}
	]`

	mockExec.ExecFunc = func(ctx context.Context, command string, args ...string) (*adapter.CommandResult, error) {
		// Check command
		if command != "systemctl" {
			return nil, fmt.Errorf("unexpected command: %s", command)
		}
		return &adapter.CommandResult{
			Stdout:   jsonOutput,
			ExitCode: 0,
		}, nil
	}

	services, err := manager.ListServices()
	if err != nil {
		t.Fatalf("ListServices failed: %v", err)
	}

	if len(services) != 2 {
		t.Errorf("Expected 2 services, got %d", len(services))
	}

	if services[0].Name != "ssh.service" || services[0].ActiveState != "active" {
		t.Errorf("Service 0 mismatch: %+v", services[0])
	}
	// JSON mock has "inactive" for nginx
	if services[1].Name != "nginx.service" {
		t.Errorf("Service 1 name mismatch: got %s, want nginx.service", services[1].Name)
	}
	if services[1].ActiveState != "inactive" {
		t.Errorf("Service 1 state mismatch: got %s, want inactive", services[1].ActiveState)
	}
}

func TestListSystemdServices_TextFallback(t *testing.T) {
	mockExec := &MockExecutor{}
	profile := &domain.SystemProfile{Distro: "ubuntu", InitSystem: "systemd"}
	manager := NewUniversalServiceManager(mockExec, profile)

	mockExec.ExecFunc = func(ctx context.Context, command string, args ...string) (*adapter.CommandResult, error) {
		// Check if it's JSON request
		isJson := false
		for _, arg := range args {
			if arg == "--output=json" {
				isJson = true
				break
			}
		}

		if isJson {
			// Simulate JSON not supported or failing
			return nil, fmt.Errorf("json not supported")
		}

		// Return Text output for fallback
		textOutput := `ssh.service     loaded active running OpenBSD Secure Shell server
nginx.service   loaded active exited  A high performance web server`
		
		return &adapter.CommandResult{
			Stdout:   textOutput,
			ExitCode: 0,
		}, nil
	}

	services, err := manager.ListServices()
	if err != nil {
		t.Fatalf("ListServices failed: %v", err)
	}

	if len(services) != 2 {
		t.Errorf("Expected 2 services, got %d", len(services))
	}

	if services[0].Name != "ssh.service" || services[0].ActiveState != "active" {
		t.Errorf("Service 0 mismatch: %+v", services[0])
	}
	// Text mock has "active" for nginx
	if services[1].Name != "nginx.service" {
		t.Errorf("Service 1 name mismatch: got %s, want nginx.service", services[1].Name)
	}
	if services[1].ActiveState != "active" {
		t.Errorf("Service 1 state mismatch: got %s, want active", services[1].ActiveState)
	}
}

func TestServiceOperations_Systemd(t *testing.T) {
mockExec := &MockExecutor{}
profile := &domain.SystemProfile{Distro: "ubuntu", InitSystem: "systemd"}
manager := NewUniversalServiceManager(mockExec, profile)

tests := []struct {
name      string
operation func() error
cmd       string
args      []string
}{
{"Start", func() error { return manager.StartService("nginx") }, "systemctl", []string{"start", "nginx"}},
{"Stop", func() error { return manager.StopService("nginx") }, "systemctl", []string{"stop", "nginx"}},
{"Restart", func() error { return manager.RestartService("nginx") }, "systemctl", []string{"restart", "nginx"}},
{"Enable", func() error { return manager.EnableService("nginx") }, "systemctl", []string{"enable", "nginx"}},
{"Disable", func() error { return manager.DisableService("nginx") }, "systemctl", []string{"disable", "nginx"}},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
mockExec.ExecFunc = func(ctx context.Context, command string, args ...string) (*adapter.CommandResult, error) {
if command != tt.cmd {
return nil, fmt.Errorf("expected command %s, got %s", tt.cmd, command)
}
// Verify args
if len(args) != len(tt.args) {
return nil, fmt.Errorf("expected %d args, got %d", len(tt.args), len(args))
}
for i, arg := range args {
if arg != tt.args[i] {
return nil, fmt.Errorf("arg %d mismatch: got %s, want %s", i, arg, tt.args[i])
}
}
return &adapter.CommandResult{ExitCode: 0}, nil
}

if err := tt.operation(); err != nil {
t.Errorf("%s failed: %v", tt.name, err)
}
})
}
}

func TestGetServiceStatus_Systemd(t *testing.T) {
mockExec := &MockExecutor{}
profile := &domain.SystemProfile{Distro: "ubuntu", InitSystem: "systemd"}
manager := NewUniversalServiceManager(mockExec, profile)

expectedOutput := "Active: active (running)"
mockExec.ExecFunc = func(ctx context.Context, command string, args ...string) (*adapter.CommandResult, error) {
if command != "systemctl" || args[0] != "status" {
return nil, fmt.Errorf("unexpected command")
}
return &adapter.CommandResult{Stdout: expectedOutput, ExitCode: 0}, nil
}

status, err := manager.GetServiceStatus("nginx")
if err != nil {
t.Fatalf("GetServiceStatus failed: %v", err)
}
if status != expectedOutput {
t.Errorf("Expected status %s, got %s", expectedOutput, status)
}
}

func TestGetServiceLogs_Systemd(t *testing.T) {
mockExec := &MockExecutor{}
profile := &domain.SystemProfile{Distro: "ubuntu", InitSystem: "systemd"}
manager := NewUniversalServiceManager(mockExec, profile)

expectedOutput := "Feb 17 10:00:00 server nginx[123]: Started nginx"
mockExec.ExecFunc = func(ctx context.Context, command string, args ...string) (*adapter.CommandResult, error) {
if command != "journalctl" {
return nil, fmt.Errorf("unexpected command")
}
// Verify args: -u nginx -n 50 --no-pager
return &adapter.CommandResult{Stdout: expectedOutput, ExitCode: 0}, nil
}

logs, err := manager.GetServiceLogs("nginx", 50)
if err != nil {
t.Fatalf("GetServiceLogs failed: %v", err)
}
if logs != expectedOutput {
t.Errorf("Expected logs %s, got %s", expectedOutput, logs)
}
}

func TestListServices_Sysvinit(t *testing.T) {
mockExec := &MockExecutor{}
profile := &domain.SystemProfile{Distro: "debian", InitSystem: "sysvinit"}
manager := NewUniversalServiceManager(mockExec, profile)

output := ` [ + ]  nginx
 [ - ]  apache2
 [ ? ]  unknown`

mockExec.ExecFunc = func(ctx context.Context, command string, args ...string) (*adapter.CommandResult, error) {
if command != "service" || args[0] != "--status-all" {
return nil, fmt.Errorf("unexpected command")
}
return &adapter.CommandResult{Stdout: output, ExitCode: 0}, nil
}

services, err := manager.ListServices()
if err != nil {
t.Fatalf("ListServices failed: %v", err)
}

if len(services) != 3 {
t.Errorf("Expected 3 services, got %d", len(services))
}

if services[0].Name != "nginx" || services[0].Status != "active" {
t.Errorf("Service 0 mismatch")
}
if services[1].Name != "apache2" || services[1].Status != "inactive" {
t.Errorf("Service 1 mismatch")
}
if services[2].Name != "unknown" || services[2].Status != "unknown" {
t.Errorf("Service 2 mismatch")
}
}

func TestServiceOperations_Sysvinit(t *testing.T) {
mockExec := &MockExecutor{}
profile := &domain.SystemProfile{Distro: "debian", InitSystem: "sysvinit"}
manager := NewUniversalServiceManager(mockExec, profile)

tests := []struct {
name      string
operation func() error
cmd       string
args      []string
}{
{"Start", func() error { return manager.StartService("nginx") }, "service", []string{"nginx", "start"}},
{"Stop", func() error { return manager.StopService("nginx") }, "service", []string{"nginx", "stop"}},
{"Restart", func() error { return manager.RestartService("nginx") }, "service", []string{"nginx", "restart"}},
{"Enable", func() error { return manager.EnableService("nginx") }, "update-rc.d", []string{"nginx", "enable"}},
{"Disable", func() error { return manager.DisableService("nginx") }, "update-rc.d", []string{"nginx", "disable"}},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
mockExec.ExecFunc = func(ctx context.Context, command string, args ...string) (*adapter.CommandResult, error) {
if command != tt.cmd {
return nil, fmt.Errorf("expected command %s, got %s", tt.cmd, command)
}
return &adapter.CommandResult{ExitCode: 0}, nil
}

if err := tt.operation(); err != nil {
t.Errorf("%s failed: %v", tt.name, err)
}
})
}
}

func TestListServices_OpenRC(t *testing.T) {
mockExec := &MockExecutor{}
profile := &domain.SystemProfile{Distro: "alpine", InitSystem: "openrc"}
manager := NewUniversalServiceManager(mockExec, profile)

mockExec.ExecFunc = func(ctx context.Context, command string, args ...string) (*adapter.CommandResult, error) {
if command != "rc-status" || args[0] != "--all" {
return nil, fmt.Errorf("unexpected command")
}
return &adapter.CommandResult{Stdout: "some output", ExitCode: 0}, nil
}

services, err := manager.ListServices()
if err != nil {
t.Fatalf("ListServices failed: %v", err)
}
if len(services) != 0 {
t.Errorf("Expected 0 services (stub), got %d", len(services))
}
}

func TestServiceOperations_OpenRC(t *testing.T) {
mockExec := &MockExecutor{}
profile := &domain.SystemProfile{Distro: "alpine", InitSystem: "openrc"}
manager := NewUniversalServiceManager(mockExec, profile)

tests := []struct {
name      string
operation func() error
cmd       string
args      []string
}{
{"Start", func() error { return manager.StartService("nginx") }, "rc-service", []string{"nginx", "start"}},
{"Stop", func() error { return manager.StopService("nginx") }, "rc-service", []string{"nginx", "stop"}},
{"Restart", func() error { return manager.RestartService("nginx") }, "rc-service", []string{"nginx", "restart"}},
{"Enable", func() error { return manager.EnableService("nginx") }, "rc-update", []string{"add", "nginx", "default"}},
{"Disable", func() error { return manager.DisableService("nginx") }, "rc-update", []string{"del", "nginx", "default"}},
{"Status", func() error { _, err := manager.GetServiceStatus("nginx"); return err }, "rc-service", []string{"nginx", "status"}},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
mockExec.ExecFunc = func(ctx context.Context, command string, args ...string) (*adapter.CommandResult, error) {
if command != tt.cmd {
return nil, fmt.Errorf("expected command %s, got %s", tt.cmd, command)
}
return &adapter.CommandResult{ExitCode: 0}, nil
}

if err := tt.operation(); err != nil {
t.Errorf("%s failed: %v", tt.name, err)
}
})
}
}

func TestServiceOperations_Runit(t *testing.T) {
mockExec := &MockExecutor{}
profile := &domain.SystemProfile{Distro: "void", InitSystem: "runit"}
manager := NewUniversalServiceManager(mockExec, profile)

tests := []struct {
name      string
operation func() error
cmd       string
args      []string
}{
{"Start", func() error { return manager.StartService("nginx") }, "sv", []string{"start", "nginx"}},
{"Stop", func() error { return manager.StopService("nginx") }, "sv", []string{"stop", "nginx"}},
{"Restart", func() error { return manager.RestartService("nginx") }, "sv", []string{"restart", "nginx"}},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
mockExec.ExecFunc = func(ctx context.Context, command string, args ...string) (*adapter.CommandResult, error) {
if command != tt.cmd {
return nil, fmt.Errorf("expected command %s, got %s", tt.cmd, command)
}
return &adapter.CommandResult{ExitCode: 0}, nil
}

if err := tt.operation(); err != nil {
t.Errorf("%s failed: %v", tt.name, err)
}
})
}
}

func TestServiceOperations_Unsupported(t *testing.T) {
mockExec := &MockExecutor{}
profile := &domain.SystemProfile{Distro: "unknown", InitSystem: "unknown"}
manager := NewUniversalServiceManager(mockExec, profile)

if _, err := manager.ListServices(); err == nil {
t.Error("Expected error for ListServices with unsupported init system")
}
if err := manager.StartService("nginx"); err == nil {
t.Error("Expected error for StartService with unsupported init system")
}
if err := manager.StopService("nginx"); err == nil {
t.Error("Expected error for StopService with unsupported init system")
}
if err := manager.RestartService("nginx"); err == nil {
t.Error("Expected error for RestartService with unsupported init system")
}
if err := manager.EnableService("nginx"); err == nil {
t.Error("Expected error for EnableService with unsupported init system")
}
if err := manager.DisableService("nginx"); err == nil {
t.Error("Expected error for DisableService with unsupported init system")
}
if _, err := manager.GetServiceStatus("nginx"); err == nil {
t.Error("Expected error for GetServiceStatus with unsupported init system")
}
if _, err := manager.GetServiceLogs("nginx", 10); err == nil {
t.Error("Expected error for GetServiceLogs with unsupported init system")
}
}
