package diagnose

import (
	"testing"

	"github.com/srediag/srediag/internal/core"
	"github.com/stretchr/testify/assert"
)

func TestManagerCreation(t *testing.T) {
	manager := NewDiagnoseManager(nil)
	assert.NotNil(t, manager)
}

func TestSystemDiagnostics(t *testing.T) {
	// Test with nil logger should not panic
	manager := NewDiagnoseManager(nil)
	assert.NotNil(t, manager)

	// Test with proper logger
	logger, err := core.NewLogger(&core.Logger{
		Level:  "info",
		Format: "console",
	})
	assert.NoError(t, err)

	manager = NewDiagnoseManager(logger)
	// This might still error due to actual system operations, but shouldn't panic
	_ = manager.RunSystem()
}
