package core

import (
	"context"
	"fmt"
	"path/filepath"
	"plugin"
	"sync"

	"go.uber.org/zap"
)

// DefaultPluginManager is the default implementation of PluginManager
type DefaultPluginManager struct {
	logger  *zap.Logger
	plugins map[string]Plugin
	mu      sync.RWMutex
	healthy bool
	running bool
}

// NewPluginManager creates a new instance of DefaultPluginManager
func NewPluginManager(logger *zap.Logger) *DefaultPluginManager {
	return &DefaultPluginManager{
		logger:  logger,
		plugins: make(map[string]Plugin),
		healthy: true,
	}
}

// Start initializes the plugin manager
func (pm *DefaultPluginManager) Start(ctx context.Context) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if pm.running {
		return fmt.Errorf("plugin manager is already running")
	}

	pm.logger.Info("starting plugin manager")
	pm.running = true
	return nil
}

// Stop stops the plugin manager and all plugins
func (pm *DefaultPluginManager) Stop(ctx context.Context) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if !pm.running {
		return fmt.Errorf("plugin manager is not running")
	}

	pm.logger.Info("stopping plugin manager")

	// Stop all plugins
	for name, p := range pm.plugins {
		if err := p.Stop(ctx); err != nil {
			pm.logger.Error("failed to stop plugin",
				zap.String("name", name),
				zap.Error(err))
		}
	}

	pm.running = false
	return nil
}

// IsHealthy returns the health status
func (pm *DefaultPluginManager) IsHealthy() bool {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.healthy
}

// LoadPlugin loads a plugin from the given path
func (pm *DefaultPluginManager) LoadPlugin(path string) (Plugin, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if !pm.running {
		return nil, fmt.Errorf("plugin manager is not running")
	}

	// Check if plugin is already loaded
	name := filepath.Base(path)
	if p, exists := pm.plugins[name]; exists {
		return p, nil
	}

	// Load plugin
	plug, err := plugin.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open plugin %s: %w", path, err)
	}

	// Look up Plugin symbol
	sym, err := plug.Lookup("Plugin")
	if err != nil {
		return nil, fmt.Errorf("plugin %s does not export 'Plugin' symbol: %w", path, err)
	}

	// Assert that the symbol is a Plugin
	p, ok := sym.(Plugin)
	if !ok {
		return nil, fmt.Errorf("plugin %s does not implement Plugin interface", path)
	}

	// Start the plugin
	if err := p.Start(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to start plugin %s: %w", path, err)
	}

	pm.plugins[name] = p
	pm.logger.Info("loaded plugin",
		zap.String("name", name),
		zap.String("type", string(p.Type())),
		zap.String("version", p.Version()))

	return p, nil
}

// UnloadPlugin unloads a plugin by name
func (pm *DefaultPluginManager) UnloadPlugin(name string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if !pm.running {
		return fmt.Errorf("plugin manager is not running")
	}

	p, exists := pm.plugins[name]
	if !exists {
		return fmt.Errorf("plugin %s not found", name)
	}

	// Stop the plugin
	if err := p.Stop(context.Background()); err != nil {
		return fmt.Errorf("failed to stop plugin %s: %w", name, err)
	}

	delete(pm.plugins, name)
	pm.logger.Info("unloaded plugin", zap.String("name", name))

	return nil
}

// GetPlugin returns a plugin by name
func (pm *DefaultPluginManager) GetPlugin(name string) (Plugin, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	p, exists := pm.plugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", name)
	}

	return p, nil
}

// ListPlugins returns all loaded plugins
func (pm *DefaultPluginManager) ListPlugins() []Plugin {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	plugins := make([]Plugin, 0, len(pm.plugins))
	for _, p := range pm.plugins {
		plugins = append(plugins, p)
	}

	return plugins
}
