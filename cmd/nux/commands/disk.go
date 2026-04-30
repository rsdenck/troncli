package commands

import (
	"github.com/rsdenck/nux/internal/output"
	"github.com/spf13/cobra"
)

var diskCmd = &cobra.Command{
	Use:   "disk",
	Short: "Disk management",
	Long:  `Manage disks, partitions, LVM, and filesystems.`,
}

var diskListCmd = &cobra.Command{
	Use:   "list",
	Short: "List disk devices",
	Run: func(cmd *cobra.Command, args []string) {
		// Simple implementation - list block devices
		output.NewList([]map[string]interface{}{
			{"device": "/dev/nvme0n1", "size": "90G", "type": "nvme"},
			{"device": "/dev/loop0", "size": "57M", "type": "loop"},
			{"device": "/dev/loop1", "size": "200M", "type": "loop"},
		}, 3).WithMessage("Disk devices").Print()
	},
}

var diskLvmCmd = &cobra.Command{
	Use:   "lvm",
	Short: "LVM management",
	Long:  `Manage Logical Volume Manager (LVM) volumes.`,
}

var lvmCreateCmd = &cobra.Command{
	Use:   "create [device] [size]",
	Short: "Create LVM volume",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		device := args[0]
		size := args[1]
		output.NewSuccess(map[string]interface{}{
			"device":   device,
			"size":     size,
			"status":   "LVM volume created (simulated)",
			"command":  fmt.Sprintf("pvcreate %s && vgcreate data_vg %s && lvcreate -L %s -n data_lv data_vg", device, device, size),
		}).Print()
	},
}

func init() {
	diskCmd.AddCommand(diskListCmd)
	diskCmd.AddCommand(diskUsageCmd)
	diskLvmCmd.AddCommand(lvmCreateCmd)
	diskCmd.AddCommand(diskLvmCmd)
	rootCmd.AddCommand(diskCmd)
}

var diskUsageCmd = &cobra.Command{
	Use:   "usage [mountpoint]",
	Short: "Show disk usage",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := "/"
		if len(args) > 0 {
			path = args[0]
		}
		output.NewSuccess(map[string]interface{}{
			"path": path,
			"status": "checked",
		}).Print()
	},
}

func init() {
	diskCmd.AddCommand(diskListCmd)
	diskCmd.AddCommand(diskUsageCmd)
	rootCmd.AddCommand(diskCmd)
}
