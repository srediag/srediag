package plugin

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"go.opentelemetry.io/collector/component"

	"github.com/srediag/srediag/internal/core"
)

// Package plugin provides plugin management, discovery, and loading logic for SREDIAG plugins.
//
// This file defines the Loader type for plugin discovery and loading, as well as methods for retrieving component factories.
//
// Usage:
//   - Use Loader to discover and load plugins from a specified directory.
//   - Use GetFactories to retrieve loaded component factories grouped by type.
//
// Best Practices:
//   - Always check for errors from LoadPlugins.
//   - Use logger for all error and status reporting.
//   - Ensure plugin directories exist before loading.
//
// TODO:
//   - Add context.Context support for cancellation and timeouts in all methods.
//   - Improve error reporting and diagnostics for plugin loading failures.

// Loader handles plugin discovery and loading.
//
// Usage:
//   - Instantiate with NewLoader, providing a logger and plugin manager.
//   - Call LoadPlugins to discover and load plugins from a directory.
//   - Call GetFactories to retrieve loaded component factories grouped by type.
type Loader struct {
	logger  *core.Logger
	manager *PluginManager
}

// NewLoader creates a new plugin loader.
//
// Parameters:
//   - logger: Logger for status and error reporting.
//   - manager: PluginManager instance for plugin orchestration.
//
// Returns:
//   - *Loader: A new Loader instance.
func NewLoader(logger *core.Logger, manager *PluginManager) *Loader {
	return &Loader{
		logger:  logger,
		manager: manager,
	}
}

// LoadPlugins loads all plugins from the specified directory.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts.
//   - pluginDir: Directory containing plugin binaries.
//
// Returns:
//   - error: If loading any plugin fails, returns a detailed error.
//
// Side Effects:
//   - Modifies internal plugin manager state.
//   - Logs status and errors.
func (l *Loader) LoadPlugins(ctx context.Context, pluginDir string) error {
	l.logger.Info("Loading plugins", core.ZapString("dir", pluginDir))

	// Ensure plugin directory exists
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		return fmt.Errorf("failed to create plugin directory: %w", err)
	}

	// Walk through plugin directories
	for _, typ := range []string{"receivers", "processors", "exporters", "extensions"} {
		typeDir := filepath.Join(pluginDir, typ)
		if err := os.MkdirAll(typeDir, 0755); err != nil {
			return fmt.Errorf("failed to create plugin type directory: %w", err)
		}

		entries, err := os.ReadDir(typeDir)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return fmt.Errorf("failed to read plugin directory: %w", err)
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			pluginName := entry.Name()
			l.logger.Info("Loading plugin",
				core.ZapString("type", typ),
				core.ZapString("name", pluginName))

			if err := l.manager.Load(ctx, core.ComponentType(typ), pluginName); err != nil {
				l.logger.Error("Failed to load plugin",
					core.ZapString("type", typ),
					core.ZapString("name", pluginName),
					core.ZapError(err))
				continue
			}
		}
	}

	return nil
}

// GetFactories returns all loaded component factories grouped by type.
//
// Returns:
//   - receivers: Map of receiver component types to factories.
//   - processors: Map of processor component types to factories.
//   - exporters: Map of exporter component types to factories.
//   - extensions: Map of extension component types to factories.
func (l *Loader) GetFactories() (
	receivers map[component.Type]component.Factory,
	processors map[component.Type]component.Factory,
	exporters map[component.Type]component.Factory,
	extensions map[component.Type]component.Factory,
) {
	receivers = make(map[component.Type]component.Factory)
	processors = make(map[component.Type]component.Factory)
	exporters = make(map[component.Type]component.Factory)
	extensions = make(map[component.Type]component.Factory)

	// Get all loaded plugins
	plugins := l.manager.List()
	for _, meta := range plugins {
		plugin, ok := l.manager.Get(meta.Name)
		if !ok {
			continue
		}

		factory, err := plugin.Factory()
		if err != nil {
			l.logger.Error("Failed to get factory",
				core.ZapString("plugin", meta.Name),
				core.ZapError(err))
			continue
		}

		// Create component type
		compType, err := component.NewType(string(meta.Type) + "/" + meta.Name)
		if err != nil {
			l.logger.Error("Failed to create component type",
				core.ZapString("plugin", meta.Name),
				core.ZapError(err))
			continue
		}

		switch meta.Type {
		case core.TypeReceiver:
			receivers[compType] = factory
		case core.TypeProcessor:
			processors[compType] = factory
		case core.TypeExporter:
			exporters[compType] = factory
		case core.TypeExtension:
			extensions[compType] = factory
		}
	}

	return receivers, processors, exporters, extensions
}
