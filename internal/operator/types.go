package operator

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SrediagSpec defines the desired configuration for a SREDIAG instance
type SrediagSpec struct {
	// Version of the SREDIAG service
	Version string `json:"version"`

	// Service configuration
	Service ServiceConfig `json:"service"`

	// OpenTelemetry configuration
	Telemetry TelemetryConfig `json:"telemetry"`

	// Plugin configuration
	Plugins PluginsConfig `json:"plugins,omitempty"`

	// Security configuration
	Security SecurityConfig `json:"security,omitempty"`

	// Pod resources
	Resources ResourceSpec `json:"resources,omitempty"`
}

// ServiceConfig defines service-level configuration
type ServiceConfig struct {
	// Service name
	Name string `json:"name"`

	// Environment (e.g., production, staging)
	Environment string `json:"environment"`
}

// TelemetryConfig defines OpenTelemetry configuration
type TelemetryConfig struct {
	// Enable telemetry collection
	Enabled bool `json:"enabled"`

	// Service name for telemetry
	ServiceName string `json:"service_name"`

	// OTLP endpoint
	Endpoint string `json:"endpoint"`

	// Protocol (grpc/http)
	Protocol string `json:"protocol"`

	// Environment
	Environment string `json:"environment"`

	// Resource attributes
	ResourceAttributes map[string]string `json:"resource_attributes,omitempty"`
}

// PluginsConfig defines plugin configuration
type PluginsConfig struct {
	// Plugin directory path
	Directory string `json:"directory"`

	// Enabled plugins
	Enabled []string `json:"enabled"`

	// Plugin settings
	Settings map[string]map[string]interface{} `json:"settings,omitempty"`
}

// SecurityConfig defines security configuration
type SecurityConfig struct {
	// TLS configuration
	TLS TLSConfig `json:"tls"`

	// Authentication configuration
	Auth AuthConfig `json:"auth"`
}

// TLSConfig defines TLS configuration
type TLSConfig struct {
	// Enable TLS
	Enabled bool `json:"enabled"`

	// Secret name containing certificates
	SecretName string `json:"secret_name,omitempty"`
}

// AuthConfig defines authentication configuration
type AuthConfig struct {
	// Authentication type
	Type string `json:"type"`

	// Secret name containing credentials
	SecretName string `json:"secret_name"`
}

// ResourceSpec defines pod resources
type ResourceSpec struct {
	// Resource limits
	Limits ResourceRequirements `json:"limits,omitempty"`

	// Resource requests
	Requests ResourceRequirements `json:"requests,omitempty"`
}

// ResourceRequirements defines resource requirements
type ResourceRequirements struct {
	// CPU in millicores
	CPU string `json:"cpu,omitempty"`

	// Memory in bytes
	Memory string `json:"memory,omitempty"`
}

// SrediagStatus defines the observed state of a SREDIAG instance
type SrediagStatus struct {
	// Current phase of the instance
	Phase string `json:"phase"`

	// Observed conditions
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Installed plugins
	InstalledPlugins []string `json:"installed_plugins,omitempty"`

	// Last time the status was updated
	LastUpdated metav1.Time `json:"last_updated,omitempty"`
}

// Srediag is the custom type for a SREDIAG instance
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type Srediag struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SrediagSpec   `json:"spec,omitempty"`
	Status SrediagStatus `json:"status,omitempty"`
}

// SrediagList contains a list of SREDIAG instances
// +kubebuilder:object:root=true
type SrediagList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Srediag `json:"items"`
}
