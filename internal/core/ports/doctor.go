package ports

// HealthStatus defines the severity of a check result
type HealthStatus string

const (
	StatusOk       HealthStatus = "OK"
	StatusWarning  HealthStatus = "WARNING"
	StatusCritical HealthStatus = "CRITICAL"
)

// HealthCheck represents a single system health check result
type HealthCheck struct {
	Name    string
	Status  HealthStatus
	Message string
	Value   string // e.g., "85%", "1.5"
}

// DoctorManager defines operations for system health checks
type DoctorManager interface {
	RunChecks() ([]HealthCheck, error)
}
