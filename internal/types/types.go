package types

import (
	"context"

	"go.opentelemetry.io/collector/component"
)

// Type represents a component type
type Type string

const (
	// TypeCore represents a core component
	TypeCore Type = "core"
	// TypePlugin represents a plugin component
	TypePlugin Type = "plugin"
	// TypeReceiver represents a receiver component
	TypeReceiver Type = "receiver"
	// TypeProcessor represents a processor component
	TypeProcessor Type = "processor"
	// TypeExporter represents an exporter component
	TypeExporter Type = "exporter"
	// TypeExtension represents an extension component
	TypeExtension Type = "extension"
	// TypeConnector represents a connector component
	TypeConnector Type = "connector"
)

// Factory defines the interface for component factories
type Factory interface {
	// Type returns the type of component created by this factory
	Type() Type
	// CreateDefaultConfig creates the default configuration for the component
	CreateDefaultConfig() interface{}
}

// Component defines the interface that all components must implement
type Component interface {
	// Start starts the component
	Start(ctx context.Context) error
	// Stop stops the component
	Stop(ctx context.Context) error
}

// Registry defines the interface for component registry
type Registry interface {
	// RegisterFactory registers a new factory
	RegisterFactory(factory Factory) error
	// GetFactory returns a factory by type
	GetFactory(typ Type) (Factory, bool)
	// GetFactories returns all registered factories
	GetFactories() map[Type]Factory
}

// Settings holds settings for components
type Settings struct {
	// BuildInfo contains build information
	BuildInfo component.BuildInfo
	// TelemetrySettings contains telemetry settings
	TelemetrySettings component.TelemetrySettings
}

// NewSettings creates new Settings
func NewSettings(buildInfo component.BuildInfo, telemetrySettings component.TelemetrySettings) Settings {
	return Settings{
		BuildInfo:         buildInfo,
		TelemetrySettings: telemetrySettings,
	}
}
