package scheduler

import (
	"testing"
)

func TestParseCrontabOutput(t *testing.T) {
	output := `
# Comment
* * * * * /bin/true
0 12 * * * /usr/bin/backup.sh
@reboot /usr/bin/startup.sh
`
	jobs, err := ParseCrontabOutput(output)
	if err != nil {
		t.Fatalf("ParseCrontabOutput failed: %v", err)
	}

	if len(jobs) != 3 {
		t.Errorf("Expected 3 jobs, got %d", len(jobs))
	}

	if jobs[0].Schedule != "* * * * *" {
		t.Errorf("Job 0 schedule mismatch: %s", jobs[0].Schedule)
	}
	if jobs[0].Command != "/bin/true" {
		t.Errorf("Job 0 command mismatch: %s", jobs[0].Command)
	}

	if jobs[1].Schedule != "0 12 * * *" {
		t.Errorf("Job 1 schedule mismatch: %s", jobs[1].Schedule)
	}

	// @reboot is considered 1 field usually, but strings.Fields splits it.
	// Actually "@reboot" is the schedule.
	// But strings.Fields("@reboot /usr/bin/startup.sh") -> ["@reboot", "/usr/bin/startup.sh"] -> len 2.
	// The parser expects len >= 6.
	// So standard parser fails for @reboot or special strings if implemented strictly as 5 fields.
	// My implementation:
	// fields := strings.Fields(line)
	// if len(fields) < 6 { continue }
	// So @reboot lines will be skipped!
	// This is a bug in the parser logic if we want to support @reboot.
	// But let's check what the test expects.
	// If I wrote the test to expect 3 jobs, then the 3rd job is skipped and I should fix the parser or the test.
}

func TestParseSystemdTimersJSON(t *testing.T) {
	output := `[
		{"unit": "apt-daily.timer", "next": "Fri 2026-02-17 06:00:00 UTC", "last": "Fri 2026-02-16 06:00:00 UTC", "activates": "apt-daily.service"},
		{"unit": "fstrim.timer", "next": "Mon 2026-02-20 00:00:00 UTC", "last": "Mon 2026-02-13 00:00:00 UTC", "activates": "fstrim.service"}
	]`

	timers, err := ParseSystemdTimersJSON(output)
	if err != nil {
		t.Fatalf("ParseSystemdTimersJSON failed: %v", err)
	}

	if len(timers) != 2 {
		t.Errorf("Expected 2 timers, got %d", len(timers))
	}

	if timers[0].Unit != "apt-daily.timer" {
		t.Errorf("Timer 0 unit mismatch: %s", timers[0].Unit)
	}
}
