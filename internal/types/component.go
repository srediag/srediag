package types

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// ComponentType represents the type of a component
type ComponentType string

const (
	// ComponentTypeCore represents a core component
	ComponentTypeCore ComponentType = "core"
	// ComponentTypePlugin represents a plugin component
	ComponentTypePlugin ComponentType = "plugin"
	// ComponentTypeExtension represents an extension component
	ComponentTypeExtension ComponentType = "extension"
	// ComponentTypeService represents a service component
	ComponentTypeService ComponentType = "service"
	// ComponentTypeReceiver represents a receiver component
	ComponentTypeReceiver ComponentType = "receiver"
	// ComponentTypeProcessor represents a processor component
	ComponentTypeProcessor ComponentType = "processor"
	// ComponentTypeExporter represents an exporter component
	ComponentTypeExporter ComponentType = "exporter"
	// ComponentTypeUnknown represents an unknown component type
	ComponentTypeUnknown ComponentType = "unknown"
)

// ComponentStatus represents the status of a component
type ComponentStatus string

const (
	// ComponentStatusUnknown represents an unknown status
	ComponentStatusUnknown ComponentStatus = "unknown"
	// ComponentStatusInitialized represents an initialized status
	ComponentStatusInitialized ComponentStatus = "initialized"
	// ComponentStatusRunning represents a running status
	ComponentStatusRunning ComponentStatus = "running"
	// ComponentStatusStopped represents a stopped status
	ComponentStatusStopped ComponentStatus = "stopped"
	// ComponentStatusError represents an error status
	ComponentStatusError ComponentStatus = "error"
)

// ComponentInfo represents component information
type ComponentInfo struct {
	Name    string        `json:"name"`
	Version string        `json:"version"`
	Type    ComponentType `json:"type"`
}

// ComponentSettings represents component configuration settings
type ComponentSettings map[string]interface{}

// GetString gets a string value from settings
func (s ComponentSettings) GetString(key string) string {
	if val, ok := s[key].(string); ok {
		return val
	}
	return ""
}

// GetInterface gets an interface value from settings
func (s ComponentSettings) GetInterface(key string) interface{} {
	return s[key]
}

// GetTelemetrySettings returns telemetry settings from settings
func (s ComponentSettings) GetTelemetrySettings() *TelemetrySettings {
	if val, ok := s["telemetry"].(*TelemetrySettings); ok {
		return val
	}
	return nil
}

// GetHost gets host from component settings
func (s ComponentSettings) GetHost() component.Host {
	if val, ok := s["host"].(component.Host); ok {
		return val
	}
	return nil
}

// TelemetrySettings holds telemetry-related settings
type TelemetrySettings struct {
	// Logger is the component's logger
	Logger *zap.Logger

	// Tracer is the default tracer instance
	Tracer trace.Tracer

	// Meter is the default meter instance
	Meter metric.Meter

	// TracerProvider is used to create tracers for spans
	TracerProvider trace.TracerProvider

	// MeterProvider is used to create meters for metrics
	MeterProvider metric.MeterProvider

	// Resource contains attributes identifying the service/component
	Resource *resource.Resource

	// InstrumentationVersion is the version of instrumentation
	InstrumentationVersion string

	// InstrumentationScope defines the scope of instrumentation
	InstrumentationScope instrumentation.Scope
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

// IComponent defines the interface that all components must implement
type IComponent interface {
	// GetName returns the component name
	GetName() string
	// GetVersion returns the component version
	GetVersion() string
	// GetType returns the component type
	GetType() ComponentType
	// GetStatus returns the component status
	GetStatus() ComponentStatus
	// IsHealthy returns the component health status
	IsHealthy() bool
	// Configure configures the component with settings
	Configure(settings ComponentSettings) error
	// Start starts the component
	Start(ctx context.Context) error
	// Stop stops the component
	Stop(ctx context.Context) error
}

// BaseComponent provides a base implementation of IComponent
type BaseComponent struct {
	name     string
	version  string
	cType    ComponentType
	status   ComponentStatus
	settings ComponentSettings
}

// NewBaseComponent creates a new BaseComponent
func NewBaseComponent(name, version string, cType ComponentType) *BaseComponent {
	return &BaseComponent{
		name:     name,
		version:  version,
		cType:    cType,
		status:   ComponentStatusUnknown,
		settings: make(ComponentSettings),
	}
}

// GetName returns the component name
func (c *BaseComponent) GetName() string {
	return c.name
}

// GetVersion returns the component version
func (c *BaseComponent) GetVersion() string {
	return c.version
}

// GetType returns the component type
func (c *BaseComponent) GetType() ComponentType {
	return c.cType
}

// GetStatus returns the component status
func (c *BaseComponent) GetStatus() ComponentStatus {
	return c.status
}

// IsHealthy returns true if the component is running without errors
func (c *BaseComponent) IsHealthy() bool {
	return c.status == ComponentStatusRunning
}

// Configure configures the component with settings
func (c *BaseComponent) Configure(settings ComponentSettings) error {
	c.settings = settings
	return nil
}

// Start starts the component
func (c *BaseComponent) Start(ctx context.Context) error {
	c.status = ComponentStatusRunning
	return nil
}

// Stop stops the component
func (c *BaseComponent) Stop(ctx context.Context) error {
	c.status = ComponentStatusStopped
	return nil
}

// String returns the string representation of the component type
func (t ComponentType) String() string {
	switch t {
	case ComponentTypeCore:
		return "core"
	case ComponentTypePlugin:
		return "plugin"
	case ComponentTypeExtension:
		return "extension"
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

// NewComponentSettings creates a new ComponentSettings instance
func NewComponentSettings(telemetry interface{}, host interface{}, name string, category string) ComponentSettings {
	return ComponentSettings{
		"telemetry": telemetry,
		"host":      host,
		"name":      name,
		"category":  category,
		"version":   "1.0.0", // This should be configurable
	}
}

// NewTelemetrySettings creates new TelemetrySettings
func NewTelemetrySettings(logger *zap.Logger, tracer trace.Tracer, meter metric.Meter, tracerProvider trace.TracerProvider, meterProvider metric.MeterProvider) TelemetrySettings {
	return TelemetrySettings{
		Logger:         logger,
		Tracer:         tracer,
		Meter:          meter,
		TracerProvider: tracerProvider,
		MeterProvider:  meterProvider,
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
	component.Component
	// GetName returns the name of the component
	GetName() string
	// GetCategory returns the category of the component
	GetCategory() PluginCategory
	// GetVersion returns the version of the component
	GetVersion() string
	// GetStatus returns the status of the component
	GetStatus() Status
	// Healthy returns true if the component is healthy
	Healthy() bool
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
