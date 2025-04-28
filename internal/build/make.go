package build

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/srediag/srediag/internal/core"
)

// Package build provides the build orchestration layer for SREDIAG.
//
// This file defines the MakeBuilder, which integrates the build process with Makefile-based workflows.
// MakeBuilder uses BuildManager for plugin orchestration and provides methods for building and installing plugins.
//
// Usage:
//   - Use MakeBuilder to coordinate plugin builds and main project builds via Makefile.
//   - All plugin orchestration is delegated to BuildManager.
//   - All config loading should use LoadBuildConfig for schema compliance and validation.
//
// Best Practices:
//   - Always check for errors from build and install methods.
//   - Use logger for all error and status reporting.
//   - Prefer environment variable SREDIAG_PLUGIN_DIR to override default plugin install location.
//
// TODO:
//   - Add context.Context to all methods for cancellation and timeouts.
//   - Add more granular error reporting for build and install steps.

// MakeBuilder handles the build process integrating with Makefile.
//
// MakeBuilder is responsible for:
//   - Building all plugins and the main project using Makefile targets.
//   - Delegating plugin orchestration to BuildManager.
//   - Installing built plugins to the appropriate directory.
//
// Usage:
//   - Instantiate with NewMakeBuilder, providing logger, working directory, config path, and output directory.
//   - Call BuildAll to build all plugins and the main project.
//   - Call BuildPlugin to build a single plugin and the main project.
//   - Call InstallPlugins to copy built plugins to the install directory.
type MakeBuilder struct {
	logger       *core.Logger
	workDir      string
	configPath   string
	outputDir    string
	buildManager *BuildManager
}

// NewMakeBuilder creates a new MakeBuilder.
//
// Parameters:
//   - logger: Logger for status and error reporting.
//   - workDir: Working directory for Makefile operations.
//   - configPath: Path to build configuration file.
//   - outputDir: Directory where build artifacts are placed.
//
// Returns:
//   - *MakeBuilder: A new MakeBuilder instance.
func NewMakeBuilder(logger *core.Logger, workDir, configPath, outputDir string) *MakeBuilder {
	return &MakeBuilder{
		logger:       logger,
		workDir:      workDir,
		configPath:   configPath,
		outputDir:    outputDir,
		buildManager: NewBuildManager(logger, outputDir),
	}
}

// BuildAll builds the main project and all plugins.
//
// This method:
//   - Builds all plugins using BuildManager.BuildAll.
//   - Runs 'make clean' and 'make build' in the working directory.
//   - Installs built plugins to the install directory.
//
// Returns:
//   - error: If any build or install step fails, returns a detailed error.
//
// Side Effects:
//   - Modifies files in the output and install directories.
//   - Runs external Makefile commands.
func (b *MakeBuilder) BuildAll() error {
	// First build plugins to ensure they're available
	if err := b.buildManager.BuildAll(); err != nil {
		b.logger.Error("Failed to build plugins", core.ZapError(err))
		return err
	}

	// Run make clean to ensure clean build
	cleanCmd := exec.Command("make", "clean")
	cleanCmd.Dir = b.workDir
	cleanCmd.Stdout = os.Stdout
	cleanCmd.Stderr = os.Stderr
	if err := cleanCmd.Run(); err != nil {
		b.logger.Error("Failed to clean project", core.ZapError(err))
		return fmt.Errorf("make clean failed: %w", err)
	}

	// Run make build
	buildCmd := exec.Command("make", "build")
	buildCmd.Dir = b.workDir
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		b.logger.Error("Failed to build project", core.ZapError(err))
		return fmt.Errorf("make build failed: %w", err)
	}

	// Copy plugins to the appropriate directory
	if err := b.InstallPlugins(); err != nil {
		b.logger.Error("Failed to install plugins", core.ZapError(err))
		return err
	}

	return nil
}

// InstallPlugins copies built plugins to the installation directory.
//
// This method:
//   - Determines the install directory (default: /etc/srediag/plugins, overridable by SREDIAG_PLUGIN_DIR).
//   - Copies all .so files from the output directory to the install directory.
//   - Logs all errors and successful installations.
//
// Returns:
//   - error: If directory creation, reading, or file copy fails, returns a detailed error.
//
// Side Effects:
//   - Creates directories and copies files on the filesystem.
func (b *MakeBuilder) InstallPlugins() error {
	// Default plugin installation directory
	installDir := filepath.Join("/etc/srediag", "plugins")
	if envDir := os.Getenv("SREDIAG_PLUGIN_DIR"); envDir != "" {
		installDir = envDir
	}

	// Create installation directory if it doesn't exist
	if err := os.MkdirAll(installDir, 0755); err != nil {
		return fmt.Errorf("failed to create plugin installation directory: %w", err)
	}

	// Copy all .so files from output directory to installation directory
	entries, err := os.ReadDir(b.outputDir)
	if err != nil {
		return fmt.Errorf("failed to read output directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".so" {
			srcPath := filepath.Join(b.outputDir, entry.Name())
			dstPath := filepath.Join(installDir, entry.Name())

			// Read source file
			data, err := os.ReadFile(srcPath)
			if err != nil {
				b.logger.Error("Failed to read plugin file",
					core.ZapString("file", srcPath),
					core.ZapError(err))
				continue
			}

			// Write to destination
			if err := os.WriteFile(dstPath, data, 0644); err != nil {
				b.logger.Error("Failed to write plugin file",
					core.ZapString("file", dstPath),
					core.ZapError(err))
				continue
			}

			b.logger.Info("Installed plugin",
				core.ZapString("name", entry.Name()),
				core.ZapString("path", dstPath))
		}
	}

	return nil
}

// BuildPlugin builds a single plugin and the main project.
//
// This method:
//   - Builds a single plugin using BuildManager.BuildPlugin.
//   - Runs 'make build' in the working directory.
//   - Installs the built plugin to the install directory.
//
// Parameters:
//   - name: Name of the plugin to build.
//   - compType: Component type of the plugin.
//
// Returns:
//   - error: If any build or install step fails, returns a detailed error.
//
// Side Effects:
//   - Modifies files in the output and install directories.
//   - Runs external Makefile commands.
func (b *MakeBuilder) BuildPlugin(name string, compType core.ComponentType) error {
	// Build single plugin using BuildManager
	if err := b.buildManager.BuildPlugin(string(compType), name); err != nil {
		b.logger.Error("Failed to build plugin",
			core.ZapString("name", name),
			core.ZapError(err))
		return err
	}

	// Run make build
	buildCmd := exec.Command("make", "build")
	buildCmd.Dir = b.workDir
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		b.logger.Error("Failed to build project", core.ZapError(err))
		return fmt.Errorf("make build failed: %w", err)
	}

	// Install the plugin
	if err := b.InstallPlugins(); err != nil {
		b.logger.Error("Failed to install plugin", core.ZapError(err))
		return err
	}

	return nil
}

// Note: YAML/go.mod sync is handled by UpdateBuilderYAMLVersions in update.go
// Note: All build config loading should use LoadBuildConfig for validation and schema compliance.
