// Package diagnostic provides diagnostic operations for SREDIAG
package diagnostic

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
	tracenoop "go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/types"
)

// Manager manages diagnostic components
type Manager struct {
	mu             sync.RWMutex
	logger         *zap.Logger
	components     map[string]Component
	settings       types.TelemetrySettings
	resource       *resource.Resource
	healthy        bool
	running        bool
	meterProvider  metric.MeterProvider
	tracerProvider trace.TracerProvider
}

// NewManager creates a new diagnostic manager
func NewManager(settings types.TelemetrySettings) *Manager {
	if settings.Logger == nil {
		settings.Logger = zap.NewNop()
	}

	return &Manager{
		logger:         settings.Logger,
		components:     make(map[string]Component),
		settings:       settings,
		resource:       settings.Resource,
		healthy:        true,
		running:        false,
		meterProvider:  settings.MeterProvider,
		tracerProvider: settings.TracerProvider,
	}
}

// NewPluginManager creates a new plugin manager instance
func NewPluginManager(logger *zap.Logger) types.IPluginManager {
	// TODO: Implement plugin manager
	return nil
}

// NewResourceMonitor creates a new resource monitor instance
func NewResourceMonitor(logger *zap.Logger, meter metric.Meter) types.IResourceMonitor {
	// TODO: Implement resource monitor
	return nil
}

// NewConfigManager creates a new config manager instance
func NewConfigManager(logger *zap.Logger) (types.IConfigManager, error) {
	// TODO: Implement config manager
	return nil, nil
}

// NewDiagnosticManager creates a new diagnostic manager instance
func NewDiagnosticManager(logger *zap.Logger) types.IDiagnosticManager {
	// TODO: Implement diagnostic manager
	return nil
}

// MeterProvider implements IMeterProvider
type MeterProvider struct {
	provider metric.MeterProvider
}

// Meter implements IMeterProvider
func (p *MeterProvider) Meter(name string, opts ...metric.MeterOption) metric.Meter {
	return p.provider.Meter(name, opts...)
}

// MeterProvider implements IMeterProvider
func (p *MeterProvider) MeterProvider() metric.MeterProvider {
	return p.provider
}

// Shutdown implements IMeterProvider
func (p *MeterProvider) Shutdown(ctx context.Context) error {
	if shutdownable, ok := p.provider.(interface{ Shutdown(context.Context) error }); ok {
		return shutdownable.Shutdown(ctx)
	}
	return nil
}

// TracerProvider implements ITracerProvider
type TracerProvider struct {
	provider trace.TracerProvider
}

// Tracer implements ITracerProvider
func (p *TracerProvider) Tracer(name string, opts ...trace.TracerOption) trace.Tracer {
	return p.provider.Tracer(name, opts...)
}

// Shutdown implements ITracerProvider
func (p *TracerProvider) Shutdown(ctx context.Context) error {
	if shutdownable, ok := p.provider.(interface{ Shutdown(context.Context) error }); ok {
		return shutdownable.Shutdown(ctx)
	}
	return nil
}

// TelemetryBridge implements ITelemetryBridge
type TelemetryBridge struct {
	logger         *zap.Logger
	resource       map[string]string
	healthy        bool
	meterProvider  metric.MeterProvider
	tracerProvider trace.TracerProvider
	meter          metric.Meter
	tracer         trace.Tracer
	running        bool
	mu             sync.RWMutex
}

// NewTelemetryBridge creates a new telemetry bridge instance
func NewTelemetryBridge(logger *zap.Logger, res map[string]string) types.ITelemetryBridge {
	if logger == nil {
		logger = zap.NewNop()
	}

	if res == nil {
		res = make(map[string]string)
	}

	return &TelemetryBridge{
		logger:         logger,
		resource:       res,
		healthy:        true,
		meterProvider:  noop.NewMeterProvider(),
		tracerProvider: tracenoop.NewTracerProvider(),
	}
}

// Start implements ITelemetryBridge
func (b *TelemetryBridge) Start(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.running = true
	return nil
}

// Stop implements ITelemetryBridge
func (b *TelemetryBridge) Stop(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.running = false
	return nil
}

// IsHealthy implements ITelemetryBridge
func (b *TelemetryBridge) IsHealthy() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.healthy
}

// GetName implements ITelemetryBridge
func (b *TelemetryBridge) GetName() string {
	return "telemetry-bridge"
}

// GetVersion implements ITelemetryBridge
func (b *TelemetryBridge) GetVersion() string {
	return "1.0.0"
}

// GetType implements ITelemetryBridge
func (b *TelemetryBridge) GetType() types.ComponentType {
	return types.ComponentTypeService
}

// Configure configures the bridge with settings
func (b *TelemetryBridge) Configure(settings types.ComponentSettings) error {
	if settings == nil {
		return fmt.Errorf("settings cannot be nil")
	}

	telemetrySettings := settings.GetTelemetrySettings()
	if telemetrySettings == nil {
		return fmt.Errorf("telemetry settings cannot be nil")
	}

	b.logger = telemetrySettings.Logger
	b.tracer = telemetrySettings.Tracer
	b.meter = telemetrySettings.Meter

	return nil
}

// SetResource implements ITelemetryBridge
func (b *TelemetryBridge) SetResource(res map[string]string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if res == nil {
		res = make(map[string]string)
	}
	b.resource = res
	return nil
}

// GetResource implements ITelemetryBridge
func (b *TelemetryBridge) GetResource() map[string]string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.resource
}

// GetMeterProvider implements ITelemetryBridge
func (b *TelemetryBridge) GetMeterProvider() types.IMeterProvider {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return &MeterProvider{provider: b.meterProvider}
}

// GetTracerProvider implements ITelemetryBridge
func (b *TelemetryBridge) GetTracerProvider() types.ITracerProvider {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return &TracerProvider{provider: b.tracerProvider}
}

// GetStatus implements IComponent
func (b *TelemetryBridge) GetStatus() types.ComponentStatus {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if !b.running {
		return types.ComponentStatusStopped
	}
	if !b.healthy {
		return types.ComponentStatusError
	}
	return types.ComponentStatusRunning
}

// Shutdown implements ITelemetryBridge
func (b *TelemetryBridge) Shutdown(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.running {
		return nil
	}

	// Stop the bridge first
	if err := b.Stop(ctx); err != nil {
		return fmt.Errorf("failed to stop bridge: %w", err)
	}

	// Shutdown meter provider
	if mp, ok := b.meterProvider.(interface{ Shutdown(context.Context) error }); ok {
		if err := mp.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown meter provider: %w", err)
		}
	}

	// Shutdown tracer provider
	if tp, ok := b.tracerProvider.(interface{ Shutdown(context.Context) error }); ok {
		if err := tp.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown tracer provider: %w", err)
		}
	}

	b.healthy = false
	return nil
}

// Start starts the manager and all registered components
func (m *Manager) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.running {
		return nil
	}

	// Start all components
	for name, component := range m.components {
		if err := component.Start(ctx); err != nil {
			m.logger.Error("Failed to start component",
				zap.String("name", name),
				zap.Error(err))
			m.healthy = false
			return fmt.Errorf("failed to start component %s: %w", name, err)
		}
	}

	m.running = true
	m.logger.Info("Started diagnostic manager")
	return nil
}

// Stop stops the manager and all registered components
func (m *Manager) Stop(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.running {
		return nil
	}

	// Stop all components
	var errs []error
	for name, component := range m.components {
		if err := component.Stop(ctx); err != nil {
			m.logger.Error("Failed to stop component",
				zap.String("name", name),
				zap.Error(err))
			errs = append(errs, fmt.Errorf("failed to stop component %s: %w", name, err))
		}
	}

	// Shutdown meter provider if it implements Shutdownable
	if mp, ok := m.meterProvider.(interface{ Shutdown(context.Context) error }); ok {
		if err := mp.Shutdown(ctx); err != nil {
			m.logger.Error("Failed to shutdown meter provider", zap.Error(err))
			errs = append(errs, fmt.Errorf("failed to shutdown meter provider: %w", err))
		}
	}

	// Shutdown tracer provider if it implements Shutdownable
	if tp, ok := m.tracerProvider.(interface{ Shutdown(context.Context) error }); ok {
		if err := tp.Shutdown(ctx); err != nil {
			m.logger.Error("Failed to shutdown tracer provider", zap.Error(err))
			errs = append(errs, fmt.Errorf("failed to shutdown tracer provider: %w", err))
		}
	}

	m.running = false
	m.logger.Info("Stopped diagnostic manager")

	if len(errs) > 0 {
		return fmt.Errorf("errors stopping components: %v", errs)
	}
	return nil
}

// Register registers a diagnostic component
func (m *Manager) Register(component Component) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	name := component.GetName()
	if _, exists := m.components[name]; exists {
		return fmt.Errorf("component %s already registered", name)
	}

	// Create component settings
	compSettings := types.ComponentSettings{
		"telemetry": &types.TelemetrySettings{
			Logger:                 m.logger.With(zap.String("component", name)),
			Tracer:                 m.settings.Tracer,
			Meter:                  m.settings.Meter,
			TracerProvider:         m.settings.TracerProvider,
			MeterProvider:          m.settings.MeterProvider,
			Resource:               m.resource,
			InstrumentationVersion: "1.0.0",
			InstrumentationScope: instrumentation.Scope{
				Name:      name,
				Version:   component.GetVersion(),
				SchemaURL: "",
			},
		},
	}

	if err := component.Configure(compSettings); err != nil {
		return fmt.Errorf("failed to configure component %s: %w", name, err)
	}

	m.components[name] = component
	m.logger.Info("Registered diagnostic component", zap.String("name", name))
	return nil
}

// Unregister unregisters a diagnostic component
func (m *Manager) Unregister(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.components[name]; !exists {
		return fmt.Errorf("component %s not registered", name)
	}

	delete(m.components, name)
	m.logger.Info("Unregistered diagnostic component", zap.String("name", name))
	return nil
}

// GetComponent returns a registered component by name
func (m *Manager) GetComponent(name string) (Component, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	component, exists := m.components[name]
	return component, exists
}

// GetComponents returns all registered components
func (m *Manager) GetComponents() map[string]Component {
	m.mu.RLock()
	defer m.mu.RUnlock()
	components := make(map[string]Component, len(m.components))
	for name, component := range m.components {
		components[name] = component
	}
	return components
}

// GetMeterProvider returns the meter provider
func (m *Manager) GetMeterProvider() metric.MeterProvider {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.meterProvider
}

// GetTracerProvider returns the tracer provider
func (m *Manager) GetTracerProvider() trace.TracerProvider {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.tracerProvider
}

// IsHealthy returns the health status
func (m *Manager) IsHealthy() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.healthy
}

// SetHealth sets the health status
func (m *Manager) SetHealth(healthy bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.healthy = healthy
}

// GetLogger returns the logger instance
func (m *Manager) GetLogger() *zap.Logger {
	return m.logger
}

// GetName returns the component name
func (m *Manager) GetName() string {
	return "diagnostic-manager"
}

// GetVersion returns the component version
func (m *Manager) GetVersion() string {
	return "1.0.0"
}

// GetType returns the component type
func (m *Manager) GetType() types.ComponentType {
	return types.ComponentTypeService
}

// Configure configures the manager with settings
func (m *Manager) Configure(settings types.ComponentSettings) error {
	if settings == nil {
		return fmt.Errorf("settings cannot be nil")
	}

	telemetrySettings := settings.GetTelemetrySettings()
	if telemetrySettings == nil {
		return fmt.Errorf("telemetry settings cannot be nil")
	}

	m.logger = telemetrySettings.Logger
	m.settings = *telemetrySettings
	m.resource = telemetrySettings.Resource
	m.meterProvider = telemetrySettings.MeterProvider
	m.tracerProvider = telemetrySettings.TracerProvider

	return nil
}
