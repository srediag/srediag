// Package app provides the main application functionality
package app

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/extension"
	"go.opentelemetry.io/collector/otelcol"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/config"
	"github.com/srediag/srediag/internal/diagnostic"
	"github.com/srediag/srediag/internal/types"
)

// SREDiag represents the main application instance
type SREDiag struct {
	logger *zap.Logger
	config *config.ConfigRoot

	pluginManager     types.IPluginManager
	resourceMonitor   types.IResourceMonitor
	configManager     types.IConfigManager
	telemetryBridge   types.ITelemetryBridge
	collector         *otelcol.Collector
	diagnosticManager types.IDiagnosticManager

	// Diagnostic components
	systemDiag     types.IDiagnostic
	kubernetesDiag types.IDiagnostic
	cloudDiag      types.IDiagnostic

	mu      sync.RWMutex
	health  bool
	running bool
}

// Ensure SREDiag implements IRunner
var _ types.IRunner = (*SREDiag)(nil)

// NewSREDiag creates a new instance of SREDiag
func NewSREDiag(logger *zap.Logger, cfg *config.ConfigRoot) (*SREDiag, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}

	s := &SREDiag{
		logger: logger,
		config: cfg,
		health: true,
	}

	var err error

	// Initialize plugin manager
	pluginManager := diagnostic.NewPluginManager(logger.Named("plugin-manager"))
	s.pluginManager = pluginManager

	// Create resource for telemetry
	res, err := cfg.CreateResource()
	if err != nil {
		return nil, fmt.Errorf("failed to create telemetry resource: %w", err)
	}

	// Initialize telemetry bridge
	telemetryBridge := diagnostic.NewTelemetryBridge(
		logger.Named("telemetry-bridge"),
		res,
	)
	s.telemetryBridge = telemetryBridge

	// Initialize resource monitor with meter from telemetry bridge
	meter := telemetryBridge.GetMeterProvider().Meter("srediag")
	resourceMonitor := diagnostic.NewResourceMonitor(
		logger.Named("resource-monitor"),
		meter,
	)
	s.resourceMonitor = resourceMonitor

	// Initialize config manager
	configManager, err := diagnostic.NewConfigManager(logger.Named("config-manager"))
	if err != nil {
		return nil, fmt.Errorf("failed to create config manager: %w", err)
	}
	s.configManager = configManager

	// Initialize diagnostic manager
	diagnosticManager := diagnostic.NewDiagnosticManager(logger.Named("diagnostic-manager"))
	s.diagnosticManager = diagnosticManager

	// Initialize collector if enabled
	if cfg.Collector.Enabled {
		if err := s.initializeCollector(); err != nil {
			return nil, fmt.Errorf("failed to initialize collector: %w", err)
		}
	}

	// Initialize diagnostic components
	if err := s.initializeDiagnostics(); err != nil {
		return nil, fmt.Errorf("failed to initialize diagnostics: %w", err)
	}

	return s, nil
}

// initializeCollector initializes the OpenTelemetry Collector
func (s *SREDiag) initializeCollector() error {
	// Create factories
	factories, err := s.createFactories()
	if err != nil {
		return fmt.Errorf("failed to create factories: %w", err)
	}

	settings := otelcol.CollectorSettings{
		BuildInfo: component.BuildInfo{
			Command:     s.config.Service.Name,
			Description: "SREDIAG Diagnostic Collector",
			Version:     s.config.Service.Version,
		},
		Factories: func() (otelcol.Factories, error) {
			return factories, nil
		},
	}

	collector, err := otelcol.NewCollector(settings)
	if err != nil {
		return fmt.Errorf("failed to create collector: %w", err)
	}

	s.collector = collector
	return nil
}

// createFactories creates and returns all required OpenTelemetry factories
func (s *SREDiag) createFactories() (otelcol.Factories, error) {
	factories := otelcol.Factories{
		Receivers:  make(map[component.Type]receiver.Factory),
		Processors: make(map[component.Type]processor.Factory),
		Exporters:  make(map[component.Type]exporter.Factory),
		Extensions: make(map[component.Type]extension.Factory),
	}

	return factories, nil
}

// initializeDiagnostics initializes all diagnostic components
func (s *SREDiag) initializeDiagnostics() error {
	meter := s.telemetryBridge.GetMeterProvider().Meter("srediag.diagnostics")

	// Initialize system diagnostics if enabled
	if s.config.Diagnostic.System.Enabled {
		s.systemDiag = diagnostic.NewSystem(
			s.logger.Named("system-diagnostic"),
			meter,
			&s.config.Diagnostic.System,
		)
		if err := s.systemDiag.Configure(&s.config.Diagnostic.System); err != nil {
			return fmt.Errorf("failed to configure system diagnostics: %w", err)
		}
	}

	// Initialize Kubernetes diagnostics if enabled
	if s.config.Diagnostic.Kubernetes.Enabled {
		s.kubernetesDiag = diagnostic.NewKubernetes(
			s.logger.Named("kubernetes-diagnostic"),
			meter,
			&s.config.Diagnostic.Kubernetes,
		)
		if err := s.kubernetesDiag.Configure(&s.config.Diagnostic.Kubernetes); err != nil {
			return fmt.Errorf("failed to configure kubernetes diagnostics: %w", err)
		}
	}

	// Initialize cloud diagnostics if enabled
	if s.config.Diagnostic.Cloud.Enabled {
		s.cloudDiag = diagnostic.NewCloud(
			s.logger.Named("cloud-diagnostic"),
			meter,
			&s.config.Diagnostic.Cloud,
		)
		if err := s.cloudDiag.Configure(&s.config.Diagnostic.Cloud); err != nil {
			return fmt.Errorf("failed to configure cloud diagnostics: %w", err)
		}
	}

	return nil
}

// Start initializes and starts all components
func (s *SREDiag) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return fmt.Errorf("already running")
	}

	// Start telemetry bridge first to ensure metrics and tracing are available
	if err := s.telemetryBridge.Start(ctx); err != nil {
		return fmt.Errorf("failed to start telemetry bridge: %w", err)
	}

	// Start resource monitor
	if err := s.resourceMonitor.Start(ctx); err != nil {
		return fmt.Errorf("failed to start resource monitor: %w", err)
	}

	// Start config manager
	if err := s.configManager.Start(ctx); err != nil {
		return fmt.Errorf("failed to start config manager: %w", err)
	}

	// Start plugin manager
	if err := s.pluginManager.StartAll(ctx); err != nil {
		return fmt.Errorf("failed to start plugin manager: %w", err)
	}

	// Start diagnostic manager
	if err := s.diagnosticManager.Start(ctx); err != nil {
		return fmt.Errorf("failed to start diagnostic manager: %w", err)
	}

	// Start collector if enabled
	if s.collector != nil {
		if err := s.collector.Run(ctx); err != nil {
			return fmt.Errorf("failed to start collector: %w", err)
		}
	}

	// Start diagnostic components
	if err := s.startDiagnostics(ctx); err != nil {
		return fmt.Errorf("failed to start diagnostics: %w", err)
	}

	s.running = true
	s.logger.Info("SREDIAG started successfully")

	return nil
}

// startDiagnostics starts all diagnostic components
func (s *SREDiag) startDiagnostics(ctx context.Context) error {
	if s.systemDiag != nil {
		if err := s.systemDiag.Start(ctx); err != nil {
			return fmt.Errorf("failed to start system diagnostics: %w", err)
		}
	}

	if s.kubernetesDiag != nil {
		if err := s.kubernetesDiag.Start(ctx); err != nil {
			return fmt.Errorf("failed to start kubernetes diagnostics: %w", err)
		}
	}

	if s.cloudDiag != nil {
		if err := s.cloudDiag.Start(ctx); err != nil {
			return fmt.Errorf("failed to start cloud diagnostics: %w", err)
		}
	}

	return nil
}

// Stop gracefully stops all components
func (s *SREDiag) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return nil
	}

	// Stop diagnostic components first
	if err := s.stopDiagnostics(ctx); err != nil {
		s.logger.Error("error stopping diagnostics", zap.Error(err))
	}

	// Stop collector if enabled
	if s.collector != nil {
		s.collector.Shutdown()
	}

	// Stop components in reverse order
	if err := s.diagnosticManager.Stop(ctx); err != nil {
		s.logger.Error("error stopping diagnostic manager", zap.Error(err))
	}

	if err := s.pluginManager.StopAll(ctx); err != nil {
		s.logger.Error("error stopping plugin manager", zap.Error(err))
	}

	if err := s.configManager.Stop(ctx); err != nil {
		s.logger.Error("error stopping config manager", zap.Error(err))
	}

	if err := s.resourceMonitor.Stop(ctx); err != nil {
		s.logger.Error("error stopping resource monitor", zap.Error(err))
	}

	if err := s.telemetryBridge.Stop(ctx); err != nil {
		s.logger.Error("error stopping telemetry bridge", zap.Error(err))
	}

	s.running = false
	s.logger.Info("SREDIAG stopped successfully")

	return nil
}

// stopDiagnostics stops all diagnostic components
func (s *SREDiag) stopDiagnostics(ctx context.Context) error {
	if s.systemDiag != nil {
		if err := s.systemDiag.Stop(ctx); err != nil {
			return fmt.Errorf("failed to stop system diagnostics: %w", err)
		}
	}

	if s.kubernetesDiag != nil {
		if err := s.kubernetesDiag.Stop(ctx); err != nil {
			return fmt.Errorf("failed to stop kubernetes diagnostics: %w", err)
		}
	}

	if s.cloudDiag != nil {
		if err := s.cloudDiag.Stop(ctx); err != nil {
			return fmt.Errorf("failed to stop cloud diagnostics: %w", err)
		}
	}

	return nil
}

// IsHealthy returns the health status of the application
func (s *SREDiag) IsHealthy() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	healthy := s.health &&
		s.pluginManager.IsHealthy() &&
		s.resourceMonitor.IsHealthy() &&
		s.configManager.IsHealthy() &&
		s.telemetryBridge.IsHealthy()

	// Check diagnostic components health
	if s.systemDiag != nil && !s.systemDiag.IsHealthy() {
		healthy = false
	}
	if s.kubernetesDiag != nil && !s.kubernetesDiag.IsHealthy() {
		healthy = false
	}
	if s.cloudDiag != nil && !s.cloudDiag.IsHealthy() {
		healthy = false
	}

	return healthy
}

// GetLogger returns the configured logger
func (s *SREDiag) GetLogger() *zap.Logger {
	return s.logger
}

// GetCollector returns the OpenTelemetry collector instance if enabled
func (s *SREDiag) GetCollector() *otelcol.Collector {
	return s.collector
}

// GetPluginManager returns the plugin manager instance
func (s *SREDiag) GetPluginManager() types.IPluginManager {
	return s.pluginManager
}

// GetResourceMonitor returns the resource monitor instance
func (s *SREDiag) GetResourceMonitor() types.IResourceMonitor {
	return s.resourceMonitor
}

// GetConfigManager returns the config manager instance
func (s *SREDiag) GetConfigManager() types.IConfigManager {
	return s.configManager
}

// GetTelemetryBridge returns the telemetry bridge instance
func (s *SREDiag) GetTelemetryBridge() types.ITelemetryBridge {
	return s.telemetryBridge
}

// GetConfig returns the configuration
func (s *SREDiag) GetConfig() types.IConfig {
	return s.config
}

// GetDiagnosticManager returns the diagnostic manager
func (s *SREDiag) GetDiagnosticManager() types.IDiagnosticManager {
	return s.diagnosticManager
}

// Configure implements IComponent
func (s *SREDiag) Configure(cfg interface{}) error {
	if cfg == nil {
		return fmt.Errorf("configuration cannot be nil")
	}

	config, ok := cfg.(*config.ConfigRoot)
	if !ok {
		return fmt.Errorf("invalid configuration type")
	}

	s.config = config
	return nil
}

// GetName implements IComponent
func (s *SREDiag) GetName() string {
	return s.config.Service.Name
}

// GetVersion implements IComponent
func (s *SREDiag) GetVersion() string {
	return s.config.Service.Version
}

// GetType implements types.IRunner
func (s *SREDiag) GetType() types.ComponentType {
	return types.ComponentTypeService
}
