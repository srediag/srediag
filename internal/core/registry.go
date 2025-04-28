// Package core provides foundational types and utilities for the SREDIAG system.
//
// This file defines the registry type, which implements the IRegistry interface for component factories.
// The registry enables extensibility and modularity by supporting dynamic registration and lookup of component factories.
//
// Usage:
//   - Use registry to manage IFactory implementations for all component types.
//   - Use NewRegistry to create a new registry for dependency injection or plugin systems.
//
// Best Practices:
//   - Always use the provided mutexes for thread safety.
//   - Register all factories at startup before using GetFactory or GetFactories.
//   - Use the returned map from GetFactories for read-only operations only.
//
// TODO:
//   - Add support for deregistration and hot-reload of component factories in the registry.
//   - Consider supporting versioned factories or attaching metadata to factories.
//   - Add lifecycle hooks for factory initialization and shutdown.
//
// Redundancy/Refactor:
//   - No redundancy; this is the canonical registry for SREDIAG component factories.
//
// TODO(C-03 Phase 1): Implement component registry with lazy load (see TODO.md C-03, ETA 2025-06-07)
// TODO(C-04 Phase 1): Implement graceful shutdown logic (flush + RocksDB close) and tie into system signals (see TODO.md C-04, ETA 2025-06-14)
package core

import (
	"fmt"
	"sync"
)

// registry implements the IRegistry interface for component factories.
//
// Usage:
//   - Use to register and retrieve IFactory implementations for all supported component types.
//   - Thread-safe for concurrent registration and lookup.
//
// Fields:
//   - mu: RWMutex for thread safety.
//   - factories: Map of ComponentType to IFactory.
type registry struct {
	mu        sync.RWMutex
	factories map[ComponentType]IFactory
}

// NewRegistry creates a new component registry.
//
// Usage:
//   - Call at application or plugin system startup to create a registry for all component factories.
//
// Best Practices:
//   - Register all required factories immediately after creation.
//   - Use the returned registry for all component registration and lookup.
func NewRegistry() IRegistry {
	return &registry{
		factories: make(map[ComponentType]IFactory),
	}
}

// RegisterFactory registers a new factory.
//
// Usage:
//   - Use to add a new IFactory to the registry for a specific ComponentType.
//   - Returns an error if the factory is nil or already registered for the type.
//
// Best Practices:
//   - Register all factories at startup before using GetFactory or GetFactories.
//
// TODO:
//   - Add support for deregistration and hot-reload of component factories in the registry.
//   - Consider supporting versioned factories or attaching metadata to factories.
//   - Add lifecycle hooks for factory initialization and shutdown.
func (r *registry) RegisterFactory(factory IFactory) error {
	if factory == nil {
		return fmt.Errorf("factory cannot be nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	typ := factory.Type()
	if _, exists := r.factories[typ]; exists {
		return fmt.Errorf("factory already registered for type %q", typ)
	}

	r.factories[typ] = factory
	return nil
}

// GetFactory returns a factory by type.
//
// Usage:
//   - Use to retrieve a registered IFactory for a given ComponentType.
//   - Returns the factory and a boolean indicating existence.
//
// Best Practices:
//   - Always check the boolean return value before using the factory.
func (r *registry) GetFactory(typ ComponentType) (IFactory, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	factory, exists := r.factories[typ]
	return factory, exists
}

// GetFactories returns all registered factories.
//
// Usage:
//   - Use to enumerate all registered IFactory implementations.
//   - Returns a copy of the factories map to avoid race conditions.
//
// Best Practices:
//   - Use the returned map for read-only operations only.
func (r *registry) GetFactories() map[ComponentType]IFactory {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Create a copy to avoid concurrent access issues
	factories := make(map[ComponentType]IFactory, len(r.factories))
	for k, v := range r.factories {
		factories[k] = v
	}

	return factories
}
