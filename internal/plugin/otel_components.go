package plugin

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/connector"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/otlpexporter"
	"go.opentelemetry.io/collector/exporter/otlphttpexporter"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/batchprocessor"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
	"go.uber.org/zap"
)

// OTelComponentLoader handles loading of OpenTelemetry components
type OTelComponentLoader struct {
	logger *zap.Logger
}

// NewOTelComponentLoader creates a new OpenTelemetry component loader
func NewOTelComponentLoader(logger *zap.Logger) *OTelComponentLoader {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &OTelComponentLoader{
		logger: logger,
	}
}

// GetBuiltinFactories returns commonly used OpenTelemetry component factories
func (l *OTelComponentLoader) GetBuiltinFactories() (
	map[component.Type]receiver.Factory,
	map[component.Type]processor.Factory,
	map[component.Type]exporter.Factory,
	map[component.Type]connector.Factory,
	error) {

	receivers := make(map[component.Type]receiver.Factory)
	processors := make(map[component.Type]processor.Factory)
	exporters := make(map[component.Type]exporter.Factory)
	connectors := make(map[component.Type]connector.Factory)

	// Add OTLP receiver
	receivers[otlpreceiver.NewFactory().Type()] = otlpreceiver.NewFactory()

	// Add batch processor
	processors[batchprocessor.NewFactory().Type()] = batchprocessor.NewFactory()

	// Add OTLP exporters
	exporters[otlpexporter.NewFactory().Type()] = otlpexporter.NewFactory()
	exporters[otlphttpexporter.NewFactory().Type()] = otlphttpexporter.NewFactory()

	l.logger.Info("Loaded built-in OpenTelemetry component factories",
		zap.Int("receivers", len(receivers)),
		zap.Int("processors", len(processors)),
		zap.Int("exporters", len(exporters)),
		zap.Int("connectors", len(connectors)))

	return receivers, processors, exporters, connectors, nil
}

// RegisterBuiltinFactories registers built-in OpenTelemetry component factories
// with the provided loader
func (l *OTelComponentLoader) RegisterBuiltinFactories(loader *Loader) error {
	receivers, processors, exporters, connectors, err := l.GetBuiltinFactories()
	if err != nil {
		return err
	}

	// Register all factories with the loader
	for _, factory := range receivers {
		if err := loader.RegisterReceiverFactory(factory); err != nil {
			return err
		}
	}

	for _, factory := range processors {
		if err := loader.RegisterProcessorFactory(factory); err != nil {
			return err
		}
	}

	for _, factory := range exporters {
		if err := loader.RegisterExporterFactory(factory); err != nil {
			return err
		}
	}

	for _, factory := range connectors {
		if err := loader.RegisterConnectorFactory(factory); err != nil {
			return err
		}
	}

	return nil
}
