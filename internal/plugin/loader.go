package plugin

import (
	"fmt"
	"plugin"

	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/types"
)

// Loader handles plugin loading and registration
type Loader struct {
	logger   *zap.Logger
	registry types.Registry
}

// NewLoader creates a new plugin loader
func NewLoader(logger *zap.Logger, registry types.Registry) *Loader {
	return &Loader{
		logger:   logger,
		registry: registry,
	}
}

// LoadCustomPlugins loads plugins from the specified directory
func (l *Loader) LoadCustomPlugins(pluginDir string) error {
	l.logger.Info("Loading custom plugins", zap.String("dir", pluginDir))
	return nil
}

// GetFactories returns all registered component factories
func (l *Loader) GetFactories() (map[component.Type]component.Factory,
	map[component.Type]component.Factory,
	map[component.Type]component.Factory,
	map[component.Type]component.Factory,
	map[component.Type]component.Factory) {
	return nil, nil, nil, nil, nil
}

// LoadPlugin loads a single plugin
func (l *Loader) LoadPlugin(path string) error {
	plug, err := plugin.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open plugin %s: %w", path, err)
	}

	// Look up a symbol (function or variable)
	symFactory, err := plug.Lookup("Factory")
	if err != nil {
		return fmt.Errorf("factory symbol not found in plugin %s: %w", path, err)
	}

	// Assert that loaded symbol is a Factory
	factory, ok := symFactory.(types.Factory)
	if !ok {
		return fmt.Errorf("symbol in %s does not implement Factory interface", path)
	}

	// Register the factory
	return l.registry.RegisterFactory(factory)
}
