package provider

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"go.opentelemetry.io/collector/confmap"
	"gopkg.in/yaml.v3"
)

// YAMLProvider implements the confmap.Provider interface for YAML configuration files
type YAMLProvider struct {
	mu       sync.RWMutex
	watchers []confmap.WatcherFunc
	files    map[string]string // Map of file paths to their content hash
}

// NewYAMLProvider creates a new YAML configuration provider
func NewYAMLProvider() *YAMLProvider {
	return &YAMLProvider{
		watchers: make([]confmap.WatcherFunc, 0),
		files:    make(map[string]string),
	}
}

// Retrieve implements confmap.Provider
func (p *YAMLProvider) Retrieve(_ context.Context, uri string, watcher confmap.WatcherFunc) (*confmap.Retrieved, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if watcher != nil {
		p.watchers = append(p.watchers, watcher)
	}

	// Remove scheme prefix if present
	filePath := uri
	if scheme := p.Scheme() + ":"; len(uri) > len(scheme) && uri[:len(scheme)] == scheme {
		filePath = uri[len(scheme):]
	}

	// Resolve relative paths
	if !filepath.IsAbs(filePath) {
		wd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get working directory: %w", err)
		}
		filePath = filepath.Join(wd, filePath)
	}

	// Read and parse YAML file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %q: %w", filePath, err)
	}

	// Parse YAML
	var config map[string]interface{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML from %q: %w", filePath, err)
	}

	// Process environment variables
	processEnvVars(config)

	return confmap.NewRetrieved(config)
}

// Scheme implements confmap.Provider
func (p *YAMLProvider) Scheme() string {
	return "yaml"
}

// Shutdown implements confmap.Provider
func (p *YAMLProvider) Shutdown(context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.watchers = nil
	p.files = nil
	return nil
}

// processEnvVars recursively processes environment variables in the configuration
func processEnvVars(config map[string]interface{}) {
	for key, value := range config {
		switch v := value.(type) {
		case string:
			config[key] = os.ExpandEnv(v)
		case map[string]interface{}:
			processEnvVars(v)
		case []interface{}:
			for i, item := range v {
				if m, ok := item.(map[string]interface{}); ok {
					processEnvVars(m)
				} else if s, ok := item.(string); ok {
					v[i] = os.ExpandEnv(s)
				}
			}
		}
	}
}
