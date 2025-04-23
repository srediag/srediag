// Package core provides core interfaces and components for SREDIAG
package core

import (
	"context"
	"fmt"
	"path/filepath"
	"plugin"
	"sync"

	"go.uber.org/zap"
)

// defaultPluginManager provides a default implementation of PluginManager
type defaultPluginManager struct {
	logger  *zap.Logger
	plugins map[string]Plugin
	mu      sync.RWMutex
	healthy bool
}

// NewPluginManager creates a new plugin manager instance
func NewPluginManager(logger *zap.Logger) PluginManager {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &defaultPluginManager{
		logger:  logger,
		plugins: make(map[string]Plugin),
		healthy: true,
	}
}

// Start implements Component
func (m *defaultPluginManager) Start(ctx context.Context) error {
	m.logger.Info("starting plugin manager")

	m.mu.Lock()
	defer m.mu.Unlock()

	// Start all loaded plugins
	for name, p := range m.plugins {
		if err := p.Start(ctx); err != nil {
			m.logger.Error("failed to start plugin",
				zap.String("name", name),
				zap.Error(err))
			m.healthy = false
			return fmt.Errorf("failed to start plugin %s: %w", name, err)
		}
	}

	return nil
}

// Stop implements Component
func (m *defaultPluginManager) Stop(ctx context.Context) error {
	m.logger.Info("stopping plugin manager")

	m.mu.Lock()
	defer m.mu.Unlock()

	// Stop all loaded plugins in reverse order
	for name, p := range m.plugins {
		if err := p.Stop(ctx); err != nil {
			m.logger.Error("failed to stop plugin",
				zap.String("name", name),
				zap.Error(err))
		}
	}

	return nil
}

// IsHealthy implements Component
func (m *defaultPluginManager) IsHealthy() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.healthy
}

// LoadPlugin implements PluginManager
func (m *defaultPluginManager) LoadPlugin(path string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Resolve absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to resolve plugin path: %w", err)
	}

	// Open plugin
	plug, err := plugin.Open(absPath)
	if err != nil {
		return fmt.Errorf("failed to open plugin: %w", err)
	}

	// Look up plugin symbol
	sym, err := plug.Lookup("NewPlugin")
	if err != nil {
		return fmt.Errorf("plugin does not export 'NewPlugin' symbol: %w", err)
	}

	// Type assert plugin constructor
	newPlugin, ok := sym.(func() Plugin)
	if !ok {
		return fmt.Errorf("invalid plugin constructor type")
	}

	// Create plugin instance
	p := newPlugin()
	if p == nil {
		return fmt.Errorf("plugin constructor returned nil")
	}

	// Store plugin
	name := p.GetName()
	if _, exists := m.plugins[name]; exists {
		return fmt.Errorf("plugin %s already loaded", name)
	}

	m.plugins[name] = p
	m.logger.Info("loaded plugin",
		zap.String("name", name),
		zap.String("version", p.GetVersion()),
		zap.String("type", p.GetType()))

	return nil
}

// UnloadPlugin implements PluginManager
func (m *defaultPluginManager) UnloadPlugin(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	p, exists := m.plugins[name]
	if !exists {
		return fmt.Errorf("plugin %s not found", name)
	}

	// Stop plugin before unloading
	if err := p.Stop(context.Background()); err != nil {
		m.logger.Error("failed to stop plugin during unload",
			zap.String("name", name),
			zap.Error(err))
	}

	delete(m.plugins, name)
	m.logger.Info("unloaded plugin", zap.String("name", name))

	return nil
}

// GetPlugin implements PluginManager
func (m *defaultPluginManager) GetPlugin(name string) (Plugin, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	p, exists := m.plugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", name)
	}

	return p, nil
}

// ListPlugins implements PluginManager
func (m *defaultPluginManager) ListPlugins() []Plugin {
	m.mu.RLock()
	defer m.mu.RUnlock()

	plugins := make([]Plugin, 0, len(m.plugins))
	for _, p := range m.plugins {
		plugins = append(plugins, p)
	}

	return plugins
}
