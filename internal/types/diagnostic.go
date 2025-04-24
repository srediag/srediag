package types

import "context"

// IDiagnostic represents a diagnostic component
type IDiagnostic interface {
	IComponent
	// GetInterval returns the interval at which the diagnostic component should run
	GetInterval() string
	// GetThresholds returns the thresholds for the diagnostic component
	GetThresholds() map[string]float64
	// Collect collects diagnostic data
	Collect(ctx context.Context) (map[string]interface{}, error)
}

// IDiagnosticManager represents the diagnostic manager interface
type IDiagnosticManager interface {
	IComponent
	// RegisterDiagnostic registers a new diagnostic component
	RegisterDiagnostic(name string, diagnostic IDiagnostic) error
	// UnregisterDiagnostic unregisters a diagnostic component
	UnregisterDiagnostic(name string) error
	// GetDiagnostic returns a diagnostic component by name
	GetDiagnostic(name string) (IDiagnostic, error)
	// ListDiagnostics returns all registered diagnostic components
	ListDiagnostics() []IDiagnostic
}

// DefaultDiagnosticConfig returns the default diagnostic configuration
func DefaultDiagnosticConfig() *DiagnosticConfig {
	return &DiagnosticConfig{
		System: SystemConfig{
			Enabled:     true,
			Interval:    "30s",
			CPULimit:    80.0,
			MemoryLimit: 90.0,
			DiskLimit:   85.0,
		},
		Kubernetes: KubernetesConfig{
			Enabled:   false,
			Clusters:  []string{},
			Namespace: "default",
		},
		Cloud: CloudConfig{
			Enabled:     false,
			Providers:   []string{},
			Credentials: make(map[string]string),
		},
	}
}
