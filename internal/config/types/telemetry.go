package types

// MetricsConfig represents the metrics configuration
type MetricsConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Endpoint string `mapstructure:"endpoint"`
}

// ResourceConfig represents the resource configuration
type ResourceConfig struct {
	Attributes map[string]string `mapstructure:"attributes"`
}

// TelemetryConfig represents the telemetry configuration
type TelemetryConfig struct {
	Metrics  MetricsConfig  `mapstructure:"metrics"`
	Traces   TracesConfig   `mapstructure:"traces"`
	Resource ResourceConfig `mapstructure:"resource"`
}

// TracesConfig represents the traces configuration
type TracesConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Endpoint string `mapstructure:"endpoint"`
}
