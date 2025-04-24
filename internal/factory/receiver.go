package factory

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"

	"github.com/srediag/srediag/internal/types"
)

// ReceiverFactory provides a base implementation for receiver factories
type ReceiverFactory struct {
	*BaseFactory
	createTraces  func(context.Context, component.TelemetrySettings, component.Config, consumer.Traces) (component.Component, error)
	createMetrics func(context.Context, component.TelemetrySettings, component.Config, consumer.Metrics) (component.Component, error)
	createLogs    func(context.Context, component.TelemetrySettings, component.Config, consumer.Logs) (component.Component, error)
}

// NewReceiverFactory creates a new receiver factory
func NewReceiverFactory(
	typ component.Type,
	version string,
	configType component.Config,
	capabilities consumer.Capabilities,
) *ReceiverFactory {
	return &ReceiverFactory{
		BaseFactory: NewBaseFactory(typ, version, configType, capabilities),
	}
}

// CreateTracesReceiver implements types.ReceiverFactory
func (f *ReceiverFactory) CreateTracesReceiver(
	ctx context.Context,
	set component.TelemetrySettings,
	cfg component.Config,
	nextConsumer consumer.Traces,
) (component.Component, error) {
	if f.createTraces == nil {
		return nil, types.ErrInvalidPluginConfig("CreateTracesReceiver not implemented")
	}
	return f.createTraces(ctx, set, cfg, nextConsumer)
}

// CreateMetricsReceiver implements types.ReceiverFactory
func (f *ReceiverFactory) CreateMetricsReceiver(
	ctx context.Context,
	set component.TelemetrySettings,
	cfg component.Config,
	nextConsumer consumer.Metrics,
) (component.Component, error) {
	if f.createMetrics == nil {
		return nil, types.ErrInvalidPluginConfig("CreateMetricsReceiver not implemented")
	}
	return f.createMetrics(ctx, set, cfg, nextConsumer)
}

// CreateLogsReceiver implements types.ReceiverFactory
func (f *ReceiverFactory) CreateLogsReceiver(
	ctx context.Context,
	set component.TelemetrySettings,
	cfg component.Config,
	nextConsumer consumer.Logs,
) (component.Component, error) {
	if f.createLogs == nil {
		return nil, types.ErrInvalidPluginConfig("CreateLogsReceiver not implemented")
	}
	return f.createLogs(ctx, set, cfg, nextConsumer)
}

// WithTracesReceiver sets the traces receiver creation function
func (f *ReceiverFactory) WithTracesReceiver(
	createFn func(context.Context, component.TelemetrySettings, component.Config, consumer.Traces) (component.Component, error),
) *ReceiverFactory {
	f.createTraces = createFn
	return f
}

// WithMetricsReceiver sets the metrics receiver creation function
func (f *ReceiverFactory) WithMetricsReceiver(
	createFn func(context.Context, component.TelemetrySettings, component.Config, consumer.Metrics) (component.Component, error),
) *ReceiverFactory {
	f.createMetrics = createFn
	return f
}

// WithLogsReceiver sets the logs receiver creation function
func (f *ReceiverFactory) WithLogsReceiver(
	createFn func(context.Context, component.TelemetrySettings, component.Config, consumer.Logs) (component.Component, error),
) *ReceiverFactory {
	f.createLogs = createFn
	return f
}
