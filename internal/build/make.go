package build

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/srediag/srediag/internal/core"
)

// MakeBuilder handles the build process integrating with Makefile
type MakeBuilder struct {
	logger        *core.Logger
	workDir       string
	configPath    string
	outputDir     string
	pluginBuilder IBuilder
}

// NewMakeBuilder creates a new make builder
func NewMakeBuilder(logger *core.Logger, workDir, configPath, outputDir string) *MakeBuilder {
	return &MakeBuilder{
		logger:        logger,
		workDir:       workDir,
		configPath:    configPath,
		outputDir:     outputDir,
		pluginBuilder: NewPluginBuilder(logger, configPath, outputDir),
	}
}

// BuildAll builds the main project and all plugins
func (b *MakeBuilder) BuildAll() error {
	// First build plugins to ensure they're available
	if err := b.pluginBuilder.BuildAll(); err != nil {
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

// InstallPlugins copies built plugins to the installation directory
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

// BuildPlugin builds a single plugin and the main project
func (b *MakeBuilder) BuildPlugin(name string, cfg PluginConfig, compType core.ComponentType) error {
	// Build single plugin
	if err := b.pluginBuilder.BuildPlugin(name, cfg, compType); err != nil {
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
