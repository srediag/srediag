package plugins

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/extension"
	"go.opentelemetry.io/collector/otelcol"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/zap"
)

// Manager handles plugin lifecycle and management
type Manager struct {
	logger    *zap.Logger
	mu        sync.RWMutex
	plugins   map[component.ID]BasePlugin
	factories map[component.Type]Factory
	host      component.Host
}

// NewManager creates a new plugin manager
func NewManager(logger *zap.Logger, host component.Host) *Manager {
	return &Manager{
		logger:    logger,
		plugins:   make(map[component.ID]BasePlugin),
		factories: make(map[component.Type]Factory),
		host:      host,
	}
}

// RegisterFactory registers a plugin factory
func (m *Manager) RegisterFactory(factory Factory) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	typ := factory.Type()
	if _, exists := m.factories[typ]; exists {
		return fmt.Errorf("factory already registered for type: %s", typ)
	}

	m.factories[typ] = factory
	m.logger.Info("registered plugin factory", zap.String("type", typ.String()))
	return nil
}

// CreatePlugin creates a new plugin instance
func (m *Manager) CreatePlugin(pluginType component.Type, config interface{}) (BasePlugin, error) {
	m.mu.RLock()
	factory, exists := m.factories[pluginType]
	m.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("no factory registered for plugin type: %s", pluginType)
	}

	plugin, err := factory.CreatePlugin(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create plugin: %w", err)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	id := component.NewID(pluginType)
	if _, exists := m.plugins[id]; exists {
		return nil, fmt.Errorf("plugin already exists: %s", id.String())
	}

	m.plugins[id] = plugin
	m.logger.Info("created plugin",
		zap.String("id", id.String()),
		zap.String("version", plugin.Version()))

	return plugin, nil
}

// GetPlugin returns a plugin by its ID
func (m *Manager) GetPlugin(id component.ID) (BasePlugin, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	plugin, exists := m.plugins[id]
	return plugin, exists
}

// ListPlugins returns all registered plugins
func (m *Manager) ListPlugins() []BasePlugin {
	m.mu.RLock()
	defer m.mu.RUnlock()

	plugins := make([]BasePlugin, 0, len(m.plugins))
	for _, p := range m.plugins {
		plugins = append(plugins, p)
	}
	return plugins
}

// GetFactories returns all registered factories
func (m *Manager) GetFactories() map[component.Type]Factory {
	m.mu.RLock()
	defer m.mu.RUnlock()

	factories := make(map[component.Type]Factory, len(m.factories))
	for k, v := range m.factories {
		factories[k] = v
	}
	return factories
}

// Start starts all plugins
func (m *Manager) Start(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for id, plugin := range m.plugins {
		if err := plugin.Start(ctx, m.host); err != nil {
			m.logger.Error("failed to start plugin",
				zap.String("id", id.String()),
				zap.Error(err))
			return fmt.Errorf("failed to start plugin %s: %w", id.String(), err)
		}
		m.logger.Info("started plugin", zap.String("id", id.String()))
	}
	return nil
}

// Shutdown stops all plugins
func (m *Manager) Shutdown(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var lastErr error
	for id, plugin := range m.plugins {
		if err := plugin.Shutdown(ctx); err != nil {
			m.logger.Error("failed to shutdown plugin",
				zap.String("id", id.String()),
				zap.Error(err))
			lastErr = err
		} else {
			m.logger.Info("shutdown plugin", zap.String("id", id.String()))
		}
	}
	return lastErr
}

// CreateCollectorFactories creates OpenTelemetry collector factories from registered plugins
func (m *Manager) CreateCollectorFactories() (otelcol.Factories, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	factories := otelcol.Factories{
		Receivers:  make(map[component.Type]receiver.Factory),
		Processors: make(map[component.Type]processor.Factory),
		Exporters:  make(map[component.Type]exporter.Factory),
		Extensions: make(map[component.Type]extension.Factory),
	}

	for typ, factory := range m.factories {
		switch f := factory.(type) {
		case receiver.Factory:
			factories.Receivers[typ] = f
		case processor.Factory:
			factories.Processors[typ] = f
		case exporter.Factory:
			factories.Exporters[typ] = f
		case extension.Factory:
			factories.Extensions[typ] = f
		}
	}

	return factories, nil
}
