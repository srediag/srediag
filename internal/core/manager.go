// Package core provides foundational types and utilities for the SREDIAG system.
//
// This file defines the ComponentManager, which handles registration and lookup of component factories for the SREDIAG system.
// ComponentManager enables modularity and dynamic extension of the system by managing component lifecycles and types.
//
// Usage:
//   - Use ComponentManager to register, look up, and manage component factories (receivers, processors, exporters, etc).
//   - Pass a ComponentManager to AppContext for dependency injection across the system.
//
// Best Practices:
//   - Always use the provided mutexes for thread safety.
//   - Register all factories at startup before using GetFactories.
//   - Use the logger for debug and error reporting on registration.
//
// TODO:
//   - Add support for deregistration and hot-reload of factories.
//   - Consider supporting versioned factories or metadata.
//   - Add lifecycle hooks for component initialization and shutdown.
//
// Redundancy/Refactor:
//   - No redundancy; this is the canonical component manager for SREDIAG.
package core

import (
	"fmt"
	"sync"

	"go.opentelemetry.io/collector/component"
)

// ComponentManager handles component registration and lifecycle.
//
// Usage:
//   - Use to register and retrieve factories for all supported component types.
//   - Thread-safe for concurrent registration and lookup.
//
// Fields:
//   - logger: Used for debug and error reporting.
//   - factories: Nested map of component type string to map of component.Type to Factory.
//   - mu: RWMutex for thread safety.
type ComponentManager struct {
	logger    *Logger
	factories map[string]map[component.Type]component.Factory // key: component type string
	mu        sync.RWMutex
}

// NewComponentManager creates a new component manager.
//
// Usage:
//   - Call at application startup to create a manager for all component factories.
//   - Pass a logger for debug and error reporting.
//
// Best Practices:
//   - Register all required factories immediately after creation.
//   - Use the returned manager for all component registration and lookup.
func NewComponentManager(logger *Logger) *ComponentManager {
	return &ComponentManager{
		logger: logger,
		factories: map[string]map[component.Type]component.Factory{
			"connector": make(map[component.Type]component.Factory),
			"exporter":  make(map[component.Type]component.Factory),
			"extension": make(map[component.Type]component.Factory),
			"processor": make(map[component.Type]component.Factory),
			"receiver":  make(map[component.Type]component.Factory),
		},
	}
}

// GetFactories returns all registered factories for a given component type.
//
// Usage:
//   - Use to retrieve all factories for a component type (e.g., "receiver").
//   - Returns a copy of the factories map to avoid race conditions.
//
// Best Practices:
//   - Always check for nil return if the component type is unknown.
//   - Use the returned map for read-only operations.
//
// TODO:
//   - Add support for filtering or querying factories by metadata.
func (m *ComponentManager) GetFactories(componentType string) map[component.Type]component.Factory {
	m.mu.RLock()
	defer m.mu.RUnlock()

	factories, ok := m.factories[componentType]
	if !ok {
		return nil
	}
	result := make(map[component.Type]component.Factory, len(factories))
	for k, v := range factories {
		result[k] = v
	}
	return result
}

// RegisterFactory registers a factory for a given component type ("receiver", "processor", etc.).
//
// Usage:
//   - Use to add a new factory to the manager for a specific component type.
//   - Returns an error if the type is unknown or already registered.
//
// Best Practices:
//   - Register all factories at startup before using GetFactories.
//   - Use the logger for debug output on successful registration.
//
// TODO:
//   - Add support for deregistration and hot-reload of factories.
//   - Consider supporting versioned factories or metadata.
func (m *ComponentManager) RegisterFactory(componentType string, factory component.Factory) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	factories, ok := m.factories[componentType]
	if !ok {
		return fmt.Errorf("unknown component type: %s", componentType)
	}
	if _, exists := factories[factory.Type()]; exists {
		return fmt.Errorf("%s %s already registered", componentType, factory.Type())
	}
	factories[factory.Type()] = factory
	m.logger.Debug("Registered factory", ZapString("component_type", componentType), ZapString("type", factory.Type().String()))
	return nil
}
