// Package plugin provides plugin types and utilities for SREDIAG
package plugin

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/discovery"
	"github.com/srediag/srediag/internal/factory"
	"github.com/srediag/srediag/internal/types"
)

// Manager handles plugin lifecycle and management
type Manager struct {
	mu             sync.RWMutex
	logger         *zap.Logger
	registry       *Registry
	discovery      *discovery.Manager
	components     map[component.ID]component.Component
	plugins        map[string]types.IPlugin
	factories      map[component.Type]component.Factory
	errors         map[component.ID][]error
	host           component.Host
	buildInfo      component.BuildInfo
	tracerProvider trace.TracerProvider
	meterProvider  metric.MeterProvider
	running        bool
}

// NewManager creates a new plugin manager
func NewManager(logger *zap.Logger, registry *Registry, discovery *discovery.Manager, host component.Host, buildInfo component.BuildInfo, tracerProvider trace.TracerProvider, meterProvider metric.MeterProvider) *Manager {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &Manager{
		logger:         logger,
		registry:       registry,
		discovery:      discovery,
		components:     make(map[component.ID]component.Component),
		plugins:        make(map[string]types.IPlugin),
		factories:      make(map[component.Type]component.Factory),
		errors:         make(map[component.ID][]error),
		host:           host,
		buildInfo:      buildInfo,
		tracerProvider: tracerProvider,
		meterProvider:  meterProvider,
	}
}

// RegisterFactory registers a component factory
func (m *Manager) RegisterFactory(f component.Factory) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	typ := f.Type()
	if _, exists := m.factories[typ]; exists {
		return fmt.Errorf("factory already registered for type %q", typ)
	}

	m.factories[typ] = f
	return nil
}

// CreateComponent creates a new component instance
func (m *Manager) CreateComponent(typ component.Type, id component.ID, cfg component.Config) (component.Component, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	factory, exists := m.factories[typ]
	if !exists {
		return nil, fmt.Errorf("no factory registered for type %q", typ)
	}

	baseSettings := component.TelemetrySettings{
		Logger:         m.logger.With(zap.String("component_id", id.String())),
		TracerProvider: m.tracerProvider,
		MeterProvider:  m.meterProvider,
	}

	var comp component.Component
	var err error

	switch f := factory.(type) {
	case receiver.Factory:
		receiverSettings := receiver.Settings{
			ID:                id,
			TelemetrySettings: baseSettings,
			BuildInfo:         m.buildInfo,
		}
		comp, err = f.CreateMetrics(context.Background(), receiverSettings, cfg, nil)
	case processor.Factory:
		processorSettings := processor.Settings{
			ID:                id,
			TelemetrySettings: baseSettings,
			BuildInfo:         m.buildInfo,
		}
		comp, err = f.CreateMetrics(context.Background(), processorSettings, cfg, nil)
	case exporter.Factory:
		exporterSettings := exporter.Settings{
			ID:                id,
			TelemetrySettings: baseSettings,
			BuildInfo:         m.buildInfo,
		}
		comp, err = f.CreateMetrics(context.Background(), exporterSettings, cfg)
	default:
		return nil, fmt.Errorf("unsupported factory type %T", factory)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create component: %w", err)
	}

	m.components[id] = comp
	return comp, nil
}

// RegisterPlugin registers a plugin with the manager
func (m *Manager) RegisterPlugin(plugin types.IPlugin) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := plugin.Validate(); err != nil {
		return fmt.Errorf("invalid plugin: %w", err)
	}

	id := plugin.GetID()
	m.plugins[id] = plugin
	return m.registry.RegisterPlugin(plugin)
}

// UnregisterPlugin removes a plugin from the manager
func (m *Manager) UnregisterPlugin(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.plugins, id)
	return m.registry.UnregisterPlugin(id)
}

// GetPlugin returns a plugin by ID
func (m *Manager) GetPlugin(id string) (types.IPlugin, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	plugin, exists := m.plugins[id]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", id)
	}

	return plugin, nil
}

// ListPlugins returns all registered plugins
func (m *Manager) ListPlugins() []types.IPlugin {
	m.mu.RLock()
	defer m.mu.RUnlock()

	plugins := make([]types.IPlugin, 0, len(m.plugins))
	for _, plugin := range m.plugins {
		plugins = append(plugins, plugin)
	}
	return plugins
}

// ListPluginsByCategory returns plugins filtered by category
func (m *Manager) ListPluginsByCategory(category types.PluginCategory) []types.IPlugin {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var plugins []types.IPlugin
	for _, plugin := range m.plugins {
		if plugin.GetCategory() == category {
			plugins = append(plugins, plugin)
		}
	}
	return plugins
}

// ListPluginsByCapability returns plugins filtered by capability
func (m *Manager) ListPluginsByCapability(capability types.PluginCapability) []types.IPlugin {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var plugins []types.IPlugin
	for _, plugin := range m.plugins {
		if plugin.GetCapabilities().HasCapability(capability) {
			plugins = append(plugins, plugin)
		}
	}
	return plugins
}

// Start starts all components and plugins
func (m *Manager) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.running {
		return fmt.Errorf("plugin manager is already running")
	}

	// Start OpenTelemetry components
	for id, comp := range m.components {
		if err := comp.Start(ctx, m.host); err != nil {
			m.logger.Error("Failed to start component",
				zap.String("id", id.String()),
				zap.Error(err))
			m.errors[id] = append(m.errors[id], err)
		}
	}

	// Start SREDIAG plugins
	for id, plugin := range m.plugins {
		if err := plugin.Start(ctx); err != nil {
			m.logger.Error("Failed to start plugin",
				zap.String("id", id),
				zap.Error(err))
		}
	}

	// Start plugin discovery watcher in a goroutine
	go m.watchDiscovery(ctx)

	m.running = true
	m.logger.Info("Started plugin manager")
	return nil
}

// Shutdown stops all components and plugins
func (m *Manager) Shutdown(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.running {
		return nil
	}

	// Stop OpenTelemetry components
	for id, comp := range m.components {
		if err := comp.Shutdown(ctx); err != nil {
			m.logger.Error("Failed to shutdown component",
				zap.String("id", id.String()),
				zap.Error(err))
			m.errors[id] = append(m.errors[id], err)
		}
	}

	// Stop SREDIAG plugins
	for id, plugin := range m.plugins {
		if err := plugin.Stop(ctx); err != nil {
			m.logger.Error("Failed to stop plugin",
				zap.String("id", id),
				zap.Error(err))
		}
	}

	m.running = false
	m.logger.Info("Stopped plugin manager")
	return nil
}

// GetHost returns the OpenTelemetry Collector host
func (m *Manager) GetHost() component.Host {
	return m.host
}

// GetBuildInfo returns the build information
func (m *Manager) GetBuildInfo() component.BuildInfo {
	return m.buildInfo
}

// GetFactories returns all plugin factories
func (m *Manager) GetFactories() map[string]*factory.Factory {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.registry.GetFactories()
}

// watchDiscovery watches for plugin discovery events
func (m *Manager) watchDiscovery(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case service := <-m.discovery.Watch():
			m.handleDiscoveredPlugin(service)
		}
	}
}

// handleDiscoveredPlugin handles a discovered plugin
func (m *Manager) handleDiscoveredPlugin(service discovery.ServiceInfo) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if plugin is already registered
	if _, err := m.registry.GetPlugin(service.ID); err == nil {
		return
	}

	m.logger.Info("Discovered new plugin",
		zap.String("id", service.ID),
		zap.String("name", service.Name),
		zap.String("version", service.Version))

	// TODO: Implement plugin loading and registration logic here
	// This could include:
	// 1. Loading the plugin from the discovered service info
	// 2. Validating the plugin
	// 3. Registering it with the manager
	// 4. Starting the plugin if the manager is running
}

// GetName implements types.IComponent
func (m *Manager) GetName() string {
	return "plugin-manager"
}

// GetVersion implements types.IComponent
func (m *Manager) GetVersion() string {
	return "1.0.0"
}

// GetType implements types.IComponent
func (m *Manager) GetType() types.ComponentType {
	return types.ComponentTypeCore
}

// Configure implements types.IComponent
func (m *Manager) Configure(cfg interface{}) error {
	return nil
}

// GetComponent returns a component by ID
func (m *Manager) GetComponent(id component.ID) (component.Component, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	comp, exists := m.components[id]
	return comp, exists
}

// ListComponents returns all registered components
func (m *Manager) ListComponents() []component.Component {
	m.mu.RLock()
	defer m.mu.RUnlock()

	components := make([]component.Component, 0, len(m.components))
	for _, comp := range m.components {
		components = append(components, comp)
	}
	return components
}
