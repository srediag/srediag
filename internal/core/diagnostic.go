package core

import (
	"context"
)

// DiagnosticManager manages diagnostic components
type DiagnosticManager interface {
	Component
	// RegisterDiagnostic registers a diagnostic component
	RegisterDiagnostic(name string, diagnostic Diagnostic) error
	// UnregisterDiagnostic unregisters a diagnostic component
	UnregisterDiagnostic(name string) error
	// GetDiagnostic returns a diagnostic component by name
	GetDiagnostic(name string) (Diagnostic, error)
	// ListDiagnostics returns all registered diagnostic components
	ListDiagnostics() []Diagnostic
}

// Diagnostic represents a diagnostic component
type Diagnostic interface {
	Component
	// GetName returns the diagnostic name
	GetName() string
	// GetType returns the diagnostic type
	GetType() Type
	// GetVersion returns the diagnostic version
	GetVersion() string
	// Configure configures the diagnostic with the given configuration
	Configure(cfg interface{}) error
	// GetInterval returns the diagnostic interval
	GetInterval() string
	// GetThresholds returns the diagnostic thresholds
	GetThresholds() map[string]float64
	// Collect collects diagnostic data
	Collect(ctx context.Context) (map[string]interface{}, error)
}
