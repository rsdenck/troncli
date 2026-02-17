package disk

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/domain"
	"github.com/mascli/troncli/internal/core/ports"
)

// UniversalDiskManager implements DiskManager using system tools
type UniversalDiskManager struct {
	executor adapter.Executor
	profile  *domain.SystemProfile
}

// NewUniversalDiskManager creates a new instance
func NewUniversalDiskManager(executor adapter.Executor, profile *domain.SystemProfile) *UniversalDiskManager {
	return &UniversalDiskManager{
		executor: executor,
		profile:  profile,
	}
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

func (m *UniversalDiskManager) ListBlockDevices() ([]ports.BlockDevice, error) {
	ctx := context.Background()
	res, err := m.executor.Exec(ctx, "lsblk", "-J", "-o", "NAME,SIZE,TYPE,MOUNTPOINT")
	if err != nil {
		return nil, fmt.Errorf("lsblk failed: %w", err)
	}

	var result lsblkOutput
	if err := json.Unmarshal([]byte(res.Stdout), &result); err != nil {
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

func (m *UniversalDiskManager) GetFilesystemUsage(path string) (ports.FilesystemUsage, error) {
	ctx := context.Background()
	// Using df -B1 --output=source,size,used,avail,pcent,target
	// But simpler: df -B1 path
	res, err := m.executor.Exec(ctx, "df", "-B1", path)
	if err != nil {
		return ports.FilesystemUsage{}, err
	}

	lines := strings.Split(res.Stdout, "\n")
	if len(lines) < 2 {
		return ports.FilesystemUsage{}, fmt.Errorf("unexpected df output")
	}

	// Filesystem     1B-blocks      Used Available Use% Mounted on
	fields := strings.Fields(lines[1])
	if len(fields) < 4 {
		return ports.FilesystemUsage{}, fmt.Errorf("unexpected df fields")
	}

	total, _ := strconv.ParseUint(fields[1], 10, 64)
	used, _ := strconv.ParseUint(fields[2], 10, 64)
	free, _ := strconv.ParseUint(fields[3], 10, 64)

	// Get Inodes
	resI, err := m.executor.Exec(ctx, "df", "-i", path)
	if err != nil {
		return ports.FilesystemUsage{}, err
	}
	linesI := strings.Split(resI.Stdout, "\n")
	if len(linesI) < 2 {
		return ports.FilesystemUsage{}, fmt.Errorf("unexpected df -i output")
	}
	fieldsI := strings.Fields(linesI[1])
	inodes, _ := strconv.ParseUint(fieldsI[1], 10, 64)
	ifree, _ := strconv.ParseUint(fieldsI[3], 10, 64)

	return ports.FilesystemUsage{
		Path:      path,
		Total:     total,
		Used:      used,
		Free:      free,
		Files:     inodes,
		FilesFree: ifree,
	}, nil
}

func (m *UniversalDiskManager) GetMounts() ([]string, error) {
	ctx := context.Background()
	// cat /proc/mounts
	res, err := m.executor.Exec(ctx, "cat", "/proc/mounts")
	if err != nil {
		return nil, err
	}

	var mounts []string
	lines := strings.Split(res.Stdout, "\n")
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			mounts = append(mounts, parts[1])
		}
	}
	return mounts, nil
}

// GetTopFiles returns largest files in path
func (m *UniversalDiskManager) GetTopFiles(path string, count int) ([]ports.FileNode, error) {
	ctx := context.Background()
	// du -am path | sort -nr | head -n count
	// Using -m for MB to avoid overflow/easy parsing, but KB (-k) is better for precision.
	// du -ak path | sort -nr | head -n count

	cmd := fmt.Sprintf("du -ak %s | sort -nr | head -n %d", path, count)
	res, err := m.executor.Exec(ctx, "sh", "-c", cmd)
	if err != nil {
		return nil, err
	}

	var files []ports.FileNode
	lines := strings.Split(res.Stdout, "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		sizeKB, err := strconv.ParseUint(fields[0], 10, 64)
		if err != nil {
			continue
		}

		filePath := strings.Join(fields[1:], " ")

		files = append(files, ports.FileNode{
			Name:  filePath,
			Path:  filePath,
			Size:  sizeKB * 1024, // Convert to bytes
			IsDir: false,         // Could check if path ends with / or stat it, but du output doesn't specify.
		})
	}
	return files, nil
}

// GetInodeUsage returns inode usage
func (m *UniversalDiskManager) GetInodeUsage(path string) (int, int, error) {
	usage, err := m.GetFilesystemUsage(path)
	if err != nil {
		return 0, 0, err
	}
	return int(usage.Files - usage.FilesFree), int(usage.Files), nil
}

// Cleanup frees disk space by cleaning package caches and logs
func (m *UniversalDiskManager) Cleanup() error {
	ctx := context.Background()

	// Package cleanup based on detected manager
	if m.profile.PackageManager == "apt" {
		m.executor.Exec(ctx, "apt-get", "clean")
		m.executor.Exec(ctx, "apt-get", "autoremove", "-y")
	} else if m.profile.PackageManager == "dnf" || m.profile.PackageManager == "yum" {
		m.executor.Exec(ctx, "dnf", "clean", "all")
	} else if m.profile.PackageManager == "pacman" {
		m.executor.Exec(ctx, "pacman", "-Sc", "--noconfirm")
	}

	// Journal cleanup
	// journalctl --vacuum-time=2d
	// Ignore error if journalctl is missing (e.g. docker)
	m.executor.Exec(ctx, "journalctl", "--vacuum-time=2d")

	return nil
}

// GetDiskHealth returns simple health status
func (m *UniversalDiskManager) GetDiskHealth() (string, error) {
	// Check root usage
	usage, err := m.GetFilesystemUsage("/")
	if err != nil {
		return "Unknown", err
	}

	if usage.Total == 0 {
		return "Unknown", nil
	}

	percent := (float64(usage.Used) / float64(usage.Total)) * 100
	if percent > 90 {
		return "Critical", nil
	}
	if percent > 80 {
		return "Warning", nil
	}
	return "Healthy", nil
}

// Mount mounts a filesystem
func (m *UniversalDiskManager) Mount(source, target, fstype, options string) error {
	args := []string{}
	if fstype != "" {
		args = append(args, "-t", fstype)
	}
	if options != "" {
		args = append(args, "-o", options)
	}
	args = append(args, source, target)

	_, err := m.executor.Exec(context.Background(), "mount", args...)
	return err
}

// Unmount unmounts a filesystem
func (m *UniversalDiskManager) Unmount(target string) error {
	_, err := m.executor.Exec(context.Background(), "umount", target)
	return err
}

// Format formats a device
func (m *UniversalDiskManager) Format(device, fstype string) error {
	if fstype == "" {
		fstype = "ext4"
	}
	// mkfs.fstype
	cmd := fmt.Sprintf("mkfs.%s", fstype)
	_, err := m.executor.Exec(context.Background(), cmd, device)
	return err
}
