// Package plugin provides plugin types and utilities for SREDIAG
package plugin

import (
	"context"

	"github.com/srediag/srediag/internal/core"
)

// Type represents the type of a plugin
type Type string

const (
	// TypeDiagnostic represents a diagnostic plugin
	TypeDiagnostic Type = "diagnostic"
	// TypeAnalysis represents an analysis plugin
	TypeAnalysis Type = "analysis"
	// TypeManagement represents a management plugin
	TypeManagement Type = "management"
	// TypeIntegration represents an integration plugin
	TypeIntegration Type = "integration"
	// TypeSecurity represents a security plugin
	TypeSecurity Type = "security"
)

// Capability represents a plugin capability
type Capability string

const (
	// CapabilityMetrics indicates the plugin can collect metrics
	CapabilityMetrics Capability = "metrics"
	// CapabilityTracing indicates the plugin can collect traces
	CapabilityTracing Capability = "tracing"
	// CapabilityLogging indicates the plugin can collect logs
	CapabilityLogging Capability = "logging"
	// CapabilityAnalysis indicates the plugin can perform analysis
	CapabilityAnalysis Capability = "analysis"
	// CapabilityManagement indicates the plugin can manage resources
	CapabilityManagement Capability = "management"
	// CapabilitySecurity indicates the plugin can perform security checks
	CapabilitySecurity Capability = "security"
)

// Status represents the status of a plugin
type Status string

const (
	// StatusUnknown indicates the plugin status is unknown
	StatusUnknown Status = "unknown"
	// StatusLoaded indicates the plugin is loaded but not started
	StatusLoaded Status = "loaded"
	// StatusRunning indicates the plugin is running
	StatusRunning Status = "running"
	// StatusStopped indicates the plugin is stopped
	StatusStopped Status = "stopped"
	// StatusError indicates the plugin is in an error state
	StatusError Status = "error"
)

// Info represents plugin information
type Info struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Type         core.Type         `json:"type"`
	Capabilities []core.Capability `json:"capabilities"`
	Status       core.Status       `json:"status"`
	Error        string            `json:"error,omitempty"`
}

// Config represents plugin configuration
type Config struct {
	Name       string                 `json:"name"`
	Type       core.Type              `json:"type"`
	Enabled    bool                   `json:"enabled"`
	Settings   map[string]interface{} `json:"settings,omitempty"`
	Depends    []string               `json:"depends,omitempty"`
	AutoStart  bool                   `json:"auto_start"`
	AutoReload bool                   `json:"auto_reload"`
}

// Metadata represents plugin metadata
type Metadata struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Description  string            `json:"description"`
	Author       string            `json:"author"`
	License      string            `json:"license"`
	Homepage     string            `json:"homepage,omitempty"`
	Repository   string            `json:"repository,omitempty"`
	Type         core.Type         `json:"type"`
	Capabilities []core.Capability `json:"capabilities"`
	Settings     []Setting         `json:"settings,omitempty"`
}

// Setting represents a plugin setting
type Setting struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Required    bool        `json:"required"`
	Default     interface{} `json:"default,omitempty"`
	Description string      `json:"description"`
}

// Plugin represents a SREDIAG plugin
type Plugin interface {
	// GetName returns the name of the plugin
	GetName() string
	// GetVersion returns the version of the plugin
	GetVersion() string
	// GetType returns the type of the plugin
	GetType() string
	// GetCapabilities returns the capabilities of the plugin
	GetCapabilities() []core.Capability
	// GetStatus returns the status of the plugin
	GetStatus() core.Status
	// Start starts the plugin
	Start(ctx context.Context) error
	// Stop stops the plugin
	Stop(ctx context.Context) error
	// IsHealthy returns true if the plugin is healthy
	IsHealthy() bool
}
