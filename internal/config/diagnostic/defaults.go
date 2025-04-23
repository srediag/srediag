package diagnostic

// DefaultConfig returns the default diagnostic configuration
func DefaultConfig() Config {
	return Config{
		System: &SystemConfig{
			Enabled:     true,
			Interval:    "30s",
			CPULimit:    80.0,
			MemoryLimit: 90.0,
			DiskLimit:   85.0,
		},
		Kubernetes: &KubernetesConfig{
			Enabled:   false,
			Clusters:  []string{},
			Namespace: "default",
		},
		Cloud: &CloudConfig{
			Enabled:     false,
			Providers:   []string{},
			Credentials: make(map[string]string),
		},
		Security: &SecurityConfig{
			Enabled:         true,
			ScanInterval:    "1h",
			Standards:       []string{"pci-dss", "hipaa"},
			ComplianceLevel: "high",
		},
	}
}
