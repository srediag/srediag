package builder

import (
	"fmt"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

// LegacyConfig represents the current otelcol-builder.yaml format
type LegacyConfig struct {
	Dist struct {
		Module           string `yaml:"module"`
		Name             string `yaml:"name"`
		Description      string `yaml:"description"`
		OutputPath       string `yaml:"output_path"`
		Version          string `yaml:"version"`
		DebugCompilation bool   `yaml:"debug_compilation"`
	} `yaml:"dist"`
	Receivers  []ModuleConfig `yaml:"receivers"`
	Processors []ModuleConfig `yaml:"processors"`
	Exporters  []ModuleConfig `yaml:"exporters"`
	Extensions []ModuleConfig `yaml:"extensions"`
	Providers  []ModuleConfig `yaml:"providers"`
}

// ModuleConfig represents a module in the legacy config
type ModuleConfig struct {
	GoMod string `yaml:"gomod"`
}

// adaptLegacyConfig converts legacy config to new format
func adaptLegacyConfig(configPath string) (*BuilderConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var legacy LegacyConfig
	if err := yaml.Unmarshal(data, &legacy); err != nil {
		return nil, fmt.Errorf("failed to parse legacy config: %w", err)
	}

	config := &BuilderConfig{
		Components: make(map[ComponentType]map[string]PluginConfig),
	}

	// Copy dist section
	config.Dist.Name = legacy.Dist.Name
	config.Dist.OutputPath = legacy.Dist.OutputPath
	config.Dist.OtelColVersion = legacy.Dist.Version

	// Initialize component maps
	config.Components[TypeReceiver] = make(map[string]PluginConfig)
	config.Components[TypeProcessor] = make(map[string]PluginConfig)
	config.Components[TypeExporter] = make(map[string]PluginConfig)
	config.Components[TypeExtension] = make(map[string]PluginConfig)

	// Convert receivers
	for _, r := range legacy.Receivers {
		name := extractComponentName(r.GoMod)
		config.Components[TypeReceiver][name] = PluginConfig{
			GoMod:  r.GoMod,
			Import: extractImportPath(r.GoMod),
		}
	}

	// Convert processors
	for _, p := range legacy.Processors {
		name := extractComponentName(p.GoMod)
		config.Components[TypeProcessor][name] = PluginConfig{
			GoMod:  p.GoMod,
			Import: extractImportPath(p.GoMod),
		}
	}

	// Convert exporters
	for _, e := range legacy.Exporters {
		name := extractComponentName(e.GoMod)
		config.Components[TypeExporter][name] = PluginConfig{
			GoMod:  e.GoMod,
			Import: extractImportPath(e.GoMod),
		}
	}

	// Convert extensions
	for _, e := range legacy.Extensions {
		name := extractComponentName(e.GoMod)
		config.Components[TypeExtension][name] = PluginConfig{
			GoMod:  e.GoMod,
			Import: extractImportPath(e.GoMod),
		}
	}

	return config, nil
}

// extractComponentName extracts component name from gomod path
func extractComponentName(gomod string) string {
	parts := strings.Split(gomod, " ")
	if len(parts) == 0 {
		return ""
	}
	base := path.Base(parts[0])
	return strings.TrimSuffix(base, "receiver")
}

// extractImportPath extracts import path from gomod
func extractImportPath(gomod string) string {
	parts := strings.Split(gomod, " ")
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}
