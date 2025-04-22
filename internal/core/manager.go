// Package interfaces defines the core interfaces for the SREDIAG framework.
package core

import "context"

// Manager defines the interface for all system components that require lifecycle management.
// Any component that needs to be started, stopped, and monitored should implement this interface.
type Manager interface {
	// Start initializes and starts the component.
	// It should be idempotent and return an error if the component is already running.
	// The context can be used to cancel the startup process.
	Start(ctx context.Context) error

	// Stop gracefully shuts down the component.
	// It should be idempotent and return an error if the component is not running.
	// The context can be used to set a deadline for the shutdown process.
	Stop(ctx context.Context) error

	// IsRunning returns the current running state of the component.
	// This method should be thread-safe.
	IsRunning() bool

	// Health returns the current health status of the component.
	// This method should be thread-safe and used for health checking.
	Health() Health
}

// Health represents the health status of a component.
type Health struct {
	// Status indicates the current health status
	Status HealthStatus

	// Message provides additional details about the health status
	Message string

	// LastChecked is the timestamp of the last health check
	LastChecked int64

	// Details contains component-specific health information
	Details map[string]interface{}
}

// HealthStatus represents the possible health states of a component.
type HealthStatus int

const (
	// StatusUnknown indicates that the health status cannot be determined
	StatusUnknown HealthStatus = iota

	// StatusHealthy indicates that the component is functioning normally
	StatusHealthy

	// StatusDegraded indicates that the component is functioning but with issues
	StatusDegraded

	// StatusUnhealthy indicates that the component is not functioning properly
	StatusUnhealthy
)
