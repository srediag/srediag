package base

import (
	"context"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/types"
)

// TelemetryComponent provides a base implementation of the TelemetryComponent interface
type TelemetryComponent struct {
	logger   *zap.Logger
	name     string
	settings types.TelemetrySettings
	resource *resource.Resource
	ctype    types.ComponentType
	healthy  bool
	running  bool
}

// NewTelemetryComponent creates a new telemetry component
func NewTelemetryComponent(settings types.TelemetrySettings, name string) *TelemetryComponent {
	return &TelemetryComponent{
		logger:   settings.Logger,
		name:     name,
		settings: settings,
		resource: settings.Resource,
		ctype:    types.ComponentTypeService,
		healthy:  true,
		running:  false,
	}
}

// Start implements Component.Start
func (t *TelemetryComponent) Start(ctx context.Context) error {
	t.logger.Info("Starting telemetry component",
		zap.String("type", t.ctype.String()),
		zap.String("name", t.name))
	t.running = true
	return nil
}

// Shutdown implements Component.Shutdown
func (t *TelemetryComponent) Shutdown(ctx context.Context) error {
	t.logger.Info("Shutting down telemetry component",
		zap.String("type", t.ctype.String()),
		zap.String("name", t.name))

	// Stop running
	t.running = false

	// Shutdown meter provider if it implements Shutdownable
	if mp, ok := t.settings.MeterProvider.(interface{ Shutdown(context.Context) error }); ok {
		if err := mp.Shutdown(ctx); err != nil {
			t.logger.Error("Failed to shutdown meter provider", zap.Error(err))
		}
	}

	// Shutdown tracer provider if it implements Shutdownable
	if tp, ok := t.settings.TracerProvider.(interface{ Shutdown(context.Context) error }); ok {
		if err := tp.Shutdown(ctx); err != nil {
			t.logger.Error("Failed to shutdown tracer provider", zap.Error(err))
		}
	}

	return nil
}

// GetName returns the component name
func (t *TelemetryComponent) GetName() string {
	return t.name
}

// GetType returns the component type
func (t *TelemetryComponent) GetType() types.ComponentType {
	return t.ctype
}

// IsHealthy returns if the component is healthy
func (t *TelemetryComponent) IsHealthy() bool {
	return t.healthy
}

// IsRunning returns if the component is running
func (t *TelemetryComponent) IsRunning() bool {
	return t.running
}

// Tracer returns the component's tracer
func (t *TelemetryComponent) Tracer() trace.Tracer {
	return t.settings.Tracer
}

// Meter returns the component's meter
func (t *TelemetryComponent) Meter() metric.Meter {
	return t.settings.Meter
}

// Resource returns the component's resource
func (t *TelemetryComponent) Resource() *resource.Resource {
	return t.resource
}

// WithLogger returns a new TelemetryComponent with the given logger
func (t *TelemetryComponent) WithLogger(logger *zap.Logger) *TelemetryComponent {
	settings := t.settings
	settings.Logger = logger
	return NewTelemetryComponent(settings, t.GetName())
}

// WithTracer returns a new TelemetryComponent with the given tracer
func (t *TelemetryComponent) WithTracer(tracer trace.Tracer) *TelemetryComponent {
	settings := t.settings
	settings.Tracer = tracer
	return NewTelemetryComponent(settings, t.GetName())
}

// WithMeter returns a new TelemetryComponent with the given meter
func (t *TelemetryComponent) WithMeter(meter metric.Meter) *TelemetryComponent {
	settings := t.settings
	settings.Meter = meter
	return NewTelemetryComponent(settings, t.GetName())
}

// WithResource returns a new TelemetryComponent with the given resource
func (t *TelemetryComponent) WithResource(res *resource.Resource) *TelemetryComponent {
	settings := t.settings
	settings.Resource = res
	return NewTelemetryComponent(settings, t.GetName())
}

// WithInstrumentationScope returns a new TelemetryComponent with the given instrumentation scope
func (t *TelemetryComponent) WithInstrumentationScope(scope instrumentation.Scope) *TelemetryComponent {
	settings := t.settings
	settings.InstrumentationScope = scope
	return NewTelemetryComponent(settings, t.GetName())
}

// WithInstrumentationVersion returns a new TelemetryComponent with the given instrumentation version
func (t *TelemetryComponent) WithInstrumentationVersion(version string) *TelemetryComponent {
	settings := t.settings
	settings.InstrumentationVersion = version
	return NewTelemetryComponent(settings, t.GetName())
}

// GetMeterProvider returns the meter provider
func (t *TelemetryComponent) GetMeterProvider() metric.MeterProvider {
	return t.settings.MeterProvider
}

// GetTracerProvider returns the tracer provider
func (t *TelemetryComponent) GetTracerProvider() trace.TracerProvider {
	return t.settings.TracerProvider
}
