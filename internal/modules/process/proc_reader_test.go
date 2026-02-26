//go:build linux

package process

import (
	"os"
	"syscall"
	"testing"
)

func TestProcReader_ReadProcessTree(t *testing.T) {
	reader := NewProcReader()
	processes, err := reader.ReadProcessTree()
	if err != nil {
		t.Fatalf("ReadProcessTree failed: %v", err)
	}

	if len(processes) == 0 {
		t.Error("Expected at least one process, got none")
	}

	// Verify that we can find the current process
	currentPID := os.Getpid()
	found := false
	for _, proc := range processes {
		if proc.PID == currentPID {
			found = true
			if proc.PPID == 0 {
				t.Error("Current process should have a parent")
			}
			if proc.Name == "" {
				t.Error("Process name should not be empty")
			}
			break
		}
	}

	if !found {
		t.Errorf("Could not find current process (PID %d) in process tree", currentPID)
	}
}

func TestProcReader_ReadProcess(t *testing.T) {
	reader := NewProcReader()
	currentPID := os.Getpid()

	proc, err := reader.readProcess(currentPID)
	if err != nil {
		t.Fatalf("readProcess failed: %v", err)
	}

	if proc.PID != currentPID {
		t.Errorf("Expected PID %d, got %d", currentPID, proc.PID)
	}

	if proc.Name == "" {
		t.Error("Process name should not be empty")
	}

	if proc.State == "" {
		t.Error("Process state should not be empty")
	}

	if proc.PPID == 0 {
		t.Error("Process should have a parent")
	}
}

func TestProcReader_ReadProcess_NonExistent(t *testing.T) {
	reader := NewProcReader()
	// Use a very high PID that likely doesn't exist
	_, err := reader.readProcess(999999)
	if err == nil {
		t.Error("Expected error for non-existent process")
	}
}

func TestProcReader_ProcessExists(t *testing.T) {
	reader := NewProcReader()
	currentPID := os.Getpid()

	if !reader.ProcessExists(currentPID) {
		t.Error("Current process should exist")
	}

	if reader.ProcessExists(999999) {
		t.Error("Non-existent process should not exist")
	}
}

func TestProcReader_GetProcessState(t *testing.T) {
	reader := NewProcReader()
	currentPID := os.Getpid()

	state, err := reader.GetProcessState(currentPID)
	if err != nil {
		t.Fatalf("GetProcessState failed: %v", err)
	}

	// Valid states: R (running), S (sleeping), D (disk sleep), Z (zombie), T (stopped)
	validStates := map[string]bool{
		"R": true, "S": true, "D": true, "Z": true, "T": true,
		"t": true, "W": true, "X": true, "x": true, "K": true,
		"P": true, "I": true,
	}

	if !validStates[state] {
		t.Logf("Warning: Unexpected process state '%s' (may be valid on this system)", state)
	}
}

func TestProcReader_ReadOpenFiles(t *testing.T) {
	reader := NewProcReader()
	currentPID := os.Getpid()

	files, err := reader.ReadOpenFiles(currentPID)
	if err != nil {
		t.Fatalf("ReadOpenFiles failed: %v", err)
	}

	// Current process should have at least stdin, stdout, stderr
	if len(files) < 3 {
		t.Logf("Warning: Expected at least 3 open files, got %d", len(files))
	}
}

func TestProcReader_ReadOpenFiles_NonExistent(t *testing.T) {
	reader := NewProcReader()
	_, err := reader.ReadOpenFiles(999999)
	if err == nil {
		t.Error("Expected error for non-existent process")
	}
}

func TestProcReader_KillProcess(t *testing.T) {
	// We can't actually kill processes in a test, but we can test the error handling
	reader := NewProcReader()

	// Try to send signal 0 (null signal) to current process
	// This checks if the process exists without actually sending a signal
	currentPID := os.Getpid()
	err := reader.KillProcess(currentPID, syscall.Signal(0))
	if err != nil {
		t.Errorf("KillProcess with signal 0 should succeed: %v", err)
	}

	// Try to kill a non-existent process
	err = reader.KillProcess(999999, syscall.SIGTERM)
	if err == nil {
		t.Error("Expected error when killing non-existent process")
	}
}

func TestProcReader_ReniceProcess(t *testing.T) {
	// We can't actually change priority without root, but we can test the error handling
	reader := NewProcReader()
	currentPID := os.Getpid()

	// Try to set priority to 0 (normal priority)
	// This may fail if we don't have permissions
	err := reader.ReniceProcess(currentPID, 0)
	if err != nil {
		// This is expected if we don't have permissions
		t.Logf("ReniceProcess failed (expected without root): %v", err)
	}

	// Try to renice a non-existent process
	err = reader.ReniceProcess(999999, 0)
	if err == nil {
		t.Error("Expected error when renicing non-existent process")
	}
}

func TestProcReader_ReadStat(t *testing.T) {
	reader := NewProcReader()
	currentPID := os.Getpid()

	var proc Process
	err := reader.readStat(currentPID, &proc)
	if err != nil {
		t.Fatalf("readStat failed: %v", err)
	}

	if proc.Name == "" {
		t.Error("Process name should not be empty")
	}

	if proc.State == "" {
		t.Error("Process state should not be empty")
	}

	if proc.PPID == 0 {
		t.Error("Process should have a parent")
	}
}

func TestProcReader_ReadCmdline(t *testing.T) {
	reader := NewProcReader()
	currentPID := os.Getpid()

	var proc Process
	proc.Name = "test" // Set a default name
	err := reader.readCmdline(currentPID, &proc)
	if err != nil {
		t.Fatalf("readCmdline failed: %v", err)
	}

	// Cmdline should not be empty for the current process
	if proc.Cmdline == "" {
		t.Error("Process cmdline should not be empty")
	}
}

func TestProcReader_ReadStatus(t *testing.T) {
	reader := NewProcReader()
	currentPID := os.Getpid()

	var proc Process
	err := reader.readStatus(currentPID, &proc)
	if err != nil {
		t.Fatalf("readStatus failed: %v", err)
	}

	if len(proc.Status) == 0 {
		t.Error("Status map should not be empty")
	}

	// Check for some common status fields
	if _, ok := proc.Status["Name"]; !ok {
		t.Error("Status should contain 'Name' field")
	}

	if _, ok := proc.Status["State"]; !ok {
		t.Error("Status should contain 'State' field")
	}
}
