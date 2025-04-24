// Package types provides plugin-related types and interfaces for SREDIAG.
package types

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/component"
)

// ErrInvalidPluginConfig represents an error with plugin configuration
type ErrInvalidPluginConfig string

func (e ErrInvalidPluginConfig) Error() string {
	return string(e)
}

// PluginCategory represents the category of a plugin
type PluginCategory string

const (
	// PluginCategoryReceiver represents a receiver plugin
	PluginCategoryReceiver PluginCategory = "receiver"
	// PluginCategoryProcessor represents a processor plugin
	PluginCategoryProcessor PluginCategory = "processor"
	// PluginCategoryExporter represents an exporter plugin
	PluginCategoryExporter PluginCategory = "exporter"
)

// PluginCapability represents a plugin capability
type PluginCapability string

const (
	// PluginCapabilityMetrics represents metrics capability
	PluginCapabilityMetrics PluginCapability = "metrics"
	// PluginCapabilityTraces represents traces capability
	PluginCapabilityTraces PluginCapability = "traces"
	// PluginCapabilityLogs represents logs capability
	PluginCapabilityLogs PluginCapability = "logs"
)

// PluginCapabilities represents a set of plugin capabilities
type PluginCapabilities map[PluginCapability]bool

// HasCapability checks if the capabilities include a specific capability
func (c PluginCapabilities) HasCapability(cap PluginCapability) bool {
	return c[cap]
}

// PluginLifecycle represents the lifecycle state of a plugin
type PluginLifecycle string

const (
	// LifecycleUnregistered represents an unregistered plugin
	LifecycleUnregistered PluginLifecycle = "unregistered"
	// LifecycleInitialized represents an initialized plugin
	LifecycleInitialized PluginLifecycle = "initialized"
	// LifecycleRunning represents a running plugin
	LifecycleRunning PluginLifecycle = "running"
	// LifecyclePaused represents a paused plugin
	LifecyclePaused PluginLifecycle = "paused"
	// LifecycleStopped represents a stopped plugin
	LifecycleStopped PluginLifecycle = "stopped"
	// LifecycleError represents a plugin in error state
	LifecycleError PluginLifecycle = "error"
)

// PluginMetadata contains metadata about a plugin
type PluginMetadata struct {
	ID           string             `json:"id"`
	Name         string             `json:"name"`
	Version      string             `json:"version"`
	Description  string             `json:"description"`
	Category     PluginCategory     `json:"category"`
	Capabilities PluginCapabilities `json:"capabilities"`
	Author       string             `json:"author"`
	License      string             `json:"license"`
	Homepage     string             `json:"homepage"`
	Repository   string             `json:"repository"`
	Tags         []string           `json:"tags"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
}

// PluginConfig represents the configuration for a plugin
type PluginConfig struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name"`
	Enabled  bool                   `json:"enabled"`
	Settings map[string]interface{} `json:"settings"`
}

// Validate validates the plugin configuration
func (c *PluginConfig) Validate() error {
	if c.ID == "" {
		return ErrInvalidPluginConfig("plugin ID is required")
	}
	if c.Name == "" {
		return ErrInvalidPluginConfig("plugin name is required")
	}
	return nil
}

// IPlugin defines the interface that all plugins must implement
type IPlugin interface {
	// GetID returns the unique identifier of the plugin
	GetID() string
	// GetCategory returns the category of the plugin
	GetCategory() PluginCategory
	// GetCapabilities returns the capabilities of the plugin
	GetCapabilities() PluginCapabilities
	// Validate validates the plugin configuration
	Validate() error
	// Start starts the plugin
	Start(ctx context.Context) error
	// Stop stops the plugin
	Stop(ctx context.Context) error
	// Component returns the underlying OpenTelemetry component
	Component() component.Component
}

// IPluginManager defines the interface for plugin management
type IPluginManager interface {
	// RegisterPlugin registers a plugin with the manager
	RegisterPlugin(plugin IPlugin) error
	// UnregisterPlugin unregisters a plugin from the manager
	UnregisterPlugin(pluginID string) error
	// GetPlugin returns a plugin by ID
	GetPlugin(pluginID string) (IPlugin, error)
	// ListPlugins returns a list of all registered plugins
	ListPlugins() []IPlugin
	// ListPluginsByCategory returns a list of plugins by category
	ListPluginsByCategory(category PluginCategory) []IPlugin
	// ListPluginsByCapability returns a list of plugins by capability
	ListPluginsByCapability(capability PluginCapability) []IPlugin
	// LoadPlugin loads a plugin with configuration
	LoadPlugin(pluginID string, config PluginConfig) error
	// UnloadPlugin unloads a plugin
	UnloadPlugin(pluginID string) error
	// StartPlugin starts a plugin
	StartPlugin(pluginID string) error
	// StopPlugin stops a plugin
	StopPlugin(pluginID string) error
	// PausePlugin pauses a plugin
	PausePlugin(pluginID string) error
	// ResumePlugin resumes a plugin
	ResumePlugin(pluginID string) error
	// GetPluginStatus returns the status of a plugin
	GetPluginStatus(pluginID string) (PluginLifecycle, error)
	// IsPluginHealthy returns the health status of a plugin
	IsPluginHealthy(pluginID string) bool
	// GetPluginErrors returns any errors for a plugin
	GetPluginErrors(pluginID string) []error
	// StartAll starts all registered plugins
	StartAll(ctx context.Context) error
	// StopAll stops all registered plugins
	StopAll(ctx context.Context) error
	// GetSystemHealth returns the health status of all plugins
	GetSystemHealth() map[string]bool
	// IsHealthy returns the overall system health status
	IsHealthy() bool
	// GetName returns the component name
	GetName() string
	// GetVersion returns the component version
	GetVersion() string
	// GetType returns the component type
	GetType() ComponentType
	// Configure configures the component
	Configure(cfg interface{}) error
}
