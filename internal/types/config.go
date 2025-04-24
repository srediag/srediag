package types

import (
	"fmt"
)

// ConfigFormat represents the configuration file format
type ConfigFormat string

const (
	// ConfigFormatYAML represents YAML configuration format
	ConfigFormatYAML ConfigFormat = "yaml"
	// ConfigFormatJSON represents JSON configuration format
	ConfigFormatJSON ConfigFormat = "json"
	// ConfigFormatTOML represents TOML configuration format
	ConfigFormatTOML ConfigFormat = "toml"
)

// ConfigLogLevel represents the logging level
type ConfigLogLevel string

const (
	// ConfigLogLevelDebug represents debug log level
	ConfigLogLevelDebug ConfigLogLevel = "debug"
	// ConfigLogLevelInfo represents info log level
	ConfigLogLevelInfo ConfigLogLevel = "info"
	// ConfigLogLevelWarn represents warn log level
	ConfigLogLevelWarn ConfigLogLevel = "warn"
	// ConfigLogLevelError represents error log level
	ConfigLogLevelError ConfigLogLevel = "error"
)

// ConfigLogFormat represents the logging format
type ConfigLogFormat string

const (
	// ConfigLogFormatJSON represents JSON log format
	ConfigLogFormatJSON ConfigLogFormat = "json"
	// ConfigLogFormatText represents text log format
	ConfigLogFormatText ConfigLogFormat = "text"
)

// Configuration Interfaces

// IConfig represents the main configuration interface that provides access to all configuration components
type IConfig interface {
	// GetVersion returns the version of the configuration
	GetVersion() string
	// GetCore returns the core configuration
	GetCore() ICoreConfig
	// GetService returns the service configuration
	GetService() IServiceConfig
	// GetCollector returns the collector configuration
	GetCollector() ICollectorConfig
	// GetTelemetry returns the telemetry configuration
	GetTelemetry() ITelemetryConfig
	// GetPlugins returns the plugins configuration
	GetPlugins() IPluginsConfig
	// GetDiagnostic returns the diagnostic configuration
	GetDiagnostic() IDiagnosticConfig
	// Validate validates the entire configuration
	Validate() error
}

// ICoreConfig represents the core configuration interface
type ICoreConfig interface {
	// GetLogLevel returns the logging level
	GetLogLevel() ConfigLogLevel
	// GetLogFormat returns the logging format
	GetLogFormat() ConfigLogFormat
	// GetVersion returns the core version
	GetVersion() string
}

// IServiceConfig represents the service configuration interface
type IServiceConfig interface {
	// GetName returns the name of the service
	GetName() string
	// GetEnvironment returns the environment of the service
	GetEnvironment() string
	// GetType returns the type of the service
	GetType() ComponentType
	// GetSecurity returns the security configuration
	GetSecurity() SecurityConfig
	// GetVersion returns the service version
	GetVersion() string
}

// ITelemetryConfig represents the telemetry configuration interface
type ITelemetryConfig interface {
	// IsEnabled returns true if telemetry is enabled
	IsEnabled() bool
	// GetMetrics returns the metrics configuration
	GetMetrics() IMetricsConfig
	// GetTraces returns the traces configuration
	GetTraces() ITracesConfig
	// GetResource returns the resource configuration
	GetResource() IResourceConfig
}

// IMetricsConfig represents the metrics configuration interface
type IMetricsConfig interface {
	// IsEnabled returns true if metrics collection is enabled
	IsEnabled() bool
	// GetEndpoint returns the metrics endpoint
	GetEndpoint() string
	// GetPort returns the metrics port
	GetPort() int
	// GetAttributes returns the metrics attributes
	GetAttributes() map[string]string
}

// ITracesConfig represents the traces configuration interface
type ITracesConfig interface {
	// IsEnabled returns true if trace collection is enabled
	IsEnabled() bool
	// GetEndpoint returns the traces endpoint
	GetEndpoint() string
	// GetPort returns the traces port
	GetPort() int
	// GetAttributes returns the traces attributes
	GetAttributes() map[string]string
}

// IResourceConfig represents the resource configuration interface
type IResourceConfig interface {
	// GetAttributes returns the resource attributes
	GetAttributes() map[string]string
}

// IPluginsConfig represents the plugins configuration interface
type IPluginsConfig interface {
	// GetDirectory returns the plugins directory
	GetDirectory() string
	// IsAutoLoadEnabled returns true if auto-loading is enabled
	IsAutoLoadEnabled() bool
	// GetSettings returns the plugin settings
	GetSettings() map[string]PluginConfig
}

// ICollectorConfig represents the collector configuration interface
type ICollectorConfig interface {
	// IsEnabled returns true if the collector is enabled
	IsEnabled() bool
	// GetConfigPath returns the path to the collector configuration
	GetConfigPath() string
	// GetPipelines returns the collector pipelines
	GetPipelines() []PipelineConfig
}

// IDiagnosticConfig represents the diagnostic configuration interface
type IDiagnosticConfig interface {
	// GetSystem returns the system diagnostic configuration
	GetSystem() ISystemConfig
	// GetKubernetes returns the kubernetes diagnostic configuration
	GetKubernetes() IKubernetesConfig
	// GetCloud returns the cloud diagnostic configuration
	GetCloud() ICloudConfig
}

// ISystemConfig represents the system diagnostic configuration interface
type ISystemConfig interface {
	// IsEnabled returns true if system diagnostics are enabled
	IsEnabled() bool
	// GetInterval returns the collection interval
	GetInterval() string
	// GetCPULimit returns the CPU limit
	GetCPULimit() float64
	// GetMemoryLimit returns the memory limit
	GetMemoryLimit() float64
	// GetDiskLimit returns the disk limit
	GetDiskLimit() float64
}

// IKubernetesConfig represents the kubernetes diagnostic configuration interface
type IKubernetesConfig interface {
	// IsEnabled returns true if kubernetes diagnostics are enabled
	IsEnabled() bool
	// GetClusters returns the kubernetes clusters
	GetClusters() []string
	// GetNamespace returns the kubernetes namespace
	GetNamespace() string
}

// ICloudConfig represents the cloud diagnostic configuration interface
type ICloudConfig interface {
	// IsEnabled returns true if cloud diagnostics are enabled
	IsEnabled() bool
	// GetProviders returns the cloud providers
	GetProviders() []string
	// GetCredentials returns the cloud credentials
	GetCredentials() map[string]string
}

// IConfigManager represents the configuration manager interface
type IConfigManager interface {
	IComponent
	// LoadConfig loads configuration from the specified path
	LoadConfig(path string) error
	// SaveConfig saves configuration to the specified path
	SaveConfig(path string) error
	// GetConfig returns the current configuration
	GetConfig() IConfig
	// SetConfig sets the current configuration
	SetConfig(cfg IConfig) error
	// ValidateConfig validates the configuration
	ValidateConfig(cfg IConfig) error
}

// Configuration Implementations

// CoreConfig represents the core service configuration
type CoreConfig struct {
	LogLevel  ConfigLogLevel  `mapstructure:"log_level" json:"log_level"`
	LogFormat ConfigLogFormat `mapstructure:"log_format" json:"log_format"`
	Version   string          `mapstructure:"version" json:"version"`
}

// SecurityConfig represents security-related configuration
type SecurityConfig struct {
	TLS      TLSConfig `mapstructure:"tls" json:"tls"`
	Enabled  bool      `mapstructure:"enabled" json:"enabled"`
	CertFile string    `mapstructure:"cert_file" json:"cert_file"`
	KeyFile  string    `mapstructure:"key_file" json:"key_file"`
}

// TLSConfig represents TLS configuration
type TLSConfig struct {
	Enabled    bool   `mapstructure:"enabled" json:"enabled"`
	CertFile   string `mapstructure:"cert_file" json:"cert_file"`
	KeyFile    string `mapstructure:"key_file" json:"key_file"`
	CAFile     string `mapstructure:"ca_file" json:"ca_file"`
	ServerName string `mapstructure:"server_name" json:"server_name"`
	SkipVerify bool   `mapstructure:"skip_verify" json:"skip_verify"`
	MinVersion string `mapstructure:"min_version" json:"min_version"`
	MaxVersion string `mapstructure:"max_version" json:"max_version"`
}

// ServiceConfig represents service-specific configuration
type ServiceConfig struct {
	Name        string         `mapstructure:"name" json:"name"`
	Version     string         `mapstructure:"version" json:"version"`
	Environment string         `mapstructure:"environment" json:"environment"`
	Type        ComponentType  `mapstructure:"type" json:"type"`
	Security    SecurityConfig `mapstructure:"security" json:"security"`
}

// TelemetryConfig represents telemetry configuration
type TelemetryConfig struct {
	Metrics  MetricsConfig  `mapstructure:"metrics" json:"metrics"`
	Traces   TracesConfig   `mapstructure:"traces" json:"traces"`
	Resource ResourceConfig `mapstructure:"resource" json:"resource"`
	Enabled  bool           `mapstructure:"enabled" json:"enabled"`
}

// MetricsConfig represents metrics configuration
type MetricsConfig struct {
	Enabled    bool              `mapstructure:"enabled" json:"enabled"`
	Endpoint   string            `mapstructure:"endpoint" json:"endpoint"`
	Port       int               `mapstructure:"port" json:"port"`
	Attributes map[string]string `mapstructure:"attributes" json:"attributes"`
}

// TracesConfig represents tracing configuration
type TracesConfig struct {
	Enabled    bool              `mapstructure:"enabled" json:"enabled"`
	Endpoint   string            `mapstructure:"endpoint" json:"endpoint"`
	Port       int               `mapstructure:"port" json:"port"`
	Attributes map[string]string `mapstructure:"attributes" json:"attributes"`
}

// ResourceConfig represents OpenTelemetry resource configuration
type ResourceConfig struct {
	Attributes map[string]string `mapstructure:"attributes" json:"attributes"`
}

// CollectorConfig represents collector configuration
type CollectorConfig struct {
	Enabled    bool             `mapstructure:"enabled" json:"enabled"`
	ConfigPath string           `mapstructure:"config_path" json:"config_path"`
	Pipelines  []PipelineConfig `mapstructure:"pipelines" json:"pipelines"`
}

// PipelineConfig represents a collector pipeline configuration
type PipelineConfig struct {
	Name       string   `mapstructure:"name" json:"name"`
	Type       string   `mapstructure:"type" json:"type"`
	Receivers  []string `mapstructure:"receivers" json:"receivers"`
	Processors []string `mapstructure:"processors" json:"processors"`
	Exporters  []string `mapstructure:"exporters" json:"exporters"`
}

// SystemConfig represents system diagnostic configuration
type SystemConfig struct {
	Enabled     bool    `mapstructure:"enabled" json:"enabled"`
	Interval    string  `mapstructure:"interval" json:"interval"`
	CPULimit    float64 `mapstructure:"cpu_limit" json:"cpu_limit"`
	MemoryLimit float64 `mapstructure:"memory_limit" json:"memory_limit"`
	DiskLimit   float64 `mapstructure:"disk_limit" json:"disk_limit"`
}

// KubernetesConfig represents Kubernetes diagnostic configuration
type KubernetesConfig struct {
	Enabled   bool     `mapstructure:"enabled" json:"enabled"`
	Clusters  []string `mapstructure:"clusters" json:"clusters"`
	Namespace string   `mapstructure:"namespace" json:"namespace"`
}

// CloudConfig represents cloud diagnostic configuration
type CloudConfig struct {
	Enabled     bool              `mapstructure:"enabled" json:"enabled"`
	Providers   []string          `mapstructure:"providers" json:"providers"`
	Credentials map[string]string `mapstructure:"credentials" json:"credentials"`
}

// DiagnosticConfig represents the complete diagnostic configuration
type DiagnosticConfig struct {
	System     SystemConfig     `mapstructure:"system" json:"system"`
	Kubernetes KubernetesConfig `mapstructure:"kubernetes" json:"kubernetes"`
	Cloud      CloudConfig      `mapstructure:"cloud" json:"cloud"`
}

// PluginsConfig represents plugins configuration
type PluginsConfig struct {
	Directory string                  `mapstructure:"directory" json:"directory"`
	AutoLoad  bool                    `mapstructure:"auto_load" json:"auto_load"`
	Settings  map[string]PluginConfig `mapstructure:"settings" json:"settings"`
}

// Config represents the complete configuration structure
type Config struct {
	Core       CoreConfig       `mapstructure:"core" json:"core"`
	Service    ServiceConfig    `mapstructure:"service" json:"service"`
	Telemetry  TelemetryConfig  `mapstructure:"telemetry" json:"telemetry"`
	Collector  CollectorConfig  `mapstructure:"collector" json:"collector"`
	Plugins    PluginsConfig    `mapstructure:"plugins" json:"plugins"`
	Diagnostic DiagnosticConfig `mapstructure:"diagnostic" json:"diagnostic"`
}

// Ensure implementations satisfy interfaces
var (
	_ IConfig           = (*Config)(nil)
	_ ICoreConfig       = (*CoreConfig)(nil)
	_ IServiceConfig    = (*ServiceConfig)(nil)
	_ ITelemetryConfig  = (*TelemetryConfig)(nil)
	_ IMetricsConfig    = (*MetricsConfig)(nil)
	_ ITracesConfig     = (*TracesConfig)(nil)
	_ IResourceConfig   = (*ResourceConfig)(nil)
	_ ICollectorConfig  = (*CollectorConfig)(nil)
	_ IPluginsConfig    = (*PluginsConfig)(nil)
	_ IDiagnosticConfig = (*DiagnosticConfig)(nil)
	_ ISystemConfig     = (*SystemConfig)(nil)
	_ IKubernetesConfig = (*KubernetesConfig)(nil)
	_ ICloudConfig      = (*CloudConfig)(nil)
)

// Config implementation
func (c *Config) GetVersion() string {
	return c.Service.Version
}

func (c *Config) GetCore() ICoreConfig {
	return &c.Core
}

func (c *Config) GetService() IServiceConfig {
	return &c.Service
}

func (c *Config) GetCollector() ICollectorConfig {
	return &c.Collector
}

func (c *Config) GetTelemetry() ITelemetryConfig {
	return &c.Telemetry
}

func (c *Config) GetPlugins() IPluginsConfig {
	return &c.Plugins
}

func (c *Config) GetDiagnostic() IDiagnosticConfig {
	return &c.Diagnostic
}

func (c *Config) Validate() error {
	if c.Service.Name == "" {
		return fmt.Errorf("service name is required")
	}
	if c.Service.Version == "" {
		return fmt.Errorf("service version is required")
	}
	if c.Core.LogLevel == "" {
		return fmt.Errorf("core log level is required")
	}
	if c.Core.LogFormat == "" {
		return fmt.Errorf("core log format is required")
	}
	return nil
}

// CoreConfig implementation
func (c *CoreConfig) GetLogLevel() ConfigLogLevel {
	return c.LogLevel
}

func (c *CoreConfig) GetLogFormat() ConfigLogFormat {
	return c.LogFormat
}

func (c *CoreConfig) GetVersion() string {
	return c.Version
}

// ServiceConfig implementation
func (c *ServiceConfig) GetName() string {
	return c.Name
}

func (c *ServiceConfig) GetEnvironment() string {
	return c.Environment
}

func (c *ServiceConfig) GetType() ComponentType {
	return c.Type
}

func (c *ServiceConfig) GetSecurity() SecurityConfig {
	return c.Security
}

func (c *ServiceConfig) GetVersion() string {
	return c.Version
}

// TelemetryConfig implementation
func (c *TelemetryConfig) IsEnabled() bool {
	return c.Enabled
}

func (c *TelemetryConfig) GetMetrics() IMetricsConfig {
	return &c.Metrics
}

func (c *TelemetryConfig) GetTraces() ITracesConfig {
	return &c.Traces
}

func (c *TelemetryConfig) GetResource() IResourceConfig {
	return &c.Resource
}

// MetricsConfig implementation
func (c *MetricsConfig) IsEnabled() bool {
	return c.Enabled
}

func (c *MetricsConfig) GetEndpoint() string {
	return c.Endpoint
}

func (c *MetricsConfig) GetPort() int {
	return c.Port
}

func (c *MetricsConfig) GetAttributes() map[string]string {
	return c.Attributes
}

// TracesConfig implementation
func (c *TracesConfig) IsEnabled() bool {
	return c.Enabled
}

func (c *TracesConfig) GetEndpoint() string {
	return c.Endpoint
}

func (c *TracesConfig) GetPort() int {
	return c.Port
}

func (c *TracesConfig) GetAttributes() map[string]string {
	return c.Attributes
}

// ResourceConfig implementation
func (c *ResourceConfig) GetAttributes() map[string]string {
	return c.Attributes
}

// CollectorConfig implementation
func (c *CollectorConfig) IsEnabled() bool {
	return c.Enabled
}

func (c *CollectorConfig) GetConfigPath() string {
	return c.ConfigPath
}

func (c *CollectorConfig) GetPipelines() []PipelineConfig {
	return c.Pipelines
}

// PluginsConfig implementation
func (c *PluginsConfig) GetDirectory() string {
	return c.Directory
}

func (c *PluginsConfig) IsAutoLoadEnabled() bool {
	return c.AutoLoad
}

func (c *PluginsConfig) GetSettings() map[string]PluginConfig {
	return c.Settings
}

// DiagnosticConfig implementation
func (c *DiagnosticConfig) GetSystem() ISystemConfig {
	return &c.System
}

func (c *DiagnosticConfig) GetKubernetes() IKubernetesConfig {
	return &c.Kubernetes
}

func (c *DiagnosticConfig) GetCloud() ICloudConfig {
	return &c.Cloud
}

// SystemConfig implementation
func (c *SystemConfig) IsEnabled() bool {
	return c.Enabled
}

func (c *SystemConfig) GetInterval() string {
	return c.Interval
}

func (c *SystemConfig) GetCPULimit() float64 {
	return c.CPULimit
}

func (c *SystemConfig) GetMemoryLimit() float64 {
	return c.MemoryLimit
}

func (c *SystemConfig) GetDiskLimit() float64 {
	return c.DiskLimit
}

// KubernetesConfig implementation
func (c *KubernetesConfig) IsEnabled() bool {
	return c.Enabled
}

func (c *KubernetesConfig) GetClusters() []string {
	return c.Clusters
}

func (c *KubernetesConfig) GetNamespace() string {
	return c.Namespace
}

// CloudConfig implementation
func (c *CloudConfig) IsEnabled() bool {
	return c.Enabled
}

func (c *CloudConfig) GetProviders() []string {
	return c.Providers
}

func (c *CloudConfig) GetCredentials() map[string]string {
	return c.Credentials
}
