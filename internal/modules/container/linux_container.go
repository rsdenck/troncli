package container

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/mascli/troncli/internal/core/ports"
)

// LinuxContainerManager implements ports.ContainerManager
type LinuxContainerManager struct {
	runtimes []string
}

// NewLinuxContainerManager creates a new container manager
func NewLinuxContainerManager() *LinuxContainerManager {
	runtimes := []string{}
	if _, err := exec.LookPath("podman"); err == nil {
		runtimes = append(runtimes, "podman")
	}
	if _, err := exec.LookPath("docker"); err == nil {
		runtimes = append(runtimes, "docker")
	}
	// Incus check could be added here
	return &LinuxContainerManager{runtimes: runtimes}
}

type containerJSON struct {
	ID    string   `json:"ID"`
	Names []string `json:"Names"`
	Image string   `json:"Image"`
	State string   `json:"State"`
}

// ListContainers returns a list of containers from all available runtimes
func (m *LinuxContainerManager) ListContainers(all bool) ([]ports.Container, error) {
	var result []ports.Container

	for _, runtime := range m.runtimes {
		// Run `runtime ps -a`
		args := []string{"ps", "--format", "{{.ID}}|{{.Names}}|{{.Image}}|{{.State}}|{{.Status}}"}
		if all {
			args = append(args, "-a")
		}
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
		return err
	}
	return exec.Command(runtime, "start", id).Run()
}

// StopContainer stops a container
func (m *LinuxContainerManager) StopContainer(id string) error {
	runtime, err := m.findRuntimeForContainer(id)
	if err != nil {
		return err
	}
	return exec.Command(runtime, "stop", id).Run()
}

// RestartContainer restarts a container
func (m *LinuxContainerManager) RestartContainer(id string) error {
	runtime, err := m.findRuntimeForContainer(id)
	if err != nil {
		return err
	}
	return exec.Command(runtime, "restart", id).Run()
}

// RemoveContainer removes a container
func (m *LinuxContainerManager) RemoveContainer(id string, force bool) error {
	runtime, err := m.findRuntimeForContainer(id)
	if err != nil {
		return err
	}
	args := []string{"rm"}
	if force {
		args = append(args, "-f")
	}
	args = append(args, id)
	return exec.Command(runtime, args...).Run()
}

// GetContainerLogs returns the logs of a container
func (m *LinuxContainerManager) GetContainerLogs(id string, tail int) (string, error) {
	runtime, err := m.findRuntimeForContainer(id)
	if err != nil {
		return "", err
	}
	args := []string{"logs", "--tail", fmt.Sprintf("%d", tail), id}
	out, err := exec.Command(runtime, args...).CombinedOutput()
	return string(out), err
}

// PruneSystem removes unused data
func (m *LinuxContainerManager) PruneSystem() (string, error) {
	var output strings.Builder
	for _, runtime := range m.runtimes {
		// system prune -f
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
		// inspect returns 0 if found
		if err := exec.Command(runtime, "inspect", id).Run(); err == nil {
			return runtime, nil
		}
	}
	return "", fmt.Errorf("container %s not found in any runtime", id)
}
