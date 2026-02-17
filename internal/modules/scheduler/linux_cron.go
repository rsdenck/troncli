package scheduler

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mascli/troncli/internal/core/ports"
)

// LinuxSchedulerManager implements ports.SchedulerManager
type LinuxSchedulerManager struct{}

// NewLinuxSchedulerManager creates a new scheduler manager
func NewLinuxSchedulerManager() *LinuxSchedulerManager {
	return &LinuxSchedulerManager{}
}

// ListCronJobs returns all cron jobs for the current user
func (m *LinuxSchedulerManager) ListCronJobs() ([]ports.CronJob, error) {
	// crontab -l
	cmd := exec.Command("crontab", "-l")
	out, err := cmd.Output()
	if err != nil {
		// If no crontab for user, it returns error usually "no crontab for user"
		return []ports.CronJob{}, nil
	}

	var jobs []ports.CronJob
	lines := strings.Split(string(out), "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		// Simple parser: first 5 fields are schedule, rest is command
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}

		schedule := strings.Join(fields[:5], " ")
		command := strings.Join(fields[5:], " ")

		jobs = append(jobs, ports.CronJob{
			ID:       fmt.Sprintf("cron-%d", i),
			Schedule: schedule,
			Command:  command,
			User:     "current", // TODO: Detect user
			File:     "user-crontab",
		})
	}
	return jobs, nil
}

// AddCronJob adds a new cron job
func (m *LinuxSchedulerManager) AddCronJob(job ports.CronJob) error {
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
}

// RemoveCronJob removes a cron job by matching command and schedule
func (m *LinuxSchedulerManager) RemoveCronJob(job ports.CronJob) error {
	current, err := m.ListCronJobs()
	if err != nil {
		return err
	}

	var newLines []string
	for _, c := range current {
		if c.Schedule == job.Schedule && c.Command == job.Command {
			continue // Skip matching job
		}
		newLines = append(newLines, fmt.Sprintf("%s %s", c.Schedule, c.Command))
	}

	return m.writeCrontab(newLines)
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
	tmpFile.Close()

	// Install new crontab
	cmd := exec.Command("crontab", tmpFile.Name())
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to install crontab: %v, output: %s", err, string(out))
	}
	return nil
}

// ListTimers returns all systemd timers
func (m *LinuxSchedulerManager) ListTimers(all bool) ([]ports.SystemdTimer, error) {
	args := []string{"list-timers", "--no-pager", "--no-legend"}
	if all {
		args = append(args, "--all")
	}

	cmd := exec.Command("systemctl", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list timers: %w", err)
	}

	var timers []ports.SystemdTimer
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		// NEXT LEFT LAST PASSED UNIT ACTIVATES
		// Fields can contain spaces (dates), so standard Fields() is tricky.
		// However, systemctl output is columnar.
		// A better way is using --output json? No, systemctl list-timers doesn't support json output in older versions.
		// Let's rely on fields being separated by multiple spaces? No.
		// Let's try to parse broadly.
		
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}
		
		// Assuming standard output format:
		// NEXT                         LEFT          LAST                         PASSED       UNIT                         ACTIVATES
		// Mon 2024-05-20 00:00:00 UTC  2h 15min left Sun 2024-05-19 00:00:00 UTC  21h ago      logrotate.timer              logrotate.service
		
		// This is hard to parse with Fields due to spaces in dates.
		// But usually the last two columns are UNIT and ACTIVATES.
		
		unit := fields[len(fields)-2]
		service := fields[len(fields)-1]
		
		// Everything before UNIT is timing info.
		// Let's just store the raw line parts for now or simplify.
		
		timers = append(timers, ports.SystemdTimer{
			Unit:    unit,
			Service: service,
			Next:    strings.Join(fields[0:3], " "), // Rough approx
			Left:    fields[len(fields)-3],          // Rough approx
		})
	}
	
	return timers, nil
}
