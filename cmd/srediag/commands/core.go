// Package commands provides the command-line interface for the SREDIAG application.
package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/srediag/srediag/internal/core"
)

const (
	// Environment variables
	envConfigPath = "SREDIAG_CONFIG"
	envLogLevel   = "SREDIAG_LOG_LEVEL"
	envLogFormat  = "SREDIAG_LOG_FORMAT"
)

var defaultConfigPaths = []string{
	"/etc/srediag/config/srediag.yaml",
	"$HOME/.srediag/config.yaml",
	"./config/srediag.yaml",
	"./srediag.yaml",
}

// OutputFormat standardizes command output formats
type OutputFormat struct {
	Format     string // json, yaml, table
	Quiet      bool   // only output essential information
	NoColor    bool   // disable color output
	OutputFile string // file to write output to
}

// NewRootCommand creates and returns the root command for SREDIAG CLI
func NewRootCommand(ctx *core.AppContext) *cobra.Command {
	var outputOpts OutputFormat
	var logLevelFlag string
	var logFormatFlag string

	cmd := &cobra.Command{
		Use:   "srediag",
		Short: "SREDIAG - System Resource and Environment Diagnostics",
		Long: `SREDIAG is a modular diagnostic and analysis platform designed for
comprehensive system monitoring and automated analysis.

It provides a flexible plugin architecture for extending monitoring capabilities
and integrates with various observability backends.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Determine log level and format (priority: flag > env > config > default)
			logLevel := logLevelFlag
			if logLevel == "" {
				logLevel = os.Getenv(envLogLevel)
			}
			if logLevel == "" {
				logLevel = ctx.Config.LogLevel
			}
			if logLevel == "" {
				logLevel = "info"
			}

			logFormat := logFormatFlag
			if logFormat == "" {
				logFormat = os.Getenv(envLogFormat)
			}
			if logFormat == "" {
				logFormat = ctx.Config.LogFormat
			}
			if logFormat == "" {
				logFormat = "console"
			}

			loggerCfg := &core.Logger{
				Level:            logLevel,
				Format:           logFormat,
				OutputPaths:      []string{"stdout"},
				ErrorOutputPaths: []string{"stderr"},
				Development:      false,
			}
			logger, err := core.NewLogger(loggerCfg)
			if err != nil {
				return fmt.Errorf("failed to create logger: %w", err)
			}
			ctx.Logger = logger

			// Find config file
			configPath := findConfigFile(ctx.GetConfig().PluginsDir)
			if configPath == "" {
				ctx.Logger.Warn("Config file not found in any of the default locations, using defaults")
			} else {
				if err := ctx.Config.Load(configPath); err != nil {
					ctx.Logger.Error(fmt.Sprintf("Failed to load config file: path=%s, err=%v", configPath, err))
					return err
				}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("please specify a subcommand, use --help for more information")
		},
		SilenceUsage: true,
	}

	// Add persistent flags
	cmd.PersistentFlags().StringVar(&ctx.Config.PluginsDir, "config", ctx.Config.PluginsDir, "path to configuration file (env: SREDIAG_CONFIG)")
	cmd.PersistentFlags().StringVar(&outputOpts.Format, "output", "table", "output format (json, yaml, table)")
	cmd.PersistentFlags().BoolVar(&outputOpts.Quiet, "quiet", false, "only output essential information")
	cmd.PersistentFlags().BoolVar(&outputOpts.NoColor, "no-color", false, "disable color output")
	cmd.PersistentFlags().StringVar(&outputOpts.OutputFile, "output-file", "", "write output to file")
	cmd.PersistentFlags().StringVar(&logLevelFlag, "log-level", "", "set log level (env: SREDIAG_LOG_LEVEL, config: log_level)")
	cmd.PersistentFlags().StringVar(&logFormatFlag, "log-format", "", "set log format: json|console (env: SREDIAG_LOG_FORMAT, config: log_format)")

	// Add subcommands
	cmd.AddCommand(
		newStartCmd(ctx),
		core.NewVersionCmd(),
		newDiagnoseCmd(ctx),
		NewBuildCmd(ctx),
		newPluginCmd(ctx),
	)

	return cmd
}

// findConfigFile searches for the config file in default locations
func findConfigFile(configPath string) string {
	// Check explicit config path first
	if configPath != "" {
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}
	}

	// Check environment variable
	if envPath := os.Getenv(envConfigPath); envPath != "" {
		if _, err := os.Stat(envPath); err == nil {
			return envPath
		}
	}

	// Expand $HOME in default paths
	home, _ := os.UserHomeDir()

	// Check default locations
	for _, path := range defaultConfigPaths {
		expandedPath := os.ExpandEnv(path)
		expandedPath = filepath.Clean(filepath.Join(filepath.Dir(expandedPath), filepath.Base(expandedPath)))
		if strings.Contains(expandedPath, "$HOME") {
			expandedPath = strings.ReplaceAll(expandedPath, "$HOME", home)
		}
		if _, err := os.Stat(expandedPath); err == nil {
			return expandedPath
		}
	}

	return ""
}

// Execute creates the root command with the given context and executes it
func Execute(ctx *core.AppContext) error {
	rootCmd := NewRootCommand(ctx)
	return rootCmd.Execute()
}
