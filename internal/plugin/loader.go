package plugin

import (
	"fmt"
	"path/filepath"
	"plugin"
	"sync"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/connector"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/extension"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/factory"
)

// Loader handles loading and registration of OpenTelemetry components
type Loader struct {
	logger     *zap.Logger
	mu         sync.RWMutex
	receivers  map[component.Type]receiver.Factory
	processors map[component.Type]processor.Factory
	exporters  map[component.Type]exporter.Factory
	connectors map[component.Type]connector.Factory
	registry   *factory.Registry
}

// NewLoader creates a new component loader instance
func NewLoader(logger *zap.Logger, registry *factory.Registry) *Loader {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &Loader{
		logger:     logger,
		receivers:  make(map[component.Type]receiver.Factory),
		processors: make(map[component.Type]processor.Factory),
		exporters:  make(map[component.Type]exporter.Factory),
		connectors: make(map[component.Type]connector.Factory),
		registry:   registry,
	}
}

// RegisterReceiverFactory registers a receiver factory
func (l *Loader) RegisterReceiverFactory(factory receiver.Factory) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if factory == nil {
		return fmt.Errorf("factory cannot be nil")
	}

	componentType := factory.Type()
	if _, exists := l.receivers[componentType]; exists {
		return fmt.Errorf("receiver factory already registered for type %q", componentType)
	}

	l.receivers[componentType] = factory
	l.logger.Debug("Registered receiver factory",
		zap.Stringer("type", componentType))
	return nil
}

// RegisterProcessorFactory registers a processor factory
func (l *Loader) RegisterProcessorFactory(factory processor.Factory) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if factory == nil {
		return fmt.Errorf("factory cannot be nil")
	}

	componentType := factory.Type()
	if _, exists := l.processors[componentType]; exists {
		return fmt.Errorf("processor factory already registered for type %q", componentType)
	}

	l.processors[componentType] = factory
	l.logger.Debug("Registered processor factory",
		zap.Stringer("type", componentType))
	return nil
}

// RegisterExporterFactory registers an exporter factory
func (l *Loader) RegisterExporterFactory(factory exporter.Factory) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if factory == nil {
		return fmt.Errorf("factory cannot be nil")
	}

	componentType := factory.Type()
	if _, exists := l.exporters[componentType]; exists {
		return fmt.Errorf("exporter factory already registered for type %q", componentType)
	}

	l.exporters[componentType] = factory
	l.logger.Debug("Registered exporter factory",
		zap.Stringer("type", componentType))
	return nil
}

// RegisterConnectorFactory registers a connector factory
func (l *Loader) RegisterConnectorFactory(factory connector.Factory) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if factory == nil {
		return fmt.Errorf("factory cannot be nil")
	}

	componentType := factory.Type()
	if _, exists := l.connectors[componentType]; exists {
		return fmt.Errorf("connector factory already registered for type %q", componentType)
	}

	l.connectors[componentType] = factory
	l.logger.Debug("Registered connector factory",
		zap.Stringer("type", componentType))
	return nil
}

// GetReceiverFactory returns a receiver factory by type
func (l *Loader) GetReceiverFactory(componentType component.Type) (receiver.Factory, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	factory, exists := l.receivers[componentType]
	return factory, exists
}

// GetProcessorFactory returns a processor factory by type
func (l *Loader) GetProcessorFactory(componentType component.Type) (processor.Factory, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	factory, exists := l.processors[componentType]
	return factory, exists
}

// GetExporterFactory returns an exporter factory by type
func (l *Loader) GetExporterFactory(componentType component.Type) (exporter.Factory, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	factory, exists := l.exporters[componentType]
	return factory, exists
}

// GetConnectorFactory returns a connector factory by type
func (l *Loader) GetConnectorFactory(componentType component.Type) (connector.Factory, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	factory, exists := l.connectors[componentType]
	return factory, exists
}

// LoadNativeComponents loads native OpenTelemetry components
func (l *Loader) LoadNativeComponents(factories map[component.Type]interface{}) error {
	for typ, f := range factories {
		switch factory := f.(type) {
		case receiver.Factory:
			if err := l.registry.RegisterReceiver(factory, "native"); err != nil {
				return fmt.Errorf("failed to register native receiver %q: %w", typ, err)
			}
		case processor.Factory:
			if err := l.registry.RegisterProcessor(factory, "native"); err != nil {
				return fmt.Errorf("failed to register native processor %q: %w", typ, err)
			}
		case exporter.Factory:
			if err := l.registry.RegisterExporter(factory, "native"); err != nil {
				return fmt.Errorf("failed to register native exporter %q: %w", typ, err)
			}
		case extension.Factory:
			if err := l.registry.RegisterExtension(factory, "native"); err != nil {
				return fmt.Errorf("failed to register native extension %q: %w", typ, err)
			}
		case connector.Factory:
			if err := l.registry.RegisterConnector(factory, "native"); err != nil {
				return fmt.Errorf("failed to register native connector %q: %w", typ, err)
			}
		default:
			return fmt.Errorf("unknown factory type for %q: %T", typ, factory)
		}
		l.logger.Info("Registered native component", zap.String("type", typ.String()))
	}
	return nil
}

// LoadCustomPlugins loads custom plugins from the specified directory
func (l *Loader) LoadCustomPlugins(pluginDir string) error {
	// Get list of .so files in plugin directory
	soFiles, err := filepath.Glob(filepath.Join(pluginDir, "*.so"))
	if err != nil {
		return fmt.Errorf("failed to list plugin files: %w", err)
	}

	for _, soFile := range soFiles {
		if err := l.loadCustomPlugin(soFile); err != nil {
			l.logger.Error("Failed to load plugin",
				zap.String("file", soFile),
				zap.Error(err))
			continue
		}
	}
	return nil
}

// loadCustomPlugin loads a single custom plugin
func (l *Loader) loadCustomPlugin(pluginPath string) error {
	// Open plugin
	plug, err := plugin.Open(pluginPath)
	if err != nil {
		return fmt.Errorf("failed to open plugin %q: %w", pluginPath, err)
	}

	// Look up factory symbol
	factorySymbol, err := plug.Lookup("Factory")
	if err != nil {
		return fmt.Errorf("plugin %q does not export 'Factory' symbol: %w", pluginPath, err)
	}

	// Register factory based on its type
	switch factory := factorySymbol.(type) {
	case receiver.Factory:
		if err := l.registry.RegisterReceiver(factory, pluginPath); err != nil {
			return fmt.Errorf("failed to register receiver from plugin %q: %w", pluginPath, err)
		}
	case processor.Factory:
		if err := l.registry.RegisterProcessor(factory, pluginPath); err != nil {
			return fmt.Errorf("failed to register processor from plugin %q: %w", pluginPath, err)
		}
	case exporter.Factory:
		if err := l.registry.RegisterExporter(factory, pluginPath); err != nil {
			return fmt.Errorf("failed to register exporter from plugin %q: %w", pluginPath, err)
		}
	case extension.Factory:
		if err := l.registry.RegisterExtension(factory, pluginPath); err != nil {
			return fmt.Errorf("failed to register extension from plugin %q: %w", pluginPath, err)
		}
	case connector.Factory:
		if err := l.registry.RegisterConnector(factory, pluginPath); err != nil {
			return fmt.Errorf("failed to register connector from plugin %q: %w", pluginPath, err)
		}
	default:
		return fmt.Errorf("unknown factory type from plugin %q: %T", pluginPath, factory)
	}

	l.logger.Info("Loaded plugin", zap.String("file", pluginPath))
	return nil
}

// GetFactories returns all registered factories
func (l *Loader) GetFactories() (map[component.Type]receiver.Factory,
	map[component.Type]processor.Factory,
	map[component.Type]exporter.Factory,
	map[component.Type]extension.Factory,
	map[component.Type]connector.Factory) {
	return l.registry.GetFactories()
}

// GetModuleInfo returns module information for all registered components
func (l *Loader) GetModuleInfo() (map[component.Type]string,
	map[component.Type]string,
	map[component.Type]string,
	map[component.Type]string,
	map[component.Type]string) {
	return l.registry.GetModuleInfo()
}
