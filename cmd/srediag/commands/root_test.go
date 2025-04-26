package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/components"
	"github.com/srediag/srediag/internal/plugin"
	"github.com/srediag/srediag/internal/settings"
)

func setupTestSettings(t *testing.T) *settings.CommandSettings {
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	t.Cleanup(func() { _ = logger.Sync() })

	return &settings.CommandSettings{
		ComponentManager: components.NewManager(logger),
		PluginManager:    plugin.NewManager(logger, "testdata/plugins"),
		Logger:           logger,
	}
}

func TestNewRootCommand(t *testing.T) {
	t.Run("nil options", func(t *testing.T) {
		cmd := NewRootCommand(nil)
		assert.NotNil(t, cmd)
		assert.Equal(t, "srediag", cmd.Use)
		assert.NotEmpty(t, cmd.Short)
		assert.NotEmpty(t, cmd.Long)
	})

	t.Run("custom options", func(t *testing.T) {
		opts := &Options{
			LogConfig: LogConfig{
				Level:  "debug",
				Format: "json",
			},
		}
		cmd := NewRootCommand(opts)
		assert.NotNil(t, cmd)

		// Verify flags are set correctly
		level, _ := cmd.PersistentFlags().GetString("log-level")
		format, _ := cmd.PersistentFlags().GetString("log-format")
		assert.Equal(t, opts.LogConfig.Level, level)
		assert.Equal(t, opts.LogConfig.Format, format)
	})
}

func TestExecute(t *testing.T) {
	// Create a temporary config directory structure
	tmpDir, err := os.MkdirTemp("", "srediag-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create the configs directory
	configDir := filepath.Join(tmpDir, "configs")
	err = os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	// Create plugin directories
	pluginDir := filepath.Join(tmpDir, ".srediag", "plugins")
	pluginTmpDir := filepath.Join(pluginDir, ".tmp")
	pluginBinDir := filepath.Join(pluginDir, "bin")
	err = os.MkdirAll(pluginTmpDir, 0755)
	require.NoError(t, err)
	err = os.MkdirAll(pluginBinDir, 0755)
	require.NoError(t, err)

	// Create the config file
	configPath := filepath.Join(configDir, "srediag.yaml")
	configContent := `
plugins:
  dir: %s
  output_dir: %s
  builder_config: configs/otelcol-builder.yaml
`
	configContent = fmt.Sprintf(configContent, pluginTmpDir, pluginBinDir)
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// Create test settings
	settings := setupTestSettings(t)

	// Test without config file - should fail with config not found
	t.Run("without config file", func(t *testing.T) {
		// Clear environment
		os.Unsetenv("SREDIAG_CONFIG")
		err = Execute(settings)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "config file not found")
	})

	// Test with config file but no subcommand
	t.Run("with config file but no subcommand", func(t *testing.T) {
		t.Setenv("SREDIAG_CONFIG", configPath)
		err = Execute(settings)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "please specify a subcommand")
	})
}
