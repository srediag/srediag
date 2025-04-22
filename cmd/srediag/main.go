package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/app"
	"github.com/srediag/srediag/internal/config"
	"github.com/srediag/srediag/internal/logger"
	"github.com/srediag/srediag/internal/telemetry"
)

// Version is set at build-time via -ldflags
var Version = "v0.1.0"

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "srediag",
		Short: "SREDIAG - SRE Diagnostics and Analysis System",
		Long: `SREDIAG is a diagnostic and analysis tool for SRE
that integrates multiple data sources and provides insights through plugins.`,
		Version: Version,
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.srediag.yaml)")
	rootCmd.AddCommand(createStartCmd())
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".srediag")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func createStartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start the SREDIAG service",
		Long:  `Start the SREDIAG service with the specified configuration`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load configuration
			cfg, err := config.LoadConfig(cfgFile)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// Initialize logger
			if err := logger.Init(false); err != nil {
				return fmt.Errorf("failed to initialize logger: %w", err)
			}
			log, err := logger.Get()
			if err != nil {
				return fmt.Errorf("failed to get logger: %w", err)
			}
			defer func() { _ = log.Sync() }()

			// Create context with cancellation
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// Setup signal handling
			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				sig := <-sigCh
				log.Info("received shutdown signal", zap.String("signal", sig.String()))
				cancel()
			}()

			// Initialize telemetry if enabled
			if cfg.Telemetry.Enabled {
				tp, err := telemetry.InitTracer(cfg.Telemetry.ServiceName)
				if err != nil {
					log.Error("failed to initialize telemetry", zap.Error(err))
				} else {
					defer func() { _ = tp.Shutdown(context.Background()) }()
				}
			}

			// Create and start application
			application := app.New(cfg, log)
			if err := application.Start(ctx); err != nil {
				return fmt.Errorf("application error: %w", err)
			}

			// Wait for context cancellation
			<-ctx.Done()
			log.Info("shutting down service")

			return nil
		},
	}
}
