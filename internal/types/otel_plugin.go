package types

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
)

// OtelPlugin represents an OpenTelemetry plugin interface
type OtelPlugin interface {
	component.Component
	// GetType returns the plugin type
	GetType() component.Type
	// GetKind returns the plugin kind
	GetKind() component.Kind
	// Capabilities returns the plugin capabilities
	Capabilities() consumer.Capabilities
}

// OtelPluginFactory represents an OpenTelemetry plugin factory interface
type OtelPluginFactory interface {
	// Type returns the plugin type
	Type() component.Type
	// CreateDefaultConfig creates the default configuration for the plugin
	CreateDefaultConfig() component.Config
	// CreateComponent creates a new instance of the plugin
	CreateComponent(ctx context.Context, cfg component.Config) (OtelPlugin, error)
}

// OtelPluginConfig represents the configuration for an OpenTelemetry plugin
type OtelPluginConfig struct {
	// Type is the plugin type
	Type component.Type `mapstructure:"type"`
	// Kind is the plugin kind
	Kind component.Kind `mapstructure:"kind"`
	// Config is the plugin-specific configuration
	Config component.Config `mapstructure:"config"`
}

// Validate validates the plugin configuration
func (c *OtelPluginConfig) Validate() error {
	if c.Type.String() == "" {
		return ErrInvalidPluginConfig("plugin type is required")
	}
	return nil
}

// OtelPluginOptions represents options for creating an OpenTelemetry plugin
type OtelPluginOptions struct {
	// Type is the plugin type
	Type component.Type
	// Kind is the plugin kind
	Kind component.Kind
	// Config is the plugin configuration
	Config component.Config
	// Settings contains component settings
	Settings OtelSettings
}

// NewOtelPluginOptions creates new OpenTelemetry plugin options
func NewOtelPluginOptions(typ component.Type, kind component.Kind, cfg component.Config, settings OtelSettings) OtelPluginOptions {
	return OtelPluginOptions{
		Type:     typ,
		Kind:     kind,
		Config:   cfg,
		Settings: settings,
	}
}
