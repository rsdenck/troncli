package ports

// CronJob represents a cron job entry
type CronJob struct {
	ID       string // hash or index
	Schedule string
	Command  string
	User     string // "root" or specific user
	File     string // source file (e.g., /var/spool/cron/root)
}

// SystemdTimer represents a systemd timer
type SystemdTimer struct {
	Unit    string
	Next    string
	Left    string
	Last    string
	Passed  string
	Service string
}

// SchedulerManager defines the interface for scheduling operations
type SchedulerManager interface {
	// ListCronJobs returns all cron jobs for the current user (and root if privileged)
	ListCronJobs() ([]CronJob, error)

	// AddCronJob adds a new cron job
	AddCronJob(job CronJob) error

	// RemoveCronJob removes a cron job by ID (or matching content)
	RemoveCronJob(job CronJob) error

	// ListTimers returns all systemd timers
	ListTimers(all bool) ([]SystemdTimer, error)
}
