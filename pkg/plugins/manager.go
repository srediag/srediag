package plugins

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"
)

// Manager handles plugin lifecycle and management
type Manager struct {
	logger    *zap.Logger
	mu        sync.RWMutex
	plugins   map[string]Plugin
	factories map[string]Factory
}

// NewManager creates a new plugin manager
func NewManager(logger *zap.Logger) *Manager {
	return &Manager{
		logger:    logger,
		plugins:   make(map[string]Plugin),
		factories: make(map[string]Factory),
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
	m.logger.Info("registered plugin factory", zap.String("type", typ))
	return nil
}

// CreatePlugin creates a new plugin instance
func (m *Manager) CreatePlugin(pluginType string, config interface{}) (Plugin, error) {
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

	id := fmt.Sprintf("%s/%s", plugin.Type(), plugin.Name())
	if _, exists := m.plugins[id]; exists {
		return nil, fmt.Errorf("plugin already exists: %s", id)
	}

	m.plugins[id] = plugin
	m.logger.Info("created plugin",
		zap.String("id", id),
		zap.String("version", plugin.Version()))

	return plugin, nil
}

// GetPlugin returns a plugin by its ID
func (m *Manager) GetPlugin(id string) (Plugin, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	plugin, exists := m.plugins[id]
	return plugin, exists
}

// ListPlugins returns all registered plugins
func (m *Manager) ListPlugins() []Plugin {
	m.mu.RLock()
	defer m.mu.RUnlock()

	plugins := make([]Plugin, 0, len(m.plugins))
	for _, p := range m.plugins {
		plugins = append(plugins, p)
	}
	return plugins
}

// mockHost implements component.Host for testing
type mockHost struct {
	logger *zap.Logger
}

func newMockHost() component.Host {
	return &mockHost{
		logger: zap.L(),
	}
}

// ReportFatalError implements component.Host
func (h *mockHost) ReportFatalError(err error) {
	h.logger.Fatal("fatal error reported", zap.Error(err))
}

// GetFactory implements component.Host
func (h *mockHost) GetFactory(component.Type) component.Factory {
	return nil
}

// GetExtensions implements component.Host
func (h *mockHost) GetExtensions() map[component.ID]component.Component {
	return nil
}

// Start starts all plugins
func (m *Manager) Start(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	host := newMockHost()
	for id, plugin := range m.plugins {
		if err := plugin.Start(ctx, host); err != nil {
			m.logger.Error("failed to start plugin",
				zap.String("id", id),
				zap.Error(err))
			return fmt.Errorf("failed to start plugin %s: %w", id, err)
		}
		m.logger.Info("started plugin", zap.String("id", id))
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
				zap.String("id", id),
				zap.Error(err))
			lastErr = err
		} else {
			m.logger.Info("shutdown plugin", zap.String("id", id))
		}
	}
	return lastErr
}
