package builder

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"go.uber.org/zap"
)

// ComponentType represents the type of OpenTelemetry component
type ComponentType string

// Component types supported by the builder
const (
	TypeReceiver  ComponentType = "receivers"
	TypeProcessor ComponentType = "processors"
	TypeExporter  ComponentType = "exporters"
	TypeExtension ComponentType = "extensions"
	TypeConnector ComponentType = "connectors"
)

// PluginConfig represents a component configuration from otelcol-builder.yaml
type PluginConfig struct {
	GoMod  string   `yaml:"gomod"`          // Go module path and version (e.g. "github.com/org/repo v1.0.0")
	Import string   `yaml:"import"`         // Import path for the component
	Path   string   `yaml:"path"`           // Local path to the component code
	Tags   []string `yaml:"tags,omitempty"` // Build tags for the component
}

// BuilderConfig represents the otelcol-builder.yaml configuration
type BuilderConfig struct {
	Dist struct {
		Name           string `yaml:"name"`            // Name of the distribution
		OutputPath     string `yaml:"output_path"`     // Output path for built artifacts
		OtelColVersion string `yaml:"otelcol_version"` // OpenTelemetry Collector version
	} `yaml:"dist"`
	Components map[ComponentType]map[string]PluginConfig `yaml:"components"` // Map of component type to component configurations
}

// PluginBuilder handles the compilation of OpenTelemetry components as plugins
type PluginBuilder struct {
	logger     *zap.Logger
	configPath string
	outputDir  string
	tempDir    string // Directory for temporary build files
}

// NewPluginBuilder creates a new plugin builder
func NewPluginBuilder(logger *zap.Logger, configPath, outputDir string) *PluginBuilder {
	return &PluginBuilder{
		logger:     logger,
		configPath: configPath,
		outputDir:  outputDir,
		tempDir:    filepath.Join(outputDir, ".tmp"),
	}
}

// loadConfig loads and adapts the otelcol-builder.yaml configuration
func (b *PluginBuilder) loadConfig() (*BuilderConfig, error) {
	return adaptLegacyConfig(b.configPath)
}

// generatePluginCode generates the Go code for a plugin
func (b *PluginBuilder) generatePluginCode(cfg PluginConfig, compType ComponentType) (string, error) {
	// Extract the actual component name from the gomod path
	parts := strings.Split(cfg.GoMod, "/")
	componentName := parts[len(parts)-1]
	if idx := strings.Index(componentName, " "); idx != -1 {
		componentName = componentName[:idx]
	}

	data := struct {
		Type string
		Name string
	}{
		Type: strings.TrimSuffix(string(compType), "s"),
		Name: componentName,
	}

	tmpl, err := template.New("plugin").Parse(pluginTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var code strings.Builder
	if err := tmpl.Execute(&code, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return code.String(), nil
}

// initGoModule initializes a Go module in the given directory with the required dependencies
func (b *PluginBuilder) initGoModule(dir string, cfg PluginConfig, compType ComponentType, name string) error {
	b.logger.Debug("Initializing Go module",
		zap.String("directory", dir),
		zap.String("gomod", cfg.GoMod))

	// Parse module path and version
	parts := strings.Split(cfg.GoMod, " ")
	if len(parts) != 2 {
		return fmt.Errorf("invalid gomod format, expected 'path version', got %q", cfg.GoMod)
	}
	modPath := parts[0]
	modVersion := parts[1]

	// Generate unique module name
	moduleName := fmt.Sprintf("github.com/srediag/srediag/plugins/%s/%s", compType, name)

	// Initialize module
	modInit := exec.Command("go", "mod", "init", moduleName)
	modInit.Dir = dir
	modInit.Stdout = os.Stdout
	modInit.Stderr = os.Stderr
	if err := modInit.Run(); err != nil {
		return fmt.Errorf("failed to initialize Go module: %w", err)
	}

	// Create go.mod content with latest stable versions
	goModContent := fmt.Sprintf(`module %s

go 1.24

require (
	%s %s
	github.com/srediag/srediag v0.0.0
)

replace github.com/srediag/srediag => ../../../../
`, moduleName, modPath, modVersion)

	// Write go.mod
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goModContent), 0644); err != nil {
		return fmt.Errorf("failed to write go.mod: %w", err)
	}

	// Tidy modules
	modTidy := exec.Command("go", "mod", "tidy")
	modTidy.Dir = dir
	modTidy.Stdout = os.Stdout
	modTidy.Stderr = os.Stderr
	if err := modTidy.Run(); err != nil {
		return fmt.Errorf("failed to tidy modules: %w", err)
	}

	return nil
}

// BuildPlugin builds a single plugin
func (b *PluginBuilder) BuildPlugin(name string, cfg PluginConfig, compType ComponentType) error {
	b.logger.Info("Building plugin",
		zap.String("name", name),
		zap.String("type", string(compType)),
		zap.String("gomod", cfg.GoMod))

	// Extract the actual component name from the gomod path
	parts := strings.Split(cfg.GoMod, "/")
	componentName := parts[len(parts)-1]
	if idx := strings.Index(componentName, " "); idx != -1 {
		componentName = componentName[:idx]
	}

	// Create temp directory for plugin build using component name
	tempDir, err := os.MkdirTemp(b.tempDir, fmt.Sprintf("%s_%s_*", compType, componentName))
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Generate plugin code
	code, err := b.generatePluginCode(cfg, compType)
	if err != nil {
		return fmt.Errorf("failed to generate code: %w", err)
	}

	// Write plugin code to file
	pluginFile := filepath.Join(tempDir, "plugin.go")
	if err := os.WriteFile(pluginFile, []byte(code), 0644); err != nil {
		return fmt.Errorf("failed to write plugin code: %w", err)
	}

	// Initialize Go module
	if err := b.initGoModule(tempDir, cfg, compType, componentName); err != nil {
		return fmt.Errorf("failed to initialize module: %w", err)
	}

	// Create output directory structure plugins/type/name
	outputDir := filepath.Join(b.outputDir, string(compType))
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Build plugin as executable
	outputFile := filepath.Join(outputDir, componentName)
	buildCmd := exec.Command("go", "build",
		"-trimpath",
		"-o", outputFile)
	buildCmd.Dir = tempDir
	buildCmd.Env = append(os.Environ(),
		"CGO_ENABLED=0", // Disable CGO for better portability
		fmt.Sprintf("GOOS=%s", os.Getenv("GOOS")),
		fmt.Sprintf("GOARCH=%s", os.Getenv("GOARCH")))
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("failed to build plugin: %w", err)
	}

	b.logger.Info("Successfully built plugin",
		zap.String("name", componentName),
		zap.String("type", string(compType)),
		zap.String("output", outputFile))

	return nil
}

// BuildAll builds all plugins defined in the configuration
func (b *PluginBuilder) BuildAll() error {
	config, err := b.loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Clean up and recreate output and temp directories
	for _, dir := range []string{b.outputDir, b.tempDir} {
		if err := os.RemoveAll(dir); err != nil {
			return fmt.Errorf("failed to clean directory %s: %w", dir, err)
		}
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Build plugins for each component type
	for compType, components := range config.Components {
		for name, cfg := range components {
			if err := b.BuildPlugin(name, cfg, compType); err != nil {
				b.logger.Error("Failed to build plugin",
					zap.String("name", name),
					zap.String("type", string(compType)),
					zap.Error(err))
				continue
			}
		}
	}

	return nil
}

// GenerateAll generates code for all plugins defined in the configuration
func (b *PluginBuilder) GenerateAll() error {
	b.logger.Info("Loading configuration", zap.String("config", b.configPath))

	config, err := b.loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create output directory
	if err := os.MkdirAll(b.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	b.logger.Info("Generating plugin code",
		zap.String("output_dir", b.outputDir),
		zap.Int("components", len(config.Components)))

	// Generate code for each component type
	for compType, components := range config.Components {
		for name, cfg := range components {
			// Create directory for plugin
			pluginDir := filepath.Join(b.outputDir, fmt.Sprintf("%s_%s", compType, name))
			if err := os.MkdirAll(pluginDir, 0755); err != nil {
				return fmt.Errorf("failed to create plugin directory: %w", err)
			}

			// Generate plugin code
			code, err := b.generatePluginCode(cfg, compType)
			if err != nil {
				return fmt.Errorf("failed to generate code for %s_%s: %w", compType, name, err)
			}

			// Write plugin code
			if err := os.WriteFile(filepath.Join(pluginDir, "plugin.go"), []byte(code), 0644); err != nil {
				return fmt.Errorf("failed to write plugin code: %w", err)
			}

			// Initialize Go module
			if err := b.initGoModule(pluginDir, cfg, compType, name); err != nil {
				return fmt.Errorf("failed to initialize Go module for %s_%s: %w", compType, name, err)
			}

			b.logger.Info("Generated plugin code",
				zap.String("name", name),
				zap.String("type", string(compType)),
				zap.String("directory", pluginDir))
		}
	}

	return nil
}
