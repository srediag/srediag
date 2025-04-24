package config

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/collector/confmap"
	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/factory"
	"github.com/srediag/srediag/internal/types"
)

// OtelManager manages configuration
type OtelManager struct {
	mu          sync.RWMutex
	logger      *zap.Logger
	provider    ConfigProvider
	factories   map[string]*factory.Factory
	settings    *types.ServiceSettings
	watchCancel context.CancelFunc
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

	// Create a new context for watching that we can cancel later
	watchCtx, cancel := context.WithCancel(ctx)
	m.watchCancel = cancel

	// Start watching for configuration changes
	go m.watchConfig(watchCtx)

	// Load initial configuration
	if err := m.loadConfig(ctx); err != nil {
		return fmt.Errorf("failed to load initial configuration: %w", err)
	}

	m.logger.Info("Started configuration manager")
	return nil
}

// Stop stops the configuration manager
func (m *OtelManager) Stop(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Stop watching for changes
	if m.watchCancel != nil {
		m.watchCancel()
		m.watchCancel = nil
	}

	// Shutdown the provider
	if err := m.provider.Shutdown(ctx); err != nil {
		m.logger.Error("Failed to shutdown config provider", zap.Error(err))
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
	// Load configuration periodically
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := m.loadConfig(ctx); err != nil {
				m.logger.Error("Failed to reload configuration", zap.Error(err))
			}
		}
	}
}

// loadConfig loads the current configuration and updates settings
func (m *OtelManager) loadConfig(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Load configuration from provider
	conf, err := m.provider.Retrieve(ctx, m.provider.Scheme()+":default", nil)
	if err != nil {
		return fmt.Errorf("failed to retrieve configuration: %w", err)
	}

	// Convert configuration to service settings
	if err := m.updateSettings(conf); err != nil {
		return fmt.Errorf("failed to update settings: %w", err)
	}

	return nil
}

// updateSettings updates the service settings from configuration
func (m *OtelManager) updateSettings(conf *confmap.Conf) error {
	raw := conf.ToStringMap()
	serviceConfig, ok := raw["service"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("service configuration not found or invalid")
	}

	// Convert security configuration
	var security types.SecurityConfig
	if secMap, ok := serviceConfig["security"].(map[string]interface{}); ok {
		tlsConfig, ok := secMap["tls"].(map[string]interface{})
		if ok {
			security.TLS = types.TLSConfig{
				Enabled:    getBool(tlsConfig, "enabled", false),
				CertFile:   getString(tlsConfig, "cert_file", ""),
				KeyFile:    getString(tlsConfig, "key_file", ""),
				CAFile:     getString(tlsConfig, "ca_file", ""),
				ServerName: getString(tlsConfig, "server_name", ""),
				SkipVerify: getBool(tlsConfig, "skip_verify", false),
				MinVersion: getString(tlsConfig, "min_version", "TLS1.2"),
				MaxVersion: getString(tlsConfig, "max_version", ""),
			}
		}
		security.Enabled = getBool(secMap, "enabled", false)
		security.CertFile = getString(secMap, "cert_file", "")
		security.KeyFile = getString(secMap, "key_file", "")
	}

	// Convert settings
	settings := getStringMap(serviceConfig, "settings")

	// Convert type to ComponentType
	typeStr := getString(serviceConfig, "type", "")
	var componentType types.ComponentType
	switch typeStr {
	case "core":
		componentType = types.ComponentTypeCore
	case "service":
		componentType = types.ComponentTypeService
	case "receiver":
		componentType = types.ComponentTypeReceiver
	case "processor":
		componentType = types.ComponentTypeProcessor
	case "exporter":
		componentType = types.ComponentTypeExporter
	case "extension":
		componentType = types.ComponentTypeExtension
	case "plugin":
		componentType = types.ComponentTypePlugin
	default:
		componentType = types.ComponentTypeUnknown
	}

	m.settings = &types.ServiceSettings{
		Name:        getString(serviceConfig, "name", ""),
		Version:     getString(serviceConfig, "version", ""),
		Environment: getString(serviceConfig, "environment", ""),
		Type:        componentType,
		Security:    security,
		Settings:    settings,
	}

	return nil
}

// getString safely gets a string value from a map
func getString(m map[string]interface{}, key string, defaultValue string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return defaultValue
}

// getBool safely gets a boolean value from a map
func getBool(m map[string]interface{}, key string, defaultValue bool) bool {
	if val, ok := m[key].(bool); ok {
		return val
	}
	return defaultValue
}

// getStringMap safely gets a string map from a map
func getStringMap(m map[string]interface{}, key string) map[string]string {
	result := make(map[string]string)
	if val, ok := m[key].(map[string]interface{}); ok {
		for k, v := range val {
			if str, ok := v.(string); ok {
				result[k] = str
			}
		}
	}
	return result
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
