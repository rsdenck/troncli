//go:build linux

package process

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/mascli/troncli/internal/core/ports"
)

// LinuxProcessManager implements ProcessManager for Linux
type LinuxProcessManager struct{}

// NewLinuxProcessManager creates a new LinuxProcessManager
func NewLinuxProcessManager() ports.ProcessManager {
	return &LinuxProcessManager{}
}

// KillProcess sends a signal to a process
func (m *LinuxProcessManager) KillProcess(pid int, signal string) error {
	// Find the process
	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process %d: %w", pid, err)
	}

	var sig syscall.Signal
	switch signal {
	case "SIGTERM":
		sig = syscall.SIGTERM
	case "SIGKILL":
		sig = syscall.SIGKILL
	case "SIGINT":
		sig = syscall.SIGINT
	default:
		return fmt.Errorf("unsupported signal: %s", signal)
	}

	// Send signal to kernel
	// Effect: Process receives signal and terminates or handles it
	// Resource: /proc/{pid}
	if err := process.Signal(sig); err != nil {
		return fmt.Errorf("failed to send signal %s to process %d: %w", signal, pid, err)
	}

	// Verification Protocol
	// Verify if process is still running (for SIGKILL)
	if sig == syscall.SIGKILL {
		// Give the kernel a moment to reap
		time.Sleep(100 * time.Millisecond)

		// Check if /proc/{pid} still exists
		if _, err := os.Stat(fmt.Sprintf("/proc/%d", pid)); !os.IsNotExist(err) {
			// It might be a zombie, check state
			state, _ := getProcessState(pid)
			if state != "Z" { // Z is Zombie
				return fmt.Errorf("process %d still exists after SIGKILL", pid)
			}
		}
	}

	return nil
}

// ReniceProcess changes the priority of a process
func (m *LinuxProcessManager) ReniceProcess(pid int, priority int) error {
	if priority < -20 || priority > 19 {
		return fmt.Errorf("priority must be between -20 and 19")
	}

	// Change priority in kernel
	// Effect: visible in 'top' or 'ps -o ni'
	// Resource: /proc/{pid}/stat (nice value)
	if err := syscall.Setpriority(syscall.PRIO_PROCESS, pid, priority); err != nil {
		return fmt.Errorf("failed to renice process %d to %d: %w", pid, priority, err)
	}

	// Verification Protocol
	// Verify new priority using 'ps'
	// Using exec.Command as requested for system tool verification
	// 'ps -o ni -p PID --no-headers' prints the nice value
	cmd := exec.Command("ps", "-o", "ni", "-p", strconv.Itoa(pid), "--no-headers")
	output, err := cmd.Output()
	if err != nil {
		// If ps fails, maybe process died?
		return fmt.Errorf("failed to verify renice for process %d: %w", pid, err)
	}

	// Output should be the nice value, e.g., "  0" or " -10"
	valStr := strings.TrimSpace(string(output))
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return fmt.Errorf("failed to parse verification output for process %d: %s", pid, valStr)
	}

	if val != priority {
		// Sometimes verification might be tricky if user lacks permissions to read other process details
		// but since we could set priority (which requires root usually for negative), we should be able to read it.
		return fmt.Errorf("verification failed: expected priority %d, got %d", priority, val)
	}

	return nil
}

func getProcessState(pid int) (string, error) {
	// Read /proc/{pid}/stat
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return "", err
	}
	// The state is the 3rd field
	// However, the filename (2nd field) can contain spaces, so we need to be careful.
	// The filename is in parentheses.
	s := string(data)
	closeParen := strings.LastIndex(s, ")")
	if closeParen == -1 || closeParen+2 >= len(s) {
		return "", fmt.Errorf("parse error")
	}
	// state is the character after ") "
func (m *LinuxProcessManager) KillZombies() (int, error) {
	// Strategy:
	// 1. Find all Zombie processes
	// 2. Identify their parents (PPID)
	// 3. Send SIGCHLD to parents to trigger reaping
	// 4. If parent is init (1), it should happen automatically, but we can verify.
	
	procs, err := os.ReadDir("/proc")
	if err != nil {
		return 0, fmt.Errorf("failed to read /proc: %w", err)
	}

	zombieCount := 0
	parentsNotified := make(map[int]bool)

	for _, p := range procs {
		if !p.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(p.Name())
		if err != nil {
			continue
		}

		state, err := getProcessState(pid)
		if err != nil {
			continue
		}

		if state == "Z" {
			// Found a zombie
			ppid, err := getPPID(pid)
			if err != nil {
				continue
			}

			if ppid > 1 {
				if !parentsNotified[ppid] {
					// Send SIGCHLD to parent
					// Effect: Parent receives signal and should call wait()
					parentProc, err := os.FindProcess(ppid)
					if err == nil {
						// Ignore error if parent is gone
						_ = parentProc.Signal(syscall.SIGCHLD)
						parentsNotified[ppid] = true
					}
				}
				zombieCount++
			}
		}
	}

	return zombieCount, nil
}

func getPPID(pid int) (int, error) {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return 0, err
	}
	s := string(data)
	closeParen := strings.LastIndex(s, ")")
	if closeParen == -1 || closeParen+2 >= len(s) {
		return 0, fmt.Errorf("parse error")
	}
	// Fields after ") ": state ppid pgrp ...
	rest := s[closeParen+2:]
	fields := strings.Fields(rest)
	if len(fields) < 2 {
		return 0, fmt.Errorf("parse error")
	}
	// 0 is state, 1 is ppid
	return strconv.Atoi(fields[1])
}

