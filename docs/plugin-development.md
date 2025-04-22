# Plugin Development Guide

This guide explains how to develop plugins for SREDIAG.

## Overview

SREDIAG uses a plugin architecture to extend its functionality. Plugins can collect metrics, logs, and traces from various sources and send them to the OBSERVO platform.

## Plugin Types

### 1. Metrics Collector

Metrics collectors gather numerical data about system or application performance.

```go
type MetricsCollector interface {
    // Initialize the collector
    Init(config map[string]interface{}) error
    
    // Start collecting metrics
    Start(ctx context.Context) error
    
    // Stop collecting metrics
    Stop(ctx context.Context) error
    
    // Get collector metadata
    Metadata() CollectorMetadata
}
```

### 2. Log Collector

Log collectors gather and process log entries from various sources.

```go
type LogCollector interface {
    // Initialize the collector
    Init(config map[string]interface{}) error
    
    // Start collecting logs
    Start(ctx context.Context) error
    
    // Stop collecting logs
    Stop(ctx context.Context) error
    
    // Get collector metadata
    Metadata() CollectorMetadata
}
```

### 3. Trace Collector

Trace collectors gather distributed tracing data.

```go
type TraceCollector interface {
    // Initialize the collector
    Init(config map[string]interface{}) error
    
    // Start collecting traces
    Start(ctx context.Context) error
    
    // Stop collecting traces
    Stop(ctx context.Context) error
    
    // Get collector metadata
    Metadata() CollectorMetadata
}
```

## Creating a Plugin

1. Create a new Go module for your plugin:

    ```bash
    mkdir my-plugin
    cd my-plugin
    go mod init github.com/username/my-plugin
    ```

2. Implement the appropriate interface for your plugin type.

3. Create a plugin manifest:

    ```yaml
    name: "my-plugin"
    version: "1.0.0"
    type: "metrics"  # or "logs" or "traces"
    author: "Your Name"
    description: "Description of your plugin"
    ```

4. Build your plugin:

```bash
go build -buildmode=plugin -o my-plugin.so
```

## Plugin Configuration

Plugins are configured through the main SREDIAG configuration file:

```yaml
plugins:
  directory: "plugins"
  enabled:
    - "my-plugin"
  settings:
    my-plugin:
      interval: "10s"
      # other plugin-specific settings
```

## Best Practices

1. **Error Handling**
   - Handle errors gracefully
   - Provide meaningful error messages
   - Don't panic in plugin code

2. **Resource Management**
   - Clean up resources in Stop()
   - Use context for cancellation
   - Monitor resource usage

3. **Configuration**
   - Validate configuration
   - Provide defaults
   - Document configuration options

4. **Testing**
   - Write unit tests
   - Test error conditions
   - Test configuration parsing

5. **Documentation**
   - Document plugin functionality
   - Provide configuration examples
   - Include usage instructions

## Example Plugin

Here's a simple example of a metrics collector plugin:

```go
package main

import (
    "context"
    "time"
)

type ExampleCollector struct {
    interval time.Duration
    // Add other fields as needed
}

func (c *ExampleCollector) Init(config map[string]interface{}) error {
    // Parse configuration
    if interval, ok := config["interval"].(string); ok {
        d, err := time.ParseDuration(interval)
        if err != nil {
            return err
        }
        c.interval = d
    }
    return nil
}

func (c *ExampleCollector) Start(ctx context.Context) error {
    // Start collection loop
    go func() {
        ticker := time.NewTicker(c.interval)
        defer ticker.Stop()
        
        for {
            select {
            case <-ctx.Done():
                return
            case <-ticker.C:
                // Collect and send metrics
            }
        }
    }()
    return nil
}

func (c *ExampleCollector) Stop(ctx context.Context) error {
    // Cleanup resources
    return nil
}

func (c *ExampleCollector) Metadata() CollectorMetadata {
    return CollectorMetadata{
        Name:        "example-collector",
        Version:     "1.0.0",
        Type:        "metrics",
        Description: "Example metrics collector",
    }
}

// Plugin entry point
var Collector ExampleCollector
```

## Debugging Plugins

1. Use the `debug` configuration option:

    ```yaml
    debug: true
    log_level: "debug"
    ```

2. Use the plugin development tools:

    ```bash
    srediag debug-plugin my-plugin.so
    ```

## Common Issues

1. **Plugin Not Loading**
   - Check file permissions
   - Verify plugin binary compatibility
   - Check for missing dependencies

2. **Configuration Issues**
   - Validate configuration format
   - Check for required fields
   - Verify configuration values

3. **Resource Leaks**
   - Ensure proper cleanup
   - Monitor goroutine leaks
   - Check file descriptor usage

## Support

For plugin development support:

- Check the [documentation](https://github.com/observo/srediag/docs)
- Ask in [GitHub Discussions](https://github.com/observo/srediag/discussions)
- Report issues in the [Issue Tracker](https://github.com/observo/srediag/issues)
