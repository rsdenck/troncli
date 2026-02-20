//go:build !linux

package scheduler

import (
	"testing"

	"github.com/mascli/troncli/internal/core/ports"
)

func TestStubMethods(t *testing.T) {
	// Compliance: Zero Mock Policy
	// We test that the real implementation fails gracefully on non-Linux
	manager := NewLinuxSchedulerManager()

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
