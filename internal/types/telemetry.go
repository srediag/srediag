// Package types provides telemetry-related interfaces for SREDIAG.
// This file contains interfaces for handling metrics, tracing, and resource monitoring
// using OpenTelemetry as the underlying telemetry framework.
package types

import (
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
)

// IMeterProvider defines an interface for creating metric instruments.
// It abstracts the OpenTelemetry meter provider to allow for different
// implementations and easier testing.
type IMeterProvider interface {
	// Meter returns a new meter instance for creating metric instruments.
	// The name parameter is used to identify the meter's source in the telemetry data.
	// Options can be used to configure the meter's behavior.
	Meter(name string, opts ...metric.MeterOption) metric.Meter
}

// ITracerProvider defines an interface for creating tracers.
// It abstracts the OpenTelemetry tracer provider to allow for different
// implementations and easier testing.
type ITracerProvider interface {
	// Tracer returns a new tracer instance for creating spans.
	// The name parameter is used to identify the tracer's source in the telemetry data.
	// Options can be used to configure the tracer's behavior.
	Tracer(name string, opts ...trace.TracerOption) trace.Tracer
}

// ITelemetryBridge serves as the main interface for telemetry operations.
// It provides access to metrics and tracing capabilities while managing
// the underlying OpenTelemetry resource attributes.
type ITelemetryBridge interface {
	// IComponent embeds the base component interface
	IComponent

	// SetResource updates the OpenTelemetry resource with new attributes.
	// The resource contains attributes that are common to all telemetry data
	// emitted by this component (e.g., service name, version, environment).
	SetResource(res *resource.Resource)

	// GetResource returns the current OpenTelemetry resource.
	// This resource contains the attributes that identify the source of
	// telemetry data.
	GetResource() *resource.Resource

	// GetMeterProvider returns the meter provider for creating metrics.
	// The meter provider is used to create new meters for recording
	// measurements and observations.
	GetMeterProvider() IMeterProvider

	// GetTracerProvider returns the tracer provider for distributed tracing.
	// The tracer provider is used to create new tracers for tracking
	// request flows and dependencies.
	GetTracerProvider() ITracerProvider
}

// IResourceMonitor provides an interface for monitoring system resources.
// It tracks resource usage and health metrics for the local system or
// specific components.
type IResourceMonitor interface {
	// IComponent embeds the base component interface
	IComponent

	// GetMetrics returns the current resource metrics as key-value pairs.
	// Keys are metric names (e.g., "cpu.usage", "memory.used")
	// Values are the corresponding measurements.
	GetMetrics() map[string]float64
}
