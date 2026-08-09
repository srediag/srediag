package plugin

import (
	"context"

	"go.opentelemetry.io/collector/component"

	"github.com/cloudwego/shmipc-go"
)

// PluginComponentFactory implements component.Factory for plugin-based components
type PluginComponentFactory struct {
	metadata PluginMetadata
	ch       *shmipc.SessionManager
}

// Type returns the component type
func (f *PluginComponentFactory) Type() component.Type {
	typ, _ := component.NewType(string(f.metadata.Type) + "/" + f.metadata.Name)
	return typ
}

// CreateDefaultConfig creates a default configuration for the component
func (f *PluginComponentFactory) CreateDefaultConfig() component.Config {
	return &PluginConfig{
		Name:    f.metadata.Name,
		Type:    string(f.metadata.Type),
		Enabled: true,
	}
}

// CreateComponent creates a generic component (placeholder for Phase 2)
func (f *PluginComponentFactory) CreateComponent(
	ctx context.Context,
	cfg component.Config,
) (component.Component, error) {
	return &PluginComponent{
		metadata: f.metadata,
		ch:       f.ch,
		config:   cfg,
	}, nil
}

// PluginComponent implements component.Component for plugin-based components
type PluginComponent struct {
	metadata PluginMetadata
	ch       *shmipc.SessionManager
	config   component.Config
}

// Start starts the plugin component
func (c *PluginComponent) Start(ctx context.Context, host component.Host) error {
	// Implementation would send Start IPC message to plugin
	return nil
}

// Shutdown shuts down the plugin component
func (c *PluginComponent) Shutdown(ctx context.Context) error {
	// Implementation would send Shutdown IPC message to plugin
	return nil
}
