package plugins

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/extension"
)

// BaseFactory provides a base implementation of Factory that other factories can embed
type BaseFactory struct {
	typ     component.Type
	version string
}

// NewBaseFactory creates a new BaseFactory
func NewBaseFactory(typ component.Type, version string) BaseFactory {
	return BaseFactory{
		typ:     typ,
		version: version,
	}
}

// Type implements Factory
func (f BaseFactory) Type() component.Type {
	return f.typ
}

// CreateDefaultConfig implements component.Factory
func (f BaseFactory) CreateDefaultConfig() component.Config {
	return &PluginConfig{
		Enabled:  true,
		Settings: make(map[string]string),
	}
}

// ReceiverFactory provides a base implementation for receiver factories
type ReceiverFactory struct {
	BaseFactory
	createTraces  consumer.Traces
	createMetrics consumer.Metrics
	createLogs    consumer.Logs
}

// NewReceiverFactory creates a new ReceiverFactory
func NewReceiverFactory(
	typ component.Type,
	version string,
	createTraces consumer.Traces,
	createMetrics consumer.Metrics,
	createLogs consumer.Logs,
) *ReceiverFactory {
	return &ReceiverFactory{
		BaseFactory:   NewBaseFactory(typ, version),
		createTraces:  createTraces,
		createMetrics: createMetrics,
		createLogs:    createLogs,
	}
}

// ProcessorFactory provides a base implementation for processor factories
type ProcessorFactory struct {
	BaseFactory
	createTraces  consumer.Traces
	createMetrics consumer.Metrics
	createLogs    consumer.Logs
}

// NewProcessorFactory creates a new ProcessorFactory
func NewProcessorFactory(
	typ component.Type,
	version string,
	createTraces consumer.Traces,
	createMetrics consumer.Metrics,
	createLogs consumer.Logs,
) *ProcessorFactory {
	return &ProcessorFactory{
		BaseFactory:   NewBaseFactory(typ, version),
		createTraces:  createTraces,
		createMetrics: createMetrics,
		createLogs:    createLogs,
	}
}

// ExporterFactory provides a base implementation for exporter factories
type ExporterFactory struct {
	BaseFactory
	createTraces  consumer.Traces
	createMetrics consumer.Metrics
	createLogs    consumer.Logs
}

// NewExporterFactory creates a new ExporterFactory
func NewExporterFactory(
	typ component.Type,
	version string,
	createTraces consumer.Traces,
	createMetrics consumer.Metrics,
	createLogs consumer.Logs,
) *ExporterFactory {
	return &ExporterFactory{
		BaseFactory:   NewBaseFactory(typ, version),
		createTraces:  createTraces,
		createMetrics: createMetrics,
		createLogs:    createLogs,
	}
}

// ExtensionFactory provides a base implementation for extension factories
type ExtensionFactory struct {
	BaseFactory
	createExtension func(context.Context, component.TelemetrySettings) (extension.Extension, error)
}

// NewExtensionFactory creates a new ExtensionFactory
func NewExtensionFactory(
	typ component.Type,
	version string,
	createExtension func(context.Context, component.TelemetrySettings) (extension.Extension, error),
) *ExtensionFactory {
	return &ExtensionFactory{
		BaseFactory:     NewBaseFactory(typ, version),
		createExtension: createExtension,
	}
}
