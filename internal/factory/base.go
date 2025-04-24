package factory

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"

	"github.com/srediag/srediag/internal/types"
)

// BaseFactory provides a base implementation of ComponentFactory
type BaseFactory struct {
	typ          component.Type
	version      string
	configType   component.Config
	capabilities consumer.Capabilities
}

// NewBaseFactory creates a new BaseFactory
func NewBaseFactory(typ component.Type, version string, configType component.Config, capabilities consumer.Capabilities) *BaseFactory {
	return &BaseFactory{
		typ:          typ,
		version:      version,
		configType:   configType,
		capabilities: capabilities,
	}
}

// Type implements ComponentFactory
func (f *BaseFactory) Type() component.Type {
	return f.typ
}

// CreateDefaultConfig implements ComponentFactory
func (f *BaseFactory) CreateDefaultConfig() component.Config {
	return f.configType
}

// Capabilities implements ComponentFactory
func (f *BaseFactory) Capabilities() consumer.Capabilities {
	return f.capabilities
}

// ValidateConfig implements ComponentFactory
func (f *BaseFactory) ValidateConfig(cfg component.Config) error {
	if cfg == nil {
		return types.ErrInvalidPluginConfig("config cannot be nil")
	}
	return nil
}

// CreateComponent implements ComponentFactory
func (f *BaseFactory) CreateComponent(ctx context.Context, set component.TelemetrySettings, cfg component.Config) (component.Component, error) {
	return nil, types.ErrInvalidPluginConfig("CreateComponent not implemented")
}

// Version returns the factory version
func (f *BaseFactory) Version() string {
	return f.version
}

// WithCapabilities sets the factory capabilities
func (f *BaseFactory) WithCapabilities(capabilities consumer.Capabilities) *BaseFactory {
	f.capabilities = capabilities
	return f
}

// WithConfigType sets the factory config type
func (f *BaseFactory) WithConfigType(configType component.Config) *BaseFactory {
	f.configType = configType
	return f
}
