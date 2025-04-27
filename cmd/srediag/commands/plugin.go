// Package commands provides the command-line interface for the SREDIAG application.
package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/srediag/srediag/internal/core"
	"github.com/srediag/srediag/internal/plugin"
)

// newPluginCmd creates a new command for managing plugins
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

func getPluginManager(ctx *core.AppContext) *plugin.PluginManager {
	return plugin.NewManager(ctx.GetLogger(), ctx.GetConfig().PluginsDir)
}

func newPluginListCmd(ctx *core.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all available plugins",
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr := getPluginManager(ctx)
			plugins := mgr.List()
			for _, p := range plugins {
				fmt.Printf("%s\t%s\t%s\n", p.Name, p.Type, p.Version)
			}
			return nil
		},
	}
}

func newPluginInfoCmd(ctx *core.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "info [name]",
		Short: "Show information about a plugin",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr := getPluginManager(ctx)
			inst, ok := mgr.Get(args[0])
			if !ok {
				return fmt.Errorf("plugin '%s' not found", args[0])
			}
			// Try to get metadata from the instance
			type metaGetter interface{ Metadata() plugin.PluginMetadata }
			if mg, ok := inst.(metaGetter); ok {
				meta := mg.Metadata()
				fmt.Printf("Name: %s\nType: %s\nVersion: %s\nDescription: %s\nCapabilities: %v\nSHA256: %s\nSignature: %s\n",
					meta.Name, meta.Type, meta.Version, meta.Description, meta.Capabilities, meta.SHA256, meta.Signature)
				return nil
			}
			// fallback: print type info
			fmt.Printf("Plugin '%s' loaded.\n", args[0])
			return nil
		},
	}
}

func newPluginEnableCmd(ctx *core.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "enable [type] [name]",
		Short: "Enable a plugin",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr := getPluginManager(ctx)
			pluginType := core.ComponentType(args[0])
			name := args[1]
			if err := mgr.Load(context.Background(), pluginType, name); err != nil {
				return fmt.Errorf("failed to enable plugin: %w", err)
			}
			fmt.Printf("Plugin '%s' enabled.\n", name)
			return nil
		},
	}
}

func newPluginDisableCmd(ctx *core.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:   "disable [name]",
		Short: "Disable a plugin",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr := getPluginManager(ctx)
			if err := mgr.Unload(context.Background(), args[0]); err != nil {
				return fmt.Errorf("failed to disable plugin: %w", err)
			}
			fmt.Printf("Plugin '%s' disabled.\n", args[0])
			return nil
		},
	}
}
