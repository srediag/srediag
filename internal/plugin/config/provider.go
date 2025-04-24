package config

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/collector/confmap"
	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/types"
)

const providerScheme = "plugin"

// pluginProvider implements confmap.Provider for plugin configuration
type pluginProvider struct {
	mu      sync.RWMutex
	logger  *zap.Logger
	plugins map[string]types.IPlugin
}

// NewPluginProviderFactory creates a new factory for the plugin configuration provider
func NewPluginProviderFactory() confmap.ProviderFactory {
	return confmap.NewProviderFactory(newPluginProvider)
}

func newPluginProvider(settings confmap.ProviderSettings) confmap.Provider {
	return &pluginProvider{
		logger:  settings.Logger,
		plugins: make(map[string]types.IPlugin),
	}
}

// RegisterPlugin registers a plugin with the provider
func (p *pluginProvider) RegisterPlugin(plugin types.IPlugin) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	id := plugin.GetName()
	if _, exists := p.plugins[id]; exists {
		return fmt.Errorf("plugin %q already registered", id)
	}

	p.plugins[id] = plugin
	return nil
}

// UnregisterPlugin removes a plugin from the provider
func (p *pluginProvider) UnregisterPlugin(id string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, exists := p.plugins[id]; !exists {
		return fmt.Errorf("plugin %q not found", id)
	}

	delete(p.plugins, id)
	return nil
}

// Retrieve implements confmap.Provider
func (p *pluginProvider) Retrieve(_ context.Context, uri string, watcher confmap.WatcherFunc) (*confmap.Retrieved, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// Remove scheme prefix if present
	pluginID := uri
	if scheme := p.Scheme() + ":"; len(uri) > len(scheme) && uri[:len(scheme)] == scheme {
		pluginID = uri[len(scheme):]
	}

	// Get plugin configuration
	plugin, exists := p.plugins[pluginID]
	if !exists {
		return nil, fmt.Errorf("plugin %q not found", pluginID)
	}

	// Get plugin configuration
	config := map[string]interface{}{
		"id":       plugin.GetName(),
		"name":     plugin.GetName(),
		"enabled":  true,
		"settings": make(map[string]interface{}),
	}

	return confmap.NewRetrieved(config)
}

// Scheme implements confmap.Provider
func (p *pluginProvider) Scheme() string {
	return providerScheme
}

// Shutdown implements confmap.Provider
func (p *pluginProvider) Shutdown(context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.plugins = nil
	return nil
}
