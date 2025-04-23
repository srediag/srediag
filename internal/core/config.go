package core

// ISREDiagConfig represents the main configuration interface
type ISREDiagConfig interface {
	// GetVersion returns the service version
	GetVersion() string
	// GetServiceConfig returns the service configuration
	GetServiceConfig() ISREDiagServiceConfig
	// GetCollectorConfig returns the collector configuration
	GetCollectorConfig() ISREDiagCollectorConfig
	// Validate validates the configuration
	Validate() error
}

// ISREDiagServiceConfig represents the service configuration interface
type ISREDiagServiceConfig interface {
	// GetName returns the service name
	GetName() string
	// GetEnvironment returns the service environment
	GetEnvironment() string
	// GetType returns the service type
	GetType() Type
}

// ISREDiagCollectorConfig represents the collector configuration interface
type ISREDiagCollectorConfig interface {
	// IsEnabled returns whether the collector is enabled
	IsEnabled() bool
	// GetConfigPath returns the collector configuration path
	GetConfigPath() string
}

// ISREDiagPluginConfig represents the plugin configuration interface
type ISREDiagPluginConfig interface {
	// GetName returns the plugin name
	GetName() string
	// GetType returns the plugin type
	GetType() Type
	// IsEnabled returns whether the plugin is enabled
	IsEnabled() bool
	// GetSettings returns the plugin settings
	GetSettings() map[string]interface{}
}

// ISREDiagDiagnosticConfig represents the diagnostic configuration interface
type ISREDiagDiagnosticConfig interface {
	// IsEnabled returns whether diagnostics are enabled
	IsEnabled() bool
	// GetType returns the diagnostic type
	GetType() Type
	// GetSettings returns the diagnostic settings
	GetSettings() map[string]interface{}
	// GetInterval returns the diagnostic interval
	GetInterval() string
	// GetThresholds returns the diagnostic thresholds
	GetThresholds() map[string]float64
}
