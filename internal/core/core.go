// Package core provides foundational types and utilities for the SREDIAG system.
//
// This file provides the CLI command for displaying version/build information and related utilities.
package core

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewVersionCmd creates a new Cobra command to display SREDIAG version information.
//
// Usage:
//   - Add this command to your root Cobra CLI to provide a `srediag version` subcommand.
//   - Prints the semantic version, commit hash, and build date of the running SREDIAG binary.
//
// Best Practices:
//   - Always ensure build info is injected at build time for production releases.
//   - Use this command for all version reporting to ensure consistency.
//
// TODO:
//   - Consider adding more detailed build metadata (Go version, platform, etc).
//   - Consider supporting JSON or machine-readable output for CI/CD.
//
// Redundancy/Refactor:
//   - No redundancy; this is the canonical version command for SREDIAG.
func NewVersionCmd() *cobra.Command {
	buildInfo := DefaultBuildInfo

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Long:  "Print detailed version information about the SREDIAG build",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("SREDIAG %s (commit: %s, built at: %s)\n",
				buildInfo.Version,
				buildInfo.Commit,
				buildInfo.Date,
			)
		},
	}

	return cmd
}
