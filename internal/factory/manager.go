package factory

import (
	"fmt"
	"sync"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/extension"
	"go.opentelemetry.io/collector/otelcol"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/types"
)

// Manager handles component factories and provides integration with OpenTelemetry Collector
type Manager struct {
	mu         sync.RWMutex
	logger     *zap.Logger
	receivers  map[component.Type]receiver.Factory
	processors map[component.Type]processor.Factory
	exporters  map[component.Type]exporter.Factory
	extensions map[component.Type]extension.Factory
	buildInfo  component.BuildInfo
}

// NewManager creates a new factory manager
func NewManager(logger *zap.Logger, buildInfo component.BuildInfo) *Manager {
	return &Manager{
		logger:     logger,
		receivers:  make(map[component.Type]receiver.Factory),
		processors: make(map[component.Type]processor.Factory),
		exporters:  make(map[component.Type]exporter.Factory),
		extensions: make(map[component.Type]extension.Factory),
		buildInfo:  buildInfo,
	}
}

// RegisterReceiverFactory registers a new receiver factory
func (m *Manager) RegisterReceiverFactory(factory receiver.Factory) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.receivers[factory.Type()]; exists {
		return fmt.Errorf("receiver factory already registered for type %s", factory.Type())
	}

	m.receivers[factory.Type()] = factory
	m.logger.Info("Registered receiver factory", zap.String("type", factory.Type().String()))
	return nil
}

// RegisterProcessorFactory registers a new processor factory
func (m *Manager) RegisterProcessorFactory(factory processor.Factory) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.processors[factory.Type()]; exists {
		return fmt.Errorf("processor factory already registered for type %s", factory.Type())
	}

	m.processors[factory.Type()] = factory
	m.logger.Info("Registered processor factory", zap.String("type", factory.Type().String()))
	return nil
}

// RegisterExporterFactory registers a new exporter factory
func (m *Manager) RegisterExporterFactory(factory exporter.Factory) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.exporters[factory.Type()]; exists {
		return fmt.Errorf("exporter factory already registered for type %s", factory.Type())
	}

	m.exporters[factory.Type()] = factory
	m.logger.Info("Registered exporter factory", zap.String("type", factory.Type().String()))
	return nil
}

// RegisterExtensionFactory registers a new extension factory
func (m *Manager) RegisterExtensionFactory(factory extension.Factory) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.extensions[factory.Type()]; exists {
		return fmt.Errorf("extension factory already registered for type %s", factory.Type())
	}

	m.extensions[factory.Type()] = factory
	m.logger.Info("Registered extension factory", zap.String("type", factory.Type().String()))
	return nil
}

// GetReceiverFactory returns a receiver factory by type
func (m *Manager) GetReceiverFactory(typ component.Type) (receiver.Factory, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	factory, exists := m.receivers[typ]
	if !exists {
		return nil, fmt.Errorf("receiver factory not found for type %s", typ)
	}

	return factory, nil
}

// GetProcessorFactory returns a processor factory by type
func (m *Manager) GetProcessorFactory(typ component.Type) (processor.Factory, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	factory, exists := m.processors[typ]
	if !exists {
		return nil, fmt.Errorf("processor factory not found for type %s", typ)
	}

	return factory, nil
}

// GetExporterFactory returns an exporter factory by type
func (m *Manager) GetExporterFactory(typ component.Type) (exporter.Factory, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	factory, exists := m.exporters[typ]
	if !exists {
		return nil, fmt.Errorf("exporter factory not found for type %s", typ)
	}

	return factory, nil
}

// GetExtensionFactory returns an extension factory by type
func (m *Manager) GetExtensionFactory(typ component.Type) (extension.Factory, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	factory, exists := m.extensions[typ]
	if !exists {
		return nil, fmt.Errorf("extension factory not found for type %s", typ)
	}

	return factory, nil
}

// GetFactories returns all registered factories
func (m *Manager) GetFactories() (otelcol.Factories, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return otelcol.Factories{
		Receivers:  m.receivers,
		Processors: m.processors,
		Exporters:  m.exporters,
		Extensions: m.extensions,
	}, nil
}

// UnregisterFactory unregisters a factory by type and category
func (m *Manager) UnregisterFactory(category types.ComponentType, typ component.Type) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch category {
	case types.ComponentTypeReceiver:
		delete(m.receivers, typ)
	case types.ComponentTypeProcessor:
		delete(m.processors, typ)
	case types.ComponentTypeExporter:
		delete(m.exporters, typ)
	case types.ComponentTypeExtension:
		delete(m.extensions, typ)
	default:
		return fmt.Errorf("unknown component category %v", category)
	}

	m.logger.Info("Unregistered factory",
		zap.String("category", category.String()),
		zap.String("type", typ.String()))
	return nil
}

// Clear removes all registered factories
func (m *Manager) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.receivers = make(map[component.Type]receiver.Factory)
	m.processors = make(map[component.Type]processor.Factory)
	m.exporters = make(map[component.Type]exporter.Factory)
	m.extensions = make(map[component.Type]extension.Factory)

	m.logger.Info("Cleared all factories")
}

// GetBuildInfo returns the build information
func (m *Manager) GetBuildInfo() component.BuildInfo {
	return m.buildInfo
}
