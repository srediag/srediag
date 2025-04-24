package telemetry

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
	tracenoop "go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/types"
)

// Bridge implements the TelemetryBridge interface
type Bridge struct {
	logger         *zap.Logger
	resource       *resource.Resource
	meterProvider  types.IMeterProvider
	tracerProvider types.ITracerProvider
	mu             sync.RWMutex
	running        bool
	healthy        bool
}

// NewBridge creates a new telemetry bridge instance
func NewBridge(logger *zap.Logger, res *resource.Resource) types.ITelemetryBridge {
	if logger == nil {
		logger = zap.NewNop()
	}

	meterProvider := &MeterProvider{
		provider: noop.NewMeterProvider(),
	}
	tracerProvider := &TracerProvider{
		provider: tracenoop.NewTracerProvider(),
	}

	return &Bridge{
		logger:         logger,
		resource:       res,
		meterProvider:  meterProvider,
		tracerProvider: tracerProvider,
		healthy:        true,
	}
}

// Start implements ITelemetryBridge
func (b *Bridge) Start(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.running {
		return nil
	}

	b.running = true
	return nil
}

// Stop implements ITelemetryBridge
func (b *Bridge) Stop(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.running {
		return nil
	}

	b.running = false
	return nil
}

// SetResource implements ITelemetryBridge
func (b *Bridge) SetResource(res *resource.Resource) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.resource = res
}

// GetResource implements ITelemetryBridge
func (b *Bridge) GetResource() *resource.Resource {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.resource
}

// GetMeterProvider implements ITelemetryBridge
func (b *Bridge) GetMeterProvider() types.IMeterProvider {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.meterProvider
}

// GetTracerProvider implements ITelemetryBridge
func (b *Bridge) GetTracerProvider() types.ITracerProvider {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.tracerProvider
}

// IsHealthy implements ITelemetryBridge
func (b *Bridge) IsHealthy() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.healthy
}

// GetName implements ITelemetryBridge
func (b *Bridge) GetName() string {
	return "telemetry-bridge"
}

// GetVersion implements ITelemetryBridge
func (b *Bridge) GetVersion() string {
	return "1.0.0"
}

// GetType implements ITelemetryBridge
func (b *Bridge) GetType() types.ComponentType {
	return types.ComponentTypeService
}

// Configure implements ITelemetryBridge
func (b *Bridge) Configure(cfg interface{}) error {
	return nil
}

// MeterProvider implements IMeterProvider
type MeterProvider struct {
	provider metric.MeterProvider
}

// Meter implements IMeterProvider
func (p *MeterProvider) Meter(name string, opts ...metric.MeterOption) metric.Meter {
	return p.provider.Meter(name, opts...)
}

// TracerProvider implements ITracerProvider
type TracerProvider struct {
	provider trace.TracerProvider
}

// Tracer implements ITracerProvider
func (p *TracerProvider) Tracer(name string, opts ...trace.TracerOption) trace.Tracer {
	return p.provider.Tracer(name, opts...)
}

func NewTracerProvider() *TracerProvider {
	return &TracerProvider{
		provider: tracenoop.NewTracerProvider(),
	}
}
