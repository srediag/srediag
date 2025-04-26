package plugin

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"
)

// Manager handles plugin lifecycle and communication
type Manager struct {
	logger     *zap.Logger
	pluginsDir string
	plugins    map[component.Type]*IPCPlugin
	mu         sync.RWMutex
}

// NewManager creates a new plugin manager
func NewManager(logger *zap.Logger, pluginsDir string) *Manager {
	return &Manager{
		logger:     logger,
		pluginsDir: pluginsDir,
		plugins:    make(map[component.Type]*IPCPlugin),
	}
}

// LoadPlugin loads and starts a plugin
func (m *Manager) LoadPlugin(ctx context.Context, pluginType component.Type, name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if plugin is already loaded
	if _, exists := m.plugins[pluginType]; exists {
		return fmt.Errorf("plugin type %s already loaded", pluginType)
	}

	// Ensure plugin directory exists
	pluginDir := filepath.Join(m.pluginsDir, pluginType.String())
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		return fmt.Errorf("failed to create plugin directory: %w", err)
	}

	// Construct plugin path
	pluginPath := filepath.Join(pluginDir, name)
	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		return fmt.Errorf("plugin not found: %s", pluginPath)
	}

	// Create new plugin instance
	plugin, err := NewIPCPlugin(ctx, pluginPath, pluginType)
	if err != nil {
		return fmt.Errorf("failed to create plugin: %w", err)
	}

	// Store plugin
	m.plugins[pluginType] = plugin
	m.logger.Info("Loaded plugin",
		zap.String("type", pluginType.String()),
		zap.String("name", name),
		zap.String("path", pluginPath))

	return nil
}

// GetFactory retrieves a component factory from a plugin
func (m *Manager) GetFactory(pluginType component.Type) (component.Factory, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	plugin, exists := m.plugins[pluginType]
	if !exists {
		return nil, fmt.Errorf("plugin type %s not loaded", pluginType)
	}

	// Get factory type
	typeResp, err := plugin.Send("type", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get factory type: %w", err)
	}

	// Get default config
	configResp, err := plugin.Send("config", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get default config: %w", err)
	}

	// Create proxy factory
	return &proxyFactory{
		plugin:        plugin,
		factoryType:   typeResp.Data.(component.Type),
		defaultConfig: configResp.Data.(component.Config),
	}, nil
}

// Close closes all plugins
func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errs []error
	for _, plugin := range m.plugins {
		if err := plugin.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to close plugins: %v", errs)
	}
	return nil
}

// proxyFactory implements component.Factory interface
type proxyFactory struct {
	plugin        *IPCPlugin
	factoryType   component.Type
	defaultConfig component.Config
}

func (f *proxyFactory) Type() component.Type {
	return f.factoryType
}

func (f *proxyFactory) CreateDefaultConfig() component.Config {
	return f.defaultConfig
}
