// Package diagnostic provides diagnostic configuration types
package diagnostic

// Config represents the diagnostic configuration
type Config struct {
	System     *SystemConfig     `mapstructure:"system"`
	Kubernetes *KubernetesConfig `mapstructure:"kubernetes"`
	Cloud      *CloudConfig      `mapstructure:"cloud"`
	Security   *SecurityConfig   `mapstructure:"security"`
}

// SystemConfig represents system diagnostic configuration
type SystemConfig struct {
	Enabled     bool    `mapstructure:"enabled"`
	Interval    string  `mapstructure:"interval"`
	CPULimit    float64 `mapstructure:"cpu_limit"`
	MemoryLimit float64 `mapstructure:"memory_limit"`
	DiskLimit   float64 `mapstructure:"disk_limit"`
}

// KubernetesConfig represents Kubernetes diagnostic configuration
type KubernetesConfig struct {
	Enabled   bool     `mapstructure:"enabled"`
	Clusters  []string `mapstructure:"clusters"`
	Namespace string   `mapstructure:"namespace"`
}

// CloudConfig represents cloud diagnostic configuration
type CloudConfig struct {
	Enabled     bool              `mapstructure:"enabled"`
	Providers   []string          `mapstructure:"providers"`
	Credentials map[string]string `mapstructure:"credentials"`
}

// SecurityConfig represents security diagnostic configuration
type SecurityConfig struct {
	Enabled         bool     `mapstructure:"enabled"`
	ScanInterval    string   `mapstructure:"scan_interval"`
	Standards       []string `mapstructure:"standards"`
	ComplianceLevel string   `mapstructure:"compliance_level"`
}
