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

// Load loads the configuration from file and environment
func Load(logger *zap.Logger) (*Config, error) {
	cfg := DefaultConfig()

	// Get config file path from environment or use default
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	configFile := filepath.Join(homeDir, ".srediag", "config", "srediag.yaml")
	if envConfig := os.Getenv("SREDIAG_CONFIG"); envConfig != "" {
		configFile = envConfig
	}

	// Create default config directory if it doesn't exist
	configDir := filepath.Dir(configFile)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		logger.Warn("Failed to create config directory", zap.String("path", configDir), zap.Error(err))
	}

	// Read config file if it exists
	if _, err := os.Stat(configFile); err == nil {
		data, err := os.ReadFile(configFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	// Override with environment variables
	if envPluginDir := os.Getenv("SREDIAG_PLUGIN_DIR"); envPluginDir != "" {
		cfg.PluginsDir = envPluginDir
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
