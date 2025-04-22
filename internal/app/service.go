package app

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/config"
	"github.com/srediag/srediag/internal/plugins"
	"github.com/srediag/srediag/internal/telemetry"
)

// Service represents the main SREDIAG service
type Service struct {
	logger     *zap.Logger
	config     *config.Config
	pluginMgr  *plugins.Manager
	telemetry  *telemetry.Manager
	mu         sync.RWMutex
	isRunning  bool
	cancelFunc context.CancelFunc
}

// NewService creates a new instance of the SREDIAG service
func NewService(cfg *config.Config, logger *zap.Logger) (*Service, error) {
	if cfg == nil {
		return nil, fmt.Errorf("configuration cannot be nil")
	}

	return &Service{
		logger: logger,
		config: cfg,
	}, nil
}

// Start initializes and starts the SREDIAG service
func (s *Service) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		return fmt.Errorf("service is already running")
	}

	ctx, cancel := context.WithCancel(ctx)
	s.cancelFunc = cancel

	// Initialize plugin manager
	pluginMgr := plugins.NewManager(s.config.Plugins, s.logger)
	s.pluginMgr = pluginMgr

	// Initialize telemetry manager
	telemetry, err := telemetry.NewManager(s.config.Telemetry, s.config.Version, s.logger)
	if err != nil {
		return fmt.Errorf("failed to initialize telemetry manager: %w", err)
	}
	s.telemetry = telemetry

	// Start components
	if err := s.pluginMgr.Start(ctx); err != nil {
		return fmt.Errorf("failed to start plugin manager: %w", err)
	}

	if err := s.telemetry.Start(ctx); err != nil {
		return fmt.Errorf("failed to start telemetry manager: %w", err)
	}

	s.isRunning = true
	s.logger.Info("SREDIAG service started successfully")

	return nil
}

// Stop gracefully stops the SREDIAG service
func (s *Service) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return nil
	}

	if s.cancelFunc != nil {
		s.cancelFunc()
	}

	// Stop components
	if err := s.pluginMgr.Stop(ctx); err != nil {
		s.logger.Error("Error stopping plugin manager", zap.Error(err))
	}

	if err := s.telemetry.Stop(ctx); err != nil {
		s.logger.Error("Error stopping telemetry manager", zap.Error(err))
	}

	s.isRunning = false
	s.logger.Info("SREDIAG service stopped successfully")

	return nil
}

// IsRunning returns whether the service is currently running
func (s *Service) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isRunning
}
