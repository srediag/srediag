// Package core provides foundational types and utilities for the SREDIAG system.
//
// This file defines the ComponentType type and constants, which enumerate the supported component types in SREDIAG.
// ComponentType is used for registration, lookup, and modularity throughout the system.
//
// Usage:
//   - Use ComponentType to identify, register, and look up component factories and instances.
//   - Used throughout the config, registry, and manager subsystems.
//
// Best Practices:
//   - Always use the defined constants for type safety and consistency.
//   - Extend with new types only when adding new component categories to SREDIAG.
//
// TODO:
//   - Consider supporting versioned or namespaced component types for plugins.
//
// Redundancy/Refactor:
//   - No redundancy; this is the canonical type system for SREDIAG components.
package core

// ComponentType represents a type of component in the SREDIAG system.
//
// Usage:
//   - Used as a key for registration, lookup, and config mapping.
type ComponentType string

const (
	// TypeCore represents a core component (internal system logic).
	TypeCore ComponentType = "core"
	// TypePlugin represents a plugin component (external extension).
	TypePlugin ComponentType = "plugin"
	// TypeConnector represents a connector component (data pipeline connector).
	TypeConnector ComponentType = "connector"
	// TypeExporter represents an exporter component (data sink).
	TypeExporter ComponentType = "exporter"
	// TypeExtension represents an extension component (optional feature).
	TypeExtension ComponentType = "extension"
	// TypeProcessor represents a processor component (data transformation).
	TypeProcessor ComponentType = "processor"
	// TypeReceiver represents a receiver component (data source).
	TypeReceiver ComponentType = "receiver"
)
