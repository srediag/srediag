// Package commands provides the command-line interface for the SREDIAG application.
package commands

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/srediag/srediag/internal/build"
	"github.com/srediag/srediag/internal/core"
)

// NewBuildCmd creates the root build command with subcommands for build management.
func NewBuildCmd(ctx *core.AppContext) *cobra.Command {
	var outputDir, configPath string
	cmd := &cobra.Command{
		Use:   "build",
		Short: "Build plugins and main project",
	}

	cmd.PersistentFlags().StringVar(&outputDir, "output-dir", "/tmp/srediag-plugins", "Output directory for generated/built artifacts")
	cmd.PersistentFlags().StringVar(&configPath, "config", "otelcol-builder.yaml", "Path to the OpenTelemetry Collector builder configuration file")

	cmd.AddCommand(
		newBuildAllCmd(ctx, &outputDir, &configPath),
		newBuildPluginCmd(ctx, &outputDir, &configPath),
		newBuildGenerateCmd(ctx, &outputDir, &configPath),
		newBuildInstallCmd(ctx, &outputDir, &configPath),
		newBuildUpdateYAMLVersionsCmd(),
	)

	return cmd
}

func newBuildAllCmd(ctx *core.AppContext, outputDir, configPath *string) *cobra.Command {
	return &cobra.Command{
		Use:   "all",
		Short: "Build all plugins and the main project",
		RunE: func(cmd *cobra.Command, args []string) error {
			workDir, err := filepath.Abs(".")
			if err != nil {
				return fmt.Errorf("failed to get absolute path: %w", err)
			}
			makeBuilder := build.NewMakeBuilder(ctx.GetLogger(), workDir, *configPath, *outputDir)
			if err := makeBuilder.BuildAll(); err != nil {
				return fmt.Errorf("build failed: %w", err)
			}
			ctx.GetLogger().Info("Build completed successfully", core.ZapString("output_dir", *outputDir))
			return nil
		},
	}
}

func newBuildPluginCmd(ctx *core.AppContext, outputDir, configPath *string) *cobra.Command {
	var pluginType, pluginName string
	cmd := &cobra.Command{
		Use:   "plugin",
		Short: "Build a single plugin by type and name",
		RunE: func(cmd *cobra.Command, args []string) error {
			if pluginType == "" || pluginName == "" {
				return fmt.Errorf("--type and --name are required")
			}
			workDir, err := filepath.Abs(".")
			if err != nil {
				return fmt.Errorf("failed to get absolute path: %w", err)
			}
			makeBuilder := build.NewMakeBuilder(ctx.GetLogger(), workDir, *configPath, *outputDir)
			pluginBuilder := build.NewPluginBuilder(ctx.GetLogger(), *configPath, *outputDir)
			cfg, err := pluginBuilder.LoadConfig()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			compType := core.ComponentType(pluginType)
			pluginCfg, ok := cfg.Components[compType][pluginName]
			if !ok {
				return fmt.Errorf("plugin %s/%s not found in config", pluginType, pluginName)
			}
			if err := makeBuilder.BuildPlugin(pluginName, pluginCfg, compType); err != nil {
				return fmt.Errorf("build plugin: %w", err)
			}
			ctx.GetLogger().Info("Plugin built", core.ZapString("type", pluginType), core.ZapString("name", pluginName), core.ZapString("output_dir", *outputDir))
			return nil
		},
	}
	cmd.Flags().StringVar(&pluginType, "type", "", "Plugin type (receiver, processor, exporter, extension, connector)")
	cmd.Flags().StringVar(&pluginName, "name", "", "Plugin name")
	return cmd
}

func newBuildGenerateCmd(ctx *core.AppContext, outputDir, configPath *string) *cobra.Command {
	return &cobra.Command{
		Use:   "generate",
		Short: "Generate plugin code without compiling",
		RunE: func(cmd *cobra.Command, args []string) error {
			pluginBuilder := build.NewPluginBuilder(ctx.GetLogger(), *configPath, *outputDir)
			if err := pluginBuilder.GenerateAll(); err != nil {
				return fmt.Errorf("generate code: %w", err)
			}
			ctx.GetLogger().Info("Code generation completed", core.ZapString("output_dir", *outputDir))
			return nil
		},
	}
}

func newBuildInstallCmd(ctx *core.AppContext, outputDir, configPath *string) *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "Install pre-built plugins into the system",
		RunE: func(cmd *cobra.Command, args []string) error {
			workDir, err := filepath.Abs(".")
			if err != nil {
				return fmt.Errorf("failed to get absolute path: %w", err)
			}
			makeBuilder := build.NewMakeBuilder(ctx.GetLogger(), workDir, *configPath, *outputDir)
			if err := makeBuilder.InstallPlugins(); err != nil {
				return fmt.Errorf("install plugins: %w", err)
			}
			ctx.GetLogger().Info("Plugin install completed", core.ZapString("output_dir", *outputDir))
			return nil
		},
	}
}

func newBuildUpdateYAMLVersionsCmd() *cobra.Command {
	var yamlPath, goModPath, pluginGenDir string
	cmd := &cobra.Command{
		Use:   "update-yaml-versions",
		Short: "Synchronise component versions between go.mod and builder YAML",
		RunE: func(cmd *cobra.Command, args []string) error {
			if yamlPath == "" || goModPath == "" || pluginGenDir == "" {
				return fmt.Errorf("all flags --yaml, --gomod and --plugin-gen are required")
			}
			if err := build.UpdateYAMLVersions(yamlPath, goModPath, pluginGenDir); err != nil {
				return err
			}
			fmt.Println("YAML updated and summary printed.")
			return nil
		},
	}
	cmd.Flags().StringVar(&yamlPath, "yaml", "configs/srediag-builder.yaml", "Path to the builder YAML file")
	cmd.Flags().StringVar(&goModPath, "gomod", "go.mod", "Path to go.mod")
	cmd.Flags().StringVar(&pluginGenDir, "plugin-gen", "plugin/generated", "Path to plugin/generated")
	return cmd
}
