// Package types provides configuration types for SREDIAG
package types

import (
	"github.com/srediag/srediag/internal/core"
)

// CollectorConfig represents the OpenTelemetry collector configuration
type CollectorConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	ConfigPath string `mapstructure:"config_path"`
}

// IsEnabled implements ISREDiagCollectorConfig
func (c *CollectorConfig) IsEnabled() bool {
	return c.Enabled
}

// GetConfigPath implements ISREDiagCollectorConfig
func (c *CollectorConfig) GetConfigPath() string {
	return c.ConfigPath
}

// Ensure CollectorConfig implements ISREDiagCollectorConfig
var _ core.ISREDiagCollectorConfig = (*CollectorConfig)(nil)
