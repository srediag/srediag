package config

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/factory"
	"github.com/srediag/srediag/internal/types"
)

// OtelManager manages configuration
type OtelManager struct {
	mu        sync.RWMutex
	logger    *zap.Logger
	provider  ConfigProvider
	factories map[string]*factory.Factory
	settings  *types.ServiceSettings
}

// NewOtelManager creates a new configuration manager
func NewOtelManager(logger *zap.Logger, provider ConfigProvider) *OtelManager {
	return &OtelManager{
		logger:    logger,
		provider:  provider,
		factories: make(map[string]*factory.Factory),
		settings:  &types.ServiceSettings{},
	}
}

// Start starts the configuration manager
func (m *OtelManager) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := m.provider.Start(ctx); err != nil {
		return fmt.Errorf("failed to start config provider: %w", err)
	}

	// Watch for configuration changes
	go m.watchConfig(ctx)

	m.logger.Info("Started configuration manager")
	return nil
}

// Stop stops the configuration manager
func (m *OtelManager) Stop(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := m.provider.Stop(ctx); err != nil {
		m.logger.Error("Failed to stop config provider", zap.Error(err))
	}

	m.logger.Info("Stopped configuration manager")
	return nil
}

// RegisterFactory registers a component factory
func (m *OtelManager) RegisterFactory(factory *factory.Factory) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := factory.GetID()
	if _, exists := m.factories[id]; exists {
		return fmt.Errorf("factory already registered for id %s", id)
	}

	m.factories[id] = factory
	m.logger.Info("Registered factory", zap.String("id", id))

	return nil
}

// UnregisterFactory removes a component factory
func (m *OtelManager) UnregisterFactory(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.factories[id]; !exists {
		return fmt.Errorf("no factory registered for id %s", id)
	}

	delete(m.factories, id)
	m.logger.Info("Unregistered factory", zap.String("id", id))

	return nil
}

// GetFactory returns a factory by id
func (m *OtelManager) GetFactory(id string) (*factory.Factory, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	factory, exists := m.factories[id]
	if !exists {
		return nil, fmt.Errorf("no factory registered for id %s", id)
	}

	return factory, nil
}

// ListFactories returns all registered factories
func (m *OtelManager) ListFactories() map[string]*factory.Factory {
	m.mu.RLock()
	defer m.mu.RUnlock()

	factories := make(map[string]*factory.Factory, len(m.factories))
	for id, factory := range m.factories {
		factories[id] = factory
	}

	return factories
}

// GetSettings returns the service settings
func (m *OtelManager) GetSettings() *types.ServiceSettings {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.settings
}

// watchConfig watches for configuration changes
func (m *OtelManager) watchConfig(ctx context.Context) {
	configChan := m.provider.Watch()

	for {
		select {
		case <-ctx.Done():
			return
		case cfg := <-configChan:
			if err := m.handleConfigChange(cfg); err != nil {
				m.logger.Error("Failed to handle configuration change",
					zap.Error(err))
			}
		}
	}
}

// handleConfigChange handles configuration changes
func (m *OtelManager) handleConfigChange(cfg *types.Config) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	serviceConfig := cfg.GetService()
	if serviceConfig != nil {
		// Convert ServiceConfig to ServiceSettings
		m.settings = &types.ServiceSettings{
			Name:        serviceConfig.GetName(),
			Version:     serviceConfig.GetVersion(),
			Environment: serviceConfig.GetEnvironment(),
			Type:        serviceConfig.GetType(),
			Security:    serviceConfig.GetSecurity(),
			Settings:    make(map[string]string),
		}
		m.logger.Info("Service configuration changed")
	}

	return nil
}

// RegisterFactories registers multiple factories at once
func (m *OtelManager) RegisterFactories(factories map[string]*factory.Factory) error {
	for _, f := range factories {
		if err := m.RegisterFactory(f); err != nil {
			return err
		}
	}
	return nil
}

// Clear removes all registered factories
func (m *OtelManager) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.factories = make(map[string]*factory.Factory)
	m.logger.Info("Cleared all factory registrations")
}
