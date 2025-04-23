package commands

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func newDiagnoseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diagnose",
		Short: "Run diagnostics",
		Long: `Run various diagnostic checks on the system, Kubernetes, or cloud resources.
		
Examples:
  # Basic system diagnostics
  srediag diagnose system
  srediag diagnose system --resource cpu
  srediag diagnose system --resource memory
  srediag diagnose system --resource disk`,
	}

	// Add subcommands
	cmd.AddCommand(
		newDiagnoseSystemCmd(),
		newDiagnoseKubernetesCmd(),
	)

	return cmd
}

func newDiagnoseSystemCmd() *cobra.Command {
	var resource string

	cmd := &cobra.Command{
		Use:   "system [--resource <resource>]",
		Short: "Run system diagnostics",
		Long:  "Run diagnostics on the local system resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := getLogger()
			logger.Info("running system diagnostics",
				zap.String("resource", resource))
			return nil
		},
	}

	cmd.Flags().StringVar(&resource, "resource", "", "resource to diagnose (cpu/memory/disk)")
	return cmd
}

func newDiagnoseKubernetesCmd() *cobra.Command {
	var cluster string

	cmd := &cobra.Command{
		Use:   "kubernetes [--cluster <cluster>]",
		Short: "Run Kubernetes diagnostics",
		Long:  "Run diagnostics on Kubernetes clusters",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := getLogger()
			logger.Info("running kubernetes diagnostics",
				zap.String("cluster", cluster))
			return nil
		},
	}

	cmd.Flags().StringVar(&cluster, "cluster", "", "target Kubernetes cluster")
	return cmd
}
