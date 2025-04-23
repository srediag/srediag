package diagnostic

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"

	"github.com/srediag/srediag/pkg/plugins"
)

const (
	pluginType    = "processor/diagnostic"
	pluginName    = "diagnostic"
	pluginVersion = "v0.1.0"
)

// Config represents the processor configuration
type Config struct {
	CPUThreshold float64 `mapstructure:"cpu_threshold"`
	MemThreshold float64 `mapstructure:"mem_threshold"`
	AlertOnWarn  bool    `mapstructure:"alert_on_warn"`
}

// Processor implements the diagnostic processor
type Processor struct {
	logger       *zap.Logger
	config       *Config
	nextConsumer consumer.Metrics
	host         component.Host
}

// NewFactory creates a new diagnostic processor factory
func NewFactory() plugins.Factory {
	return &factory{}
}

type factory struct{}

func (f *factory) Type() string { return pluginType }

func (f *factory) CreatePlugin(cfg interface{}) (plugins.Plugin, error) {
	config, ok := cfg.(*Config)
	if !ok {
		return nil, fmt.Errorf("invalid configuration type")
	}

	return &Processor{
		config: config,
		logger: zap.L().Named("diagnostic-processor"),
	}, nil
}

// Type returns the plugin type
func (p *Processor) Type() string { return pluginType }

// Name returns the plugin name
func (p *Processor) Name() string { return pluginName }

// Version returns the plugin version
func (p *Processor) Version() string { return pluginVersion }

// Start initializes the processor
func (p *Processor) Start(ctx context.Context, host component.Host) error {
	p.host = host
	return nil
}

// Shutdown stops the processor
func (p *Processor) Shutdown(ctx context.Context) error {
	return nil
}

// ConsumeMetrics processes the metrics data
func (p *Processor) ConsumeMetrics(ctx context.Context, md pmetric.Metrics) error {
	rms := md.ResourceMetrics()
	for i := 0; i < rms.Len(); i++ {
		rm := rms.At(i)
		sms := rm.ScopeMetrics()
		for j := 0; j < sms.Len(); j++ {
			sm := sms.At(j)
			metrics := sm.Metrics()
			for k := 0; k < metrics.Len(); k++ {
				metric := metrics.At(k)
				p.analyzeMetric(metric)
			}
		}
	}

	return p.nextConsumer.ConsumeMetrics(ctx, md)
}

func (p *Processor) analyzeMetric(metric pmetric.Metric) {
	switch metric.Name() {
	case "system.cpu.usage":
		p.analyzeCPUMetric(metric)
	case "system.memory.usage":
		p.analyzeMemoryMetric(metric)
	}
}

func (p *Processor) analyzeCPUMetric(metric pmetric.Metric) {
	if metric.Type() != pmetric.MetricTypeGauge {
		return
	}

	dps := metric.Gauge().DataPoints()
	for i := 0; i < dps.Len(); i++ {
		dp := dps.At(i)
		value := dp.DoubleValue()
		if value > p.config.CPUThreshold {
			p.logger.Warn("CPU usage exceeds threshold",
				zap.Float64("value", value),
				zap.Float64("threshold", p.config.CPUThreshold))
		}
	}
}

func (p *Processor) analyzeMemoryMetric(metric pmetric.Metric) {
	if metric.Type() != pmetric.MetricTypeGauge {
		return
	}

	dps := metric.Gauge().DataPoints()
	for i := 0; i < dps.Len(); i++ {
		dp := dps.At(i)
		value := dp.DoubleValue()
		if value > p.config.MemThreshold {
			p.logger.Warn("Memory usage exceeds threshold",
				zap.Float64("value", value),
				zap.Float64("threshold", p.config.MemThreshold))
		}
	}
}

// SetNextConsumer sets the next consumer in the pipeline
func (p *Processor) SetNextConsumer(nextConsumer consumer.Metrics) {
	p.nextConsumer = nextConsumer
}
