package core

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

	// Diagnostic component types
	TypeSystem     Type = "system"
	TypeKubernetes Type = "kubernetes"
	TypeCloud      Type = "cloud"
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
