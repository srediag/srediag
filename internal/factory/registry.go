package factory

import (
	"fmt"
	"sync"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/connector"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/extension"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/receiver"
)

// Registry manages both native OpenTelemetry components and custom components
type Registry struct {
	mu sync.RWMutex

	receivers  map[component.Type]receiver.Factory
	processors map[component.Type]processor.Factory
	exporters  map[component.Type]exporter.Factory
	extensions map[component.Type]extension.Factory
	connectors map[component.Type]connector.Factory

	// Module information
	receiverModules  map[component.Type]string
	processorModules map[component.Type]string
	exporterModules  map[component.Type]string
	extensionModules map[component.Type]string
	connectorModules map[component.Type]string
}

// NewRegistry creates a new Registry
func NewRegistry() *Registry {
	return &Registry{
		receivers:        make(map[component.Type]receiver.Factory),
		processors:       make(map[component.Type]processor.Factory),
		exporters:        make(map[component.Type]exporter.Factory),
		extensions:       make(map[component.Type]extension.Factory),
		connectors:       make(map[component.Type]connector.Factory),
		receiverModules:  make(map[component.Type]string),
		processorModules: make(map[component.Type]string),
		exporterModules:  make(map[component.Type]string),
		extensionModules: make(map[component.Type]string),
		connectorModules: make(map[component.Type]string),
	}
}

// RegisterReceiver registers a receiver factory
func (r *Registry) RegisterReceiver(f receiver.Factory, module string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.receivers[f.Type()]; exists {
		return fmt.Errorf("receiver factory %q already registered", f.Type())
	}
	r.receivers[f.Type()] = f
	r.receiverModules[f.Type()] = module
	return nil
}

// RegisterProcessor registers a processor factory
func (r *Registry) RegisterProcessor(f processor.Factory, module string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.processors[f.Type()]; exists {
		return fmt.Errorf("processor factory %q already registered", f.Type())
	}
	r.processors[f.Type()] = f
	r.processorModules[f.Type()] = module
	return nil
}

// RegisterExporter registers an exporter factory
func (r *Registry) RegisterExporter(f exporter.Factory, module string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.exporters[f.Type()]; exists {
		return fmt.Errorf("exporter factory %q already registered", f.Type())
	}
	r.exporters[f.Type()] = f
	r.exporterModules[f.Type()] = module
	return nil
}

// RegisterExtension registers an extension factory
func (r *Registry) RegisterExtension(f extension.Factory, module string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.extensions[f.Type()]; exists {
		return fmt.Errorf("extension factory %q already registered", f.Type())
	}
	r.extensions[f.Type()] = f
	r.extensionModules[f.Type()] = module
	return nil
}

// RegisterConnector registers a connector factory
func (r *Registry) RegisterConnector(f connector.Factory, module string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.connectors[f.Type()]; exists {
		return fmt.Errorf("connector factory %q already registered", f.Type())
	}
	r.connectors[f.Type()] = f
	r.connectorModules[f.Type()] = module
	return nil
}

// GetFactories returns all registered factories
func (r *Registry) GetFactories() (map[component.Type]receiver.Factory,
	map[component.Type]processor.Factory,
	map[component.Type]exporter.Factory,
	map[component.Type]extension.Factory,
	map[component.Type]connector.Factory) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	receivers := make(map[component.Type]receiver.Factory, len(r.receivers))
	for k, v := range r.receivers {
		receivers[k] = v
	}

	processors := make(map[component.Type]processor.Factory, len(r.processors))
	for k, v := range r.processors {
		processors[k] = v
	}

	exporters := make(map[component.Type]exporter.Factory, len(r.exporters))
	for k, v := range r.exporters {
		exporters[k] = v
	}

	extensions := make(map[component.Type]extension.Factory, len(r.extensions))
	for k, v := range r.extensions {
		extensions[k] = v
	}

	connectors := make(map[component.Type]connector.Factory, len(r.connectors))
	for k, v := range r.connectors {
		connectors[k] = v
	}

	return receivers, processors, exporters, extensions, connectors
}

// GetModuleInfo returns module information for all registered components
func (r *Registry) GetModuleInfo() (map[component.Type]string,
	map[component.Type]string,
	map[component.Type]string,
	map[component.Type]string,
	map[component.Type]string) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	receiverModules := make(map[component.Type]string, len(r.receiverModules))
	for k, v := range r.receiverModules {
		receiverModules[k] = v
	}

	processorModules := make(map[component.Type]string, len(r.processorModules))
	for k, v := range r.processorModules {
		processorModules[k] = v
	}

	exporterModules := make(map[component.Type]string, len(r.exporterModules))
	for k, v := range r.exporterModules {
		exporterModules[k] = v
	}

	extensionModules := make(map[component.Type]string, len(r.extensionModules))
	for k, v := range r.extensionModules {
		extensionModules[k] = v
	}

	connectorModules := make(map[component.Type]string, len(r.connectorModules))
	for k, v := range r.connectorModules {
		connectorModules[k] = v
	}

	return receiverModules, processorModules, exporterModules, extensionModules, connectorModules
}
