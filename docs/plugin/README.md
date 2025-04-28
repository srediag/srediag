# SREDIAG Plugin System

## Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [OpenTelemetry Integration](#opentelemetry-integration)
4. [Development Guide](#development-guide)
5. [Security](#security)
6. [Best Practices](#best-practices)

## Overview

SREDIAG's plugin system extends the OpenTelemetry Collector architecture, allowing seamless integration of both standard OpenTelemetry components and custom SREDIAG-specific plugins. This hybrid approach enables:

- Direct use of OpenTelemetry Collector core and contrib components
- Custom SREDIAG plugins for specialized diagnostic capabilities
- Unified configuration and management
- Consistent telemetry pipeline

## Architecture

### Component Integration

```ascii
+------------------------+     +-------------------------+
|   SREDIAG Core         |     |   OpenTelemetry         |
|                        |     |   Collector             |
| +------------------+   |     |                         |
| |  Plugin Manager  |   |     |  +-----------------+    |
| |   - Loading      |<---------->| Component       |    |
| |   - Lifecycle    |   |     |  | Factory         |    |
| |   - Config       |   |     |  +-----------------+    |
| +------------------+   |     |          ↑              |
|          ↑             |     |    +---------------+    |
|    +---------------+   |     |    |  Components   |    |
|    |   Plugin API  |   |     |    |  Registry     |    |
|    +---------------+   |     |    +---------------+    |
|          ↑             |     |          ↑              |
| +------------------+   |     | +------------------+    |
| |     Plugins      |   |     | |    Components    |    |
| | - Processors     |   |     | | - Processors     |    |
| | - Receivers      |   |     | | - Receivers      |    |
| | - Exporters      |   |     | | - Exporters      |    |
| +------------------+   |     | +------------------+    |
+------------------------+     +-------------------------+
```

### Plugin Types Alignment

SREDIAG plugins align with OpenTelemetry Collector components while adding diagnostic-specific capabilities:

1. **Receivers** (OpenTelemetry + SREDIAG)
   - Standard OpenTelemetry receivers
   - Custom diagnostic data collectors
   - System-specific data gatherers

2. **Processors** (OpenTelemetry + SREDIAG)
   - Standard OpenTelemetry processors
   - Diagnostic analysis processors
   - Pattern detection and correlation

3. **Exporters** (OpenTelemetry + SREDIAG)
   - Standard OpenTelemetry exporters
   - Diagnostic report generators
   - Custom visualization exporters

## OpenTelemetry Integration

### Component Compatibility

```yaml
receivers:
  # OpenTelemetry Standard Receivers
  otlp:
    protocols:
      grpc:
        endpoint: localhost:4317
  
  # SREDIAG Custom Receivers
  srediag_system:
    type: system_diagnostics
    collection_interval: 10s

processors:
  # OpenTelemetry Standard Processors
  batch:
    timeout: 1s
    send_batch_size: 1024

  # SREDIAG Custom Processors
  srediag_analyzer:
    type: pattern_detection
    rules:
      - pattern: "error_spike"
        threshold: 100

exporters:
  # OpenTelemetry Standard Exporters
  otlp:
    endpoint: "otelcol:4317"
    
  # SREDIAG Custom Exporters
  srediag_report:
    type: diagnostic_report
    format: pdf
```

### Data Model Integration

SREDIAG extends OpenTelemetry's data model while maintaining compatibility:

```go
type DiagnosticData struct {
    // OpenTelemetry Standard Fields
    Resource     *resource.Resource
    ScopeLogs    []ScopeLog
    ScopeMetrics []ScopeMetric
    ScopeTraces  []ScopeTrace

    // SREDIAG Extensions
    Diagnostics  []DiagnosticResult
    Patterns     []DetectedPattern
    Correlations []CorrelationGroup
}
```

## Development Guide

### Creating OpenTelemetry-Compatible Plugins

1. **Implement Standard Interfaces**

   ```go
   type DiagnosticReceiver interface {
      component.Receiver // OpenTelemetry interface
      DiagnosticCapabilities() []DiagnosticType
   }

   type DiagnosticProcessor interface {
      component.Processor // OpenTelemetry interface
      AnalyzePatterns(context.Context, DiagnosticData) error
   }
   ```

2. **Configuration Integration**

   ```yaml
   srediag:
   plugins:
      path: /etc/srediag/plugins
      auto_discovery: true
      compatibility:
         otel_version: "0.81.0"
         srediag_version: "1.0.0"
   ```

   ### Plugin Lifecycle

   ```go
   func (p *DiagnosticPlugin) Start(ctx context.Context, host component.Host) error {
      // Initialize OpenTelemetry components
      if err := p.initOtel(ctx, host); err != nil {
         return fmt.Errorf("failed to init OTel: %w", err)
      }
      
      // Start SREDIAG-specific functionality
      if err := p.startDiagnostics(ctx); err != nil {
         return fmt.Errorf("failed to start diagnostics: %w", err)
      }
      
      return nil
   }
   ```

## Security

### Plugin Isolation

```yaml
security:
  plugins:
    isolation:
      mode: "container"  # or "process"
      capabilities:
        drop: ["ALL"]
        add: ["NET_ADMIN"]  # If needed
      resources:
        limits:
          cpu: "1"
          memory: "512Mi"
```

### Access Control

```yaml
security:
  plugins:
    permissions:
      otel_components: ["otlp", "batch", "memory_limiter"]
      srediag_components: ["system_diagnostics", "pattern_analyzer"]
      capabilities:
        metrics: true
        traces: true
        logs: true
        diagnostics: true
```

## Best Practices

1. **OpenTelemetry Compatibility**
   - Follow OpenTelemetry data model conventions
   - Implement standard component interfaces
   - Use OpenTelemetry SDK for telemetry

2. **Performance**
   - Implement batching for data processing
   - Use appropriate buffer sizes
   - Follow OpenTelemetry performance guidelines

3. **Error Handling**
   - Implement graceful degradation
   - Provide detailed error context
   - Follow OpenTelemetry error handling patterns

4. **Configuration**
   - Use YAML for consistency with OpenTelemetry
   - Provide validation
   - Support dynamic updates where possible

## See Also

- [OpenTelemetry Collector Documentation](https://opentelemetry.io/docs/collector/)
- [Plugin Development Guide](development.md)
- [Plugin Examples](examples/README.md)
- [API Reference](../reference/api.md)
- [Architecture Overview](../architecture/overview.md)
