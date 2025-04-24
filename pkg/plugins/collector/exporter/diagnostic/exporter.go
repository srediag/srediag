package diagnostic

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

// Config defines configuration for diagnostic exporter.
type Config struct {
	FlushInterval time.Duration `mapstructure:"flush_interval"`
	PrettyPrint   bool          `mapstructure:"pretty_print"`
}

// Exporter implements the OpenTelemetry metrics exporter
type Exporter struct {
	settings component.TelemetrySettings
	config   *Config
}

// NewFactory creates a factory for the diagnostic exporter
func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		component.MustNewType("diagnostic"),
		createDefaultConfig,
		exporter.WithMetrics(createMetricsExporter, component.StabilityLevelDevelopment),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		FlushInterval: 30 * time.Second,
		PrettyPrint:   true,
	}
}

// createMetricsExporter creates a new instance of metrics exporter.
func createMetricsExporter(
	_ context.Context,
	set exporter.Settings,
	cfg component.Config,
) (exporter.Metrics, error) {
	config, ok := cfg.(*Config)
	if !ok {
		return nil, fmt.Errorf("invalid configuration type")
	}

	exp := &Exporter{
		settings: set.TelemetrySettings,
		config:   config,
	}
	return exp, nil
}

// Start tells the exporter to start. The exporter may prepare for exporting
// by connecting to the endpoint. Host parameter can be used for communicating
// with the host after Start() has already returned.
func (e *Exporter) Start(ctx context.Context, host component.Host) error {
	return nil
}

// Shutdown is invoked during shutdown.
func (e *Exporter) Shutdown(ctx context.Context) error {
	return nil
}

// Capabilities implements consumer.Metrics
func (e *Exporter) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

// ConsumeMetrics receives metrics data and exports it into a file.
func (e *Exporter) ConsumeMetrics(ctx context.Context, md pmetric.Metrics) error {
	return nil
}
