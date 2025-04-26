package builder

// Builder defines the interface for building OpenTelemetry plugins
type Builder interface {
	// BuildAll builds all plugins defined in the configuration
	BuildAll() error

	// BuildPlugin builds a single plugin
	BuildPlugin(name string, cfg PluginConfig, compType ComponentType) error
}

// MakeBuilderInterface extends Builder with make integration
type MakeBuilderInterface interface {
	Builder

	// installPlugins installs built plugins to the system
	installPlugins() error
}
