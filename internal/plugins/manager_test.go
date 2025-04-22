package plugins

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/srediag/srediag/internal/config"
)

func TestNewManager(t *testing.T) {
	logger := zaptest.NewLogger(t)

	tests := []struct {
		name    string
		config  config.PluginsConfig
		wantErr bool
	}{
		{
			name: "valid basic configuration",
			config: config.PluginsConfig{
				Directory: "testdata/plugins",
				Enabled:   []string{"plugin1", "plugin2"},
				Settings: map[string]map[string]interface{}{
					"plugin1": {"key": "value"},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewManager(tt.config, logger)
			assert.NotNil(t, manager)
			assert.Equal(t, tt.config, manager.config)
		})
	}
}

func TestManager_StartStop(t *testing.T) {
	logger := zaptest.NewLogger(t)
	ctx := context.Background()

	// Create temporary directory for test plugins
	tmpDir, err := os.MkdirTemp("", "plugin-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	config := config.PluginsConfig{
		Directory: tmpDir,
		Enabled:   []string{},
	}

	t.Run("start and stop without plugins", func(t *testing.T) {
		manager := NewManager(config, logger)

		err := manager.Start(ctx)
		require.NoError(t, err)
		assert.True(t, manager.IsRunning())

		defer func() {
			if err := manager.Stop(ctx); err != nil {
				t.Errorf("failed to stop manager: %v", err)
			}
			assert.False(t, manager.IsRunning())
		}()
	})
}

func TestManager_LoadPlugin(t *testing.T) {
	logger := zaptest.NewLogger(t)

	// Create temporary directory for test plugins
	tmpDir, err := os.MkdirTemp("", "plugin-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	config := config.PluginsConfig{
		Directory: tmpDir,
	}

	t.Run("load invalid plugin", func(t *testing.T) {
		manager := NewManager(config, logger)
		err := manager.LoadPlugin(filepath.Join(tmpDir, "invalid.so"))
		assert.Error(t, err)
	})
}

func TestManager_GetPlugin(t *testing.T) {
	logger := zaptest.NewLogger(t)

	config := config.PluginsConfig{
		Directory: "testdata/plugins",
	}

	manager := NewManager(config, logger)

	t.Run("get nonexistent plugin", func(t *testing.T) {
		plugin, exists := manager.GetPlugin("nonexistent")
		assert.False(t, exists)
		assert.Nil(t, plugin)
	})
}

func TestManager_ListPlugins(t *testing.T) {
	logger := zaptest.NewLogger(t)

	config := config.PluginsConfig{
		Directory: "testdata/plugins",
	}

	manager := NewManager(config, logger)

	t.Run("list empty plugins", func(t *testing.T) {
		plugins := manager.ListPlugins()
		assert.Empty(t, plugins)
	})
}

func TestManager_InvalidStates(t *testing.T) {
	logger := zaptest.NewLogger(t)
	ctx := context.Background()

	config := config.PluginsConfig{
		Directory: "testdata/plugins",
	}

	t.Run("try to stop manager that is not running", func(t *testing.T) {
		manager := NewManager(config, logger)
		err := manager.Stop(ctx)
		assert.Error(t, err)
	})

	t.Run("try to start manager that is already running", func(t *testing.T) {
		manager := NewManager(config, logger)

		err := manager.Start(ctx)
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
		manager := NewManager(config, logger)

		err := manager.Start(ctx)
		require.NoError(t, err)

		err = manager.Stop(ctx)
		require.NoError(t, err)

		err = manager.Stop(ctx)
		assert.Error(t, err)
	})
}

func TestManager_GetPluginConfig(t *testing.T) {
	logger := zaptest.NewLogger(t)

	config := config.PluginsConfig{
		Directory: "testdata/plugins",
		Settings: map[string]map[string]interface{}{
			"plugin1": {
				"key": "value",
			},
		},
	}

	manager := NewManager(config, logger)

	t.Run("get existing plugin config", func(t *testing.T) {
		cfg := manager.GetPluginConfig("plugin1")
		assert.NotNil(t, cfg)
		assert.Equal(t, "value", cfg["key"])
	})

	t.Run("get nonexistent plugin config", func(t *testing.T) {
		cfg := manager.GetPluginConfig("nonexistent")
		assert.Nil(t, cfg)
	})
}
