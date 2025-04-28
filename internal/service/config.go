package service

import (
	"fmt"

	"github.com/srediag/srediag/internal/core"
)

// Package service provides service configuration management and helpers for the SREDIAG collector service.
//
// This file defines the canonical service config structure and helpers for loading and overlaying YAML config.
//
// Usage:
//   - Use ServiceConfig to represent the canonical collector config structure for SREDIAG.
//   - Use LoadServiceConfig to load and overlay service config from YAML, environment, and CLI flags.
//
// Best Practices:
//   - Always validate required fields after loading config.
//   - Use canonical types for all new code.
//
// TODO:
//   - Add context.Context support for cancellation and timeouts.
//   - Add stricter schema validation and error reporting.

// ServiceConfig represents the canonical service (collector) config structure for SREDIAG.
//
// Fields:
//   - Service: Top-level service configuration (e.g., pipelines, telemetry, extensions).
//   - Receivers: Map of receiver component configurations.
//   - Processors: Map of processor component configurations.
//   - Exporters: Map of exporter component configurations.
//   - Extensions: Map of extension component configurations.
type ServiceConfig struct {
	// Service contains top-level service configuration (e.g., pipelines, telemetry, extensions).
	Service map[string]interface{} `yaml:"service"`
	// Receivers contains receiver component configurations.
	Receivers map[string]interface{} `yaml:"receivers"`
	// Processors contains processor component configurations.
	Processors map[string]interface{} `yaml:"processors"`
	// Exporters contains exporter component configurations.
	Exporters map[string]interface{} `yaml:"exporters"`
	// Extensions contains extension component configurations.
	Extensions map[string]interface{} `yaml:"extensions"`
}

// LoadServiceConfig loads and overlays the service config from configs/srediag-service.yaml.
//
// Precedence: CLI flags > env > YAML > built-ins.
//
// Parameters:
//   - cliFlags: Map of CLI flag overrides.
//
// Returns:
//   - *ServiceConfig: The loaded service configuration.
//   - error: If loading or validation fails, returns a detailed error.
func LoadServiceConfig(cliFlags map[string]string) (*ServiceConfig, error) {
	var cfg ServiceConfig
	if err := core.LoadConfigWithOverlay(&cfg, cliFlags, core.WithConfigPathSuffix("service.yaml")); err != nil {
		return nil, fmt.Errorf("failed to load service config: %w", err)
	}
	return &cfg, nil
}
