package plugins

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/extension"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/receiver"
)

// BasePlugin represents the base interface that all plugins must implement
type BasePlugin interface {
	component.Component
	Type() component.Type
	Name() string
	Version() string
}

// ReceiverPlugin represents a plugin that can receive telemetry data
type ReceiverPlugin interface {
	BasePlugin
	receiver.Traces
	receiver.Metrics
	receiver.Logs
}

// ProcessorPlugin represents a plugin that can process telemetry data
type ProcessorPlugin interface {
	BasePlugin
	processor.Traces
	processor.Metrics
	processor.Logs
}

// ExporterPlugin represents a plugin that can export telemetry data
type ExporterPlugin interface {
	BasePlugin
	exporter.Traces
	exporter.Metrics
	exporter.Logs
}

// ExtensionPlugin represents a plugin that extends collector functionality
type ExtensionPlugin interface {
	BasePlugin
	extension.Extension
}

// DiagnosticPlugin represents a plugin that performs system diagnostics
type DiagnosticPlugin interface {
	BasePlugin
	Diagnose(ctx context.Context, target string, options map[string]interface{}) (DiagnosticResult, error)
}

// DiagnosticResult represents the result of a diagnostic operation
type DiagnosticResult struct {
	Status   string                 `json:"status"`
	Message  string                 `json:"message"`
	Details  map[string]interface{} `json:"details,omitempty"`
	Severity string                 `json:"severity,omitempty"`
	Metrics  pmetric.Metrics        `json:"metrics,omitempty"`
	Traces   ptrace.Traces          `json:"traces,omitempty"`
}

// CollectorPlugin represents a plugin that integrates with OpenTelemetry Collector
type CollectorPlugin interface {
	BasePlugin
	component.Component
	consumer.Traces
	consumer.Metrics
	consumer.Logs
}

// IntegrationPlugin represents a plugin that provides integration with external systems
type IntegrationPlugin interface {
	BasePlugin
	Connect(ctx context.Context, config map[string]interface{}) error
	Execute(ctx context.Context, action string, params map[string]interface{}) (interface{}, error)
}

// ManagementPlugin represents a plugin that manages system resources or configurations
type ManagementPlugin interface {
	BasePlugin
	Apply(ctx context.Context, target string, config map[string]interface{}) error
	Validate(ctx context.Context, config map[string]interface{}) error
	Status(ctx context.Context, target string) (ManagementStatus, error)
}

// ManagementStatus represents the status of a managed resource
type ManagementStatus struct {
	State    string                 `json:"state"`
	Health   string                 `json:"health"`
	Details  map[string]interface{} `json:"details,omitempty"`
	Warnings []string               `json:"warnings,omitempty"`
}

// Factory creates new instances of plugins
type Factory interface {
	component.Factory
	CreatePlugin(config interface{}) (BasePlugin, error)
}

// PluginConfig represents the base configuration for all plugins
type PluginConfig struct {
	component.Config `mapstructure:",squash"`
	Enabled          bool              `mapstructure:"enabled"`
	Settings         map[string]string `mapstructure:"settings"`
}
