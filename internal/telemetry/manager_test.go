package telemetry

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/srediag/srediag/internal/config"
)

func TestNewManager(t *testing.T) {
	logger := zaptest.NewLogger(t)

	tests := []struct {
		name       string
		config     config.TelemetryConfig
		version    string
		wantErr    bool
		telEnabled bool
	}{
		{
			name: "telemetry disabled",
			config: config.TelemetryConfig{
				Enabled: false,
			},
			version:    "1.0.0",
			wantErr:    false,
			telEnabled: false,
		},
		{
			name: "telemetry enabled with valid config",
			config: config.TelemetryConfig{
				Enabled:     true,
				ServiceName: "test-service",
				Environment: "test",
				Endpoint:    "localhost:4317",
				Sampling: config.SamplingConfig{
					Type: "always_on",
					Rate: 1.0,
				},
			},
			version:    "1.0.0",
			wantErr:    false,
			telEnabled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, err := NewManager(tt.config, tt.version, logger)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, manager)
			assert.Equal(t, tt.telEnabled, tt.config.Enabled)
		})
	}
}

func TestManager_StartStop(t *testing.T) {
	logger := zaptest.NewLogger(t)
	ctx := context.Background()

	tests := []struct {
		name    string
		config  config.TelemetryConfig
		version string
		wantErr bool
	}{
		{
			name: "start and stop with telemetry disabled",
			config: config.TelemetryConfig{
				Enabled: false,
			},
			version: "1.0.0",
			wantErr: false,
		},
		{
			name: "start and stop with telemetry enabled",
			config: config.TelemetryConfig{
				Enabled:     true,
				ServiceName: "test-service",
				Environment: "test",
				Endpoint:    "localhost:4317",
				Sampling: config.SamplingConfig{
					Type: "always_on",
					Rate: 1.0,
				},
			},
			version: "1.0.0",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, err := NewManager(tt.config, tt.version, logger)
			require.NoError(t, err)

			// Test Start
			err = manager.Start(ctx)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.True(t, manager.IsRunning())

			// Small pause to ensure everything is initialized
			time.Sleep(100 * time.Millisecond)

			// Test Stop
			defer func() {
				if err := manager.Stop(ctx); err != nil {
					t.Errorf("failed to stop manager: %v", err)
				}
				// Assert running state after stopping
				assert.False(t, manager.IsRunning())
			}()
		})
	}
}

func TestManager_TracerMeter(t *testing.T) {
	logger := zaptest.NewLogger(t)

	tests := []struct {
		name    string
		config  config.TelemetryConfig
		version string
	}{
		{
			name: "get tracer and meter with telemetry disabled",
			config: config.TelemetryConfig{
				Enabled: false,
			},
			version: "1.0.0",
		},
		{
			name: "get tracer and meter with telemetry enabled",
			config: config.TelemetryConfig{
				Enabled:     true,
				ServiceName: "test-service",
				Environment: "test",
				Endpoint:    "localhost:4317",
				Sampling: config.SamplingConfig{
					Type: "always_on",
					Rate: 1.0,
				},
			},
			version: "1.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, err := NewManager(tt.config, tt.version, logger)
			require.NoError(t, err)

			if tt.config.Enabled {
				err = manager.Start(context.Background())
				require.NoError(t, err)
				defer func() {
					if err := manager.Stop(context.Background()); err != nil {
						t.Errorf("failed to stop manager: %v", err)
					}
				}()
			}

			// Test Tracer
			tracer := manager.Tracer("test")
			assert.NotNil(t, tracer)

			// Test Meter
			meter := manager.Meter("test")
			assert.NotNil(t, meter)
		})
	}
}

func TestManager_InvalidStates(t *testing.T) {
	logger := zaptest.NewLogger(t)
	ctx := context.Background()

	config := config.TelemetryConfig{
		Enabled:     true,
		ServiceName: "test-service",
		Environment: "test",
		Endpoint:    "localhost:4317",
	}

	t.Run("try to stop manager that is not running", func(t *testing.T) {
		manager, err := NewManager(config, "1.0.0", logger)
		require.NoError(t, err)

		err = manager.Stop(ctx)
		assert.Error(t, err)
	})

	t.Run("try to start manager that is already running", func(t *testing.T) {
		manager, err := NewManager(config, "1.0.0", logger)
		require.NoError(t, err)

		err = manager.Start(ctx)
		require.NoError(t, err)
		defer func() {
			if err := manager.Stop(ctx); err != nil {
				t.Errorf("failed to stop manager: %v", err)
			}
		}()

		err = manager.Start(ctx)
		assert.Error(t, err)
	})

	t.Run("try to stop manager that is already stopped", func(t *testing.T) {
		manager, err := NewManager(config, "1.0.0", logger)
		require.NoError(t, err)

		err = manager.Start(ctx)
		require.NoError(t, err)

		err = manager.Stop(ctx)
		require.NoError(t, err)

		err = manager.Stop(ctx)
		assert.Error(t, err)
	})
}
