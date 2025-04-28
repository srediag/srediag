// Package build provides the build orchestration layer for SREDIAG.
//
// This file defines the CLI entrypoints for build operations, wiring Cobra commands to internal build logic.
// Each function extracts parameters from the CLI context, instantiates the BuildManager, and delegates to the appropriate method.
//
// Usage:
//   - These functions are called by the CLI wiring in cmd/srediag/commands/build.go.
//   - They are not intended for direct use outside CLI command handlers.
//
// Best Practices:
//   - Always validate required flags and parameters before calling BuildManager methods.
//   - Log all errors and important events for traceability.
//   - Use context-aware logging and error handling for better diagnostics.
//
// TODO: Add support for context.Context to all CLI functions to enable cancellation and timeouts.
// TODO: Refactor to reduce repeated logger fallback logic in CLI functions.
// TODO: Consider extracting common flag extraction logic to helper functions.
// TODO: Implement 'srediag build all' and 'srediag build plugin' commands using the YAML specification.
// TODO: Output static binary and plugin bundles as described in the build specification.
// TODO: Integrate SBOM generation, cosign signing, and SLSA attestation into the build pipeline.
// TODO: Implement standardized error codes for the build CLI.
// TODO: Provide actionable error messages for common build pitfalls in the CLI.
//
// Redundancy/Refactor:
//   - Logger fallback and outputDir extraction are repeated in every function; could be DRYed up.
//   - All CLI entrypoints follow a similar pattern; consider a generic wrapper for error/log handling.
//
// TODO: Implement 'srediag build all' and 'srediag build plugin' commands using YAML spec (see architecture/build.md §2)
// TODO: Output static binary and plugin bundles as described (see architecture/build.md §2)
// TODO: Integrate SBOM, cosign signing, and SLSA attestation in build pipeline (see architecture/build.md §3)
// TODO: Implement standardized error codes for build CLI (see architecture/build.md §7)
// TODO: Provide actionable error messages for common build pitfalls (see architecture/build.md §7)
// TODO: Implement granular validation for all nested config sections.
// TODO: Add a Redact/Sanitize method to Config to prevent logging secrets.
// TODO: Integrate StrictYAMLUnmarshal into LoadConfigWithOverlay for schema safety.
// TODO: Add support for multi-error reporting in ValidateConfig.
// TODO: Add tests for environment variable overlays and malformed YAML.
// TODO: Allow explicit override of default directories for containerized or custom environments.
package build

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/srediag/srediag/internal/core"
)

var updateBuilderFunc = UpdateBuilder
var newBuildManagerFunc = func(logger *core.Logger, outputDir string) BuildManagerInterface {
	return NewBuildManager(logger, outputDir)
}

// CLI_BuildAll is the entrypoint for 'srediag build all'.
//
// Usage:
//   - Called by the 'all' subcommand to build the agent and all plugins.
//   - Extracts config/output-dir from flags/env, instantiates BuildManager, and calls BuildAll.
//
// Best Practices:
//   - Always log errors and completion status.
//   - Validate outputDir before use.
//
// TODO:
//   - Add support for context.Context for cancellation.
//   - Refactor logger fallback to a helper.
func CLI_BuildAll(ctx *core.AppContext, cmd *cobra.Command, args []string) error {
	logger := ctx.Logger
	if logger == nil {
		var err error
		logger, err = core.NewLogger(nil)
		if err != nil {
			return fmt.Errorf("failed to create fallback logger: %w", err)
		}
	}
	outputDir := viper.GetString("build.output_dir")
	if outputDir == "" {
		outputDir = "/tmp/srediag-build"
	}
	bm := newBuildManagerFunc(logger, outputDir)
	if err := bm.BuildAll(); err != nil {
		logger.Error("BuildAll failed", core.ZapError(err))
		return fmt.Errorf("build all failed: %w", err)
	}
	logger.Info("BuildAll completed successfully")
	return nil
}

// CLI_BuildPlugin is the entrypoint for 'srediag build plugin'.
//
// Usage:
//   - Called by the 'plugin' subcommand to build a single plugin.
//   - Extracts type/name/output-dir from flags/env, instantiates BuildManager, and calls BuildPlugin.
//
// Best Practices:
//   - Always validate required flags (type, name).
//   - Log errors and completion status.
//
// TODO:
//   - Add support for context.Context for cancellation.
//   - Refactor logger fallback to a helper.
func CLI_BuildPlugin(ctx *core.AppContext, cmd *cobra.Command, args []string) error {
	logger := ctx.Logger
	if logger == nil {
		var err error
		logger, err = core.NewLogger(nil)
		if err != nil {
			return fmt.Errorf("failed to create fallback logger: %w", err)
		}
	}
	outputDir := viper.GetString("build.output_dir")
	if outputDir == "" {
		outputDir = "/tmp/srediag-build"
	}
	pluginType, err := cmd.Flags().GetString("type")
	if err != nil || pluginType == "" {
		return fmt.Errorf("--type flag is required")
	}
	pluginName, err := cmd.Flags().GetString("name")
	if err != nil || pluginName == "" {
		return fmt.Errorf("--name flag is required")
	}
	bm := newBuildManagerFunc(logger, outputDir)
	if err := bm.BuildPlugin(pluginType, pluginName); err != nil {
		logger.Error("BuildPlugin failed", core.ZapError(err))
		return fmt.Errorf("build plugin failed: %w", err)
	}
	logger.Info("BuildPlugin completed successfully", core.ZapString("type", pluginType), core.ZapString("name", pluginName))
	return nil
}

// CLI_Generate is the entrypoint for 'srediag build generate'.
//
// Usage:
//   - Called by the 'generate' subcommand to scaffold plugin code.
//   - Extracts type/name/output-dir from flags/env, instantiates BuildManager, and calls Generate.
//
// Best Practices:
//   - Log errors and completion status.
//   - Validate pluginType and pluginName if required by business logic.
//
// TODO:
//   - Add support for context.Context for cancellation.
//   - Refactor logger fallback to a helper.
func CLI_Generate(ctx *core.AppContext, cmd *cobra.Command, args []string) error {
	logger := ctx.Logger
	if logger == nil {
		var err error
		logger, err = core.NewLogger(nil)
		if err != nil {
			return fmt.Errorf("failed to create fallback logger: %w", err)
		}
	}
	outputDir := viper.GetString("build.output_dir")
	if outputDir == "" {
		outputDir = "/tmp/srediag-build"
	}
	pluginType, _ := cmd.Flags().GetString("type")
	pluginName, _ := cmd.Flags().GetString("name")
	bm := newBuildManagerFunc(logger, outputDir)
	if err := bm.Generate(pluginType, pluginName); err != nil {
		logger.Error("Generate failed", core.ZapError(err))
		return fmt.Errorf("generate failed: %w", err)
	}
	logger.Info("Generate completed successfully", core.ZapString("type", pluginType), core.ZapString("name", pluginName))
	return nil
}

// CLI_InstallPlugins is the entrypoint for 'srediag build install'.
//
// Usage:
//   - Called by the 'install' subcommand to copy pre-built plugins to the exec dir.
//   - Extracts output-dir from flags/env, instantiates BuildManager, and calls InstallPlugins.
//
// Best Practices:
//   - Log errors and completion status.
//   - Validate outputDir before use.
//
// TODO:
//   - Add support for context.Context for cancellation.
//   - Refactor logger fallback to a helper.
func CLI_InstallPlugins(ctx *core.AppContext, cmd *cobra.Command, args []string) error {
	logger := ctx.Logger
	if logger == nil {
		var err error
		logger, err = core.NewLogger(nil)
		if err != nil {
			return fmt.Errorf("failed to create fallback logger: %w", err)
		}
	}
	outputDir := viper.GetString("build.output_dir")
	if outputDir == "" {
		outputDir = "/tmp/srediag-build"
	}
	bm := newBuildManagerFunc(logger, outputDir)
	if err := bm.InstallPlugins(); err != nil {
		logger.Error("InstallPlugins failed", core.ZapError(err))
		return fmt.Errorf("install plugins failed: %w", err)
	}
	logger.Info("InstallPlugins completed successfully")
	return nil
}

// CLI_UpdateBuilder is the entrypoint for 'srediag build update'.
//
// Usage:
//   - Called by the 'update' subcommand to sync builder YAML with go.mod.
//   - Extracts yaml, gomod, and plugin-gen flags, and calls UpdateBuilder.
//
// Best Practices:
//   - Always validate required flags (yaml, gomod, plugin-gen).
//   - Log errors and completion status.
//
// TODO:
//   - Add support for context.Context for cancellation.
//   - Refactor logger fallback to a helper.
func CLI_UpdateBuilder(ctx *core.AppContext, cmd *cobra.Command, args []string) error {
	logger := ctx.Logger
	if logger == nil {
		var err error
		logger, err = core.NewLogger(nil)
		if err != nil {
			return fmt.Errorf("failed to create fallback logger: %w", err)
		}
	}
	yamlPath, err := cmd.Flags().GetString("yaml")
	if err != nil || yamlPath == "" {
		return fmt.Errorf("--yaml flag is required")
	}
	gomodPath, err := cmd.Flags().GetString("gomod")
	if err != nil || gomodPath == "" {
		return fmt.Errorf("--gomod flag is required")
	}
	pluginGenDir, err := cmd.Flags().GetString("plugin-gen")
	if err != nil || pluginGenDir == "" {
		return fmt.Errorf("--plugin-gen flag is required")
	}
	if err := updateBuilderFunc(yamlPath, gomodPath, pluginGenDir); err != nil {
		logger.Error("UpdateBuilder failed", core.ZapError(err))
		return fmt.Errorf("update builder failed: %w", err)
	}
	logger.Info("UpdateBuilder completed successfully")
	return nil
}
