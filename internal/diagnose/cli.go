package diagnose

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/srediag/srediag/internal/core"
)

// TODO(D-01 Phase 3): Implement system diagnostics plugin (CPU, mem, IO, net) (see TODO.md D-01, Phase 3)
// TODO(D-02 Phase 3): Implement Kubernetes diagnostics plugin (cluster, node, pod) (see TODO.md D-02, Phase 3)
// TODO(D-03 Phase 5): Implement cloud provider stubs (AWS, Azure, GCP) (see TODO.md D-03, Phase 5)
// TODO(D-04 Phase 5): Implement IaC analyzers (Terraform, K8s manifests, Helm) (see TODO.md D-04, Phase 5)

// Package diagnose provides diagnostic operations for SREDIAG, including system, performance, and security diagnostics.
//
// This file defines CLI entrypoints for diagnostic commands, wiring Cobra commands to internal diagnostic logic.
//
// Usage:
//   - Use these CLI functions as entrypoints for 'srediag diagnose' subcommands.
//   - Each function extracts parameters from the CLI context, instantiates the DiagnoseManager, and delegates to the appropriate method.
//
// Best Practices:
//   - Always validate required flags and parameters before calling DiagnoseManager methods.
//   - Log all errors and important events for traceability.
//   - Use context-aware logging and error handling for better diagnostics.
//
// TODO:
//   - Add context.Context support for cancellation and timeouts.
//   - Refactor to reduce repeated logger fallback logic.

// CLI_SystemDiagnostics is the entrypoint for 'srediag diagnose system'.
//
// Parameters:
//   - ctx: Application context containing logger and configuration.
//   - cmd: Cobra command instance.
//   - args: Command-line arguments.
//
// Returns:
//   - error: If system diagnostics fail, returns a detailed error.
func CLI_SystemDiagnostics(ctx *core.AppContext, cmd *cobra.Command, args []string) error {
	logger := ctx.Logger
	if logger == nil {
		var err error
		logger, err = core.NewLogger(nil)
		if err != nil {
			return fmt.Errorf("failed to create fallback logger: %w", err)
		}
	}
	mgr := NewDiagnoseManager(logger)
	if err := mgr.RunSystem(); err != nil {
		logger.Error("System diagnostics failed", core.ZapError(err))
		return fmt.Errorf("system diagnostics failed: %w", err)
	}
	logger.Info("System diagnostics completed successfully")
	return nil
}

// CLI_PerformanceDiagnostics is the entrypoint for 'srediag diagnose performance'.
//
// Parameters:
//   - ctx: Application context containing logger and configuration.
//   - cmd: Cobra command instance.
//   - args: Command-line arguments.
//
// Returns:
//   - error: If performance diagnostics fail, returns a detailed error.
func CLI_PerformanceDiagnostics(ctx *core.AppContext, cmd *cobra.Command, args []string) error {
	logger := ctx.Logger
	if logger == nil {
		var err error
		logger, err = core.NewLogger(nil)
		if err != nil {
			return fmt.Errorf("failed to create fallback logger: %w", err)
		}
	}
	mgr := NewDiagnoseManager(logger)
	if err := mgr.RunPerformance(); err != nil {
		logger.Error("Performance diagnostics failed", core.ZapError(err))
		return fmt.Errorf("performance diagnostics failed: %w", err)
	}
	logger.Info("Performance diagnostics completed successfully")
	return nil
}

// CLI_SecurityDiagnostics is the entrypoint for 'srediag diagnose security'.
//
// Parameters:
//   - ctx: Application context containing logger and configuration.
//   - cmd: Cobra command instance.
//   - args: Command-line arguments.
//
// Returns:
//   - error: If security diagnostics fail, returns a detailed error.
func CLI_SecurityDiagnostics(ctx *core.AppContext, cmd *cobra.Command, args []string) error {
	logger := ctx.Logger
	if logger == nil {
		var err error
		logger, err = core.NewLogger(nil)
		if err != nil {
			return fmt.Errorf("failed to create fallback logger: %w", err)
		}
	}
	mgr := NewDiagnoseManager(logger)
	if err := mgr.RunSecurity(); err != nil {
		logger.Error("Security diagnostics failed", core.ZapError(err))
		return fmt.Errorf("security diagnostics failed: %w", err)
	}
	logger.Info("Security diagnostics completed successfully")
	return nil
}
