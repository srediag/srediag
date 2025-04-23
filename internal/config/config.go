// Package config provides configuration management for SREDIAG
package config

import (
	"fmt"

	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"

	"github.com/srediag/srediag/internal/config/diagnostic"
	"github.com/srediag/srediag/internal/config/types"
	"github.com/srediag/srediag/internal/core"
)

// Root represents the main SREDIAG configuration
type Root struct {
	Core       types.CoreConfig      `mapstructure:"core"`
	Service    types.ServiceConfig   `mapstructure:"service"`
	Telemetry  types.TelemetryConfig `mapstructure:"telemetry"`
	Collector  types.CollectorConfig `mapstructure:"collector"`
	Diagnostic diagnostic.Config     `mapstructure:"diagnostic"`
}

// Load loads configuration from the specified file
func Load(filePath string) (*Root, error) {
	v := viper.New()
	v.SetConfigFile(filePath)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("config file not found: %s", filePath)
		}
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	cfg := Defaults()
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("error decoding configuration: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// Defaults returns the default configuration
func Defaults() *Root {
	return &Root{
		Core: types.CoreConfig{
			LogLevel:  "info",
			LogFormat: "json",
			Security: types.SecurityConfig{
				TLS: types.TLSConfig{
					Enabled:  false,
					CertFile: "/etc/srediag/certs/server.crt",
					KeyFile:  "/etc/srediag/certs/server.key",
				},
			},
		},
		Service: types.ServiceConfig{
			Name:        "srediag",
			Version:     "dev",
			Environment: "development",
		},
		Telemetry: types.TelemetryConfig{
			Metrics: types.MetricsConfig{
				Enabled:  true,
				Endpoint: "localhost:8888",
			},
			Traces: types.TracesConfig{
				Enabled:  true,
				Endpoint: "localhost:4317",
			},
			Resource: types.ResourceConfig{
				Attributes: map[string]string{
					"service.name":        "srediag",
					"service.version":     "dev",
					"service.environment": "development",
				},
			},
		},
		Collector: types.CollectorConfig{
			Enabled:    true,
			ConfigPath: "/etc/srediag/config/collector.yaml",
		},
		Diagnostic: diagnostic.DefaultConfig(),
	}
}

// Validate performs configuration validation
func (c *Root) Validate() error {
	if c.Service.Name == "" {
		return fmt.Errorf("service name is required")
	}
	if c.Service.Version == "" {
		return fmt.Errorf("service version is required")
	}
	if c.Service.Environment == "" {
		return fmt.Errorf("service environment is required")
	}
	return nil
}

// CreateResource creates an OpenTelemetry resource from the configuration
func (c *Root) CreateResource() (*resource.Resource, error) {
	attrs := []attribute.KeyValue{
		semconv.ServiceName(c.Service.Name),
		semconv.ServiceVersion(c.Service.Version),
		semconv.DeploymentEnvironment(c.Service.Environment),
	}

	for k, v := range c.Telemetry.Resource.Attributes {
		attrs = append(attrs, attribute.String(k, v))
	}

	return resource.Merge(
		resource.Default(),
		resource.NewWithAttributes("", attrs...),
	)
}

// GetVersion implements ISREDiagConfig
func (c *Root) GetVersion() string {
	return c.Service.Version
}

// GetServiceConfig implements ISREDiagConfig
func (c *Root) GetServiceConfig() core.ISREDiagServiceConfig {
	return &c.Service
}

// GetCollectorConfig implements ISREDiagConfig
func (c *Root) GetCollectorConfig() core.ISREDiagCollectorConfig {
	return &c.Collector
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Root {
	return &Root{
		Core: types.CoreConfig{
			LogLevel:  "info",
			LogFormat: "json",
			Security: types.SecurityConfig{
				TLS: types.TLSConfig{
					Enabled:  false,
					CertFile: "/etc/srediag/certs/server.crt",
					KeyFile:  "/etc/srediag/certs/server.key",
				},
			},
		},
		Service: types.ServiceConfig{
			Name:        "srediag",
			Version:     "dev",
			Environment: "development",
			Type:        core.TypeDiagnostic,
		},
		Telemetry: types.TelemetryConfig{
			Metrics: types.MetricsConfig{
				Enabled:  true,
				Endpoint: "localhost:8888",
			},
			Traces: types.TracesConfig{
				Enabled:  true,
				Endpoint: "localhost:4317",
			},
			Resource: types.ResourceConfig{
				Attributes: map[string]string{
					"service.name":        "srediag",
					"service.version":     "dev",
					"service.environment": "development",
				},
			},
		},
		Collector: types.CollectorConfig{
			Enabled:    true,
			ConfigPath: "/etc/srediag/config/collector.yaml",
		},
		Diagnostic: diagnostic.DefaultConfig(),
	}
}
