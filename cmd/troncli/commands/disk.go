package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/mascli/troncli/internal/core/adapter"
	"github.com/mascli/troncli/internal/core/ports"
	"github.com/mascli/troncli/internal/core/services"
	"github.com/mascli/troncli/internal/modules/disk"
	"github.com/mascli/troncli/internal/ui/console"
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

			table := console.NewBoxTable(os.Stdout)
			table.SetTitle("TRONCLI - USO DE DISCO (DEVICES)")
			table.SetHeaders([]string{"NAME", "SIZE", "TYPE", "MOUNTPOINT"})

			var addDev func(d ports.BlockDevice, indent string)
			addDev = func(d ports.BlockDevice, indent string) {
				table.AddRow([]string{indent + d.Name, d.Size, d.Type, d.MountPoint})
				for _, child := range d.Children {
					addDev(child, indent+"  ")
				}
			}
			for _, d := range devices {
				addDev(d, "")
			}
			table.Render()
		} else {
			usage, err := manager.GetFilesystemUsage(args[0])
			if err != nil {
				fmt.Printf("Error getting usage for %s: %v\n", args[0], err)
				os.Exit(1)
			}

			table := console.NewBoxTable(os.Stdout)
			table.SetTitle(fmt.Sprintf("TRONCLI - USO DE DISCO: %s", args[0]))
			table.SetHeaders([]string{"METRIC", "VALUE"})

			table.AddRow([]string{"Path", usage.Path})
			table.AddRow([]string{"Total", fmt.Sprintf("%d bytes", usage.Total)})
			table.AddRow([]string{"Used", fmt.Sprintf("%d bytes", usage.Used)})
			table.AddRow([]string{"Free", fmt.Sprintf("%d bytes", usage.Free)})
			table.AddRow([]string{"Inodes", fmt.Sprintf("%d (Free: %d)", usage.Files, usage.FilesFree)})

			table.Render()
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

		table := console.NewBoxTable(os.Stdout)
		table.SetTitle("TRONCLI - SAÚDE DO DISCO")
		table.SetHeaders([]string{"STATUS", "HEALTH"})
		table.AddRow([]string{icon, status})
		table.Render()
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

		table := console.NewBoxTable(os.Stdout)
		table.SetTitle(fmt.Sprintf("TRONCLI - TOP %d ARQUIVOS: %s", count, path))
		table.SetHeaders([]string{"SIZE", "PATH"})

		for _, f := range files {
			p := f.Path
			if len(p) > 60 {
				p = "..." + p[len(p)-57:]
			}
			table.AddRow([]string{formatBytes(f.Size), p})
		}
		table.Render()
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

		table := console.NewBoxTable(os.Stdout)
		table.SetTitle(fmt.Sprintf("TRONCLI - INODES: %s", path))
		table.SetHeaders([]string{"METRIC", "VALUE"})
		table.AddRow([]string{"Used", fmt.Sprintf("%d", used)})
		table.AddRow([]string{"Total", fmt.Sprintf("%d", total)})
		table.AddRow([]string{"Usage %", fmt.Sprintf("%.2f%%", percent)})
		table.Render()
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
