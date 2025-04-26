// Package commands provides the command-line interface for the SREDIAG application.
package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/srediag/srediag/internal/diagnostics"
)

// newDiagnoseCmd creates a new command for running system diagnostics
func newDiagnoseCmd(opts *Options) *cobra.Command {
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
		newSystemDiagCmd(opts),
		newPerformanceDiagCmd(opts),
		newSecurityDiagCmd(opts),
	)

	return cmd
}

func runDiagnose(cmd *cobra.Command, args []string) error {
	fmt.Println("Please specify a diagnostic type to run (use --help for available types)")
	return cmd.Help()
}

// newSystemDiagCmd creates a command for running system diagnostics
func newSystemDiagCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "system",
		Short: "Run system diagnostics",
		Long:  `Check system health, resource usage, and configuration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			diag := diagnostics.NewSystemDiagnostics(opts.Settings.GetLogger())
			return diag.Run(cmd.Context())
		},
	}
}

// newPerformanceDiagCmd creates a command for running performance diagnostics
func newPerformanceDiagCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "performance",
		Short: "Run performance diagnostics",
		Long:  `Analyze system and application performance metrics.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			diag := diagnostics.NewPerformanceDiagnostics(opts.Settings.GetLogger())
			return diag.Run(cmd.Context())
		},
	}
}

// newSecurityDiagCmd creates a command for running security diagnostics
func newSecurityDiagCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "security",
		Short: "Run security diagnostics",
		Long:  `Check security configurations and potential vulnerabilities.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			diag := diagnostics.NewSecurityDiagnostics(opts.Settings.GetLogger())
			return diag.Run(cmd.Context())
		},
	}
}
