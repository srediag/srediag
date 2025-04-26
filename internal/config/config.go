package config

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	LogLevel   string `yaml:"log_level"`
	LogFormat  string `yaml:"log_format"`
	PluginsDir string `yaml:"plugins_dir"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}

	return &Config{
		LogLevel:   "info",
		LogFormat:  "console",
		PluginsDir: filepath.Join(homeDir, ".srediag", "plugins"),
	}
}

// tryLoadConfigFile attempts to load config from a specific file
func tryLoadConfigFile(logger *zap.Logger, cfg *Config, path string) bool {
	if _, err := os.Stat(path); err == nil {
		data, err := os.ReadFile(path)
		if err != nil {
			logger.Debug("Failed to read config file",
				zap.String("path", path),
				zap.Error(err))
			return false
		}

		if err := yaml.Unmarshal(data, cfg); err != nil {
			logger.Debug("Failed to parse config file",
				zap.String("path", path),
				zap.Error(err))
			return false
		}

		logger.Info("Loaded configuration from file", zap.String("path", path))
		return true
	}
	return false
}

// Load loads the configuration from file and environment
func Load(logger *zap.Logger) (*Config, error) {
	cfg := DefaultConfig()

	// Try loading from environment variable first
	if envConfig := os.Getenv("SREDIAG_CONFIG"); envConfig != "" {
		if !tryLoadConfigFile(logger, cfg, envConfig) {
			return nil, fmt.Errorf("failed to load config from SREDIAG_CONFIG=%s", envConfig)
		}
	} else {
		// Try loading from default locations
		configLocations := []string{
			filepath.Join("configs", "srediag.yaml"), // Project configs directory
			"srediag.yaml",                           // Current directory
			".srediag.yaml",                          // Hidden file in current directory
		}

		// Add home directory config
		if homeDir, err := os.UserHomeDir(); err == nil {
			configLocations = append(configLocations,
				filepath.Join(homeDir, ".srediag", "config", "srediag.yaml"),
				filepath.Join(homeDir, ".srediag.yaml"),
			)
		}

		// Try each location
		configFound := false
		for _, loc := range configLocations {
			if tryLoadConfigFile(logger, cfg, loc) {
				configFound = true
				break
			}
		}

		if !configFound {
			logger.Info("No config file found in default locations, using defaults")
		}
	}

	// Override with environment variables
	if envPluginDir := os.Getenv("SREDIAG_PLUGIN_DIR"); envPluginDir != "" {
		cfg.PluginsDir = envPluginDir
		logger.Info("Using plugin directory from environment",
			zap.String("SREDIAG_PLUGIN_DIR", envPluginDir))
	}

	// Validate and normalize paths
	if cfg.PluginsDir != "" {
		absPath, err := filepath.Abs(cfg.PluginsDir)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve plugin path: %w", err)
		}
		cfg.PluginsDir = absPath

		// Create plugin directory if it doesn't exist
		if err := os.MkdirAll(cfg.PluginsDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create plugin directory: %w", err)
		}
	}

	return cfg, nil
}

// ValidateConfig validates the configuration
func ValidateConfig(cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("config cannot be nil")
	}

	// Validate plugin path
	if cfg.PluginsDir == "" {
		return fmt.Errorf("plugin path cannot be empty")
	}

	return nil
}
