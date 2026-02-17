//go:build linux

package disk

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/mascli/troncli/internal/core/ports"
)

type LinuxDiskManager struct{}

func NewLinuxDiskManager() ports.DiskManager {
	return &LinuxDiskManager{}
}

// lsblkOutput represents the JSON structure from lsblk -J
type lsblkOutput struct {
	BlockDevices []blockDeviceJSON `json:"blockdevices"`
}

type blockDeviceJSON struct {
	Name       string            `json:"name"`
	Size       string            `json:"size"`
	Type       string            `json:"type"`
	MountPoint string            `json:"mountpoint"`
	Children   []blockDeviceJSON `json:"children,omitempty"`
}

func (m *LinuxDiskManager) ListBlockDevices() ([]ports.BlockDevice, error) {
	cmd := exec.Command("lsblk", "-J", "-o", "NAME,SIZE,TYPE,MOUNTPOINT")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("lsblk failed: %w", err)
	}

	var result lsblkOutput
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse lsblk json: %w", err)
	}

	return convertDevices(result.BlockDevices), nil
}

func convertDevices(devices []blockDeviceJSON) []ports.BlockDevice {
	var result []ports.BlockDevice
	for _, dev := range devices {
		bd := ports.BlockDevice{
			Name:       dev.Name,
			Size:       dev.Size,
			Type:       dev.Type,
			MountPoint: dev.MountPoint,
		}
		if len(dev.Children) > 0 {
			bd.Children = convertDevices(dev.Children)
		}
		result = append(result, bd)
	}
	return result
}

func (m *LinuxDiskManager) GetUsage(path string) (*ports.DiskUsage, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return nil, err
	}

	// Bsize is uint32 on some archs (e.g. arm), uint64 on others. Cast to uint64 to be safe.
	bsize := uint64(stat.Bsize)
	total := stat.Blocks * bsize
	free := stat.Bfree * bsize
	used := total - free

	usagePercent := 0.0
	if total > 0 {
		usagePercent = float64(used) / float64(total) * 100
	}

	return &ports.DiskUsage{
		Path:         path,
		TotalBytes:   total,
		UsedBytes:    used,
		FreeBytes:    free,
		UsagePercent: usagePercent,
	}, nil
}

func (m *LinuxDiskManager) Cleanup() error {
	// No specific cleanup needed for LinuxDiskManager yet
	return nil
}

func (m *LinuxDiskManager) GetMounts() ([]string, error) {
	data, err := os.ReadFile("/proc/mounts")
	if err != nil {
		return nil, err
	}

	var mounts []string
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			mounts = append(mounts, parts[1])
		}
	}
	return mounts, nil
}
