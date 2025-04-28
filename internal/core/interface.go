// Package core provides foundational types and interfaces for the SREDIAG system.
//
// This file defines the core interfaces for components, factories, and registries.
// These interfaces enable extensibility and modularity in the SREDIAG architecture.
package core

import (
	"context"
)

// IComponent defines the interface that all SREDIAG components must implement.
//
// Usage:
//   - Implement this interface for any component that needs explicit lifecycle management.
//   - Used by the component manager to start and stop components.
//
// Best Practices:
//   - Ensure Start and Stop are idempotent and handle repeated calls gracefully.
//   - Always use context for cancellation and deadlines.
//
// TODO:
//   - Consider adding a HealthCheck method for better observability.
//
// Redundancy/Refactor:
//   - No redundancy; this is the canonical component interface.
type IComponent interface {
	// Start starts the component.
	Start(ctx context.Context) error
	// Stop stops the component.
	Stop(ctx context.Context) error
}

// IFactory defines the interface for component factories.
//
// Usage:
//   - Implement this interface to provide new component types to the registry.
//   - Used for dynamic instantiation and config generation.
//
// Best Practices:
//   - Ensure Type returns a unique ComponentType for each factory.
//   - CreateDefaultConfig should return a fully valid config struct.
//
// TODO:
//   - Consider supporting versioning or metadata for factories.
//
// Redundancy/Refactor:
//   - No redundancy; this is the canonical factory interface.
type IFactory interface {
	// Type returns the type of component created by this factory.
	Type() ComponentType
	// CreateDefaultConfig creates the default configuration for the component.
	CreateDefaultConfig() interface{}
}

// IRegistry defines the interface for component registries.
//
// Usage:
//   - Use to register and look up factories for dynamic component management.
//   - Supports extensibility and modularity in the SREDIAG system.
//
// Best Practices:
//   - Always check for existence before registering a new factory.
//   - Use GetFactories to enumerate all available types.
//
// TODO:
//   - Consider supporting deregistration or hot-reload of factories.
//
// Redundancy/Refactor:
//   - No redundancy; this is the canonical registry interface.
type IRegistry interface {
	// RegisterFactory registers a new factory.
	RegisterFactory(factory IFactory) error
	// GetFactory returns a factory by type.
	GetFactory(typ ComponentType) (IFactory, bool)
	// GetFactories returns all registered factories.
	GetFactories() map[ComponentType]IFactory
}
