// Package commands provides the command-line interface for the SREDIAG application.
package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/srediag/srediag/internal/core"
)

// OutputFormat standardizes command output formats
// Only CLI wiring and context setup should be present in this file.
type OutputFormat struct {
	Format     string // json, yaml, table
	Quiet      bool   // only output essential information
	NoColor    bool   // disable color output
	OutputFile string // file to write output to
}

// NewRootCommand creates and returns the root command for SREDIAG CLI
func NewRootCommand(ctx *core.AppContext) *cobra.Command {
	var printConfig bool

	cmd := &cobra.Command{
		Use:   "srediag",
		Short: "SREDIAG - System Resource and Environment Diagnostics",
		Long: `SREDIAG is a modular diagnostic and analysis platform designed for
comprehensive system monitoring and automated analysis.

It provides a flexible plugin architecture for extending monitoring capabilities
and integrates with various observability backends.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// 1. Bind Viper to Cobra flags
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				return fmt.Errorf("failed to bind flags: %w", err)
			}
			if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
				return fmt.Errorf("failed to bind persistent flags: %w", err)
			}

			// 2. Bind unique env vars for root config
			if err := viper.BindEnv("srediag.config", "SREDIAG_CONFIG"); err != nil {
				return fmt.Errorf("failed to bind env SREDIAG_CONFIG: %w", err)
			}

			// 3. Load config with overlays (flags > env > YAML > built-ins)
			var config core.Config
			if err := core.LoadConfigWithOverlay(&config, viperAllSettings()); err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			if err := core.ValidateConfig(&config); err != nil {
				return fmt.Errorf("invalid config: %w", err)
			}
			ctx.Config = &config

			// 4. Bootstrap logger
			logger, err := core.NewLogger(&core.Logger{
				Level:            viper.GetString("log-level"),
				Format:           viper.GetString("log-format"),
				OutputPaths:      []string{"stdout"},
				ErrorOutputPaths: []string{"stderr"},
			})
			if err != nil {
				return fmt.Errorf("failed to initialize logger: %w", err)
			}
			ctx.Logger = logger
			ctx.BuildInfo = core.DefaultBuildInfo
			ctx.ComponentManager = core.NewComponentManager(logger)

			// 5. Print config and exit if requested
			if printConfig {
				if err := core.PrintEffectiveConfig(&config); err != nil {
					return fmt.Errorf("failed to print config: %w", err)
				}
				os.Exit(0)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("please specify a subcommand, use --help for more information")
		},
		SilenceUsage: true,
	}

	// Add persistent flags (all global/common options)
	cmd.PersistentFlags().String("config", "", "path to SREDIAG configuration file (env: SREDIAG_CONFIG)")
	cmd.PersistentFlags().String("output", "table", "output format (json, yaml, table)")
	cmd.PersistentFlags().Bool("quiet", false, "only output essential information")
	cmd.PersistentFlags().Bool("no-color", false, "disable color output")
	cmd.PersistentFlags().String("output-file", "", "write output to file")
	cmd.PersistentFlags().String("log-level", "", "set log level (env: SREDIAG_LOG_LEVEL, config: log_level)")
	cmd.PersistentFlags().String("log-format", "", "set log format: json|console (env: SREDIAG_LOG_FORMAT, config: log_format)")
	cmd.PersistentFlags().BoolVar(&printConfig, "print-config", false, "print the effective merged config and exit")

	if err := viper.BindPFlag("srediag.config", cmd.PersistentFlags().Lookup("config")); err != nil {
		return nil // or handle/log the error as appropriate for your context
	}

	// Add subcommands
	cmd.AddCommand(
		NewBuildCmd(ctx),
		newDiagnoseCmd(ctx),
		newPluginCmd(ctx),
		NewServiceCmd(ctx),
	)

	return cmd
}

// viperAllSettings returns a map of all viper settings for overlay.
func viperAllSettings() map[string]string {
	settings := make(map[string]string)
	for k, v := range viper.AllSettings() {
		if s, ok := v.(string); ok {
			settings[k] = s
		}
	}
	return settings
}

// Execute creates the root command with the given context and executes it
func Execute(ctx *core.AppContext) error {
	rootCmd := NewRootCommand(ctx)
	return rootCmd.Execute()
}
