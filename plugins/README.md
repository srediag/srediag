# SREDIAG Plugins

This directory contains custom plugins for SREDIAG. Plugins extend the functionality of SREDIAG by adding new components that can be used in your diagnostics and monitoring pipelines.

## Plugin Types

SREDIAG supports the following types of plugins:

- **Receivers**: Collect data from various sources
- **Processors**: Process and transform data
- **Exporters**: Send data to various destinations
- **Extensions**: Add functionality to the collector
- **Connectors**: Connect different components together

## Creating a Plugin

To create a plugin:

1. Create a new Go module for your plugin
2. Implement the appropriate factory interface from OpenTelemetry:
   - `receiver.Factory` for receivers
   - `processor.Factory` for processors
   - `exporter.Factory` for exporters
   - `extension.Factory` for extensions
   - `connector.Factory` for connectors

3. Export your factory as a symbol named `Factory`
4. Build your plugin as a shared object (`.so` file)
5. Place the `.so` file in this directory

Example plugin structure:

```go
package main

import (
    "go.opentelemetry.io/collector/component"
    "go.opentelemetry.io/collector/receiver"
)

// Factory is the symbol that will be loaded by SREDIAG
var Factory receiver.Factory

func init() {
    // Initialize your factory here
    Factory = receiver.NewFactory(
        "your_receiver",
        createDefaultConfig,
        receiver.WithMetrics(createMetricsReceiver, stability),
    )
}

func createDefaultConfig() component.Config {
    return &Config{
        // Your default configuration
    }
}

func createMetricsReceiver(
    ctx context.Context,
    params receiver.CreateSettings,
    cfg component.Config,
    consumer consumer.Metrics,
) (receiver.Metrics, error) {
    // Create and return your receiver
}
```

## Building a Plugin

To build your plugin as a shared object:

```bash
go build -buildmode=plugin -o plugins/your_plugin.so your_plugin_dir/main.go
```

## Plugin Loading

SREDIAG will automatically load plugins from this directory at startup. You can specify a different plugin directory using the `SREDIAG_PLUGIN_DIR` environment variable.

## Plugin Configuration

Configure your plugins in the SREDIAG configuration file (`srediag.yaml`). Example:

```yaml
receivers:
  your_receiver:
    # Your receiver configuration

processors:
  your_processor:
    # Your processor configuration

exporters:
  your_exporter:
    # Your exporter configuration

extensions:
  your_extension:
    # Your extension configuration

service:
  pipelines:
    metrics:
      receivers: [your_receiver]
      processors: [your_processor]
      exporters: [your_exporter]
  extensions: [your_extension]
```

## Best Practices

1. Use semantic versioning for your plugins
2. Document your plugin's configuration options
3. Include example configurations
4. Test your plugin thoroughly
5. Handle errors gracefully
6. Follow OpenTelemetry's stability guidelines
7. Use appropriate logging and telemetry

## Troubleshooting

If a plugin fails to load:

1. Check the plugin file permissions
2. Verify the plugin was built with the correct Go version
3. Ensure all dependencies are available
4. Check SREDIAG logs for error messages
5. Verify the plugin implements the correct interface
6. Make sure the `Factory` symbol is exported correctly
