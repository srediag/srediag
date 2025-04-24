package core

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/otel/metric/noop"
	tracenoop "go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/discovery"
	"github.com/srediag/srediag/internal/plugin"
)

const defaultDiscoveryInterval = 30 * time.Second

// Collector represents the core collector instance
type Collector struct {
	logger  *zap.Logger
	manager *plugin.Manager
	host    component.Host
}

// NewCollector creates a new collector instance
func NewCollector(logger *zap.Logger) (*Collector, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	host := componenttest.NewNopHost()
	registry := plugin.NewRegistry(logger)
	discoveryMgr := discovery.NewManager(logger, defaultDiscoveryInterval)
	buildInfo := component.NewDefaultBuildInfo()
	tracerProvider := tracenoop.NewTracerProvider()
	meterProvider := noop.NewMeterProvider()

	manager := plugin.NewManager(
		logger,
		registry,
		discoveryMgr,
		host,
		buildInfo,
		tracerProvider,
		meterProvider,
	)

	return &Collector{
		logger:  logger,
		manager: manager,
		host:    host,
	}, nil
}

// Start initializes and starts all components
func (c *Collector) Start(ctx context.Context) error {
	c.logger.Info("Starting collector...")

	if err := c.manager.Start(ctx); err != nil {
		return fmt.Errorf("failed to start plugins: %w", err)
	}

	c.logger.Info("Collector started successfully")
	return nil
}

// Shutdown stops all components
func (c *Collector) Shutdown(ctx context.Context) error {
	c.logger.Info("Shutting down collector...")

	if err := c.manager.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown plugins: %w", err)
	}

	c.logger.Info("Collector shutdown complete")
	return nil
}

// RegisterFactory registers a new component factory
func (c *Collector) RegisterFactory(factory component.Factory) error {
	return c.manager.RegisterFactory(factory)
}

// GetComponent retrieves a component by ID
func (c *Collector) GetComponent(id component.ID) (component.Component, bool) {
	return c.manager.GetComponent(id)
}

// ListComponents returns all registered components
func (c *Collector) ListComponents() []component.Component {
	return c.manager.ListComponents()
}
