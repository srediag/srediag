package core

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"

	"github.com/srediag/srediag/pkg/plugins/collector/processor/diagnostic"
	"github.com/srediag/srediag/pkg/plugins/collector/receiver/system"
)

// Runner manages the collector lifecycle
type Runner struct {
	logger    *zap.Logger
	collector *Collector
	config    *Config
}

// NewRunner creates a new runner instance
func NewRunner(configPath string) (*Runner, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	collector, err := NewCollector(logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create collector: %w", err)
	}

	return &Runner{
		logger:    logger,
		collector: collector,
		config:    config,
	}, nil
}

// Run starts the collector and manages its lifecycle
func (r *Runner) Run(ctx context.Context) error {
	// Register built-in factories
	if err := r.registerFactories(); err != nil {
		return fmt.Errorf("failed to register factories: %w", err)
	}

	// Create and initialize components
	if err := r.initializeComponents(); err != nil {
		return fmt.Errorf("failed to initialize components: %w", err)
	}

	// Setup signal handling
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Start the collector
	if err := r.collector.Start(ctx); err != nil {
		return fmt.Errorf("failed to start collector: %w", err)
	}

	// Wait for shutdown signal
	<-signalChan
	r.logger.Info("Received shutdown signal")

	// Initiate graceful shutdown
	if err := r.collector.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown collector: %w", err)
	}

	return nil
}

// registerFactories registers all built-in component factories
func (r *Runner) registerFactories() error {
	factories := []struct {
		name    string
		factory component.Factory
	}{
		{"diagnostic", diagnostic.NewFactory()},
		{"system", system.NewFactory()},
	}

	for _, f := range factories {
		if err := r.collector.RegisterFactory(f.factory); err != nil {
			return fmt.Errorf("failed to register %s factory: %w", f.name, err)
		}
		r.logger.Info("Registered factory", zap.String("name", f.name))
	}

	return nil
}

// initializeComponents creates component instances based on configuration
func (r *Runner) initializeComponents() error {
	for name, compConfig := range r.config.Components {
		if !compConfig.Enabled {
			r.logger.Info("Skipping disabled component", zap.String("name", name))
			continue
		}

		id := component.NewIDWithName(component.MustNewType(compConfig.Type), name)
		if _, exists := r.collector.GetComponent(id); !exists {
			return fmt.Errorf("failed to initialize component %s: component not found", name)
		}

		r.logger.Info("Initialized component",
			zap.String("name", name),
			zap.String("type", compConfig.Type))
	}

	return nil
}
