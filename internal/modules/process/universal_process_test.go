package process

import (
"context"
"fmt"
"testing"

"github.com/mascli/troncli/internal/core/adapter"
"github.com/mascli/troncli/internal/core/domain"
"github.com/mascli/troncli/internal/core/ports"
)

// MockExecutor for testing
type MockExecutor struct {
ExecFunc func(ctx context.Context, command string, args ...string) (*adapter.CommandResult, error)
}

func (m *MockExecutor) Exec(ctx context.Context, command string, args ...string) (*adapter.CommandResult, error) {
if m.ExecFunc != nil {
return m.ExecFunc(ctx, command, args...)
}
return nil, fmt.Errorf("unexpected command: %s", command)
}

func (m *MockExecutor) ExecWithInput(ctx context.Context, input string, command string, args ...string) (*adapter.CommandResult, error) {
// Not used in these tests
return m.Exec(ctx, command, args...)
}

func TestUniversalProcessManager_KillProcess(t *testing.T) {
mockExec := &MockExecutor{}
profile := &domain.SystemProfile{Distro: "ubuntu"}
manager := NewUniversalProcessManager(mockExec, profile)

pid := 1234
signal := "SIGKILL"

mockExec.ExecFunc = func(ctx context.Context, command string, args ...string) (*adapter.CommandResult, error) {
if command == "kill" {
if len(args) == 2 && args[0] == "-9" && args[1] == "1234" {
return &adapter.CommandResult{ExitCode: 0}, nil
}
return nil, fmt.Errorf("unexpected args: %v", args)
}
return nil, fmt.Errorf("unexpected command: %s", command)
}

if err := manager.KillProcess(pid, signal); err != nil {
t.Errorf("KillProcess failed: %v", err)
}
}

func TestUniversalProcessManager_ReniceProcess(t *testing.T) {
mockExec := &MockExecutor{}
profile := &domain.SystemProfile{Distro: "ubuntu"}
manager := NewUniversalProcessManager(mockExec, profile)

pid := 1234
priority := 10

mockExec.ExecFunc = func(ctx context.Context, command string, args ...string) (*adapter.CommandResult, error) {
if command == "renice" {
// expect: renice -n 10 -p 1234
expectedArgs := []string{"-n", "10", "-p", "1234"}
if len(args) != 4 {
return nil, fmt.Errorf("unexpected args count: %d", len(args))
}
for i, arg := range args {
if arg != expectedArgs[i] {
return nil, fmt.Errorf("arg mismatch at %d: got %s, want %s", i, arg, expectedArgs[i])
}
}
return &adapter.CommandResult{ExitCode: 0}, nil
}
return nil, fmt.Errorf("unexpected command: %s", command)
}

if err := manager.ReniceProcess(pid, priority); err != nil {
t.Errorf("ReniceProcess failed: %v", err)
}
}

func TestUniversalProcessManager_KillZombies(t *testing.T) {
mockExec := &MockExecutor{}
profile := &domain.SystemProfile{Distro: "ubuntu"}
manager := NewUniversalProcessManager(mockExec, profile)

// Mock ps output with zombies
// ps -A -o stat,ppid
psOutput := `STAT PPID
S    1
Z    1234
Z    5678
R    1
`

mockExec.ExecFunc = func(ctx context.Context, command string, args ...string) (*adapter.CommandResult, error) {
if command == "ps" {
return &adapter.CommandResult{Stdout: psOutput, ExitCode: 0}, nil
}
if command == "kill" {
// expect: kill -s CHLD ppid
if len(args) < 3 {
return nil, fmt.Errorf("kill args too short: %v", args)
}
ppid := args[2] // kill -s CHLD <ppid>
if ppid == "1234" || ppid == "5678" {
return &adapter.CommandResult{ExitCode: 0}, nil
}
return nil, fmt.Errorf("unexpected ppid to kill: %s", ppid)
}
return nil, fmt.Errorf("unexpected command: %s", command)
}

count, err := manager.KillZombies()
if err != nil {
t.Errorf("KillZombies failed: %v", err)
}
if count != 2 {
t.Errorf("Expected 2 zombies, got %d", count)
}
}

func TestUniversalProcessManager_GetProcessTree(t *testing.T) {
mockExec := &MockExecutor{}
profile := &domain.SystemProfile{Distro: "ubuntu"}
manager := NewUniversalProcessManager(mockExec, profile)

// Mock ps output
// PID PPID USER STAT COMM
psOutput := `PID PPID USER STAT COMM
1 0 root S init
2 0 root S kthreadd
100 1 root S systemd-journal
200 1 root S sshd
201 200 user S sshd
`

mockExec.ExecFunc = func(ctx context.Context, command string, args ...string) (*adapter.CommandResult, error) {
if command == "ps" {
return &adapter.CommandResult{Stdout: psOutput, ExitCode: 0}, nil
}
return nil, fmt.Errorf("unexpected command: %s", command)
}

tree, err := manager.GetProcessTree()
if err != nil {
t.Fatalf("GetProcessTree failed: %v", err)
}

// Expect 2 roots: 1 and 2 (PPID 0)
// Actually implementation treats PPID 0 as root, but also nodes whose PPID is not in map.
// In psOutput, PID 1 has PPID 0. PID 2 has PPID 0.
// Since 0 is not in map, they are roots.
if len(tree) != 2 {
t.Errorf("Expected 2 roots, got %d", len(tree))
}

// Check children of PID 1
var pid1 *ports.ProcessNode
for i := range tree {
if tree[i].PID == 1 {
pid1 = &tree[i]
break
}
}

if pid1 == nil {
t.Fatal("PID 1 not found in roots")
}

// PID 1 children: 100, 200
if len(pid1.Children) != 2 {
t.Errorf("Expected 2 children for PID 1, got %d", len(pid1.Children))
}
}

func TestUniversalProcessManager_GetOpenFiles(t *testing.T) {
mockExec := &MockExecutor{}
profile := &domain.SystemProfile{Distro: "ubuntu"}
manager := NewUniversalProcessManager(mockExec, profile)

lsOutput := `total 0
lrwx------ 1 root root 64 Feb 17 10:00 0 -> /dev/null
lrwx------ 1 root root 64 Feb 17 10:00 1 -> /var/log/syslog
lrwx------ 1 root root 64 Feb 17 10:00 2 -> socket:[12345]
`

mockExec.ExecFunc = func(ctx context.Context, command string, args ...string) (*adapter.CommandResult, error) {
if command == "ls" {
return &adapter.CommandResult{Stdout: lsOutput, ExitCode: 0}, nil
}
return nil, fmt.Errorf("unexpected command: %s", command)
}

files, err := manager.GetOpenFiles(1234)
if err != nil {
t.Fatalf("GetOpenFiles failed: %v", err)
}

if len(files) != 3 {
t.Errorf("Expected 3 files, got %d", len(files))
}
if files[0] != "/dev/null" {
t.Errorf("Expected /dev/null, got %s", files[0])
}
}

func TestUniversalProcessManager_GetAllListeningPorts(t *testing.T) {
mockExec := &MockExecutor{}
profile := &domain.SystemProfile{Distro: "ubuntu"}
manager := NewUniversalProcessManager(mockExec, profile)

// ss -nltu
// Netid State Recv-Q Send-Q Local_Address:Port Peer_Address:Port
ssOutput := `Netid State Recv-Q Send-Q Local_Address:Port Peer_Address:Port
tcp LISTEN 0 128 0.0.0.0:22 0.0.0.0:*
udp UNCONN 0 0 127.0.0.53%lo:53 0.0.0.0:*
`

mockExec.ExecFunc = func(ctx context.Context, command string, args ...string) (*adapter.CommandResult, error) {
if command == "ss" {
return &adapter.CommandResult{Stdout: ssOutput, ExitCode: 0}, nil
}
return nil, fmt.Errorf("unexpected command: %s", command)
}

portsList, err := manager.GetAllListeningPorts()
if err != nil {
t.Fatalf("GetAllListeningPorts failed: %v", err)
}

if len(portsList) != 2 {
t.Errorf("Expected 2 ports, got %d", len(portsList))
}
// Format: Netid/Local_Address:Port
if portsList[0] != "tcp/0.0.0.0:22" {
t.Errorf("Expected tcp/0.0.0.0:22, got %s", portsList[0])
}
}