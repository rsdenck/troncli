package container

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/mascli/troncli/internal/core/ports"
)

type DockerManager struct {
	// We could add configuration here if needed
}

func NewDockerManager() ports.ContainerManager {
	return &DockerManager{}
}

func (m *DockerManager) ListContainers(all bool) ([]ports.Container, error) {
	args := []string{"ps", "--format", "{{.ID}}|{{.Names}}|{{.Image}}|{{.State}}|{{.Status}}"}
	if all {
		args = append(args, "-a")
	}

	cmd := exec.Command("docker", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Check if docker is installed/running
		if strings.Contains(string(output), "Is the docker daemon running") {
			return nil, fmt.Errorf("docker daemon is not running")
		}
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	var containers []ports.Container
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) < 5 {
			continue
		}

		containers = append(containers, ports.Container{
			ID:      parts[0],
			Names:   strings.Split(parts[1], ","),
			Image:   parts[2],
			State:   parts[3],
			Status:  parts[4],
			Runtime: "docker",
		})
	}

	return containers, nil
}

func (m *DockerManager) StartContainer(id string) error {
	return exec.Command("docker", "start", id).Run()
}

func (m *DockerManager) StopContainer(id string) error {
	return exec.Command("docker", "stop", id).Run()
}

func (m *DockerManager) RestartContainer(id string) error {
	return exec.Command("docker", "restart", id).Run()
}

func (m *DockerManager) RemoveContainer(id string, force bool) error {
	args := []string{"rm"}
	if force {
		args = append(args, "-f")
	}
	args = append(args, id)
	return exec.Command("docker", args...).Run()
}

func (m *DockerManager) GetContainerLogs(id string, tail int) (string, error) {
	args := []string{"logs"}
	if tail > 0 {
		args = append(args, fmt.Sprintf("--tail=%d", tail))
	}
	args = append(args, id)
	
	output, err := exec.Command("docker", args...).CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func (m *DockerManager) PruneSystem() (string, error) {
	// Force prune without confirmation
	output, err := exec.Command("docker", "system", "prune", "-f").CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
