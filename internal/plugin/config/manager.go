package config

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/types"
)

// Manager handles plugin configuration management
type Manager struct {
	mu      sync.RWMutex
	logger  *zap.Logger
	configs map[component.ID]*types.PluginConfig
}

// NewManager creates a new plugin configuration manager
func NewManager(logger *zap.Logger) *Manager {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &Manager{
		logger:  logger,
		configs: make(map[component.ID]*types.PluginConfig),
	}
}

// LoadPluginConfig loads and validates plugin configuration
func (m *Manager) LoadPluginConfig(ctx context.Context, plugin types.IPlugin) error {
	if plugin == nil {
		return fmt.Errorf("plugin cannot be nil")
	}

	settings := types.ComponentSettings{
		"name":     plugin.GetName(),
		"version":  plugin.GetVersion(),
		"type":     plugin.GetType(),
		"category": plugin.GetCategory(),
	}

	return plugin.Configure(settings)
}

// GetPluginConfig retrieves a plugin configuration by ID
func (m *Manager) GetPluginConfig(id component.ID) (*types.PluginConfig, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	config, ok := m.configs[id]
	if !ok {
		m.logger.Debug("Plugin configuration not found",
			zap.String("plugin_id", id.String()))
		return nil, types.ErrPluginNotFound
	}

	return config, nil
}

// DeletePluginConfig removes a plugin configuration
func (m *Manager) DeletePluginConfig(id component.ID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.configs[id]; !ok {
		m.logger.Debug("Cannot delete: plugin configuration not found",
			zap.String("plugin_id", id.String()))
		return types.ErrPluginNotFound
	}

	delete(m.configs, id)
	m.logger.Info("Deleted plugin configuration",
		zap.String("plugin_id", id.String()))
	return nil
}

// ListPluginConfigs returns all plugin configurations
func (m *Manager) ListPluginConfigs() map[component.ID]*types.PluginConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Return a copy to prevent external modifications
	configs := make(map[component.ID]*types.PluginConfig, len(m.configs))
	for id, config := range m.configs {
		configs[id] = config
	}

	return configs
}

// Shutdown performs cleanup when the manager is stopped
func (m *Manager) Shutdown(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.configs = make(map[component.ID]*types.PluginConfig)
	m.logger.Info("Plugin configuration manager shutdown complete")
	return nil
}
