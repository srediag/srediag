// Package types provides core types and interfaces for SREDIAG (Site Reliability Engineering Diagnostics)
// This package defines the fundamental types and constants used throughout the application
// for component identification, status tracking, and capability definition.
package types

// Type represents a component type in the SREDIAG system.
// It is used to identify different parts of the system and their roles.
type Type string

const (
	// TypeUnknown represents an unidentified component type.
	// Used as a default value when the type cannot be determined.
	TypeUnknown Type = "unknown"

	// TypeService represents a service component that provides core functionality.
	// Services are long-running processes that manage other components.
	TypeService Type = "service"

	// TypePlugin represents a plugin component that extends system functionality.
	// Plugins are dynamically loaded modules that add features to SREDIAG.
	TypePlugin Type = "plugin"

	// TypeDiagnostic represents a diagnostic component that performs system analysis.
	// Diagnostic components collect and analyze system health and performance data.
	TypeDiagnostic Type = "diagnostic"

	// TypeCollector represents a data collection component.
	// Collectors gather metrics, logs, and other telemetry data from various sources.
	TypeCollector Type = "collector"

	// TypeProcessor represents a data processing component.
	// Processors transform, filter, or aggregate collected data before analysis.
	TypeProcessor Type = "processor"

	// TypeExporter represents a data export component.
	// Exporters send processed data to external systems or storage.
	TypeExporter Type = "exporter"

	// TypeSystem represents system-level diagnostic components.
	// These components monitor host system resources and performance.
	TypeSystem Type = "system"

	// TypeKubernetes represents Kubernetes-specific diagnostic components.
	// These components monitor Kubernetes clusters and workloads.
	TypeKubernetes Type = "kubernetes"

	// TypeCloud represents cloud platform diagnostic components.
	// These components monitor cloud resources and services (AWS, GCP, Azure).
	TypeCloud Type = "cloud"

	// TypeSecurity represents security-focused diagnostic components.
	// These components monitor security-related aspects and compliance.
	TypeSecurity Type = "security"
)

// Status represents the operational state of a component.
// Used to track the lifecycle of components in the system.
type Status string

const (
	// StatusUnknown indicates the component's state cannot be determined.
	StatusUnknown Status = "unknown"

	// StatusLoaded indicates the component is loaded but not yet running.
	// This typically occurs after initialization but before Start().
	StatusLoaded Status = "loaded"

	// StatusRunning indicates the component is active and functioning normally.
	// The component has been successfully started and is performing its tasks.
	StatusRunning Status = "running"

	// StatusStopped indicates the component has been gracefully stopped.
	// The component has completed its shutdown procedure.
	StatusStopped Status = "stopped"

	// StatusError indicates the component has encountered an error.
	// The component may need intervention to resume normal operation.
	StatusError Status = "error"
)

// Capability represents a specific functionality that a component can provide.
// Used to advertise and discover component features dynamically.
type Capability string

const (
	// CapabilityMetrics indicates the ability to handle metric data.
	// Components with this capability can collect, process, or export metrics.
	CapabilityMetrics Capability = "metrics"

	// CapabilityTracing indicates the ability to handle distributed tracing.
	// Components with this capability can work with trace spans and context.
	CapabilityTracing Capability = "tracing"

	// CapabilityLogging indicates the ability to handle log data.
	// Components with this capability can process or manage log entries.
	CapabilityLogging Capability = "logging"

	// CapabilityDiagnostic indicates the ability to perform system diagnostics.
	// Components with this capability can analyze system health and performance.
	CapabilityDiagnostic Capability = "diagnostic"
)

// DiagnosticType represents specific types of diagnostic components.
// This type is used to categorize different diagnostic capabilities.
type DiagnosticType Type

const (
	// DiagnosticTypeSystem represents system-level diagnostics.
	// These diagnostics focus on host system metrics and health.
	DiagnosticTypeSystem DiagnosticType = DiagnosticType(TypeSystem)

	// DiagnosticTypeKubernetes represents Kubernetes diagnostics.
	// These diagnostics focus on Kubernetes cluster health and performance.
	DiagnosticTypeKubernetes DiagnosticType = DiagnosticType(TypeKubernetes)

	// DiagnosticTypeCloud represents cloud platform diagnostics.
	// These diagnostics focus on cloud resource utilization and health.
	DiagnosticTypeCloud DiagnosticType = DiagnosticType(TypeCloud)

	// DiagnosticTypeSecurity represents security diagnostics.
	// These diagnostics focus on security posture and compliance.
	DiagnosticTypeSecurity DiagnosticType = DiagnosticType(TypeSecurity)
)
