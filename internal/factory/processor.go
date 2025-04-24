package factory

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"

	"github.com/srediag/srediag/internal/types"
)

// ProcessorFactory provides a base implementation for processor factories
type ProcessorFactory struct {
	*BaseFactory
	createTraces  func(context.Context, component.TelemetrySettings, component.Config, consumer.Traces) (component.Component, error)
	createMetrics func(context.Context, component.TelemetrySettings, component.Config, consumer.Metrics) (component.Component, error)
	createLogs    func(context.Context, component.TelemetrySettings, component.Config, consumer.Logs) (component.Component, error)
}

// NewProcessorFactory creates a new processor factory
func NewProcessorFactory(
	typ component.Type,
	version string,
	configType component.Config,
	capabilities consumer.Capabilities,
) *ProcessorFactory {
	return &ProcessorFactory{
		BaseFactory: NewBaseFactory(typ, version, configType, capabilities),
	}
}

// CreateTracesProcessor implements types.ProcessorFactory
func (f *ProcessorFactory) CreateTracesProcessor(
	ctx context.Context,
	set component.TelemetrySettings,
	cfg component.Config,
	nextConsumer consumer.Traces,
) (component.Component, error) {
	if f.createTraces == nil {
		return nil, types.ErrInvalidPluginConfig("CreateTracesProcessor not implemented")
	}
	return f.createTraces(ctx, set, cfg, nextConsumer)
}

// CreateMetricsProcessor implements types.ProcessorFactory
func (f *ProcessorFactory) CreateMetricsProcessor(
	ctx context.Context,
	set component.TelemetrySettings,
	cfg component.Config,
	nextConsumer consumer.Metrics,
) (component.Component, error) {
	if f.createMetrics == nil {
		return nil, types.ErrInvalidPluginConfig("CreateMetricsProcessor not implemented")
	}
	return f.createMetrics(ctx, set, cfg, nextConsumer)
}

// CreateLogsProcessor implements types.ProcessorFactory
func (f *ProcessorFactory) CreateLogsProcessor(
	ctx context.Context,
	set component.TelemetrySettings,
	cfg component.Config,
	nextConsumer consumer.Logs,
) (component.Component, error) {
	if f.createLogs == nil {
		return nil, types.ErrInvalidPluginConfig("CreateLogsProcessor not implemented")
	}
	return f.createLogs(ctx, set, cfg, nextConsumer)
}

// WithTracesProcessor sets the traces processor creation function
func (f *ProcessorFactory) WithTracesProcessor(
	createFn func(context.Context, component.TelemetrySettings, component.Config, consumer.Traces) (component.Component, error),
) *ProcessorFactory {
	f.createTraces = createFn
	return f
}

// WithMetricsProcessor sets the metrics processor creation function
func (f *ProcessorFactory) WithMetricsProcessor(
	createFn func(context.Context, component.TelemetrySettings, component.Config, consumer.Metrics) (component.Component, error),
) *ProcessorFactory {
	f.createMetrics = createFn
	return f
}

// WithLogsProcessor sets the logs processor creation function
func (f *ProcessorFactory) WithLogsProcessor(
	createFn func(context.Context, component.TelemetrySettings, component.Config, consumer.Logs) (component.Component, error),
) *ProcessorFactory {
	f.createLogs = createFn
	return f
}
