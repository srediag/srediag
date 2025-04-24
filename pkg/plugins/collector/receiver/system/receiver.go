package system

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver"
	semconv "go.opentelemetry.io/collector/semconv/v1.26.0"
	"go.uber.org/zap"
)

const (
	typeStr = "system"
	// Default values
	defaultCollectInterval = 30 * time.Second
	minCollectInterval     = 1 * time.Second
)

var receiverType = component.MustNewType(typeStr)

// Config represents the receiver configuration for collecting system metrics
type Config struct {
	CollectInterval time.Duration `mapstructure:"collect_interval"`
}

// Validate checks if the receiver configuration is valid
func (cfg *Config) Validate() error {
	if cfg.CollectInterval < minCollectInterval {
		return fmt.Errorf("collect_interval must be greater than or equal to %v", minCollectInterval)
	}
	return nil
}

// Receiver implements the system metrics receiver
type Receiver struct {
	logger       *zap.Logger
	config       *Config
	nextConsumer consumer.Metrics
	shutdownChan chan struct{}
}

// NewFactory creates a new system receiver factory
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		receiverType,
		createDefaultConfig,
		receiver.WithMetrics(createMetricsReceiver, component.StabilityLevelBeta),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		CollectInterval: defaultCollectInterval,
	}
}

func createMetricsReceiver(
	_ context.Context,
	set receiver.Settings,
	cfg component.Config,
	nextConsumer consumer.Metrics,
) (receiver.Metrics, error) {
	config, ok := cfg.(*Config)
	if !ok {
		return nil, fmt.Errorf("invalid configuration type")
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &Receiver{
		config:       config,
		nextConsumer: nextConsumer,
		logger:       set.Logger,
		shutdownChan: make(chan struct{}),
	}, nil
}

// Start initializes the receiver
func (r *Receiver) Start(ctx context.Context, host component.Host) error {
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

	// Get hostname once at startup
	hostname, err := os.Hostname()
	if err != nil {
		r.logger.Warn("Failed to get hostname", zap.Error(err))
		hostname = "unknown"
	}

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

			// Add resource attributes following OTel semantic conventions
			attrs := rm.Resource().Attributes()
			attrs.PutStr(semconv.AttributeServiceName, "system-receiver")
			attrs.PutStr(semconv.AttributeServiceVersion, "v0.1.0")
			attrs.PutStr(semconv.AttributeHostName, hostname)

			// Create metrics
			ms := rm.ScopeMetrics().AppendEmpty().Metrics()
			now := time.Now()
			startTime := now.Add(-r.config.CollectInterval)

			// Collect CPU metrics
			if cpuPercent, err := cpu.Percent(0, false); err == nil {
				cpuMetric := ms.AppendEmpty()
				cpuMetric.SetName("system.cpu.utilization")
				cpuMetric.SetDescription("CPU utilization percentage")
				cpuMetric.SetUnit("%")
				cpuDP := cpuMetric.SetEmptyGauge().DataPoints().AppendEmpty()
				cpuDP.SetStartTimestamp(pcommon.NewTimestampFromTime(startTime))
				cpuDP.SetTimestamp(pcommon.NewTimestampFromTime(now))
				if len(cpuPercent) > 0 {
					cpuDP.SetDoubleValue(cpuPercent[0])
				}
			} else {
				r.logger.Warn("Failed to collect CPU metrics", zap.Error(err))
			}

			// Collect Memory metrics
			if vmStat, err := mem.VirtualMemory(); err == nil {
				memMetric := ms.AppendEmpty()
				memMetric.SetName("system.memory.utilization")
				memMetric.SetDescription("Memory utilization percentage")
				memMetric.SetUnit("%")
				memDP := memMetric.SetEmptyGauge().DataPoints().AppendEmpty()
				memDP.SetStartTimestamp(pcommon.NewTimestampFromTime(startTime))
				memDP.SetTimestamp(pcommon.NewTimestampFromTime(now))
				memDP.SetDoubleValue(vmStat.UsedPercent)

				// Add memory usage in bytes
				memUsageMetric := ms.AppendEmpty()
				memUsageMetric.SetName("system.memory.usage")
				memUsageMetric.SetDescription("Memory usage in bytes")
				memUsageMetric.SetUnit("By")
				memUsageDP := memUsageMetric.SetEmptyGauge().DataPoints().AppendEmpty()
				memUsageDP.SetStartTimestamp(pcommon.NewTimestampFromTime(startTime))
				memUsageDP.SetTimestamp(pcommon.NewTimestampFromTime(now))
				memUsageDP.SetDoubleValue(float64(vmStat.Used))
			} else {
				r.logger.Warn("Failed to collect memory metrics", zap.Error(err))
			}

			if err := r.nextConsumer.ConsumeMetrics(ctx, metrics); err != nil {
				r.logger.Error("Failed to consume metrics", zap.Error(err))
			}
		}
	}
}
