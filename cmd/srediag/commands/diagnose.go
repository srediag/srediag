// Package commands provides the command-line interface for the SREDIAG application.
package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/srediag/srediag/internal/core"
	"github.com/srediag/srediag/internal/diagnose"
)

// newDiagnoseCmd creates a new command for running system diagnostics
// Only CLI wiring is present here; all business logic is delegated to internal/diagnostic CLI_* functions.
func newDiagnoseCmd(ctx *core.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diagnose [type]",
		Short: "Run system diagnostics",
		Long: `The diagnose command runs various diagnostic checks to identify
potential issues and provide insights about the system.

Available diagnostic types:
  - system: Check system health, resource usage, and configuration
  - performance: Analyze system and application performance metrics
  - security: Check security configurations and potential vulnerabilities`,
		RunE: runDiagnose,
	}

	// Add subcommands
	cmd.AddCommand(
		newSystemDiagCmd(ctx),
		newPerformanceDiagCmd(ctx),
		newSecurityDiagCmd(ctx),
	)

	return cmd
}

func runDiagnose(cmd *cobra.Command, args []string) error {
	fmt.Println("Please specify a diagnostic type to run (use --help for available types)")
	return cmd.Help()
}

// newSystemDiagCmd wires the 'system' subcommand to diagnostic.CLI_SystemDiagnostics.
func newSystemDiagCmd(ctx *core.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "system",
		Short: "Run system diagnostics",
		Long:  `Check system health, resource usage, and configuration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return diagnose.CLI_SystemDiagnostics(ctx, cmd, args)
		},
	}
}

// newPerformanceDiagCmd wires the 'performance' subcommand to diagnostic.CLI_PerformanceDiagnostics.
func newPerformanceDiagCmd(ctx *core.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "performance",
		Short: "Run performance diagnostics",
		Long:  `Analyze system and application performance metrics.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return diagnose.CLI_PerformanceDiagnostics(ctx, cmd, args)
		},
	}
}

// newSecurityDiagCmd wires the 'security' subcommand to diagnostic.CLI_SecurityDiagnostics.
func newSecurityDiagCmd(ctx *core.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "security",
		Short: "Run security diagnostics",
		Long:  `Check security configurations and potential vulnerabilities.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return diagnose.CLI_SecurityDiagnostics(ctx, cmd, args)
		},
	}
}
