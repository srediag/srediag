package plugin

import (
	"context"
	"testing"

	"github.com/srediag/srediag/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPluginRegistry(t *testing.T) {
	logger, err := core.NewLogger(&core.Logger{
		Level:  "info",
		Format: "console",
	})
	require.NoError(t, err)

	registry := NewPluginRegistry(logger)
	assert.NotNil(t, registry)

	// Test registering built-in plugins
	err = registry.RegisterBuiltinPlugins()
	require.NoError(t, err)

	// Test listing available plugins
	available := registry.ListAvailablePlugins()
	assert.Len(t, available, 2)

	// Test getting plugin counts
	availableCount, loadedCount, runningCount := registry.GetPluginCount()
	assert.Equal(t, 2, availableCount)
	assert.Equal(t, 0, loadedCount)
	assert.Equal(t, 0, runningCount)

	// Test loading a plugin
	instance, err := registry.LoadPlugin("simple-plugin", "/fake/path/simple-plugin")
	require.NoError(t, err)
	assert.NotNil(t, instance)
	assert.Equal(t, "simple-plugin", instance.Metadata.Name)
	assert.Equal(t, StatusLoaded, instance.Status)

	// Test getting plugin counts after loading
	availableCount, loadedCount, runningCount = registry.GetPluginCount()
	assert.Equal(t, 2, availableCount)
	assert.Equal(t, 1, loadedCount)
	assert.Equal(t, 0, runningCount)

	// Test starting plugins
	err = registry.StartAllPlugins(context.Background())
	require.NoError(t, err)

	// Test getting plugin counts after starting
	availableCount, loadedCount, runningCount = registry.GetPluginCount()
	assert.Equal(t, 2, availableCount)
	assert.Equal(t, 1, loadedCount)
	assert.Equal(t, 1, runningCount)

	// Test health check
	health := registry.HealthCheck()
	assert.Len(t, health, 1)
	assert.Contains(t, health, "simple-plugin")
	assert.NoError(t, health["simple-plugin"])

	// Test getting plugins by type
	examplePlugins := registry.GetPluginsByType(ExamplePlugin)
	assert.Len(t, examplePlugins, 1)
	assert.Equal(t, "simple-plugin", examplePlugins[0].Name)

	receiverPlugins := registry.GetPluginsByType(ReceiverPlugin)
	assert.Len(t, receiverPlugins, 1)
	assert.Equal(t, "debug-receiver", receiverPlugins[0].Name)

	// Test stopping plugins
	err = registry.StopAllPlugins()
	require.NoError(t, err)

	// Test unloading plugin
	err = registry.UnloadPlugin("simple-plugin")
	require.NoError(t, err)

	// Test getting plugin counts after unloading
	availableCount, loadedCount, runningCount = registry.GetPluginCount()
	assert.Equal(t, 2, availableCount)
	assert.Equal(t, 0, loadedCount)
	assert.Equal(t, 0, runningCount)
}

func TestPluginRegistryErrors(t *testing.T) {
	logger, err := core.NewLogger(&core.Logger{
		Level:  "info",
		Format: "console",
	})
	require.NoError(t, err)

	registry := NewPluginRegistry(logger)

	// Test registering plugin with empty name
	err = registry.RegisterPlugin(&PluginMetadata{
		Name: "",
		Type: ExamplePlugin,
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "plugin name cannot be empty")

	// Test loading non-existent plugin
	_, err = registry.LoadPlugin("non-existent", "/path")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found in registry")

	// Test loading plugin twice
	registry.RegisterPlugin(&PluginMetadata{
		Name: "test-plugin",
		Type: ExamplePlugin,
	})

	_, err = registry.LoadPlugin("test-plugin", "/path1")
	require.NoError(t, err)

	_, err = registry.LoadPlugin("test-plugin", "/path2")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already loaded")

	// Test unloading non-loaded plugin
	err = registry.UnloadPlugin("non-loaded")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "is not loaded")

	// Test getting non-existent plugin
	_, err = registry.GetPlugin("non-existent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "is not loaded")
}
