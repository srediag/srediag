package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"

	"github.com/srediag/srediag/internal/types"
)

const (
	// DefaultConfigPath is the default path to the configuration file
	DefaultConfigPath = "/etc/srediag/config/srediag.yaml"
)

// InitializeConfig initializes the configuration system
func InitializeConfig(configPath string) {
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		// Search for config in default locations
		viper.AddConfigPath("/etc/srediag/config")
		viper.AddConfigPath("$HOME/.srediag")
		viper.AddConfigPath(".")
		viper.SetConfigName("srediag")
		viper.SetConfigType("yaml")
	}

	// Set environment variables
	viper.SetEnvPrefix("SREDIAG")
	viper.AutomaticEnv()

	// Bind environment variables
	bindEnvs := []struct {
		key      string
		envVar   string
		required bool
	}{
		{"config", "SREDIAG_CONFIG", false},
		{"format", "SREDIAG_OUTPUT_FORMAT", false},
		{"log_level", "SREDIAG_LOG_LEVEL", false},
		{"api_key", "SREDIAG_API_KEY", false},
	}

	for _, env := range bindEnvs {
		if err := viper.BindEnv(env.key, env.envVar); err != nil {
			if env.required {
				fmt.Fprintf(os.Stderr, "Error binding environment variable %s: %v\n", env.envVar, err)
				os.Exit(1)
			}
			fmt.Fprintf(os.Stderr, "Warning: failed to bind environment variable %s: %v\n", env.envVar, err)
		}
	}
}

// Load loads the configuration from file
func Load(configPath string) (*types.Config, error) {
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var cfg types.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	return &cfg, nil
}
