package scheduler

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mascli/troncli/internal/core/ports"
)

// ParseCrontabOutput parses the output of `crontab -l`
func ParseCrontabOutput(output string) ([]ports.CronJob, error) {
	var jobs []ports.CronJob
	lines := strings.Split(output, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		
		var schedule, command string
		
		// Handle special strings like @reboot
		if strings.HasPrefix(fields[0], "@") {
			if len(fields) < 2 {
				continue
			}
			schedule = fields[0]
			command = strings.Join(fields[1:], " ")
		} else {
			// Standard 5 fields
			if len(fields) < 6 {
				continue
			}
			schedule = strings.Join(fields[:5], " ")
			command = strings.Join(fields[5:], " ")
		}

		jobs = append(jobs, ports.CronJob{
			ID:       fmt.Sprintf("cron-%d", i),
			Schedule: schedule,
			Command:  command,
			User:     "current", // TODO: Detect user if possible
			File:     "user-crontab",
		})
	}
	return jobs, nil
}

// SystemdTimerEntry for JSON unmarshalling
type SystemdTimerEntry struct {
	Unit      string `json:"unit"`
	Activates string `json:"activates"`
	Next      string `json:"next"`
	Left      string `json:"left"`
	Last      string `json:"last"`
	Passed    string `json:"passed"`
}

// ParseSystemdTimersJSON parses the output of `systemctl list-timers --output=json`
func ParseSystemdTimersJSON(output string) ([]ports.SystemdTimer, error) {
	var entries []SystemdTimerEntry
	if err := json.Unmarshal([]byte(output), &entries); err != nil {
		return nil, err
	}

	var timers []ports.SystemdTimer
	for _, entry := range entries {
		timers = append(timers, ports.SystemdTimer{
			Unit:    entry.Unit,
			Next:    entry.Next,
			Last:    entry.Last,
			Left:    entry.Left,
			Passed:  entry.Passed,
			Service: entry.Activates,
		})
	}
	return timers, nil
}
