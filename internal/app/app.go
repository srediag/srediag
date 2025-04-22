package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/config"
	"github.com/srediag/srediag/internal/plugins"
)

type App struct {
	cfg     *config.Config
	logger  *zap.Logger
	plugins *plugins.Manager
}

func New(cfg *config.Config, logger *zap.Logger) *App {
	return &App{
		cfg:     cfg,
		logger:  logger,
		plugins: plugins.NewManager(cfg.Plugins, logger),
	}
}

func (a *App) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Setup signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Start plugins
	if err := a.plugins.Start(ctx); err != nil {
		return err
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
	if err := a.plugins.Stop(ctx); err != nil {
		a.logger.Error("error stopping plugins", zap.Error(err))
	}

	return nil
}

func (a *App) Stop() error {
	a.logger.Info("stopping application")
	return nil
}
