// Package plugin provides a plugin management system for OpenTelemetry components
package plugin

import (
	"context"

	"go.opentelemetry.io/collector/component"
)

// Instance represents a plugin instance and its lifecycle.
type IPluginInstance interface {
	// Initialize sets up the plugin with provided metadata
	Initialize(ctx context.Context, metadata PluginMetadata) error
	// Start begins plugin operation
	Start(ctx context.Context) error
	// Stop gracefully shuts down the plugin
	Stop(ctx context.Context) error
	// HealthCheck returns the current plugin status
	HealthCheck(ctx context.Context) (*PluginHealth, error)
	// Factory returns the OpenTelemetry component factory
	Factory() (component.Factory, error)
}
