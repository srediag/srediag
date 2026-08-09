package plugin

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/srediag/srediag/internal/core"
)

// Package plugin provides plugin management and CLI entrypoints for SREDIAG plugins.
//
// This file defines CLI entrypoints for plugin commands, wiring Cobra commands to internal plugin logic.
//
// Usage:
//   - Use these CLI functions as entrypoints for 'srediag plugin' subcommands.
//   - Each function extracts parameters from the CLI context, instantiates the appropriate manager, and delegates to the correct method.
//
// Best Practices:
//   - Always validate required flags and parameters before calling plugin manager methods.
//   - Log all errors and important events for traceability.
//   - Use context-aware logging and error handling for better diagnostics.
//

// CLI_List is the entrypoint for 'srediag plugin list'.
//
// Parameters:
//   - ctx: Application context containing logger and configuration.
//   - cmd: Cobra command instance.
//   - args: Command-line arguments.
//
// Returns:
//   - error: If plugin listing fails or is not implemented, returns a detailed error.
func CLI_List(ctx *core.AppContext, cmd *cobra.Command, args []string) error {
	logger := ctx.Logger
	if logger == nil {
		var err error
		logger, err = core.NewLogger(nil)
		if err != nil {
			return fmt.Errorf("failed to create fallback logger: %w", err)
		}
	}
	// Create plugin registry
	registry := NewPluginRegistry(logger)

	// Register built-in plugins
	if err := registry.RegisterBuiltinPlugins(); err != nil {
		return fmt.Errorf("failed to register built-in plugins: %w", err)
	}

	// List available plugins
	available := registry.ListAvailablePlugins()
	loaded := registry.ListLoadedPlugins()

	logger.Info("Available plugins:",
		core.ZapInt("count", len(available)))

	for _, plugin := range available {
		logger.Info("Plugin",
			core.ZapString("name", plugin.Name),
			core.ZapString("type", string(plugin.Type)),
			core.ZapString("version", plugin.Version),
			core.ZapString("description", plugin.Description))
	}

	logger.Info("Loaded plugins:",
		core.ZapInt("count", len(loaded)))

	for _, plugin := range loaded {
		logger.Info("Loaded plugin",
			core.ZapString("name", plugin.Metadata.Name),
			core.ZapString("status", string(plugin.Status)),
			core.ZapString("path", plugin.Path))
	}

	return nil
}

// CLI_Info is the entrypoint for 'srediag plugin info [name]'.
//
// Parameters:
//   - ctx: Application context containing logger and configuration.
//   - cmd: Cobra command instance.
//   - args: Command-line arguments.
//
// Returns:
//   - error: If plugin info retrieval fails or is not implemented, returns a detailed error.
func CLI_Info(ctx *core.AppContext, cmd *cobra.Command, args []string) error {
	logger := ctx.Logger
	if logger == nil {
		var err error
		logger, err = core.NewLogger(nil)
		if err != nil {
			return fmt.Errorf("failed to create fallback logger: %w", err)
		}
	}
	if len(args) == 0 {
		return fmt.Errorf("plugin name is required")
	}

	pluginName := args[0]

	// Create plugin registry
	registry := NewPluginRegistry(logger)

	// Register built-in plugins
	if err := registry.RegisterBuiltinPlugins(); err != nil {
		return fmt.Errorf("failed to register built-in plugins: %w", err)
	}

	// Get plugin information
	available := registry.ListAvailablePlugins()
	for _, plugin := range available {
		if plugin.Name == pluginName {
			logger.Info("Plugin information",
				core.ZapString("name", plugin.Name),
				core.ZapString("type", string(plugin.Type)),
				core.ZapString("version", plugin.Version),
				core.ZapString("description", plugin.Description),
				core.ZapReflect("capabilities", plugin.Capabilities),
				core.ZapString("sha256", plugin.SHA256),
				core.ZapString("signature", plugin.Signature))
			return nil
		}
	}

	return fmt.Errorf("plugin %s not found", pluginName)
}

// CLI_Enable is the entrypoint for 'srediag plugin enable [type] [name]'.
//
// Parameters:
//   - ctx: Application context containing logger and configuration.
//   - cmd: Cobra command instance.
//   - args: Command-line arguments.
//
// Returns:
//   - error: If plugin enable fails or is not implemented, returns a detailed error.
func CLI_Enable(ctx *core.AppContext, cmd *cobra.Command, args []string) error {
	logger := ctx.Logger
	if logger == nil {
		var err error
		logger, err = core.NewLogger(nil)
		if err != nil {
			return fmt.Errorf("failed to create fallback logger: %w", err)
		}
	}
	if len(args) < 2 {
		return fmt.Errorf("plugin type and name are required")
	}

	pluginType := args[0]
	pluginName := args[1]

	// Create plugin registry
	registry := NewPluginRegistry(logger)

	// Register built-in plugins
	if err := registry.RegisterBuiltinPlugins(); err != nil {
		return fmt.Errorf("failed to register built-in plugins: %w", err)
	}

	// Load plugin
	pluginPath := fmt.Sprintf("/tmp/srediag/plugins/%ss/%s", pluginType, pluginName)
	instance, err := registry.LoadPlugin(pluginName, pluginPath)
	if err != nil {
		return fmt.Errorf("failed to load plugin: %w", err)
	}

	logger.Info("Plugin enabled",
		core.ZapString("name", instance.Metadata.Name),
		core.ZapString("type", string(instance.Metadata.Type)),
		core.ZapString("status", string(instance.Status)))

	return nil
}

// CLI_Disable is the entrypoint for 'srediag plugin disable [name]'.
//
// Parameters:
//   - ctx: Application context containing logger and configuration.
//   - cmd: Cobra command instance.
//   - args: Command-line arguments.
//
// Returns:
//   - error: If plugin disable fails or is not implemented, returns a detailed error.
func CLI_Disable(ctx *core.AppContext, cmd *cobra.Command, args []string) error {
	logger := ctx.Logger
	if logger == nil {
		var err error
		logger, err = core.NewLogger(nil)
		if err != nil {
			return fmt.Errorf("failed to create fallback logger: %w", err)
		}
	}
	if len(args) == 0 {
		return fmt.Errorf("plugin name is required")
	}

	pluginName := args[0]

	// Create plugin registry
	registry := NewPluginRegistry(logger)

	// Register built-in plugins
	if err := registry.RegisterBuiltinPlugins(); err != nil {
		return fmt.Errorf("failed to register built-in plugins: %w", err)
	}

	// Unload plugin
	if err := registry.UnloadPlugin(pluginName); err != nil {
		return fmt.Errorf("failed to unload plugin: %w", err)
	}

	logger.Info("Plugin disabled",
		core.ZapString("name", pluginName))

	return nil
}
