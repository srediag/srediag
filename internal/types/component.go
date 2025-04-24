package types

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// ComponentType represents the type of a component
type ComponentType int

// ComponentStatus represents the current state of a component
type ComponentStatus int

// ComponentInfo represents component information
type ComponentInfo struct {
	Name    string        `json:"name"`
	Version string        `json:"version"`
	Type    ComponentType `json:"type"`
}

// ComponentSettings holds common component settings
type ComponentSettings struct {
	Logger *zap.Logger
	Name   string
	Type   ComponentType
}

// TelemetrySettings holds telemetry-related settings
type TelemetrySettings struct {
	Logger *zap.Logger
	Tracer trace.Tracer
	Meter  metric.Meter
}

// FactorySettings holds factory-related settings
type FactorySettings struct {
	Logger        *zap.Logger
	Host          component.Host
	DefaultConfig component.Config
}

// ComponentTelemetrySettings holds telemetry-related settings for components
type ComponentTelemetrySettings struct {
	Logger *zap.Logger
	Tracer trace.Tracer
	Meter  metric.Meter
}

// IComponent represents a base component interface
type IComponent interface {
	// GetName returns the name of the component
	GetName() string
	// GetVersion returns the version of the component
	GetVersion() string
	// GetType returns the type of the component
	GetType() ComponentType
	// Start starts the component
	Start(ctx context.Context) error
	// Stop stops the component
	Stop(ctx context.Context) error
	// IsHealthy returns true if the component is healthy
	IsHealthy() bool
	// Configure configures the component
	Configure(cfg interface{}) error
}

const (
	// ComponentStatusUnknown indicates the component status is not known
	ComponentStatusUnknown ComponentStatus = iota
	// ComponentStatusStarting indicates the component is starting up
	ComponentStatusStarting
	// ComponentStatusRunning indicates the component is running normally
	ComponentStatusRunning
	// ComponentStatusStopping indicates the component is shutting down
	ComponentStatusStopping
	// ComponentStatusStopped indicates the component has stopped
	ComponentStatusStopped
	// ComponentStatusError indicates the component encountered an error
	ComponentStatusError
)

const (
	// ComponentTypeUnknown represents an unknown component type
	ComponentTypeUnknown ComponentType = iota
	// ComponentTypeCore represents a core component type
	ComponentTypeCore
	// ComponentTypeService represents a service component
	ComponentTypeService
	// ComponentTypeReceiver represents a receiver component
	ComponentTypeReceiver
	// ComponentTypeProcessor represents a processor component
	ComponentTypeProcessor
	// ComponentTypeExporter represents an exporter component
	ComponentTypeExporter
	// ComponentTypeExtension represents an extension component
	ComponentTypeExtension
	// ComponentTypePlugin represents a plugin component
	ComponentTypePlugin
)

// String returns the string representation of the component type
func (t ComponentType) String() string {
	switch t {
	case ComponentTypeCore:
		return "core"
	case ComponentTypeService:
		return "service"
	case ComponentTypeReceiver:
		return "receiver"
	case ComponentTypeProcessor:
		return "processor"
	case ComponentTypeExporter:
		return "exporter"
	case ComponentTypeExtension:
		return "extension"
	case ComponentTypePlugin:
		return "plugin"
	default:
		return "unknown"
	}
}

// ToOtelType converts a ComponentType to an OpenTelemetry component.Type
func (t ComponentType) ToOtelType() component.Type {
	dataType := ComponentTypeToDataType(t)
	return dataType.ToOtelComponentType()
}

// FromOtelType converts an OpenTelemetry component.Type to a ComponentType
func FromOtelType(t component.Type) ComponentType {
	dataType := OtelComponentTypeToDataType(t)
	return dataType.ToComponentType()
}

// NewComponentInfo creates a new ComponentInfo instance
func NewComponentInfo(typ ComponentType, name string) ComponentInfo {
	return ComponentInfo{
		Name: name,
		Type: typ,
	}
}

// NewComponentSettings creates new ComponentSettings
func NewComponentSettings(logger *zap.Logger, name string, typ ComponentType) ComponentSettings {
	return ComponentSettings{
		Logger: logger,
		Name:   name,
		Type:   typ,
	}
}

// NewTelemetrySettings creates new TelemetrySettings
func NewTelemetrySettings(logger *zap.Logger, tracer trace.Tracer, meter metric.Meter) TelemetrySettings {
	return TelemetrySettings{
		Logger: logger,
		Tracer: tracer,
		Meter:  meter,
	}
}

// NewFactorySettings creates new FactorySettings
func NewFactorySettings(logger *zap.Logger, host component.Host, defaultConfig component.Config) FactorySettings {
	return FactorySettings{
		Logger:        logger,
		Host:          host,
		DefaultConfig: defaultConfig,
	}
}

// NewComponentTelemetrySettings creates new ComponentTelemetrySettings
func NewComponentTelemetrySettings(logger *zap.Logger, tracer trace.Tracer, meter metric.Meter) ComponentTelemetrySettings {
	return ComponentTelemetrySettings{
		Logger: logger,
		Tracer: tracer,
		Meter:  meter,
	}
}

// Component represents a base component interface
type Component interface {
	// Start starts the component
	Start(ctx context.Context) error
	// Shutdown stops the component
	Shutdown(ctx context.Context) error
}

// ConfigurableComponent represents a component that can be configured
type ConfigurableComponent interface {
	Component
	// Configure configures the component with the given configuration
	Configure(cfg *confmap.Conf) error
}

// TelemetryComponent represents a component that provides telemetry
type TelemetryComponent interface {
	Component
	// Logger returns the component's logger
	Logger() *zap.Logger
	// Tracer returns the component's tracer
	Tracer() trace.Tracer
	// Meter returns the component's meter
	Meter() metric.Meter
}

// FactoryComponent represents a component that can create other components
type FactoryComponent interface {
	Component
	// Type returns the component type
	Type() component.Type
	// CreateDefaultConfig creates a default configuration for the component
	CreateDefaultConfig() component.Config
}

// HostComponent represents a component that can host other components
type HostComponent interface {
	Component
	// GetComponent returns a component by its ID
	GetComponent(id component.ID) (Component, bool)
	// GetComponents returns all components of a given type
	GetComponents(typ component.Type) []Component
}
