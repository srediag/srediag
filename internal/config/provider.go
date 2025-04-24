package config

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.uber.org/zap"
)

// Provider manages configuration loading and reloading
type Provider struct {
	mu            sync.RWMutex
	logger        *zap.Logger
	providers     map[string]confmap.Provider
	watchers      []confmap.WatcherFunc
	validators    []ConfigValidator
	settings      component.TelemetrySettings
	reloadPeriod  time.Duration
	lastReload    time.Time
	currentConfig *confmap.Conf
}

// NewProvider creates a new configuration provider
func NewProvider(logger *zap.Logger, settings component.TelemetrySettings) *Provider {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &Provider{
		logger:       logger,
		providers:    make(map[string]confmap.Provider),
		validators:   make([]ConfigValidator, 0),
		settings:     settings,
		reloadPeriod: 30 * time.Second,
	}
}

// RegisterProvider registers a configuration provider
func (p *Provider) RegisterProvider(name string, provider confmap.Provider) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, exists := p.providers[name]; exists {
		return fmt.Errorf("provider %q already registered", name)
	}

	p.providers[name] = provider
	p.logger.Info("Registered configuration provider", zap.String("name", name))
	return nil
}

// UnregisterProvider removes a configuration provider
func (p *Provider) UnregisterProvider(name string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, exists := p.providers[name]; !exists {
		return fmt.Errorf("provider %q not found", name)
	}

	delete(p.providers, name)
	p.logger.Info("Unregistered configuration provider", zap.String("name", name))
	return nil
}

// AddValidator adds a configuration validator
func (p *Provider) AddValidator(validator ConfigValidator) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.validators = append(p.validators, validator)
}

// AddWatcher adds a configuration change watcher
func (p *Provider) AddWatcher(watcher confmap.WatcherFunc) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.watchers = append(p.watchers, watcher)
}

// Load loads configuration from all registered providers
func (p *Provider) Load(ctx context.Context) (*confmap.Conf, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	var configs []*confmap.Conf
	for name, provider := range p.providers {
		retrieved, err := provider.Retrieve(ctx, provider.Scheme()+":"+name, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to load config from provider %q: %w", name, err)
		}

		raw, err := retrieved.AsRaw()
		if err != nil {
			return nil, fmt.Errorf("failed to get raw config from provider %q: %w", name, err)
		}

		if rawMap, ok := raw.(map[string]interface{}); ok {
			conf := confmap.NewFromStringMap(rawMap)
			configs = append(configs, conf)
		} else {
			return nil, fmt.Errorf("invalid configuration format from provider %q", name)
		}
	}

	if len(configs) == 0 {
		return confmap.New(), nil
	}

	merged := configs[0]
	for _, conf := range configs[1:] {
		if err := merged.Merge(conf); err != nil {
			return nil, fmt.Errorf("failed to merge configurations: %w", err)
		}
	}

	// Validate configuration
	for _, validator := range p.validators {
		if err := validator.Validate(merged); err != nil {
			return nil, fmt.Errorf("configuration validation failed: %w", err)
		}
	}

	if err := p.updateConfig(merged); err != nil {
		return nil, err
	}

	return merged, nil
}

// updateConfig updates the current configuration and notifies watchers
func (p *Provider) updateConfig(conf *confmap.Conf) error {
	p.currentConfig = conf
	p.lastReload = time.Now()

	// Notify watchers
	for _, watcher := range p.watchers {
		watcher(&confmap.ChangeEvent{
			Error: nil,
		})
	}

	return nil
}

// Get returns the current configuration
func (p *Provider) Get() *confmap.Conf {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.currentConfig
}

// Watch starts watching for configuration changes
func (p *Provider) Watch(ctx context.Context) error {
	ticker := time.NewTicker(p.reloadPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if _, err := p.Load(ctx); err != nil {
				p.logger.Error("Failed to reload configuration",
					zap.Error(err))
			}
		}
	}
}

// SetReloadPeriod sets the configuration reload period
func (p *Provider) SetReloadPeriod(period time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.reloadPeriod = period
}

// LastReloadTime returns the time of the last configuration reload
func (p *Provider) LastReloadTime() time.Time {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.lastReload
}

// Shutdown cleans up resources
func (p *Provider) Shutdown(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	for name, provider := range p.providers {
		if err := provider.Shutdown(ctx); err != nil {
			p.logger.Error("Failed to shutdown provider",
				zap.String("name", name),
				zap.Error(err))
		}
	}

	p.providers = make(map[string]confmap.Provider)
	p.watchers = nil
	p.validators = nil
	p.currentConfig = nil

	return nil
}
