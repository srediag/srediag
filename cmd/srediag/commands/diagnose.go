package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/srediag/srediag/internal/diagnostics"
)

func newDiagnoseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diagnose",
		Short: "Run system diagnostics",
		Long: `The diagnose command runs various diagnostic checks to identify
potential issues and provide insights about the system.`,
		RunE: runDiagnose,
	}

	// Add subcommands
	cmd.AddCommand(
		newSystemDiagCmd(),
		newPerformanceDiagCmd(),
		newSecurityDiagCmd(),
	)

	return cmd
}

func runDiagnose(cmd *cobra.Command, args []string) error {
	fmt.Println("Please specify a diagnostic type to run")
	return cmd.Help()
}

func newSystemDiagCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "system",
		Short: "Run system diagnostics",
		Long:  `Check system health, resource usage, and configuration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			diag := diagnostics.NewSystemDiagnostics(cmdSettings.GetLogger())
			return diag.Run(context.Background())
		},
	}
}

func newPerformanceDiagCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "performance",
		Short: "Run performance diagnostics",
		Long:  `Analyze system and application performance metrics.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			diag := diagnostics.NewPerformanceDiagnostics(cmdSettings.GetLogger())
			return diag.Run(context.Background())
		},
	}
}

func newSecurityDiagCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "security",
		Short: "Run security diagnostics",
		Long:  `Check security configurations and potential vulnerabilities.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			diag := diagnostics.NewSecurityDiagnostics(cmdSettings.GetLogger())
			return diag.Run(context.Background())
		},
	}
}
