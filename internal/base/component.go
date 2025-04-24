package base

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/types"
)

// BaseComponent extends OpenTelemetry's base component with SREDIAG-specific functionality
type BaseComponent struct {
	component.StartFunc
	component.ShutdownFunc
	logger    *zap.Logger
	name      string
	category  types.PluginCategory
	version   string
	status    types.Status
	host      component.Host
	telemetry component.TelemetrySettings
}

// NewBaseComponent creates a new base component
func NewBaseComponent(settings types.ComponentSettings) *BaseComponent {
	telemetrySettings, ok := settings.GetInterface("telemetry").(component.TelemetrySettings)
	if !ok {
		// Fallback to empty telemetry settings if not provided
		telemetrySettings = component.TelemetrySettings{}
	}

	host, ok := settings.GetInterface("host").(component.Host)
	if !ok {
		// Fallback to nil host if not provided
		host = nil
	}

	return &BaseComponent{
		logger:    telemetrySettings.Logger,
		name:      settings.GetString("name"),
		category:  types.PluginCategory(settings.GetString("category")),
		version:   settings.GetString("version"),
		status:    types.StatusLoaded,
		host:      host,
		telemetry: telemetrySettings,
	}
}

// Start implements component.Component
func (b *BaseComponent) Start(ctx context.Context) error {
	b.logger.Info("Starting component",
		zap.String("name", b.name),
		zap.String("category", string(b.category)))
	b.status = types.StatusRunning
	return nil
}

// Shutdown implements component.Component
func (b *BaseComponent) Shutdown(ctx context.Context) error {
	b.logger.Info("Shutting down component",
		zap.String("name", b.name),
		zap.String("category", string(b.category)))
	b.status = types.StatusStopped
	return nil
}

// GetName returns the component's name
func (b *BaseComponent) GetName() string {
	return b.name
}

// GetCategory returns the component's category
func (b *BaseComponent) GetCategory() types.PluginCategory {
	return b.category
}

// GetVersion returns the component's version
func (b *BaseComponent) GetVersion() string {
	return b.version
}

// GetStatus returns the component's status
func (b *BaseComponent) GetStatus() types.Status {
	return b.status
}

// GetHost returns the component's host
func (b *BaseComponent) GetHost() component.Host {
	return b.host
}

// GetTelemetrySettings returns the component's telemetry settings
func (b *BaseComponent) GetTelemetrySettings() component.TelemetrySettings {
	return b.telemetry
}

// Logger returns the component's logger
func (b *BaseComponent) Logger() *zap.Logger {
	return b.logger
}

// SetStatus sets the component's status
func (b *BaseComponent) SetStatus(status types.Status) {
	b.status = status
}

// Healthy returns true if the component is healthy
func (b *BaseComponent) Healthy() bool {
	return b.status == types.StatusRunning
}

// ConfigurableComponent provides a base implementation of the ConfigurableComponent interface
type ConfigurableComponent struct {
	*BaseComponent
	settings types.ComponentSettings
	config   *confmap.Conf
}

// NewConfigurableComponent creates a new configurable component
func NewConfigurableComponent(settings types.ComponentSettings) *ConfigurableComponent {
	return &ConfigurableComponent{
		BaseComponent: NewBaseComponent(settings),
		settings:      settings,
		config:        nil,
	}
}

// Configure implements ConfigurableComponent.Configure
func (c *ConfigurableComponent) Configure(cfg *confmap.Conf) error {
	c.logger.Info("Configuring component",
		zap.String("name", c.GetName()),
		zap.String("category", string(c.GetCategory())))
	c.config = cfg
	return nil
}

// GetConfig returns the component's configuration
func (c *ConfigurableComponent) GetConfig() *confmap.Conf {
	return c.config
}

// GetSettings returns the component's settings
func (c *ConfigurableComponent) GetSettings() types.ComponentSettings {
	return c.settings
}
