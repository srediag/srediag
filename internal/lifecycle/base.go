// Package lifecycle provides base implementations for component lifecycle management.
package lifecycle

import (
	"context"
	"sync/atomic"
)

// Component defines the interface for a lifecycle-managed component
type Component interface {
	// Start initializes and starts the component
	Start(ctx context.Context) error
	// Stop gracefully stops the component
	Stop(ctx context.Context) error
	// IsHealthy returns the health status of the component
	IsHealthy() bool
}

// BaseManager provides a base implementation for lifecycle management
type BaseManager struct {
	running atomic.Bool
	health  atomic.Value
}

// NewBaseManager creates a new instance of BaseManager
func NewBaseManager() *BaseManager {
	bm := &BaseManager{}
	bm.health.Store(true)
	return bm
}

// Start initializes and starts the component
func (bm *BaseManager) Start(ctx context.Context) error {
	if err := bm.CheckRunningState(false); err != nil {
		return err
	}
	bm.SetRunning(true)
	return nil
}

// Stop gracefully stops the component
func (bm *BaseManager) Stop(ctx context.Context) error {
	if err := bm.CheckRunningState(true); err != nil {
		return err
	}
	bm.SetRunning(false)
	return nil
}

// IsRunning returns true if the component is currently running
func (bm *BaseManager) IsRunning() bool {
	return bm.running.Load()
}

// SetRunning sets the running state of the component
func (bm *BaseManager) SetRunning(state bool) {
	bm.running.Store(state)
}

// CheckRunningState verifies if the component is in the expected running state
func (bm *BaseManager) CheckRunningState(expectedRunning bool) error {
	isRunning := bm.IsRunning()
	if expectedRunning && !isRunning {
		return ErrNotRunning
	}
	if !expectedRunning && isRunning {
		return ErrAlreadyRunning
	}
	return nil
}

// IsHealthy returns the current health status of the component
func (bm *BaseManager) IsHealthy() bool {
	return bm.health.Load().(bool)
}

// UpdateHealth updates the health status of the component
func (bm *BaseManager) UpdateHealth(status bool) {
	bm.health.Store(status)
}
