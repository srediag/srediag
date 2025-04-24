package diagnostic

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/base"
)

var typeStr = component.MustNewType("diagnostic")

// CreateFactory creates a factory for the diagnostic exporter
func CreateFactory() exporter.Factory {
	return exporter.NewFactory(
		typeStr,
		createDefaultConfig,
		exporter.WithMetrics(createMetricsExporter, component.StabilityLevelBeta),
	)
}

type diagnosticExporter struct {
	*base.BaseComponent
	logger *zap.Logger
	cfg    *Config
}

// Start implements component.Component
func (e *diagnosticExporter) Start(_ context.Context, _ component.Host) error {
	return nil
}

// Shutdown implements component.Component
func (e *diagnosticExporter) Shutdown(context.Context) error {
	return nil
}

// Capabilities implements exporter.Metrics
func (e *diagnosticExporter) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

// ConsumeMetrics implements exporter.Metrics
func (e *diagnosticExporter) ConsumeMetrics(_ context.Context, md pmetric.Metrics) error {
	// TODO: Implement metric consumption logic
	return nil
}
