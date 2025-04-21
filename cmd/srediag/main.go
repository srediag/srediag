package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/logger"
	"github.com/srediag/srediag/internal/telemetry"
)

// Version is set at build-time via -ldflags
var Version = "v0.1.0"

func main() {
	// Initialize root command
	debug := viper.GetBool("debug")
	if err := logger.Init(debug); err != nil {
		fmt.Fprintf(os.Stderr, "logger init error: %v\n", err)
		os.Exit(1)
	}
	logger, _ := zap.NewProduction()
	defer func() { _ = logger.Sync() }()

	tp, err := telemetry.InitTracer("srediag")
	if err != nil {
		logger.Sugar().Fatalf("telemetry init: %v", err)
	}

	ctx := context.Background()
	defer func() { _ = tp.Shutdown(ctx) }()
	rootCmd := &cobra.Command{
		Use:     "srediag",
		Short:   "SRE Diagnostics agent for OBSERVO",
		Long:    "srediag is a modular observability agent, extensible via plugins, collecting metrics, logs, and traces.",
		Version: Version,
	}

	// Global flags
	rootCmd.PersistentFlags().StringP("config", "c", "configs/config.yaml", "Path to config file")
	if err := viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config")); err != nil {
		logger.Warn("failed to bind flag", zap.Error(err))
	}

	// Subcommands
	rootCmd.AddCommand(startCmd(logger))
	rootCmd.AddCommand(versionCmd())

	if err := rootCmd.Execute(); err != nil {
		logger.Fatal("command execution failed", zap.Error(err))
		os.Exit(1)
	}
}

func startCmd(logger *zap.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start the Srediag agent",
		Run: func(cmd *cobra.Command, args []string) {
			configPath := viper.GetString("config")
			logger.Info("Loading config", zap.String("path", configPath))

			// TODO: load YAML into struct via viper
			// Initialize telemetry, plugin loader, main loop

			fmt.Println("Starting srediag agent...")
		},
	}
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of srediag",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("srediag version %s\n", Version)
		},
	}
}
