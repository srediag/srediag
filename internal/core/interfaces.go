package core

import (
	"context"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// Component represents a core component of the application
type Component interface {
	// Start starts the component
	Start(ctx context.Context) error
	// Stop stops the component
	Stop(ctx context.Context) error
	// IsHealthy returns the health status of the component
	IsHealthy() bool
}

// Plugin represents a SREDIAG plugin
type Plugin interface {
	Component
	// Type returns the plugin type
	Type() PluginType
	// Name returns the plugin name
	Name() string
	// Version returns the plugin version
	Version() string
	// Configure configures the plugin with the given configuration
	Configure(cfg interface{}) error
}

// PluginType represents the type of a plugin
type PluginType string

const (
	// PluginTypeDiagnostic represents a diagnostic plugin
	PluginTypeDiagnostic PluginType = "diagnostic"
	// PluginTypeAnalysis represents an analysis plugin
	PluginTypeAnalysis PluginType = "analysis"
	// PluginTypeManagement represents a management plugin
	PluginTypeManagement PluginType = "management"
)

// PluginManager manages the lifecycle of plugins
type PluginManager interface {
	Component
	// LoadPlugin loads a plugin from the given path
	LoadPlugin(path string) (Plugin, error)
	// UnloadPlugin unloads a plugin by name
	UnloadPlugin(name string) error
	// GetPlugin returns a plugin by name
	GetPlugin(name string) (Plugin, error)
	// ListPlugins returns all loaded plugins
	ListPlugins() []Plugin
}

// ResourceMonitor monitors system resources
type ResourceMonitor interface {
	Component
	// CollectMetrics collects system metrics
	CollectMetrics(ctx context.Context) ([]Metric, error)
	// GetResourceUsage returns current resource usage
	GetResourceUsage() ResourceUsage
	// SetThresholds sets resource usage thresholds
	SetThresholds(thresholds ResourceThresholds) error
}

// Metric represents a system metric
type Metric struct {
	Name       string
	Value      float64
	Labels     map[string]string
	Timestamp  int64
	MetricType MetricType
}

// MetricType represents the type of a metric
type MetricType string

const (
	// MetricTypeGauge represents a gauge metric
	MetricTypeGauge MetricType = "gauge"
	// MetricTypeCounter represents a counter metric
	MetricTypeCounter MetricType = "counter"
	// MetricTypeHistogram represents a histogram metric
	MetricTypeHistogram MetricType = "histogram"
)

// ResourceUsage represents system resource usage
type ResourceUsage struct {
	CPU    float64
	Memory float64
	Disk   float64
}

// ResourceThresholds represents resource usage thresholds
type ResourceThresholds struct {
	CPUThreshold    float64
	MemoryThreshold float64
	DiskThreshold   float64
}

// EventProcessor processes system events
type EventProcessor interface {
	Component
	// ProcessEvent processes a system event
	ProcessEvent(ctx context.Context, event Event) error
	// GetEventTypes returns supported event types
	GetEventTypes() []string
}

// Event represents a system event
type Event struct {
	Type      string
	Source    string
	Timestamp int64
	Severity  string
	Message   string
	Data      map[string]interface{}
}

// TelemetryBridge manages telemetry data
type TelemetryBridge interface {
	Component
	// GetMeterProvider returns the OpenTelemetry meter provider
	GetMeterProvider() metric.MeterProvider
	// GetTracerProvider returns the OpenTelemetry tracer provider
	GetTracerProvider() trace.TracerProvider
	// GetLogger returns the configured logger
	GetLogger() *zap.Logger
}

// ConfigManager manages configuration
type ConfigManager interface {
	Component
	// LoadConfig loads configuration from the given path
	LoadConfig(path string) error
	// GetConfig returns the current configuration
	GetConfig() interface{}
	// ValidateConfig validates the given configuration
	ValidateConfig(cfg interface{}) error
	// WatchConfig watches for configuration changes
	WatchConfig(ctx context.Context) (<-chan interface{}, error)
}
