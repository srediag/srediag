// Package config provides configuration management for SREDIAG
package config

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"

	"github.com/srediag/srediag/internal/types"
)

const (
	// DefaultConfigPath is the default path for the main configuration file
	DefaultConfigPath = "/etc/srediag/config/config.yaml"
	// DefaultPluginsPath is the default path for plugins
	DefaultPluginsPath = "/etc/srediag/plugins"
	// DefaultCertsPath is the default path for certificates
	DefaultCertsPath = "/etc/srediag/certs"
)

// ConfigRoot represents the main SREDIAG configuration
type ConfigRoot struct {
	// Core configuration
	Core types.CoreConfig `mapstructure:"core" json:"core"`
	// Service configuration
	Service types.ServiceConfig `mapstructure:"service" json:"service"`
	// Telemetry configuration
	Telemetry types.TelemetryConfig `mapstructure:"telemetry" json:"telemetry"`
	// Collector configuration
	Collector types.CollectorConfig `mapstructure:"collector" json:"collector"`
	// Plugins configuration
	Plugins types.PluginsConfig `mapstructure:"plugins" json:"plugins"`
	// Diagnostic configuration
	Diagnostic types.DiagnosticConfig `mapstructure:"diagnostic" json:"diagnostic"`
}

// Ensure ConfigRoot implements types.IConfig
var _ types.IConfig = (*ConfigRoot)(nil)

// LoadOptions represents options for loading configuration
type LoadOptions struct {
	// ConfigPath is the path to the configuration file
	ConfigPath string
	// EnvPrefix is the prefix for environment variables
	EnvPrefix string
	// AllowEnvOverride allows environment variables to override config file values
	AllowEnvOverride bool
}

// Load loads configuration from the specified file with options
func Load(opts LoadOptions) (*ConfigRoot, error) {
	if opts.ConfigPath == "" {
		opts.ConfigPath = DefaultConfigPath
	}

	v := viper.New()
	v.SetConfigFile(opts.ConfigPath)

	// Set up environment variables support
	if opts.AllowEnvOverride {
		v.SetEnvPrefix(opts.EnvPrefix)
		v.AutomaticEnv()
	}

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("config file not found: %s", opts.ConfigPath)
		}
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	cfg := DefaultConfig()
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("error decoding configuration: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// DefaultConfig returns the default configuration with secure defaults
func DefaultConfig() *ConfigRoot {
	return &ConfigRoot{
		Core: types.CoreConfig{
			LogLevel:  types.ConfigLogLevelInfo,
			LogFormat: types.ConfigLogFormatJSON,
			Version:   "1.0.0",
		},
		Service: types.ServiceConfig{
			Name:        "srediag",
			Version:     "dev",
			Environment: "development",
			Type:        types.ComponentTypeService,
			Security: types.SecurityConfig{
				TLS: types.TLSConfig{
					Enabled:    false,
					CertFile:   filepath.Join(DefaultCertsPath, "server.crt"),
					KeyFile:    filepath.Join(DefaultCertsPath, "server.key"),
					CAFile:     filepath.Join(DefaultCertsPath, "ca.crt"),
					MinVersion: "TLS1.2",
					SkipVerify: false,
				},
			},
		},
		Telemetry: types.TelemetryConfig{
			Enabled: true,
			Metrics: types.MetricsConfig{
				Enabled:  true,
				Endpoint: "localhost:8888",
				Port:     8888,
				Attributes: map[string]string{
					"deployment.environment": "development",
				},
			},
			Traces: types.TracesConfig{
				Enabled:  true,
				Endpoint: "localhost:4317",
				Port:     4317,
				Attributes: map[string]string{
					"deployment.environment": "development",
				},
			},
			Resource: types.ResourceConfig{
				Attributes: map[string]string{
					"service.name":        "srediag",
					"service.version":     "dev",
					"service.environment": "development",
					"service.namespace":   "monitoring",
				},
			},
		},
		Collector: types.CollectorConfig{
			Enabled:    true,
			ConfigPath: filepath.Join(filepath.Dir(DefaultConfigPath), "collector.yaml"),
			Pipelines: []types.PipelineConfig{
				{
					Name:       "metrics",
					Type:       "metrics",
					Receivers:  []string{"otlp"},
					Processors: []string{"batch"},
					Exporters:  []string{"prometheus"},
				},
				{
					Name:       "traces",
					Type:       "traces",
					Receivers:  []string{"otlp"},
					Processors: []string{"batch"},
					Exporters:  []string{"jaeger"},
				},
			},
		},
		Plugins: types.PluginsConfig{
			Directory: DefaultPluginsPath,
			AutoLoad:  true,
			Settings:  make(map[string]types.PluginConfig),
		},
		Diagnostic: types.DiagnosticConfig{
			System: types.SystemConfig{
				Enabled:     true,
				Interval:    "30s",
				CPULimit:    80.0,
				MemoryLimit: 90.0,
				DiskLimit:   85.0,
			},
			Kubernetes: types.KubernetesConfig{
				Enabled:   false,
				Clusters:  []string{},
				Namespace: "default",
			},
			Cloud: types.CloudConfig{
				Enabled:     false,
				Providers:   []string{},
				Credentials: make(map[string]string),
			},
		},
	}
}

// GetVersion implements types.IConfig
func (c *ConfigRoot) GetVersion() string {
	return c.Service.Version
}

// GetCore implements types.IConfig
func (c *ConfigRoot) GetCore() types.ICoreConfig {
	return &c.Core
}

// GetService implements types.IConfig
func (c *ConfigRoot) GetService() types.IServiceConfig {
	return &c.Service
}

// GetCollector implements types.IConfig
func (c *ConfigRoot) GetCollector() types.ICollectorConfig {
	return &c.Collector
}

// GetTelemetry implements types.IConfig
func (c *ConfigRoot) GetTelemetry() types.ITelemetryConfig {
	return &c.Telemetry
}

// GetPlugins implements types.IConfig
func (c *ConfigRoot) GetPlugins() types.IPluginsConfig {
	return &c.Plugins
}

// GetDiagnostic implements types.IConfig
func (c *ConfigRoot) GetDiagnostic() types.IDiagnosticConfig {
	return &c.Diagnostic
}

// Validate implements types.IConfig
func (c *ConfigRoot) Validate() error {
	if c.Service.Name == "" {
		return fmt.Errorf("service name is required")
	}
	if c.Service.Version == "" {
		return fmt.Errorf("service version is required")
	}
	if c.Service.Environment == "" {
		return fmt.Errorf("service environment is required")
	}

	// Validate TLS configuration if enabled
	if c.Service.Security.TLS.Enabled {
		if c.Service.Security.TLS.CertFile == "" {
			return fmt.Errorf("TLS certificate file is required when TLS is enabled")
		}
		if c.Service.Security.TLS.KeyFile == "" {
			return fmt.Errorf("TLS key file is required when TLS is enabled")
		}
	}

	// Validate telemetry configuration if enabled
	if c.Telemetry.IsEnabled() {
		if c.Telemetry.Metrics.IsEnabled() && c.Telemetry.Metrics.Port == 0 {
			return fmt.Errorf("metrics port is required when metrics are enabled")
		}
		if c.Telemetry.Traces.IsEnabled() && c.Telemetry.Traces.Port == 0 {
			return fmt.Errorf("traces port is required when traces are enabled")
		}
	}

	// Validate collector configuration if enabled
	if c.Collector.IsEnabled() && c.Collector.ConfigPath == "" {
		return fmt.Errorf("collector config path is required when collector is enabled")
	}

	return nil
}

// CreateResource creates an OpenTelemetry resource from the configuration
func (c *ConfigRoot) CreateResource() (*resource.Resource, error) {
	attrs := []attribute.KeyValue{
		semconv.ServiceName(c.Service.Name),
		semconv.ServiceVersion(c.Service.Version),
		semconv.DeploymentEnvironment(c.Service.Environment),
	}

	// Add all configured resource attributes
	for k, v := range c.Telemetry.Resource.Attributes {
		attrs = append(attrs, attribute.String(k, v))
	}

	// Add service namespace if available
	if ns, ok := c.Telemetry.Resource.Attributes["service.namespace"]; ok {
		attrs = append(attrs, semconv.ServiceNamespace(ns))
	}

	return resource.Merge(
		resource.Default(),
		resource.NewWithAttributes("", attrs...),
	)
}
