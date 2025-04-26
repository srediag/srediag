package factory

import (
	"fmt"
	"sync"

	"github.com/srediag/srediag/internal/types"
)

// registry implements the types.Registry interface
type registry struct {
	mu        sync.RWMutex
	factories map[types.Type]types.Factory
}

// NewRegistry creates a new component registry
func NewRegistry() types.Registry {
	return &registry{
		factories: make(map[types.Type]types.Factory),
	}
}

// RegisterFactory registers a new factory
func (r *registry) RegisterFactory(factory types.Factory) error {
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

// GetFactory returns a factory by type
func (r *registry) GetFactory(typ types.Type) (types.Factory, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	factory, exists := r.factories[typ]
	return factory, exists
}

// GetFactories returns all registered factories
func (r *registry) GetFactories() map[types.Type]types.Factory {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Create a copy to avoid concurrent access issues
	factories := make(map[types.Type]types.Factory, len(r.factories))
	for k, v := range r.factories {
		factories[k] = v
	}

	return factories
}
