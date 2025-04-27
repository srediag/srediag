package main

import (
	"fmt"
	"os"

	"github.com/srediag/srediag/cmd/srediag/commands"
	"github.com/srediag/srediag/internal/core"
)

// main is the entry point for the srediag CLI.
func main() {
	// 1. Bootstrap logger with env or default
	logLevel := os.Getenv("SREDIAG_LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	logFormat := os.Getenv("SREDIAG_LOG_FORMAT")
	if logFormat == "" {
		logFormat = "console"
	}
	bootstrapLogger, err := core.NewLogger(&core.Logger{
		Level:            logLevel,
		Format:           logFormat,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	// 2. Load config (using bootstrap logger)
	config, err := core.LoadConfig(bootstrapLogger.UnderlyingZap())
	if err != nil {
		bootstrapLogger.Error("Failed to load config", core.ZapError(err))
		os.Exit(1)
	}

	// 3. Re-initialize logger with config values (if different)
	finalLogger, err := core.NewLogger(&core.Logger{
		Level:            config.LogLevel,
		Format:           config.LogFormat,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	})
	if err != nil {
		bootstrapLogger.Error("Failed to re-initialize logger", core.ZapError(err))
		os.Exit(1)
	}

	// 4. Initialize component manager
	componentManager := core.NewComponentManager(finalLogger)

	// 5. Build AppContext
	ctx := &core.AppContext{
		Logger:           finalLogger,
		ComponentManager: componentManager,
		Config:           config,
		BuildInfo:        core.DefaultBuildInfo,
	}

	// 6. Execute root command
	if err := commands.Execute(ctx); err != nil {
		finalLogger.Error("Failed to execute command", core.ZapError(err))
		os.Exit(1)
	}
}
