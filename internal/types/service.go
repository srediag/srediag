package types

// ServiceSettings represents service-level settings
type ServiceSettings struct {
	Name        string            `mapstructure:"name" json:"name"`
	Version     string            `mapstructure:"version" json:"version"`
	Environment string            `mapstructure:"environment" json:"environment"`
	Type        ComponentType     `mapstructure:"type" json:"type"`
	Security    SecurityConfig    `mapstructure:"security" json:"security"`
	Settings    map[string]string `mapstructure:"settings" json:"settings"`
}

// Ensure ServiceSettings implements IServiceConfig
var _ IServiceConfig = (*ServiceSettings)(nil)

// GetName returns the name of the service
func (s *ServiceSettings) GetName() string {
	return s.Name
}

// GetEnvironment returns the environment of the service
func (s *ServiceSettings) GetEnvironment() string {
	return s.Environment
}

// GetType returns the type of the service
func (s *ServiceSettings) GetType() ComponentType {
	return s.Type
}

// GetSecurity returns the security configuration
func (s *ServiceSettings) GetSecurity() SecurityConfig {
	return s.Security
}

// GetVersion returns the service version
func (s *ServiceSettings) GetVersion() string {
	return s.Version
}
