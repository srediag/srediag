package telemetry

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/types"
)

// Manager manages telemetry for the application
type Manager struct {
	mu         sync.RWMutex
	logger     *zap.Logger
	config     *types.TelemetryConfig
	tracer     trace.Tracer
	meter      metric.Meter
	host       component.Host
	buildInfo  component.BuildInfo
	components map[string]component.Component
}

// NewManager creates a new telemetry manager
func NewManager(logger *zap.Logger, config *types.TelemetryConfig, host component.Host, buildInfo component.BuildInfo) *Manager {
	return &Manager{
		logger:     logger,
		config:     config,
		host:       host,
		buildInfo:  buildInfo,
		components: make(map[string]component.Component),
	}
}

// Start starts the telemetry manager
func (m *Manager) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.config.IsEnabled() {
		m.logger.Info("Telemetry is disabled")
		return nil
	}

	// Initialize tracer
	tracerProvider := otel.GetTracerProvider()
	m.tracer = tracerProvider.Tracer(
		"srediag",
		trace.WithInstrumentationVersion(m.buildInfo.Version),
	)

	// Initialize meter
	meterProvider := otel.GetMeterProvider()
	m.meter = meterProvider.Meter(
		"srediag",
		metric.WithInstrumentationVersion(m.buildInfo.Version),
	)

	m.logger.Info("Started telemetry manager",
		zap.String("version", m.buildInfo.Version))
	return nil
}

// Stop stops the telemetry manager
func (m *Manager) Stop(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for name, component := range m.components {
		if err := component.Shutdown(ctx); err != nil {
			m.logger.Error("Failed to stop component",
				zap.String("name", name),
				zap.Error(err))
		}
	}

	m.logger.Info("Stopped telemetry manager")
	return nil
}

// RegisterComponent registers a telemetry component
func (m *Manager) RegisterComponent(name string, component component.Component) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.components[name]; exists {
		return fmt.Errorf("component %s already registered", name)
	}

	m.components[name] = component
	m.logger.Info("Registered telemetry component",
		zap.String("name", name))

	return nil
}

// UnregisterComponent removes a telemetry component
func (m *Manager) UnregisterComponent(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.components[name]; !exists {
		return fmt.Errorf("component %s not registered", name)
	}

	delete(m.components, name)
	m.logger.Info("Unregistered telemetry component",
		zap.String("name", name))

	return nil
}

// GetComponent returns a telemetry component by name
func (m *Manager) GetComponent(name string) (component.Component, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	component, exists := m.components[name]
	if !exists {
		return nil, fmt.Errorf("component %s not found", name)
	}

	return component, nil
}

// ListComponents returns all registered telemetry components
func (m *Manager) ListComponents() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	components := make([]string, 0, len(m.components))
	for name := range m.components {
		components = append(components, name)
	}

	return components
}

// GetTracer returns the tracer
func (m *Manager) GetTracer() trace.Tracer {
	return m.tracer
}

// GetMeter returns the meter
func (m *Manager) GetMeter() metric.Meter {
	return m.meter
}

// GetHost returns the OpenTelemetry Collector host
func (m *Manager) GetHost() component.Host {
	return m.host
}

// GetBuildInfo returns the build information
func (m *Manager) GetBuildInfo() component.BuildInfo {
	return m.buildInfo
}

// UpdateConfig updates the telemetry configuration
func (m *Manager) UpdateConfig(config *types.TelemetryConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.config = config
	m.logger.Info("Updated telemetry configuration")
	return nil
}

// GetConfig returns the current telemetry configuration
func (m *Manager) GetConfig() *types.TelemetryConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.config
}
