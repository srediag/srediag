package types

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.uber.org/zap"
)

// Settings represents component settings
type Settings struct {
	// Logger provides logging functionality
	Logger *zap.Logger
	// BuildInfo contains build information
	BuildInfo component.BuildInfo
	// Telemetry contains OpenTelemetry telemetry settings
	Telemetry component.TelemetrySettings
}

// NewSettings creates new component settings
func NewSettings(logger *zap.Logger, buildInfo component.BuildInfo, telemetry component.TelemetrySettings) Settings {
	return Settings{
		Logger:    logger,
		BuildInfo: buildInfo,
		Telemetry: telemetry,
	}
}

// ComponentConfig represents component configuration
type ComponentConfig struct {
	// ID is the component identifier
	ID component.ID
	// Settings contains component settings
	Settings Settings
	// Config contains component-specific configuration
	Config confmap.Conf
}

// NewComponentConfig creates a new component configuration
func NewComponentConfig(id component.ID, settings Settings, cfg confmap.Conf) ComponentConfig {
	return ComponentConfig{
		ID:       id,
		Settings: settings,
		Config:   cfg,
	}
}

// ToTelemetrySettings converts Settings to component.TelemetrySettings
func (s Settings) ToTelemetrySettings() component.TelemetrySettings {
	return s.Telemetry
}
