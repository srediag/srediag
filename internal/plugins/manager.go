package plugins

import (
	"context"
	"fmt"
	"path/filepath"
	"plugin"
	"sync"

	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/config"
	"github.com/srediag/srediag/internal/lifecycle"
)

// Plugin represents the SREDIAG plugin interface
type Plugin interface {
	// Init initializes the plugin with configuration
	Init(config map[string]interface{}) error

	// Start starts the plugin
	Start(ctx context.Context) error

	// Stop stops the plugin
	Stop(ctx context.Context) error

	// Info returns plugin metadata
	Info() Info
}

// Info contains plugin metadata
type Info struct {
	Name        string
	Version     string
	Type        string
	Description string
	Author      string
}

// Manager manages plugin lifecycle
type Manager struct {
	*lifecycle.BaseManager
	config  config.PluginsConfig
	logger  *zap.Logger
	plugins map[string]Plugin
	mu      sync.RWMutex
}

// NewManager creates a new plugin manager instance
func NewManager(cfg config.PluginsConfig, logger *zap.Logger) *Manager {
	return &Manager{
		BaseManager: lifecycle.NewBaseManager(),
		config:      cfg,
		logger:      logger,
		plugins:     make(map[string]Plugin),
	}
}

// Start initializes and starts all enabled plugins
func (m *Manager) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := m.CheckRunningState(false); err != nil {
		return err
	}

	// Load plugins from configured directory
	if err := m.loadPlugins(); err != nil {
		return fmt.Errorf("failed to initialize: %v", err)
	}

	// Initialize and start each enabled plugin
	for name, p := range m.plugins {
		if !m.isEnabled(name) {
			m.logger.Info("skipping disabled plugin", zap.String("name", name))
			continue
		}

		if err := p.Init(m.config.Settings[name]); err != nil {
			return fmt.Errorf("failed to initialize plugin %s: %v", name, err)
		}

		if err := p.Start(ctx); err != nil {
			return fmt.Errorf("failed to start plugin %s: %v", name, err)
		}

		m.logger.Info("plugin started", zap.String("name", name))
	}

	m.SetRunning(true)
	m.logger.Info("plugin manager started successfully")

	return nil
}

// Stop gracefully stops all plugins
func (m *Manager) Stop(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := m.CheckRunningState(true); err != nil {
		return err
	}

	var errs []error

	// Stop each enabled plugin
	for name, p := range m.plugins {
		if !m.isEnabled(name) {
			continue
		}

		if err := p.Stop(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to stop plugin %s: %v", name, err))
		}
		m.logger.Info("plugin stopped", zap.String("name", name))
	}

	m.SetRunning(false)

	if len(errs) > 0 {
		return fmt.Errorf("failed to stop plugins: %v", errs)
	}

	m.logger.Info("plugin manager stopped successfully")
	return nil
}

// LoadPlugin loads a plugin from the filesystem
func (m *Manager) LoadPlugin(path string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Open plugin
	plug, err := plugin.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open plugin %s: %w", path, err)
	}

	// Look up plugin symbol
	sym, err := plug.Lookup("New")
	if err != nil {
		return fmt.Errorf("plugin %s does not export 'New' symbol: %w", path, err)
	}

	// Verify plugin interface
	p, ok := sym.(Plugin)
	if !ok {
		return fmt.Errorf("plugin %s does not implement Plugin interface", path)
	}

	// Get plugin info
	info := p.Info()

	// Register plugin
	m.plugins[info.Name] = p
	m.logger.Info("plugin loaded successfully",
		zap.String("name", info.Name),
		zap.String("version", info.Version),
		zap.String("type", info.Type),
	)

	return nil
}

// loadPlugins loads all plugins from the configured directory
func (m *Manager) loadPlugins() error {
	pattern := filepath.Join(m.config.Directory, "*.so")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("failed to list plugins: %w", err)
	}

	for _, path := range matches {
		if err := m.LoadPlugin(path); err != nil {
			m.logger.Error("failed to load plugin",
				zap.String("path", path),
				zap.Error(err),
			)
		}
	}

	return nil
}

// GetPlugin returns a plugin by name
func (m *Manager) GetPlugin(name string) (Plugin, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.plugins[name]
	return p, ok
}

// ListPlugins returns a list of all loaded plugins
func (m *Manager) ListPlugins() []Info {
	m.mu.RLock()
	defer m.mu.RUnlock()

	plugins := make([]Info, 0, len(m.plugins))
	for _, p := range m.plugins {
		plugins = append(plugins, p.Info())
	}
	return plugins
}

// isEnabled checks if a plugin is enabled in the configuration
func (m *Manager) isEnabled(name string) bool {
	if len(m.config.Enabled) == 0 {
		return true // if no plugins are explicitly enabled, all are enabled
	}
	for _, enabled := range m.config.Enabled {
		if enabled == name {
			return true
		}
	}
	return false
}

// GetPluginConfig returns the configuration for a specific plugin
func (m *Manager) GetPluginConfig(name string) map[string]interface{} {
	return m.config.Settings[name]
}
