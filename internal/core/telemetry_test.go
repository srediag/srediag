package core

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
)

func TestTelemetryManagerCreation(t *testing.T) {
	logger, err := NewLogger(&Logger{
		Level:  "info",
		Format: "console",
	})
	require.NoError(t, err)

	config := &TelemetryConfig{
		ServiceName:    "test-service",
		ServiceVersion: "1.0.0",
		Environment:    "test",
		TracingEnabled: false, // Disable to avoid external dependencies
		MetricsEnabled: false,
	}

	tm, err := NewTelemetryManager(config, logger)
	assert.NoError(t, err)
	assert.NotNil(t, tm)
}

func TestTelemetryManagerWithNilConfig(t *testing.T) {
	logger, err := NewLogger(&Logger{
		Level:  "info",
		Format: "console",
	})
	require.NoError(t, err)

	tm, err := NewTelemetryManager(nil, logger)
	assert.NoError(t, err)
	assert.NotNil(t, tm)
	assert.Equal(t, "srediag", tm.config.ServiceName)
}

func TestStartSpanWithoutTracing(t *testing.T) {
	logger, err := NewLogger(&Logger{
		Level:  "info",
		Format: "console",
	})
	require.NoError(t, err)

	config := &TelemetryConfig{
		TracingEnabled: false,
		MetricsEnabled: false,
	}

	tm, err := NewTelemetryManager(config, logger)
	require.NoError(t, err)

	ctx := context.Background()
	newCtx, span := tm.StartSpan(ctx, "test-span")

	assert.NotNil(t, newCtx)
	assert.NotNil(t, span)

	span.End()
}

func TestRecordRequestWithoutMetrics(t *testing.T) {
	logger, err := NewLogger(&Logger{
		Level:  "info",
		Format: "console",
	})
	require.NoError(t, err)

	config := &TelemetryConfig{
		TracingEnabled: false,
		MetricsEnabled: false,
	}

	tm, err := NewTelemetryManager(config, logger)
	require.NoError(t, err)

	ctx := context.Background()

	// This should not panic when metrics are disabled
	tm.RecordRequest(ctx, "GET", "/test", 200, 100*time.Millisecond)
	tm.IncrementActiveConnections(ctx, 1)
	tm.AddSpanEvent(ctx, "test-event", attribute.String("key", "value"))
	tm.SetSpanAttributes(ctx, attribute.String("key", "value"))
}

func TestInstrumentHTTPHandler(t *testing.T) {
	logger, err := NewLogger(&Logger{
		Level:  "info",
		Format: "console",
	})
	require.NoError(t, err)

	config := &TelemetryConfig{
		TracingEnabled: false,
		MetricsEnabled: false,
	}

	tm, err := NewTelemetryManager(config, logger)
	require.NoError(t, err)

	// Test successful handler
	successHandler := tm.InstrumentHTTPHandler("test-handler", func(ctx context.Context) error {
		return nil
	})

	ctx := context.Background()
	err = successHandler(ctx)
	assert.NoError(t, err)

	// Test error handler
	errorHandler := tm.InstrumentHTTPHandler("error-handler", func(ctx context.Context) error {
		return assert.AnError
	})

	err = errorHandler(ctx)
	assert.Error(t, err)
}

func TestTelemetryShutdown(t *testing.T) {
	logger, err := NewLogger(&Logger{
		Level:  "info",
		Format: "console",
	})
	require.NoError(t, err)

	config := &TelemetryConfig{
		TracingEnabled: false,
		MetricsEnabled: false,
	}

	tm, err := NewTelemetryManager(config, logger)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = tm.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestGetInstrumentation(t *testing.T) {
	logger, err := NewLogger(&Logger{
		Level:  "info",
		Format: "console",
	})
	require.NoError(t, err)

	config := &TelemetryConfig{
		ServiceName:    "test-service",
		ServiceVersion: "1.2.3",
		TracingEnabled: false,
		MetricsEnabled: false,
	}

	tm, err := NewTelemetryManager(config, logger)
	require.NoError(t, err)

	instrumentation := tm.GetInstrumentation()
	assert.Equal(t, "test-service", instrumentation.Name)
	assert.Equal(t, "1.2.3", instrumentation.Version)
}
