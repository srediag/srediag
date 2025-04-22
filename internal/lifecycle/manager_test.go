package lifecycle

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockManager implements a mock lifecycle manager for testing
type MockManager struct {
	*BaseManager
	startCalled bool
	stopCalled  bool
	startError  error
	stopError   error
}

func NewMockManager() *MockManager {
	return &MockManager{
		BaseManager: NewBaseManager(),
	}
}

func (m *MockManager) Start(ctx context.Context) error {
	m.startCalled = true
	if m.startError != nil {
		return m.startError
	}
	m.SetRunning(true)
	return nil
}

func (m *MockManager) Stop(ctx context.Context) error {
	m.stopCalled = true
	if m.stopError != nil {
		return m.stopError
	}
	m.SetRunning(false)
	return nil
}

func TestMockManager_Lifecycle(t *testing.T) {
	ctx := context.Background()
	manager := NewMockManager()

	// Test Start
	err := manager.Start(ctx)
	assert.NoError(t, err)
	assert.True(t, manager.startCalled)
	assert.True(t, manager.IsRunning())

	// Test Stop
	err = manager.Stop(ctx)
	assert.NoError(t, err)
	assert.True(t, manager.stopCalled)
	assert.False(t, manager.IsRunning())
}

func TestMockManager_Errors(t *testing.T) {
	ctx := context.Background()
	manager := NewMockManager()

	// Test Start error
	expectedErr := ErrInvalidState
	manager.startError = expectedErr
	err := manager.Start(ctx)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)

	// Test Stop error
	manager.stopError = expectedErr
	err = manager.Stop(ctx)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}
