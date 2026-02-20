package doctor

// Package doctor provides system diagnostic capabilities.

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/domain"
	"github.com/mascli/troncli/internal/core/ports"
)

type UniversalDoctorManager struct {
	executor adapter.Executor
	profile  *domain.SystemProfile
}

func NewUniversalDoctorManager(executor adapter.Executor, profile *domain.SystemProfile) *UniversalDoctorManager {
	return &UniversalDoctorManager{
		executor: executor,
		profile:  profile,
	}
}

func (m *UniversalDoctorManager) RunChecks() ([]ports.HealthCheck, error) {
	var checks []ports.HealthCheck

	// 1. Load Average vs Cores
	loadCheck, err := m.checkLoad()
	if err == nil {
		checks = append(checks, loadCheck)
	}

	// 2. Swap Usage
	swapCheck, err := m.checkSwap()
	if err == nil {
		checks = append(checks, swapCheck)
	}

	// 3. Disk Usage (Root)
	diskCheck, err := m.checkDiskRoot()
	if err == nil {
		checks = append(checks, diskCheck)
	}

	// 4. TCP CLOSE_WAIT
	tcpCheck, err := m.checkTCPCloseWait()
	if err == nil {
		checks = append(checks, tcpCheck)
	}

	return checks, nil
}

func (m *UniversalDoctorManager) checkLoad() (ports.HealthCheck, error) {
	ctx := context.Background()
	// Get cores
	resCores, err := m.executor.Exec(ctx, "nproc")
	if err != nil {
		return ports.HealthCheck{}, err
	}
	cores, _ := strconv.Atoi(strings.TrimSpace(resCores.Stdout))
	if cores == 0 { cores = 1 }

	// Get Load
	resLoad, err := m.executor.Exec(ctx, "cat", "/proc/loadavg")
	if err != nil {
		return ports.HealthCheck{}, err
	}
	parts := strings.Fields(resLoad.Stdout)
	if len(parts) < 1 {
		return ports.HealthCheck{}, fmt.Errorf("invalid loadavg")
	}
	load1, _ := strconv.ParseFloat(parts[0], 64)

	status := ports.StatusOk
	msg := fmt.Sprintf("Load %.2f is within capacity (%d cores)", load1, cores)
	if load1 > float64(cores) {
		status = ports.StatusWarning
		msg = fmt.Sprintf("Load %.2f exceeds core count (%d)", load1, cores)
	}
	if load1 > float64(cores)*2 {
		status = ports.StatusCritical
		msg = fmt.Sprintf("Load %.2f is critically high (%d cores)", load1, cores)
	}

	return ports.HealthCheck{
		Name:    "System Load (1m)",
		Status:  status,
		Message: msg,
		Value:   fmt.Sprintf("%.2f", load1),
	}, nil
}

func (m *UniversalDoctorManager) checkSwap() (ports.HealthCheck, error) {
	ctx := context.Background()
	// free -m
	res, err := m.executor.Exec(ctx, "free", "-m")
	if err != nil {
		return ports.HealthCheck{}, err
	}
	lines := strings.Split(res.Stdout, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Swap:") {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				total, _ := strconv.ParseFloat(fields[1], 64)
				used, _ := strconv.ParseFloat(fields[2], 64)
				
				if total == 0 {
					return ports.HealthCheck{
						Name:    "Swap Usage",
						Status:  ports.StatusOk,
						Message: "Swap not configured",
						Value:   "0/0 MB",
					}, nil
				}

				percent := (used / total) * 100
				status := ports.StatusOk
				if percent > 50 { status = ports.StatusWarning }
				if percent > 80 { status = ports.StatusCritical }

				return ports.HealthCheck{
					Name:    "Swap Usage",
					Status:  status,
					Message: fmt.Sprintf("Swap usage is %.1f%%", percent),
					Value:   fmt.Sprintf("%.0f/%.0f MB", used, total),
				}, nil
			}
		}
	}
	return ports.HealthCheck{}, fmt.Errorf("swap info not found")
}

func (m *UniversalDoctorManager) checkDiskRoot() (ports.HealthCheck, error) {
	ctx := context.Background()
	// df -h /
	res, err := m.executor.Exec(ctx, "df", "-h", "/")
	if err != nil {
		return ports.HealthCheck{}, err
	}
	lines := strings.Split(res.Stdout, "\n")
	if len(lines) < 2 { return ports.HealthCheck{}, fmt.Errorf("df output error") }
	
	fields := strings.Fields(lines[1])
	if len(fields) < 5 { return ports.HealthCheck{}, fmt.Errorf("df fields error") }
	
	useStr := strings.TrimSuffix(fields[4], "%")
	use, _ := strconv.Atoi(useStr)
	
	status := ports.StatusOk
	if use > 80 { status = ports.StatusWarning }
	if use > 90 { status = ports.StatusCritical }

	return ports.HealthCheck{
		Name:    "Root Disk Usage",
		Status:  status,
		Message: fmt.Sprintf("Root partition usage is %d%%", use),
		Value:   fields[4],
	}, nil
}

func (m *UniversalDoctorManager) checkTCPCloseWait() (ports.HealthCheck, error) {
	ctx := context.Background()
	// ss -tn state close-wait | wc -l
	// We use shell for pipe
	res, err := m.executor.Exec(ctx, "sh", "-c", "ss -tn state close-wait | wc -l")
	if err != nil {
		// Fallback or ignore
		return ports.HealthCheck{}, err
	}
	
	countStr := strings.TrimSpace(res.Stdout)
	count, _ := strconv.Atoi(countStr)
	// 'ss' prints header line? "Recv-Q Send-Q ..."
	// If header is present, wc -l is count+1.
	// But ss with state filter might output header if no matches?
	// Actually ss usually prints header.
	// Let's assume header is always there, so count-1 is the number of sockets.
	// Or check output content.
	// safer: "ss -tnH state close-wait | wc -l" (-H no header)
	
	resSafe, err := m.executor.Exec(ctx, "sh", "-c", "ss -tnH state close-wait | wc -l")
	if err == nil {
		countStr = strings.TrimSpace(resSafe.Stdout)
		count, _ = strconv.Atoi(countStr)
	}

	status := ports.StatusOk
	if count > 50 { status = ports.StatusWarning }
	if count > 200 { status = ports.StatusCritical }

	return ports.HealthCheck{
		Name:    "TCP CLOSE_WAIT Sockets",
		Status:  status,
		Message: fmt.Sprintf("Found %d sockets in CLOSE_WAIT", count),
		Value:   strconv.Itoa(count),
	}, nil
}
