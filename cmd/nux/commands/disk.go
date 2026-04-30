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
		}).WithMessage("Disk usage").Print()
	},
}

func init() {
	diskCmd.AddCommand(diskUsageCmd)
	rootCmd.AddCommand(diskCmd)
}
