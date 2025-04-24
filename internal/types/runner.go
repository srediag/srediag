package types

// IRunner represents the main application runner interface
type IRunner interface {
	IComponent
	// GetConfig returns the configuration
	GetConfig() IConfig
	// GetPluginManager returns the plugin manager
	GetPluginManager() IPluginManager
	// GetResourceMonitor returns the resource monitor
	GetResourceMonitor() IResourceMonitor
	// GetConfigManager returns the config manager
	GetConfigManager() IConfigManager
	// GetTelemetryBridge returns the telemetry bridge
	GetTelemetryBridge() ITelemetryBridge
	// GetDiagnosticManager returns the diagnostic manager
	GetDiagnosticManager() IDiagnosticManager
}

// RunnerStatus represents the current state of the runner
type RunnerStatus int

const (
	// RunnerStatusUnknown indicates the runner status is not known
	RunnerStatusUnknown RunnerStatus = iota
	// RunnerStatusStarting indicates the runner is starting up
	RunnerStatusStarting
	// RunnerStatusRunning indicates the runner is running normally
	RunnerStatusRunning
	// RunnerStatusStopping indicates the runner is shutting down
	RunnerStatusStopping
	// RunnerStatusStopped indicates the runner has stopped
	RunnerStatusStopped
	// RunnerStatusError indicates the runner encountered an error
	RunnerStatusError
)

// RunnerConfig represents the runner configuration
type RunnerConfig struct {
	// LogLevel defines the logging level for the runner
	LogLevel string `json:"log_level" mapstructure:"log_level"`
	// LogFormat defines the logging format (json, text)
	LogFormat string `json:"log_format" mapstructure:"log_format"`
	// ConfigPath defines the path to the configuration file
	ConfigPath string `json:"config_path" mapstructure:"config_path"`
	// PluginsPath defines the path to the plugins directory
	PluginsPath string `json:"plugins_path" mapstructure:"plugins_path"`
	// TelemetryEnabled enables or disables telemetry collection
	TelemetryEnabled bool `json:"telemetry_enabled" mapstructure:"telemetry_enabled"`
	// DiagnosticsEnabled enables or disables diagnostics collection
	DiagnosticsEnabled bool `json:"diagnostics_enabled" mapstructure:"diagnostics_enabled"`
}

// RunnerInfo represents information about the runner
type RunnerInfo struct {
	// Version is the version of the runner
	Version string `json:"version"`
	// Status is the current status of the runner
	Status RunnerStatus `json:"status"`
	// Error contains any error message if the runner is in error state
	Error string `json:"error,omitempty"`
	// Config contains the current runner configuration
	Config RunnerConfig `json:"config"`
	// StartTime is when the runner was started
	StartTime string `json:"start_time"`
	// Uptime is how long the runner has been running
	Uptime string `json:"uptime"`
}
