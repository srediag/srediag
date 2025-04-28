// Package core provides foundational types and utilities for the SREDIAG system.
//
// This file defines the AppContext type, which holds global application state, logger, config, and component manager for SREDIAG.
// AppContext is used for dependency injection and lifecycle management across all major subsystems.
//
// Usage:
//   - Use AppContext to pass global state, logger, config, and managers to all major subsystems and CLI entrypoints.
//   - Use GetLogger and GetConfig for safe access to logger and config with sensible fallbacks.
//
// Best Practices:
//   - Always prefer AppContext over legacy Settings/CommandSettings.
//   - Avoid storing mutable global state outside AppContext.
//
// TODO:
//   - Remove legacy Settings and CommandSettings after full migration.
//   - Add context cancellation and shutdown hooks if needed.
//
// Redundancy/Refactor:
//   - AppContext supersedes Settings and CommandSettings; those are deprecated.
package core

import (
	"go.opentelemetry.io/collector/component"
)

// AppContext holds global app state, logger, component manager, build info, telemetry, and config.
//
// Usage:
//   - Pass AppContext to all CLI entrypoints and service initializations.
//   - Use for dependency injection and lifecycle management.
//
// Fields:
//   - Logger: Main logger for the application.
//   - ComponentManager: Manages component factories and lifecycles.
//   - BuildInfo: Build/version metadata.
//   - TelemetrySettings: OpenTelemetry collector settings.
//   - Config: Loaded configuration for the application.
type AppContext struct {
	Logger            *Logger
	ComponentManager  *ComponentManager
	BuildInfo         BuildInfo
	TelemetrySettings component.TelemetrySettings
	Config            *Config
}

// GetLogger returns the logger from the context, or a no-op logger if nil.
//
// Usage:
//   - Use to safely access the logger in any subsystem.
//   - Returns a no-op logger if Logger is nil.
func (ctx *AppContext) GetLogger() *Logger {
	if ctx.Logger == nil {
		return &Logger{} // No-op logger
	}
	return ctx.Logger
}

// GetConfig returns the config from the context, or a new default config if nil.
//
// Usage:
//   - Use to safely access the config in any subsystem.
//   - Returns a new default config if Config is nil.
func (ctx *AppContext) GetConfig() *Config {
	if ctx.Config == nil {
		return NewConfig()
	}
	return ctx.Config
}

// Settings is deprecated. Use AppContext instead.
//
// Usage:
//   - Legacy struct for build and telemetry settings.
//   - Do not use in new code.
//
// Redundancy/Refactor:
//   - Superseded by AppContext.
type Settings struct {
	// BuildInfo contains build information
	BuildInfo component.BuildInfo
	// TelemetrySettings contains telemetry settings
	TelemetrySettings component.TelemetrySettings
}

// CommandSettings is deprecated. Use AppContext instead.
//
// Usage:
//   - Legacy struct for component manager and logger.
//   - Do not use in new code.
//
// Redundancy/Refactor:
//   - Superseded by AppContext.
type CommandSettings struct {
	ComponentManager *ComponentManager
	Logger           *Logger
}

// GetLogger returns the logger from CommandSettings, or a no-op logger if nil.
//
// Usage:
//   - Legacy method for safe logger access.
//   - Do not use in new code; prefer AppContext.GetLogger.
func (s *CommandSettings) GetLogger() *Logger {
	if s.Logger == nil {
		return &Logger{} // No-op logger
	}
	return s.Logger
}
