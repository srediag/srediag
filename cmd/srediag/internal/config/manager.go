package config

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/collector/confmap"
	"go.uber.org/zap"

	baseconfig "github.com/srediag/srediag/internal/config"
	"github.com/srediag/srediag/internal/types"
)

// Manager manages the configuration for the CLI
type Manager struct {
	mu       sync.RWMutex
	logger   *zap.Logger
	provider *baseconfig.Provider
	config   *types.Config
}

// NewManager creates a new configuration manager
func NewManager(logger *zap.Logger) *Manager {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &Manager{
		logger: logger,
	}
}

// Initialize initializes the configuration manager
func (m *Manager) Initialize(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Create the base provider
	m.provider = baseconfig.NewProvider(m.logger, DefaultSettings())

	// Create and register the CLI provider
	cliProvider := NewCLIProvider()
	if err := m.provider.RegisterProvider("cli", cliProvider); err != nil {
		return fmt.Errorf("failed to register CLI provider: %w", err)
	}

	// Add configuration validator
	m.provider.AddValidator(baseconfig.NewBaseValidator(baseconfig.ConfigVersion{
		Major: 1,
		Minor: 0,
		Patch: 0,
	}))

	// Add configuration watcher
	m.provider.AddWatcher(func(event *confmap.ChangeEvent) {
		if event.Error != nil {
			m.logger.Error("Configuration change error",
				zap.Error(event.Error))
			return
		}

		m.logger.Info("Configuration changed")
		if err := m.reload(ctx); err != nil {
			m.logger.Error("Failed to reload configuration",
				zap.Error(err))
		}
	})

	// Initial load
	return m.reload(ctx)
}

// reload reloads the configuration
func (m *Manager) reload(ctx context.Context) error {
	conf, err := m.provider.Load(ctx)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	var cfg types.Config
	if err := conf.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	m.config = &cfg
	return nil
}

// Get returns the current configuration
func (m *Manager) Get() *types.Config {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.config
}

// Shutdown shuts down the configuration manager
func (m *Manager) Shutdown(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.provider != nil {
		if err := m.provider.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown provider: %w", err)
		}
	}

	return nil
}
