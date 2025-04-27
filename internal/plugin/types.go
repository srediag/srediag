// Package plugin provides a plugin management system for OpenTelemetry components
package plugin

import (
	"os/exec"
	"time"

	"github.com/cloudwego/shmipc-go"

	"github.com/srediag/srediag/internal/core"
)

// Metadata contains plugin metadata and identity.
type PluginMetadata struct {
	// Name is the unique identifier of the plugin
	Name string
	// Type indicates the plugin category
	Type core.ComponentType
	// Version of the plugin
	Version string
	// Description provides details about the plugin's functionality
	Description string
	// Capabilities lists the plugin's supported features
	Capabilities []string
	// SHA256 checksum of the plugin binary
	SHA256 string
	// Signature for plugin verification (optional)
	Signature string
}

// Health represents plugin health status.
type PluginHealth struct {
	// Status indicates the current state: "healthy", "degraded", "failed"
	Status string
	// LastCheck timestamp of the most recent health check
	LastCheck time.Time
	// Message provides additional status information (optional)
	Message string
	// Error details if the plugin is in a failed state (optional)
	Error string
}

// pluginInstance represents a running plugin
type pluginInstance struct {
	metadata PluginMetadata
	ch       *shmipc.SessionManager
	cmd      *exec.Cmd
}
