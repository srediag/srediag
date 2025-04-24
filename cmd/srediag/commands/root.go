package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/srediag/srediag/cmd/srediag/internal/config"
	"github.com/srediag/srediag/cmd/srediag/internal/version"
	"github.com/srediag/srediag/internal/types"
)

// Exit codes
const (
	ExitSuccess       = 0
	ExitGeneralError  = 1
	ExitConfigError   = 2
	ExitPermissionErr = 3
	ExitNotFoundError = 4
	ExitTimeoutError  = 5
)

var (
	// Global flags
	configPath   string
	outputFormat string
	verbose      bool
	quiet        bool
	outputFile   string
	showVersion  bool

	// Configuration manager
	configManager *config.Manager

	// Settings
	settings Settings

	// Root command
	rootCmd = &cobra.Command{
		Use:   "srediag",
		Short: "SREDIAG - SRE Diagnostics Tool",
		Long: `SREDIAG is a tool for SRE diagnostics and monitoring.
It helps identify and diagnose issues in your infrastructure.`,
		RunE: run,
	}
)

// Execute executes the root command
func Execute(s Settings) error {
	settings = s
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&configPath, "config", config.DefaultConfigPath, "path to configuration file")
	rootCmd.PersistentFlags().StringVar(&outputFormat, "format", "table", "output format (json/yaml/table)")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "enable verbose logging")
	rootCmd.PersistentFlags().BoolVar(&quiet, "quiet", false, "suppress non-error output")
	rootCmd.PersistentFlags().StringVar(&outputFile, "output", "", "output file path")
	rootCmd.PersistentFlags().BoolVar(&showVersion, "version", false, "print version information")

	// Add all subcommands
	rootCmd.AddCommand(
		newDiagnoseCmd(),
		newAnalyzeCmd(),
		newMonitorCmd(),
		newSecurityCmd(),
	)
}

func initConfig() {
	// Initialize Viper configuration
	config.InitializeConfig(configPath)

	// Create configuration manager with appropriate logger
	if settings.Logger != nil {
		configManager = config.NewManager(settings.Logger)
	} else {
		// Create a default logger if none is provided in settings
		logger, _ := zap.NewProduction()
		configManager = config.NewManager(logger)
	}
}

func run(cmd *cobra.Command, args []string) error {
	if showVersion {
		fmt.Printf("SREDIAG %s (commit: %s, built: %s)\n",
			version.Version,
			version.GitCommit,
			version.BuildDate)
		return nil
	}

	// Initialize configuration manager
	if err := configManager.Initialize(cmd.Context()); err != nil {
		return fmt.Errorf("failed to initialize configuration: %w", err)
	}
	defer func() {
		if err := configManager.Shutdown(context.Background()); err != nil {
			settings.Logger.Error("Failed to shutdown configuration manager",
				zap.Error(err))
		}
	}()

	// Get configuration
	cfg := configManager.Get()

	// Initialize logger with configuration
	logger, err := initLogger(cfg)
	if err != nil {
		return err
	}

	logger.Info("starting SREDIAG",
		zap.String("version", version.Version),
		zap.String("config", configPath),
	)

	return nil
}

func initLogger(cfg *types.Config) (*zap.Logger, error) {
	var config zap.Config

	if cfg.Core.LogFormat == "console" {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		config = zap.NewProductionConfig()
	}

	level := string(cfg.Core.LogLevel)
	if verbose {
		level = string(types.ConfigLogLevelDebug)
	} else if quiet {
		level = string(types.ConfigLogLevelError)
	}

	parsedLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level %q: %w", level, err)
	}
	config.Level = zap.NewAtomicLevelAt(parsedLevel)

	logger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	return logger.Named("srediag"), nil
}
