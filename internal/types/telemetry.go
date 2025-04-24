// Package types provides telemetry-related interfaces for SREDIAG.
// This file contains interfaces for handling metrics, tracing, and resource monitoring
// using OpenTelemetry as the underlying telemetry framework.
package types

import (
	"context"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// IMeterProvider defines the interface for a meter provider
type IMeterProvider interface {
	// Meter creates a new Meter with the given name and version
	Meter(name string, opts ...metric.MeterOption) metric.Meter
	// MeterProvider returns the underlying meter provider
	MeterProvider() metric.MeterProvider
	// Shutdown shuts down the meter provider
	Shutdown(ctx context.Context) error
}

// ITracerProvider defines the interface for a tracer provider
type ITracerProvider interface {
	// Tracer creates a new Tracer with the given name and version
	Tracer(name string, opts ...trace.TracerOption) trace.Tracer
	// Shutdown shuts down the tracer provider
	Shutdown(ctx context.Context) error
}

// MeterProviderWrapper wraps a meter provider
type MeterProviderWrapper struct {
	provider metric.MeterProvider
}

// NewMeterProviderWrapper creates a new meter provider wrapper
func NewMeterProviderWrapper(provider metric.MeterProvider) *MeterProviderWrapper {
	return &MeterProviderWrapper{
		provider: provider,
	}
}

// Meter creates a new Meter with the given name and version
func (w *MeterProviderWrapper) Meter(name string, opts ...metric.MeterOption) metric.Meter {
	return w.provider.Meter(name, opts...)
}

// MeterProvider returns the underlying meter provider
func (w *MeterProviderWrapper) MeterProvider() metric.MeterProvider {
	return w.provider
}

// Shutdown shuts down the meter provider
func (w *MeterProviderWrapper) Shutdown(ctx context.Context) error {
	if provider, ok := w.provider.(interface{ Shutdown(context.Context) error }); ok {
		return provider.Shutdown(ctx)
	}
	return nil
}

// TracerProviderWrapper wraps a tracer provider
type TracerProviderWrapper struct {
	provider trace.TracerProvider
}

// NewTracerProviderWrapper creates a new tracer provider wrapper
func NewTracerProviderWrapper(provider trace.TracerProvider) *TracerProviderWrapper {
	return &TracerProviderWrapper{
		provider: provider,
	}
}

// Tracer creates a new Tracer with the given name and version
func (w *TracerProviderWrapper) Tracer(name string, opts ...trace.TracerOption) trace.Tracer {
	return w.provider.Tracer(name, opts...)
}

// Shutdown shuts down the tracer provider
func (w *TracerProviderWrapper) Shutdown(ctx context.Context) error {
	if provider, ok := w.provider.(interface{ Shutdown(context.Context) error }); ok {
		return provider.Shutdown(ctx)
	}
	return nil
}

// ITelemetryBridge defines the interface for a telemetry bridge
type ITelemetryBridge interface {
	// GetMeterProvider returns the meter provider
	GetMeterProvider() IMeterProvider
	// GetTracerProvider returns the tracer provider
	GetTracerProvider() ITracerProvider
	// SetResource sets the resource attributes
	SetResource(attributes map[string]string) error
	// GetResource returns the current resource attributes
	GetResource() map[string]string
	// Start starts the telemetry bridge
	Start(ctx context.Context) error
	// Stop stops the telemetry bridge
	Stop(ctx context.Context) error
	// Shutdown shuts down the telemetry bridge
	Shutdown(ctx context.Context) error
	// IsHealthy returns true if the telemetry bridge is healthy
	IsHealthy() bool
}

// IResourceMonitor defines the interface for a resource monitor
type IResourceMonitor interface {
	// GetMetrics returns the metrics for the resource
	GetMetrics() map[string]interface{}
	// GetThresholds returns the thresholds for the resource
	GetThresholds() map[string]interface{}
	// GetInterval returns the interval for the resource
	GetInterval() string
	// Start starts the resource monitor
	Start(ctx context.Context) error
	// Stop stops the resource monitor
	Stop(ctx context.Context) error
	// IsHealthy returns true if the resource monitor is healthy
	IsHealthy() bool
}

// ITelemetryComponent defines a simple interface for components with telemetry
type ITelemetryComponent interface {
	// GetName returns the name of the component
	GetName() string
	// GetType returns the type of the component
	GetType() string
}
