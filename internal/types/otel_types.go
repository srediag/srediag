package types

import (
	"go.opentelemetry.io/collector/component"
)

// OtelDataType represents OpenTelemetry data types
type OtelDataType string

const (
	// OtelDataTypeTraces represents trace data
	OtelDataTypeTraces OtelDataType = "traces"
	// OtelDataTypeMetrics represents metric data
	OtelDataTypeMetrics OtelDataType = "metrics"
	// OtelDataTypeLogs represents log data
	OtelDataTypeLogs OtelDataType = "logs"
	// OtelDataTypeNone represents no data type
	OtelDataTypeNone OtelDataType = ""
)

// OtelComponentKind represents OpenTelemetry component kinds
type OtelComponentKind string

const (
	// OtelComponentKindReceiver represents a receiver component
	OtelComponentKindReceiver OtelComponentKind = "receiver"
	// OtelComponentKindProcessor represents a processor component
	OtelComponentKindProcessor OtelComponentKind = "processor"
	// OtelComponentKindExporter represents an exporter component
	OtelComponentKindExporter OtelComponentKind = "exporter"
	// OtelComponentKindExtension represents an extension component
	OtelComponentKindExtension OtelComponentKind = "extension"
)

// ToComponentType converts an OtelDataType to a ComponentType
func (t OtelDataType) ToComponentType() ComponentType {
	switch t {
	case OtelDataTypeTraces:
		return ComponentTypeReceiver
	case OtelDataTypeMetrics:
		return ComponentTypeProcessor
	case OtelDataTypeLogs:
		return ComponentTypeExporter
	default:
		return ComponentTypeUnknown
	}
}

// ComponentTypeToDataType converts a ComponentType to an OtelDataType
func ComponentTypeToDataType(t ComponentType) OtelDataType {
	switch t {
	case ComponentTypeReceiver:
		return OtelDataTypeTraces
	case ComponentTypeProcessor:
		return OtelDataTypeMetrics
	case ComponentTypeExporter:
		return OtelDataTypeLogs
	default:
		return OtelDataTypeNone
	}
}

// ToOtelComponentType converts an OtelDataType to a component.Type
func (t OtelDataType) ToOtelComponentType() component.Type {
	typ, _ := component.NewType(string(t))
	return typ
}

// OtelComponentTypeToDataType converts a component.Type to an OtelDataType
func OtelComponentTypeToDataType(t component.Type) OtelDataType {
	return OtelDataType(t.String())
}

// ToOtelKind converts an OtelComponentKind to a component.Kind
func (k OtelComponentKind) ToOtelKind() component.Kind {
	switch k {
	case OtelComponentKindReceiver:
		return component.KindReceiver
	case OtelComponentKindProcessor:
		return component.KindProcessor
	case OtelComponentKindExporter:
		return component.KindExporter
	case OtelComponentKindExtension:
		return component.KindExtension
	default:
		return component.KindReceiver // Default to receiver as there's no Unknown kind
	}
}

// FromOtelKind converts a component.Kind to an OtelComponentKind
func FromOtelKind(k component.Kind) OtelComponentKind {
	switch k {
	case component.KindReceiver:
		return OtelComponentKindReceiver
	case component.KindProcessor:
		return OtelComponentKindProcessor
	case component.KindExporter:
		return OtelComponentKindExporter
	case component.KindExtension:
		return OtelComponentKindExtension
	default:
		return ""
	}
}
