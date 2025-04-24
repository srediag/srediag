package telemetry

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/trace"
	tracenoop "go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/types"
)

// Bridge implements the TelemetryBridge interface
type Bridge struct {
	logger         *zap.Logger
	resource       map[string]string
	meterProvider  metric.MeterProvider
	tracerProvider trace.TracerProvider
	mu             sync.RWMutex
	running        bool
	healthy        bool
	meter          metric.Meter
	tracer         trace.Tracer
}

// NewBridge creates a new telemetry bridge instance
func NewBridge(logger *zap.Logger, res map[string]string) types.ITelemetryBridge {
	if logger == nil {
		logger = zap.NewNop()
	}

	if res == nil {
		res = make(map[string]string)
	}

	return &Bridge{
		logger:         logger,
		resource:       res,
		meterProvider:  noop.NewMeterProvider(),
		tracerProvider: tracenoop.NewTracerProvider(),
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

// GetResource implements ITelemetryBridge
func (b *Bridge) GetResource() map[string]string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.resource
}

// SetResource implements ITelemetryBridge
func (b *Bridge) SetResource(res map[string]string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if res == nil {
		res = make(map[string]string)
	}
	b.resource = res
	return nil
}

// GetMeterProvider implements ITelemetryBridge
func (b *Bridge) GetMeterProvider() types.IMeterProvider {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return &MeterProvider{provider: b.meterProvider}
}

// GetTracerProvider implements ITelemetryBridge
func (b *Bridge) GetTracerProvider() types.ITracerProvider {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return &TracerProvider{provider: b.tracerProvider}
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

// Configure configures the bridge with settings
func (b *Bridge) Configure(settings types.ComponentSettings) error {
	if settings == nil {
		return fmt.Errorf("settings cannot be nil")
	}

	telemetrySettings := settings.GetTelemetrySettings()
	if telemetrySettings == nil {
		return fmt.Errorf("telemetry settings cannot be nil")
	}

	b.logger = telemetrySettings.Logger
	b.tracer = telemetrySettings.Tracer
	b.meter = telemetrySettings.Meter

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

// MeterProvider implements IMeterProvider
func (p *MeterProvider) MeterProvider() metric.MeterProvider {
	return p.provider
}

// Shutdown implements IMeterProvider
func (p *MeterProvider) Shutdown(ctx context.Context) error {
	if shutdownable, ok := p.provider.(interface{ Shutdown(context.Context) error }); ok {
		return shutdownable.Shutdown(ctx)
	}
	return nil
}

// TracerProvider implements ITracerProvider
type TracerProvider struct {
	provider trace.TracerProvider
}

// Tracer implements ITracerProvider
func (p *TracerProvider) Tracer(name string, opts ...trace.TracerOption) trace.Tracer {
	return p.provider.Tracer(name, opts...)
}

// Shutdown implements ITracerProvider
func (p *TracerProvider) Shutdown(ctx context.Context) error {
	if shutdownable, ok := p.provider.(interface{ Shutdown(context.Context) error }); ok {
		return shutdownable.Shutdown(ctx)
	}
	return nil
}

// GetStatus implements IComponent
func (b *Bridge) GetStatus() types.ComponentStatus {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if !b.running {
		return types.ComponentStatusStopped
	}
	if !b.healthy {
		return types.ComponentStatusError
	}
	return types.ComponentStatusRunning
}

// Shutdown implements ITelemetryBridge
func (b *Bridge) Shutdown(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.running {
		return nil
	}

	// Stop the bridge first
	if err := b.Stop(ctx); err != nil {
		return fmt.Errorf("failed to stop bridge: %w", err)
	}

	// Shutdown meter provider
	if mp, ok := b.meterProvider.(interface{ Shutdown(context.Context) error }); ok {
		if err := mp.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown meter provider: %w", err)
		}
	}

	// Shutdown tracer provider
	if tp, ok := b.tracerProvider.(interface{ Shutdown(context.Context) error }); ok {
		if err := tp.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown tracer provider: %w", err)
		}
	}

	b.healthy = false
	return nil
}
