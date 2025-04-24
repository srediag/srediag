package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/types"
)

// ConfigProvider defines the interface for dynamic configuration providers
type ConfigProvider interface {
	// Start starts the configuration provider
	Start(ctx context.Context) error
	// Stop stops the configuration provider
	Stop(ctx context.Context) error
	// Get retrieves the current configuration
	Get() (*types.Config, error)
	// Watch returns a channel that receives configuration updates
	Watch() <-chan *types.Config
	// IsWatching returns true if the provider is actively watching for changes
	IsWatching() bool
}

// DynamicProvider implements a dynamic configuration provider
type DynamicProvider struct {
	mu            sync.RWMutex
	logger        *zap.Logger
	config        *types.Config
	watchChan     chan *types.Config
	watchInterval time.Duration
	watching      bool
	configPath    string
}

// NewDynamicProvider creates a new dynamic configuration provider
func NewDynamicProvider(logger *zap.Logger, configPath string, watchInterval time.Duration) *DynamicProvider {
	if watchInterval == 0 {
		watchInterval = 30 * time.Second
	}

	return &DynamicProvider{
		logger:        logger,
		configPath:    configPath,
		watchInterval: watchInterval,
		watchChan:     make(chan *types.Config),
	}
}

// Start implements ConfigProvider
func (p *DynamicProvider) Start(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.watching {
		return fmt.Errorf("provider already started")
	}

	// Load initial configuration
	if err := p.loadConfig(); err != nil {
		return err
	}

	p.watching = true
	go p.watchConfig(ctx)

	p.logger.Info("Started dynamic configuration provider",
		zap.String("config_path", p.configPath),
		zap.Duration("watch_interval", p.watchInterval))

	return nil
}

// Stop implements ConfigProvider
func (p *DynamicProvider) Stop(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.watching {
		return nil
	}

	p.watching = false
	close(p.watchChan)

	p.logger.Info("Stopped dynamic configuration provider")
	return nil
}

// Get implements ConfigProvider
func (p *DynamicProvider) Get() (*types.Config, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.config == nil {
		return nil, fmt.Errorf("configuration not loaded")
	}

	return p.config, nil
}

// Watch implements ConfigProvider
func (p *DynamicProvider) Watch() <-chan *types.Config {
	return p.watchChan
}

// IsWatching implements ConfigProvider
func (p *DynamicProvider) IsWatching() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.watching
}

// loadConfig loads the configuration from file
func (p *DynamicProvider) loadConfig() error {
	data, err := os.ReadFile(p.configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var newConfig types.Config
	if err := json.Unmarshal(data, &newConfig); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := newConfig.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	p.config = &newConfig
	return nil
}

// watchConfig watches for configuration changes
func (p *DynamicProvider) watchConfig(ctx context.Context) {
	ticker := time.NewTicker(p.watchInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.mu.Lock()
			if !p.watching {
				p.mu.Unlock()
				return
			}

			oldConfig := p.config
			if err := p.loadConfig(); err != nil {
				p.logger.Error("Failed to reload configuration",
					zap.Error(err))
				p.mu.Unlock()
				continue
			}

			if p.configChanged(oldConfig, p.config) {
				p.logger.Info("Configuration changed, notifying subscribers")
				p.watchChan <- p.config
			}
			p.mu.Unlock()
		}
	}
}

// configChanged checks if the configuration has changed
func (p *DynamicProvider) configChanged(old, new *types.Config) bool {
	if old == nil || new == nil {
		return old != new
	}

	// Add your comparison logic here
	// For now, we'll consider any change as significant
	return true
}
