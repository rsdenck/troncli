//go:build linux

package disk

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/mascli/troncli/internal/core/ports"
)

// ReadMountPoints reads mount points from /proc/mounts
func (r *SysReader) ReadMountPoints() ([]ports.MountPoint, error) {
	data, err := os.ReadFile("/proc/mounts")
	if err != nil {
		return nil, fmt.Errorf("failed to read /proc/mounts: %w", err)
	}

	var mounts []ports.MountPoint
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		mount := ports.MountPoint{
			Device:     fields[0],
			MountPoint: fields[1],
			FSType:     fields[2],
			Options:    fields[3],
		}
		mounts = append(mounts, mount)
	}

	return mounts, nil
}

// GetFilesystemUsage uses syscall.Statfs to get filesystem usage statistics
func (r *SysReader) GetFilesystemUsage(path string) (ports.FilesystemUsage, error) {
	var stat syscall.Statfs_t
	err := syscall.Statfs(path, &stat)
	if err != nil {
		return ports.FilesystemUsage{}, fmt.Errorf("failed to get filesystem stats for %s: %w", path, err)
	}

	// Calculate usage statistics
	// stat.Blocks = total data blocks in filesystem
	// stat.Bfree = free blocks in filesystem
	// stat.Bavail = free blocks available to unprivileged user
	// stat.Bsize = optimal transfer block size

	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bfree * uint64(stat.Bsize)
	available := stat.Bavail * uint64(stat.Bsize)
	used := total - free

	usage := ports.FilesystemUsage{
		Path:      path,
		Total:     total,
		Free:      available, // Use Bavail (available to non-root) for Free
		Used:      used,
		Files:     stat.Files,
		FilesFree: stat.Ffree,
	}

	return usage, nil
}
