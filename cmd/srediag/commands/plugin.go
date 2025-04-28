// Package commands provides the command-line interface for the SREDIAG application.
package commands

import (
	"github.com/spf13/cobra"

	"github.com/srediag/srediag/internal/core"
	"github.com/srediag/srediag/internal/plugin"
)

// newPluginCmd creates the root command for managing plugins in the SREDIAG application.
// This command serves as a parent for subcommands that allow users to list, enable, disable,
// and retrieve information about plugins. The function wires these subcommands to their
// respective handlers in the internal/plugin package, ensuring that all business logic
// is encapsulated within the internal package. This design promotes separation of concerns
// by keeping the CLI wiring distinct from the core application logic.
// The command hierarchy is structured as follows:
// - plugin: The root command for plugin management.
//   - list: Lists all available plugins.
//   - info [name]: Displays detailed information about a specific plugin.
//   - enable [type] [name]: Enables a plugin of a specified type and name.
//   - disable [name]: Disables a plugin by its name.
//
// This function takes an AppContext as input, which provides the necessary context and
// dependencies for executing the commands.
func newPluginCmd(ctx *core.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugin",
		Short: "Manage plugins",
		Long:  `The plugin command allows you to list, enable, disable, and get information about plugins.`,
	}

	cmd.AddCommand(
		newPluginListCmd(ctx),
		newPluginInfoCmd(ctx),
		newPluginEnableCmd(ctx),
		newPluginDisableCmd(ctx),
	)

	return cmd
}

// newPluginListCmd wires the 'list' subcommand to plugin.CLI_List.
func newPluginListCmd(ctx *core.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all available plugins",
		RunE: func(cmd *cobra.Command, args []string) error {
			return plugin.CLI_List(ctx, cmd, args)
		},
	}
}

// newPluginInfoCmd wires the 'info' subcommand to plugin.CLI_Info.
func newPluginInfoCmd(ctx *core.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "info [name]",
		Short: "Show information about a plugin",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return plugin.CLI_Info(ctx, cmd, args)
		},
	}
}

// newPluginEnableCmd wires the 'enable' subcommand to plugin.CLI_Enable.
func newPluginEnableCmd(ctx *core.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "enable [type] [name]",
		Short: "Enable a plugin",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return plugin.CLI_Enable(ctx, cmd, args)
		},
	}
}

// newPluginDisableCmd wires the 'disable' subcommand to plugin.CLI_Disable.
func newPluginDisableCmd(ctx *core.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "disable [name]",
		Short: "Disable a plugin",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return plugin.CLI_Disable(ctx, cmd, args)
		},
	}
}
