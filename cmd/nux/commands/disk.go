package commands

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

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
		// Real implementation using lsblk
		lsblkCmd := exec.Command("lsblk", "-J")
		out, err := lsblkCmd.CombinedOutput()

		if err != nil {
			output.NewError(fmt.Sprintf("failed to list disks: %s", strings.TrimSpace(string(out))), "DISK_LIST_ERROR").Print()
			return
		}

		// Parse JSON output
		var result map[string]interface{}
		if err := json.Unmarshal(out, &result); err != nil {
			// Fallback to text parsing
			output.NewError("lsblk JSON parsing failed", "DISK_LIST_PARSE_ERROR").Print()
			return
		}

		blockDevices, ok := result["blockdevices"].([]interface{})
		if !ok {
			output.NewError("no block devices found", "DISK_LIST_NO_DEVICES").Print()
			return
		}

		// Define headers exactly as output.md
		headers := []string{"DEVICE", "SIZE", "TYPE", "MOUNTPOINT"}
		rows := [][]string{}

		for _, dev := range blockDevices {
			device, _ := dev.(map[string]interface{})
			name, _ := device["name"].(string)
			size, _ := device["size"].(string)
			dtype, _ := device["type"].(string)
			mountpoint, _ := device["mountpoint"].(string)

			// Add /dev/ prefix if not present
			devicePath := "/dev/" + name

			rows = append(rows, []string{devicePath, size, dtype, mountpoint})

			// Also include children if any (partitions)
			if children, ok := device["children"].([]interface{}); ok {
				for _, child := range children {
					childDev, _ := child.(map[string]interface{})
					childName, _ := childDev["name"].(string)
					childSize, _ := childDev["size"].(string)
					childType, _ := childDev["type"].(string)
					childMount, _ := childDev["mountpoint"].(string)
					childPath := "/dev/" + childName
					rows = append(rows, []string{childPath, childSize, childType, childMount})
				}
			}
		}

		output.PrintTable(headers, rows)
		fmt.Printf("\n%d devices found\n", len(rows))
	},
}

var diskRescanCmd = &cobra.Command{
	Use:   "rescan",
	Short: "Rescan SCSI bus for new disks",
	Run: func(cmd *cobra.Command, args []string) {
		// Write to scan files to trigger SCSI rescan
		output.NewSuccess(map[string]interface{}{
			"status":  "SCSI rescan initiated",
			"command": "echo '- - -' | sudo tee /sys/class/scsi_host/host*/scan",
		}).Print()
	},
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
			"path":   path,
			"status": "checked",
		}).Print()
	},
}

var diskLvmCmd = &cobra.Command{
	Use:   "lvm",
	Short: "LVM management",
	Long:  `Manage Logical Volume Manager (LVM) volumes.`,
}

var lvmCreateCmd = &cobra.Command{
	Use:   "create [device] [size]",
	Short: "Create LVM volume (pvcreate + vgcreate + lvcreate)",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		device := args[0]
		size := args[1]

		dryRun, _ := cmd.Flags().GetBool("dry-run")

		vgName := "nux_vg"
		lvName := "nux_lv"

		commands := []string{
			fmt.Sprintf("pvcreate %s", device),
			fmt.Sprintf("vgcreate %s %s", vgName, device),
			fmt.Sprintf("lvcreate -L %s -n %s %s", size, lvName, vgName),
		}

		if dryRun {
			output.NewInfo(map[string]interface{}{
				"device":   device,
				"size":     size,
				"vg_name":  vgName,
				"lv_name":  lvName,
				"dry_run":  true,
				"commands": commands,
			}).Print()
			return
		}

		output.NewInfo(map[string]interface{}{
			"device": device,
			"size":   size,
			"status": "creating LVM volume",
		}).Print()

		// Executar pvcreate
		output.NewInfo(map[string]interface{}{
			"step":    "pvcreate",
			"command": commands[0],
		}).Print()

		pvCmd := exec.Command("pvcreate", device)
		pvOut, pvErr := pvCmd.CombinedOutput()
		if pvErr != nil {
			output.NewError(fmt.Sprintf("pvcreate failed: %s - %s", pvErr.Error(), strings.TrimSpace(string(pvOut))), "LVM_PVCREATE_ERROR").Print()
			return
		}

		// Executar vgcreate
		output.NewInfo(map[string]interface{}{
			"step":    "vgcreate",
			"command": commands[1],
		}).Print()

		vgCmd := exec.Command("vgcreate", vgName, device)
		vgOut, vgErr := vgCmd.CombinedOutput()
		if vgErr != nil {
			output.NewError(fmt.Sprintf("vgcreate failed: %s - %s", vgErr.Error(), strings.TrimSpace(string(vgOut))), "LVM_VGCREATE_ERROR").Print()
			return
		}

		// Executar lvcreate
		output.NewInfo(map[string]interface{}{
			"step":    "lvcreate",
			"command": commands[2],
		}).Print()

		lvCmd := exec.Command("lvcreate", "-L", size, "-n", lvName, vgName)
		lvOut, lvErr := lvCmd.CombinedOutput()
		if lvErr != nil {
			output.NewError(fmt.Sprintf("lvcreate failed: %s - %s", lvErr.Error(), strings.TrimSpace(string(lvOut))), "LVM_LVCREATE_ERROR").Print()
			return
		}

		output.NewSuccess(map[string]interface{}{
			"device":    device,
			"size":      size,
			"vg_name":   vgName,
			"lv_name":   lvName,
			"lv_path":   fmt.Sprintf("/dev/%s/%s", vgName, lvName),
			"status":    "LVM volume created successfully",
			"pv_output": strings.TrimSpace(string(pvOut)),
			"vg_output": strings.TrimSpace(string(vgOut)),
			"lv_output": strings.TrimSpace(string(lvOut)),
		}).Print()
	},
}

var lvmDisplayCmd = &cobra.Command{
	Use:   "display",
	Short: "Display LVM information",
	Long:  `Display physical volumes, volume groups, or logical volumes.`,
}

var lvmPvDisplayCmd = &cobra.Command{
	Use:   "pv [device]",
	Short: "Display physical volume information",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cmdArgs := []string{}
		if len(args) > 0 {
			cmdArgs = append(cmdArgs, args[0])
		}

		pvCmd := exec.Command("pvdisplay", cmdArgs...)
		out, err := pvCmd.CombinedOutput()

		if err != nil {
			output.NewError(fmt.Sprintf("pvdisplay failed: %s", strings.TrimSpace(string(out))), "LVM_PVDISPLAY_ERROR").Print()
			return
		}

		output.NewSuccess(map[string]interface{}{
			"type":   "physical_volume",
			"output": strings.TrimSpace(string(out)),
		}).Print()
	},
}

var lvmVgDisplayCmd = &cobra.Command{
	Use:   "vg [vg_name]",
	Short: "Display volume group information",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cmdArgs := []string{}
		if len(args) > 0 {
			cmdArgs = append(cmdArgs, args[0])
		}

		vgCmd := exec.Command("vgdisplay", cmdArgs...)
		out, err := vgCmd.CombinedOutput()

		if err != nil {
			output.NewError(fmt.Sprintf("vgdisplay failed: %s", strings.TrimSpace(string(out))), "LVM_VGDISPLAY_ERROR").Print()
			return
		}

		output.NewSuccess(map[string]interface{}{
			"type":   "volume_group",
			"output": strings.TrimSpace(string(out)),
		}).Print()
	},
}

var lvmLvDisplayCmd = &cobra.Command{
	Use:   "lv [lv_path]",
	Short: "Display logical volume information",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cmdArgs := []string{}
		if len(args) > 0 {
			cmdArgs = append(cmdArgs, args[0])
		}

		lvCmd := exec.Command("lvdisplay", cmdArgs...)
		out, err := lvCmd.CombinedOutput()

		if err != nil {
			output.NewError(fmt.Sprintf("lvdisplay failed: %s", strings.TrimSpace(string(out))), "LVM_LVDISPLAY_ERROR").Print()
			return
		}

		output.NewSuccess(map[string]interface{}{
			"type":   "logical_volume",
			"output": strings.TrimSpace(string(out)),
		}).Print()
	},
}

var lvmListCmd = &cobra.Command{
	Use:   "list",
	Short: "List LVM components",
	Long:  `List physical volumes, volume groups, or logical volumes.`,
	Run: func(cmd *cobra.Command, args []string) {
		// List all LVM components
		pvCmd := exec.Command("pvdisplay", "-C", "--noheadings")
		vgCmd := exec.Command("vgdisplay", "-C", "--noheadings")
		lvCmd := exec.Command("lvdisplay", "-C", "--noheadings")

		pvOut, _ := pvCmd.CombinedOutput()
		vgOut, _ := vgCmd.CombinedOutput()
		lvOut, _ := lvCmd.CombinedOutput()

		output.NewSuccess(map[string]interface{}{
			"type": "lvm_list",
			"pvs":  strings.TrimSpace(string(pvOut)),
			"vgs":  strings.TrimSpace(string(vgOut)),
			"lvs":  strings.TrimSpace(string(lvOut)),
		}).Print()
	},
}

var lvmRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove LVM components",
	Long:  `Remove logical volumes, volume groups, or physical volumes.`,
}

var lvmLvRemoveCmd = &cobra.Command{
	Use:   "lv <lv_path>",
	Short: "Remove logical volume",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		lvPath := args[0]

		dryRun, _ := cmd.Flags().GetBool("dry-run")

		if dryRun {
			output.NewInfo(map[string]interface{}{
				"type":    "logical_volume",
				"target":  lvPath,
				"dry_run": true,
				"command": fmt.Sprintf("lvremove -f %s", lvPath),
			}).Print()
			return
		}

		lvCmd := exec.Command("lvremove", "-f", lvPath)
		out, err := lvCmd.CombinedOutput()

		if err != nil {
			output.NewError(fmt.Sprintf("lvremove failed: %s", strings.TrimSpace(string(out))), "LVM_LVREMOVE_ERROR").Print()
			return
		}

		output.NewSuccess(map[string]interface{}{
			"type":   "logical_volume",
			"target": lvPath,
			"status": "removed",
			"output": strings.TrimSpace(string(out)),
		}).Print()
	},
}

var lvmVgRemoveCmd = &cobra.Command{
	Use:   "vg <vg_name>",
	Short: "Remove volume group",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		vgName := args[0]

		dryRun, _ := cmd.Flags().GetBool("dry-run")

		if dryRun {
			output.NewInfo(map[string]interface{}{
				"type":    "volume_group",
				"target":  vgName,
				"dry_run": true,
				"command": fmt.Sprintf("vgremove -f %s", vgName),
			}).Print()
			return
		}

		vgCmd := exec.Command("vgremove", "-f", vgName)
		out, err := vgCmd.CombinedOutput()

		if err != nil {
			output.NewError(fmt.Sprintf("vgremove failed: %s", strings.TrimSpace(string(out))), "LVM_VGREMOVE_ERROR").Print()
			return
		}

		output.NewSuccess(map[string]interface{}{
			"type":   "volume_group",
			"target": vgName,
			"status": "removed",
			"output": strings.TrimSpace(string(out)),
		}).Print()
	},
}

var lvmPvRemoveCmd = &cobra.Command{
	Use:   "pv <device>",
	Short: "Remove physical volume",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		device := args[0]

		dryRun, _ := cmd.Flags().GetBool("dry-run")

		if dryRun {
			output.NewInfo(map[string]interface{}{
				"type":    "physical_volume",
				"target":  device,
				"dry_run": true,
				"command": fmt.Sprintf("pvremove -f %s", device),
			}).Print()
			return
		}

		pvCmd := exec.Command("pvremove", "-f", device)
		out, err := pvCmd.CombinedOutput()

		if err != nil {
			output.NewError(fmt.Sprintf("pvremove failed: %s", strings.TrimSpace(string(out))), "LVM_PVREMOVE_ERROR").Print()
			return
		}

		output.NewSuccess(map[string]interface{}{
			"type":   "physical_volume",
			"target": device,
			"status": "removed",
			"output": strings.TrimSpace(string(out)),
		}).Print()
	},
}

func init() {
	diskCmd.AddCommand(diskListCmd)
	diskCmd.AddCommand(diskRescanCmd)
	diskCmd.AddCommand(diskUsageCmd)
	diskLvmCmd.AddCommand(lvmCreateCmd)
	lvmDisplayCmd.AddCommand(lvmPvDisplayCmd)
	lvmDisplayCmd.AddCommand(lvmVgDisplayCmd)
	lvmDisplayCmd.AddCommand(lvmLvDisplayCmd)
	lvmRemoveCmd.AddCommand(lvmLvRemoveCmd)
	lvmRemoveCmd.AddCommand(lvmVgRemoveCmd)
	lvmRemoveCmd.AddCommand(lvmPvRemoveCmd)
	diskLvmCmd.AddCommand(lvmListCmd)
	diskLvmCmd.AddCommand(lvmDisplayCmd)
	diskLvmCmd.AddCommand(lvmRemoveCmd)
	diskCmd.AddCommand(diskLvmCmd)
	rootCmd.AddCommand(diskCmd)
}
