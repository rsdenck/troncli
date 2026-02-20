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

func (m *LinuxSchedulerManager) ListCronJobs() ([]ports.CronJob, error) {
	// List for current user
	cmd := exec.Command("crontab", "-l")
	output, err := cmd.CombinedOutput()
	if err != nil {
		// crontab -l returns error if no crontab for user, which is fine
		if strings.Contains(string(output), "no crontab for") {
			return []ports.CronJob{}, nil
		}
		return nil, err
	}

	var jobs []ports.CronJob
	lines := strings.Split(string(output), "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Basic parsing, assumed standard cron format
		// * * * * * command
		// We could split by space to validate but for now we trust the line

		// To distinguish command from schedule is hard without complex regex
		// We'll store the whole line as Command for now if parsing is complex,
		// but let's try a simple heuristic: first 5 fields are schedule
		fields := strings.Fields(line)
		if len(fields) >= 6 {
			schedule := strings.Join(fields[:5], " ")
			command := strings.Join(fields[5:], " ")
			jobs = append(jobs, ports.CronJob{
				ID:       fmt.Sprintf("user-%d", i),
				Schedule: schedule,
				Command:  command,
				User:     "current",
				File:     "crontab",
			})
		} else {
			// Special @reboot etc
			if strings.HasPrefix(line, "@") && len(fields) >= 2 {
				jobs = append(jobs, ports.CronJob{
					ID:       fmt.Sprintf("user-%d", i),
					Schedule: fields[0],
					Command:  strings.Join(fields[1:], " "),
					User:     "current",
					File:     "crontab",
				})
			}
		}
	}
	return jobs, nil
}

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

	cmd := exec.Command("crontab", "-")
	cmd.Stdin = strings.NewReader(input)
	return cmd.Run()
}

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

	cmd := exec.Command("crontab", "-")
	cmd.Stdin = strings.NewReader(input)
	return cmd.Run()
}

func (m *LinuxSchedulerManager) ListTimers(all bool) ([]ports.SystemdTimer, error) {
	args := []string{"list-timers", "--no-pager", "--no-legend"}
	if all {
		args = append(args, "--all")
	}

	cmd := exec.Command("systemctl", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

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
		// NEXT LEFT LAST PASSED UNIT ACTIVATES
		// We'll try to grab the unit name which ends in .timer

		var unit string
		// var next, left, last string

		// Heuristic parsing: usually unit is 2nd to last or last
		for i, f := range fields {
			if strings.HasSuffix(f, ".timer") {
				unit = f
				// Try to map relative fields if possible
				if i >= 1 {
					// This is very rough approximation
				}
				break
			}
		}

		// If we can't parse easily, we just fill Unit.
		// A better way would be using `systemctl list-timers -o json` if available (newer systemd)
		// but strict text parsing is hard without knowing column widths.
		// Let's assume standard output format and just fill Unit for now,
		// and maybe raw line parts for others if needed.

		if unit != "" {
			timers = append(timers, ports.SystemdTimer{
				Unit: unit,
				Next: "see systemctl", // difficult to parse reliably without json
				Left: "see systemctl",
				Last: "see systemctl",
			})
		}
	}
	return timers, nil
}
