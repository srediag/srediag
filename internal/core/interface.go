package core

import (
	"context"
)

// IComponent defines the interface that all components must implement.
type IComponent interface {
	// Start starts the component.
	Start(ctx context.Context) error
	// Stop stops the component.
	Stop(ctx context.Context) error
}

// IFactory defines the interface for component factories.
type IFactory interface {
	// Type returns the type of component created by this factory.
	Type() ComponentType
	// CreateDefaultConfig creates the default configuration for the component.
	CreateDefaultConfig() interface{}
}

// IRegistry defines the interface for component registry.
type IRegistry interface {
	// RegisterFactory registers a new factory.
	RegisterFactory(factory IFactory) error
	// GetFactory returns a factory by type.
	GetFactory(typ ComponentType) (IFactory, bool)
	// GetFactories returns all registered factories.
	GetFactories() map[ComponentType]IFactory
}
