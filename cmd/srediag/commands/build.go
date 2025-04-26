// Package commands provides the command-line interface for the SREDIAG application.
package commands

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/builder"
)

// newBuildCmd creates a new command for building plugins and the main project
func newBuildCmd(opts *Options) *cobra.Command {
	var outputDir string

	cmd := &cobra.Command{
		Use:   "build",
		Short: "Build plugins and main project",
		Long: `Build all components of the SREDIAG project including plugins.
The build process uses the configuration from otelcol-builder.yaml and
outputs the built artifacts to the specified directory.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get workspace directory
			workDir, err := filepath.Abs(".")
			if err != nil {
				return fmt.Errorf("failed to get absolute path: %w", err)
			}

			// Create builder
			makeBuilder := builder.NewMakeBuilder(
				opts.Settings.GetLogger(),
				workDir,
				"otelcol-builder.yaml",
				outputDir,
			)

			// Build everything
			if err := makeBuilder.BuildAll(); err != nil {
				return fmt.Errorf("build failed: %w", err)
			}

			opts.Settings.GetLogger().Info("Build completed successfully",
				zap.String("output_dir", outputDir),
			)
			return nil
		},
	}

	// Add flags
	cmd.Flags().StringVar(&outputDir, "output-dir", "/tmp/srediag-plugins",
		"Output directory for built plugins and artifacts")

	return cmd
}
