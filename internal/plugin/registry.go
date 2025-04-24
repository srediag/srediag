package plugin

import (
	"fmt"
	"sync"

	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/factory"
	"github.com/srediag/srediag/internal/types"
)

// Registry manages plugin registration and factories
type Registry struct {
	mu        sync.RWMutex
	logger    *zap.Logger
	plugins   map[string]types.IPlugin
	factories map[string]*factory.Factory
}

// NewRegistry creates a new plugin registry
func NewRegistry(logger *zap.Logger) *Registry {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &Registry{
		logger:    logger,
		plugins:   make(map[string]types.IPlugin),
		factories: make(map[string]*factory.Factory),
	}
}

// RegisterPlugin registers a plugin with the registry
func (r *Registry) RegisterPlugin(plugin types.IPlugin) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := plugin.Validate(); err != nil {
		return fmt.Errorf("invalid plugin: %w", err)
	}

	id := plugin.GetName()
	if _, exists := r.plugins[id]; exists {
		return fmt.Errorf("plugin %s already registered", id)
	}

	r.plugins[id] = plugin
	r.logger.Info("Registered plugin",
		zap.String("id", id),
		zap.String("category", string(plugin.GetCategory())))
	return nil
}

// UnregisterPlugin removes a plugin from the registry
func (r *Registry) UnregisterPlugin(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.plugins[id]; !exists {
		return fmt.Errorf("plugin %s not found", id)
	}

	delete(r.plugins, id)
	r.logger.Info("Unregistered plugin", zap.String("id", id))
	return nil
}

// GetPlugin returns a plugin by ID
func (r *Registry) GetPlugin(id string) (types.IPlugin, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	plugin, exists := r.plugins[id]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", id)
	}

	return plugin, nil
}

// ListPlugins returns all registered plugins
func (r *Registry) ListPlugins() []types.IPlugin {
	r.mu.RLock()
	defer r.mu.RUnlock()

	plugins := make([]types.IPlugin, 0, len(r.plugins))
	for _, plugin := range r.plugins {
		plugins = append(plugins, plugin)
	}
	return plugins
}

// RegisterFactory registers a factory with the registry
func (r *Registry) RegisterFactory(factory *factory.Factory) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if factory == nil {
		return fmt.Errorf("factory cannot be nil")
	}

	id := factory.GetID()
	if _, exists := r.factories[id]; exists {
		return fmt.Errorf("factory %s already registered", id)
	}

	r.factories[id] = factory
	r.logger.Info("Registered factory", zap.String("id", id))
	return nil
}

// UnregisterFactory removes a factory from the registry
func (r *Registry) UnregisterFactory(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.factories[id]; !exists {
		return fmt.Errorf("factory %s not found", id)
	}

	delete(r.factories, id)
	r.logger.Info("Unregistered factory", zap.String("id", id))
	return nil
}

// GetFactory returns a factory by ID
func (r *Registry) GetFactory(id string) (*factory.Factory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	factory, exists := r.factories[id]
	if !exists {
		return nil, fmt.Errorf("factory %s not found", id)
	}

	return factory, nil
}

// GetFactories returns all registered factories
func (r *Registry) GetFactories() map[string]*factory.Factory {
	r.mu.RLock()
	defer r.mu.RUnlock()

	factories := make(map[string]*factory.Factory, len(r.factories))
	for id, factory := range r.factories {
		factories[id] = factory
	}
	return factories
}

// Clear removes all plugins and factories from the registry
func (r *Registry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.plugins = make(map[string]types.IPlugin)
	r.factories = make(map[string]*factory.Factory)
	r.logger.Info("Cleared registry")
}
