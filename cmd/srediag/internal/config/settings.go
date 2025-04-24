package config

import (
	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"
)

// DefaultSettings returns the default telemetry settings
func DefaultSettings() component.TelemetrySettings {
	return component.TelemetrySettings{
		Logger: zap.NewNop(),
	}
}
