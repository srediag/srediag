// Package plugin provides a plugin management system for OpenTelemetry components.
//
// This file defines the core types for plugin metadata, health, and internal plugin instance representation.
//
// Usage:
//   - Use PluginMetadata to describe plugin identity, version, and security attributes for registration, discovery, and validation.
//   - Use PluginHealth to represent plugin health status, diagnostics, and error reporting for monitoring and orchestration.
//   - pluginInstance is used internally for managing running plugin processes and IPC communication.
//
// Best Practices:
//   - Always populate all required fields in PluginMetadata and PluginHealth.
//   - Use SHA256 and Signature fields for plugin integrity and verification.
//   - Keep LastCheck updated for accurate health monitoring.
package plugin

import (
	"os/exec"
	"time"

	"github.com/cloudwego/shmipc-go"

	"github.com/srediag/srediag/internal/core"
)

// PluginMetadata describes all identity, version, and security attributes for a plugin.
//
// This struct is used for plugin registration, discovery, validation, and orchestration.
// It is the canonical source of truth for plugin identity and capabilities in the SREDIAG system.
//
// Fields:
//   - Name: Globally unique identifier for the plugin. Used for registration, lookup, and orchestration. Must not be empty.
//   - Type: Plugin category (receiver, processor, exporter, extension). Used for routing, compatibility, and grouping.
//   - Version: Semantic version string (e.g., "v1.2.3"). Used for compatibility checks, upgrades, and reporting. Should follow semver.
//   - Description: Human-readable details about the plugin's functionality, purpose, and usage. Used for operator visibility and documentation.
//   - Capabilities: List of supported features (e.g., "metrics", "logs"). Used for plugin discovery, compatibility, and feature negotiation.
//   - SHA256: Hex-encoded SHA256 checksum of the plugin binary. Used for integrity verification and supply chain security. Should be validated before loading.
//   - Signature: Optional cryptographic signature for plugin authenticity. Used for trust validation and secure plugin distribution. May be empty if unsigned.
type PluginMetadata struct {
	// Name is the globally unique identifier of the plugin.
	// This must be unique within the SREDIAG deployment and is used for registration, lookup, and orchestration.
	Name string
	// Type indicates the plugin category (receiver, processor, exporter, extension).
	// This is used for routing, compatibility, and grouping in the plugin manager.
	Type core.ComponentType
	// Version is the semantic version of the plugin (e.g., "v1.2.3").
	// Used for compatibility checks, upgrades, and reporting. Should follow semantic versioning.
	Version string
	// Description provides human-readable details about the plugin's functionality, purpose, and usage.
	// This is used for operator visibility, documentation, and diagnostics.
	Description string
	// Capabilities lists the plugin's supported features (e.g., "metrics", "logs").
	// Used for plugin discovery, compatibility, and feature negotiation.
	Capabilities []string
	// SHA256 is the hex-encoded SHA256 checksum of the plugin binary.
	// Used for integrity verification and supply chain security. Should be validated before loading the plugin.
	SHA256 string
	// Signature is an optional cryptographic signature for plugin authenticity.
	// Used for trust validation and secure plugin distribution. May be empty if the plugin is unsigned.
	Signature string
}

// PluginHealth represents the health status, diagnostics, and error reporting for a plugin.
//
// This struct is used for monitoring, orchestration, and operator visibility. It is updated by health checks and diagnostic routines.
//
// Fields:
//   - Status: Current state of the plugin. Expected values: "healthy", "degraded", "failed". Used for orchestration and alerting.
//   - LastCheck: Timestamp of the most recent health check. Should be updated on every health probe. Used for staleness detection.
//   - Message: Additional status information or diagnostics. Optional, for operator visibility and troubleshooting.
//   - Error: Error details if the plugin is in a failed or degraded state. Optional, for troubleshooting and root cause analysis.
type PluginHealth struct {
	// Status indicates the current state of the plugin: "healthy", "degraded", or "failed".
	// This field is set by health checks and is used for orchestration and alerting.
	Status string
	// LastCheck is the timestamp of the most recent health check.
	// Should be set to time.Now() on each probe. Used for staleness detection and monitoring.
	LastCheck time.Time
	// Message provides additional status information or diagnostics.
	// This field is optional and is used for operator visibility and troubleshooting.
	Message string
	// Error contains error details if the plugin is in a failed or degraded state.
	// This field is optional and is used for troubleshooting and root cause analysis.
	Error string
}

// pluginInstance represents a running plugin process and its associated IPC session.
//
// This struct is used internally by the plugin manager to track the state, communication channel, and process handle for each running plugin.
// It is not exported and should not be used outside the plugin management subsystem.
type pluginInstance struct {
	// metadata contains the identity and capabilities of the running plugin.
	metadata PluginMetadata
	// ch is the shmipc SessionManager for IPC communication with the plugin process.
	ch *shmipc.SessionManager
	// cmd is the exec.Cmd handle for the running plugin process.
	cmd *exec.Cmd
}
