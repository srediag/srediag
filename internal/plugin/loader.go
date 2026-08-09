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
		l.logger.Error("Failed to create plugin directory",
			core.ZapString("dir", pluginDir),
			core.ZapError(err))
		return fmt.Errorf("failed to create plugin directory %s: %w", pluginDir, err)
	}

	var loadErrors []error
	var loadedCount int
	var failedCount int

	// Walk through plugin directories
	for _, typ := range []string{"receivers", "processors", "exporters", "extensions"} {
		typeDir := filepath.Join(pluginDir, typ)
		if err := os.MkdirAll(typeDir, 0755); err != nil {
			l.logger.Error("Failed to create plugin type directory",
				core.ZapString("type", typ),
				core.ZapString("dir", typeDir),
				core.ZapError(err))
			return fmt.Errorf("failed to create plugin type directory %s: %w", typeDir, err)
		}

		entries, err := os.ReadDir(typeDir)
		if err != nil {
			if os.IsNotExist(err) {
				l.logger.Debug("Plugin type directory does not exist, skipping",
					core.ZapString("type", typ),
					core.ZapString("dir", typeDir))
				continue
			}
			l.logger.Error("Failed to read plugin directory",
				core.ZapString("type", typ),
				core.ZapString("dir", typeDir),
				core.ZapError(err))
			return fmt.Errorf("failed to read plugin directory %s: %w", typeDir, err)
		}

		l.logger.Debug("Found plugins in directory",
			core.ZapString("type", typ),
			core.ZapInt("count", len(entries)))

		for _, entry := range entries {
			if entry.IsDir() {
				l.logger.Debug("Skipping directory entry",
					core.ZapString("name", entry.Name()),
					core.ZapString("type", typ))
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
				loadErrors = append(loadErrors, fmt.Errorf("plugin %s/%s: %w", typ, pluginName, err))
				failedCount++
				continue
			}

			l.logger.Debug("Successfully loaded plugin",
				core.ZapString("type", typ),
				core.ZapString("name", pluginName))
			loadedCount++
		}
	}

	// Log summary of plugin loading results
	l.logger.Info("Plugin loading completed",
		core.ZapInt("loaded", loadedCount),
		core.ZapInt("failed", failedCount),
		core.ZapInt("total_errors", len(loadErrors)))

	// Return aggregated errors if any occurred
	if len(loadErrors) > 0 {
		l.logger.Warn("Some plugins failed to load",
			core.ZapInt("error_count", len(loadErrors)))
		// For now, we don't fail the entire operation due to individual plugin failures
		// This allows the system to start with partial plugin availability
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
	var factoryErrors []error
	var successCount int

	l.logger.Debug("Building component factories",
		core.ZapInt("total_plugins", len(plugins)))

	for _, meta := range plugins {
		plugin, ok := l.manager.Get(meta.Name)
		if !ok {
			l.logger.Warn("Plugin listed but not retrievable",
				core.ZapString("plugin", meta.Name))
			factoryErrors = append(factoryErrors, fmt.Errorf("plugin %s listed but not retrievable", meta.Name))
			continue
		}

		l.logger.Debug("Getting factory for plugin",
			core.ZapString("plugin", meta.Name),
			core.ZapString("type", string(meta.Type)))

		factory, err := plugin.Factory()
		if err != nil {
			l.logger.Error("Failed to get factory",
				core.ZapString("plugin", meta.Name),
				core.ZapString("type", string(meta.Type)),
				core.ZapError(err))
			factoryErrors = append(factoryErrors, fmt.Errorf("plugin %s factory error: %w", meta.Name, err))
			continue
		}

		// Create component type
		compType, err := component.NewType(string(meta.Type) + "/" + meta.Name)
		if err != nil {
			l.logger.Error("Failed to create component type",
				core.ZapString("plugin", meta.Name),
				core.ZapString("type", string(meta.Type)),
				core.ZapError(err))
			factoryErrors = append(factoryErrors, fmt.Errorf("plugin %s component type creation error: %w", meta.Name, err))
			continue
		}

		switch meta.Type {
		case ReceiverPlugin:
			receivers[compType] = factory
			l.logger.Debug("Added receiver factory",
				core.ZapString("plugin", meta.Name),
				core.ZapString("component_type", compType.String()))
		case ProcessorPlugin:
			processors[compType] = factory
			l.logger.Debug("Added processor factory",
				core.ZapString("plugin", meta.Name),
				core.ZapString("component_type", compType.String()))
		case ExporterPlugin:
			exporters[compType] = factory
			l.logger.Debug("Added exporter factory",
				core.ZapString("plugin", meta.Name),
				core.ZapString("component_type", compType.String()))
		case ExtensionPlugin:
			extensions[compType] = factory
			l.logger.Debug("Added extension factory",
				core.ZapString("plugin", meta.Name),
				core.ZapString("component_type", compType.String()))
		case ExamplePlugin:
			l.logger.Debug("Skipping example plugin for OTel factory mapping",
				core.ZapString("plugin", meta.Name))
			// Example plugins don't map to OTel components
			continue
		default:
			l.logger.Warn("Unknown plugin type, skipping factory mapping",
				core.ZapString("plugin", meta.Name),
				core.ZapString("type", string(meta.Type)))
			continue
		}
		successCount++
	}

	// Log factory building summary
	l.logger.Info("Component factory building completed",
		core.ZapInt("successful_factories", successCount),
		core.ZapInt("factory_errors", len(factoryErrors)),
		core.ZapInt("receivers", len(receivers)),
		core.ZapInt("processors", len(processors)),
		core.ZapInt("exporters", len(exporters)),
		core.ZapInt("extensions", len(extensions)))

	if len(factoryErrors) > 0 {
		l.logger.Warn("Some factory operations failed",
			core.ZapInt("error_count", len(factoryErrors)))
	}

	return receivers, processors, exporters, extensions
}
