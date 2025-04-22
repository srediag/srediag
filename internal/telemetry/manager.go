package telemetry

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/trace"
	nooptrace "go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/config"
	"github.com/srediag/srediag/internal/lifecycle"
)

// Manager handles telemetry operations
type Manager struct {
	*lifecycle.BaseManager
	config     config.TelemetryConfig
	version    string
	logger     *zap.Logger
	tracerProv trace.TracerProvider
	meterProv  metric.MeterProvider
}

// NewManager creates a new telemetry manager
func NewManager(cfg config.TelemetryConfig, version string, logger *zap.Logger) (*Manager, error) {
	if !cfg.Enabled {
		logger.Info("telemetry is disabled")
		return &Manager{
			BaseManager: lifecycle.NewBaseManager(),
			config:      cfg,
			version:     version,
			logger:      logger,
		}, nil
	}

	return &Manager{
		BaseManager: lifecycle.NewBaseManager(),
		config:      cfg,
		version:     version,
		logger:      logger,
	}, nil
}

// Start initializes the telemetry providers
func (m *Manager) Start(ctx context.Context) error {
	if err := m.CheckRunningState(false); err != nil {
		return err
	}

	if !m.config.Enabled {
		m.SetRunning(true)
		return nil
	}

	// Initialize providers
	if err := m.initProviders(ctx); err != nil {
		return fmt.Errorf("failed to initialize providers: %w", err)
	}

	m.SetRunning(true)
	m.logger.Info("telemetry manager started")
	return nil
}

// Stop shuts down the telemetry providers
func (m *Manager) Stop(ctx context.Context) error {
	if err := m.CheckRunningState(true); err != nil {
		return err
	}

	if !m.config.Enabled {
		m.SetRunning(false)
		return nil
	}

	// Shutdown providers
	if err := m.shutdownProviders(ctx); err != nil {
		return fmt.Errorf("failed to shutdown providers: %w", err)
	}

	m.SetRunning(false)
	m.logger.Info("telemetry manager stopped")
	return nil
}

// Tracer returns a named tracer
func (m *Manager) Tracer(name string) trace.Tracer {
	if !m.config.Enabled || !m.IsRunning() {
		return nooptrace.NewTracerProvider().Tracer(name)
	}
	return m.tracerProv.Tracer(name)
}

// Meter returns a named meter
func (m *Manager) Meter(name string) metric.Meter {
	if !m.config.Enabled || !m.IsRunning() {
		return noop.NewMeterProvider().Meter(name)
	}
	return m.meterProv.Meter(name)
}

// initProviders initializes the telemetry providers
func (m *Manager) initProviders(ctx context.Context) error {
	// Check if context is already done before proceeding
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("context error before initializing providers: %w", err)
	}

	if m.config.Traces.Enabled {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled while initializing tracer provider: %w", ctx.Err())
		default:
			// Use default tracer provider for now
			// TODO: Implement custom tracer provider configuration
			m.tracerProv = otel.GetTracerProvider()
			m.logger.Info("traces enabled, using default tracer provider")
		}
	} else {
		m.tracerProv = nooptrace.NewTracerProvider()
		m.logger.Info("traces disabled, using noop tracer provider")
	}

	if m.config.Metrics.Enabled {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled while initializing meter provider: %w", ctx.Err())
		default:
			// Use default meter provider for now
			// TODO: Implement custom meter provider configuration
			m.meterProv = otel.GetMeterProvider()
			m.logger.Info("metrics enabled, using default meter provider")
		}
	} else {
		m.meterProv = noop.NewMeterProvider()
		m.logger.Info("metrics disabled, using noop meter provider")
	}

	// Set global providers
	otel.SetTracerProvider(m.tracerProv)
	otel.SetMeterProvider(m.meterProv)

	return nil
}

// shutdownProviders shuts down the telemetry providers
func (m *Manager) shutdownProviders(ctx context.Context) error {
	var errs []error

	// Shutdown tracer provider if it implements Shutdownable interface
	if provider, ok := m.tracerProv.(interface{ Shutdown(context.Context) error }); ok {
		if err := provider.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to shutdown tracer provider: %w", err))
		} else {
			m.logger.Info("tracer provider shutdown successfully")
		}
	}

	// Shutdown meter provider if it implements Shutdownable interface
	if provider, ok := m.meterProv.(interface{ Shutdown(context.Context) error }); ok {
		if err := provider.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to shutdown meter provider: %w", err))
		} else {
			m.logger.Info("meter provider shutdown successfully")
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors during provider shutdown: %v", errs)
	}
	return nil
}

// SetConfig updates the telemetry configuration
func (m *Manager) SetConfig(cfg config.TelemetryConfig) {
	m.config = cfg
}
