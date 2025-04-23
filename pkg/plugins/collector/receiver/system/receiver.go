package system

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"

	"github.com/srediag/srediag/pkg/plugins"
)

const (
	pluginType    = "receiver/system"
	pluginName    = "system"
	pluginVersion = "v0.1.0"
)

// Config represents the receiver configuration
type Config struct {
	CollectInterval time.Duration `mapstructure:"collect_interval"`
}

// Receiver implements the system metrics receiver
type Receiver struct {
	logger       *zap.Logger
	config       *Config
	nextConsumer consumer.Metrics
	host         component.Host
	shutdownChan chan struct{}
}

// NewFactory creates a new system receiver factory
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

	return &Receiver{
		config:       config,
		shutdownChan: make(chan struct{}),
	}, nil
}

// Type returns the plugin type
func (r *Receiver) Type() string { return pluginType }

// Name returns the plugin name
func (r *Receiver) Name() string { return pluginName }

// Version returns the plugin version
func (r *Receiver) Version() string { return pluginVersion }

// Start initializes the receiver
func (r *Receiver) Start(ctx context.Context, host component.Host) error {
	r.host = host
	r.logger = zap.L().Named("system-receiver")

	go r.collectMetrics(ctx)
	return nil
}

// Shutdown stops the receiver
func (r *Receiver) Shutdown(ctx context.Context) error {
	close(r.shutdownChan)
	return nil
}

// collectMetrics collects system metrics periodically
func (r *Receiver) collectMetrics(ctx context.Context) {
	ticker := time.NewTicker(r.config.CollectInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-r.shutdownChan:
			return
		case <-ticker.C:
			metrics := pmetric.NewMetrics()
			rms := metrics.ResourceMetrics()
			rm := rms.AppendEmpty()

			// Add resource attributes
			attrs := rm.Resource().Attributes()
			attrs.PutStr("service.name", "system-receiver")
			attrs.PutStr("service.version", pluginVersion)

			// Create metrics
			ms := rm.ScopeMetrics().AppendEmpty().Metrics()

			// CPU metrics
			cpuMetric := ms.AppendEmpty()
			cpuMetric.SetName("system.cpu.utilization")
			cpuMetric.SetDescription("CPU utilization")
			cpuMetric.SetUnit("%")
			cpuDP := cpuMetric.SetEmptyGauge().DataPoints().AppendEmpty()
			cpuDP.SetStartTimestamp(pcommon.NewTimestampFromTime(time.Now().Add(-r.config.CollectInterval)))
			cpuDP.SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))

			// Memory metrics
			memMetric := ms.AppendEmpty()
			memMetric.SetName("system.memory.utilization")
			memMetric.SetDescription("Memory utilization")
			memMetric.SetUnit("%")
			memDP := memMetric.SetEmptyGauge().DataPoints().AppendEmpty()
			memDP.SetStartTimestamp(pcommon.NewTimestampFromTime(time.Now().Add(-r.config.CollectInterval)))
			memDP.SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))

			if err := r.nextConsumer.ConsumeMetrics(ctx, metrics); err != nil {
				r.logger.Error("failed to consume metrics", zap.Error(err))
			}
		}
	}
}

// SetNextConsumer sets the next consumer in the pipeline
func (r *Receiver) SetNextConsumer(nextConsumer consumer.Metrics) {
	r.nextConsumer = nextConsumer
}
