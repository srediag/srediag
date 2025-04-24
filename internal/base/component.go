package base

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/types"
)

// BaseComponent provides a base implementation of the Component interface
type BaseComponent struct {
	info   types.ComponentInfo
	logger *zap.Logger
}

// NewBaseComponent creates a new base component
func NewBaseComponent(logger *zap.Logger, typ interface{}, name string) *BaseComponent {
	var componentType types.ComponentType
	switch t := typ.(type) {
	case types.ComponentType:
		componentType = t
	case component.Type:
		// Map OpenTelemetry component types to our component types
		componentType = types.ComponentTypePlugin
	default:
		componentType = types.ComponentTypeUnknown
	}

	return &BaseComponent{
		info:   types.NewComponentInfo(componentType, name),
		logger: logger,
	}
}

// Start implements Component.Start
func (b *BaseComponent) Start(ctx context.Context) error {
	b.logger.Info("Starting component",
		zap.String("type", b.info.Type.String()),
		zap.String("name", b.info.Name))
	return nil
}

// Shutdown implements Component.Shutdown
func (b *BaseComponent) Shutdown(ctx context.Context) error {
	b.logger.Info("Shutting down component",
		zap.String("type", b.info.Type.String()),
		zap.String("name", b.info.Name))
	return nil
}

// Logger returns the component's logger
func (b *BaseComponent) Logger() *zap.Logger {
	return b.logger
}

// Type returns the component's type
func (b *BaseComponent) Type() types.ComponentType {
	return b.info.Type
}

// Name returns the component's name
func (b *BaseComponent) Name() string {
	return b.info.Name
}

// ConfigurableComponent provides a base implementation of the ConfigurableComponent interface
type ConfigurableComponent struct {
	*BaseComponent
	settings types.ComponentSettings
}

// NewConfigurableComponent creates a new configurable component
func NewConfigurableComponent(settings types.ComponentSettings, typ interface{}, name string) *ConfigurableComponent {
	return &ConfigurableComponent{
		BaseComponent: NewBaseComponent(settings.Logger, typ, name),
		settings:      settings,
	}
}

// Configure implements ConfigurableComponent.Configure
func (c *ConfigurableComponent) Configure(cfg *confmap.Conf) error {
	c.logger.Info("Configuring component",
		zap.String("type", c.Type().String()),
		zap.String("name", c.Name()))
	return nil
}

// GetConfig returns the component's configuration
func (c *ConfigurableComponent) GetConfig() *confmap.Conf {
	return nil
}
