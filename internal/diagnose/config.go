// Package diagnose provides diagnostic operations for SREDIAG, including system, performance, and security diagnostics.
//
// This file defines the configuration structures and helpers for diagnostics, including loading and overlaying YAML config.
//
// Usage:
//   - Use DiagnosticsConfig to represent the canonical diagnostics config structure.
//   - Use LoadDiagnosticsConfig to load and overlay diagnostics config from YAML, environment, and CLI flags.
//
// Best Practices:
//   - Always validate required fields after loading config.
//   - Use canonical types for all new code.
//
// TODO:
//   - Add context.Context support for cancellation and timeouts.
//   - Add stricter schema validation and error reporting.
package diagnose

import (
	"fmt"

	"github.com/srediag/srediag/internal/core"
)

// DiagnosticsConfig represents the canonical diagnostics config structure.
//
// Usage:
//   - Use this struct to load and validate diagnostics configuration for plugin orchestration and diagnostics operations.
//   - All config loading should use LoadDiagnosticsConfig for schema compliance and validation.
type DiagnosticsConfig struct {
	Defaults struct {
		OutputFormat string `yaml:"output_format"`
		Timeout      string `yaml:"timeout"`
		MaxRetries   int    `yaml:"max_retries"`
	} `yaml:"defaults"`
	Plugins map[string]interface{} `yaml:"plugins"` // Plugin configs are loaded from plugins.d/diagnostics
}

// LoadDiagnosticsConfig loads and overlays the diagnostics config from configs/srediag-diagnose.yaml.
//
// Precedence: CLI flags > env > YAML > built-ins.
//
// Parameters:
//   - cliFlags: Map of CLI flag overrides.
//
// Returns:
//   - *DiagnosticsConfig: The loaded diagnostics configuration.
//   - error: If loading or validation fails, returns a detailed error.
func LoadDiagnosticsConfig(cliFlags map[string]string) (*DiagnosticsConfig, error) {
	var cfg DiagnosticsConfig
	if err := core.LoadConfigWithOverlay(&cfg, cliFlags, core.WithConfigPathSuffix("diagnose")); err != nil {
		return nil, fmt.Errorf("failed to load diagnostics config: %w", err)
	}
	return &cfg, nil
}
