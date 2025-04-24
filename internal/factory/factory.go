package factory

import (
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/types"
)

// Factory represents a base component factory that can be extended for specific component types
type Factory struct {
	logger        *zap.Logger
	id            string
	componentType types.ComponentType
	config        interface{}
	createFunc    func(settings component.TelemetrySettings, cfg interface{}) (component.Component, error)
	verifyFunc    func(cfg interface{}) error
}

// FactoryOption represents an option for configuring a factory
type FactoryOption func(*Factory)

// WithCreateFunc sets the component creation function
func WithCreateFunc(fn func(settings component.TelemetrySettings, cfg interface{}) (component.Component, error)) FactoryOption {
	return func(f *Factory) {
		f.createFunc = fn
	}
}

// WithVerifyFunc sets the configuration verification function
func WithVerifyFunc(fn func(cfg interface{}) error) FactoryOption {
	return func(f *Factory) {
		f.verifyFunc = fn
	}
}

// WithLogger sets the factory logger
func WithLogger(logger *zap.Logger) FactoryOption {
	return func(f *Factory) {
		f.logger = logger
	}
}

// WithID sets the factory ID
func WithID(id string) FactoryOption {
	return func(f *Factory) {
		f.id = id
	}
}

// WithComponentType sets the component type
func WithComponentType(componentType types.ComponentType) FactoryOption {
	return func(f *Factory) {
		f.componentType = componentType
	}
}

// WithConfig sets the factory configuration
func WithConfig(config interface{}) FactoryOption {
	return func(f *Factory) {
		f.config = config
	}
}

// NewFactory creates a new factory with the given options
func NewFactory(opts ...FactoryOption) *Factory {
	f := &Factory{}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

// GetID returns the factory ID
func (f *Factory) GetID() string {
	return f.id
}

// GetType returns the factory type
func (f *Factory) GetType() types.ComponentType {
	return f.componentType
}

// CreateDefaultConfig creates the default configuration for components
func (f *Factory) CreateDefaultConfig() interface{} {
	return f.config
}

// VerifyConfig verifies the configuration
func (f *Factory) VerifyConfig(cfg interface{}) error {
	if f.verifyFunc != nil {
		return f.verifyFunc(cfg)
	}
	return nil
}

// CreateComponent creates a component instance
func (f *Factory) CreateComponent(settings component.TelemetrySettings, cfg interface{}) (component.Component, error) {
	if f.createFunc == nil {
		return nil, fmt.Errorf("create function not set for factory %s", f.id)
	}

	if err := f.VerifyConfig(cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration for factory %s: %w", f.id, err)
	}

	return f.createFunc(settings, cfg)
}

// UnmarshalConfig unmarshals configuration using confmap
func (f *Factory) UnmarshalConfig(conf *confmap.Conf, target interface{}) error {
	if err := conf.Unmarshal(target); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return nil
}

// CreateSettings creates component.TelemetrySettings with the factory's logger
func (f *Factory) CreateSettings(buildInfo component.BuildInfo) component.TelemetrySettings {
	return component.TelemetrySettings{
		Logger: f.logger,
	}
}

// Type returns the factory type
func (f *Factory) Type() types.ComponentType {
	return f.componentType
}

// Config returns the factory configuration
func (f *Factory) Config() interface{} {
	return f.config
}
