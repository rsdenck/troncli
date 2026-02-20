//go:build linux

package process

import (
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"
)

func TestLinuxProcessManager_KillProcess(t *testing.T) {
	// Start a sleep process
	cmd := exec.Command("sleep", "100")
	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start sleep process: %v", err)
	}
	pid := cmd.Process.Pid

	// Ensure it's running
	if err := cmd.Process.Signal(syscall.Signal(0)); err != nil {
		t.Fatalf("process not running: %v", err)
	}

	mgr := NewLinuxProcessManager()

	// Kill it
	if err := mgr.KillProcess(pid, "SIGKILL"); err != nil {
		t.Errorf("KillProcess failed: %v", err)
	}

	// Verify it's gone
	// Wait a bit for the system to clean up
	time.Sleep(100 * time.Millisecond)
	
	// We need to reap the process to ensure it doesn't stay as zombie forever during test
	// But KillProcess verification logic handles zombies.
	// Let's verify that KillProcess didn't return error.
	
	// Now let's manually verify
	// If we send signal 0, it might succeed if it's a zombie.
	// If we check state, it should be Z or the process should be gone.
	
	if _, err := os.Stat("/proc/" + string(rune(pid))); os.IsNotExist(err) {
		// Process is gone, good
	} else {
		// Process exists, check if zombie
		state, err := getProcessState(pid)
		if err == nil && state != "Z" {
			t.Errorf("Process %d still exists and is not a zombie (state: %s)", pid, state)
		}
	}
	
	// Cleanup
	cmd.Process.Kill()
	cmd.Wait()
}

func TestLinuxProcessManager_ReniceProcess(t *testing.T) {
	// Start a sleep process
	cmd := exec.Command("sleep", "100")
	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start sleep process: %v", err)
	}
	pid := cmd.Process.Pid
	defer func() {
		cmd.Process.Kill()
		cmd.Wait()
	}()

	mgr := NewLinuxProcessManager()

	// Renice to 10 (Lower priority / Higher niceness)
	// We can only increase niceness as non-root usually.
	if err := mgr.ReniceProcess(pid, 10); err != nil {
		t.Errorf("ReniceProcess failed: %v", err)
	}

	// Verify priority
	prio, err := syscall.Getpriority(syscall.PRIO_PROCESS, pid)
	if err != nil {
		t.Fatalf("failed to get priority: %v", err)
	}
	
	if prio != 10 {
		t.Errorf("Expected priority 10, got %d", prio)
	}
}
