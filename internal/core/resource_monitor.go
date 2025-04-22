package core

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"
)

// DefaultResourceMonitor is the default implementation of ResourceMonitor
type DefaultResourceMonitor struct {
	logger     *zap.Logger
	meter      metric.Meter
	thresholds ResourceThresholds
	usage      ResourceUsage
	mu         sync.RWMutex
	healthy    bool
	running    bool
	stopChan   chan struct{}
}

// NewResourceMonitor creates a new instance of DefaultResourceMonitor
func NewResourceMonitor(logger *zap.Logger, meter metric.Meter) *DefaultResourceMonitor {
	return &DefaultResourceMonitor{
		logger:   logger,
		meter:    meter,
		healthy:  true,
		stopChan: make(chan struct{}),
	}
}

// Start initializes the resource monitor
func (rm *DefaultResourceMonitor) Start(ctx context.Context) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if rm.running {
		return fmt.Errorf("resource monitor is already running")
	}

	rm.logger.Info("starting resource monitor")
	rm.running = true

	// Start monitoring in background
	go rm.monitorResources(ctx)

	return nil
}

// Stop stops the resource monitor
func (rm *DefaultResourceMonitor) Stop(ctx context.Context) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if !rm.running {
		return fmt.Errorf("resource monitor is not running")
	}

	rm.logger.Info("stopping resource monitor")
	close(rm.stopChan)
	rm.running = false

	return nil
}

// IsHealthy returns the health status
func (rm *DefaultResourceMonitor) IsHealthy() bool {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return rm.healthy
}

// CollectMetrics collects system metrics
func (rm *DefaultResourceMonitor) CollectMetrics(ctx context.Context) ([]Metric, error) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	if !rm.running {
		return nil, fmt.Errorf("resource monitor is not running")
	}

	metrics := make([]Metric, 0)
	now := time.Now().UnixNano()

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Memory metrics
	metrics = append(metrics, []Metric{
		{
			Name:       "system.memory.alloc",
			Value:      float64(memStats.Alloc),
			MetricType: MetricTypeGauge,
			Timestamp:  now,
			Labels:     map[string]string{"unit": "bytes"},
		},
		{
			Name:       "system.memory.total",
			Value:      float64(memStats.TotalAlloc),
			MetricType: MetricTypeGauge,
			Timestamp:  now,
			Labels:     map[string]string{"unit": "bytes"},
		},
		{
			Name:       "system.memory.heap",
			Value:      float64(memStats.HeapAlloc),
			MetricType: MetricTypeGauge,
			Timestamp:  now,
			Labels:     map[string]string{"unit": "bytes"},
		},
	}...)

	// CPU metrics
	metrics = append(metrics, Metric{
		Name:       "system.cpu.goroutines",
		Value:      float64(runtime.NumGoroutine()),
		MetricType: MetricTypeGauge,
		Timestamp:  now,
		Labels:     map[string]string{"unit": "count"},
	})

	return metrics, nil
}

// GetResourceUsage returns current resource usage
func (rm *DefaultResourceMonitor) GetResourceUsage() ResourceUsage {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return rm.usage
}

// SetThresholds sets resource usage thresholds
func (rm *DefaultResourceMonitor) SetThresholds(thresholds ResourceThresholds) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if thresholds.CPUThreshold < 0 || thresholds.CPUThreshold > 100 {
		return fmt.Errorf("invalid CPU threshold: must be between 0 and 100")
	}
	if thresholds.MemoryThreshold < 0 || thresholds.MemoryThreshold > 100 {
		return fmt.Errorf("invalid memory threshold: must be between 0 and 100")
	}
	if thresholds.DiskThreshold < 0 || thresholds.DiskThreshold > 100 {
		return fmt.Errorf("invalid disk threshold: must be between 0 and 100")
	}

	rm.thresholds = thresholds
	return nil
}

// monitorResources continuously monitors system resources
func (rm *DefaultResourceMonitor) monitorResources(ctx context.Context) {
	// Create OpenTelemetry instruments
	memoryGauge, err := rm.meter.Float64ObservableGauge("system.memory.usage",
		metric.WithDescription("System memory usage"),
		metric.WithUnit("bytes"))
	if err != nil {
		rm.logger.Error("failed to create memory gauge", zap.Error(err))
		return
	}

	cpuGauge, err := rm.meter.Float64ObservableGauge("system.cpu.goroutines",
		metric.WithDescription("Number of goroutines"),
		metric.WithUnit("count"))
	if err != nil {
		rm.logger.Error("failed to create CPU gauge", zap.Error(err))
		return
	}

	_, err = rm.meter.RegisterCallback(func(ctx context.Context, o metric.Observer) error {
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)

		o.ObserveFloat64(memoryGauge, float64(memStats.Alloc),
			metric.WithAttributes(attribute.String("type", "alloc")))
		o.ObserveFloat64(memoryGauge, float64(memStats.HeapAlloc),
			metric.WithAttributes(attribute.String("type", "heap")))
		o.ObserveFloat64(cpuGauge, float64(runtime.NumGoroutine()))

		rm.mu.Lock()
		rm.usage = ResourceUsage{
			Memory: float64(memStats.Alloc) / float64(memStats.Sys) * 100,
			CPU:    float64(runtime.NumGoroutine()),
		}

		// Check thresholds
		rm.healthy = true
		if rm.usage.Memory > rm.thresholds.MemoryThreshold {
			rm.healthy = false
			rm.logger.Warn("memory usage above threshold",
				zap.Float64("usage", rm.usage.Memory),
				zap.Float64("threshold", rm.thresholds.MemoryThreshold))
		}
		rm.mu.Unlock()

		return nil
	}, memoryGauge, cpuGauge)

	if err != nil {
		rm.logger.Error("failed to register callback", zap.Error(err))
		return
	}

	<-rm.stopChan
}
