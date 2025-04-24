package main

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver"
)

const (
	typeStr = "simple"
)

// Factory is the symbol that will be loaded by the plugin system
var Factory receiver.Factory

// Config represents the receiver config settings within the collector's config.yaml
type Config struct {
	Metric string `mapstructure:"metric"`
	Value  int64  `mapstructure:"value"`
}

// Receiver is the type that provides metric telemetry data to the collector
type Receiver struct {
	consumer consumer.Metrics
	config   *Config
}

// init registers the receiver factory
func init() {
	componentType := component.MustNewType(typeStr)
	Factory = receiver.NewFactory(
		componentType,
		func() component.Config {
			return &Config{
				Metric: "example.metric",
				Value:  0,
			}
		},
		receiver.WithMetrics(createMetricsReceiver, component.StabilityLevelDevelopment),
	)
}

// createMetricsReceiver creates a metrics receiver based on provided config.
func createMetricsReceiver(
	_ context.Context,
	set receiver.Settings,
	cfg component.Config,
	nextConsumer consumer.Metrics,
) (receiver.Metrics, error) {
	if nextConsumer == nil {
		return nil, fmt.Errorf("next consumer cannot be nil")
	}

	rCfg, ok := cfg.(*Config)
	if !ok {
		return nil, fmt.Errorf("failed to cast configuration to %s", typeStr)
	}

	return &Receiver{
		config:   rCfg,
		consumer: nextConsumer,
	}, nil
}

// Start implements the component.Component interface.
func (r *Receiver) Start(_ context.Context, _ component.Host) error {
	// Create a new metric record
	metrics := pmetric.NewMetrics()
	rm := metrics.ResourceMetrics().AppendEmpty()
	sm := rm.ScopeMetrics().AppendEmpty()
	m := sm.Metrics().AppendEmpty()
	m.SetName(r.config.Metric)
	dp := m.SetEmptyGauge().DataPoints().AppendEmpty()
	dp.SetIntValue(r.config.Value)

	// Send the metric
	return r.consumer.ConsumeMetrics(context.Background(), metrics)
}

// Shutdown implements the component.Component interface.
func (r *Receiver) Shutdown(context.Context) error {
	return nil
}
