package lifecycle

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBaseManager_InitialState(t *testing.T) {
	bm := NewBaseManager()
	if bm.IsRunning() {
		t.Error("new manager should not be running")
	}
	if !bm.IsHealthy() {
		t.Error("new manager should be healthy")
	}
}

func TestBaseManager_Start(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name      string
		setup     func(*BaseManager)
		wantError bool
	}{
		{
			name:      "start when stopped",
			setup:     func(bm *BaseManager) {},
			wantError: false,
		},
		{
			name: "start when already running",
			setup: func(bm *BaseManager) {
				bm.SetRunning(true)
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bm := NewBaseManager()
			tt.setup(bm)

			err := bm.Start(ctx)
			if (err != nil) != tt.wantError {
				t.Errorf("Start() error = %v, wantError %v", err, tt.wantError)
			}

			if !tt.wantError && !bm.IsRunning() {
				t.Error("Start() should set state to running")
			}
		})
	}
}

func TestBaseManager_Stop(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name      string
		setup     func(*BaseManager)
		wantError bool
	}{
		{
			name: "stop when running",
			setup: func(bm *BaseManager) {
				bm.SetRunning(true)
			},
			wantError: false,
		},
		{
			name:      "stop when already stopped",
			setup:     func(bm *BaseManager) {},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bm := NewBaseManager()
			tt.setup(bm)

			err := bm.Stop(ctx)
			if (err != nil) != tt.wantError {
				t.Errorf("Stop() error = %v, wantError %v", err, tt.wantError)
			}

			if !tt.wantError && bm.IsRunning() {
				t.Error("Stop() should set state to not running")
			}
		})
	}
}

func TestBaseManager_Health(t *testing.T) {
	bm := NewBaseManager()

	// Test initial state
	if !bm.IsHealthy() {
		t.Error("initial health state should be true")
	}

	// Test state updates
	bm.UpdateHealth(false)
	if bm.IsHealthy() {
		t.Error("health should be false after UpdateHealth(false)")
	}

	bm.UpdateHealth(true)
	if !bm.IsHealthy() {
		t.Error("health should be true after UpdateHealth(true)")
	}
}

func TestBaseManager_ConcurrentAccess(t *testing.T) {
	bm := NewBaseManager()
	done := make(chan bool)
	iterations := 1000

	// Test concurrent access to running state
	go func() {
		for i := 0; i < iterations; i++ {
			bm.SetRunning(true)
			bm.SetRunning(false)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < iterations; i++ {
			bm.IsRunning()
		}
		done <- true
	}()

	// Test concurrent access to health state
	go func() {
		for i := 0; i < iterations; i++ {
			bm.UpdateHealth(true)
			bm.UpdateHealth(false)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < iterations; i++ {
			bm.IsHealthy()
		}
		done <- true
	}()

	// Add timeout to prevent deadlock
	timeout := time.After(5 * time.Second)
	for i := 0; i < 4; i++ {
		select {
		case <-done:
			continue
		case <-timeout:
			t.Fatal("timeout waiting for goroutines to finish")
		}
	}
}

func TestNewBaseManager(t *testing.T) {
	manager := NewBaseManager()
	assert.NotNil(t, manager)
	assert.False(t, manager.IsRunning())
}

func TestBaseManager_SetRunning(t *testing.T) {
	manager := NewBaseManager()

	// Test setting to running
	manager.SetRunning(true)
	assert.True(t, manager.IsRunning())

	// Test setting to stopped
	manager.SetRunning(false)
	assert.False(t, manager.IsRunning())
}

func TestBaseManager_CheckRunningState(t *testing.T) {
	manager := NewBaseManager()

	// Test expecting an error when checking for running state (wantRunning=true) while stopped
	err := manager.CheckRunningState(true) // Should error: Expected running, but is stopped
	assert.Error(t, err)
	assert.Equal(t, ErrManagerNotRunning, err)

	// Test expecting no error when checking for stopped state (wantRunning=false) while stopped
	err = manager.CheckRunningState(false) // Should not error: Expected stopped, and is stopped
	assert.NoError(t, err)

	// Set to running
	manager.SetRunning(true)

	// Test expecting an error when checking for stopped state (wantRunning=false) while running
	err = manager.CheckRunningState(false) // Should error: Expected stopped, but is running
	assert.Error(t, err)
	assert.Equal(t, ErrManagerAlreadyRunning, err)

	// Test expecting no error when checking for running state (wantRunning=true) while running
	err = manager.CheckRunningState(true) // Should not error: Expected running, and is running
	assert.NoError(t, err)
}

func TestBaseManager_IsRunning(t *testing.T) {
	manager := NewBaseManager()

	// Initial state
	assert.False(t, manager.IsRunning())

	// After setting to running
	manager.SetRunning(true)
	assert.True(t, manager.IsRunning())

	// After setting to stopped
	manager.SetRunning(false)
	assert.False(t, manager.IsRunning())
}

func TestBaseManager_Lifecycle(t *testing.T) {
	manager := NewBaseManager()

	// Test initial state
	assert.False(t, manager.IsRunning())

	// Test setting running state
	manager.SetRunning(true)
	assert.True(t, manager.IsRunning())

	manager.SetRunning(false)
	assert.False(t, manager.IsRunning())
}
