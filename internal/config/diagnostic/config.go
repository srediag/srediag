package diagnostic

// Config contains diagnostic configurations
type Config struct {
	System     *SystemConfig     `mapstructure:"system"`
	Kubernetes *KubernetesConfig `mapstructure:"kubernetes"`
	Cloud      *CloudConfig      `mapstructure:"cloud"`
	Security   *SecurityConfig   `mapstructure:"security"`
}

// SystemConfig contains system diagnostic configurations
type SystemConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

// KubernetesConfig contains kubernetes diagnostic configurations
type KubernetesConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

// CloudConfig contains cloud diagnostic configurations
type CloudConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

// SecurityConfig contains security diagnostic configurations
type SecurityConfig struct {
	Enabled bool `mapstructure:"enabled"`
}
