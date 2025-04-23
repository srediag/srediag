// Package types provides configuration types for SREDIAG
package types

import (
	"github.com/srediag/srediag/internal/core"
)

// CoreConfig represents the core configuration
type CoreConfig struct {
	LogLevel  string         `mapstructure:"log_level"`
	LogFormat string         `mapstructure:"log_format"`
	Security  SecurityConfig `mapstructure:"security"`
}

// SecurityConfig represents the security configuration
type SecurityConfig struct {
	TLS TLSConfig `mapstructure:"tls"`
}

// TLSConfig represents the TLS configuration
type TLSConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	CertFile string `mapstructure:"cert_file"`
	KeyFile  string `mapstructure:"key_file"`
	CAFile   string `mapstructure:"ca_file,omitempty"`
}

// ServiceConfig represents the service configuration
type ServiceConfig struct {
	Name        string    `mapstructure:"name"`
	Version     string    `mapstructure:"version"`
	Environment string    `mapstructure:"environment"`
	Type        core.Type `mapstructure:"type"`
}

// Ensure ServiceConfig implements ISREDiagServiceConfig
var _ core.ISREDiagServiceConfig = (*ServiceConfig)(nil)

// GetName implements ISREDiagServiceConfig
func (s *ServiceConfig) GetName() string {
	return s.Name
}

// GetEnvironment implements ISREDiagServiceConfig
func (s *ServiceConfig) GetEnvironment() string {
	return s.Environment
}

// GetType implements ISREDiagServiceConfig
func (s *ServiceConfig) GetType() core.Type {
	return s.Type
}

// PluginsConfig represents the plugins configuration
type PluginsConfig struct {
	Directory string                 `mapstructure:"directory"`
	AutoLoad  bool                   `mapstructure:"autoload"`
	Enabled   []string               `mapstructure:"enabled"`
	Settings  map[string]interface{} `mapstructure:"settings"`
}

// Ensure PluginsConfig implements ISREDiagPluginConfig
var _ core.ISREDiagPluginConfig = (*PluginsConfig)(nil)

// IsEnabled implements ISREDiagPluginConfig
func (p *PluginsConfig) IsEnabled() bool {
	return p.AutoLoad
}

// GetName implements ISREDiagPluginConfig
func (p *PluginsConfig) GetName() string {
	return "plugins"
}

// GetSettings implements ISREDiagPluginConfig
func (p *PluginsConfig) GetSettings() map[string]interface{} {
	return p.Settings
}

// GetType implements ISREDiagPluginConfig
func (p *PluginsConfig) GetType() core.Type {
	return core.TypeManagement
}
