//go:build linux

package scheduler

// Package scheduler provides task scheduling capabilities.

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/ports"
)

// LinuxSchedulerManager implements ports.SchedulerManager
type LinuxSchedulerManager struct {
	executor adapter.Executor
}

// NewLinuxSchedulerManager creates a new scheduler manager
func NewLinuxSchedulerManager(executor adapter.Executor) *LinuxSchedulerManager {
	return &LinuxSchedulerManager{
		executor: executor,
	}
}

// withLock executes a function with a file lock to prevent race conditions
func (m *LinuxSchedulerManager) withLock(action func() error) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home dir: %w", err)
	}
	lockDir := filepath.Join(home, ".troncli")
	if err := os.MkdirAll(lockDir, 0700); err != nil {
		return fmt.Errorf("failed to create lock dir: %w", err)
	}
	lockFile := filepath.Join(lockDir, "cron.lock")

	file, err := os.OpenFile(lockFile, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return fmt.Errorf("failed to open lock file: %w", err)
	}
	defer file.Close()

	// Exclusive lock, blocking
	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX); err != nil {
		return fmt.Errorf("failed to acquire lock: %w", err)
	}
	defer syscall.Flock(int(file.Fd()), syscall.LOCK_UN)

	return action()
}

// ListCronJobs returns all cron jobs for the current user
func (m *LinuxSchedulerManager) ListCronJobs() ([]ports.CronJob, error) {
	// crontab -l is safe to read usually, but for consistency we could lock.
	// However, `crontab -l` doesn't support locking mechanism from outside easily unless we manage the file directly.
	// Since we rely on the `crontab` command, we accept a small read race, but write must be locked.
	
	cmd := exec.Command("crontab", "-l")
	out, err := cmd.Output()
	if err != nil {
		// If no crontab for user, it returns error usually "no crontab for user"
		return []ports.CronJob{}, nil
	}

	return ParseCrontabOutput(string(out))
}

// AddCronJob adds a new cron job
func (m *LinuxSchedulerManager) AddCronJob(job ports.CronJob) error {
	return m.withLock(func() error {
		// 1. Get current crontab
		current, _ := m.ListCronJobs() // ignore error, might be empty

		// 2. Append new job
		newLine := fmt.Sprintf("%s %s", job.Schedule, job.Command)
		
		var lines []string
		for _, c := range current {
			lines = append(lines, fmt.Sprintf("%s %s", c.Schedule, c.Command))
		}
		lines = append(lines, newLine)

		// 3. Write back
		return m.writeCrontab(lines)
	})
}

// RemoveCronJob removes a cron job by matching command and schedule
func (m *LinuxSchedulerManager) RemoveCronJob(job ports.CronJob) error {
	return m.withLock(func() error {
		current, err := m.ListCronJobs()
		if err != nil {
			return err
		}

		var newLines []string
		found := false
		for _, c := range current {
			if c.Schedule == job.Schedule && c.Command == job.Command {
				found = true
				continue // Skip matching job
			}
			newLines = append(newLines, fmt.Sprintf("%s %s", c.Schedule, c.Command))
		}
		
		if !found {
			return fmt.Errorf("job not found")
		}

		return m.writeCrontab(newLines)
	})
}

func (m *LinuxSchedulerManager) writeCrontab(lines []string) error {
	content := strings.Join(lines, "\n") + "\n"
	
	// Write to temp file
	tmpFile, err := os.CreateTemp("", "cron-")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		return err
	}
	if err := tmpFile.Close(); err != nil {
		return err
	}

	// Load into crontab
	cmd := exec.Command("crontab", tmpFile.Name())
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to install crontab: %s, output: %s", err, string(out))
	}
	return nil
}

// ListTimers returns all systemd timers
func (m *LinuxSchedulerManager) ListTimers(all bool) ([]ports.SystemdTimer, error) {
	// Try JSON format first (systemd v233+)
	args := []string{"list-timers", "--output=json", "--no-pager"}
	if all {
		args = append(args, "--all")
	}

	cmd := exec.Command("systemctl", args...)
	out, err := cmd.Output()
	
	// If JSON supported and successful
	if err == nil && len(out) > 0 && json.Valid(out) {
		return ParseSystemdTimersJSON(string(out))
	}

	// Fallback to text parsing
	// Reset args for text output
	args = []string{"list-timers", "--no-pager", "--no-legend"}
	if all {
		args = append(args, "--all")
	}

	cmd = exec.Command("systemctl", args...)
	out, err = cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list timers: %w", err)
	}

	var timers []ports.SystemdTimer
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		
		// Text format:
		// NEXT                         LEFT          LAST                         PASSED       UNIT                         ACTIVATES
		// Mon 2024-05-20 00:00:00 UTC  2h 15min left Sun 2024-05-19 00:00:00 UTC  21h ago      logrotate.timer              logrotate.service
		
		// Simplified fallback: just get unit and service from end
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}
		
		unit := fields[len(fields)-2]
		service := fields[len(fields)-1]
		
		timers = append(timers, ports.SystemdTimer{
			Unit:    unit,
			Service: service,
			Next:    "unknown (upgrade systemd for details)", 
			Left:    "unknown",
		})
	}
	
	return timers, nil
}
