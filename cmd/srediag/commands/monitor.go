package commands

import (
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func newMonitorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "monitor",
		Short: "Monitor resources",
		Long: `Monitor various resources in real-time.
		
Examples:
  # Real-time system monitoring
  srediag monitor system --interval 5s`,
	}

	// Add subcommands
	cmd.AddCommand(
		newMonitorSystemCmd(),
	)

	return cmd
}

func newMonitorSystemCmd() *cobra.Command {
	var interval time.Duration

	cmd := &cobra.Command{
		Use:   "system [--interval <duration>]",
		Short: "Monitor system",
		Long:  "Monitor system resources in real-time",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := getLogger()
			logger.Info("monitoring system",
				zap.Duration("interval", interval))
			return nil
		},
	}

	cmd.Flags().DurationVar(&interval, "interval", 5*time.Second, "monitoring interval")
	return cmd
}
