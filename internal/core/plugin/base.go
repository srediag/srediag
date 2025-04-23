// Package plugin provides plugin types and utilities for SREDIAG
package plugin

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/core"
)

// Base provides a base implementation of core.Plugin
type Base struct {
	logger       *zap.Logger
	meter        metric.Meter
	metadata     Metadata
	config       Config
	status       core.Status
	capabilities []core.Capability
	mu           sync.RWMutex
	healthy      bool
	running      bool
}

// NewBase creates a new base plugin instance
func NewBase(logger *zap.Logger, meter metric.Meter, metadata Metadata) *Base {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &Base{
		logger:       logger,
		meter:        meter,
		metadata:     metadata,
		status:       core.StatusLoaded,
		capabilities: metadata.Capabilities,
		healthy:      true,
	}
}

// Start implements core.Plugin
func (b *Base) Start(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.running {
		return fmt.Errorf("plugin is already running")
	}

	b.logger.Info("starting plugin",
		zap.String("name", b.metadata.Name),
		zap.String("version", b.metadata.Version))

	b.running = true
	b.status = core.StatusRunning
	return nil
}

// Stop implements core.Plugin
func (b *Base) Stop(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.running {
		return nil
	}

	b.logger.Info("stopping plugin",
		zap.String("name", b.metadata.Name),
		zap.String("version", b.metadata.Version))

	b.running = false
	b.status = core.StatusStopped
	return nil
}

// IsHealthy implements core.Plugin
func (b *Base) IsHealthy() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.healthy
}

// GetName implements core.Plugin
func (b *Base) GetName() string {
	return b.metadata.Name
}

// GetVersion implements core.Plugin
func (b *Base) GetVersion() string {
	return b.metadata.Version
}

// GetType implements core.Plugin
func (b *Base) GetType() string {
	return string(b.metadata.Type)
}

// Configure implements core.Plugin
func (b *Base) Configure(cfg interface{}) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	config, ok := cfg.(Config)
	if !ok {
		return fmt.Errorf("invalid configuration type")
	}

	if config.Name != b.metadata.Name {
		return fmt.Errorf("configuration name mismatch")
	}

	if config.Type != b.metadata.Type {
		return fmt.Errorf("configuration type mismatch")
	}

	b.config = config
	return nil
}

// GetInfo returns plugin information
func (b *Base) GetInfo() Info {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return Info{
		Name:         b.metadata.Name,
		Version:      b.metadata.Version,
		Type:         b.metadata.Type,
		Capabilities: b.capabilities,
		Status:       b.status,
	}
}

// GetMetadata returns plugin metadata
func (b *Base) GetMetadata() Metadata {
	return b.metadata
}

// GetConfig returns plugin configuration
func (b *Base) GetConfig() Config {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.config
}

// GetStatus returns plugin status
func (b *Base) GetStatus() core.Status {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.status
}

// GetCapabilities returns plugin capabilities
func (b *Base) GetCapabilities() []core.Capability {
	return b.capabilities
}

// HasCapability checks if the plugin has a specific capability
func (b *Base) HasCapability(capability core.Capability) bool {
	for _, c := range b.capabilities {
		if c == capability {
			return true
		}
	}
	return false
}

// SetStatus sets the plugin status
func (b *Base) SetStatus(status core.Status) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.status = status
}

// SetHealth sets the plugin health status
func (b *Base) SetHealth(healthy bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.healthy = healthy
}

// GetLogger returns the plugin logger
func (b *Base) GetLogger() *zap.Logger {
	return b.logger
}

// GetMeter returns the plugin meter
func (b *Base) GetMeter() metric.Meter {
	return b.meter
}
