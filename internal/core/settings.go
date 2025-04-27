package core

import (
	"go.opentelemetry.io/collector/component"
)

// AppContext holds global app state, logger, component manager, build info, telemetry, and config.
type AppContext struct {
	Logger            *Logger
	ComponentManager  *ComponentManager
	BuildInfo         BuildInfo
	TelemetrySettings component.TelemetrySettings
	Config            *Config
}

// GetLogger returns the logger from the context, or a no-op logger if nil.
func (ctx *AppContext) GetLogger() *Logger {
	if ctx.Logger == nil {
		return &Logger{} // No-op logger
	}
	return ctx.Logger
}

// GetConfig returns the config from the context, or a new default config if nil.
func (ctx *AppContext) GetConfig() *Config {
	if ctx.Config == nil {
		return NewConfig()
	}
	return ctx.Config
}

// Deprecated: use AppContext instead.
type Settings struct {
	// BuildInfo contains build information
	BuildInfo component.BuildInfo
	// TelemetrySettings contains telemetry settings
	TelemetrySettings component.TelemetrySettings
}

// Deprecated: use AppContext instead.
type CommandSettings struct {
	ComponentManager *ComponentManager
	Logger           *Logger
}

// Deprecated: use AppContext.GetLogger instead.
func (s *CommandSettings) GetLogger() *Logger {
	if s.Logger == nil {
		return &Logger{} // No-op logger
	}
	return s.Logger
}
