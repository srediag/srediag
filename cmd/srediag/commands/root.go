package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/srediag/srediag/internal/settings"
)

var (
	// Global flags
	configPath string

	// Settings
	cmdSettings *settings.CommandSettings

	// Root command
	rootCmd = &cobra.Command{
		Use:   "srediag",
		Short: "SREDIAG - System Resource and Environment Diagnostics",
		Long: `SREDIAG is a modular diagnostic and analysis platform designed for 
comprehensive system monitoring and automated analysis.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("please specify a subcommand")
		},
	}
)

// Execute executes the root command
func Execute(s *settings.CommandSettings) error {
	cmdSettings = s

	// Add subcommands
	rootCmd.AddCommand(
		newStartCmd(),
		newVersionCmd(),
		newDiagnoseCmd(),
		newBuildCmd(),
		NewPluginCmd(),
	)

	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "/etc/srediag/config/srediag.yaml", "path to configuration file")
}
