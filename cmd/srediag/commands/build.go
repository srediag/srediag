// Package commands provides the command-line interface for the SREDIAG application.
package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/srediag/srediag/internal/build"
	"github.com/srediag/srediag/internal/core"
)

// NewBuildCmd creates the root build command with subcommands for build management.
// Only CLI wiring is present here; all business logic is delegated to internal/build CLI_* functions.
func NewBuildCmd(ctx *core.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "build",
		Short:   "Build agent and plugins (wraps otelcol-builder)",
		Long:    `Wraps otelcol-builder and helper scripts to compile the agent or individual plugins.\n\nArtifacts:\n  * Agent   → <output-dir>/srediag\n  * Plugins → <output-dir>/plugins/<name>/<name> (copy to plugins.exec_dir or use install command)`,
		Example: `  srediag build all --build-config build/srediag-build.yaml\n  srediag build plugin --type exporter --name clickhouseexporter\n  srediag build generate --type processor --name myprocessor`,
	}

	// Persistent flags for build config and output-dir, as per build.md
	cmd.PersistentFlags().String("build-config", "build/srediag-build.yaml", "Path to builder YAML (env: SREDIAG_BUILD_CONFIG)")
	cmd.PersistentFlags().String("output-dir", "/tmp/srediag-build", "Where artefacts are stored (env: SREDIAG_BUILD_OUTPUT_DIR)")
	if err := viper.BindPFlag("build.config", cmd.PersistentFlags().Lookup("build-config")); err != nil {
		return nil
	}
	if err := viper.BindPFlag("build.output_dir", cmd.PersistentFlags().Lookup("output-dir")); err != nil {
		return nil
	}
	if err := viper.BindEnv("build.config", "SREDIAG_BUILD_CONFIG"); err != nil {
		return nil
	}
	if err := viper.BindEnv("build.output_dir", "SREDIAG_BUILD_OUTPUT_DIR"); err != nil {
		return nil
	}

	cmd.AddCommand(
		newBuildAllCmd(ctx),
		newBuildGenerateCmd(ctx),
		newBuildInstallCmd(ctx),
		newBuildPluginCmd(ctx),
		newBuildUpdateCmd(ctx),
	)

	return cmd
}

// newBuildAllCmd wires the 'all' subcommand to build.CLI_BuildAll.
func newBuildAllCmd(ctx *core.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:     "all",
		Short:   "Build agent and every plugin declared in the YAML",
		Long:    "Builds the main agent and all plugins as defined in the builder YAML.",
		Example: "srediag build all --build-config build/srediag-build.yaml",
		RunE: func(cmd *cobra.Command, args []string) error {
			return build.CLI_BuildAll(ctx, cmd, args)
		},
	}
}

// newBuildPluginCmd wires the 'plugin' subcommand to build.CLI_BuildPlugin.
func newBuildPluginCmd(ctx *core.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "plugin",
		Short:   "Build a single plugin (requires --type and --name)",
		Long:    "Builds a single plugin by type and name as defined in the builder YAML.",
		Example: "srediag build plugin --type exporter --name clickhouseexporter",
		RunE: func(cmd *cobra.Command, args []string) error {
			return build.CLI_BuildPlugin(ctx, cmd, args)
		},
	}
	cmd.Flags().String("type", "", "Plugin type (receiver, processor, exporter, extension, connector)")
	cmd.Flags().String("name", "", "Plugin name")
	return cmd
}

// newBuildGenerateCmd wires the 'generate' subcommand to build.CLI_Generate.
func newBuildGenerateCmd(ctx *core.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate",
		Short:   "Produce plugin scaffold code (no compile)",
		Long:    "Generates Go code scaffolding for plugins. If --type and --name are provided, only that plugin is generated; otherwise, all plugins are generated.",
		Example: "srediag build generate --type processor --name myprocessor",
		RunE: func(cmd *cobra.Command, args []string) error {
			return build.CLI_Generate(ctx, cmd, args)
		},
	}
	cmd.Flags().String("type", "", "Plugin type (receiver, processor, exporter, extension, connector)")
	cmd.Flags().String("name", "", "Plugin name")
	return cmd
}

// newBuildInstallCmd wires the 'install' subcommand to build.CLI_InstallPlugins.
func newBuildInstallCmd(ctx *core.AppContext) *cobra.Command {
	return &cobra.Command{
		Use:     "install",
		Short:   "Copy pre-built plugins into plugins.exec_dir",
		Long:    "Copies all pre-built plugins from the output directory into the plugins.exec_dir directory.",
		Example: "srediag build install --output-dir /tmp/srediag-build",
		RunE: func(cmd *cobra.Command, args []string) error {
			return build.CLI_InstallPlugins(ctx, cmd, args)
		},
	}
}

// newBuildUpdateCmd wires the 'update' subcommand to build.CLI_UpdateBuilder.
func newBuildUpdateCmd(ctx *core.AppContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update",
		Short:   "Sync builder versions with go.mod",
		Long:    "Synchronise component versions between go.mod and builder, commenting out unknowns.",
		Example: "srediag build update --yaml configs/srediag-builder.yaml --gomod go.mod --plugin-gen plugin/generated",
		RunE: func(cmd *cobra.Command, args []string) error {
			return build.CLI_UpdateBuilder(ctx, cmd, args)
		},
	}
	cmd.Flags().String("yaml", "configs/srediag-builder.yaml", "Path to the builder YAML file")
	cmd.Flags().String("gomod", "go.mod", "Path to go.mod")
	cmd.Flags().String("plugin-gen", "plugin/generated", "Path to plugin/generated")
	return cmd
}
