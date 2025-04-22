package config

import (
	"fmt"

	"github.com/srediag/srediag/internal/core"
)

// Config represents the main SREDIAG configuration
type Config struct {
	Core       CoreConfig       `mapstructure:"core"`
	Service    ServiceConfig    `mapstructure:"service"`
	Telemetry  TelemetryConfig  `mapstructure:"telemetry"`
	Plugins    PluginsConfig    `mapstructure:"plugins"`
	Collector  CollectorConfig  `mapstructure:"collector"`
	Diagnostic DiagnosticConfig `mapstructure:"diagnostic"`
}

// Ensure Config implements ISREDiagConfig
var _ core.ISREDiagConfig = (*Config)(nil)

// CoreConfig represents the core configuration
type CoreConfig struct {
	LogLevel  string         `mapstructure:"log_level"`
	LogFormat string         `mapstructure:"log_format"`
	Security  SecurityConfig `mapstructure:"security"`
}

// ServiceConfig represents the service configuration
type ServiceConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"`
}

// Ensure ServiceConfig implements ISREDiagServiceConfig
var _ core.ISREDiagServiceConfig = (*ServiceConfig)(nil)

// TelemetryConfig represents the telemetry configuration
type TelemetryConfig struct {
	Metrics  MetricsConfig  `mapstructure:"metrics"`
	Traces   TracesConfig   `mapstructure:"traces"`
	Resource ResourceConfig `mapstructure:"resource"`
}

// PluginsConfig represents the plugins configuration
type PluginsConfig struct {
	Directory string         `mapstructure:"directory"`
	AutoLoad  bool           `mapstructure:"autoload"`
	Enabled   []string       `mapstructure:"enabled"`
	Settings  map[string]any `mapstructure:"settings"`
}

// CollectorConfig represents the OpenTelemetry collector configuration
type CollectorConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	ConfigPath string `mapstructure:"config_path"`
}

// Ensure CollectorConfig implements ISREDiagCollectorConfig
var _ core.ISREDiagCollectorConfig = (*CollectorConfig)(nil)

// DiagnosticConfig represents the diagnostic configuration
type DiagnosticConfig struct {
	System     SystemDiagConfig     `mapstructure:"system"`
	Kubernetes KubernetesDiagConfig `mapstructure:"kubernetes"`
	Cloud      CloudDiagConfig      `mapstructure:"cloud"`
	Security   SecurityDiagConfig   `mapstructure:"security"`
}

// SystemDiagConfig represents system diagnostic configuration
type SystemDiagConfig struct {
	Enabled     bool    `mapstructure:"enabled"`
	Interval    string  `mapstructure:"interval"`
	CPULimit    float64 `mapstructure:"cpu_limit"`
	MemoryLimit float64 `mapstructure:"memory_limit"`
	DiskLimit   float64 `mapstructure:"disk_limit"`
}

// KubernetesDiagConfig represents Kubernetes diagnostic configuration
type KubernetesDiagConfig struct {
	Enabled   bool     `mapstructure:"enabled"`
	Clusters  []string `mapstructure:"clusters"`
	Namespace string   `mapstructure:"namespace"`
}

// CloudDiagConfig represents cloud provider diagnostic configuration
type CloudDiagConfig struct {
	Enabled     bool              `mapstructure:"enabled"`
	Providers   []string          `mapstructure:"providers"`
	Credentials map[string]string `mapstructure:"credentials"`
}

// SecurityDiagConfig represents security diagnostic configuration
type SecurityDiagConfig struct {
	Enabled         bool     `mapstructure:"enabled"`
	ScanInterval    string   `mapstructure:"scan_interval"`
	Standards       []string `mapstructure:"standards"`
	ComplianceLevel string   `mapstructure:"compliance_level"`
}

// MetricsConfig represents the metrics configuration
type MetricsConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Endpoint string `mapstructure:"endpoint"`
}

// TracesConfig represents the traces configuration
type TracesConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Endpoint string `mapstructure:"endpoint"`
}

// ResourceConfig represents the resource configuration
type ResourceConfig struct {
	Attributes map[string]string `mapstructure:"attributes"`
}

// SecurityConfig represents the security configuration
type SecurityConfig struct {
	TLS TLSConfig `mapstructure:"tls"`
}

// TLSConfig represents the TLS configuration
type TLSConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	CertFile string `mapstructure:"cert_file"`
	KeyFile  string `mapstructure:"key_file"`
	CAFile   string `mapstructure:"ca_file,omitempty"`
}

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

// GetName implements ISREDiagServiceConfig
func (s *ServiceConfig) GetName() string {
	return s.Name
}

// GetEnvironment implements ISREDiagServiceConfig
func (s *ServiceConfig) GetEnvironment() string {
	return s.Environment
}

// IsEnabled implements ISREDiagCollectorConfig
func (c *CollectorConfig) IsEnabled() bool {
	return c.Enabled
}

// GetConfigPath implements ISREDiagCollectorConfig
func (c *CollectorConfig) GetConfigPath() string {
	return c.ConfigPath
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Core: CoreConfig{
			LogLevel:  "info",
			LogFormat: "json",
			Security: SecurityConfig{
				TLS: TLSConfig{
					Enabled:  false,
					CertFile: "/etc/srediag/certs/server.crt",
					KeyFile:  "/etc/srediag/certs/server.key",
				},
			},
		},
		Service: ServiceConfig{
			Name:        "srediag",
			Version:     "dev",
			Environment: "development",
		},
		Telemetry: TelemetryConfig{
			Metrics: MetricsConfig{
				Enabled:  true,
				Endpoint: "localhost:8888",
			},
			Traces: TracesConfig{
				Enabled:  true,
				Endpoint: "localhost:4317",
			},
			Resource: ResourceConfig{
				Attributes: map[string]string{
					"service.name":        "srediag",
					"service.version":     "dev",
					"service.environment": "development",
				},
			},
		},
		Plugins: PluginsConfig{
			Directory: "/etc/srediag/plugins",
			AutoLoad:  true,
			Enabled:   []string{},
			Settings:  make(map[string]any),
		},
		Collector: CollectorConfig{
			Enabled:    true,
			ConfigPath: "/etc/srediag/config/collector.yaml",
		},
		Diagnostic: DiagnosticConfig{
			System: SystemDiagConfig{
				Enabled:     true,
				Interval:    "30s",
				CPULimit:    80,
				MemoryLimit: 90,
				DiskLimit:   85,
			},
			Kubernetes: KubernetesDiagConfig{
				Enabled:   false,
				Clusters:  []string{},
				Namespace: "default",
			},
			Cloud: CloudDiagConfig{
				Enabled:     false,
				Providers:   []string{},
				Credentials: make(map[string]string),
			},
			Security: SecurityDiagConfig{
				Enabled:         true,
				ScanInterval:    "1h",
				Standards:       []string{"pci-dss", "hipaa"},
				ComplianceLevel: "high",
			},
		},
	}
}
