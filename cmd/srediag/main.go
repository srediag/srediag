package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/srediag/srediag/cmd/srediag/commands"
	"github.com/srediag/srediag/internal/factory"
	"github.com/srediag/srediag/internal/plugin"
)

func main() {
	// Initialize logger with a custom sync function
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := config.Build(zap.WithCaller(true))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create logger: %v\n", err)
		os.Exit(1)
	}

	// Ensure we sync the logger on exit, ignoring specific known errors
	defer func() {
		// Ignore sync errors as they are generally harmless
		// See: https://github.com/uber-go/zap/issues/880
		_ = logger.Sync()
	}()

	// Create factory registry
	registry := factory.NewRegistry()

	// Create plugin loader
	loader := plugin.NewLoader(logger, registry)

	// Create and use OpenTelemetry component loader
	otelLoader := plugin.NewOTelComponentLoader(logger)
	if err := otelLoader.RegisterBuiltinFactories(loader); err != nil {
		logger.Error("Failed to register built-in OpenTelemetry components", zap.Error(err))
		// Continue execution even if some components fail to load
	}

	// Load custom plugins
	pluginDir := getPluginDir()
	if err := loader.LoadCustomPlugins(pluginDir); err != nil {
		// Check if this is a version mismatch error
		if strings.Contains(err.Error(), "different version") {
			logger.Warn("Some plugins were built with different versions and cannot be loaded. Please rebuild the plugins with the current version.",
				zap.String("plugin_dir", pluginDir),
				zap.Error(err))
		} else {
			logger.Error("Failed to load custom plugins",
				zap.String("plugin_dir", pluginDir),
				zap.Error(err))
		}
		// Continue execution even if some plugins fail to load
	}

	// Get all factories
	receivers, processors, exporters, extensions, connectors := loader.GetFactories()

	// Initialize command executor with factories
	if err := commands.Execute(commands.Settings{
		Receivers:  receivers,
		Processors: processors,
		Exporters:  exporters,
		Extensions: extensions,
		Connectors: connectors,
		Logger:     logger,
	}); err != nil {
		logger.Fatal("Failed to execute command", zap.Error(err))
	}
}

// getPluginDir returns the directory containing custom plugins
func getPluginDir() string {
	// First check environment variable
	if dir := os.Getenv("SREDIAG_PLUGIN_DIR"); dir != "" {
		return dir
	}

	// Then check default locations
	candidates := []string{
		"/etc/srediag/plugins",
		filepath.Join(os.Getenv("HOME"), ".srediag/plugins"),
		"./plugins",
	}

	for _, dir := range candidates {
		if _, err := os.Stat(dir); err == nil {
			return dir
		}
	}

	// Return default location if none exists
	return "./plugins"
}
