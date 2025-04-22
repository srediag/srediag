package app

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/otelcol"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/config"
	"github.com/srediag/srediag/internal/core"
	"github.com/srediag/srediag/internal/core/diagnostic"
)

// SREDiag represents the main application instance
type SREDiag struct {
	logger *zap.Logger
	config *config.SREDiagConfig

	pluginManager   core.PluginManager
	resourceMonitor core.ResourceMonitor
	configManager   core.ConfigManager
	telemetryBridge core.TelemetryBridge
	collector       *otelcol.Collector

	// Diagnostic components
	systemDiag     diagnostic.Component
	kubernetesDiag diagnostic.Component
	cloudDiag      diagnostic.Component
	securityDiag   diagnostic.Component

	mu      sync.RWMutex
	health  bool
	running bool
}

// Ensure SREDiag implements ISREDiagRunner
var _ core.ISREDiagRunner = (*SREDiag)(nil)

// NewSREDiag creates a new instance of SREDiag
func NewSREDiag(logger *zap.Logger, cfg *config.SREDiagConfig) (*SREDiag, error) {
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
	pluginManager := core.NewPluginManager(logger.Named("plugin-manager"))
	s.pluginManager = pluginManager

	// Create resource for telemetry
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			"", // Empty schema version as it's not needed
			semconv.ServiceName(cfg.Service.Name),
			semconv.ServiceVersion(cfg.Service.Version),
			semconv.DeploymentEnvironment(cfg.Service.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create telemetry resource: %w", err)
	}

	// Initialize telemetry bridge
	telemetryBridge := core.NewTelemetryBridge(
		logger.Named("telemetry-bridge"),
		res,
	)
	s.telemetryBridge = telemetryBridge

	// Initialize resource monitor with meter from telemetry bridge
	meter := telemetryBridge.GetMeterProvider().Meter("srediag")
	resourceMonitor := core.NewResourceMonitor(
		logger.Named("resource-monitor"),
		meter,
	)
	s.resourceMonitor = resourceMonitor

	// Initialize config manager
	configManager, err := core.NewConfigManager(logger.Named("config-manager"))
	if err != nil {
		return nil, fmt.Errorf("failed to create config manager: %w", err)
	}
	s.configManager = configManager

	// Initialize collector if enabled
	if cfg.Collector.IsEnabled() {
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
	settings := otelcol.CollectorSettings{
		BuildInfo: component.BuildInfo{
			Command:     s.config.Service.Name,
			Description: "SREDIAG Diagnostic Collector",
			Version:     s.config.Service.Version,
		},
		Factories: func() (otelcol.Factories, error) {
			return otelcol.Factories{}, nil
		},
	}

	collector, err := otelcol.NewCollector(settings)
	if err != nil {
		return fmt.Errorf("failed to create collector: %w", err)
	}

	s.collector = collector
	return nil
}

// initializeDiagnostics initializes all diagnostic components
func (s *SREDiag) initializeDiagnostics() error {
	meter := s.telemetryBridge.GetMeterProvider().Meter("srediag.diagnostics")

	// Initialize system diagnostics if enabled
	if s.config.Diagnostic.System.Enabled {
		s.systemDiag = diagnostic.NewSystem(
			s.logger.Named("system-diagnostic"),
			meter,
			s.config.Diagnostic.System,
		)
		if err := s.systemDiag.Configure(s.config.Diagnostic.System); err != nil {
			return fmt.Errorf("failed to configure system diagnostics: %w", err)
		}
	}

	// Initialize Kubernetes diagnostics if enabled
	if s.config.Diagnostic.Kubernetes.Enabled {
		s.kubernetesDiag = diagnostic.NewKubernetes(
			s.logger.Named("kubernetes-diagnostic"),
			meter,
			s.config.Diagnostic.Kubernetes,
		)
		if err := s.kubernetesDiag.Configure(s.config.Diagnostic.Kubernetes); err != nil {
			return fmt.Errorf("failed to configure kubernetes diagnostics: %w", err)
		}
	}

	// Initialize cloud diagnostics if enabled
	if s.config.Diagnostic.Cloud.Enabled {
		s.cloudDiag = diagnostic.NewCloud(
			s.logger.Named("cloud-diagnostic"),
			meter,
			s.config.Diagnostic.Cloud,
		)
		if err := s.cloudDiag.Configure(s.config.Diagnostic.Cloud); err != nil {
			return fmt.Errorf("failed to configure cloud diagnostics: %w", err)
		}
	}

	// Initialize security diagnostics if enabled
	if s.config.Diagnostic.Security.Enabled {
		s.securityDiag = diagnostic.NewSecurity(
			s.logger.Named("security-diagnostic"),
			meter,
			s.config.Diagnostic.Security,
		)
		if err := s.securityDiag.Configure(s.config.Diagnostic.Security); err != nil {
			return fmt.Errorf("failed to configure security diagnostics: %w", err)
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

	// Start config manager
	if err := s.configManager.Start(ctx); err != nil {
		return fmt.Errorf("failed to start config manager: %w", err)
	}

	// Start resource monitor
	if err := s.resourceMonitor.Start(ctx); err != nil {
		return fmt.Errorf("failed to start resource monitor: %w", err)
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

	// Start plugin manager last as it may depend on other components
	if err := s.pluginManager.Start(ctx); err != nil {
		return fmt.Errorf("failed to start plugin manager: %w", err)
	}

	s.running = true
	s.logger.Info("srediag started successfully",
		zap.String("version", s.config.Service.Version),
		zap.String("environment", s.config.Service.Environment))
	return nil
}

// startDiagnostics starts all enabled diagnostic components
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

	if s.securityDiag != nil {
		if err := s.securityDiag.Start(ctx); err != nil {
			return fmt.Errorf("failed to start security diagnostics: %w", err)
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

	// Stop components in reverse order
	if err := s.pluginManager.Stop(ctx); err != nil {
		s.logger.Error("failed to stop plugin manager", zap.Error(err))
	}

	if err := s.stopDiagnostics(ctx); err != nil {
		s.logger.Error("failed to stop diagnostics", zap.Error(err))
	}

	if s.collector != nil {
		s.collector.Shutdown()
	}

	if err := s.resourceMonitor.Stop(ctx); err != nil {
		s.logger.Error("failed to stop resource monitor", zap.Error(err))
	}

	if err := s.configManager.Stop(ctx); err != nil {
		s.logger.Error("failed to stop config manager", zap.Error(err))
	}

	if err := s.telemetryBridge.Stop(ctx); err != nil {
		s.logger.Error("failed to stop telemetry bridge", zap.Error(err))
	}

	s.running = false
	s.logger.Info("srediag stopped successfully")
	return nil
}

// stopDiagnostics stops all enabled diagnostic components
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

	if s.securityDiag != nil {
		if err := s.securityDiag.Stop(ctx); err != nil {
			return fmt.Errorf("failed to stop security diagnostics: %w", err)
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
	if s.securityDiag != nil && !s.securityDiag.IsHealthy() {
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
func (s *SREDiag) GetPluginManager() core.PluginManager {
	return s.pluginManager
}

// GetResourceMonitor returns the resource monitor instance
func (s *SREDiag) GetResourceMonitor() core.ResourceMonitor {
	return s.resourceMonitor
}

// GetConfigManager returns the config manager instance
func (s *SREDiag) GetConfigManager() core.ConfigManager {
	return s.configManager
}

// GetTelemetryBridge returns the telemetry bridge instance
func (s *SREDiag) GetTelemetryBridge() core.TelemetryBridge {
	return s.telemetryBridge
}
