// Package commands provides the command-line interface for the SREDIAG application.
package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/srediag/srediag/internal/core"
	"github.com/srediag/srediag/internal/diagnose"
)

// newDiagnoseCmd creates a new command for running system diagnostics
// Only CLI wiring is present here; all business logic is delegated to internal/diagnostic functions:
// CLI_SystemDiagnostics, CLI_PerformanceDiagnostics, and CLI_SecurityDiagnostics.
func newDiagnoseCmd(ctx *core.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diagnose [type]",
		Short: "Run system diagnostics",
		Long: `The diagnose command runs diagnostic checks to identify potential issues.
It provides insights about the system's health, performance, and security.

Available diagnostic types:
  - system: Check system health, resource usage, and configuration.
  - performance: Analyze system and application performance metrics.
  - security: Check security configurations and potential vulnerabilities.`,
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
	if len(args) == 0 {
		fmt.Println("Please specify a diagnostic type (e.g., 'system', 'performance', 'security'). Use --help for more details.")
		return cmd.Help()
	}

	validTypes := map[string]bool{
		"system":      true,
		"performance": true,
		"security":    true,
	}

	if !validTypes[args[0]] {
		fmt.Printf("Error: Invalid diagnostic type '%s'. Please specify a valid type (use --help for available types).\n", args[0])
		return cmd.Help()
	}

	fmt.Printf("Running '%s' diagnostics...\n", args[0])
	return nil
}

// newSystemDiagCmd wires the 'system' subcommand to diagnostic.CLI_SystemDiagnostics.
func newSystemDiagCmd(ctx *core.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "system",
		Short: "Run system diagnostics",
		Long:  `Check system health, resource usage, and configuration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := diagnose.CLI_SystemDiagnostics(ctx, cmd, args)
			if err != nil {
				fmt.Printf("Error running system diagnostics: %v\n", err)
			}
			return err
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
			err := diagnose.CLI_PerformanceDiagnostics(ctx, cmd, args)
			if err != nil {
				fmt.Printf("Error running performance diagnostics: %v\n", err)
			}
			return err
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
			err := diagnose.CLI_SecurityDiagnostics(ctx, cmd, args)
			if err != nil {
				fmt.Printf("Error running security diagnostics: %v\n", err)
			}
			return err
		},
	}
}
