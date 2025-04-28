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
// This file defines the BuildManager, which is the central orchestrator for all build operations, including plugin builds, configuration loading, code generation, and plugin installation.
// BuildManager delegates low-level plugin build and code generation to the plugin.go module.
//
// Usage:
//   - Use BuildManager to coordinate all build-related operations in SREDIAG.
//   - Instantiate with NewBuildManager, providing a logger and output directory.
//   - All config loading should use LoadBuildConfig for schema compliance and validation.
//
// Best Practices:
//   - Always check for errors from build and install methods.
//   - Use logger for all error and status reporting.
//
// TODO:
//   - Implement all methods to perform actual build, generate, and install logic.
//   - Add context.Context to all methods for cancellation and timeouts.

// BuildManager is the central orchestrator for all build operations (plugin build, config load, etc).
//
// BuildManager is responsible for:
//   - Loading build configuration.
//   - Building all plugins or a single plugin.
//   - Generating plugin scaffold code.
//   - Installing built plugins to the installation directory.
//
// Usage:
//   - Instantiate with NewBuildManager, providing a logger and output directory.
//   - Call BuildAll to build all plugins.
//   - Call BuildPlugin to build a single plugin.
//   - Call Generate to scaffold plugin code.
//   - Call InstallPlugins to copy built plugins to the install directory.
type BuildManager struct {
	logger    *core.Logger
	outputDir string
}

// NewBuildManager creates a new BuildManager.
//
// Parameters:
//   - logger: Logger for status and error reporting.
//   - outputDir: Directory where build artifacts are placed.
//
// Returns:
//   - *BuildManager: A new BuildManager instance.
func NewBuildManager(logger *core.Logger, outputDir string) *BuildManager {
	return &BuildManager{
		logger:    logger,
		outputDir: outputDir,
	}
}

// LoadConfig loads the builder configuration using the unified loader.
//
// Returns:
//   - *BuilderConfig: The loaded builder configuration.
//   - error: If loading or validation fails, returns a detailed error.
func (m *BuildManager) LoadConfig() (*BuilderConfig, error) {
	return LoadBuildConfig(map[string]string{})
}

// BuildAll builds all plugins defined in the configuration.
//
// Returns:
//   - error: If any build step fails, returns a detailed error.
//
// Side Effects:
//   - Intended to modify files in the output directory (not yet implemented).
func (m *BuildManager) BuildAll() error {
	cfg, err := LoadBuildConfig(nil)
	if err != nil {
		return fmt.Errorf("failed to load build config: %w", err)
	}
	outputDir := cfg.Dist.OutputPath
	if outputDir == "" {
		outputDir = "bin"
	}
	pluginRoot := filepath.Join(outputDir, "plugins")
	failed := false

	for _, typ := range []core.ComponentType{core.TypeReceiver, core.TypeProcessor, core.TypeExporter, core.TypeExtension} {
		for name, plugin := range cfg.Components[typ] {
			if plugin.Path == "" {
				m.logger.Error("No path specified for plugin", core.ZapString("type", string(typ)), core.ZapString("name", name))
				failed = true
				continue
			}
			outDir := filepath.Join(pluginRoot, string(typ), name)
			if err := os.MkdirAll(outDir, 0755); err != nil {
				m.logger.Error("Failed to create dir", core.ZapString("dir", outDir), core.ZapError(err))
				failed = true
				continue
			}
			outFile := filepath.Join(outDir, name+".so")
			cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", outFile, plugin.Path)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			m.logger.Info("Building plugin", core.ZapString("type", string(typ)), core.ZapString("name", name))
			if err := cmd.Run(); err != nil {
				m.logger.Error("Failed to build plugin", core.ZapString("type", string(typ)), core.ZapString("name", name), core.ZapError(err))
				failed = true
				continue
			}
			m.logger.Info("Built plugin", core.ZapString("type", string(typ)), core.ZapString("name", name), core.ZapString("output", outFile))
		}
	}

	if failed {
		return fmt.Errorf("one or more plugins failed to build")
	}
	m.logger.Info("All plugins built successfully", core.ZapString("dir", pluginRoot))
	return nil
}

// BuildPlugin builds a single plugin by name and type.
//
// Parameters:
//   - pluginType: The type of the plugin to build.
//   - pluginName: The name of the plugin to build.
//
// Returns:
//   - error: If the build fails or plugin is not found, returns a detailed error.
//
// Side Effects:
//   - Intended to modify files in the output directory (not yet implemented).
func (m *BuildManager) BuildPlugin(pluginType, pluginName string) error {
	cfg, err := LoadBuildConfig(nil)
	if err != nil {
		return fmt.Errorf("failed to load build config: %w", err)
	}
	outputDir := cfg.Dist.OutputPath
	if outputDir == "" {
		outputDir = "bin"
	}
	pluginRoot := filepath.Join(outputDir, "plugins")
	typ := core.ComponentType(pluginType)
	plugin, ok := cfg.Components[typ][pluginName]
	if !ok {
		return fmt.Errorf("plugin %s/%s not found in config", pluginType, pluginName)
	}
	if plugin.Path == "" {
		return fmt.Errorf("no path specified for plugin %s/%s", pluginType, pluginName)
	}
	outDir := filepath.Join(pluginRoot, pluginType, pluginName)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("failed to create dir %s: %w", outDir, err)
	}
	outFile := filepath.Join(outDir, pluginName+".so")
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", outFile, plugin.Path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	m.logger.Info("Building plugin", core.ZapString("type", pluginType), core.ZapString("name", pluginName))
	if err := cmd.Run(); err != nil {
		m.logger.Error("Failed to build plugin", core.ZapString("type", pluginType), core.ZapString("name", pluginName), core.ZapError(err))
		return fmt.Errorf("failed to build plugin %s/%s: %w", pluginType, pluginName, err)
	}
	m.logger.Info("Built plugin", core.ZapString("type", pluginType), core.ZapString("name", pluginName), core.ZapString("output", outFile))
	return nil
}

// Generate produces plugin scaffold code (no compile).
//
// Parameters:
//   - pluginType: The type of the plugin to generate.
//   - pluginName: The name of the plugin to generate.
//
// Returns:
//   - error: If generation fails, returns a detailed error.
//
// Side Effects:
//   - Intended to create scaffold files in the output directory (not yet implemented).
func (m *BuildManager) Generate(pluginType, pluginName string) error {
	return fmt.Errorf("generate not yet implemented in plugin.go")
}

// InstallPlugins copies built plugins to the installation directory.
//
// Returns:
//   - error: If installation fails, returns a detailed error.
//
// Side Effects:
//   - Intended to copy files to the install directory (not yet implemented).
func (m *BuildManager) InstallPlugins() error {
	return fmt.Errorf("installPlugins not yet implemented in plugin.go")
}
