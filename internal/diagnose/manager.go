package diagnose

import (
	"context"

	"github.com/srediag/srediag/internal/core"
)

// Package diagnose provides diagnostic operations for SREDIAG, including system, performance, and security diagnostics.
//
// This file defines the DiagnoseManager, which orchestrates all diagnostic operations.
//
// Usage:
//   - Use DiagnoseManager to coordinate system, performance, and security diagnostics.
//   - Instantiate with NewDiagnoseManager, providing a logger.
//   - Call RunSystem, RunPerformance, or RunSecurity to execute diagnostics.
//
// Best Practices:
//   - Always check for errors from diagnostic methods.
//   - Use logger for all error and status reporting.
//
// TODO:
//   - Add context.Context to all methods for cancellation and timeouts.
// TODO: Implement plugin manager to inject CLI subcommands (see architecture/diagnose.md §2)
// TODO: Implement plugin execution in main process, with optional cmdhelper for heavy collectors (see architecture/diagnose.md §2)
// TODO: Implement control-plane feedback loop for remote diagnostics (see architecture/diagnose.md §6)
// TODO: Enforce error & exit-code semantics for diagnostics (see architecture/diagnose.md §7)
// TODO: Implement diagnostics metrics contract (see architecture/diagnose.md §8)

// DiagnoseManager orchestrates all diagnostic operations (system, performance, security).
//
// Usage:
//   - Instantiate with NewDiagnoseManager, providing a logger.
//   - Call RunSystem, RunPerformance, or RunSecurity to execute diagnostics.
type DiagnoseManager struct {
	logger *core.Logger
}

// NewDiagnoseManager creates a new DiagnoseManager.
//
// Parameters:
//   - logger: Logger for status and error reporting.
//
// Returns:
//   - *DiagnoseManager: A new DiagnoseManager instance.
func NewDiagnoseManager(logger *core.Logger) *DiagnoseManager {
	return &DiagnoseManager{logger: logger}
}

// RunSystem runs system diagnostics.
//
// Returns:
//   - error: If system diagnostics fail, returns a detailed error.
func (m *DiagnoseManager) RunSystem() error {
	d := NewSystemDiagnostics(m.logger)
	return d.Run(context.Background())
}

// RunPerformance runs performance diagnostics.
//
// Returns:
//   - error: If performance diagnostics fail, returns a detailed error.
func (m *DiagnoseManager) RunPerformance() error {
	d := NewPerformanceDiagnostics(m.logger)
	return d.Run(context.Background())
}

// RunSecurity runs security diagnostics.
//
// Returns:
//   - error: If security diagnostics fail, returns a detailed error.
func (m *DiagnoseManager) RunSecurity() error {
	d := NewSecurityDiagnostics(m.logger)
	return d.Run(context.Background())
}
