package main

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"

	"github.com/srediag/srediag/cmd/srediag/commands"
	"github.com/srediag/srediag/internal/components"
	"github.com/srediag/srediag/internal/config"
	"github.com/srediag/srediag/internal/plugin"
	"github.com/srediag/srediag/internal/settings"
)

// Component types
const (
	typeReceiver  = "receiver"
	typeProcessor = "processor"
	typeExporter  = "exporter"
	typeExtension = "extension"
)

func main() {
	ctx := context.Background()

	// Initialize logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to sync logger: %v\n", err)
		}
	}()

	// Load configuration
	cfg, err := config.Load(logger)
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Initialize managers
	componentManager := components.NewManager(logger)
	pluginManager := plugin.NewManager(logger, cfg.PluginsDir)

	// Initialize settings
	settings := &settings.CommandSettings{
		ComponentManager: componentManager,
		PluginManager:    pluginManager,
		Logger:           logger,
	}

	// Only load components if not running plugin generate command
	if !isPluginGenerateCommand() {
		// Load and register components
		if err := initializeComponents(ctx, settings); err != nil {
			logger.Fatal("Failed to initialize components", zap.Error(err))
		}
	}

	// Execute root command
	if err := commands.Execute(settings); err != nil {
		logger.Fatal("Failed to execute command", zap.Error(err))
	}
}

// isPluginGenerateCommand checks if we're running the plugin generate command
func isPluginGenerateCommand() bool {
	if len(os.Args) < 2 {
		return false
	}
	return os.Args[1] == "plugin" && len(os.Args) > 2 && os.Args[2] == "generate"
}

// initializeComponents loads and registers all components
func initializeComponents(ctx context.Context, settings *settings.CommandSettings) error {
	// Load core components
	if err := loadCoreComponents(ctx, settings.PluginManager); err != nil {
		return fmt.Errorf("failed to load core components: %w", err)
	}

	// Register components
	if err := registerPluginComponents(ctx, settings.PluginManager, settings.ComponentManager); err != nil {
		return fmt.Errorf("failed to register components: %w", err)
	}

	return nil
}

// loadCoreComponents loads the core OpenTelemetry components
func loadCoreComponents(ctx context.Context, pm *plugin.Manager) error {
	coreComponents := []struct {
		typ  string
		name string
	}{
		{typeReceiver, "otlp"},
		{typeProcessor, "batch"},
		{typeExporter, "otlp"},
		{typeExtension, "zpages"},
	}

	for _, comp := range coreComponents {
		var typ component.Type
		var err error
		switch comp.typ {
		case typeReceiver:
			typ, err = component.NewType("receiver")
		case typeProcessor:
			typ, err = component.NewType("processor")
		case typeExporter:
			typ, err = component.NewType("exporter")
		case typeExtension:
			typ, err = component.NewType("extension")
		default:
			return fmt.Errorf("unknown component type: %s", comp.typ)
		}
		if err != nil {
			return fmt.Errorf("failed to create component type %s: %w", comp.typ, err)
		}

		if err := pm.LoadPlugin(ctx, typ, comp.name); err != nil {
			return fmt.Errorf("failed to load %s plugin %s: %w", comp.typ, comp.name, err)
		}
	}

	return nil
}

// registerPluginComponents registers components from loaded plugins
func registerPluginComponents(ctx context.Context, pm *plugin.Manager, cm *components.Manager) error {
	// Create component types
	var componentTypes []component.Type
	for _, t := range []string{typeReceiver, typeProcessor, typeExporter, typeExtension} {
		typ, err := component.NewType(t)
		if err != nil {
			return fmt.Errorf("failed to create component type %s: %w", t, err)
		}
		componentTypes = append(componentTypes, typ)
	}

	// Register factories for each type
	for _, typ := range componentTypes {
		factory, err := pm.GetFactory(typ)
		if err != nil {
			// Skip if plugin not loaded
			continue
		}

		switch typ.String() {
		case typeReceiver:
			if err := cm.RegisterReceiver(factory); err != nil {
				return fmt.Errorf("failed to register receiver: %w", err)
			}
		case typeProcessor:
			if err := cm.RegisterProcessor(factory); err != nil {
				return fmt.Errorf("failed to register processor: %w", err)
			}
		case typeExporter:
			if err := cm.RegisterExporter(factory); err != nil {
				return fmt.Errorf("failed to register exporter: %w", err)
			}
		case typeExtension:
			if err := cm.RegisterExtension(factory); err != nil {
				return fmt.Errorf("failed to register extension: %w", err)
			}
		}
	}

	return nil
}
