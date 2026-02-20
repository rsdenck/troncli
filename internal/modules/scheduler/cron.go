package scheduler

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/mascli/troncli/internal/core/ports"
)

type LinuxSchedulerManager struct{}

func NewLinuxSchedulerManager() ports.SchedulerManager {
	return &LinuxSchedulerManager{}
}

// ListCronJobs returns a list of cron jobs
func (m *LinuxSchedulerManager) ListCronJobs() ([]ports.CronJob, error) {
	// List for current user
	// Effect: Read /var/spool/cron/crontabs/{user}
	// Resource: /var/spool/cron/crontabs/{user}
	//nolint:gosec // G204: Constant args
	cmd := exec.Command("crontab", "-l")
	output, err := cmd.CombinedOutput()
	if err != nil {
		// crontab -l returns error if no crontab for user, which is fine
		if strings.Contains(string(output), "no crontab for") {
			return []ports.CronJob{}, nil
		}
		return nil, fmt.Errorf("kernel execution failed (crontab -l): %w", err)
	}

	return ParseCrontabOutput(string(output))
}

// AddCronJob adds a new cron job
func (m *LinuxSchedulerManager) AddCronJob(job ports.CronJob) error {
	current, _ := m.ListCronJobs() // ignore error, start fresh if empty

	// Reconstruct crontab
	var lines []string
	for _, c := range current {
		lines = append(lines, fmt.Sprintf("%s %s", c.Schedule, c.Command))
	}
	// Add new job
	lines = append(lines, fmt.Sprintf("%s %s", job.Schedule, job.Command))

	input := strings.Join(lines, "\n") + "\n"

	// Write new crontab
	// Effect: Update /var/spool/cron/crontabs/{user}
	// Resource: /var/spool/cron/crontabs/{user}
	//nolint:gosec // G204: Constant args
	cmd := exec.Command("crontab", "-")
	cmd.Stdin = strings.NewReader(input)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("kernel execution failed (crontab -): %w", err)
	}

	// Verify addition
	// Effect: Read back crontab to confirm existence
	verified, err := m.ListCronJobs()
	if err != nil {
		return fmt.Errorf("verification failed: %w", err)
	}
	found := false
	for _, v := range verified {
		if v.Schedule == job.Schedule && v.Command == job.Command {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("verification failed: job not found after addition")
	}

	return nil
}

// RemoveCronJob removes a cron job
func (m *LinuxSchedulerManager) RemoveCronJob(job ports.CronJob) error {
	current, err := m.ListCronJobs()
	if err != nil {
		return err
	}

	var lines []string
	removed := false
	for _, c := range current {
		// Simple match
		if c.Schedule == job.Schedule && c.Command == job.Command {
			removed = true
			continue
		}
		lines = append(lines, fmt.Sprintf("%s %s", c.Schedule, c.Command))
	}

	if !removed {
		return fmt.Errorf("job not found")
	}

	input := strings.Join(lines, "\n") + "\n"

	// Write new crontab
	// Effect: Update /var/spool/cron/crontabs/{user}
	// Resource: /var/spool/cron/crontabs/{user}
	//nolint:gosec // G204: Constant args
	cmd := exec.Command("crontab", "-")
	cmd.Stdin = strings.NewReader(input)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("kernel execution failed (crontab -): %w", err)
	}

	// Verify removal
	// Effect: Read back crontab to confirm absence
	verified, err := m.ListCronJobs()
	if err != nil {
		return fmt.Errorf("verification failed: %w", err)
	}
	for _, v := range verified {
		if v.Schedule == job.Schedule && v.Command == job.Command {
			return fmt.Errorf("verification failed: job still exists after removal")
		}
	}

	return nil
}

func (m *LinuxSchedulerManager) ListTimers(all bool) ([]ports.SystemdTimer, error) {
	args := []string{"list-timers", "--no-pager", "--no-legend"}
	if all {
		args = append(args, "--all")
	}

	// ListTimers returns a list of systemd timers
	// Effect: Reads systemd timer state
	// Resource: /run/systemd/system/*
	//nolint:gosec // G204: Arguments are controlled
	cmd := exec.Command("systemctl", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("kernel execution failed (list-timers): %w", err)
	}

	// Try parsing as JSON first (modern systemd)
	if timers, err := ParseSystemdTimersJSON(string(output)); err == nil {
		return timers, nil
	}

	// Fallback to text parsing if JSON fails (older systemd)
	var timers []ports.SystemdTimer
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}

		// Heuristic: Unit is usually last or second to last
		// We'll just take the one ending in .timer
		var unit string
		for _, f := range fields {
			if strings.HasSuffix(f, ".timer") {
				unit = f
				break
			}
		}

		if unit != "" {
			timers = append(timers, ports.SystemdTimer{
				Unit:    unit,
				Next:    "see systemctl",
				Left:    "see systemctl",
				Last:    "see systemctl",
				Service: "see systemctl",
			})
		}
	}
	return timers, nil
}
