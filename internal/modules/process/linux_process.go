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

func (m *LinuxProcessManager) ReniceProcess(pid int, priority int) error {
	// renice -n priority -p pid
	return exec.Command("renice", "-n", fmt.Sprintf("%d", priority), "-p", fmt.Sprintf("%d", pid)).Run()
}

func (m *LinuxProcessManager) GetProcessTree() ([]ports.ProcessNode, error) {
	// ps -axo pid,ppid,comm,user,stat --sort=ppid
	cmd := exec.Command("ps", "-axo", "pid,ppid,comm,user,stat", "--sort=ppid")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	// Skip header
	if len(lines) > 0 {
		lines = lines[1:]
	}

	var nodes []ports.ProcessNode
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 5 {
			pid, _ := strconv.Atoi(fields[0])
			ppid, _ := strconv.Atoi(fields[1])
			nodes = append(nodes, ports.ProcessNode{
				PID:   pid,
				PPID:  ppid,
				Name:  fields[2],
				User:  fields[3],
				State: fields[4],
			})
		}
	}
	return nodes, nil
}

func (m *LinuxProcessManager) GetOpenFiles(pid int) ([]string, error) {
	// lsof -p pid -F n
	// Or ls -l /proc/pid/fd
	dir := fmt.Sprintf("/proc/%d/fd", pid)
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var openFiles []string
	for _, f := range files {
		target, err := os.Readlink(fmt.Sprintf("%s/%s", dir, f.Name()))
		if err == nil {
			openFiles = append(openFiles, target)
		}
	}
	return openFiles, nil
}

func (m *LinuxProcessManager) GetProcessPorts(pid int) ([]string, error) {
	// ss -nlp | grep pid=pid
	cmd := exec.Command("ss", "-nlp")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var ports []string
	pidStr := fmt.Sprintf("pid=%d", pid)
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, pidStr) {
			ports = append(ports, line)
		}
	}
	return ports, nil
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
	return string(s[closeParen+2]), nil
}

// KillZombies identifies zombie processes and attempts to eliminate them

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

func (m *LinuxProcessManager) GetAllListeningPorts() ([]string, error) {
	// ss -tuln
	cmd := exec.Command("ss", "-tuln")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var ports []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		// Parse logic here (simplified)
		if strings.Contains(line, "LISTEN") {
			ports = append(ports, line)
		}
	}
	return ports, nil
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
