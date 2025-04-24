package plugin

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/types"
)

// OtelBase provides a base implementation of types.OtelPlugin interface
type OtelBase struct {
	*Base
	typ      component.Type
	kind     component.Kind
	settings types.OtelSettings
}

// NewOtelBase creates a new OpenTelemetry plugin base
func NewOtelBase(logger *zap.Logger, opts types.OtelPluginOptions) *OtelBase {
	metadata := types.PluginMetadata{
		ID:          opts.Type.String(),
		Name:        opts.Type.String(),
		Version:     "1.0.0",
		Category:    types.PluginCategory(opts.Kind.String()),
		Author:      "SREDIAG",
		Description: "OpenTelemetry plugin",
	}

	return &OtelBase{
		Base:     NewBase(logger, nil, metadata),
		typ:      opts.Type,
		kind:     opts.Kind,
		settings: opts.Settings,
	}
}

// GetType implements types.OtelPlugin
func (b *OtelBase) GetType() component.Type {
	return b.typ
}

// GetKind implements types.OtelPlugin
func (b *OtelBase) GetKind() component.Kind {
	return b.kind
}

// Capabilities implements types.OtelPlugin
func (b *OtelBase) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{
		MutatesData: false,
	}
}

// Start implements component.Component
func (b *OtelBase) Start(_ context.Context, _ component.Host) error {
	return b.Base.Start(context.Background())
}

// Shutdown implements component.Component
func (b *OtelBase) Shutdown(_ context.Context) error {
	return b.Stop(context.Background())
}

// OtelPluginFactory provides a base implementation of types.OtelPluginFactory
type OtelPluginFactory struct {
	typ           component.Type
	kind          component.Kind
	createFunc    func(context.Context, types.OtelPluginOptions) (types.OtelPlugin, error)
	defaultConfig component.Config
	settings      types.OtelSettings
}

// NewOtelPluginFactory creates a new OpenTelemetry plugin factory
func NewOtelPluginFactory(typ component.Type, kind component.Kind, defaultConfig component.Config) *OtelPluginFactory {
	return &OtelPluginFactory{
		typ:           typ,
		kind:          kind,
		defaultConfig: defaultConfig,
	}
}

// Type implements types.OtelPluginFactory
func (f *OtelPluginFactory) Type() component.Type {
	return f.typ
}

// CreateDefaultConfig implements types.OtelPluginFactory
func (f *OtelPluginFactory) CreateDefaultConfig() component.Config {
	return f.defaultConfig
}

// CreateComponent implements types.OtelPluginFactory
func (f *OtelPluginFactory) CreateComponent(ctx context.Context, cfg component.Config) (types.OtelPlugin, error) {
	if f.createFunc == nil {
		return nil, types.ErrInvalidPluginConfig("create function not set")
	}

	opts := types.OtelPluginOptions{
		Type:     f.typ,
		Kind:     f.kind,
		Config:   cfg,
		Settings: f.settings,
	}

	return f.createFunc(ctx, opts)
}

// WithCreateFunc sets the component creation function
func (f *OtelPluginFactory) WithCreateFunc(fn func(context.Context, types.OtelPluginOptions) (types.OtelPlugin, error)) *OtelPluginFactory {
	f.createFunc = fn
	return f
}

// WithSettings sets the component settings
func (f *OtelPluginFactory) WithSettings(settings types.OtelSettings) *OtelPluginFactory {
	f.settings = settings
	return f
}
