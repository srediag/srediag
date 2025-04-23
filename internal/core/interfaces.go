// Package core provides core interfaces and components for SREDIAG
package core

import (
	"context"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// Component represents a base component interface
type Component interface {
	// Start starts the component
	Start(ctx context.Context) error
	// Stop stops the component
	Stop(ctx context.Context) error
	// IsHealthy returns the health status of the component
	IsHealthy() bool
}

// Plugin represents a plugin interface
type Plugin interface {
	Component
	// GetName returns the plugin name
	GetName() string
	// GetVersion returns the plugin version
	GetVersion() string
	// GetType returns the plugin type
	GetType() string
	// GetCapabilities returns the capabilities of the plugin
	GetCapabilities() []Capability
	// GetStatus returns the status of the plugin
	GetStatus() Status
	// Configure configures the plugin with the given configuration
	Configure(cfg interface{}) error
}

// ResourceMonitor monitors system and application resources
type ResourceMonitor interface {
	Component
	// GetMetrics returns the current resource metrics
	GetMetrics() map[string]float64
	// SetThreshold sets a threshold for a metric
	SetThreshold(metric string, value float64) error
	// GetThresholds returns all configured thresholds
	GetThresholds() ResourceThresholds
}

// ConfigManager manages configuration loading and validation
type ConfigManager interface {
	Component
	// LoadConfig loads configuration from the given path
	LoadConfig(path string) error
	// SaveConfig saves configuration to the given path
	SaveConfig(path string) error
	// GetConfig returns the current configuration
	GetConfig() interface{}
	// ValidateConfig validates the given configuration
	ValidateConfig(cfg interface{}) error
}

// TelemetryBridge provides OpenTelemetry integration
type TelemetryBridge interface {
	Component
	// GetMeterProvider returns the OpenTelemetry meter provider
	GetMeterProvider() metric.MeterProvider
	// GetTracerProvider returns the OpenTelemetry tracer provider
	GetTracerProvider() trace.TracerProvider
	// GetLogger returns the configured logger
	GetLogger() *zap.Logger
}

// Runner represents a runnable component
type Runner interface {
	Component
	// GetLogger returns the configured logger
	GetLogger() *zap.Logger
}

// MetricsProvider provides metrics functionality
type MetricsProvider interface {
	// GetMeter returns a meter with the given name
	GetMeter(name string) metric.Meter
	// RegisterCallback registers a callback for metrics collection
	RegisterCallback(callback func(context.Context) error) error
}

// TracingProvider provides tracing functionality
type TracingProvider interface {
	// GetTracer returns a tracer with the given name
	GetTracer(name string) trace.Tracer
	// StartSpan starts a new span
	StartSpan(ctx context.Context, name string) (context.Context, trace.Span)
}

// LoggingProvider provides logging functionality
type LoggingProvider interface {
	// GetLogger returns a logger with the given name
	GetLogger(name string) *zap.Logger
	// WithFields returns a new logger with the given fields
	WithFields(fields ...zap.Field) *zap.Logger
}

// PluginProvider provides plugin management functionality
type PluginProvider interface {
	// Start starts the plugin provider
	Start(ctx context.Context) error
	// Stop stops the plugin provider
	Stop(ctx context.Context) error
	// IsHealthy returns the health status of the plugin provider
	IsHealthy() bool
	// GetPlugins returns a list of available plugins
	GetPlugins() ([]string, error)
	// GetPluginInfo returns information about a plugin
	GetPluginInfo(name string) (interface{}, error)
	// LoadPlugin loads a plugin
	LoadPlugin(name string) error
	// UnloadPlugin unloads a plugin
	UnloadPlugin(name string) error
}

// ISREDiagRunner defines the system runner interface
type ISREDiagRunner interface {
	Component
	// GetLogger returns the configured logger
	GetLogger() *zap.Logger
	// GetConfig returns the system configuration
	GetConfig() ISREDiagConfig
	// GetPluginManager returns the plugin manager
	GetPluginManager() PluginManager
	// GetTelemetryBridge returns the telemetry bridge
	GetTelemetryBridge() TelemetryBridge
	// GetDiagnosticManager returns the diagnostic manager
	GetDiagnosticManager() DiagnosticManager
}

// PluginManager manages plugin lifecycle
type PluginManager interface {
	Component
	// LoadPlugin loads a plugin from the given path
	LoadPlugin(path string) error
	// UnloadPlugin unloads a plugin by name
	UnloadPlugin(name string) error
	// GetPlugin returns a plugin by name
	GetPlugin(name string) (Plugin, error)
	// ListPlugins returns all loaded plugins
	ListPlugins() []Plugin
}
