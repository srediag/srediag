package config

import (
	"context"

	"go.opentelemetry.io/collector/confmap"
)

// ConfigProvider defines the interface for configuration providers
type ConfigProvider interface {
	// Retrieve retrieves the configuration from the provider.
	// The uri parameter specifies the location of the configuration.
	// The watcher parameter is optional and can be nil.
	Retrieve(ctx context.Context, uri string, watcher confmap.WatcherFunc) (*confmap.Conf, error)

	// Scheme returns the scheme that this provider supports (e.g., "file", "env", "yaml")
	Scheme() string

	// Shutdown signals that the provider should close all connections and release resources
	Shutdown(ctx context.Context) error
}

// ConfigVersion represents a configuration schema version
type ConfigVersion struct {
	Major int
	Minor int
	Patch int
}

// ConfigValidator defines the interface for configuration validators
type ConfigValidator interface {
	// Validate validates the configuration and returns an error if invalid
	Validate(conf *confmap.Conf) error

	// Version returns the configuration schema version this validator supports
	Version() ConfigVersion
}

// ConfigWatcher defines the interface for configuration change watchers
type ConfigWatcher interface {
	// OnConfigChange is called when the configuration changes
	OnConfigChange(conf *confmap.Conf) error
}
