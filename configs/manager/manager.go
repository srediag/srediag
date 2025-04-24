package manager

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.uber.org/zap"

	"github.com/srediag/srediag/configs/provider"
	"github.com/srediag/srediag/internal/config"
	"github.com/srediag/srediag/internal/types"
)

// ConfigManager manages multiple configuration files
type ConfigManager struct {
	mu       sync.RWMutex
	logger   *zap.Logger
	provider *config.Provider
	configs  map[string]*confmap.Conf
}

// NewConfigManager creates a new configuration manager
func NewConfigManager(logger *zap.Logger) *ConfigManager {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &ConfigManager{
		logger:  logger,
		configs: make(map[string]*confmap.Conf),
	}
}

// Initialize initializes the configuration manager
func (m *ConfigManager) Initialize(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Create the base provider
	settings := component.TelemetrySettings{
		Logger: m.logger,
	}
	m.provider = config.NewProvider(m.logger, settings)

	// Create and register the YAML provider
	yamlProvider := provider.NewYAMLProvider()
	if err := m.provider.RegisterProvider("yaml", yamlProvider); err != nil {
		return fmt.Errorf("failed to register YAML provider: %w", err)
	}

	// Add configuration validator
	m.provider.AddValidator(config.NewBaseValidator(config.ConfigVersion{
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

	// Load all configuration files
	configFiles := DefaultConfigFiles()

	for _, file := range configFiles {
		// Find configuration file
		filePath, err := FindConfigFile(file.Name)
		if err != nil {
			if file.Required {
				return fmt.Errorf("required configuration file not found: %w", err)
			}
			m.logger.Warn("Optional configuration file not found",
				zap.String("file", file.Name),
				zap.Error(err))
			continue
		}

		// Validate configuration file
		if err := ValidateConfigFile(filePath); err != nil {
			return fmt.Errorf("invalid configuration file %q: %w", file.Name, err)
		}

		// Load configuration
		if err := m.loadConfig(ctx, filePath); err != nil {
			return fmt.Errorf("failed to load %q: %w", file.Name, err)
		}
	}

	return nil
}

// loadConfig loads a specific configuration file
func (m *ConfigManager) loadConfig(ctx context.Context, filename string) error {
	// Get absolute path
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for %q: %w", filename, err)
	}

	// Register the file with the YAML provider
	if err := m.provider.RegisterProvider(absPath, provider.NewYAMLProvider()); err != nil {
		return fmt.Errorf("failed to register provider for %q: %w", absPath, err)
	}

	// Load configuration
	conf, err := m.provider.Load(ctx)
	if err != nil {
		return fmt.Errorf("failed to load configuration from %q: %w", absPath, err)
	}

	m.configs[filename] = conf
	return nil
}

// GetConfig returns a specific configuration
func (m *ConfigManager) GetConfig(filename string) (*confmap.Conf, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	conf, ok := m.configs[filename]
	return conf, ok
}

// GetSREDiagConfig returns the main SREDIAG configuration
func (m *ConfigManager) GetSREDiagConfig() (*types.Config, error) {
	conf, ok := m.GetConfig("srediag.yaml")
	if !ok {
		return nil, fmt.Errorf("SREDIAG configuration not found")
	}

	var cfg types.Config
	if err := conf.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal SREDIAG configuration: %w", err)
	}

	return &cfg, nil
}

// GetOTelConfig returns the OpenTelemetry configuration
func (m *ConfigManager) GetOTelConfig() (*confmap.Conf, error) {
	conf, ok := m.GetConfig("otel-config.yaml")
	if !ok {
		return nil, fmt.Errorf("OpenTelemetry configuration not found")
	}

	return conf, nil
}

// reload reloads all configurations
func (m *ConfigManager) reload(ctx context.Context) error {
	for filename := range m.configs {
		if err := m.loadConfig(ctx, filename); err != nil {
			return fmt.Errorf("failed to reload %q: %w", filename, err)
		}
	}
	return nil
}

// Shutdown shuts down the configuration manager
func (m *ConfigManager) Shutdown(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.provider != nil {
		if err := m.provider.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown provider: %w", err)
		}
	}

	return nil
}
