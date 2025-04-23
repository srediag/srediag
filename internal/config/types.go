package config

import (
	"fmt"

	"github.com/srediag/srediag/internal/config/diagnostic"
	"github.com/srediag/srediag/internal/config/types"
	"github.com/srediag/srediag/internal/core"
)

// Config represents the main SREDIAG configuration
type Config struct {
	Core       types.CoreConfig      `mapstructure:"core"`
	Service    types.ServiceConfig   `mapstructure:"service"`
	Telemetry  types.TelemetryConfig `mapstructure:"telemetry"`
	Plugins    types.PluginsConfig   `mapstructure:"plugins"`
	Collector  types.CollectorConfig `mapstructure:"collector"`
	Diagnostic diagnostic.Config     `mapstructure:"diagnostic"`
}

// Ensure Config implements ISREDiagConfig
var _ core.ISREDiagConfig = (*Config)(nil)

// GetVersion implements ISREDiagConfig
func (c *Config) GetVersion() string {
	return c.Service.Version
}

// GetServiceConfig implements ISREDiagConfig
func (c *Config) GetServiceConfig() core.ISREDiagServiceConfig {
	return &c.Service
}

// GetCollectorConfig implements ISREDiagConfig
func (c *Config) GetCollectorConfig() core.ISREDiagCollectorConfig {
	return &c.Collector
}

// Validate implements ISREDiagConfig
func (c *Config) Validate() error {
	if c.Service.Name == "" {
		return fmt.Errorf("service name is required")
	}
	if c.Service.Version == "" {
		return fmt.Errorf("service version is required")
	}
	return nil
}
