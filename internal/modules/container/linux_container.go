package container

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/mascli/troncli/internal/core/ports"
)

// LinuxContainerManager implements ports.ContainerManager
// Compliance: Kernel Executor Standards
type LinuxContainerManager struct {
	runtimes []string
}

// NewLinuxContainerManager creates a new container manager
// Effect: Probes system PATH for container runtimes
func NewLinuxContainerManager() *LinuxContainerManager {
	runtimes := []string{}
	// Verify podman availability
	if _, err := exec.LookPath("podman"); err == nil {
		runtimes = append(runtimes, "podman")
	}
	// Verify docker availability
	if _, err := exec.LookPath("docker"); err == nil {
		runtimes = append(runtimes, "docker")
	}
	return &LinuxContainerManager{runtimes: runtimes}
}

// ListContainers returns a list of containers from all available runtimes
func (m *LinuxContainerManager) ListContainers(all bool) ([]ports.Container, error) {
	var result []ports.Container

	for _, runtime := range m.runtimes {
		// Run `runtime ps -a`
		// Effect: Reads container state from disk/memory
		args := []string{"ps", "--format", "{{.ID}}|{{.Names}}|{{.Image}}|{{.State}}|{{.Status}}"}
		if all {
			args = append(args, "-a")
		}
		//nolint:gosec // G204: Runtime determined by LookPath, args are not shell-executed
		cmd := exec.Command(runtime, args...)
		out, err := cmd.Output()
		if err != nil {
			continue // Skip failed runtimes
		}

		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) == "" {
				continue
			}
			parts := strings.Split(line, "|")
			if len(parts) < 5 {
				continue
			}

			result = append(result, ports.Container{
				ID:      parts[0],
				Names:   strings.Split(parts[1], ","),
				Image:   parts[2],
				State:   parts[3],
				Status:  parts[4],
				Runtime: runtime,
			})
		}
	}

	return result, nil
}

// StartContainer starts a container
func (m *LinuxContainerManager) StartContainer(id string) error {
	runtime, err := m.findRuntimeForContainer(id)
	if err != nil {
		return fmt.Errorf("pre-flight check failed: %w", err)
	}

	// Create/Start container process in new namespaces
	// Effect: fork/exec, clone(CLONE_NEWNS|CLONE_NEWPID...), cgroup creation
	// Resource: /proc/{pid}, /sys/fs/cgroup/{runtime}/{id}
	//nolint:gosec // G204: Runtime determined by LookPath, args are not shell-executed
	if err := exec.Command(runtime, "start", id).Run(); err != nil {
		return fmt.Errorf("kernel execution failed (start): %w", err)
	}

	// Verify state
	// Effect: Query runtime state to confirm process running
	//nolint:gosec // G204: Runtime determined by LookPath
	out, err := exec.Command(runtime, "inspect", "--format", "{{.State.Running}}", id).Output()
	if err != nil {
		return fmt.Errorf("verification failed: could not inspect container: %w", err)
	}
	if strings.TrimSpace(string(out)) != "true" {
		return fmt.Errorf("verification failed: container state is not 'Running'")
	}

	return nil
}

// StopContainer stops a container
func (m *LinuxContainerManager) StopContainer(id string) error {
	runtime, err := m.findRuntimeForContainer(id)
	if err != nil {
		return fmt.Errorf("pre-flight check failed: %w", err)
	}

	// Send signal to container process
	// Effect: kill(pid, SIGTERM/SIGKILL)
	// Resource: /proc/{pid}
	//nolint:gosec // G204: Runtime determined by LookPath, args are not shell-executed
	if err := exec.Command(runtime, "stop", id).Run(); err != nil {
		return fmt.Errorf("kernel execution failed (stop): %w", err)
	}

	// Verify state
	// Effect: Query runtime state to confirm process stopped
	//nolint:gosec // G204: Runtime determined by LookPath
	out, err := exec.Command(runtime, "inspect", "--format", "{{.State.Running}}", id).Output()
	if err != nil {
		return fmt.Errorf("verification failed: could not inspect container: %w", err)
	}
	if strings.TrimSpace(string(out)) == "true" {
		return fmt.Errorf("verification failed: container state is still 'Running'")
	}

	return nil
}

// RestartContainer restarts a container
func (m *LinuxContainerManager) RestartContainer(id string) error {
	runtime, err := m.findRuntimeForContainer(id)
	if err != nil {
		return fmt.Errorf("pre-flight check failed: %w", err)
	}

	// Restart container process
	// Effect: stop -> start cycle (kill + fork/exec)
	//nolint:gosec // G204: Runtime determined by LookPath, args are not shell-executed
	if err := exec.Command(runtime, "restart", id).Run(); err != nil {
		return fmt.Errorf("kernel execution failed (restart): %w", err)
	}

	// Verify state
	// Effect: Query runtime state to confirm process running
	//nolint:gosec // G204: Runtime determined by LookPath
	out, err := exec.Command(runtime, "inspect", "--format", "{{.State.Running}}", id).Output()
	if err != nil {
		return fmt.Errorf("verification failed: could not inspect container: %w", err)
	}
	if strings.TrimSpace(string(out)) != "true" {
		return fmt.Errorf("verification failed: container state is not 'Running'")
	}

	return nil
}

// RemoveContainer removes a container
func (m *LinuxContainerManager) RemoveContainer(id string, force bool) error {
	runtime, err := m.findRuntimeForContainer(id)
	if err != nil {
		return fmt.Errorf("pre-flight check failed: %w", err)
	}

	args := []string{"rm"}
	if force {
		args = append(args, "-f")
	}
	args = append(args, id)

	// Remove container resources
	// Effect: Delete cgroups, namespaces, and storage layers
	// Resource: /var/lib/{runtime}/containers/{id}
	//nolint:gosec // G204: Runtime determined by LookPath, args are not shell-executed
	if err := exec.Command(runtime, args...).Run(); err != nil {
		return fmt.Errorf("kernel execution failed (remove): %w", err)
	}

	// Verify removal
	// Effect: Ensure inspect fails or returns empty
	//nolint:gosec // G204: Runtime determined by LookPath
	if err := exec.Command(runtime, "inspect", id).Run(); err == nil {
		return fmt.Errorf("verification failed: container %s still exists", id)
	}

	return nil
}

// GetContainerLogs returns the logs of a container
func (m *LinuxContainerManager) GetContainerLogs(id string, tail int) (string, error) {
	runtime, err := m.findRuntimeForContainer(id)
	if err != nil {
		return "", fmt.Errorf("pre-flight check failed: %w", err)
	}

	// Read container logs
	// Effect: Read from log driver (json-file/journald)
	// Resource: /var/lib/{runtime}/containers/{id}/{id}-json.log
	args := []string{"logs", "--tail", fmt.Sprintf("%d", tail), id}
	//nolint:gosec // G204: Runtime determined by LookPath, args are not shell-executed
	out, err := exec.Command(runtime, args...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("kernel execution failed (logs): %w", err)
	}
	return string(out), nil
}

// PruneSystem removes unused data
func (m *LinuxContainerManager) PruneSystem() (string, error) {
	var output strings.Builder
	for _, runtime := range m.runtimes {
		// Prune unused resources
		// Effect: Delete dangling images, stopped containers, unused networks
		// Resource: /var/lib/{runtime}/*
		//nolint:gosec // G204: Runtime determined by LookPath, args are not shell-executed
		cmd := exec.Command(runtime, "system", "prune", "-f")
		out, err := cmd.CombinedOutput()
		if err != nil {
			output.WriteString(fmt.Sprintf("%s prune failed: %v\n", runtime, err))
		} else {
			output.WriteString(fmt.Sprintf("%s prune:\n%s\n", runtime, string(out)))
		}
	}
	return output.String(), nil
}

// Helper to find which runtime a container belongs to
func (m *LinuxContainerManager) findRuntimeForContainer(id string) (string, error) {
	for _, runtime := range m.runtimes {
		// Inspect container existence
		// Effect: Stat container config in /var/lib/{runtime}
		//nolint:gosec // G204: Runtime determined by LookPath, args are not shell-executed
		if err := exec.Command(runtime, "inspect", id).Run(); err == nil {
			return runtime, nil
		}
	}
	return "", fmt.Errorf("container %s not found in any runtime", id)
}
