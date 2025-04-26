// Package commands provides the command-line interface for the SREDIAG application.
package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// BuildInfo holds version information about the build
type BuildInfo struct {
	Version string
	Commit  string
	Date    string
}

// DefaultBuildInfo provides default build information when not set during compilation
var DefaultBuildInfo = BuildInfo{
	Version: "dev",
	Commit:  "none",
	Date:    "unknown",
}

// newVersionCmd creates a new command to display version information
func newVersionCmd() *cobra.Command {
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
