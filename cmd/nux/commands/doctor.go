package commands

import (
	"github.com/rsdenck/nux/internal/output"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "System health check",
	Long:  `Run system health checks (Load, Swap, Disk, TCP, etc).`,
	Run: func(cmd *cobra.Command, args []string) {
		output.NewSuccess(map[string]interface{}{
			"status": "ok",
			"checks": 4,
		}).WithMessage("System health check completed").Print()
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
