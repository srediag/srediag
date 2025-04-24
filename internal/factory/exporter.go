package factory

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"

	"github.com/srediag/srediag/internal/types"
)

// ExporterFactory provides a base implementation for exporter factories
type ExporterFactory struct {
	*BaseFactory
	createTraces  func(context.Context, component.TelemetrySettings, component.Config) (component.Component, error)
	createMetrics func(context.Context, component.TelemetrySettings, component.Config) (component.Component, error)
	createLogs    func(context.Context, component.TelemetrySettings, component.Config) (component.Component, error)
}

// NewExporterFactory creates a new exporter factory
func NewExporterFactory(
	typ component.Type,
	version string,
	configType component.Config,
	capabilities consumer.Capabilities,
) *ExporterFactory {
	return &ExporterFactory{
		BaseFactory: NewBaseFactory(typ, version, configType, capabilities),
	}
}

// CreateTracesExporter implements types.ExporterFactory
func (f *ExporterFactory) CreateTracesExporter(
	ctx context.Context,
	set component.TelemetrySettings,
	cfg component.Config,
) (component.Component, error) {
	if f.createTraces == nil {
		return nil, types.ErrInvalidPluginConfig("CreateTracesExporter not implemented")
	}
	return f.createTraces(ctx, set, cfg)
}

// CreateMetricsExporter implements types.ExporterFactory
func (f *ExporterFactory) CreateMetricsExporter(
	ctx context.Context,
	set component.TelemetrySettings,
	cfg component.Config,
) (component.Component, error) {
	if f.createMetrics == nil {
		return nil, types.ErrInvalidPluginConfig("CreateMetricsExporter not implemented")
	}
	return f.createMetrics(ctx, set, cfg)
}

// CreateLogsExporter implements types.ExporterFactory
func (f *ExporterFactory) CreateLogsExporter(
	ctx context.Context,
	set component.TelemetrySettings,
	cfg component.Config,
) (component.Component, error) {
	if f.createLogs == nil {
		return nil, types.ErrInvalidPluginConfig("CreateLogsExporter not implemented")
	}
	return f.createLogs(ctx, set, cfg)
}

// WithTracesExporter sets the traces exporter creation function
func (f *ExporterFactory) WithTracesExporter(
	createFn func(context.Context, component.TelemetrySettings, component.Config) (component.Component, error),
) *ExporterFactory {
	f.createTraces = createFn
	return f
}

// WithMetricsExporter sets the metrics exporter creation function
func (f *ExporterFactory) WithMetricsExporter(
	createFn func(context.Context, component.TelemetrySettings, component.Config) (component.Component, error),
) *ExporterFactory {
	f.createMetrics = createFn
	return f
}

// WithLogsExporter sets the logs exporter creation function
func (f *ExporterFactory) WithLogsExporter(
	createFn func(context.Context, component.TelemetrySettings, component.Config) (component.Component, error),
) *ExporterFactory {
	f.createLogs = createFn
	return f
}
