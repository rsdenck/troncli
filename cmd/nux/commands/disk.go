package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rsdenck/nux/internal/core"
	"github.com/rsdenck/nux/internal/output"
	"github.com/spf13/cobra"
)

var diskExecutor core.Executor = &core.RealExecutor{}

var diskCmd = &cobra.Command{
	Use:   "disk",
	Short: "Disk management",
	Long:  `Manage disks, partitions, LVM, and filesystems.`,
}

var diskListCmd = &cobra.Command{
	Use:   "list",
	Short: "List disk devices",
	Run: func(cmd *cobra.Command, args []string) {
		out, err := diskExecutor.CombinedOutput("lsblk", "-J")
		if err != nil {
			output.NewError(fmt.Sprintf("failed to list disks: %s", err.Error()), "DISK_LIST_ERROR").Print()
			return
		}

		out = core.SanitizeInput(out)

		var result map[string]interface{}
		if err := json.Unmarshal([]byte(out), &result); err != nil {
			output.NewError("lsblk JSON parsing failed", "DISK_LIST_PARSE_ERROR").Print()
			return
		}

		blockDevices, ok := result["blockdevices"].([]interface{})
		if !ok {
			output.NewError("no block devices found", "DISK_LIST_NO_DEVICES").Print()
			return
		}

		headers := []string{"DEVICE", "SIZE", "TYPE", "MOUNTPOINT"}
		rows := [][]string{}

		for _, dev := range blockDevices {
			device := dev.(map[string]interface{})
			name := device["name"].(string)
			size := device["size"].(string)
			dtype := device["type"].(string)
			mountpoint := ""
			if mp, ok := device["mountpoint"].(string); ok {
				mountpoint = mp
			}

			devicePath := "/dev/" + name
			rows = append(rows, []string{devicePath, size, dtype, mountpoint})

			if children, ok := device["children"].([]interface{}); ok {
				for _, child := range children {
					childDev := child.(map[string]interface{})
					childName := childDev["name"].(string)
					childSize := childDev["size"].(string)
					childType := childDev["type"].(string)
					childMount := ""
					if cm, ok := childDev["mountpoint"].(string); ok {
						childMount = cm
					}
					childPath := "/dev/" + childName
					rows = append(rows, []string{childPath, childSize, childType, childMount})
				}
			}
		}

		output.PrintTable(headers, rows)
		fmt.Printf("\n%d devices found\n", len(rows))
	},
}

var diskUsageCmd = &cobra.Command{
	Use:   "usage [mountpoint]",
	Short: "Show disk usage",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := "/"
		if len(args) > 0 {
			path = core.SanitizeInput(args[0])
			if !core.ValidatePath(path) {
				output.NewError(fmt.Sprintf("invalid path: %s", path), "DISK_USAGE_INVALID_PATH").Print()
				return
			}
		}

		out, err := diskExecutor.CombinedOutput("df", "-h", path)
		if err != nil {
			output.NewError(fmt.Sprintf("failed to get disk usage: %s", err.Error()), "DISK_USAGE_ERROR").Print()
			return
		}

		lines := strings.Split(out, "\n")
		if len(lines) < 2 {
			output.NewError("no disk usage data", "DISK_USAGE_NO_DATA").Print()
			return
		}

		headers := []string{"FILESYSTEM", "SIZE", "USED", "AVAIL", "USE%", "MOUNTED"}
		rows := [][]string{}

		for i, line := range lines {
			if i == 0 || strings.TrimSpace(line) == "" {
				continue
			}
			fields := strings.Fields(line)
			if len(fields) >= 6 {
				rows = append(rows, fields[:6])
			}
		}

		output.PrintTable(headers, rows)
	},
}

var diskRescanCmd = &cobra.Command{
	Use:   "rescan",
	Short: "Rescan SCSI bus for new disks",
	Run: func(cmd *cobra.Command, args []string) {
		err := diskExecutor.RunSilent("bash", "-c", "echo '- - -' | tee /sys/class/scsi_host/host*/scan")
		if err != nil {
			output.NewError("failed to rescan SCSI bus", "DISK_RESCAN_ERROR").Print()
			return
		}
		output.NewSuccess(map[string]interface{}{
			"status":  "SCSI rescan initiated",
			"command": "echo '- - -' | sudo tee /sys/class/scsi_host/host*/scan",
		}).Print()
	},
}

func init() {
	diskCmd.AddCommand(diskListCmd)
	diskCmd.AddCommand(diskUsageCmd)
	diskCmd.AddCommand(diskRescanCmd)
	// TODO: Add LVM commands with executor
	rootCmd.AddCommand(diskCmd)
}
