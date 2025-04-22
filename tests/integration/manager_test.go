package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/srediag/srediag/internal/config"
	"github.com/srediag/srediag/internal/plugins"
	"github.com/srediag/srediag/internal/telemetry"
)

func TestManagersIntegration(t *testing.T) {
	logger := zaptest.NewLogger(t)
	ctx := context.Background()

	// Telemetry manager configuration
	telConfig := config.TelemetryConfig{
		Enabled:     true,
		ServiceName: "integration-test",
		Environment: "test",
		Endpoint:    "localhost:4317",
		Sampling: config.SamplingConfig{
			Type: "always_on",
			Rate: 1.0,
		},
	}

	// Plugin manager configuration
	pluginConfig := config.PluginsConfig{
		Directory: "testdata/plugins",
		Enabled:   []string{},
	}

	t.Run("coordinated initialization and shutdown", func(t *testing.T) {
		// Create managers
		telManager, err := telemetry.NewManager(telConfig, "1.0.0", logger)
		require.NoError(t, err)

		pluginManager := plugins.NewManager(pluginConfig, logger)

		// Start managers in correct order
		err = telManager.Start(ctx)
		require.NoError(t, err)
		assert.True(t, telManager.IsRunning())

		err = pluginManager.Start(ctx)
		require.NoError(t, err)
		assert.True(t, pluginManager.IsRunning())

		// Small pause to ensure everything is initialized
		time.Sleep(100 * time.Millisecond)

		// Stop managers in reverse order
		err = pluginManager.Stop(ctx)
		require.NoError(t, err)
		assert.False(t, pluginManager.IsRunning())

		err = telManager.Stop(ctx)
		require.NoError(t, err)
		assert.False(t, telManager.IsRunning())
	})

	t.Run("telemetry with plugins", func(t *testing.T) {
		// Create managers
		telManager, err := telemetry.NewManager(telConfig, "1.0.0", logger)
		require.NoError(t, err)

		pluginManager := plugins.NewManager(pluginConfig, logger)

		// Start telemetry
		err = telManager.Start(ctx)
		require.NoError(t, err)

		// Get tracer and meter for plugin use
		tracer := telManager.Tracer("test-plugin")
		meter := telManager.Meter("test-plugin")

		assert.NotNil(t, tracer)
		assert.NotNil(t, meter)

		// Start plugins
		err = pluginManager.Start(ctx)
		require.NoError(t, err)

		// Stop everything
		err = pluginManager.Stop(ctx)
		require.NoError(t, err)

		err = telManager.Stop(ctx)
		require.NoError(t, err)
	})
}

func TestManagersErrorHandling(t *testing.T) {
	logger := zaptest.NewLogger(t)
	ctx := context.Background()

	// Invalid telemetry configuration
	invalidTelConfig := config.TelemetryConfig{
		Enabled:     true,
		ServiceName: "integration-test",
		Environment: "test",
		Endpoint:    "invalid:port", // Invalid endpoint
	}

	t.Run("telemetry initialization failure", func(t *testing.T) {
		telManager, err := telemetry.NewManager(invalidTelConfig, "1.0.0", logger)
		require.NoError(t, err)

		// Use a cancelled context to force an initialization error
		cancelledCtx, cancel := context.WithCancel(ctx)
		cancel() // Cancel the context immediately

		err = telManager.Start(cancelledCtx)
		assert.Error(t, err)
		assert.False(t, telManager.IsRunning())
	})

	t.Run("recovery after failure", func(t *testing.T) {
		telManager, err := telemetry.NewManager(invalidTelConfig, "1.0.0", logger)
		require.NoError(t, err)

		// First attempt with invalid configuration - use cancelled context
		cancelledCtx, cancel := context.WithCancel(ctx)
		cancel()
		err = telManager.Start(cancelledCtx)
		assert.Error(t, err)
		assert.False(t, telManager.IsRunning(), "Manager should not be running after a start failure")

		// Fix configuration
		validTelConfig := config.TelemetryConfig{
			Enabled:     true,
			ServiceName: "integration-test",
			Environment: "test",
			Endpoint:    "localhost:4317",
		}
		telManager.SetConfig(validTelConfig)

		// Second attempt with valid configuration - use background context
		err = telManager.Start(ctx)
		require.NoError(t, err)
		assert.True(t, telManager.IsRunning())

		// Cleanup
		err = telManager.Stop(ctx)
		require.NoError(t, err)
	})
}
