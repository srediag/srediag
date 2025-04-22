package core

const (
	// DefaultMetricInterval is the default interval for metric collection
	DefaultMetricInterval = 10 // seconds

	// DefaultSamplingRate is the default sampling rate for traces
	DefaultSamplingRate = 1.0

	// SamplingTypeAlwaysOn indicates that all traces should be collected
	SamplingTypeAlwaysOn = "always_on"

	// SamplingTypeAlwaysOff indicates that no traces should be collected
	SamplingTypeAlwaysOff = "always_off"

	// SamplingTypeProbabilistic indicates that traces should be collected probabilistically
	SamplingTypeProbabilistic = "probabilistic"

	// DefaultPluginPattern is the pattern for plugin files
	DefaultPluginPattern = "*.so"

	// DefaultPluginSymbol is the symbol that must be exported by plugins
	DefaultPluginSymbol = "Plugin"
)
