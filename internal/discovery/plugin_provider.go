package discovery

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"plugin"
	"sync"

	"go.uber.org/zap"

	"github.com/srediag/srediag/internal/types"
)

// PluginProvider implements service discovery for plugins
type PluginProvider struct {
	mu         sync.RWMutex
	logger     *zap.Logger
	pluginDir  string
	discovered map[string]ServiceInfo
	extensions []string
	symbolName string
}

// NewPluginProvider creates a new plugin discovery provider
func NewPluginProvider(logger *zap.Logger, pluginDir string) *PluginProvider {
	return &PluginProvider{
		logger:     logger,
		pluginDir:  pluginDir,
		discovered: make(map[string]ServiceInfo),
		extensions: []string{".so", ".dylib", ".dll"},
		symbolName: "Plugin",
	}
}

// Start implements DiscoveryProvider
func (p *PluginProvider) Start(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, err := os.Stat(p.pluginDir); os.IsNotExist(err) {
		return fmt.Errorf("plugin directory %s does not exist", p.pluginDir)
	}

	p.logger.Info("Started plugin discovery provider",
		zap.String("plugin_dir", p.pluginDir))
	return nil
}

// Stop implements DiscoveryProvider
func (p *PluginProvider) Stop(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.discovered = make(map[string]ServiceInfo)
	p.logger.Info("Stopped plugin discovery provider")
	return nil
}

// Discover implements DiscoveryProvider
func (p *PluginProvider) Discover(ctx context.Context) ([]ServiceInfo, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	discovered := make(map[string]ServiceInfo)

	err := filepath.WalkDir(p.pluginDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		if !p.isValidExtension(ext) {
			return nil
		}

		pluginInfo, err := p.loadPlugin(path)
		if err != nil {
			p.logger.Error("Failed to load plugin",
				zap.String("path", path),
				zap.Error(err))
			return nil
		}

		discovered[pluginInfo.ID] = pluginInfo
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk plugin directory: %w", err)
	}

	// Update discovered plugins
	p.discovered = discovered

	// Convert to slice
	services := make([]ServiceInfo, 0, len(discovered))
	for _, service := range discovered {
		services = append(services, service)
	}

	return services, nil
}

// isValidExtension checks if the file extension is valid for plugins
func (p *PluginProvider) isValidExtension(ext string) bool {
	for _, validExt := range p.extensions {
		if ext == validExt {
			return true
		}
	}
	return false
}

// loadPlugin loads a plugin and extracts its information
func (p *PluginProvider) loadPlugin(path string) (ServiceInfo, error) {
	plug, err := plugin.Open(path)
	if err != nil {
		return ServiceInfo{}, fmt.Errorf("failed to open plugin: %w", err)
	}

	symbol, err := plug.Lookup(p.symbolName)
	if err != nil {
		return ServiceInfo{}, fmt.Errorf("plugin symbol %s not found: %w", p.symbolName, err)
	}

	plugin, ok := symbol.(types.IPlugin)
	if !ok {
		return ServiceInfo{}, fmt.Errorf("invalid plugin type: %T", symbol)
	}

	return ServiceInfo{
		ID:      plugin.GetName(),
		Name:    plugin.GetName(),
		Type:    string(plugin.GetCategory()),
		Version: plugin.GetVersion(),
		Metadata: map[string]string{
			"category": string(plugin.GetCategory()),
		},
	}, nil
}

// SetExtensions sets the valid plugin file extensions
func (p *PluginProvider) SetExtensions(extensions []string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.extensions = extensions
}

// SetSymbolName sets the plugin symbol name to look for
func (p *PluginProvider) SetSymbolName(name string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.symbolName = name
}
