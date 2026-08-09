package core

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// TelemetryConfig holds configuration for OpenTelemetry
type TelemetryConfig struct {
	ServiceName    string `yaml:"service_name"`
	ServiceVersion string `yaml:"service_version"`
	Environment    string `yaml:"environment"`

	// Tracing configuration
	TracingEnabled  bool    `yaml:"tracing_enabled"`
	TraceEndpoint   string  `yaml:"trace_endpoint"`
	TraceSampleRate float64 `yaml:"trace_sample_rate"`

	// Metrics configuration
	MetricsEnabled bool `yaml:"metrics_enabled"`
	MetricsPort    int  `yaml:"metrics_port"`

	// Resource attributes
	ResourceAttributes map[string]string `yaml:"resource_attributes"`
}

// TelemetryManager manages OpenTelemetry instrumentation for SREDIAG
type TelemetryManager struct {
	config *TelemetryConfig
	logger *Logger

	// OpenTelemetry providers
	tracerProvider *sdktrace.TracerProvider
	meterProvider  *sdkmetric.MeterProvider

	// Telemetry instruments
	tracer trace.Tracer
	meter  metric.Meter

	// Custom metrics
	requestCounter    metric.Int64Counter
	requestDuration   metric.Float64Histogram
	activeConnections metric.Int64UpDownCounter
	memoryUsage       metric.Int64ObservableGauge
	cpuUsage          metric.Float64ObservableGauge
}

// NewTelemetryManager creates a new telemetry manager
func NewTelemetryManager(config *TelemetryConfig, logger *Logger) (*TelemetryManager, error) {
	if config == nil {
		config = &TelemetryConfig{
			ServiceName:     "srediag",
			ServiceVersion:  "1.0.0",
			Environment:     "development",
			TracingEnabled:  true,
			TraceEndpoint:   "http://localhost:4318/v1/traces",
			TraceSampleRate: 1.0,
			MetricsEnabled:  true,
			MetricsPort:     8889,
		}
	}

	tm := &TelemetryManager{
		config: config,
		logger: logger,
	}

	if err := tm.initializeTracing(); err != nil {
		return nil, fmt.Errorf("failed to initialize tracing: %w", err)
	}

	if err := tm.initializeMetrics(); err != nil {
		return nil, fmt.Errorf("failed to initialize metrics: %w", err)
	}

	if err := tm.setupCustomMetrics(); err != nil {
		return nil, fmt.Errorf("failed to setup custom metrics: %w", err)
	}

	return tm, nil
}

// initializeTracing sets up OpenTelemetry tracing
func (tm *TelemetryManager) initializeTracing() error {
	if !tm.config.TracingEnabled {
		tm.logger.Info("Tracing disabled")
		return nil
	}

	// Create resource
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(tm.config.ServiceName),
		semconv.ServiceVersion(tm.config.ServiceVersion),
		semconv.DeploymentEnvironment(tm.config.Environment),
	)

	// Add custom resource attributes
	if tm.config.ResourceAttributes != nil {
		attrs := make([]attribute.KeyValue, 0, len(tm.config.ResourceAttributes))
		for k, v := range tm.config.ResourceAttributes {
			attrs = append(attrs, attribute.String(k, v))
		}
		customRes := resource.NewWithAttributes(semconv.SchemaURL, attrs...)
		mergedRes, err := resource.Merge(res, customRes)
		if err != nil {
			return fmt.Errorf("failed to add custom attributes: %w", err)
		}
		res = mergedRes
	}

	// Create OTLP trace exporter
	traceExporter, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpoint(tm.config.TraceEndpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Create trace provider
	tm.tracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(tm.config.TraceSampleRate)),
	)

	// Set global tracer provider
	otel.SetTracerProvider(tm.tracerProvider)

	// Set propagator
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Get tracer
	tm.tracer = tm.tracerProvider.Tracer(
		tm.config.ServiceName,
		trace.WithInstrumentationVersion(tm.config.ServiceVersion),
	)

	tm.logger.Info("Tracing initialized",
		zap.String("endpoint", tm.config.TraceEndpoint),
		zap.Float64("sample_rate", tm.config.TraceSampleRate))

	return nil
}

// initializeMetrics sets up OpenTelemetry metrics
func (tm *TelemetryManager) initializeMetrics() error {
	if !tm.config.MetricsEnabled {
		tm.logger.Info("Metrics disabled")
		return nil
	}

	// Create Prometheus exporter
	prometheusExporter, err := prometheus.New()
	if err != nil {
		return fmt.Errorf("failed to create Prometheus exporter: %w", err)
	}

	// Create meter provider
	tm.meterProvider = sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(prometheusExporter),
		sdkmetric.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(tm.config.ServiceName),
			semconv.ServiceVersion(tm.config.ServiceVersion),
		)),
	)

	// Set global meter provider
	otel.SetMeterProvider(tm.meterProvider)

	// Get meter
	tm.meter = tm.meterProvider.Meter(
		tm.config.ServiceName,
		metric.WithInstrumentationVersion(tm.config.ServiceVersion),
	)

	tm.logger.Info("Metrics initialized",
		zap.Int("port", tm.config.MetricsPort))

	return nil
}

// setupCustomMetrics creates custom metrics instruments
func (tm *TelemetryManager) setupCustomMetrics() error {
	if !tm.config.MetricsEnabled {
		return nil
	}

	var err error

	// Request counter
	tm.requestCounter, err = tm.meter.Int64Counter(
		"srediag_requests_total",
		metric.WithDescription("Total number of requests processed"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return fmt.Errorf("failed to create request counter: %w", err)
	}

	// Request duration histogram
	tm.requestDuration, err = tm.meter.Float64Histogram(
		"srediag_request_duration_seconds",
		metric.WithDescription("Request processing duration"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return fmt.Errorf("failed to create request duration histogram: %w", err)
	}

	// Active connections
	tm.activeConnections, err = tm.meter.Int64UpDownCounter(
		"srediag_active_connections",
		metric.WithDescription("Number of active connections"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return fmt.Errorf("failed to create active connections counter: %w", err)
	}

	// Memory usage gauge
	tm.memoryUsage, err = tm.meter.Int64ObservableGauge(
		"srediag_memory_usage_bytes",
		metric.WithDescription("Memory usage in bytes"),
		metric.WithUnit("By"),
	)
	if err != nil {
		return fmt.Errorf("failed to create memory usage gauge: %w", err)
	}

	// CPU usage gauge
	tm.cpuUsage, err = tm.meter.Float64ObservableGauge(
		"srediag_cpu_usage_percent",
		metric.WithDescription("CPU usage percentage"),
		metric.WithUnit("%"),
	)
	if err != nil {
		return fmt.Errorf("failed to create CPU usage gauge: %w", err)
	}

	tm.logger.Info("Custom metrics initialized")
	return nil
}

// StartSpan starts a new trace span
func (tm *TelemetryManager) StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	if tm.tracer == nil {
		return ctx, trace.SpanFromContext(ctx)
	}
	return tm.tracer.Start(ctx, name, opts...)
}

// RecordRequest records metrics for a request
func (tm *TelemetryManager) RecordRequest(ctx context.Context, method, endpoint string, status int, duration time.Duration) {
	if !tm.config.MetricsEnabled {
		return
	}

	labels := []attribute.KeyValue{
		attribute.String("method", method),
		attribute.String("endpoint", endpoint),
		attribute.Int("status", status),
	}

	// Increment request counter
	tm.requestCounter.Add(ctx, 1, metric.WithAttributes(labels...))

	// Record request duration
	tm.requestDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(labels...))
}

// IncrementActiveConnections increments the active connections counter
func (tm *TelemetryManager) IncrementActiveConnections(ctx context.Context, delta int64) {
	if tm.activeConnections != nil {
		tm.activeConnections.Add(ctx, delta)
	}
}

// AddSpanEvent adds an event to the current span
func (tm *TelemetryManager) AddSpanEvent(ctx context.Context, name string, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(name, trace.WithAttributes(attrs...))
}

// SetSpanAttributes sets attributes on the current span
func (tm *TelemetryManager) SetSpanAttributes(ctx context.Context, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attrs...)
}

// RecordError records an error on the current span
func (tm *TelemetryManager) RecordError(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

// GetTracer returns the tracer instance
func (tm *TelemetryManager) GetTracer() trace.Tracer {
	return tm.tracer
}

// GetMeter returns the meter instance
func (tm *TelemetryManager) GetMeter() metric.Meter {
	return tm.meter
}

// Shutdown gracefully shuts down the telemetry providers
func (tm *TelemetryManager) Shutdown(ctx context.Context) error {
	var errs []error

	if tm.tracerProvider != nil {
		if err := tm.tracerProvider.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to shutdown tracer provider: %w", err))
		}
	}

	if tm.meterProvider != nil {
		if err := tm.meterProvider.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to shutdown meter provider: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("shutdown errors: %v", errs)
	}

	tm.logger.Info("Telemetry shutdown completed")
	return nil
}

// InstrumentHTTPHandler wraps an HTTP handler with telemetry
func (tm *TelemetryManager) InstrumentHTTPHandler(handlerName string, handler func(ctx context.Context) error) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		start := time.Now()

		// Start span
		ctx, span := tm.StartSpan(ctx, handlerName)
		defer span.End()

		// Execute handler
		err := handler(ctx)

		// Record metrics and span status
		status := 200
		if err != nil {
			status = 500
			tm.RecordError(ctx, err)
		}

		tm.RecordRequest(ctx, "HTTP", handlerName, status, time.Since(start))

		return err
	}
}

// GetInstrumentation returns instrumentation information
func (tm *TelemetryManager) GetInstrumentation() instrumentation.Scope {
	return instrumentation.Scope{
		Name:    tm.config.ServiceName,
		Version: tm.config.ServiceVersion,
	}
}
