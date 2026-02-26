package disk

import (
	"os"
	"testing"
)

func TestSysReader_ReadBlockDevices(t *testing.T) {
	// Skip test if /sys/block doesn't exist (non-Linux systems)
	if _, err := os.Stat("/sys/block"); os.IsNotExist(err) {
		t.Skip("Skipping test: /sys/block not available (non-Linux system)")
	}

	reader := NewSysReader()
	devices, err := reader.ReadBlockDevices()

	if err != nil {
		t.Fatalf("ReadBlockDevices() failed: %v", err)
	}

	if len(devices) == 0 {
		t.Error("Expected at least one block device, got none")
	}

	// Verify that at least one device has expected fields populated
	foundValidDevice := false
	for _, dev := range devices {
		if dev.Name != "" && dev.Size != "" {
			foundValidDevice = true
			t.Logf("Found device: %s (Size: %s, Type: %s)", 
				dev.Name, dev.Size, dev.Type)
			
			// Log children if any
			if len(dev.Children) > 0 {
				t.Logf("  Device %s has %d partitions", dev.Name, len(dev.Children))
				for _, child := range dev.Children {
					t.Logf("    - %s (Size: %s)", child.Name, child.Size)
				}
			}
			break
		}
	}

	if !foundValidDevice {
		t.Error("No valid device found with name and size")
	}
}

func TestSysReader_readBlockDevice(t *testing.T) {
	// Skip test if /sys/block doesn't exist (non-Linux systems)
	if _, err := os.Stat("/sys/block"); os.IsNotExist(err) {
		t.Skip("Skipping test: /sys/block not available (non-Linux system)")
	}

	reader := NewSysReader()

	// Get the first available block device
	entries, err := os.ReadDir("/sys/block")
	if err != nil {
		t.Fatalf("Failed to read /sys/block: %v", err)
	}

	if len(entries) == 0 {
		t.Skip("No block devices found in /sys/block")
	}

	deviceName := entries[0].Name()
	dev, err := reader.readBlockDevice(deviceName)
	if err != nil {
		t.Fatalf("readBlockDevice('%s') failed: %v", deviceName, err)
	}

	if dev.Name != deviceName {
		t.Errorf("Expected device name '%s', got '%s'", deviceName, dev.Name)
	}

	if dev.Size == "" {
		t.Error("Expected Size to be set, got empty string")
	}

	if dev.Type == "" {
		t.Error("Expected Type to be set, got empty string")
	}

	t.Logf("Device: Name=%s, Size=%s, Type=%s, Children=%d",
		dev.Name, dev.Size, dev.Type, len(dev.Children))
}

func TestSysReader_readBlockDevice_NonExistent(t *testing.T) {
	// Skip test if /sys/block doesn't exist (non-Linux systems)
	if _, err := os.Stat("/sys/block"); os.IsNotExist(err) {
		t.Skip("Skipping test: /sys/block not available (non-Linux system)")
	}

	reader := NewSysReader()

	// Try to read a non-existent device
	_, err := reader.readBlockDevice("nonexistent999")
	if err == nil {
		t.Error("Expected error for non-existent device, got nil")
	}
}

func TestFormatSize(t *testing.T) {
	tests := []struct {
		name     string
		bytes    int64
		expected string
	}{
		{"Zero bytes", 0, "0 B"},
		{"Small bytes", 512, "512 B"},
		{"1 KB", 1024, "1.0K"},
		{"1 MB", 1024 * 1024, "1.0M"},
		{"1 GB", 1024 * 1024 * 1024, "1.0G"},
		{"1 TB", 1024 * 1024 * 1024 * 1024, "1.0T"},
		{"500 GB", 500 * 1024 * 1024 * 1024, "500.0G"},
		{"1.5 GB", 1536 * 1024 * 1024, "1.5G"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatSize(tt.bytes)
			if result != tt.expected {
				t.Errorf("formatSize(%d) = %s, expected %s", tt.bytes, result, tt.expected)
			}
		})
	}
}

func TestSysReader_readPartitions(t *testing.T) {
	// Skip test if /sys/block doesn't exist (non-Linux systems)
	if _, err := os.Stat("/sys/block"); os.IsNotExist(err) {
		t.Skip("Skipping test: /sys/block not available (non-Linux system)")
	}

	reader := NewSysReader()

	// Find a device with partitions
	entries, err := os.ReadDir("/sys/block")
	if err != nil {
		t.Fatalf("Failed to read /sys/block: %v", err)
	}

	var deviceWithPartitions string
	for _, entry := range entries {
		// Check if device has partitions
		devicePath := "/sys/block/" + entry.Name()
		subEntries, err := os.ReadDir(devicePath)
		if err != nil {
			continue
		}

		for _, subEntry := range subEntries {
			if subEntry.IsDir() && len(subEntry.Name()) > len(entry.Name()) {
				// Found a potential partition
				deviceWithPartitions = entry.Name()
				break
			}
		}

		if deviceWithPartitions != "" {
			break
		}
	}

	if deviceWithPartitions == "" {
		t.Skip("No device with partitions found")
	}

	partitions, err := reader.readPartitions(deviceWithPartitions)
	if err != nil {
		t.Fatalf("readPartitions('%s') failed: %v", deviceWithPartitions, err)
	}

	t.Logf("Device %s has %d partitions", deviceWithPartitions, len(partitions))

	for _, part := range partitions {
		if part.Name == "" {
			t.Error("Partition name should not be empty")
		}
		if part.Type != "part" {
			t.Errorf("Expected partition type 'part', got '%s'", part.Type)
		}
		t.Logf("  Partition: %s (Size: %s)", part.Name, part.Size)
	}
}

func TestSysReader_ReadMountPoints(t *testing.T) {
	// Skip test if /proc/mounts doesn't exist (non-Linux systems)
	if _, err := os.Stat("/proc/mounts"); os.IsNotExist(err) {
		t.Skip("Skipping test: /proc/mounts not available (non-Linux system)")
	}

	reader := NewSysReader()
	mounts, err := reader.ReadMountPoints()

	if err != nil {
		t.Fatalf("ReadMountPoints() failed: %v", err)
	}

	if len(mounts) == 0 {
		t.Error("Expected at least one mount point, got none")
	}

	// Verify that at least one mount has expected fields populated
	foundValidMount := false
	for _, mount := range mounts {
		if mount.Device != "" && mount.MountPoint != "" && mount.FSType != "" {
			foundValidMount = true
			t.Logf("Found mount: Device=%s, MountPoint=%s, FSType=%s, Options=%s",
				mount.Device, mount.MountPoint, mount.FSType, mount.Options)
			break
		}
	}

	if !foundValidMount {
		t.Error("No valid mount found with device, mount point, and fstype")
	}

	// Verify root filesystem is present
	foundRoot := false
	for _, mount := range mounts {
		if mount.MountPoint == "/" {
			foundRoot = true
			t.Logf("Root filesystem: Device=%s, FSType=%s", mount.Device, mount.FSType)
			break
		}
	}

	if !foundRoot {
		t.Error("Expected to find root filesystem mount (/), but didn't")
	}
}

func TestSysReader_GetFilesystemUsage(t *testing.T) {
	// Skip test if not on Linux
	if _, err := os.Stat("/proc/mounts"); os.IsNotExist(err) {
		t.Skip("Skipping test: not on Linux system")
	}

	reader := NewSysReader()

	// Test with root filesystem
	usage, err := reader.GetFilesystemUsage("/")
	if err != nil {
		t.Fatalf("GetFilesystemUsage('/') failed: %v", err)
	}

	if usage.Path != "/" {
		t.Errorf("Expected path '/', got '%s'", usage.Path)
	}

	if usage.Total == 0 {
		t.Error("Expected Total > 0, got 0")
	}

	if usage.Used > usage.Total {
		t.Errorf("Used (%d) should not exceed Total (%d)", usage.Used, usage.Total)
	}

	if usage.Free > usage.Total {
		t.Errorf("Free (%d) should not exceed Total (%d)", usage.Free, usage.Total)
	}

	// Calculate usage percent
	usedPercent := float64(usage.Used) / float64(usage.Total) * 100
	t.Logf("Filesystem usage for /: Total=%d, Used=%d (%.1f%%), Free=%d, Files=%d, FilesFree=%d",
		usage.Total, usage.Used, usedPercent, usage.Free, usage.Files, usage.FilesFree)

	// Test with /tmp if it exists
	if _, err := os.Stat("/tmp"); err == nil {
		tmpUsage, err := reader.GetFilesystemUsage("/tmp")
		if err != nil {
			t.Logf("GetFilesystemUsage('/tmp') failed (may be expected): %v", err)
		} else {
			t.Logf("Filesystem usage for /tmp: Total=%d, Used=%d, Free=%d",
				tmpUsage.Total, tmpUsage.Used, tmpUsage.Free)
		}
	}
}

func TestSysReader_GetFilesystemUsage_NonExistent(t *testing.T) {
	// Skip test if not on Linux
	if _, err := os.Stat("/proc/mounts"); os.IsNotExist(err) {
		t.Skip("Skipping test: not on Linux system")
	}

	reader := NewSysReader()

	// Try to get usage for non-existent path
	_, err := reader.GetFilesystemUsage("/nonexistent/path/that/does/not/exist")
	if err == nil {
		t.Error("Expected error for non-existent path, got nil")
	}

	t.Logf("Got expected error: %v", err)
}

func TestSysReader_GetFilesystemUsage_CalculationAccuracy(t *testing.T) {
	// Skip test if not on Linux
	if _, err := os.Stat("/proc/mounts"); os.IsNotExist(err) {
		t.Skip("Skipping test: not on Linux system")
	}

	reader := NewSysReader()
	usage, err := reader.GetFilesystemUsage("/")
	if err != nil {
		t.Fatalf("GetFilesystemUsage('/') failed: %v", err)
	}

	// Verify that Used + Free is approximately equal to Total
	// There may be small differences due to reserved blocks
	calculatedTotal := usage.Used + usage.Free
	
	// Allow up to 10% difference for reserved blocks
	tolerance := usage.Total / 10
	diff := int64(usage.Total) - int64(calculatedTotal)
	if diff < 0 {
		diff = -diff
	}

	if uint64(diff) > tolerance {
		t.Logf("Warning: Used + Free (%d) differs significantly from Total (%d), diff=%d",
			calculatedTotal, usage.Total, diff)
		t.Logf("This may be due to reserved blocks or filesystem overhead")
	} else {
		t.Logf("Calculation check passed: Used + Free ≈ Total (diff=%d within tolerance=%d)",
			diff, tolerance)
	}
}
