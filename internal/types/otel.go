package types

import (
	"go.opentelemetry.io/collector/component"
)

// OtelComponentType represents an OpenTelemetry component type
type OtelComponentType string

const (
	// OtelComponentTypeReceiver represents a receiver component
	OtelComponentTypeReceiver OtelComponentType = "receiver"
	// OtelComponentTypeProcessor represents a processor component
	OtelComponentTypeProcessor OtelComponentType = "processor"
	// OtelComponentTypeExporter represents an exporter component
	OtelComponentTypeExporter OtelComponentType = "exporter"
	// OtelComponentTypeExtension represents an extension component
	OtelComponentTypeExtension OtelComponentType = "extension"
)

// ToComponentType converts an OtelComponentType to a ComponentType
func (t OtelComponentType) ToComponentType() ComponentType {
	switch t {
	case OtelComponentTypeReceiver:
		return ComponentTypeReceiver
	case OtelComponentTypeProcessor:
		return ComponentTypeProcessor
	case OtelComponentTypeExporter:
		return ComponentTypeExporter
	case OtelComponentTypeExtension:
		return ComponentTypeExtension
	default:
		return ComponentTypeUnknown
	}
}

// FromComponentType converts a ComponentType to an OtelComponentType
func FromComponentType(t ComponentType) OtelComponentType {
	switch t {
	case ComponentTypeReceiver:
		return OtelComponentTypeReceiver
	case ComponentTypeProcessor:
		return OtelComponentTypeProcessor
	case ComponentTypeExporter:
		return OtelComponentTypeExporter
	case ComponentTypeExtension:
		return OtelComponentTypeExtension
	default:
		return ""
	}
}

// OtelComponent represents an OpenTelemetry component interface
type OtelComponent interface {
	component.Component
}

// OtelFactory represents an OpenTelemetry factory interface
type OtelFactory interface {
	Type() component.Type
	CreateDefaultConfig() interface{}
}

// OtelSettings represents OpenTelemetry component settings
type OtelSettings struct {
	// BuildInfo contains build information
	BuildInfo component.BuildInfo
	// TelemetrySettings contains telemetry settings
	TelemetrySettings component.TelemetrySettings
}

// NewOtelSettings creates new OpenTelemetry settings
func NewOtelSettings(buildInfo component.BuildInfo, telemetrySettings component.TelemetrySettings) OtelSettings {
	return OtelSettings{
		BuildInfo:         buildInfo,
		TelemetrySettings: telemetrySettings,
	}
}
