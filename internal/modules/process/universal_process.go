package process

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/domain"
	"github.com/mascli/troncli/internal/core/ports"
)

// UniversalProcessManager implements ProcessManager using system tools
type UniversalProcessManager struct {
	executor adapter.Executor
	profile  *domain.SystemProfile
}

// NewUniversalProcessManager creates a new instance
func NewUniversalProcessManager(executor adapter.Executor, profile *domain.SystemProfile) *UniversalProcessManager {
	return &UniversalProcessManager{
		executor: executor,
		profile:  profile,
	}
}

// KillProcess sends a signal to a process
func (m *UniversalProcessManager) KillProcess(pid int, signal string) error {
	ctx := context.Background()
	// Using "kill" command which is universally available
	// Maps standard signals
	sig := "-15" // SIGTERM default
	switch signal {
	case "SIGKILL":
		sig = "-9"
	case "SIGINT":
		sig = "-2"
	case "SIGHUP":
		sig = "-1"
	}

	_, err := m.executor.Exec(ctx, "kill", sig, strconv.Itoa(pid))
	if err != nil {
		return fmt.Errorf("failed to kill process %d: %w", pid, err)
	}
	return nil
}

// ReniceProcess changes process priority
func (m *UniversalProcessManager) ReniceProcess(pid int, priority int) error {
	ctx := context.Background()
	// renice -n priority -p pid
	_, err := m.executor.Exec(ctx, "renice", "-n", strconv.Itoa(priority), "-p", strconv.Itoa(pid))
	if err != nil {
		return fmt.Errorf("failed to renice process %d: %w", pid, err)
	}
	return nil
}

// KillZombies finds and kills zombie processes
func (m *UniversalProcessManager) KillZombies() (int, error) {
	// Strategy: Find Z state processes, get PPID, kill -SIGCHLD PPID
	ctx := context.Background()
	// ps -A -o stat,ppid
	res, err := m.executor.Exec(ctx, "ps", "-A", "-o", "stat,ppid")
	if err != nil {
		return 0, err
	}

	lines := strings.Split(res.Stdout, "\n")
	count := 0
	parents := make(map[string]bool)

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		state := fields[0]
		if strings.HasPrefix(state, "Z") {
			ppid := fields[1]
			if ppid != "1" {
				parents[ppid] = true
				count++
			}
		}
	}

	for ppid := range parents {
		// Send SIGCHLD (17 on Linux x86/ARM usually, but kill -17 might vary?
		// Best to rely on 'kill -s SIGCHLD' if supported, or just let init handle it eventually)
		// Standard kill doesn't always support named signals easily across busybox/gnu.
		// Let's try "kill -s CHLD"
		m.executor.Exec(ctx, "kill", "-s", "CHLD", ppid)
	}

	return count, nil
}

// GetProcessTree returns a tree of processes
func (m *UniversalProcessManager) GetProcessTree() ([]ports.ProcessNode, error) {
	ctx := context.Background()
	// ps -e -o pid,ppid,user,stat,comm
	res, err := m.executor.Exec(ctx, "ps", "-e", "-o", "pid,ppid,user,stat,comm")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(res.Stdout, "\n")
	nodeMap := make(map[int]*ports.ProcessNode)
	var rootNodes []ports.ProcessNode

	// First pass: create nodes
	for i, line := range lines {
		if i == 0 {
			continue
		} // Skip header
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		pid, _ := strconv.Atoi(fields[0])
		ppid, _ := strconv.Atoi(fields[1])
		user := fields[2]
		stat := fields[3]
		comm := strings.Join(fields[4:], " ")

		node := &ports.ProcessNode{
			PID:   pid,
			PPID:  ppid,
			Name:  comm,
			User:  user,
			State: stat,
		}
		nodeMap[pid] = node
	}

	// Second pass: return all nodes
	// Returning flat list is sufficient for TUI to build tree structure using PPID
	var allNodes []ports.ProcessNode
	for _, node := range nodeMap {
		allNodes = append(allNodes, *node)
	}
	return allNodes, nil
}

// GetOpenFiles returns list of open files for a PID
func (m *UniversalProcessManager) GetOpenFiles(pid int) ([]string, error) {
	ctx := context.Background()
	// lsof -p PID -F n
	// Or read /proc/pid/fd

	// /proc is faster and universal on Linux
	fdPath := fmt.Sprintf("/proc/%d/fd", pid)
	// We can use "ls -l /proc/pid/fd" via executor
	res, err := m.executor.Exec(ctx, "ls", "-l", fdPath)
	if err != nil {
		return nil, err
	}

	var files []string
	lines := strings.Split(res.Stdout, "\n")
	for _, line := range lines {
		// lrwx------ 1 root root 64 Feb 17 10:00 0 -> /dev/pts/0
		parts := strings.Split(line, " -> ")
		if len(parts) == 2 {
			files = append(files, parts[1])
		}
	}
	return files, nil
}

// GetProcessPorts returns listening ports for a PID
func (m *UniversalProcessManager) GetProcessPorts(pid int) ([]string, error) {
	ctx := context.Background()
	// ss -lptn | grep pid
	res, err := m.executor.Exec(ctx, "ss", "-lptn")
	if err != nil {
		return nil, err
	}

	var ports []string
	lines := strings.Split(res.Stdout, "\n")
	pidStr := fmt.Sprintf("pid=%d,", pid)

	for _, line := range lines {
		if strings.Contains(line, pidStr) {
			// LISTEN 0 128 0.0.0.0:22 0.0.0.0:* users:(("sshd",pid=123,fd=3))
			fields := strings.Fields(line)
			if len(fields) >= 4 {
				ports = append(ports, fields[3]) // Local Address:Port
			}
		}
	}
	return ports, nil
}
