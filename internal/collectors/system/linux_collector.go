//go:build linux

package system

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/mascli/troncli/internal/core/ports"
)

// LinuxSystemMonitor implements SystemMonitor for Linux using /proc
type LinuxSystemMonitor struct {
	lastCPUTotal uint64
	lastCPUIdle  uint64
	lastReadTime time.Time

	// Network delta state
	lastNetReadBytes  uint64
	lastNetWriteBytes uint64
	lastNetTime       time.Time
}

// NewSystemMonitor creates a new LinuxSystemMonitor
func NewSystemMonitor() ports.SystemMonitor {
	return &LinuxSystemMonitor{
		lastNetTime: time.Now(),
	}
}

func (m *LinuxSystemMonitor) GetMetrics() (ports.SystemMetrics, error) {
	metrics := ports.SystemMetrics{}
	var err error
	now := time.Now()

	// 1. Load Average
	metrics.LoadAvg, err = readLoadAvg()
	if err != nil {
		// Log error but continue
	}

	// 2. Memory
	mem, swap, err := readMemInfo()
	if err == nil {
		metrics.MemTotal = mem[0]
		metrics.MemUsed = mem[1]
		metrics.SwapTotal = swap[0]
		metrics.SwapUsed = swap[1]
	}

	// 3. CPU Usage (Delta calculation)
	metrics.CPUUsage = m.calculateCPUUsage()

	// 4. Disk IO
	metrics.DiskIO = readDiskIO()

	// 5. Network IO (Delta calculation)
	metrics.NetworkIO = m.calculateNetworkIO(now)

	// 6. Top Processes
	metrics.TopProcesses = readTopProcesses()

	return metrics, nil
}

func readLoadAvg() ([3]float64, error) {
	data, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return [3]float64{}, err
	}
	parts := strings.Fields(string(data))
	if len(parts) < 3 {
		return [3]float64{}, fmt.Errorf("invalid loadavg format")
	}
	var load [3]float64
	for i := 0; i < 3; i++ {
		load[i], _ = strconv.ParseFloat(parts[i], 64)
	}
	return load, nil
}

func readMemInfo() ([2]uint64, [2]uint64, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return [2]uint64{}, [2]uint64{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var memTotal, memFree, memAvailable, swapTotal, swapFree uint64

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		val, _ := strconv.ParseUint(parts[1], 10, 64)
		// /proc/meminfo is in kB
		val *= 1024

		switch {
		case strings.HasPrefix(parts[0], "MemTotal"):
			memTotal = val
		case strings.HasPrefix(parts[0], "MemFree"):
			memFree = val
		case strings.HasPrefix(parts[0], "MemAvailable"):
			memAvailable = val
		case strings.HasPrefix(parts[0], "SwapTotal"):
			swapTotal = val
		case strings.HasPrefix(parts[0], "SwapFree"):
			swapFree = val
		}
	}

	// MemUsed = Total - Available (if present) or Total - Free
	used := memTotal - memFree
	if memAvailable > 0 {
		used = memTotal - memAvailable
	}

	swapUsed := swapTotal - swapFree

	return [2]uint64{memTotal, used}, [2]uint64{swapTotal, swapUsed}, nil
}

func (m *LinuxSystemMonitor) calculateCPUUsage() float64 {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return 0.0
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 5 || fields[0] != "cpu" {
			return 0.0
		}

		var total uint64
		var idle uint64
		for i, field := range fields[1:] {
			val, _ := strconv.ParseUint(field, 10, 64)
			total += val
			if i == 3 { // idle is the 4th field (index 3 in 0-based slice of values)
				idle = val
			}
		}

		diffTotal := total - m.lastCPUTotal
		diffIdle := idle - m.lastCPUIdle

		m.lastCPUTotal = total
		m.lastCPUIdle = idle

		if diffTotal == 0 {
			return 0.0
		}

		return (float64(diffTotal-diffIdle) / float64(diffTotal)) * 100.0
	}
	return 0.0
}

func readDiskIO() ports.DiskIO {
	file, err := os.Open("/proc/diskstats")
	if err != nil {
		return ports.DiskIO{}
	}
	defer file.Close()

	var readBytes, writeBytes, iops uint64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 14 {
			continue
		}
		// Field 3: reads completed successfully
		reads, _ := strconv.ParseUint(fields[3], 10, 64)
		// Field 5: sectors read
		sectorsRead, _ := strconv.ParseUint(fields[5], 10, 64)
		// Field 7: writes completed
		writes, _ := strconv.ParseUint(fields[7], 10, 64)
		// Field 9: sectors written
		sectorsWritten, _ := strconv.ParseUint(fields[9], 10, 64)

		iops += reads + writes
		readBytes += sectorsRead * 512
		writeBytes += sectorsWritten * 512
	}
	return ports.DiskIO{
		ReadBytes:  readBytes,
		WriteBytes: writeBytes,
		IOPS:       iops,
	}
}

func (m *LinuxSystemMonitor) calculateNetworkIO(now time.Time) ports.NetworkIO {
	file, err := os.Open("/proc/net/dev")
	if err != nil {
		return ports.NetworkIO{}
	}
	defer file.Close()

	var totalRead, totalWrite uint64
	scanner := bufio.NewScanner(file)
	// Skip header lines
	scanner.Scan()
	scanner.Scan()

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) < 10 {
			continue
		}

		// interface name is parts[0], usually ends with ':'
		// RX bytes is parts[1]
		// TX bytes is parts[9]

		rx, _ := strconv.ParseUint(parts[1], 10, 64)
		tx, _ := strconv.ParseUint(parts[9], 10, 64)

		totalRead += rx
		totalWrite += tx
	}

	// Calculate delta
	timeDiff := now.Sub(m.lastNetTime).Seconds()
	if timeDiff <= 0 {
		timeDiff = 1 // Prevent division by zero
	}

	rxRate := float64(totalRead-m.lastNetReadBytes) / timeDiff
	txRate := float64(totalWrite-m.lastNetWriteBytes) / timeDiff

	// Update state
	m.lastNetReadBytes = totalRead
	m.lastNetWriteBytes = totalWrite
	m.lastNetTime = now

	// If this is the first run (deltas are huge or zero), return 0 rate but total bytes
	if rxRate < 0 {
		rxRate = 0
	}
	if txRate < 0 {
		txRate = 0
	}

	// For the very first run, rates might be skewed, but self-corrects on next tick

	return ports.NetworkIO{
		RxBytes: totalRead,
		TxBytes: totalWrite,
		RxRate:  uint64(rxRate),
		TxRate:  uint64(txRate),
	}
}

func readTopProcesses() []ports.ProcessInfo {
	matches, err := filepath.Glob("/proc/[0-9]*")
	if err != nil {
		return nil
	}

	// Read uptime for CPU calculation
	uptimeData, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return nil
	}
	uptimeFields := strings.Fields(string(uptimeData))
	if len(uptimeFields) < 1 {
		return nil
	}
	uptimeSeconds, _ := strconv.ParseFloat(uptimeFields[0], 64)
	clkTck := 100.0 // Default for Linux x86/ARM usually

	var procs []ports.ProcessInfo

	for _, path := range matches {
		pidStr := filepath.Base(path)
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			continue
		}

		statPath := filepath.Join(path, "stat")
		data, err := os.ReadFile(statPath)
		if err != nil {
			continue
		}

		// Parse stat file
		// Format: pid (comm) state ppid ... utime stime ...
		fields := strings.Fields(string(data))
		if len(fields) < 24 {
			continue
		}

		name := strings.Trim(fields[1], "()")
		utime, _ := strconv.ParseFloat(fields[13], 64)
		stime, _ := strconv.ParseFloat(fields[14], 64)
		starttime, _ := strconv.ParseFloat(fields[21], 64)
		rss, _ := strconv.ParseUint(fields[23], 10, 64)

		// Calculate average CPU usage since process start
		// total_time / elapsed_time
		totalTime := utime + stime
		elapsedTicks := (uptimeSeconds * clkTck) - starttime
		if elapsedTicks <= 0 {
			elapsedTicks = 1
		}
		
		// Usage as percentage of one core (can exceed 100%? No, because totalTime is ticks)
		// Wait, if multi-core, totalTime can increase faster than wall clock?
		// No, ticks are ticks.
		// Percentage = (ticks / elapsed_ticks) * 100
		cpuUsage := (totalTime / elapsedTicks) * 100.0

		procs = append(procs, ports.ProcessInfo{
			PID:    pid,
			Name:   name,
			CPU:    cpuUsage,
			Memory: rss * 4096, // Page size usually 4KB
			User:   "",         // Will be filled for top processes
		})
	}

	// Sort by CPU usage (descending)
	sort.Slice(procs, func(i, j int) bool {
		return procs[i].CPU > procs[j].CPU
	})

	// Keep top 10
	if len(procs) > 10 {
		procs = procs[:10]
	}

	// Enrich with User (UID)
	for i := range procs {
		procs[i].User = getProcessOwner(procs[i].PID)
	}

	return procs
}

func getProcessOwner(pid int) string {
	path := filepath.Join("/proc", strconv.Itoa(pid), "status")
	file, err := os.Open(path)
	if err != nil {
		return "?"
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Uid:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				return fields[1] // Real UID
			}
		}
	}
	return "?"
}
