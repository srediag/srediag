// Package build provides the build orchestration layer for SREDIAG.
//
// This file defines the configuration types and helpers for the build system, including plugin and builder config schemas.
// It also provides helpers for loading, validating, and converting build configs from legacy to canonical formats.
//
// Usage:
//   - Use these types to load and validate build configuration for plugin orchestration and builder operations.
//   - All config loading should use LoadBuildConfig for schema compliance and validation.
//
// Best Practices:
//   - Always validate required fields after loading config.
//   - Use canonical types (BuilderConfig, PluginConfig) for all new code.
//   - Avoid direct use of LegacyConfig except for migration/compatibility.
//
// TODO: Remove LegacyConfig after the migration to canonical config is complete.
// TODO: Add support for context.Context to all methods for cancellation and timeouts.
// TODO: Add stricter schema validation and improve error reporting for build configuration.
// TODO(C-01 Phase 0): Upgrade to OTel v0.124.0 / API v1.30.0 (see TODO.md C-01, ETA 2025-05-10)
// TODO(C-02 Phase 0): Implement pipeline builder that converts Go configuration to YAML (see TODO.md C-02, ETA 2025-05-24)
// TODO: Enforce that all components in build YAML use the exact Go module path and pinned version.
// TODO: Fail service startup if any unrecognized components are found in the pipeline YAML.
// TODO: Implement a utility to synchronize versions between go.mod and build YAML.
// TODO: Validate module versions and provide a diff/apply feature with --write flag.
// TODO: Enforce exact Go module path and pinned version for all components in build YAML (see docs/architecture/build.md §1.1)
// TODO: Fail service startup if unrecognized components appear in pipeline YAML (see docs/architecture/build.md §1.1)
// TODO: Implement version synchronization utility between go.mod and build YAML (see docs/architecture/build.md §4)
// TODO: Validate module versions and generate diff/apply with --write (see docs/architecture/build.md §4)
package build

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v3"

	"github.com/srediag/srediag/internal/core"
)

// Note: YAML/go.mod sync is handled by UpdateBuilderYAMLVersions in update.go

// PluginConfig represents a module/component in the build config (extensible).
//
// Usage:
//   - Used as the canonical config for a plugin or component in BuilderConfig.
//   - Use for all plugin orchestration and build operations.
type PluginConfig struct {
	GoMod  string   `yaml:"gomod"`
	Import string   `yaml:"import,omitempty"`
	Path   string   `yaml:"path,omitempty"`
	Tags   []string `yaml:"tags,omitempty"`
}

// BuilderConfig represents the otelcol-builder.yaml configuration.
// Used as the canonical build config for plugin orchestration and validated/converted from LegacyConfig.
//
// Usage:
//   - Use for all new build orchestration and plugin management code.
//   - Always validate required fields after loading.
type BuilderConfig struct {
	Dist struct {
		Name           string `yaml:"name"`            // Name of the distribution
		OutputPath     string `yaml:"output_path"`     // Output path for built artifacts
		OtelColVersion string `yaml:"otelcol_version"` // OpenTelemetry Collector version
	} `yaml:"dist"`
	Components map[core.ComponentType]map[string]PluginConfig `yaml:"components"` // Map of component type to component configurations
}

// LegacyConfig represents the current otelcol-builder.yaml format.
//
// Usage:
//   - Transitional struct for migration from legacy to canonical config.
//   - Avoid in new code; use BuilderConfig instead.
type LegacyConfig struct {
	Dist struct {
		Module           string `yaml:"module"`
		Name             string `yaml:"name"`
		Description      string `yaml:"description"`
		OutputPath       string `yaml:"output_path"`
		Version          string `yaml:"version"`
		DebugCompilation bool   `yaml:"debug_compilation"`
	} `yaml:"dist"`
	Exporters  []ModuleConfig `yaml:"exporters"`
	Extensions []ModuleConfig `yaml:"extensions"`
	Receivers  []ModuleConfig `yaml:"receivers"`
	Processors []ModuleConfig `yaml:"processors"`
	Providers  []ModuleConfig `yaml:"providers"`
}

// ModuleConfig represents a module/component in the build config (legacy/extensible).
//
// Usage:
//   - Transitional struct for migration from legacy to canonical config.
//   - Avoid in new code; use PluginConfig instead.
type ModuleConfig struct {
	GoMod  string   `yaml:"gomod"`
	Import string   `yaml:"import,omitempty"`
	Path   string   `yaml:"path,omitempty"`
	Tags   []string `yaml:"tags,omitempty"`
}

// LoadBuildConfig loads, validates, and converts the build config for the builder.
// It uses core.LoadConfigWithOverlay for discovery and precedence (flags > env > YAML > built-ins).
// The canonical search root is the configs/ directory, matching the modular structure.
//
// Usage:
//   - Use to load and validate build config for all builder operations.
//   - Returns a canonical BuilderConfig for plugin orchestration.
//
// Best Practices:
//   - Always check for required fields (dist.name, dist.version, dist.output_path, at least one component).
//   - Use the returned BuilderConfig for all downstream build logic.
//
// TODO:
//   - Add stricter schema validation and error reporting.
//   - Remove legacy conversion after migration.
func LoadBuildConfig(cliFlags map[string]string) (*BuilderConfig, error) {
	var raw map[string]interface{}
	if err := core.LoadConfigWithOverlay(&raw, cliFlags, core.WithConfigPathSuffix("build")); err != nil {
		return nil, err
	}

	// Try canonical format first
	b, err := yaml.Marshal(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config for canonical check: %w", err)
	}
	var canonical BuilderConfig
	if err := yaml.Unmarshal(b, &canonical); err == nil {
		dist := canonical.Dist
		if dist.Name != "" && dist.OtelColVersion != "" && dist.OutputPath != "" && len(canonical.Components) > 0 {
			return &canonical, nil
		}
	}

	// Fallback: try legacy format
	var legacy LegacyConfig
	if err := yaml.Unmarshal(b, &legacy); err != nil {
		return nil, fmt.Errorf("failed to parse build config as canonical or legacy: %w", err)
	}
	if legacy.Dist.Name == "" || legacy.Dist.Version == "" || legacy.Dist.OutputPath == "" {
		return nil, fmt.Errorf("dist.name, dist.version, and dist.output_path are required in build config")
	}
	if len(legacy.Receivers)+len(legacy.Processors)+len(legacy.Exporters)+len(legacy.Extensions) == 0 {
		return nil, fmt.Errorf("at least one component (receiver, processor, exporter, extension) must be defined in build config")
	}

	config := &BuilderConfig{}
	config.Dist.Name = legacy.Dist.Name
	config.Dist.OutputPath = legacy.Dist.OutputPath
	config.Dist.OtelColVersion = legacy.Dist.Version
	config.Components = make(map[core.ComponentType]map[string]PluginConfig)
	config.Components[core.TypeReceiver] = make(map[string]PluginConfig)
	config.Components[core.TypeProcessor] = make(map[string]PluginConfig)
	config.Components[core.TypeExporter] = make(map[string]PluginConfig)
	config.Components[core.TypeExtension] = make(map[string]PluginConfig)

	// Convert receivers
	for _, r := range legacy.Receivers {
		name := extractComponentName(r.GoMod)
		config.Components[core.TypeReceiver][name] = PluginConfig(r)
	}
	// Convert processors
	for _, p := range legacy.Processors {
		name := extractComponentName(p.GoMod)
		config.Components[core.TypeProcessor][name] = PluginConfig(p)
	}
	// Convert exporters
	for _, e := range legacy.Exporters {
		name := extractComponentName(e.GoMod)
		config.Components[core.TypeExporter][name] = PluginConfig(e)
	}
	// Convert extensions
	for _, e := range legacy.Extensions {
		name := extractComponentName(e.GoMod)
		config.Components[core.TypeExtension][name] = PluginConfig(e)
	}

	return config, nil
}

// extractComponentName extracts the component name from a Go module path (last path segment, minus version).
//
// Usage:
//   - Used internally to normalize and extract plugin/component names from Go module paths.
//
// Best Practices:
//   - Always use this helper when converting legacy configs.
func extractComponentName(gomod string) string {
	parts := strings.Fields(gomod)
	if len(parts) == 0 {
		return ""
	}
	modPath := parts[0]
	base := path.Base(modPath)
	// Remove version suffix if present (e.g., v0.124.0)
	base = strings.TrimSuffix(base, filepath.Ext(base))
	return base
}
