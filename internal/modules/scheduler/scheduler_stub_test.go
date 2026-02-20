//go:build !linux

package scheduler

import (
"context"
"testing"

"github.com/mascli/troncli/internal/core/adapter"
"github.com/mascli/troncli/internal/core/ports"
)

// MockExecutor implements Executor for testing
type MockExecutor struct {
ExecFunc func(ctx context.Context, command string, args ...string) (*adapter.CommandResult, error)
}

func (m *MockExecutor) Exec(ctx context.Context, command string, args ...string) (*adapter.CommandResult, error) {
if m.ExecFunc != nil {
return m.ExecFunc(ctx, command, args...)
}
return &adapter.CommandResult{ExitCode: 0}, nil
}

func (m *MockExecutor) ExecWithInput(ctx context.Context, input string, command string, args ...string) (*adapter.CommandResult, error) {
return m.Exec(ctx, command, args...)
}

func TestStubMethods(t *testing.T) {
mockExec := &MockExecutor{}
manager := NewLinuxSchedulerManager(mockExec)

if _, err := manager.ListCronJobs(); err == nil {
t.Error("ListCronJobs should return error on non-Linux")
}

if err := manager.AddCronJob(ports.CronJob{}); err == nil {
t.Error("AddCronJob should return error on non-Linux")
}

if err := manager.RemoveCronJob(ports.CronJob{}); err == nil {
t.Error("RemoveCronJob should return error on non-Linux")
}

if _, err := manager.ListTimers(true); err == nil {
t.Error("ListTimers should return error on non-Linux")
}
}