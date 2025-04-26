package config

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// ServiceConfig holds service-level configuration
type ServiceConfig struct {
	LogLevel string `yaml:"logLevel"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
}

// ComponentsConfig holds component configuration
type ComponentsConfig struct {
	Receivers  map[string]interface{} `yaml:"receivers"`
	Processors map[string]interface{} `yaml:"processors"`
	Exporters  map[string]interface{} `yaml:"exporters"`
	Extensions map[string]interface{} `yaml:"extensions"`
	Connectors map[string]interface{} `yaml:"connectors"`
}

// Manager handles configuration loading and validation
type Manager struct {
	logger *zap.Logger
	config Config
}

// NewManager creates a new configuration manager
func NewManager(logger *zap.Logger) *Manager {
	return &Manager{
		logger: logger,
	}
}

// Load loads configuration from a file
func (m *Manager) Load(configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, &m.config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	return nil
}

// GetConfig returns the loaded configuration
func (m *Manager) GetConfig() Config {
	return m.config
}
