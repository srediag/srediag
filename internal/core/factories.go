package core

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/otlpexporter"
	"go.opentelemetry.io/collector/extension"
	"go.opentelemetry.io/collector/otelcol"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/batchprocessor"
	"go.opentelemetry.io/collector/processor/memorylimiterprocessor"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
)

// CreateFactories creates the default set of component factories for the collector
func CreateFactories() (otelcol.Factories, error) {
	factories := otelcol.Factories{
		Receivers:        make(map[component.Type]receiver.Factory),
		Processors:       make(map[component.Type]processor.Factory),
		Exporters:        make(map[component.Type]exporter.Factory),
		Extensions:       make(map[component.Type]extension.Factory),
		ReceiverModules:  make(map[component.Type]string),
		ProcessorModules: make(map[component.Type]string),
		ExporterModules:  make(map[component.Type]string),
		ExtensionModules: make(map[component.Type]string),
	}

	// Add default receiver factories
	otlpRcvFactory := otlpreceiver.NewFactory()
	factories.Receivers[otlpRcvFactory.Type()] = otlpRcvFactory
	factories.ReceiverModules[otlpRcvFactory.Type()] = "go.opentelemetry.io/collector/receiver/otlpreceiver"

	// Add default processor factories
	batchProcFactory := batchprocessor.NewFactory()
	factories.Processors[batchProcFactory.Type()] = batchProcFactory
	factories.ProcessorModules[batchProcFactory.Type()] = "go.opentelemetry.io/collector/processor/batchprocessor"

	memLimiterProcFactory := memorylimiterprocessor.NewFactory()
	factories.Processors[memLimiterProcFactory.Type()] = memLimiterProcFactory
	factories.ProcessorModules[memLimiterProcFactory.Type()] = "go.opentelemetry.io/collector/processor/memorylimiterprocessor"

	// Add default exporter factories
	otlpExpFactory := otlpexporter.NewFactory()
	factories.Exporters[otlpExpFactory.Type()] = otlpExpFactory
	factories.ExporterModules[otlpExpFactory.Type()] = "go.opentelemetry.io/collector/exporter/otlpexporter"

	return factories, nil
}
