package main

import (
	"context"
	"fmt"

	"github.com/srediag/srediag/internal/plugins"
)

// MockPlugin is a test plugin
type MockPlugin struct {
	info    plugins.Info
	config  map[string]interface{}
	running bool
}

// New creates a new instance of the mock plugin
func New() plugins.Plugin {
	return &MockPlugin{
		info: plugins.Info{
			Name:        "mock-plugin",
			Version:     "1.0.0",
			Type:        "test",
			Description: "Test plugin for unit testing",
			Author:      "SREDIAG Team",
		},
	}
}

// Init initializes the plugin with the provided configuration
func (p *MockPlugin) Init(config map[string]interface{}) error {
	if config == nil {
		return fmt.Errorf("configuration cannot be nil")
	}
	p.config = config
	return nil
}

// Start starts the plugin
func (p *MockPlugin) Start(ctx context.Context) error {
	if p.running {
		return fmt.Errorf("plugin is already running")
	}
	p.running = true
	return nil
}

// Stop stops the plugin
func (p *MockPlugin) Stop(ctx context.Context) error {
	if !p.running {
		return fmt.Errorf("plugin is not running")
	}
	p.running = false
	return nil
}

// Info returns the plugin metadata
func (p *MockPlugin) Info() plugins.Info {
	return p.info
}
