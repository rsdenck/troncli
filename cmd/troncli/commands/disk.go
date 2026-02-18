package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/ports"
	"github.com/mascli/troncli/internal/core/services"
	"github.com/mascli/troncli/internal/modules/disk"
	"github.com/spf13/cobra"
)

var diskCmd = &cobra.Command{
	Use:   "disk",
	Short: "Gerenciamento de Disco",
	Long:  `Gerencie discos, partições, e uso de espaço (usage, cleanup, health).`,
}

var diskUsageCmd = &cobra.Command{
	Use:   "usage [mountpoint]",
	Short: "Exibe uso do disco",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := getDiskManager()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if len(args) == 0 {
			// List block devices for overview
			devices, err := manager.ListBlockDevices()
			if err != nil {
				fmt.Printf("Error listing devices: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("%-10s %-10s %-10s %-20s\n", "NAME", "SIZE", "TYPE", "MOUNTPOINT")
			var printDev func(d ports.BlockDevice, indent string)
			printDev = func(d ports.BlockDevice, indent string) {
				fmt.Printf("%-10s %-10s %-10s %-20s\n", indent+d.Name, d.Size, d.Type, d.MountPoint)
				for _, child := range d.Children {
					printDev(child, indent+"  ")
				}
			}
			for _, d := range devices {
				printDev(d, "")
			}
		} else {
			usage, err := manager.GetFilesystemUsage(args[0])
			if err != nil {
				fmt.Printf("Error getting usage for %s: %v\n", args[0], err)
				os.Exit(1)
			}
			fmt.Printf("Path: %s\n", usage.Path)
			fmt.Printf("Total: %d bytes\n", usage.Total)
			fmt.Printf("Used:  %d bytes\n", usage.Used)
			fmt.Printf("Free:  %d bytes\n", usage.Free)
			fmt.Printf("Inodes: %d (Free: %d)\n", usage.Files, usage.FilesFree)
		}
	},
}

var diskCleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Limpa arquivos temporários e caches",
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := getDiskManager()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Cleaning up system...")
		if err := manager.Cleanup(); err != nil {
			fmt.Printf("Error during cleanup: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Cleanup completed successfully.")
	},
}

var diskHealthCmd = &cobra.Command{
	Use:   "health",
	Short: "Verifica saúde do disco",
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := getDiskManager()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		status, err := manager.GetDiskHealth()
		if err != nil {
			fmt.Printf("Error checking disk health: %v\n", err)
			os.Exit(1)
		}
		var icon string
		switch status {
		case "Warning":
			icon = "⚠️"
		case "Critical":
			icon = "❌"
		default:
			icon = "✅"
		}
		fmt.Printf("%s Disk Health: %s\n", icon, status)
	},
}

var diskTopFilesCmd = &cobra.Command{
	Use:   "top-files [path] [count]",
	Short: "Lista maiores arquivos",
	Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		path := "/"
		count := 10
		if len(args) >= 1 {
			path = args[0]
		}
		if len(args) >= 2 {
			c, err := strconv.Atoi(args[1])
			if err == nil {
				count = c
			}
		}

		manager, err := getDiskManager()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Scanning top %d files in %s...\n", count, path)
		files, err := manager.GetTopFiles(path, count)
		if err != nil {
			fmt.Printf("Error scanning files: %v\n", err)
			os.Exit(1)
		}
		for _, f := range files {
			fmt.Printf("%-10s %s\n", formatBytes(f.Size), f.Path)
		}
	},
}

var diskInodesCmd = &cobra.Command{
	Use:   "inodes [path]",
	Short: "Exibe uso de inodes",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := "/"
		if len(args) >= 1 {
			path = args[0]
		}
		manager, err := getDiskManager()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		used, total, err := manager.GetInodeUsage(path)
		if err != nil {
			fmt.Printf("Error getting inode usage: %v\n", err)
			os.Exit(1)
		}
		percent := 0.0
		if total > 0 {
			percent = (float64(used) / float64(total)) * 100
		}
		fmt.Printf("Inodes on %s: %d / %d (%.2f%%)\n", path, used, total, percent)
	},
}

func init() {
	rootCmd.AddCommand(diskCmd)
	diskCmd.AddCommand(diskUsageCmd)
	diskCmd.AddCommand(diskCleanupCmd)
	diskCmd.AddCommand(diskHealthCmd)
	diskCmd.AddCommand(diskTopFilesCmd)
	diskCmd.AddCommand(diskInodesCmd)
}

func getDiskManager() (ports.DiskManager, error) {
	executor := adapter.NewExecutor()
	profileEngine := services.NewProfileEngine(executor)
	profile, err := profileEngine.DetectProfile()
	if err != nil {
		return nil, fmt.Errorf("failed to detect system profile: %w", err)
	}

	return disk.NewUniversalDiskManager(executor, profile), nil
}

func formatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
