package base

import (
	"context"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/types"
)

// TelemetryComponent provides a base implementation of the TelemetryComponent interface
type TelemetryComponent struct {
	*BaseComponent
	settings types.TelemetrySettings
}

// NewTelemetryComponent creates a new telemetry component
func NewTelemetryComponent(settings types.TelemetrySettings, name string) *TelemetryComponent {
	return &TelemetryComponent{
		BaseComponent: NewBaseComponent(settings.Logger, types.ComponentTypeService, name),
		settings:      settings,
	}
}

// Start implements Component.Start
func (t *TelemetryComponent) Start(ctx context.Context) error {
	t.logger.Info("Starting telemetry component",
		zap.String("type", t.Type().String()),
		zap.String("name", t.Name()))
	return nil
}

// Shutdown implements Component.Shutdown
func (t *TelemetryComponent) Shutdown(ctx context.Context) error {
	t.logger.Info("Shutting down telemetry component",
		zap.String("type", t.Type().String()),
		zap.String("name", t.Name()))
	return nil
}

// Tracer returns the component's tracer
func (t *TelemetryComponent) Tracer() trace.Tracer {
	return t.settings.Tracer
}

// Meter returns the component's meter
func (t *TelemetryComponent) Meter() metric.Meter {
	return t.settings.Meter
}

// WithLogger returns a new TelemetryComponent with the given logger
func (t *TelemetryComponent) WithLogger(logger *zap.Logger) *TelemetryComponent {
	settings := t.settings
	settings.Logger = logger
	return NewTelemetryComponent(settings, t.Name())
}

// WithTracer returns a new TelemetryComponent with the given tracer
func (t *TelemetryComponent) WithTracer(tracer trace.Tracer) *TelemetryComponent {
	settings := t.settings
	settings.Tracer = tracer
	return NewTelemetryComponent(settings, t.Name())
}

// WithMeter returns a new TelemetryComponent with the given meter
func (t *TelemetryComponent) WithMeter(meter metric.Meter) *TelemetryComponent {
	settings := t.settings
	settings.Meter = meter
	return NewTelemetryComponent(settings, t.Name())
}
