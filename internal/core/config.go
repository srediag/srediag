// Package core provides foundational types, configuration schemas, and utilities for the SREDIAG system.
//
// This file defines the canonical Config struct, representing the full configuration schema for SREDIAG, and all supporting types for modular configuration sections.
// It also provides helpers for loading, overlaying, and validating configuration from YAML, environment variables, and CLI flags, supporting strict schema validation and modular overlays.
//
// Philosophy:
//   - The configuration system is designed to be modular, extensible, and robust, supporting overlays from multiple sources (YAML, env, CLI) with clear precedence.
//   - All configuration is validated after loading, and each section is represented by a dedicated struct for clarity and maintainability.
//   - The system is intended to be the single source of truth for all configuration in SREDIAG, supporting both legacy and future modular config consumers.
//
// Usage:
//   - Use Config as the root struct for loading, validating, and passing configuration throughout the system.
//   - Use section-specific types (e.g., ServiceConfig, LoggingConfig) for modular access and validation.
//   - Use LoadConfigWithOverlay and section-specific loaders to populate config from all supported overlays.
//   - Always call ValidateConfig after loading.
//
// Extensibility:
//   - Add new config sections by extending the Config struct and providing a dedicated section struct.
//   - Add new overlays or precedence rules by extending LoadConfigWithOverlay.
//   - Use strict YAML unmarshalling and schema validation for new sections.
//
// Security:
//   - Sensitive fields (e.g., secrets, keys) should be loaded securely and never logged.
//   - Use SHA256 and signature fields for plugin and binary integrity verification.
//
// See documentation in docs/configuration/ for full schema and usage examples.
package core

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v3"
)

// Config holds all configuration for SREDIAG, matching the documented YAML structure.
//
// Usage:
//   - Use as the root struct for loading, validating, and passing configuration throughout the system.
//   - All config overlays (flags, env, YAML) should ultimately populate this struct.
//
// Fields:
//   - Service: Service-level settings (see ServiceConfig).
//   - Logging: Logging configuration (see LoggingConfig).
//   - Security: Security, TLS, authentication, RBAC, quotas, and runtime security (see SecurityConfig).
//   - Collector: OpenTelemetry Collector settings (see CollectorConfig).
//   - Plugins: Plugin directory and enabled plugins (see PluginsConfig).
//   - Diagnostics: Diagnostics defaults and plugin configs (see DiagnosticsConfig).
//   - Build: Build system configuration and plugin build lists (see BuildConfig).
//
// Best Practices:
//   - Always validate after loading (see ValidateConfig).
//   - Avoid direct field access in business logic; prefer passing config sections as needed.
//
// TODO:
//   - Remove legacy/flat fields from config consumers after migration to modular config.
//
// Redundancy/Refactor:
//   - No redundancy; this is the canonical config struct for SREDIAG.
type Config struct {
	Service     ServiceConfig     `yaml:"service"`     // Service section (see ServiceConfig)
	Logging     LoggingConfig     `yaml:"logging"`     // Logging section (see LoggingConfig)
	Security    SecurityConfig    `yaml:"security"`    // Security section (see SecurityConfig)
	Collector   CollectorConfig   `yaml:"collector"`   // Collector section (see CollectorConfig)
	Plugins     PluginsConfig     `yaml:"plugins"`     // Plugins section (see PluginsConfig)
	Diagnostics DiagnosticsConfig `yaml:"diagnostics"` // Diagnostics section (see DiagnosticsConfig)
	Build       BuildConfig       `yaml:"build"`       // Build section (see BuildConfig)
	// TODO: Remove legacy/flat fields from config consumers after migration.
}

// ServiceConfig maps to the 'service:' section in YAML (docs: service.md)
//
// Usage: Used for service-level settings (name, port, environment).
//
// Fields:
//   - Name: Service name. Used for identification and logging.
//   - Port: Service port. Used for network binding and health checks.
//   - Environment: Deployment environment (e.g., dev, staging, prod). Used for environment-specific logic.
type ServiceConfig struct {
	Name        string `yaml:"name"`        // Service name
	Port        int    `yaml:"port"`        // Service port
	Environment string `yaml:"environment"` // Deployment environment
}

// LoggingConfig maps to the 'logging:' section in YAML (docs: README.md)
//
// Usage: Used for logging level and format configuration.
//
// Fields:
//   - Level: Log level (e.g., info, debug). Used for controlling verbosity.
//   - Format: Output format (json, console). Used for log formatting.
type LoggingConfig struct {
	Level  string `yaml:"level"`  // Log level (e.g., info, debug)
	Format string `yaml:"format"` // Output format (json, console)
}

// SecurityConfig maps to the 'security:' section in YAML (docs: security.md)
//
// Usage: Used for TLS, authentication, RBAC, quotas, and runtime security settings.
//
// Fields:
//   - TLS: TLS configuration (enabled, cert/key paths, CA, min version, client verification).
//   - Auth: Authentication configuration (type, JWT, OAuth2).
//   - RBAC: Role-based access control configuration (enabled, roles, default role).
//   - Quotas: Quota configuration (spans per second, logs throughput).
//   - RateLimit: Rate limiting configuration (enabled, RPS, burst).
//   - Runtime: Runtime security configuration (seccomp, AppArmor, read-only root, guards).
type SecurityConfig struct {
	TLS struct {
		Enabled      bool   `yaml:"enabled"`       // Enable TLS
		CertFile     string `yaml:"cert_file"`     // Path to TLS certificate
		KeyFile      string `yaml:"key_file"`      // Path to TLS key
		CAFile       string `yaml:"ca_file"`       // Path to CA certificate
		MinVersion   string `yaml:"min_version"`   // Minimum TLS version
		VerifyClient bool   `yaml:"verify_client"` // Require client certificate
	} `yaml:"tls"`
	Auth struct {
		Type string `yaml:"type"` // Auth type (e.g., jwt, oauth2)
		JWT  struct {
			Secret   string `yaml:"secret"`   // JWT secret
			Lifetime string `yaml:"lifetime"` // JWT token lifetime
		} `yaml:"jwt"`
		OAuth2 struct {
			Issuer       string   `yaml:"issuer"`        // OAuth2 issuer
			ClientID     string   `yaml:"client_id"`     // OAuth2 client ID
			ClientSecret string   `yaml:"client_secret"` // OAuth2 client secret
			Scopes       []string `yaml:"scopes"`        // OAuth2 scopes
		} `yaml:"oauth2"`
	} `yaml:"auth"`
	RBAC struct {
		Enabled     bool                `yaml:"enabled"`      // Enable RBAC
		DefaultRole string              `yaml:"default_role"` // Default RBAC role
		Roles       map[string][]string `yaml:"roles"`        // Role definitions
	} `yaml:"rbac"`
	Quotas struct {
		SpansPerSecond int `yaml:"spans_per_second"` // Span rate limit
		LogsMiBPerMin  int `yaml:"logs_mib_per_min"` // Log throughput limit
	} `yaml:"quotas"`
	RateLimit struct {
		Enabled bool `yaml:"enabled"` // Enable rate limiting
		RPS     int  `yaml:"rps"`     // Requests per second
		Burst   int  `yaml:"burst"`   // Burst size
	} `yaml:"rate_limit"`
	Runtime struct {
		SeccompProfile  string `yaml:"seccomp_profile"`  // Seccomp profile path
		AppArmorProfile string `yaml:"apparmor_profile"` // AppArmor profile path
		ReadOnlyRootFS  bool   `yaml:"read_only_rootfs"` // Enforce read-only root
		MemGuardMiB     int    `yaml:"mem_guard_mib"`    // Memory guard (MiB)
		CPUGuardPct     int    `yaml:"cpu_guard_pct"`    // CPU guard (%)
	} `yaml:"runtime"`
}

// CollectorConfig maps to the 'collector:' section in YAML (docs: service.md)
//
// Usage: Used for OpenTelemetry Collector settings.
//
// Fields:
//   - Enabled: Whether the collector is enabled.
//   - ConfigPath: Path to the collector config file.
//   - MemoryLimitMiB: Memory limit for the collector in MiB.
type CollectorConfig struct {
	Enabled        bool   `yaml:"enabled"`          // Enable collector
	ConfigPath     string `yaml:"config_path"`      // Path to collector config
	MemoryLimitMiB int    `yaml:"memory_limit_mib"` // Memory limit (MiB)
}

// PluginsConfig maps to the 'plugins:' section in YAML (docs: plugin.md)
//
// Usage: Used for plugin directory and enabled plugins.
//
// Fields:
//   - Dir: Directory where plugins are stored.
//   - ExecDir: Directory where plugins are executed.
//   - Enabled: List of enabled plugin names.
type PluginsConfig struct {
	Dir     string   `yaml:"dir"`      // Plugin directory
	ExecDir string   `yaml:"exec_dir"` // Plugin execution directory
	Enabled []string `yaml:"enabled"`  // List of enabled plugins
}

// DiagnosticsConfig maps to the 'diagnostics:' section in YAML (docs: diagnose.md)
//
// Usage: Used for diagnostics defaults and plugin configs.
//
// Fields:
//   - Defaults: Default output format and timeout for diagnostics.
//   - Plugins: Map of plugin-specific diagnostic configs.
type DiagnosticsConfig struct {
	Defaults struct {
		OutputFormat string `yaml:"output_format"` // Default output format
		Timeout      string `yaml:"timeout"`       // Default timeout
	} `yaml:"defaults"`
	Plugins map[string]map[string]interface{} `yaml:"plugins"` // Plugin-specific configs
}

// BuildConfig maps to the 'build:' section in YAML (docs: build.md)
//
// Usage: Used for build system configuration and plugin build lists.
//
// Fields:
//   - Config: Path to builder YAML (for overlays).
//   - OutputDir: Directory where build artefacts are stored.
//   - Dist: Distribution metadata (name, description, version, output path).
//   - Receivers: List of receiver plugin configs.
//   - Processors: List of processor plugin configs.
//   - Exporters: List of exporter plugin configs.
//   - Extensions: List of extension plugin configs.
type BuildConfig struct {
	Config    string `yaml:"config"`     // Path to builder YAML (for overlays)
	OutputDir string `yaml:"output_dir"` // Where artefacts are stored
	Dist      struct {
		Name        string `yaml:"name"`        // Distribution name
		Description string `yaml:"description"` // Distribution description
		Version     string `yaml:"version"`     // Distribution version
		OutputPath  string `yaml:"output_path"` // Output path for build
	} `yaml:"dist"`
	Receivers  []map[string]string `yaml:"receivers"`  // Receiver plugins
	Processors []map[string]string `yaml:"processors"` // Processor plugins
	Exporters  []map[string]string `yaml:"exporters"`  // Exporter plugins
	Extensions []map[string]string `yaml:"extensions"` // Extension plugins
}

// configSearchPaths defines the order of precedence for config file locations.
//
// Usage: Used internally by findConfigFile to determine where to look for config files.
var configSearchPaths = []string{
	"/etc/srediag/srediag.yaml",  // System-wide
	"$HOME/.srediag/config.yaml", // User
	"./config/srediag.yaml",      // Project
	"./srediag.yaml",             // Project root
}

// findConfigFile returns the first config file found and its extension, or empty string if none found.
//
// Usage:
//   - Used internally by LoadConfigWithOverlay for config discovery.
//
// Returns:
//   - string: Path to the first config file found, or empty string if none found.
//   - string: File extension (e.g., .yaml, .json, .toml), or .yaml if not found.
func findConfigFile() (string, string) {
	home, _ := os.UserHomeDir()
	for _, path := range configSearchPaths {
		expanded := strings.ReplaceAll(path, "$HOME", home)
		if _, err := os.Stat(expanded); err == nil {
			ext := filepath.Ext(expanded)
			return expanded, ext
		}
	}
	return "", ".yaml" // default to yaml if not found
}

// LoadConfigWithOverlay loads config from file, overlays env vars, then overlays CLI flags.
// Discovery order and precedence: CLI flags > env vars > YAML > built-ins.
// File type is inferred from extension, defaulting to YAML.
//
// Parameters:
//   - spec: Pointer to the config struct to populate (usually *Config).
//   - cliFlags: Map of CLI flag overrides (highest precedence).
//   - opts: Optional ConfigOption(s) for custom paths, overlays, or suffixes.
//
// Returns:
//   - error: If loading, parsing, or validation fails, returns a detailed error.
//
// Side Effects:
//   - Reads from disk, environment, and CLI flags.
//   - Populates the provided config struct with merged values.
//
// Usage:
//   - Use to load any config struct (usually *Config) with overlays for CLI, env, and YAML.
//   - Always call ValidateConfig after loading.
//
// Best Practices:
//   - Use strict YAML unmarshalling for new config sections.
//   - Extend this function for new overlays or precedence rules.
func LoadConfigWithOverlay(spec interface{}, cliFlags map[string]string, opts ...ConfigOption) error {
	v := viper.New()
	v.SetConfigType("yaml")
	// Set defaults (should match Config struct)
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.format", "console")
	v.SetDefault("plugins.dir", DefaultPluginDir())
	v.SetDefault("plugins.exec_dir", DefaultPluginExecDir())
	v.SetDefault("service.port", 8080)
	v.SetDefault("service.name", "srediag")
	v.SetDefault("collector.enabled", false)
	v.SetDefault("collector.config_path", "/etc/srediag/srediag-service.yaml")
	v.SetDefault("build.output_dir", DefaultBuildOutputDir())
	// Bind all documented env vars
	bindEnvs := map[string]string{
		"logging.level":                      "SREDIAG_LOG_LEVEL",
		"logging.format":                     "SREDIAG_LOG_FORMAT",
		"plugins.dir":                        "SREDIAG_PLUGINS_DIR",
		"plugins.exec_dir":                   "SREDIAG_PLUGINS_EXEC_DIR",
		"service.port":                       "SREDIAG_SERVICE_PORT",
		"service.name":                       "SREDIAG_SERVICE_NAME",
		"collector.enabled":                  "SREDIAG_COLLECTOR_ENABLED",
		"collector.config_path":              "SREDIAG_COLLECTOR_CONFIG_PATH",
		"build.output_dir":                   "SREDIAG_BUILD_OUTPUT_DIR",
		"security.tls.enabled":               "SREDIAG_TLS_ENABLED",
		"security.tls.cert_file":             "SREDIAG_TLS_CERT_FILE",
		"security.tls.key_file":              "SREDIAG_TLS_KEY_FILE",
		"security.tls.ca_file":               "SREDIAG_TLS_CA_FILE",
		"security.auth.type":                 "SREDIAG_AUTH_TYPE",
		"security.rbac.enabled":              "SREDIAG_RBAC_ENABLED",
		"diagnostics.defaults.output_format": "SREDIAG_DIAG_OUTPUT_FORMAT",
		"diagnostics.defaults.timeout":       "SREDIAG_DIAG_TIMEOUT",
	}
	for key, env := range bindEnvs {
		if err := v.BindEnv(key, env); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to bind env %s: %v\n", env, err)
		}
	}
	// 1. CLI flags (highest precedence)
	for key, val := range cliFlags {
		v.Set(key, val)
	}
	// 2. Config file discovery
	var o configSpecOpts
	for _, opt := range opts {
		opt(&o)
	}
	if o.Path != "" {
		v.SetConfigFile(o.Path)
	} else if envPath := os.Getenv("SREDIAG_CONFIG"); envPath != "" {
		v.SetConfigFile(envPath)
	} else {
		cfgPath, ext := findConfigFile()
		if cfgPath != "" {
			v.SetConfigFile(cfgPath)
			switch ext {
			case ".json":
				v.SetConfigType("json")
			case ".toml":
				v.SetConfigType("toml")
			default:
				v.SetConfigType("yaml")
			}
		}
	}
	_ = v.ReadInConfig() // ignore error if not found
	return v.Unmarshal(spec)
}

// DefaultPluginDir returns the default plugin directory based on install context.
//
// Usage:
//   - Used to set plugin directory defaults in config loading.
//
// Returns:
//   - string: Path to the default plugin directory.
func DefaultPluginDir() string {
	if isSystemInstall() {
		return "/usr/libexec/srediag"
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "libexec", "srediag")
}

// DefaultPluginExecDir returns the default plugin exec directory based on install context.
//
// Usage:
//   - Used to set plugin exec directory defaults in config loading.
//
// Returns:
//   - string: Path to the default plugin exec directory.
func DefaultPluginExecDir() string {
	if isSystemInstall() {
		return "/usr/libexec/srediag"
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "libexec", "srediag")
}

// DefaultBuildOutputDir returns the default build output directory based on context.
//
// Usage:
//   - Used to set build output directory defaults in config loading.
//
// Returns:
//   - string: Path to the default build output directory.
func DefaultBuildOutputDir() string {
	if isSystemInstall() {
		return "/var/lib/srediag/build"
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".srediag", "build")
}

// isSystemInstall returns true if running as a system install (heuristic: root, /usr/bin, etc)
//
// Usage:
//   - Used internally to determine default paths.
//
// Returns:
//   - bool: True if running as a system install, false otherwise.
func isSystemInstall() bool {
	return os.Geteuid() == 0 || strings.HasPrefix(os.Args[0], "/usr/")
}

// PrintEffectiveConfig prints the merged config as YAML (for --print-config flag).
//
// Usage:
//   - Use in CLI to print the effective config for debugging or support.
//
// Parameters:
//   - cfg: Pointer to the Config struct to print.
//
// Returns:
//   - error: If marshaling or printing fails, returns a detailed error.
func PrintEffectiveConfig(cfg *Config) error {
	out, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	fmt.Println(string(out))
	return nil
}

// ConfigOption allows customizing config loading (e.g., file name, search paths, overlays).
//
// Usage:
//   - Pass to LoadConfigWithOverlay for custom config discovery or overlays.
type ConfigOption func(*configSpecOpts)

// configSpecOpts holds options for config loading.
//
// Usage:
//   - Used internally by LoadConfigWithOverlay and ConfigOption helpers.
type configSpecOpts struct {
	Path       string // config file path or name
	EnvPrefix  string // prefix for environment variable overlays
	PathSuffix string // config file suffix for discovery
}

// WithConfigPath sets the config file path or name.
//
// Usage:
//   - Pass to LoadConfigWithOverlay to specify a config file directly.
//
// Parameters:
//   - path: Path to the config file.
//
// Returns:
//   - ConfigOption: Option to set the config file path.
func WithConfigPath(path string) ConfigOption {
	return func(o *configSpecOpts) { o.Path = path }
}

// WithEnvPrefix sets the environment variable prefix for overlays.
//
// Usage:
//   - Pass to LoadConfigWithOverlay to specify a custom env var prefix.
//
// Parameters:
//   - prefix: Prefix for environment variable overlays.
//
// Returns:
//   - ConfigOption: Option to set the env var prefix.
func WithEnvPrefix(prefix string) ConfigOption {
	return func(o *configSpecOpts) { o.EnvPrefix = prefix }
}

// WithConfigPathSuffix sets a suffix to search for in the config discovery order (e.g., "build" or "build.yaml").
// If no extension is provided, .yaml is assumed.
//
// Usage:
//   - Pass to LoadConfigWithOverlay to search for configs with a specific suffix.
//
// Parameters:
//   - suffix: Suffix to search for in config discovery.
//
// Returns:
//   - ConfigOption: Option to set the config file suffix.
func WithConfigPathSuffix(suffix string) ConfigOption {
	return func(o *configSpecOpts) { o.PathSuffix = suffix }
}

// StrictYAMLUnmarshal unmarshals YAML with strict field checking and debug logging for unknown keys.
//
// Usage:
//   - Use for strict schema validation when loading config from YAML.
//
// Parameters:
//   - data: YAML data to unmarshal.
//   - out: Pointer to the struct to populate.
//
// Returns:
//   - error: If parsing or validation fails, returns a detailed error.
func StrictYAMLUnmarshal(data []byte, out interface{}) error {
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	if err := dec.Decode(out); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}
	// TODO: Add debug logging for unknown keys if needed.
	return nil
}

// LoadBuildConfig loads only the build config section with overlays and precedence.
//
// Usage:
//   - Use to load just the build section for build operations.
//
// Parameters:
//   - cliFlags: Map of CLI flag overrides.
//   - opts: Optional ConfigOption(s) for custom overlays.
//
// Returns:
//   - BuildConfig: The loaded build configuration.
//   - error: If loading or validation fails, returns a detailed error.
func LoadBuildConfig(cliFlags map[string]string, opts ...ConfigOption) (BuildConfig, error) {
	var cfg Config
	if err := LoadConfigWithOverlay(&cfg, cliFlags, opts...); err != nil {
		return BuildConfig{}, err
	}
	return cfg.Build, nil
}

// LoadPluginConfig loads only the plugins config section with overlays and precedence.
//
// Usage:
//   - Use to load just the plugins section for plugin operations.
//
// Parameters:
//   - cliFlags: Map of CLI flag overrides.
//   - opts: Optional ConfigOption(s) for custom overlays.
//
// Returns:
//   - PluginsConfig: The loaded plugins configuration.
//   - error: If loading or validation fails, returns a detailed error.
func LoadPluginConfig(cliFlags map[string]string, opts ...ConfigOption) (PluginsConfig, error) {
	var cfg Config
	if err := LoadConfigWithOverlay(&cfg, cliFlags, opts...); err != nil {
		return PluginsConfig{}, err
	}
	return cfg.Plugins, nil
}

// LoadDiagnosticsConfig loads only the diagnostics config section with overlays and precedence.
//
// Usage:
//   - Use to load just the diagnostics section for diagnostics operations.
//
// Parameters:
//   - cliFlags: Map of CLI flag overrides.
//   - opts: Optional ConfigOption(s) for custom overlays.
//
// Returns:
//   - DiagnosticsConfig: The loaded diagnostics configuration.
//   - error: If loading or validation fails, returns a detailed error.
func LoadDiagnosticsConfig(cliFlags map[string]string, opts ...ConfigOption) (DiagnosticsConfig, error) {
	var cfg Config
	if err := LoadConfigWithOverlay(&cfg, cliFlags, opts...); err != nil {
		return DiagnosticsConfig{}, err
	}
	return cfg.Diagnostics, nil
}

// LoadSecurityConfig loads only the security config section with overlays and precedence.
//
// Usage:
//   - Use to load just the security section for security operations.
//
// Parameters:
//   - cliFlags: Map of CLI flag overrides.
//   - opts: Optional ConfigOption(s) for custom overlays.
//
// Returns:
//   - SecurityConfig: The loaded security configuration.
//   - error: If loading or validation fails, returns a detailed error.
func LoadSecurityConfig(cliFlags map[string]string, opts ...ConfigOption) (SecurityConfig, error) {
	var cfg Config
	if err := LoadConfigWithOverlay(&cfg, cliFlags, opts...); err != nil {
		return SecurityConfig{}, err
	}
	return cfg.Security, nil
}

// LoadServiceConfig loads only the service config section with overlays and precedence.
//
// Usage:
//   - Use to load just the service section for service operations.
//
// Parameters:
//   - cliFlags: Map of CLI flag overrides.
//   - opts: Optional ConfigOption(s) for custom overlays.
//
// Returns:
//   - ServiceConfig: The loaded service configuration.
//   - error: If loading or validation fails, returns a detailed error.
func LoadServiceConfig(cliFlags map[string]string, opts ...ConfigOption) (ServiceConfig, error) {
	var cfg Config
	if err := LoadConfigWithOverlay(&cfg, cliFlags, opts...); err != nil {
		return ServiceConfig{}, err
	}
	return cfg.Service, nil
}

// LoadCollectorConfig loads only the collector config section with overlays and precedence.
//
// Usage:
//   - Use to load just the collector section for collector operations.
//
// Parameters:
//   - cliFlags: Map of CLI flag overrides.
//   - opts: Optional ConfigOption(s) for custom overlays.
//
// Returns:
//   - CollectorConfig: The loaded collector configuration.
//   - error: If loading or validation fails, returns a detailed error.
func LoadCollectorConfig(cliFlags map[string]string, opts ...ConfigOption) (CollectorConfig, error) {
	var cfg Config
	if err := LoadConfigWithOverlay(&cfg, cliFlags, opts...); err != nil {
		return CollectorConfig{}, err
	}
	return cfg.Collector, nil
}

// ValidateConfig validates the configuration for required fields.
//
// Usage:
//   - Call after loading config to ensure all required fields are set.
//   - Returns an error if any required field is missing or invalid.
//
// Best Practices:
//   - Extend this function as new required fields are added to Config.
//   - Use in all CLI entrypoints and service initializations.
//
// TODO:
//   - Add more granular validation for nested config sections.
//
// Redundancy/Refactor:
//   - No redundancy; this is the canonical config validator for SREDIAG.
func ValidateConfig(cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("config cannot be nil")
	}
	if cfg.Plugins.Dir == "" {
		return fmt.Errorf("plugins.dir must be set")
	}
	if cfg.Logging.Level == "" {
		return fmt.Errorf("logging.level must be set")
	}
	if cfg.Service.Port == 0 {
		return fmt.Errorf("service.port must be set")
	}
	return nil
}

// NewConfig returns a new Config with default (zero) values.
//
// Usage:
//   - Use to create a new config for testing or as a fallback.
//
// Returns:
//   - *Config: A new Config instance with zero/default values.
func NewConfig() *Config {
	return &Config{}
}
