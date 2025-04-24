package factory

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"

	"github.com/srediag/srediag/internal/types"
)

// ExtensionFactory provides a base implementation for extension factories
type ExtensionFactory struct {
	*BaseFactory
	createExtension func(context.Context, component.TelemetrySettings, component.Config) (component.Component, error)
}

// NewExtensionFactory creates a new extension factory
func NewExtensionFactory(
	typ component.Type,
	version string,
	configType component.Config,
	capabilities consumer.Capabilities,
) *ExtensionFactory {
	return &ExtensionFactory{
		BaseFactory: NewBaseFactory(typ, version, configType, capabilities),
	}
}

// CreateExtension implements types.ExtensionFactory
func (f *ExtensionFactory) CreateExtension(
	ctx context.Context,
	set component.TelemetrySettings,
	cfg component.Config,
) (component.Component, error) {
	if f.createExtension == nil {
		return nil, types.ErrInvalidPluginConfig("CreateExtension not implemented")
	}
	return f.createExtension(ctx, set, cfg)
}

// WithExtension sets the extension creation function
func (f *ExtensionFactory) WithExtension(
	createFn func(context.Context, component.TelemetrySettings, component.Config) (component.Component, error),
) *ExtensionFactory {
	f.createExtension = createFn
	return f
}
