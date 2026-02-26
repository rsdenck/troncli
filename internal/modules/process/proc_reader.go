//go:build linux

package process

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

// ProcReader reads process information directly from /proc filesystem
type ProcReader struct{}

// NewProcReader creates a new ProcReader instance
func NewProcReader() *ProcReader {
	return &ProcReader{}
}

// Process represents a process with information from /proc
type Process struct {
	PID     int
	PPID    int
	Name    string
	State   string
	Cmdline string
	User    string
	Status  map[string]string
}

// ReadProcessTree reads the process hierarchy from /proc filesystem
// Returns a list of all processes with their parent-child relationships
func (r *ProcReader) ReadProcessTree() ([]Process, error) {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, fmt.Errorf("failed to read /proc: %w", err)
	}

	var processes []Process
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Check if directory name is a PID (numeric)
		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue // Not a process directory
		}

		proc, err := r.readProcess(pid)
		if err != nil {
			// Process may have exited between listing and reading
			continue
		}
		processes = append(processes, proc)
	}

	return processes, nil
}

// readProcess reads process information from /proc/[pid]/
func (r *ProcReader) readProcess(pid int) (Process, error) {
	proc := Process{PID: pid}

	// Read /proc/[pid]/stat for basic info
	if err := r.readStat(pid, &proc); err != nil {
		return proc, fmt.Errorf("failed to read stat: %w", err)
	}

	// Read /proc/[pid]/cmdline for full command
	if err := r.readCmdline(pid, &proc); err != nil {
		// Non-fatal: some processes may not have cmdline
		proc.Cmdline = ""
	}

	// Read /proc/[pid]/status for additional info
	if err := r.readStatus(pid, &proc); err != nil {
		// Non-fatal: continue with partial information
		proc.Status = make(map[string]string)
	}

	return proc, nil
}

// readStat parses /proc/[pid]/stat file
// Format: pid (comm) state ppid pgrp session tty_nr tpgid flags ...
// See proc(5) man page for complete format
func (r *ProcReader) readStat(pid int, proc *Process) error {
	statPath := fmt.Sprintf("/proc/%d/stat", pid)
	data, err := os.ReadFile(statPath)
	if err != nil {
		return err
	}

	// Parse the stat file carefully
	// The command name is in parentheses and may contain spaces
	// Format: pid (comm) state ppid ...
	str := string(data)

	// Find the last ')' to handle command names with parentheses
	lastParen := strings.LastIndex(str, ")")
	if lastParen == -1 {
		return fmt.Errorf("invalid stat format: no closing parenthesis")
	}

	// Extract command name (between first '(' and last ')')
	firstParen := strings.Index(str, "(")
	if firstParen == -1 || firstParen >= lastParen {
		return fmt.Errorf("invalid stat format: malformed parentheses")
	}
	proc.Name = str[firstParen+1 : lastParen]

	// Parse fields after the command name
	// Fields after ')': state ppid pgrp session tty_nr tpgid flags ...
	fieldsStr := strings.TrimSpace(str[lastParen+1:])
	fields := strings.Fields(fieldsStr)

	if len(fields) < 2 {
		return fmt.Errorf("invalid stat format: insufficient fields")
	}

	// Field 0: state (R, S, D, Z, T, etc.)
	proc.State = fields[0]

	// Field 1: ppid
	ppid, err := strconv.Atoi(fields[1])
	if err != nil {
		return fmt.Errorf("invalid ppid: %w", err)
	}
	proc.PPID = ppid

	return nil
}

// readCmdline parses /proc/[pid]/cmdline file
// The cmdline file contains the command line arguments separated by null bytes
func (r *ProcReader) readCmdline(pid int, proc *Process) error {
	cmdlinePath := fmt.Sprintf("/proc/%d/cmdline", pid)
	data, err := os.ReadFile(cmdlinePath)
	if err != nil {
		return err
	}

	// Replace null bytes with spaces
	cmdline := string(bytes.ReplaceAll(data, []byte{0}, []byte(" ")))
	proc.Cmdline = strings.TrimSpace(cmdline)

	// If cmdline is empty, use the name from stat
	if proc.Cmdline == "" {
		proc.Cmdline = proc.Name
	}

	return nil
}

// readStatus parses /proc/[pid]/status file
// This file contains various process information in key-value format
func (r *ProcReader) readStatus(pid int, proc *Process) error {
	statusPath := fmt.Sprintf("/proc/%d/status", pid)
	data, err := os.ReadFile(statusPath)
	if err != nil {
		return err
	}

	proc.Status = make(map[string]string)
	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		// Format: "Key:\tValue"
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		proc.Status[key] = value

		// Extract user information from Uid field
		if key == "Uid" {
			// Format: "Uid: real effective saved filesystem"
			uidFields := strings.Fields(value)
			if len(uidFields) > 0 {
				proc.User = uidFields[0] // Use real UID
			}
		}
	}

	return nil
}

// KillProcess sends a signal to a process using kill(2) syscall
func (r *ProcReader) KillProcess(pid int, signal syscall.Signal) error {
	err := syscall.Kill(pid, signal)
	if err != nil {
		return fmt.Errorf("failed to send signal %v to process %d: %w", signal, pid, err)
	}
	return nil
}

// ReniceProcess changes process priority using setpriority(2) syscall
func (r *ProcReader) ReniceProcess(pid int, priority int) error {
	err := syscall.Setpriority(syscall.PRIO_PROCESS, pid, priority)
	if err != nil {
		return fmt.Errorf("failed to set priority %d for process %d: %w", priority, pid, err)
	}
	return nil
}

// ReadOpenFiles reads open file descriptors from /proc/[pid]/fd
func (r *ProcReader) ReadOpenFiles(pid int) ([]string, error) {
	fdDir := fmt.Sprintf("/proc/%d/fd", pid)
	entries, err := os.ReadDir(fdDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", fdDir, err)
	}

	var files []string
	for _, entry := range entries {
		fdPath := filepath.Join(fdDir, entry.Name())
		link, err := os.Readlink(fdPath)
		if err != nil {
			// Some file descriptors may not be readable
			continue
		}
		files = append(files, link)
	}

	return files, nil
}

// ProcessExists checks if a process exists by checking /proc/[pid]
func (r *ProcReader) ProcessExists(pid int) bool {
	_, err := os.Stat(fmt.Sprintf("/proc/%d", pid))
	return err == nil
}

// GetProcessState returns the state of a process (R, S, D, Z, T, etc.)
func (r *ProcReader) GetProcessState(pid int) (string, error) {
	statPath := fmt.Sprintf("/proc/%d/stat", pid)
	data, err := os.ReadFile(statPath)
	if err != nil {
		return "", fmt.Errorf("failed to read stat: %w", err)
	}

	str := string(data)
	lastParen := strings.LastIndex(str, ")")
	if lastParen == -1 || lastParen+2 >= len(str) {
		return "", fmt.Errorf("invalid stat format")
	}

	// State is the first field after the closing parenthesis
	fieldsStr := strings.TrimSpace(str[lastParen+1:])
	fields := strings.Fields(fieldsStr)
	if len(fields) < 1 {
		return "", fmt.Errorf("invalid stat format: no state field")
	}

	return fields[0], nil
}
