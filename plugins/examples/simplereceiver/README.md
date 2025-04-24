# Simple Receiver Example Plugin

This is an example plugin that demonstrates the structure of a custom OpenTelemetry receiver plugin for SREDIAG.

## Overview

The Simple Receiver is a basic example that would generate metrics at regular intervals specified in the configuration. In a real implementation, it would collect actual metrics from a source system.

## Building the Plugin

To build this plugin:

```bash
go build -buildmode=plugin -o ../../simplereceiver.so .
```

This will create a `simplereceiver.so` file in the plugins directory that can be loaded by SREDIAG.

## Configuration

In your SREDIAG configuration file, you can configure the Simple Receiver as follows:

```yaml
receivers:
  simple:
    interval: 30s  # Metrics generation interval (default: 15s)

service:
  pipelines:
    metrics:
      receivers: [simple]
      processors: []
      exporters: [otlp]
```

## Implementation Notes

This example demonstrates:

1. How to export a `Factory` symbol that can be loaded by SREDIAG
2. Basic receiver structure with Start and Shutdown methods
3. Configuration handling

In a real implementation, you would:

1. Use the actual OpenTelemetry helper functions to create factories
2. Implement actual metric collection logic
3. Add proper error handling and lifecycle management
4. Include comprehensive tests

## Limitations

This example is intentionally simplified and incomplete to demonstrate the structure without requiring specific OpenTelemetry dependencies. A real plugin would need to address:

1. Proper imports from OpenTelemetry packages
2. Complete implementation of the receiver interfaces
3. Generation and export of actual metrics to the consumer
4. Proper stability and performance considerations
