package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"

	"github.com/srediag/srediag/internal/components"
	"github.com/srediag/srediag/internal/plugin"
	"github.com/srediag/srediag/internal/settings"
)

func setupPluginDir(t *testing.T, logger *zap.Logger) (string, *plugin.Manager) {
	// Create temporary directory for plugins
	tmpDir, err := os.MkdirTemp("", "srediag-test-*")
	require.NoError(t, err)
	t.Cleanup(func() { os.RemoveAll(tmpDir) })

	// Create plugin manager
	pm := plugin.NewManager(logger, filepath.Join(tmpDir, "bin"))
	return tmpDir, pm
}

func TestInitializeComponents(t *testing.T) {
	// Create a test logger
	logger := zaptest.NewLogger(t)

	// Setup plugin directory and manager
	_, pm := setupPluginDir(t, logger)

	// Create test settings
	settings := &settings.CommandSettings{
		ComponentManager: components.NewManager(logger),
		PluginManager:    pm,
		Logger:           logger,
	}

	// Test component initialization
	err := initializeComponents(context.Background(), settings)
	assert.Error(t, err) // Expected to fail since no plugins exist
	assert.Contains(t, err.Error(), "plugin not found")

	// Verify managers are still valid
	assert.NotNil(t, settings.ComponentManager)
	assert.NotNil(t, settings.PluginManager)
}

func TestLoadCoreComponents(t *testing.T) {
	// Create a test logger
	logger := zaptest.NewLogger(t)

	// Setup plugin directory and manager
	tmpDir, pm := setupPluginDir(t, logger)

	// Test loading core components
	err := loadCoreComponents(context.Background(), pm)
	assert.Error(t, err) // Expected to fail since no plugins exist
	assert.Contains(t, err.Error(), "plugin not found")

	// Create mock plugin files (empty files are sufficient for the test)
	pluginTypes := []string{"receivers", "processors", "exporters", "extensions"}
	for _, typ := range pluginTypes {
		dir := filepath.Join(tmpDir, "bin", typ)
		require.NoError(t, os.MkdirAll(dir, 0755))

		// Create empty plugin files
		pluginPath := filepath.Join(dir, "otlp")
		require.NoError(t, os.WriteFile(pluginPath, []byte{}, 0644))
	}

	// Test loading core components again
	err = loadCoreComponents(context.Background(), pm)
	assert.Error(t, err) // Expected to fail since plugins are empty
	assert.Contains(t, err.Error(), "failed to create plugin")
}

func TestRegisterPluginComponents(t *testing.T) {
	// Create a test logger
	logger := zaptest.NewLogger(t)

	// Setup plugin directory and manager
	_, pm := setupPluginDir(t, logger)

	// Create component manager
	cm := components.NewManager(logger)

	// Test registering components
	err := registerPluginComponents(context.Background(), pm, cm)
	assert.NoError(t, err) // Should succeed since GetFactory returns nil for non-existent plugins
}

func TestMain(m *testing.M) {
	// Setup test environment
	logger, err := zap.NewDevelopment()
	if err != nil {
		os.Exit(1)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to sync logger: %v\n", err)
		}
	}()

	// Run tests
	os.Exit(m.Run())
}
