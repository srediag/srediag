package factory

import (
	"fmt"
	"sync"

	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/types"
)

// Registry manages component factories
type Registry struct {
	mu        sync.RWMutex
	logger    *zap.Logger
	factories map[types.ComponentType]map[string]*Factory
}

// NewRegistry creates a new factory registry
func NewRegistry(logger *zap.Logger) *Registry {
	return &Registry{
		logger:    logger,
		factories: make(map[types.ComponentType]map[string]*Factory),
	}
}

// RegisterFactory registers a factory for a specific component type
func (r *Registry) RegisterFactory(componentType types.ComponentType, factory *Factory) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.factories[componentType]; !exists {
		r.factories[componentType] = make(map[string]*Factory)
	}

	id := factory.GetID()
	if _, exists := r.factories[componentType][id]; exists {
		return fmt.Errorf("factory already registered for type %s with id %s", componentType, id)
	}

	r.factories[componentType][id] = factory
	r.logger.Info("Registered factory",
		zap.String("component_type", componentType.String()),
		zap.String("factory_id", id))

	return nil
}

// GetFactory returns a factory for a specific component and id
func (r *Registry) GetFactory(componentType types.ComponentType, id string) (*Factory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if factories, exists := r.factories[componentType]; exists {
		if factory, exists := factories[id]; exists {
			return factory, nil
		}
	}

	return nil, fmt.Errorf("no factory registered for component type %s and id %s",
		componentType.String(), id)
}

// ListFactories returns all registered factories for a component type
func (r *Registry) ListFactories(componentType types.ComponentType) []*Factory {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var factories []*Factory
	if componentFactories, exists := r.factories[componentType]; exists {
		for _, factory := range componentFactories {
			factories = append(factories, factory)
		}
	}

	return factories
}

// UnregisterFactory removes a factory registration
func (r *Registry) UnregisterFactory(componentType types.ComponentType, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if factories, exists := r.factories[componentType]; exists {
		if _, exists := factories[id]; exists {
			delete(factories, id)
			r.logger.Info("Unregistered factory",
				zap.String("component_type", componentType.String()),
				zap.String("factory_id", id))
			return nil
		}
	}

	return fmt.Errorf("no factory registered for component type %s and id %s",
		componentType.String(), id)
}

// Clear removes all factory registrations
func (r *Registry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.factories = make(map[types.ComponentType]map[string]*Factory)
	r.logger.Info("Cleared all factory registrations")
}
