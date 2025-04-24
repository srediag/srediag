package base

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/types"
)

// FactoryComponent provides a base implementation of the FactoryComponent interface
type FactoryComponent struct {
	*BaseComponent
	settings types.FactorySettings
}

// NewFactoryComponent creates a new factory component
func NewFactoryComponent(settings types.FactorySettings, name string) *FactoryComponent {
	return &FactoryComponent{
		BaseComponent: NewBaseComponent(settings.Logger, types.ComponentTypeService, name),
		settings:      settings,
	}
}

// Start implements Component.Start
func (f *FactoryComponent) Start(ctx context.Context) error {
	f.logger.Info("Starting factory component",
		zap.String("type", f.Type().String()),
		zap.String("name", f.Name()))
	return nil
}

// Shutdown implements Component.Shutdown
func (f *FactoryComponent) Shutdown(ctx context.Context) error {
	f.logger.Info("Shutting down factory component",
		zap.String("type", f.Type().String()),
		zap.String("name", f.Name()))
	return nil
}

// WithLogger returns a new FactoryComponent with the given logger
func (f *FactoryComponent) WithLogger(logger *zap.Logger) *FactoryComponent {
	settings := f.settings
	settings.Logger = logger
	return NewFactoryComponent(settings, f.Name())
}

// CreateDefaultConfig implements FactoryComponent.CreateDefaultConfig
func (f *FactoryComponent) CreateDefaultConfig() component.Config {
	return f.settings.DefaultConfig
}

// WithHost returns a new FactoryComponent with the given host
func (f *FactoryComponent) WithHost(host component.Host) *FactoryComponent {
	settings := f.settings
	settings.Host = host
	return NewFactoryComponent(settings, f.Name())
}

// Host returns the component's host
func (f *FactoryComponent) Host() component.Host {
	return f.settings.Host
}
