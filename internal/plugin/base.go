// Package plugin provides plugin types and utilities for SREDIAG
package plugin

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/types"
)

// Base provides a base implementation of types.IPlugin interface
type Base struct {
	logger       *zap.Logger
	meter        metric.Meter
	metadata     types.PluginMetadata
	config       types.PluginConfig
	lifecycle    types.PluginLifecycle
	capabilities types.PluginCapabilities
	mu           sync.RWMutex
	healthy      bool
	running      bool
	lastError    error
}

// NewBase creates a new base plugin instance
func NewBase(logger *zap.Logger, meter metric.Meter, metadata types.PluginMetadata) *Base {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &Base{
		logger:       logger,
		meter:        meter,
		metadata:     metadata,
		lifecycle:    types.LifecycleUnregistered,
		capabilities: metadata.Capabilities,
		healthy:      true,
	}
}

// GetID implements types.IPlugin
func (b *Base) GetID() string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.metadata.ID
}

// GetName implements types.IPlugin
func (b *Base) GetName() string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.metadata.Name
}

// GetVersion implements types.IPlugin
func (b *Base) GetVersion() string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.metadata.Version
}

// GetCategory implements types.IPlugin
func (b *Base) GetCategory() types.PluginCategory {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.metadata.Category
}

// GetCapabilities implements types.IPlugin
func (b *Base) GetCapabilities() types.PluginCapabilities {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.capabilities
}

// Initialize implements types.IPlugin
func (b *Base) Initialize(ctx context.Context, config types.PluginConfig) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.lifecycle != types.LifecycleUnregistered {
		return fmt.Errorf("plugin must be in unregistered state to initialize")
	}

	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	b.config = config
	b.lifecycle = types.LifecycleInitialized
	return nil
}

// Start implements types.IPlugin
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
	b.lifecycle = types.LifecycleRunning
	return nil
}

// Stop implements types.IPlugin
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
	b.lifecycle = types.LifecycleStopped
	return nil
}

// Pause implements types.IPlugin
func (b *Base) Pause(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.running {
		return fmt.Errorf("plugin is not running")
	}

	b.lifecycle = types.LifecyclePaused
	return nil
}

// Resume implements types.IPlugin
func (b *Base) Resume(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.lifecycle != types.LifecyclePaused {
		return fmt.Errorf("plugin is not paused")
	}

	b.lifecycle = types.LifecycleRunning
	return nil
}

// GetLifecycle implements types.IPlugin
func (b *Base) GetLifecycle() types.PluginLifecycle {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.lifecycle
}

// IsHealthy implements types.IPlugin
func (b *Base) IsHealthy() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.healthy
}

// GetLastError implements types.IPlugin
func (b *Base) GetLastError() error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.lastError
}

// GetConfig implements types.IPlugin
func (b *Base) GetConfig() types.PluginConfig {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.config
}

// GetMetadata implements types.IPlugin
func (b *Base) GetMetadata() types.PluginMetadata {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.metadata
}

// Validate implements types.IPlugin
func (b *Base) Validate() error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.metadata.ID == "" {
		return fmt.Errorf("plugin ID is required")
	}
	if b.metadata.Name == "" {
		return fmt.Errorf("plugin name is required")
	}
	if b.metadata.Version == "" {
		return fmt.Errorf("plugin version is required")
	}
	return nil
}

// SetHealth sets the plugin health status
func (b *Base) SetHealth(healthy bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.healthy = healthy
}

// SetError sets the last error encountered by the plugin
func (b *Base) SetError(err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.lastError = err
}

// GetLogger returns the plugin logger
func (b *Base) GetLogger() *zap.Logger {
	return b.logger
}

// GetMeter returns the plugin meter
func (b *Base) GetMeter() metric.Meter {
	return b.meter
}
