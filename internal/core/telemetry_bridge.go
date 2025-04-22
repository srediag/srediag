package core

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// DefaultTelemetryBridge is the default implementation of TelemetryBridge
type DefaultTelemetryBridge struct {
	logger         *zap.Logger
	meterProvider  metric.MeterProvider
	tracerProvider trace.TracerProvider
	mu             sync.RWMutex
	healthy        bool
	running        bool
}

// NewTelemetryBridge creates a new instance of DefaultTelemetryBridge
func NewTelemetryBridge(logger *zap.Logger, res *resource.Resource) *DefaultTelemetryBridge {
	return &DefaultTelemetryBridge{
		logger:  logger,
		healthy: true,
	}
}

// Start initializes the telemetry bridge
func (tb *DefaultTelemetryBridge) Start(ctx context.Context) error {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	if tb.running {
		return fmt.Errorf("telemetry bridge is already running")
	}

	tb.logger.Info("starting telemetry bridge")

	// Set global providers
	otel.SetMeterProvider(tb.meterProvider)
	otel.SetTracerProvider(tb.tracerProvider)

	tb.running = true
	return nil
}

// Stop stops the telemetry bridge
func (tb *DefaultTelemetryBridge) Stop(ctx context.Context) error {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	if !tb.running {
		return fmt.Errorf("telemetry bridge is not running")
	}

	tb.logger.Info("stopping telemetry bridge")
	tb.running = false
	return nil
}

// IsHealthy returns the health status
func (tb *DefaultTelemetryBridge) IsHealthy() bool {
	tb.mu.RLock()
	defer tb.mu.RUnlock()
	return tb.healthy
}

// GetMeterProvider returns the OpenTelemetry meter provider
func (tb *DefaultTelemetryBridge) GetMeterProvider() metric.MeterProvider {
	return tb.meterProvider
}

// GetTracerProvider returns the OpenTelemetry tracer provider
func (tb *DefaultTelemetryBridge) GetTracerProvider() trace.TracerProvider {
	return tb.tracerProvider
}

// GetLogger returns the configured logger
func (tb *DefaultTelemetryBridge) GetLogger() *zap.Logger {
	return tb.logger
}

// SetMeterProvider sets the meter provider
func (tb *DefaultTelemetryBridge) SetMeterProvider(provider metric.MeterProvider) {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.meterProvider = provider
	if tb.running {
		otel.SetMeterProvider(provider)
	}
}

// SetTracerProvider sets the tracer provider
func (tb *DefaultTelemetryBridge) SetTracerProvider(provider trace.TracerProvider) {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.tracerProvider = provider
	if tb.running {
		otel.SetTracerProvider(provider)
	}
}
