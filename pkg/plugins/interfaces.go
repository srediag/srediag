package plugins

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
)

// Plugin represents the base interface that all plugins must implement
type Plugin interface {
	component.Component
	Type() string
	Name() string
	Version() string
}

// DiagnosticPlugin represents a plugin that performs system diagnostics
type DiagnosticPlugin interface {
	Plugin
	Diagnose(ctx context.Context, target string, options map[string]interface{}) (DiagnosticResult, error)
}

// DiagnosticResult represents the result of a diagnostic operation
type DiagnosticResult struct {
	Status   string                 `json:"status"`
	Message  string                 `json:"message"`
	Details  map[string]interface{} `json:"details,omitempty"`
	Severity string                 `json:"severity,omitempty"`
}

// CollectorPlugin represents a plugin that integrates with OpenTelemetry Collector
type CollectorPlugin interface {
	Plugin
	component.Component
	consumer.Traces
	consumer.Metrics
	consumer.Logs
}

// IntegrationPlugin represents a plugin that provides integration with external systems
type IntegrationPlugin interface {
	Plugin
	Connect(ctx context.Context, config map[string]interface{}) error
	Execute(ctx context.Context, action string, params map[string]interface{}) (interface{}, error)
}

// ManagementPlugin represents a plugin that manages system resources or configurations
type ManagementPlugin interface {
	Plugin
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
	Type() string
	CreatePlugin(config interface{}) (Plugin, error)
}
