// internal/config/config.go
package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"

	"github.com/srediag/srediag/internal/core"
)

// Config represents the main SREDIAG configuration
type Config struct {
	Version   string          `mapstructure:"version"`
	Debug     bool            `mapstructure:"debug"`
	LogLevel  string          `mapstructure:"log_level"`
	Service   ServiceConfig   `mapstructure:"service"`
	Telemetry TelemetryConfig `mapstructure:"telemetry"`
	Plugins   PluginsConfig   `mapstructure:"plugins"`
	Security  SecurityConfig  `mapstructure:"security"`
}

// ServiceConfig represents service-level configuration
type ServiceConfig struct {
	Name        string `mapstructure:"name"`
	Environment string `mapstructure:"environment"`
}

// TelemetryConfig represents OpenTelemetry configuration
type TelemetryConfig struct {
	Enabled            bool              `mapstructure:"enabled"`
	ServiceName        string            `mapstructure:"service_name"`
	Endpoint           string            `mapstructure:"endpoint"`
	Protocol           string            `mapstructure:"protocol"`
	Environment        string            `mapstructure:"environment"`
	ResourceAttributes map[string]string `mapstructure:"resource_attributes"`
	Sampling           SamplingConfig    `mapstructure:"sampling"`
}

// SamplingConfig represents trace sampling configuration
type SamplingConfig struct {
	Type string  `mapstructure:"type"`
	Rate float64 `mapstructure:"rate"`
}

// PluginsConfig represents plugin configuration
type PluginsConfig struct {
	Directory string                            `mapstructure:"directory"`
	Enabled   []string                          `mapstructure:"enabled"`
	Settings  map[string]map[string]interface{} `mapstructure:"settings"`
}

// SecurityConfig represents security configuration
type SecurityConfig struct {
	TLS  TLSConfig  `mapstructure:"tls"`
	Auth AuthConfig `mapstructure:"auth"`
}

// TLSConfig represents TLS configuration
type TLSConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	CertFile string `mapstructure:"cert_file"`
	KeyFile  string `mapstructure:"key_file"`
	CAFile   string `mapstructure:"ca_file"`
}

// AuthConfig represents authentication configuration
type AuthConfig struct {
	Type      string          `mapstructure:"type"`
	TokenFile string          `mapstructure:"token_file"`
	Basic     BasicAuthConfig `mapstructure:"basic"`
	OAuth     OAuthConfig     `mapstructure:"oauth"`
}

// BasicAuthConfig represents basic authentication configuration
type BasicAuthConfig struct {
	Username     string `mapstructure:"username"`
	PasswordFile string `mapstructure:"password_file"`
}

// OAuthConfig represents OAuth configuration
type OAuthConfig struct {
	ClientID         string   `mapstructure:"client_id"`
	ClientSecretFile string   `mapstructure:"client_secret_file"`
	TokenURL         string   `mapstructure:"token_url"`
	Scopes           []string `mapstructure:"scopes"`
}

// Load loads the configuration from file and environment variables
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// Set default values
	setDefaults(v)

	// Configure viper
	v.SetConfigFile(configPath)
	v.SetEnvPrefix("SREDIAG")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	// Unmarshal config
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}

// setDefaults sets default values for configuration
func setDefaults(v *viper.Viper) {
	v.SetDefault("version", "v0.1.0")
	v.SetDefault("debug", false)
	v.SetDefault("log_level", "info")

	v.SetDefault("service.name", "srediag")
	v.SetDefault("service.environment", "production")

	v.SetDefault("telemetry.enabled", true)
	v.SetDefault("telemetry.service_name", "srediag")
	v.SetDefault("telemetry.endpoint", "http://localhost:4317")
	v.SetDefault("telemetry.protocol", "grpc")
	v.SetDefault("telemetry.environment", "production")
	v.SetDefault("telemetry.sampling.type", core.SamplingTypeProbabilistic)
	v.SetDefault("telemetry.sampling.rate", core.DefaultSamplingRate)

	v.SetDefault("plugins.directory", "plugins")

	v.SetDefault("security.tls.enabled", false)
	v.SetDefault("security.tls.cert_file", "/etc/srediag/certs/server.crt")
	v.SetDefault("security.tls.key_file", "/etc/srediag/certs/server.key")
	v.SetDefault("security.tls.ca_file", "/etc/srediag/certs/ca.crt")
	v.SetDefault("security.auth.type", "none")
}
