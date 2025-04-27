package plugin

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"go.opentelemetry.io/collector/component"

	"github.com/srediag/srediag/internal/core"
)

// Loader handles plugin discovery and loading.
type Loader struct {
	logger  *core.Logger
	manager *PluginManager
}

// NewLoader creates a new plugin loader.
func NewLoader(logger *core.Logger, manager *PluginManager) *Loader {
	return &Loader{
		logger:  logger,
		manager: manager,
	}
}

// LoadPlugins loads all plugins from the specified directory.
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
