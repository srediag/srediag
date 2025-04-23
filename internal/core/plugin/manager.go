// Package plugin provides plugin types and utilities for SREDIAG
package plugin

import (
	"context"
	"fmt"
	"path/filepath"
	"plugin"
	"sync"

	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/core"
)

// Manager manages plugin lifecycle and configuration
type Manager interface {
	core.Component
	// LoadPlugin loads a plugin from the given path
	LoadPlugin(path string) error
	// UnloadPlugin unloads a plugin by name
	UnloadPlugin(name string) error
	// GetPlugin returns a plugin by name
	GetPlugin(name string) (core.Plugin, error)
	// ListPlugins returns a list of loaded plugins
	ListPlugins() []core.Plugin
	// GetPluginInfo returns information about a plugin
	GetPluginInfo(name string) (*Info, error)
	// GetPluginMetadata returns metadata about a plugin
	GetPluginMetadata(name string) (*Metadata, error)
	// GetPluginConfig returns the configuration of a plugin
	GetPluginConfig(name string) (*Config, error)
	// ConfigurePlugin configures a plugin with the given configuration
	ConfigurePlugin(name string, config *Config) error
}

// defaultManager provides a default implementation of Manager
type defaultManager struct {
	logger  *zap.Logger
	plugins map[string]core.Plugin
	mu      sync.RWMutex
	healthy bool
}

// NewManager creates a new plugin manager instance
func NewManager(logger *zap.Logger) Manager {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &defaultManager{
		logger:  logger.Named("plugin-manager"),
		plugins: make(map[string]core.Plugin),
		healthy: true,
	}
}

// Start implements core.Component
func (m *defaultManager) Start(ctx context.Context) error {
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

// Stop implements core.Component
func (m *defaultManager) Stop(ctx context.Context) error {
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

// IsHealthy implements core.Component
func (m *defaultManager) IsHealthy() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.healthy
}

// LoadPlugin implements Manager
func (m *defaultManager) LoadPlugin(path string) error {
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
	newPlugin, ok := sym.(func() core.Plugin)
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

// UnloadPlugin implements Manager
func (m *defaultManager) UnloadPlugin(name string) error {
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

// GetPlugin implements Manager
func (m *defaultManager) GetPlugin(name string) (core.Plugin, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	p, exists := m.plugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", name)
	}

	return p, nil
}

// ListPlugins implements Manager
func (m *defaultManager) ListPlugins() []core.Plugin {
	m.mu.RLock()
	defer m.mu.RUnlock()

	plugins := make([]core.Plugin, 0, len(m.plugins))
	for _, p := range m.plugins {
		plugins = append(plugins, p)
	}

	return plugins
}

// GetPluginInfo implements Manager
func (m *defaultManager) GetPluginInfo(name string) (*Info, error) {
	p, err := m.GetPlugin(name)
	if err != nil {
		return nil, err
	}

	if bp, ok := p.(*Base); ok {
		info := bp.GetInfo()
		return &info, nil
	}

	return nil, fmt.Errorf("plugin %s does not support info retrieval", name)
}

// GetPluginMetadata implements Manager
func (m *defaultManager) GetPluginMetadata(name string) (*Metadata, error) {
	p, err := m.GetPlugin(name)
	if err != nil {
		return nil, err
	}

	if bp, ok := p.(*Base); ok {
		metadata := bp.GetMetadata()
		return &metadata, nil
	}

	return nil, fmt.Errorf("plugin %s does not support metadata retrieval", name)
}

// GetPluginConfig implements Manager
func (m *defaultManager) GetPluginConfig(name string) (*Config, error) {
	p, err := m.GetPlugin(name)
	if err != nil {
		return nil, err
	}

	if bp, ok := p.(*Base); ok {
		config := bp.GetConfig()
		return &config, nil
	}

	return nil, fmt.Errorf("plugin %s does not support config retrieval", name)
}

// ConfigurePlugin implements Manager
func (m *defaultManager) ConfigurePlugin(name string, config *Config) error {
	p, err := m.GetPlugin(name)
	if err != nil {
		return err
	}

	return p.Configure(config)
}
