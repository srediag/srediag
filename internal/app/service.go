// Package app provides the main application functionality
package app

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/config"
	"github.com/srediag/srediag/internal/diagnostic"
	"github.com/srediag/srediag/internal/types"
)

// Service represents the main SREDIAG service
type Service struct {
	logger     *zap.Logger
	config     *config.ConfigRoot
	pluginMgr  types.IPluginManager
	telemetry  types.ITelemetryBridge
	mu         sync.RWMutex
	isRunning  bool
	cancelFunc context.CancelFunc
}

// NewService creates a new instance of the SREDIAG service
func NewService(cfg *config.ConfigRoot, logger *zap.Logger) (*Service, error) {
	if cfg == nil {
		return nil, fmt.Errorf("configuration cannot be nil")
	}
	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
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
	pluginMgr := diagnostic.NewPluginManager(s.logger.Named("plugin-manager"))
	s.pluginMgr = pluginMgr

	// Initialize telemetry bridge
	telemetry := diagnostic.NewTelemetryBridge(
		s.logger.Named("telemetry-bridge"),
		nil, // Resource will be created in SREDiag
	)
	s.telemetry = telemetry

	// Start components
	if err := s.pluginMgr.StartAll(ctx); err != nil {
		return fmt.Errorf("failed to start plugin manager: %w", err)
	}

	if err := s.telemetry.Start(ctx); err != nil {
		return fmt.Errorf("failed to start telemetry bridge: %w", err)
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

	// Stop components in reverse order
	if err := s.telemetry.Stop(ctx); err != nil {
		s.logger.Error("error stopping telemetry bridge", zap.Error(err))
	}

	if err := s.pluginMgr.StopAll(ctx); err != nil {
		s.logger.Error("error stopping plugin manager", zap.Error(err))
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
