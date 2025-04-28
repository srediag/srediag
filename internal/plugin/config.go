package plugin

import (
	"fmt"

	"github.com/srediag/srediag/internal/core"
)

// Package plugin provides plugin management and configuration logic for SREDIAG plugins.
//
// This file defines the configuration structures and helpers for the plugin manager, including loading and overlaying YAML config.
//
// Usage:
//   - Use PluginManagerConfig to represent the canonical plugin manager config structure.
//   - Use LoadPluginManagerConfig to load and overlay plugin manager config from YAML, environment, and CLI flags.
//
// Best Practices:
//   - Always validate required fields after loading config.
//   - Use canonical types for all new code.
//
// TODO:
//   - Add context.Context support for cancellation and timeouts.
//   - Add stricter schema validation and error reporting.

// PluginManagerConfig represents the canonical plugin manager config structure.
//
// Usage:
//   - Use this struct to load and validate plugin manager configuration for plugin orchestration and management operations.
//   - All config loading should use LoadPluginManagerConfig for schema compliance and validation.
type PluginManagerConfig struct {
	Plugins struct {
		Dir     string   `yaml:"dir"`
		ExecDir string   `yaml:"exec_dir"`
		Enabled []string `yaml:"enabled"`
	} `yaml:"plugins"`
}

// LoadPluginManagerConfig loads and overlays the plugin manager config from configs/srediag.yaml.
//
// Precedence: CLI flags > env > YAML > built-ins.
//
// Parameters:
//   - cliFlags: Map of CLI flag overrides.
//
// Returns:
//   - *PluginManagerConfig: The loaded plugin manager configuration.
//   - error: If loading or validation fails, returns a detailed error.
func LoadPluginManagerConfig(cliFlags map[string]string) (*PluginManagerConfig, error) {
	var cfg PluginManagerConfig
	if err := core.LoadConfigWithOverlay(&cfg, cliFlags, core.WithConfigPathSuffix("plugin")); err != nil {
		return nil, fmt.Errorf("failed to load plugin manager config: %w", err)
	}
	return &cfg, nil
}
