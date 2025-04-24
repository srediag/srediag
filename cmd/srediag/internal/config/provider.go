package config

import (
	"context"
	"sync"

	"github.com/spf13/viper"
	"go.opentelemetry.io/collector/confmap"
)

// CLIProvider implements the confmap.Provider interface for CLI configuration
type CLIProvider struct {
	mu       sync.RWMutex
	watchers []confmap.WatcherFunc
}

// NewCLIProvider creates a new CLI configuration provider
func NewCLIProvider() *CLIProvider {
	return &CLIProvider{
		watchers: make([]confmap.WatcherFunc, 0),
	}
}

// Retrieve implements confmap.Provider
func (p *CLIProvider) Retrieve(_ context.Context, _ string, watcher confmap.WatcherFunc) (*confmap.Retrieved, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if watcher != nil {
		p.watchers = append(p.watchers, watcher)
	}

	// Get all settings from viper
	settings := viper.AllSettings()
	if len(settings) == 0 {
		return confmap.NewRetrieved(nil)
	}

	return confmap.NewRetrieved(settings)
}

// Scheme implements confmap.Provider
func (p *CLIProvider) Scheme() string {
	return "cli"
}

// Shutdown implements confmap.Provider
func (p *CLIProvider) Shutdown(context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.watchers = nil
	return nil
}

// NotifyConfigChange notifies all watchers of a configuration change
func (p *CLIProvider) NotifyConfigChange() {
	p.mu.RLock()
	defer p.mu.RUnlock()

	for _, watcher := range p.watchers {
		watcher(&confmap.ChangeEvent{})
	}
}
