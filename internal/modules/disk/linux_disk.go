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

func (m *LinuxDiskManager) GetFilesystemUsage(path string) (ports.FilesystemUsage, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return ports.FilesystemUsage{}, err
	}

	// Blocks * BlockSize = Total Bytes
	total := uint64(stat.Blocks) * uint64(stat.Bsize)
	free := uint64(stat.Bfree) * uint64(stat.Bsize)
	available := uint64(stat.Bavail) * uint64(stat.Bsize)
	used := total - free

	return ports.FilesystemUsage{
		Path:      path,
		Total:     total,
		Used:      used,
		Free:      available, // Available to unprivileged users
		Files:     uint64(stat.Files),
		FilesFree: uint64(stat.Ffree),
	}, nil
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
