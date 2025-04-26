package commands

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/srediag/srediag/internal/builder"
)

func newBuildCmd() *cobra.Command {
	var (
		outputDir string
	)

	cmd := &cobra.Command{
		Use:   "build",
		Short: "Build plugins and main project",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get workspace directory
			workDir, err := filepath.Abs(".")
			if err != nil {
				return err
			}

			// Create builder
			makeBuilder := builder.NewMakeBuilder(
				cmdSettings.GetLogger(),
				workDir,
				"otelcol-builder.yaml",
				outputDir,
			)

			// Build everything
			return makeBuilder.BuildAll()
		},
	}

	// Add flags
	cmd.Flags().StringVar(&outputDir, "output-dir", "/tmp/srediag-plugins", "Output directory for plugins")

	return cmd
}
