package diagnostic

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

// TypeStr is the unique identifier for the Diagnostic exporter.
var typeStr = component.MustNewType("diagnostic")

// CreateFactory creates a factory for the Diagnostic exporter.
func CreateFactory() exporter.Factory {
	return exporter.NewFactory(
		typeStr,
		createDefaultConfig,
		exporter.WithMetrics(createMetricsExporter, component.StabilityLevelBeta),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		OutputDir:     "diagnostic_reports",
		FlushInterval: 5 * time.Minute,
		PrettyPrint:   true,
	}
}

type diagnosticExporter struct {
	cfg    *Config
	logger *zap.Logger
}

func createMetricsExporter(
	_ context.Context,
	params exporter.Settings,
	cfg component.Config,
) (exporter.Metrics, error) {
	eCfg := cfg.(*Config)
	return &diagnosticExporter{
		cfg:    eCfg,
		logger: params.Logger,
	}, nil
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
