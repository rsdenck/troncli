package disk

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mascli/troncli/internal/core/ports"
)

// SysReader reads disk information directly from /sys/block
type SysReader struct{}

// NewSysReader creates a new SysReader instance
func NewSysReader() *SysReader {
	return &SysReader{}
}

// ReadBlockDevices reads block devices from /sys/block
func (r *SysReader) ReadBlockDevices() ([]ports.BlockDevice, error) {
	entries, err := os.ReadDir("/sys/block")
	if err != nil {
		return nil, fmt.Errorf("failed to read /sys/block: %w", err)
	}

	var devices []ports.BlockDevice
	for _, entry := range entries {
		dev, err := r.readBlockDevice(entry.Name())
		if err != nil {
			// Log error but continue with other devices
			continue
		}
		devices = append(devices, dev)
	}

	return devices, nil
}

// readBlockDevice reads device details from /sys/block/[name]
func (r *SysReader) readBlockDevice(name string) (ports.BlockDevice, error) {
	dev := ports.BlockDevice{Name: name}
	basePath := fmt.Sprintf("/sys/block/%s", name)

	// Read size (in 512-byte sectors)
	sizeData, err := os.ReadFile(filepath.Join(basePath, "size"))
	if err == nil {
		sectors, err := strconv.ParseInt(strings.TrimSpace(string(sizeData)), 10, 64)
		if err == nil {
			// Convert sectors to human-readable size
			sizeBytes := sectors * 512
			dev.Size = formatSize(sizeBytes)
		}
	}

	// Read removable flag
	removableData, err := os.ReadFile(filepath.Join(basePath, "removable"))
	if err == nil {
		removable := strings.TrimSpace(string(removableData)) == "1"
		if removable {
			dev.Type = "removable"
		} else {
			dev.Type = "disk"
		}
	}

	// Read model (if available) - may not exist for all devices
	modelData, err := os.ReadFile(filepath.Join(basePath, "device/model"))
	if err == nil {
		model := strings.TrimSpace(string(modelData))
		// Store model in a comment or extend the struct if needed
		// For now, we'll use Type field to indicate if it's a special device
		if model != "" {
			_ = model // Model read successfully but not stored in current struct
		}
	}

	// Check if it's a partition or whole disk
	// Partitions have a "partition" file in their sysfs directory
	if _, err := os.Stat(filepath.Join(basePath, "partition")); err == nil {
		dev.Type = "part"
	}

	// Read child partitions if this is a whole disk
	if dev.Type == "disk" || dev.Type == "removable" {
		children, err := r.readPartitions(name)
		if err == nil {
			dev.Children = children
		}
	}

	return dev, nil
}

// readPartitions reads partition information for a block device
func (r *SysReader) readPartitions(deviceName string) ([]ports.BlockDevice, error) {
	basePath := fmt.Sprintf("/sys/block/%s", deviceName)
	entries, err := os.ReadDir(basePath)
	if err != nil {
		return nil, err
	}

	var partitions []ports.BlockDevice
	for _, entry := range entries {
		// Partitions are subdirectories that start with the device name
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), deviceName) {
			continue
		}

		partPath := filepath.Join(basePath, entry.Name())
		
		// Verify it's actually a partition by checking for "partition" file
		if _, err := os.Stat(filepath.Join(partPath, "partition")); err != nil {
			continue
		}

		partition := ports.BlockDevice{
			Name: entry.Name(),
			Type: "part",
		}

		// Read partition size
		sizeData, err := os.ReadFile(filepath.Join(partPath, "size"))
		if err == nil {
			sectors, err := strconv.ParseInt(strings.TrimSpace(string(sizeData)), 10, 64)
			if err == nil {
				sizeBytes := sectors * 512
				partition.Size = formatSize(sizeBytes)
			}
		}

		partitions = append(partitions, partition)
	}

	return partitions, nil
}

// formatSize converts bytes to human-readable format
func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	units := []string{"K", "M", "G", "T", "P", "E"}
	if exp >= len(units) {
		exp = len(units) - 1
	}

	return fmt.Sprintf("%.1f%s", float64(bytes)/float64(div), units[exp])
}

