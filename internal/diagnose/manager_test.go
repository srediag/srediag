package diagnose

import (
	"testing"

	"github.com/srediag/srediag/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiagnoseManagerCreation(t *testing.T) {
	logger, err := core.NewLogger(&core.Logger{
		Level:  "info",
		Format: "console",
	})
	require.NoError(t, err)

	mgr := NewDiagnoseManager(logger)
	assert.NotNil(t, mgr)
	assert.NotNil(t, mgr.logger)
}

func TestRunSystem(t *testing.T) {
	logger, err := core.NewLogger(&core.Logger{
		Level:  "info",
		Format: "console",
	})
	require.NoError(t, err)

	mgr := NewDiagnoseManager(logger)

	// This should return an error indicating unimplemented functionality
	err = mgr.RunSystem()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not yet implemented")
}

func TestRunPerformance(t *testing.T) {
	logger, err := core.NewLogger(&core.Logger{
		Level:  "info",
		Format: "console",
	})
	require.NoError(t, err)

	mgr := NewDiagnoseManager(logger)

	// This should succeed as it's basic performance diagnostics
	err = mgr.RunPerformance()
	assert.NoError(t, err)
}

func TestRunSecurity(t *testing.T) {
	logger, err := core.NewLogger(&core.Logger{
		Level:  "info",
		Format: "console",
	})
	require.NoError(t, err)

	mgr := NewDiagnoseManager(logger)

	// This should succeed as it's basic security diagnostics
	err = mgr.RunSecurity()
	assert.NoError(t, err)
}

func TestDiagnoseManagerMultipleCalls(t *testing.T) {
	logger, err := core.NewLogger(&core.Logger{
		Level:  "info",
		Format: "console",
	})
	require.NoError(t, err)

	mgr := NewDiagnoseManager(logger)

	// Multiple calls should not interfere with each other
	err1 := mgr.RunSystem()
	err2 := mgr.RunPerformance()
	err3 := mgr.RunSecurity()

	// System diagnostics now properly return error for unimplemented functionality
	assert.Error(t, err1)
	assert.Contains(t, err1.Error(), "not yet implemented")
	assert.NoError(t, err2) // Performance and security are still placeholders that return nil
	assert.NoError(t, err3)
}
