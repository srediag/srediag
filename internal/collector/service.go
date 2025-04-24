package collector

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/config"
	"github.com/srediag/srediag/internal/factory"
	"github.com/srediag/srediag/internal/telemetry"
	"github.com/srediag/srediag/internal/types"
)

// Service represents the collector service
type Service struct {
	logger            *zap.Logger
	factories         *factory.Factory
	config            *config.Manager
	telemetry         *telemetry.Manager
	settings          *types.ServiceSettings
	buildInfo         component.BuildInfo
	telemetrySettings component.TelemetrySettings
}

// Options holds service configuration options
type Options struct {
	ConfigPath string
	Settings   *types.ServiceSettings
	BuildInfo  component.BuildInfo
}

// New creates a new collector service
func New(opts Options) (*Service, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	// Create factory with build info
	factoryOpts := []factory.FactoryOption{
		factory.WithLogger(logger),
		factory.WithID("collector"),
		factory.WithComponentType(types.ComponentTypeService),
	}
	factoryMgr := factory.NewFactory(factoryOpts...)

	// Create default service settings if not provided
	settings := opts.Settings
	if settings == nil {
		settings = &types.ServiceSettings{
			Name:    "collector",
			Version: "1.0.0",
			Type:    types.ComponentTypeService,
		}
	}

	configMgr := config.NewManager(logger, opts.ConfigPath)

	telemetryConfig := &types.TelemetryConfig{
		Enabled: true,
		Metrics: types.MetricsConfig{
			Enabled: true,
		},
	}

	telemetryMgr := telemetry.NewManager(logger, telemetryConfig, nil, opts.BuildInfo)

	return &Service{
		logger:    logger,
		factories: factoryMgr,
		config:    configMgr,
		telemetry: telemetryMgr,
		settings:  settings,
		buildInfo: opts.BuildInfo,
	}, nil
}

// Start initializes and starts the collector service
func (s *Service) Start(ctx context.Context) error {
	s.logger.Info("Starting collector service...")

	// Create and store component settings
	s.telemetrySettings = s.factories.CreateSettings(s.buildInfo)

	// Load configuration
	conf, err := s.config.Load(ctx)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Start telemetry
	if err := s.telemetry.Start(ctx); err != nil {
		return fmt.Errorf("failed to start telemetry: %w", err)
	}

	// TODO: Initialize pipelines with configuration
	if conf != nil {
		s.logger.Info("Configuration loaded successfully")
	}

	s.logger.Info("Collector service started successfully")
	return nil
}

// Shutdown stops the collector service
func (s *Service) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down collector service...")

	// Stop telemetry
	if err := s.telemetry.Stop(ctx); err != nil {
		s.logger.Error("Failed to stop telemetry", zap.Error(err))
	}

	s.logger.Info("Collector service shutdown complete")
	return nil
}
