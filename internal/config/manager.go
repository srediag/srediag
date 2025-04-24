package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// Manager handles configuration loading and validation
type Manager struct {
	logger     *zap.Logger
	configPath string
	buildInfo  component.BuildInfo
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

// NewManager creates a new configuration manager
func NewManager(logger *zap.Logger, configPath string, buildInfo component.BuildInfo) *Manager {
	return &Manager{
		logger:     logger,
		configPath: configPath,
		buildInfo:  buildInfo,
	}
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
	return m.buildInfo
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
