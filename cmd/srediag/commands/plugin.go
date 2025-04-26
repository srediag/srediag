// Package commands provides the command-line interface for the SREDIAG application.
package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/builder"
)

// NewPluginCmd creates a new command for plugin management
func NewPluginCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugin",
		Short: "Plugin management commands",
		Long: `Commands for managing SREDIAG plugins including code generation
and plugin lifecycle management.`,
	}

	cmd.AddCommand(NewPluginGenerateCmd(opts))
	return cmd
}

// NewPluginGenerateCmd creates a new command for generating plugin code
func NewPluginGenerateCmd(opts *Options) *cobra.Command {
	var (
		configPath string
		outputDir  string
	)

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate plugin code from otelcol-builder.yaml",
		Long: `Generate plugin code based on the configuration specified in
otelcol-builder.yaml. The generated code will be written to the
specified output directory.

The generator creates:
- Plugin interface implementations
- Component factories
- Registration code`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create plugin builder
			pluginBuilder := builder.NewPluginBuilder(
				opts.Settings.GetLogger(),
				configPath,
				outputDir,
			)

			// Generate plugin code
			if err := pluginBuilder.GenerateAll(); err != nil {
				return fmt.Errorf("failed to generate plugin code: %w", err)
			}

			opts.Settings.GetLogger().Info("Plugin code generation completed",
				zap.String("config", configPath),
				zap.String("output_dir", outputDir),
			)
			return nil
		},
	}

	// Add flags
	cmd.Flags().StringVar(&configPath, "config", "otelcol-builder.yaml",
		"Path to the OpenTelemetry Collector builder configuration file")
	cmd.Flags().StringVar(&outputDir, "output-dir", "",
		"Output directory for the generated plugin code")

	// Mark output-dir as required
	if err := cmd.MarkFlagRequired("output-dir"); err != nil {
		opts.Settings.GetLogger().Error("Failed to mark flag as required",
			zap.String("flag", "output-dir"),
			zap.Error(err),
		)
	}

	return cmd
}
