package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/srediag/srediag/internal/config"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		wantErr  bool
	}{
		{
			name:     "load valid configuration",
			filePath: filepath.Join("testdata", "config.yaml"),
			wantErr:  false,
		},
		{
			name:     "nonexistent file",
			filePath: "nonexistent.yaml",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := config.LoadConfig(tt.filePath)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, cfg)

			// Check service configuration
			assert.Equal(t, "srediag-test", cfg.Service.Name)
			assert.Equal(t, "test", cfg.Service.Environment)

			// Check telemetry configuration
			assert.True(t, cfg.Telemetry.Enabled)
			assert.Equal(t, "srediag-test", cfg.Telemetry.ServiceName)
			assert.Equal(t, "test", cfg.Telemetry.Environment)
			assert.Equal(t, "localhost:4317", cfg.Telemetry.Endpoint)
			assert.Equal(t, "always_on", cfg.Telemetry.Sampling.Type)
			assert.Equal(t, float64(1.0), cfg.Telemetry.Sampling.Rate)

			// Check plugins configuration
			assert.Equal(t, "testdata/plugins", cfg.Plugins.Directory)
			assert.Contains(t, cfg.Plugins.Enabled, "mock-plugin")

			mockPluginCfg := cfg.Plugins.Settings["mock-plugin"]
			assert.NotNil(t, mockPluginCfg)
			assert.Equal(t, "value", mockPluginCfg["key"])
			assert.Equal(t, true, mockPluginCfg["enabled"])

			// Check logging configuration
			assert.Equal(t, "debug", cfg.Logging.Level)
			assert.Equal(t, "console", cfg.Logging.Format)
			assert.Equal(t, "stdout", cfg.Logging.Output)
		})
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  config.Config
		wantErr bool
	}{
		{
			name: "valid configuration",
			config: config.Config{
				Service: config.ServiceConfig{
					Name:        "test",
					Environment: "test",
				},
				Telemetry: config.TelemetryConfig{
					Enabled:     true,
					ServiceName: "test",
					Environment: "test",
					Endpoint:    "localhost:4317",
				},
				Plugins: config.PluginsConfig{
					Directory: "plugins",
				},
				Logging: config.LoggingConfig{
					Level:  "debug",
					Format: "console",
					Output: "stdout",
				},
			},
			wantErr: false,
		},
		{
			name: "telemetry enabled without endpoint",
			config: config.Config{
				Telemetry: config.TelemetryConfig{
					Enabled: true,
				},
			},
			wantErr: true,
		},
		{
			name: "empty plugins directory",
			config: config.Config{
				Plugins: config.PluginsConfig{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// If the config expects a plugin directory, create it for the test
			if tt.config.Plugins.Directory != "" {
				if err := os.MkdirAll(tt.config.Plugins.Directory, 0750); err != nil {
					t.Fatalf("failed to create temp plugin dir %s: %v", tt.config.Plugins.Directory, err)
				}
				// Clean up the directory after the test run
				defer os.RemoveAll(tt.config.Plugins.Directory)
			}

			err := tt.config.Validate()

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
