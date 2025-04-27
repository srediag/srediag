package core

import (
	"fmt"
	"sync"

	"go.opentelemetry.io/collector/component"
)

// ComponentManager handles component registration and lifecycle.
type ComponentManager struct {
	logger    *Logger
	factories map[string]map[component.Type]component.Factory // key: component type string
	mu        sync.RWMutex
}

// NewComponentManager creates a new component manager.
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
