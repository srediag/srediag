// Package plugin provides a plugin management system for OpenTelemetry components.
//
// This file defines the IPluginInstance interface for plugin lifecycle management and integration with OpenTelemetry Collector components.
//
// Usage:
//   - Implement IPluginInstance to represent a plugin instance and its lifecycle.
//   - Use the interface methods to initialize, start, stop, check health, and obtain the component factory for a plugin.
//
// Best Practices:
//   - Always check for errors from lifecycle methods.
//   - Use context.Context for all operations to support cancellation and timeouts.
//   - Document all interface methods with expected side effects and error handling.
package plugin

import (
	"context"

	"go.opentelemetry.io/collector/component"
)

// IPluginInstance represents a plugin instance and its lifecycle.
//
// Implement this interface to manage the lifecycle of a plugin and integrate with OpenTelemetry Collector components.
type IPluginInstance interface {
	// Initialize sets up the plugin with provided metadata.
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeouts.
	//   - metadata: PluginMetadata containing plugin identity and configuration.
	//
	// Returns:
	//   - error: If initialization fails, returns a detailed error.
	Initialize(ctx context.Context, metadata PluginMetadata) error
	// Start begins plugin operation.
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeouts.
	//
	// Returns:
	//   - error: If starting fails, returns a detailed error.
	Start(ctx context.Context) error
	// Stop gracefully shuts down the plugin.
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeouts.
	//
	// Returns:
	//   - error: If stopping fails, returns a detailed error.
	Stop(ctx context.Context) error
	// HealthCheck returns the current plugin status.
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeouts.
	//
	// Returns:
	//   - *PluginHealth: Pointer to the current health status.
	//   - error: If health check fails, returns a detailed error.
	HealthCheck(ctx context.Context) (*PluginHealth, error)
	// Factory returns the OpenTelemetry component factory.
	//
	// Returns:
	//   - component.Factory: The OpenTelemetry component factory for this plugin.
	//   - error: If factory retrieval fails, returns a detailed error.
	Factory() (component.Factory, error)
}
