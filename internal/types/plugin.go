// Package types provides plugin-related types and interfaces for SREDIAG.
package types

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Common errors
var (
	// ErrPluginNotFound is returned when a plugin cannot be found
	ErrPluginNotFound = errors.New("plugin not found")
)

// ErrInvalidPluginConfig represents an error with plugin configuration
type ErrInvalidPluginConfig string

func (e ErrInvalidPluginConfig) Error() string {
	return string(e)
}

// PluginCategory represents the category of a plugin
type PluginCategory string

const (
	// PluginCategoryDiagnostic represents a diagnostic plugin
	PluginCategoryDiagnostic PluginCategory = "diagnostic"
	// PluginCategoryTelemetry represents a telemetry plugin
	PluginCategoryTelemetry PluginCategory = "telemetry"
	// PluginCategoryCollector represents a collector plugin
	PluginCategoryCollector PluginCategory = "collector"
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
	Name        string         `json:"name"`
	Version     string         `json:"version"`
	Type        ComponentType  `json:"type"`
	Category    PluginCategory `json:"category"`
	Settings    interface{}    `json:"settings,omitempty"`
	Enabled     bool           `json:"enabled"`
	Description string         `json:"description,omitempty"`
}

// Validate validates the plugin configuration
func (c *PluginConfig) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("plugin name is required")
	}
	if c.Version == "" {
		return fmt.Errorf("plugin version is required")
	}
	if c.Type == "" {
		return fmt.Errorf("plugin type is required")
	}
	if c.Category == "" {
		return fmt.Errorf("plugin category is required")
	}
	return nil
}

// IPlugin defines the interface that all plugins must implement
type IPlugin interface {
	// GetName returns the plugin name
	GetName() string
	// GetVersion returns the plugin version
	GetVersion() string
	// GetType returns the plugin type
	GetType() ComponentType
	// GetCategory returns the plugin category
	GetCategory() PluginCategory
	// Configure configures the plugin with settings
	Configure(settings ComponentSettings) error
	// Start starts the plugin
	Start(ctx context.Context) error
	// Stop stops the plugin
	Stop(ctx context.Context) error
	// GetStatus returns the plugin status
	GetStatus() ComponentStatus
	// Validate validates the plugin configuration
	Validate() error
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
