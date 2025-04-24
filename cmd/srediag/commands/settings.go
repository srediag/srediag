package commands

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/connector"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/extension"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/zap"
)

// Settings contains the configuration for the command executor
type Settings struct {
	// Component factories
	Receivers  map[component.Type]receiver.Factory
	Processors map[component.Type]processor.Factory
	Exporters  map[component.Type]exporter.Factory
	Extensions map[component.Type]extension.Factory
	Connectors map[component.Type]connector.Factory

	// Logger
	Logger *zap.Logger
}
