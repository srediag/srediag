// Package lifecycle provides base implementations for component lifecycle management.
package lifecycle

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

var (
	// ErrManagerNotRunning is returned when an operation requires the manager to be running, but it's not.
	ErrManagerNotRunning = fmt.Errorf("manager is not running")
	// ErrManagerAlreadyRunning is returned when an operation requires the manager to be stopped, but it's running.
	ErrManagerAlreadyRunning = fmt.Errorf("manager is already running")
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

// BaseManager provides basic lifecycle management functionality
type BaseManager struct {
	running bool
	mu      sync.RWMutex
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
	bm.mu.Lock()
	defer bm.mu.Unlock()

	if bm.running {
		return fmt.Errorf("manager is already running")
	}

	bm.running = true
	return nil
}

// Stop gracefully stops the component
func (bm *BaseManager) Stop(ctx context.Context) error {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	if !bm.running {
		return fmt.Errorf("manager is not running")
	}

	bm.running = false
	return nil
}

// IsRunning returns whether the manager is currently running
func (m *BaseManager) IsRunning() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.running
}

// SetRunning sets the running state of the manager
func (m *BaseManager) SetRunning(state bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.running = state
}

// CheckRunningState validates if the current running state matches the expected state
func (m *BaseManager) CheckRunningState(expectedRunning bool) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.running != expectedRunning {
		if expectedRunning {
			return ErrManagerNotRunning
		}
		return ErrManagerAlreadyRunning
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
