package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	"github.com/srediag/srediag/internal/types"
)

// Manager manages configuration
type Manager struct {
	mu         sync.RWMutex
	logger     *zap.Logger
	configs    map[string]interface{}
	configPath string
}

// PipelineConfig represents a data pipeline configuration
type PipelineConfig struct {
	Receivers  []string `yaml:"receivers"`
	Processors []string `yaml:"processors"`
	Exporters  []string `yaml:"exporters"`
}

// ServiceConfig represents the service configuration
type ServiceConfig struct {
	Pipelines map[string]PipelineConfig `yaml:"pipelines"`
}

// CollectorConfig represents the full collector configuration
type CollectorConfig struct {
	Receivers  map[string]interface{} `yaml:"receivers"`
	Processors map[string]interface{} `yaml:"processors"`
	Exporters  map[string]interface{} `yaml:"exporters"`
	Service    ServiceConfig          `yaml:"service"`
}

// NewManager creates a new config manager
func NewManager(logger *zap.Logger, configPath string) *Manager {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &Manager{
		logger:     logger,
		configs:    make(map[string]interface{}),
		configPath: configPath,
	}
}

// LoadConfig loads configuration from file
func (m *Manager) LoadConfig() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.configPath == "" {
		return fmt.Errorf("config path not set")
	}

	data, err := os.ReadFile(m.configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var config map[string]interface{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	m.configs = config
	return nil
}

// SaveConfig saves configuration to file
func (m *Manager) SaveConfig() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.configPath == "" {
		return fmt.Errorf("config path not set")
	}

	data, err := yaml.Marshal(m.configs)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(m.configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(m.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Load loads and validates configuration
func (m *Manager) Load(ctx context.Context) (*confmap.Conf, error) {
	// Verify config file exists and is accessible
	absPath, err := m.verifyConfigFile()
	if err != nil {
		return nil, err
	}

	// Read configuration file
	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML configuration
	var config CollectorConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate configuration
	if err := m.validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// Convert to confmap.Conf
	conf := confmap.NewFromStringMap(map[string]interface{}{
		"receivers":  config.Receivers,
		"processors": config.Processors,
		"exporters":  config.Exporters,
		"service":    config.Service,
	})

	return conf, nil
}

// BuildInfo returns the collector build information
func (m *Manager) BuildInfo() component.BuildInfo {
	return component.BuildInfo{}
}

func (m *Manager) verifyConfigFile() (string, error) {
	if m.configPath == "" {
		return "", fmt.Errorf("config path is required")
	}

	absPath, err := filepath.Abs(m.configPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return "", fmt.Errorf("config file not found: %s", absPath)
	}

	return absPath, nil
}

func (m *Manager) validateConfig(config *CollectorConfig) error {
	if len(config.Service.Pipelines) == 0 {
		return fmt.Errorf("at least one pipeline must be configured")
	}

	for name, pipeline := range config.Service.Pipelines {
		if err := m.validatePipeline(name, pipeline); err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) validatePipeline(name string, pipeline PipelineConfig) error {
	if len(pipeline.Receivers) == 0 {
		return fmt.Errorf("pipeline %s must have at least one receiver", name)
	}

	if len(pipeline.Exporters) == 0 {
		return fmt.Errorf("pipeline %s must have at least one exporter", name)
	}

	return nil
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

// DeletePluginConfig deletes plugin configuration
func (m *Manager) DeletePluginConfig(ctx context.Context, name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.configs, name)
	m.logger.Info("Deleted plugin configuration", zap.String("name", name))
	return nil
}

// Initialize initializes the config manager
func (m *Manager) Initialize(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.logger.Info("Initializing config manager")
	return nil
}

// Shutdown shuts down the config manager
func (m *Manager) Shutdown(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.logger.Info("Shutting down config manager")
	return nil
}
