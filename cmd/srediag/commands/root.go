// Package commands provides the command-line interface for the SREDIAG application.
package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/srediag/srediag/internal/settings"
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

// Options holds command-line options that are shared across commands
type Options struct {
	ConfigPath string
	Settings   *settings.CommandSettings
	LogConfig  LogConfig
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level  string // debug, info, warn, error
	Format string // json, console
}

// OutputFormat standardizes command output formats
type OutputFormat struct {
	Format     string // json, yaml, table
	Quiet      bool   // only output essential information
	NoColor    bool   // disable color output
	OutputFile string // file to write output to
}

// NewRootCommand creates and returns the root command for SREDIAG CLI
func NewRootCommand(opts *Options) *cobra.Command {
	if opts == nil {
		opts = &Options{
			LogConfig: LogConfig{
				Level:  "info",
				Format: "console",
			},
		}
	}

	var outputOpts OutputFormat

	cmd := &cobra.Command{
		Use:   "srediag",
		Short: "SREDIAG - System Resource and Environment Diagnostics",
		Long: `SREDIAG is a modular diagnostic and analysis platform designed for 
comprehensive system monitoring and automated analysis.

It provides a flexible plugin architecture for extending monitoring capabilities
and integrates with various observability backends.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Configure logging
			logCfg := zap.NewProductionConfig()

			// Set log level from flag or env var
			if envLevel := os.Getenv(envLogLevel); envLevel != "" {
				opts.LogConfig.Level = envLevel
			}
			if level := opts.LogConfig.Level; level != "" {
				logLevel, err := zapcore.ParseLevel(level)
				if err != nil {
					return fmt.Errorf("invalid log level %q: %w", level, err)
				}
				logCfg.Level = zap.NewAtomicLevelAt(logLevel)
			}

			// Set log format from flag or env var
			if envFormat := os.Getenv(envLogFormat); envFormat != "" {
				opts.LogConfig.Format = envFormat
			}
			if opts.LogConfig.Format == "console" {
				logCfg.Encoding = "console"
				logCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
			}

			// Create logger
			logger, err := logCfg.Build()
			if err != nil {
				return fmt.Errorf("failed to create logger: %w", err)
			}

			// Update settings with new logger
			if opts.Settings != nil {
				opts.Settings.Logger = logger
			}

			// Find config file
			configPath := findConfigFile(opts.ConfigPath)
			if configPath == "" {
				return fmt.Errorf("config file not found in any of the default locations")
			}
			opts.ConfigPath = configPath

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("please specify a subcommand, use --help for more information")
		},
		SilenceUsage: true,
	}

	// Add persistent flags
	cmd.PersistentFlags().StringVar(&opts.ConfigPath, "config", "", "path to configuration file (env: SREDIAG_CONFIG)")
	cmd.PersistentFlags().StringVar(&opts.LogConfig.Level, "log-level", "info", "log level (debug, info, warn, error)")
	cmd.PersistentFlags().StringVar(&opts.LogConfig.Format, "log-format", "console", "log format (json, console)")

	// Add output format flags
	cmd.PersistentFlags().StringVar(&outputOpts.Format, "output", "table", "output format (json, yaml, table)")
	cmd.PersistentFlags().BoolVar(&outputOpts.Quiet, "quiet", false, "only output essential information")
	cmd.PersistentFlags().BoolVar(&outputOpts.NoColor, "no-color", false, "disable color output")
	cmd.PersistentFlags().StringVar(&outputOpts.OutputFile, "output-file", "", "write output to file")

	// Add subcommands
	cmd.AddCommand(
		newStartCmd(opts),
		newVersionCmd(),
		newDiagnoseCmd(opts),
		newBuildCmd(opts),
		NewPluginCmd(opts),
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

// Execute creates the root command with the given settings and executes it
func Execute(s *settings.CommandSettings) error {
	opts := &Options{
		Settings: s,
	}

	rootCmd := NewRootCommand(opts)
	return rootCmd.Execute()
}
