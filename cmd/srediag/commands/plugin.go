package commands

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/builder"
)

// NewPluginCmd creates a new plugin command
func NewPluginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugin",
		Short: "Plugin management commands",
	}

	cmd.AddCommand(NewPluginGenerateCmd())
	return cmd
}

// NewPluginGenerateCmd creates a new plugin generate command
func NewPluginGenerateCmd() *cobra.Command {
	var (
		configPath string
		outputDir  string
	)

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate plugin code from otelcol-builder.yaml",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create plugin builder
			pluginBuilder := builder.NewPluginBuilder(
				cmdSettings.GetLogger(),
				configPath,
				outputDir,
			)

			// Generate plugin code
			if err := pluginBuilder.GenerateAll(); err != nil {
				return err
			}
			return nil
		},
	}

	// Add flags
	cmd.Flags().StringVar(&configPath, "config", "otelcol-builder.yaml", "Path to otelcol-builder.yaml")
	cmd.Flags().StringVar(&outputDir, "output-dir", "", "Output directory for generated code")

	// Mark output-dir as required
	if err := cmd.MarkFlagRequired("output-dir"); err != nil {
		cmdSettings.GetLogger().Error("Failed to mark flag as required", zap.Error(err))
	}

	return cmd
}
