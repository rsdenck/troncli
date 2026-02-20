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

// lsblkOutputLinux represents the JSON structure from lsblk -J
type lsblkOutputLinux struct {
	BlockDevices []blockDeviceJSONLinux `json:"blockdevices"`
}

type blockDeviceJSONLinux struct {
	Name       string                 `json:"name"`
	Size       string                 `json:"size"`
	Type       string                 `json:"type"`
	MountPoint string                 `json:"mountpoint"`
	Children   []blockDeviceJSONLinux `json:"children,omitempty"`
}

func (m *LinuxDiskManager) ListBlockDevices() ([]ports.BlockDevice, error) {
	cmd := exec.Command("lsblk", "-J", "-o", "NAME,SIZE,TYPE,MOUNTPOINT")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("lsblk failed: %w", err)
	}

	var result lsblkOutputLinux
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse lsblk json: %w", err)
	}

	return convertDevicesLinux(result.BlockDevices), nil
}

func convertDevicesLinux(devices []blockDeviceJSONLinux) []ports.BlockDevice {
	var result []ports.BlockDevice
	for _, d := range devices {
		dev := ports.BlockDevice{
			Name:       d.Name,
			Size:       d.Size,
			Type:       d.Type,
			MountPoint: d.MountPoint,
		}
		if len(d.Children) > 0 {
			dev.Children = convertDevicesLinux(d.Children)
		}
		result = append(result, dev)
	}
	return result
}

func (m *LinuxDiskManager) GetFilesystemUsage(path string) (ports.FilesystemUsage, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return ports.FilesystemUsage{}, err
	}

	// Bsize is uint32 on some archs (e.g. arm), uint64 on others. Cast to uint64 to be safe.
	bsize := uint64(stat.Bsize)
	total := stat.Blocks * bsize
	free := stat.Bfree * bsize
	used := total - free

	// Get Inodes
	files := stat.Files
	filesFree := stat.Ffree

	return ports.FilesystemUsage{
		Path:      path,
		Total:     total,
		Used:      used,
		Free:      free,
		Files:     files,
		FilesFree: filesFree,
	}, nil
}

func (m *LinuxDiskManager) Format(device, fstype string) error {
	// mkfs.fstype device
	cmd := exec.Command(fmt.Sprintf("mkfs.%s", fstype), device)
	return cmd.Run()
}

func (m *LinuxDiskManager) GetInodeUsage(path string) (int, int, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return 0, 0, err
	}
	return int(stat.Files - stat.Ffree), int(stat.Files), nil
}

func (m *LinuxDiskManager) GetTopFiles(path string, count int) ([]ports.FileNode, error) {
	// find path -type f -exec du -h {} + | sort -rh | head -n count
	// Simplified: just return empty for now or use 'du'
	return nil, nil
}

func (m *LinuxDiskManager) Mount(source, target, fstype, options string) error {
	args := []string{"mount"}
	if fstype != "" {
		args = append(args, "-t", fstype)
	}
	if options != "" {
		args = append(args, "-o", options)
	}
	args = append(args, source, target)
	return exec.Command(args[0], args[1:]...).Run()
}

func (m *LinuxDiskManager) Unmount(target string) error {
	return exec.Command("umount", target).Run()
}

func (m *LinuxDiskManager) GetDiskHealth() (string, error) {
	// smartctl -H /dev/sda (requires root and knowing the device)
	// For now return "unknown" or check logs
	return "unknown", nil
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
