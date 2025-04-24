package core

import (
	"fmt"

	"go.opentelemetry.io/collector/component"
)

// ComponentConfig represents the configuration for a component
type ComponentConfig struct {
	Type    string           `mapstructure:"type"`
	Enabled bool             `mapstructure:"enabled"`
	Config  component.Config `mapstructure:"config"`
}

// Config represents the main configuration structure
type Config struct {
	Components map[string]*ComponentConfig `mapstructure:"components"`
}

// LoadConfig loads the configuration from a file
func LoadConfig(path string) (*Config, error) {
	// TODO: Implement configuration loading
	return &Config{
		Components: make(map[string]*ComponentConfig),
	}, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Components == nil {
		return fmt.Errorf("components configuration is required")
	}

	for name, comp := range c.Components {
		if comp.Type == "" {
			return fmt.Errorf("component %s: type is required", name)
		}
	}

	return nil
}
