// Package core provides foundational types and utilities for the SREDIAG system.
//
// This file defines the BuildInfo struct, which holds versioning and build metadata for the SREDIAG binary.
// BuildInfo is typically set at build time via ldflags or similar mechanisms, and is used for diagnostics, support, and traceability.
package core

// BuildInfo holds version information about the SREDIAG build.
//
// Usage:
//
//	Use BuildInfo to report version, commit, and build date in CLI commands, logs, and diagnostics endpoints.
//	It is idiomatic to inject this struct at build time using Go ldflags for reproducible builds.
//
// Best Practices:
//   - Always set these fields via build tooling for production releases.
//   - Avoid hardcoding values except for local/dev builds.
//   - Use this struct for all version reporting to ensure consistency.
//
// TODO:
//   - Consider adding additional fields (e.g., Go version, build host) if needed for support.
//
// Redundancy/Refactor:
//   - No redundancy detected. This struct is the canonical source for build metadata in SREDIAG.
//
// This struct is used throughout the CLI and diagnostics to report the running version.
type BuildInfo struct {
	Version string // Semantic version (e.g., v1.2.3)
	Commit  string // Git commit hash or VCS identifier
	Date    string // Build timestamp (RFC3339 or similar)
}

// DefaultBuildInfo provides fallback build information when not set during compilation.
//
// Usage:
//
//	Used as a fallback if no build metadata is injected (e.g., during local development or CI without ldflags).
//	Should not be used in production builds—override via ldflags or build system.
//
// Best Practices:
//   - Always override these values for release builds.
//   - Use DefaultBuildInfo only as a last resort.
//
// TODO:
//   - None. This is a minimal, non-redundant fallback.
//
// Redundancy/Refactor:
//   - No redundancy detected. This is the only fallback for BuildInfo.
var DefaultBuildInfo = BuildInfo{
	Version: "dev",
	Commit:  "none",
	Date:    "unknown",
}
