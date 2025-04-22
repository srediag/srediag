package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pluginCmd())
	rootCmd.AddCommand(versionCmd())
}

func pluginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugin",
		Short: "Manage SREDIAG plugins",
		Long:  `Commands to manage SREDIAG plugins: list, install, remove, etc.`,
	}

	// Plugin management subcommands
	cmd.AddCommand(
		&cobra.Command{
			Use:   "list",
			Short: "List all installed plugins",
			Run: func(cmd *cobra.Command, args []string) {
				// TODO: Implement plugin listing
				fmt.Println("Listing plugins...")
			},
		},
		&cobra.Command{
			Use:   "install [plugin]",
			Short: "Install a new plugin",
			Args:  cobra.ExactArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				// TODO: Implement plugin installation
				fmt.Printf("Installing plugin %s...\n", args[0])
			},
		},
		&cobra.Command{
			Use:   "remove [plugin]",
			Short: "Remove an installed plugin",
			Args:  cobra.ExactArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				// TODO: Implement plugin removal
				fmt.Printf("Removing plugin %s...\n", args[0])
			},
		},
	)

	return cmd
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show SREDIAG version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("SREDIAG version %s\n", Version)
		},
	}
}
