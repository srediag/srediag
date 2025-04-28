// Package build provides the build orchestration layer for SREDIAG.
//
// This file defines the interfaces for build operations, including the BuildManager interface and related abstractions.
// These interfaces enable mocking, testing, and alternative implementations for build orchestration.
//
// Usage:
//   - Use these interfaces to decouple build orchestration logic from CLI and implementation details.
//   - Implement BuildManager for custom build workflows or testing.
//
// Best Practices:
//   - Always use interfaces for dependency injection and testing.
//   - Document all interface methods with expected side effects and error handling.
//
// TODO:
//   - Add context.Context to all interface methods for cancellation and timeouts.
//   - Consider splitting large interfaces into smaller, focused ones.
//
// Redundancy/Refactor:
//   - If only one implementation exists, consider if the interface is necessary.
//
// TODO(P-02 Phase 1): Implement Manifest v1 JSON schema (see TODO.md P-02, ETA 2025-05-31)
// TODO: Enforce manifest generation for plugins (SHA-256, cosign signature reference) (see architecture/build.md §3)
// TODO: Validate ABI compatibility using Go symbol tables (see architecture/build.md §3)
// TODO: Fail build if manifest or ABI check fails (see architecture/build.md §3)
package build

import "github.com/srediag/srediag/internal/core"

// Builder defines the interface for building OpenTelemetry plugins
// Only orchestration methods are exposed; config details are handled by BuildManager.
type IBuilder interface {
	// BuildAll builds all plugins defined in the configuration
	BuildAll() error

	// BuildPlugin builds a single plugin by name and type
	BuildPlugin(name string, compType core.ComponentType) error
}

// MakeBuilderInterface extends Builder with make integration
// Note: YAML update is handled by UpdateBuilderYAMLVersions, not this interface.
// All config loading should use LoadBuildConfig for schema compliance and validation.
type MakeBuilderInterface interface {
	IBuilder

	// InstallPlugins installs built plugins to the system
	InstallPlugins() error
}

// BuildManagerInterface defines the interface for orchestrating build operations.
//
// Usage:
//   - Use this interface to abstract build orchestration logic for CLI, tests, or alternative implementations.
//   - Implement for custom build workflows or mocking in tests.
//
// Best Practices:
//   - All methods should return detailed errors for diagnostics.
//   - Document side effects (filesystem, network, etc).
//
// TODO:
//   - Add context.Context to all methods for cancellation.
type BuildManagerInterface interface {
	// BuildAll builds the agent and all plugins.
	// Returns error if any build step fails.
	BuildAll() error

	// BuildPlugin builds a single plugin by type and name.
	// Returns error if the build fails or plugin is not found.
	BuildPlugin(pluginType, pluginName string) error

	// Generate scaffolds plugin code for a given type and name.
	// Returns error if generation fails.
	Generate(pluginType, pluginName string) error

	// InstallPlugins copies pre-built plugins to the execution directory.
	// Returns error if installation fails.
	InstallPlugins() error
}
