// Package app provides the main application functionality
package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/config"
)

// App represents the main application
type App struct {
	cfg     *config.ConfigRoot
	logger  *zap.Logger
	srediag *SREDiag
}

// New creates a new instance of App
func New(cfg *config.ConfigRoot, logger *zap.Logger) (*App, error) {
	if cfg == nil {
		return nil, fmt.Errorf("configuration cannot be nil")
	}
	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	srediag, err := NewSREDiag(logger, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create srediag instance: %w", err)
	}

	return &App{
		cfg:     cfg,
		logger:  logger,
		srediag: srediag,
	}, nil
}

// Start initializes and starts the application
func (a *App) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Setup signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Start SREDiag
	if err := a.srediag.Start(ctx); err != nil {
		return fmt.Errorf("failed to start srediag: %w", err)
	}

	// Wait for shutdown signal
	select {
	case sig := <-sigCh:
		a.logger.Info("received shutdown signal", zap.String("signal", sig.String()))
	case <-ctx.Done():
		a.logger.Info("context cancelled")
	}

	// Graceful shutdown
	a.logger.Info("initiating graceful shutdown")
	if err := a.srediag.Stop(ctx); err != nil {
		a.logger.Error("error during shutdown", zap.Error(err))
	}

	return nil
}

// Stop gracefully stops the application
func (a *App) Stop(ctx context.Context) error {
	a.logger.Info("stopping application")
	return a.srediag.Stop(ctx)
}
