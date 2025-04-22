// internal/config/config.go
package config

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/srediag/srediag/internal/config/diagnostic"
)

// SREDiagConfig represents the main SREDIAG configuration
type SREDiagConfig struct {
	Version    string                  `mapstructure:"version"`
	Service    *SREDiagServiceConfig   `mapstructure:"service"`
	Collector  *SREDiagCollectorConfig `mapstructure:"collector"`
	Diagnostic *diagnostic.Config      `mapstructure:"diagnostic"`
}

// SREDiagServiceConfig contains service configurations
type SREDiagServiceConfig struct {
	Name        string `mapstructure:"name"`
	Environment string `mapstructure:"environment"`
	Version     string `mapstructure:"version"`
}

// SREDiagCollectorConfig contains collector configurations
type SREDiagCollectorConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	ConfigPath string `mapstructure:"config_path"`
}

// IsEnabled returns whether the collector is enabled
func (c *SREDiagCollectorConfig) IsEnabled() bool {
	return c.Enabled
}

// GetConfigPath returns the collector configuration path
func (c *SREDiagCollectorConfig) GetConfigPath() string {
	return c.ConfigPath
}

// LoadConfig loads configuration from file
func LoadConfig(filePath string) (*SREDiagConfig, error) {
	v := viper.New()
	v.SetConfigFile(filePath)

	setDefaultConfig(v)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("config file not found: %s", filePath)
		}
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg SREDiagConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error decoding configuration: %w", err)
	}

	if err := validateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// validateConfig validates the configuration
func validateConfig(cfg *SREDiagConfig) error {
	if cfg.Service == nil {
		return fmt.Errorf("service configuration is required")
	}
	if cfg.Service.Name == "" {
		return fmt.Errorf("service name is required")
	}
	if cfg.Service.Environment == "" {
		return fmt.Errorf("service environment is required")
	}
	if cfg.Service.Version == "" {
		return fmt.Errorf("service version is required")
	}
	if cfg.Collector == nil {
		return fmt.Errorf("collector configuration is required")
	}
	if cfg.Diagnostic == nil {
		return fmt.Errorf("diagnostic configuration is required")
	}
	return nil
}

// setDefaultConfig sets default values for all configurations
func setDefaultConfig(v *viper.Viper) {
	// Service defaults
	v.SetDefault("version", "v0.1.0")
	v.SetDefault("service.name", "srediag")
	v.SetDefault("service.environment", "production")
	v.SetDefault("service.version", "v0.1.0")

	// Collector defaults
	v.SetDefault("collector.enabled", true)
	v.SetDefault("collector.config_path", "configs/otel-config.yaml")

	// Diagnostic defaults
	v.SetDefault("diagnostic.system.enabled", true)
	v.SetDefault("diagnostic.kubernetes.enabled", false)
	v.SetDefault("diagnostic.cloud.enabled", false)
	v.SetDefault("diagnostic.security.enabled", false)
}
