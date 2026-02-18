//go:build !linux

package scheduler

import (
	"fmt"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/ports"
)

// LinuxSchedulerManager stub for non-Linux OS
type LinuxSchedulerManager struct{}

// NewLinuxSchedulerManager creates a new scheduler manager stub
func NewLinuxSchedulerManager(executor adapter.Executor) *LinuxSchedulerManager {
	return &LinuxSchedulerManager{}
}

func (m *LinuxSchedulerManager) ListCronJobs() ([]ports.CronJob, error) {
	return nil, fmt.Errorf("scheduler management not supported on this OS")
}

func (m *LinuxSchedulerManager) AddCronJob(job ports.CronJob) error {
	return fmt.Errorf("scheduler management not supported on this OS")
}

func (m *LinuxSchedulerManager) RemoveCronJob(job ports.CronJob) error {
	return fmt.Errorf("scheduler management not supported on this OS")
}

func (m *LinuxSchedulerManager) ListTimers(all bool) ([]ports.SystemdTimer, error) {
	return nil, fmt.Errorf("scheduler management not supported on this OS")
}
