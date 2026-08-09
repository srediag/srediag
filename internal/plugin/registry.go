package plugin

import (
	"context"
	"fmt"
	"sync"

	"github.com/srediag/srediag/internal/core"
)

// PluginRegistry manages available and loaded plugins
type PluginRegistry struct {
	logger           *core.Logger
	availablePlugins map[string]*PluginMetadata
	loadedPlugins    map[string]*PluginInstance
	mu               sync.RWMutex
}

// NewPluginRegistry creates a new plugin registry
func NewPluginRegistry(logger *core.Logger) *PluginRegistry {
	return &PluginRegistry{
		logger:           logger,
		availablePlugins: make(map[string]*PluginMetadata),
		loadedPlugins:    make(map[string]*PluginInstance),
	}
}

// RegisterPlugin registers a plugin metadata
func (r *PluginRegistry) RegisterPlugin(metadata *PluginMetadata) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if metadata.Name == "" {
		return fmt.Errorf("plugin name cannot be empty")
	}

	r.availablePlugins[metadata.Name] = metadata
	r.logger.Info("Plugin registered",
		core.ZapString("name", metadata.Name),
		core.ZapString("version", metadata.Version),
		core.ZapReflect("type", metadata.Type))

	return nil
}

// LoadPlugin loads a plugin instance
func (r *PluginRegistry) LoadPlugin(name string, path string) (*PluginInstance, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	metadata, exists := r.availablePlugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found in registry", name)
	}

	if _, loaded := r.loadedPlugins[name]; loaded {
		return nil, fmt.Errorf("plugin %s is already loaded", name)
	}

	instance := &PluginInstance{
		Metadata: metadata,
		Path:     path,
		Status:   StatusLoaded,
	}

	r.loadedPlugins[name] = instance
	r.logger.Info("Plugin loaded",
		core.ZapString("name", name),
		core.ZapString("path", path))

	return instance, nil
}

// UnloadPlugin unloads a plugin instance
func (r *PluginRegistry) UnloadPlugin(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	instance, exists := r.loadedPlugins[name]
	if !exists {
		return fmt.Errorf("plugin %s is not loaded", name)
	}

	instance.Status = StatusStopped
	delete(r.loadedPlugins, name)

	r.logger.Info("Plugin unloaded", core.ZapString("name", name))
	return nil
}

// GetPlugin returns a loaded plugin instance
func (r *PluginRegistry) GetPlugin(name string) (*PluginInstance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	instance, exists := r.loadedPlugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin %s is not loaded", name)
	}

	return instance, nil
}

// ListAvailablePlugins returns all available plugins
func (r *PluginRegistry) ListAvailablePlugins() []*PluginMetadata {
	r.mu.RLock()
	defer r.mu.RUnlock()

	plugins := make([]*PluginMetadata, 0, len(r.availablePlugins))
	for _, metadata := range r.availablePlugins {
		plugins = append(plugins, metadata)
	}

	return plugins
}

// ListLoadedPlugins returns all loaded plugins
func (r *PluginRegistry) ListLoadedPlugins() []*PluginInstance {
	r.mu.RLock()
	defer r.mu.RUnlock()

	plugins := make([]*PluginInstance, 0, len(r.loadedPlugins))
	for _, instance := range r.loadedPlugins {
		plugins = append(plugins, instance)
	}

	return plugins
}

// StartAllPlugins starts all loaded plugins
func (r *PluginRegistry) StartAllPlugins(ctx context.Context) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for name, instance := range r.loadedPlugins {
		if instance.Status == StatusLoaded {
			instance.Status = StatusRunning
			r.logger.Info("Plugin started", core.ZapString("name", name))
		}
	}

	return nil
}

// StopAllPlugins stops all running plugins
func (r *PluginRegistry) StopAllPlugins() error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for name, instance := range r.loadedPlugins {
		if instance.Status == StatusRunning {
			instance.Status = StatusStopped
			r.logger.Info("Plugin stopped", core.ZapString("name", name))
		}
	}

	return nil
}

// GetPluginsByType returns plugins of a specific type
func (r *PluginRegistry) GetPluginsByType(pluginType PluginType) []*PluginMetadata {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var plugins []*PluginMetadata
	for _, metadata := range r.availablePlugins {
		if metadata.Type == pluginType {
			plugins = append(plugins, metadata)
		}
	}

	return plugins
}

// HealthCheck checks the health of all loaded plugins
func (r *PluginRegistry) HealthCheck() map[string]error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	results := make(map[string]error)
	for name, instance := range r.loadedPlugins {
		if instance.Status == StatusRunning {
			// In a real implementation, this would call the plugin's health check
			results[name] = nil
		}
	}

	return results
}

// GetPluginCount returns the count of available and loaded plugins
func (r *PluginRegistry) GetPluginCount() (available int, loaded int, running int) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	available = len(r.availablePlugins)
	loaded = len(r.loadedPlugins)

	for _, instance := range r.loadedPlugins {
		if instance.Status == StatusRunning {
			running++
		}
	}

	return available, loaded, running
}

// RegisterBuiltinPlugins registers built-in plugins
func (r *PluginRegistry) RegisterBuiltinPlugins() error {
	// Register the simple plugin
	simplePlugin := &PluginMetadata{
		Name:         "simple-plugin",
		Type:         ExamplePlugin,
		Version:      "1.0.0",
		Description:  "A simple example plugin for SREDIAG",
		Capabilities: []string{"logging", "metrics"},
		SHA256:       "", // Would be calculated for real plugins
		Signature:    "", // Would be provided for signed plugins
	}

	if err := r.RegisterPlugin(simplePlugin); err != nil {
		return fmt.Errorf("failed to register simple plugin: %w", err)
	}

	// Register debug receiver plugin
	debugReceiver := &PluginMetadata{
		Name:         "debug-receiver",
		Type:         ReceiverPlugin,
		Version:      "1.0.0",
		Description:  "Debug receiver that generates synthetic telemetry data",
		Capabilities: []string{"traces", "metrics", "logs"},
		SHA256:       "",
		Signature:    "",
	}

	if err := r.RegisterPlugin(debugReceiver); err != nil {
		return fmt.Errorf("failed to register debug receiver: %w", err)
	}

	r.logger.Info("Built-in plugins registered successfully",
		core.ZapInt("count", 2))

	return nil
}
