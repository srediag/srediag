package core

import (
	"context"

	"go.opentelemetry.io/collector/otelcol"
	"go.uber.org/zap"
)

// ISREDiagConfig defines the system configuration interface
type ISREDiagConfig interface {
	// GetVersion returns the system version
	GetVersion() string
	// GetServiceConfig returns the service configuration
	GetServiceConfig() ISREDiagServiceConfig
	// GetCollectorConfig returns the collector configuration
	GetCollectorConfig() ISREDiagCollectorConfig
	// Validate validates the configuration
	Validate() error
}

// ISREDiagServiceConfig defines the service configuration interface
type ISREDiagServiceConfig interface {
	// GetName returns the service name
	GetName() string
	// GetEnvironment returns the service environment
	GetEnvironment() string
}

// ISREDiagCollectorConfig defines the collector configuration interface
type ISREDiagCollectorConfig interface {
	// IsEnabled returns if the collector is enabled
	IsEnabled() bool
	// GetConfigPath returns the configuration file path
	GetConfigPath() string
}

// ISREDiagRunner defines the system runner interface
type ISREDiagRunner interface {
	// Start starts the runner with the given context
	Start(ctx context.Context) error
	// Stop stops the runner with the given context
	Stop(ctx context.Context) error
	// GetLogger returns the configured logger
	GetLogger() *zap.Logger
	// GetCollector returns the OpenTelemetry collector instance if enabled
	GetCollector() *otelcol.Collector
}
