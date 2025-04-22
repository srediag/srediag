package config

import (
	"fmt"

	"github.com/spf13/viper"
	"go.opentelemetry.io/collector/confmap"
)

// OtelCollectorConfig represents the OpenTelemetry Collector configuration
type OtelCollectorConfig struct {
	Receivers  map[string]interface{} `mapstructure:"receivers"`
	Processors map[string]interface{} `mapstructure:"processors"`
	Exporters  map[string]interface{} `mapstructure:"exporters"`
	Service    ServicePipelines       `mapstructure:"service"`
}

// ServicePipelines define the service pipelines of the collector
type ServicePipelines struct {
	Pipelines map[string]Pipeline `mapstructure:"pipelines"`
}

// Pipeline defines the structure of an individual pipeline
type Pipeline struct {
	Receivers  []string `mapstructure:"receivers"`
	Processors []string `mapstructure:"processors"`
	Exporters  []string `mapstructure:"exporters"`
}

// NewDefaultConfig returns a default collector configuration
func NewDefaultConfig() *OtelCollectorConfig {
	return &OtelCollectorConfig{
		Receivers: map[string]interface{}{
			"otlp": map[string]interface{}{
				"protocols": map[string]interface{}{
					"grpc": map[string]interface{}{},
					"http": map[string]interface{}{},
				},
			},
		},
		Processors: map[string]interface{}{
			"batch": map[string]interface{}{},
		},
		Exporters: map[string]interface{}{
			"logging": map[string]interface{}{
				"verbosity": "detailed",
			},
		},
		Service: ServicePipelines{
			Pipelines: map[string]Pipeline{
				"traces": {
					Receivers:  []string{"otlp"},
					Processors: []string{"batch"},
					Exporters:  []string{"logging"},
				},
				"metrics": {
					Receivers:  []string{"otlp"},
					Processors: []string{"batch"},
					Exporters:  []string{"logging"},
				},
			},
		},
	}
}

// ValidateConfig validates the collector configuration
func (c *OtelCollectorConfig) ValidateConfig() error {
	if len(c.Service.Pipelines) == 0 {
		return fmt.Errorf("at least one pipeline must be configured")
	}

	for name, pipeline := range c.Service.Pipelines {
		if len(pipeline.Receivers) == 0 {
			return fmt.Errorf("pipeline %s must have at least one receiver", name)
		}
		if len(pipeline.Exporters) == 0 {
			return fmt.Errorf("pipeline %s must have at least one exporter", name)
		}

		// Validate that all referenced components exist
		for _, recv := range pipeline.Receivers {
			if _, exists := c.Receivers[recv]; !exists {
				return fmt.Errorf("pipeline %s references non-existent receiver %s", name, recv)
			}
		}

		for _, proc := range pipeline.Processors {
			if _, exists := c.Processors[proc]; !exists {
				return fmt.Errorf("pipeline %s references non-existent processor %s", name, proc)
			}
		}

		for _, exp := range pipeline.Exporters {
			if _, exists := c.Exporters[exp]; !exists {
				return fmt.Errorf("pipeline %s references non-existent exporter %s", name, exp)
			}
		}
	}

	return nil
}

// LoadCollectorConfig loads the collector configuration from a YAML file
func LoadCollectorConfig(configPath string) (*OtelCollectorConfig, error) {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading configuration file: %w", err)
	}

	var cfg OtelCollectorConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error decoding configuration: %w", err)
	}

	if err := cfg.ValidateConfig(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// ToConfMap converts the configuration to a confmap.Conf
func (c *OtelCollectorConfig) ToConfMap() (*confmap.Conf, error) {
	data := map[string]interface{}{
		"receivers":  c.Receivers,
		"processors": c.Processors,
		"exporters":  c.Exporters,
		"service":    c.Service,
	}

	return confmap.NewFromStringMap(data), nil
}
