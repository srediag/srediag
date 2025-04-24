package types

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
)

// ComponentFactory defines the interface for creating components in SREDIAG
type ComponentFactory interface {
	// Type returns the type of the component created by this factory
	Type() component.Type

	// CreateDefaultConfig creates the default configuration for the component
	CreateDefaultConfig() component.Config

	// CreateComponent creates a new instance of the component
	CreateComponent(ctx context.Context, set component.TelemetrySettings, cfg component.Config) (component.Component, error)

	// Capabilities returns the capabilities of the component
	Capabilities() consumer.Capabilities

	// ValidateConfig validates the component configuration
	ValidateConfig(cfg component.Config) error
}

// ReceiverFactory defines the interface for creating receiver components
type ReceiverFactory interface {
	ComponentFactory

	// CreateTracesReceiver creates a trace receiver component
	CreateTracesReceiver(ctx context.Context, set component.TelemetrySettings, cfg component.Config, nextConsumer consumer.Traces) (component.Component, error)

	// CreateMetricsReceiver creates a metrics receiver component
	CreateMetricsReceiver(ctx context.Context, set component.TelemetrySettings, cfg component.Config, nextConsumer consumer.Metrics) (component.Component, error)

	// CreateLogsReceiver creates a logs receiver component
	CreateLogsReceiver(ctx context.Context, set component.TelemetrySettings, cfg component.Config, nextConsumer consumer.Logs) (component.Component, error)
}

// ProcessorFactory defines the interface for creating processor components
type ProcessorFactory interface {
	ComponentFactory

	// CreateTracesProcessor creates a trace processor component
	CreateTracesProcessor(ctx context.Context, set component.TelemetrySettings, cfg component.Config, nextConsumer consumer.Traces) (component.Component, error)

	// CreateMetricsProcessor creates a metrics processor component
	CreateMetricsProcessor(ctx context.Context, set component.TelemetrySettings, cfg component.Config, nextConsumer consumer.Metrics) (component.Component, error)

	// CreateLogsProcessor creates a logs processor component
	CreateLogsProcessor(ctx context.Context, set component.TelemetrySettings, cfg component.Config, nextConsumer consumer.Logs) (component.Component, error)
}

// ExporterFactory defines the interface for creating exporter components
type ExporterFactory interface {
	ComponentFactory

	// CreateTracesExporter creates a trace exporter component
	CreateTracesExporter(ctx context.Context, set component.TelemetrySettings, cfg component.Config) (component.Component, error)

	// CreateMetricsExporter creates a metrics exporter component
	CreateMetricsExporter(ctx context.Context, set component.TelemetrySettings, cfg component.Config) (component.Component, error)

	// CreateLogsExporter creates a logs exporter component
	CreateLogsExporter(ctx context.Context, set component.TelemetrySettings, cfg component.Config) (component.Component, error)
}

// ExtensionFactory defines the interface for creating extension components
type ExtensionFactory interface {
	ComponentFactory

	// CreateExtension creates an extension component
	CreateExtension(ctx context.Context, set component.TelemetrySettings, cfg component.Config) (component.Component, error)
}
