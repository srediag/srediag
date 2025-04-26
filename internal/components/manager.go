package components

import (
	"fmt"
	"sync"

	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"
)

// Manager handles component registration and lifecycle
type Manager struct {
	logger     *zap.Logger
	receivers  map[component.Type]component.Factory
	processors map[component.Type]component.Factory
	exporters  map[component.Type]component.Factory
	extensions map[component.Type]component.Factory
	mu         sync.RWMutex
}

// NewManager creates a new component manager
func NewManager(logger *zap.Logger) *Manager {
	return &Manager{
		logger:     logger,
		receivers:  make(map[component.Type]component.Factory),
		processors: make(map[component.Type]component.Factory),
		exporters:  make(map[component.Type]component.Factory),
		extensions: make(map[component.Type]component.Factory),
	}
}

// RegisterReceiver registers a receiver factory
func (m *Manager) RegisterReceiver(factory component.Factory) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.receivers[factory.Type()]; exists {
		return fmt.Errorf("receiver %s already registered", factory.Type())
	}

	m.receivers[factory.Type()] = factory
	m.logger.Debug("Registered receiver", zap.String("type", factory.Type().String()))
	return nil
}

// RegisterProcessor registers a processor factory
func (m *Manager) RegisterProcessor(factory component.Factory) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.processors[factory.Type()]; exists {
		return fmt.Errorf("processor %s already registered", factory.Type())
	}

	m.processors[factory.Type()] = factory
	m.logger.Debug("Registered processor", zap.String("type", factory.Type().String()))
	return nil
}

// RegisterExporter registers an exporter factory
func (m *Manager) RegisterExporter(factory component.Factory) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.exporters[factory.Type()]; exists {
		return fmt.Errorf("exporter %s already registered", factory.Type())
	}

	m.exporters[factory.Type()] = factory
	m.logger.Debug("Registered exporter", zap.String("type", factory.Type().String()))
	return nil
}

// RegisterExtension registers an extension factory
func (m *Manager) RegisterExtension(factory component.Factory) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.extensions[factory.Type()]; exists {
		return fmt.Errorf("extension %s already registered", factory.Type())
	}

	m.extensions[factory.Type()] = factory
	m.logger.Debug("Registered extension", zap.String("type", factory.Type().String()))
	return nil
}

// GetReceivers returns all registered receiver factories
func (m *Manager) GetReceivers() map[component.Type]component.Factory {
	m.mu.RLock()
	defer m.mu.RUnlock()

	receivers := make(map[component.Type]component.Factory, len(m.receivers))
	for k, v := range m.receivers {
		receivers[k] = v
	}
	return receivers
}

// GetProcessors returns all registered processor factories
func (m *Manager) GetProcessors() map[component.Type]component.Factory {
	m.mu.RLock()
	defer m.mu.RUnlock()

	processors := make(map[component.Type]component.Factory, len(m.processors))
	for k, v := range m.processors {
		processors[k] = v
	}
	return processors
}

// GetExporters returns all registered exporter factories
func (m *Manager) GetExporters() map[component.Type]component.Factory {
	m.mu.RLock()
	defer m.mu.RUnlock()

	exporters := make(map[component.Type]component.Factory, len(m.exporters))
	for k, v := range m.exporters {
		exporters[k] = v
	}
	return exporters
}

// GetExtensions returns all registered extension factories
func (m *Manager) GetExtensions() map[component.Type]component.Factory {
	m.mu.RLock()
	defer m.mu.RUnlock()

	extensions := make(map[component.Type]component.Factory, len(m.extensions))
	for k, v := range m.extensions {
		extensions[k] = v
	}
	return extensions
}
