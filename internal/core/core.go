package core

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewVersionCmd creates a new command to display version information
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
