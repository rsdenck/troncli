package ports

// SystemMetrics represents a snapshot of system health
type SystemMetrics struct {
	LoadAvg      [3]float64
	CPUUsage     float64 // Percentage
	MemTotal     uint64  // Bytes
	MemUsed      uint64  // Bytes
	SwapTotal    uint64  // Bytes
	SwapUsed     uint64  // Bytes
	DiskIO       DiskIO
	NetworkIO    NetworkIO
	TopProcesses []ProcessInfo
}

type DiskIO struct {
	ReadBytes  uint64
	WriteBytes uint64
	IOPS       uint64
}

type NetworkIO struct {
	RxBytes uint64
	TxBytes uint64
	RxRate  uint64 // Bytes per second
	TxRate  uint64 // Bytes per second
}

type ProcessInfo struct {
	PID    int
	Name   string
	CPU    float64
	Memory uint64
	User   string
}

// SystemMonitor defines operations for gathering system metrics
type SystemMonitor interface {
	GetMetrics() (SystemMetrics, error)
}
