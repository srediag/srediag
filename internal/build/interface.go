package build

import "github.com/srediag/srediag/internal/core"

// Builder defines the interface for building OpenTelemetry plugins
type IBuilder interface {
	// BuildAll builds all plugins defined in the configuration
	BuildAll() error

	// BuildPlugin builds a single plugin
	BuildPlugin(name string, cfg PluginConfig, compType core.ComponentType) error
}

// MakeBuilderInterface extends Builder with make integration
type MakeBuilderInterface interface {
	IBuilder

	// InstallPlugins installs built plugins to the system
	InstallPlugins() error
}
