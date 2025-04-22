// internal/config/config.go
package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config represents the global application configuration
type Config struct {
	Version   string          `mapstructure:"version"`
	Service   ServiceConfig   `mapstructure:"service"`
	Telemetry TelemetryConfig `mapstructure:"telemetry"`
	Plugins   PluginsConfig   `mapstructure:"plugins"`
	Logging   LoggingConfig   `mapstructure:"logging"`
	Security  SecurityConfig  `mapstructure:"security"`
}

// ServiceConfig represents service configuration
type ServiceConfig struct {
	Name        string `mapstructure:"name"`
	Environment string `mapstructure:"environment"`
}

// TelemetryConfig represents OpenTelemetry configuration
type TelemetryConfig struct {
	Enabled            bool              `mapstructure:"enabled"`
	ServiceName        string            `mapstructure:"service_name"`
	Endpoint           string            `mapstructure:"endpoint"`
	Protocol           string            `mapstructure:"protocol"`
	Environment        string            `mapstructure:"environment"`
	ResourceAttributes map[string]string `mapstructure:"resource_attributes"`
	Sampling           SamplingConfig    `mapstructure:"sampling"`
	Traces             TracesConfig      `mapstructure:"traces"`
	Metrics            MetricsConfig     `mapstructure:"metrics"`
}

// SamplingConfig represents trace sampling configuration
type SamplingConfig struct {
	Type string  `mapstructure:"type"`
	Rate float64 `mapstructure:"rate"`
}

// PluginsConfig represents plugin configuration
type PluginsConfig struct {
	Directory string                            `mapstructure:"directory"`
	Enabled   []string                          `mapstructure:"enabled"`
	Settings  map[string]map[string]interface{} `mapstructure:"settings"`
}

// LoggingConfig represents logging configuration
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

// SecurityConfig represents security configuration
type SecurityConfig struct {
	TLS  TLSConfig  `mapstructure:"tls"`
	Auth AuthConfig `mapstructure:"auth"`
}

// TLSConfig represents TLS configuration
type TLSConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	CertFile string `mapstructure:"cert_file"`
	KeyFile  string `mapstructure:"key_file"`
	CAFile   string `mapstructure:"ca_file"`
}

// AuthConfig represents authentication configuration
type AuthConfig struct {
	Type      string          `mapstructure:"type"`
	TokenFile string          `mapstructure:"token_file"`
	Basic     BasicAuthConfig `mapstructure:"basic"`
	OAuth     OAuthConfig     `mapstructure:"oauth"`
}

// BasicAuthConfig represents basic authentication configuration
type BasicAuthConfig struct {
	Username     string `mapstructure:"username"`
	PasswordFile string `mapstructure:"password_file"`
}

// OAuthConfig represents OAuth configuration
type OAuthConfig struct {
	ClientID         string   `mapstructure:"client_id"`
	ClientSecretFile string   `mapstructure:"client_secret_file"`
	TokenURL         string   `mapstructure:"token_url"`
	Scopes           []string `mapstructure:"scopes"`
}

// TracesConfig represents trace-specific configuration
type TracesConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

// MetricsConfig represents metrics-specific configuration
type MetricsConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

// LoadConfig loads configuration from a file
func LoadConfig(filePath string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(filePath)

	// Set default values before reading the config file
	setDefaults(v)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("config file not found: %s", filePath)
		}
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error decoding configuration: %w", err)
	}

	// Set default values for telemetry configuration
	cfg.Telemetry.setDefaults()

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Validate telemetry configuration
	if c.Telemetry.Enabled {
		if c.Telemetry.ServiceName == "" {
			return fmt.Errorf("service name is required when telemetry is enabled")
		}
		if c.Telemetry.Endpoint == "" {
			return fmt.Errorf("endpoint is required when telemetry is enabled")
		}
	}

	// Validate plugins configuration
	if c.Plugins.Directory == "" {
		return fmt.Errorf("plugins directory is required")
	}
	if _, err := os.Stat(c.Plugins.Directory); os.IsNotExist(err) {
		return fmt.Errorf("plugins directory does not exist: %s", c.Plugins.Directory)
	}

	// Validate logging configuration
	if c.Logging.Level == "" {
		c.Logging.Level = "info" // default value
	}
	if c.Logging.Format == "" {
		c.Logging.Format = "console" // default value
	}
	if c.Logging.Output == "" {
		c.Logging.Output = "stdout" // default value
	}

	return nil
}

// Validate validates telemetry configuration
func (c *TelemetryConfig) Validate() error {
	if !c.Enabled {
		return nil
	}

	if c.ServiceName == "" {
		return fmt.Errorf("service name is required when telemetry is enabled")
	}

	if c.Endpoint == "" {
		return fmt.Errorf("endpoint is required when telemetry is enabled")
	}

	if c.Protocol == "" {
		return fmt.Errorf("protocol is required when telemetry is enabled")
	}

	if !c.Traces.Enabled && !c.Metrics.Enabled {
		return fmt.Errorf("at least one of traces or metrics must be enabled when telemetry is enabled")
	}

	return nil
}

// setDefaults sets default values for configuration
func setDefaults(v *viper.Viper) {
	v.SetDefault("version", "v0.1.0")
	v.SetDefault("service.name", "srediag")
	v.SetDefault("service.environment", "production")

	v.SetDefault("telemetry.enabled", true)
	v.SetDefault("telemetry.service_name", "srediag")
	v.SetDefault("telemetry.endpoint", "http://localhost:4317")
	v.SetDefault("telemetry.protocol", "grpc")
	v.SetDefault("telemetry.environment", "production")
	v.SetDefault("telemetry.sampling.type", "probabilistic")
	v.SetDefault("telemetry.sampling.rate", 0.1)

	v.SetDefault("plugins.directory", "plugins")

	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.format", "console")
	v.SetDefault("logging.output", "stdout")

	v.SetDefault("security.tls.enabled", false)
	v.SetDefault("security.tls.cert_file", "/etc/srediag/certs/server.crt")
	v.SetDefault("security.tls.key_file", "/etc/srediag/certs/server.key")
	v.SetDefault("security.tls.ca_file", "/etc/srediag/certs/ca.crt")
	v.SetDefault("security.auth.type", "none")
}

// setDefaults sets default values for configuration
func (c *TelemetryConfig) setDefaults() {
	if c.Protocol == "" {
		c.Protocol = "grpc"
	}

	if c.ServiceName == "" {
		c.ServiceName = "srediag"
	}

	if c.Endpoint == "" {
		c.Endpoint = "localhost:4317"
	}

	// Default values for traces
	if c.Traces.Enabled {
		if c.Sampling.Rate == 0 {
			c.Sampling.Rate = 0.1
		}
		if c.Sampling.Type == "" {
			c.Sampling.Type = "probabilistic"
		}
	}

	// Default values for metrics
	if c.Metrics.Enabled {
		if c.ResourceAttributes == nil {
			c.ResourceAttributes = map[string]string{
				"service.name": c.ServiceName,
				"environment":  "production",
			}
		}
	}
}
