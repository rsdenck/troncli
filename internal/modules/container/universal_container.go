package container

// Package container provides container management capabilities.

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/domain"
	"github.com/mascli/troncli/internal/core/ports"
)

type UniversalContainerManager struct {
	executor adapter.Executor
	profile  *domain.SystemProfile
	runtimes []string // "docker", "podman"
}

func NewUniversalContainerManager(executor adapter.Executor, profile *domain.SystemProfile) *UniversalContainerManager {
	m := &UniversalContainerManager{
		executor: executor,
		profile:  profile,
		runtimes: []string{},
	}
	// Detect runtimes
	if _, err := executor.Exec(context.Background(), "docker", "--version"); err == nil {
		m.runtimes = append(m.runtimes, "docker")
	}
	if _, err := executor.Exec(context.Background(), "podman", "--version"); err == nil {
		m.runtimes = append(m.runtimes, "podman")
	}
	return m
}

func (m *UniversalContainerManager) ListContainers(all bool) ([]ports.Container, error) {
	var result []ports.Container
	ctx := context.Background()

	for _, runtime := range m.runtimes {
		args := []string{"ps", "--format", "{{json .}}"}
		if all {
			args = append(args, "-a")
		}

		res, err := m.executor.Exec(ctx, runtime, args...)
		if err != nil {
			continue // Skip failed runtime
		}

		lines := strings.Split(res.Stdout, "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) == "" {
				continue
			}
			var c map[string]interface{}
			if err := json.Unmarshal([]byte(line), &c); err != nil {
				continue
			}

			id, _ := c["ID"].(string)
			image, _ := c["Image"].(string)
			state, _ := c["State"].(string)
			status, _ := c["Status"].(string)
			namesStr, _ := c["Names"].(string)

			result = append(result, ports.Container{
				ID:      id,
				Names:   strings.Split(namesStr, ","),
				Image:   image,
				State:   state,
				Status:  status,
				Runtime: runtime,
			})
		}
	}
	return result, nil
}

func (m *UniversalContainerManager) StartContainer(id string) error {
	return m.runOnAnyRuntime("start", id)
}

func (m *UniversalContainerManager) StopContainer(id string) error {
	return m.runOnAnyRuntime("stop", id)
}

func (m *UniversalContainerManager) RestartContainer(id string) error {
	return m.runOnAnyRuntime("restart", id)
}

func (m *UniversalContainerManager) RemoveContainer(id string, force bool) error {
	args := []string{"rm"}
	if force {
		args = append(args, "-f")
	}
	args = append(args, id)
	return m.runOnAnyRuntimeArgs(args...)
}

func (m *UniversalContainerManager) GetContainerLogs(id string, tail int) (string, error) {
	ctx := context.Background()
	var lastErr error
	for _, runtime := range m.runtimes {
		res, err := m.executor.Exec(ctx, runtime, "logs", "--tail", fmt.Sprintf("%d", tail), id)
		if err == nil {
			return res.Stdout, nil
		}
		lastErr = err
	}
	return "", fmt.Errorf("failed to get logs from any runtime: %v", lastErr)
}

func (m *UniversalContainerManager) PruneSystem() (string, error) {
	// Prune on all runtimes
	var output strings.Builder
	for _, runtime := range m.runtimes {
		res, err := m.executor.Exec(context.Background(), runtime, "system", "prune", "-f")
		if err != nil {
			output.WriteString(fmt.Sprintf("%s prune failed: %v\n", runtime, err))
		} else {
			output.WriteString(fmt.Sprintf("%s prune: %s\n", runtime, res.Stdout))
		}
	}
	return output.String(), nil
}

func (m *UniversalContainerManager) runOnAnyRuntime(action, id string) error {
	return m.runOnAnyRuntimeArgs(action, id)
}

func (m *UniversalContainerManager) runOnAnyRuntimeArgs(args ...string) error {
	ctx := context.Background()
	var lastErr error
	// This approach is naive (try all). Correct way: check if container exists in runtime first.
	// But for "start", "stop", etc., usually the ID is unique enough or the operation fails fast.
	for _, runtime := range m.runtimes {
		_, err := m.executor.Exec(ctx, runtime, args...)
		if err == nil {
			return nil
		}
		lastErr = err
	}
	return fmt.Errorf("operation failed on all runtimes: %v", lastErr)
}
