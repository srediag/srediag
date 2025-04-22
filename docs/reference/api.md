# SREDIAG Plugin API Reference

## Table of Contents

1. [Overview](#overview)
2. [Core Interfaces](#core-interfaces)
3. [Data Types](#data-types)
4. [Plugin Lifecycle](#plugin-lifecycle)
5. [Error Handling](#error-handling)
6. [Configuration](#configuration)
7. [Metrics and Telemetry](#metrics-and-telemetry)

## Overview

The SREDIAG Plugin API provides interfaces and types for developing custom plugins. This document details the complete API reference for plugin developers.

## Core Interfaces

### Plugin Interface

The base interface that all plugins must implement:

```go
type Plugin interface {
    // Initialize sets up the plugin with the provided configuration
    Initialize(ctx context.Context, config *Config) error
    
    // Start begins plugin operation
    Start(ctx context.Context) error
    
    // Stop gracefully shuts down the plugin
    Stop(ctx context.Context) error
    
    // Status returns the current plugin status
    Status() *PluginStatus
}
```

### Processor Interface

Interface for plugins that process data:

```go
type Processor interface {
    Plugin
    
    // Process handles incoming data and returns processed result
    Process(ctx context.Context, data []byte) ([]byte, error)
    
    // Validate checks if data format is valid
    Validate(data []byte) error
}
```

### Collector Interface

Interface for plugins that collect metrics:

```go
type Collector interface {
    Plugin
    
    // Collect gathers metrics and returns them
    Collect(ctx context.Context) ([]Metric, error)
    
    // GetMetricTypes returns supported metric types
    GetMetricTypes() []MetricType
}
```

### Exporter Interface

Interface for plugins that export data:

```go
type Exporter interface {
    Plugin
    
    // Export sends metrics to external systems
    Export(ctx context.Context, metrics []Metric) error
    
    // Flush ensures all metrics are exported
    Flush(ctx context.Context) error
}
```

## Data Types

### Config

Configuration structure for plugins:

```go
type Config struct {
    // Name is the unique identifier of the plugin
    Name string `yaml:"name"`
    
    // Type specifies the plugin type (processor/collector/exporter)
    Type string `yaml:"type"`
    
    // Version of the plugin
    Version string `yaml:"version"`
    
    // Settings contains plugin-specific configuration
    Settings map[string]interface{} `yaml:"settings"`
    
    // Tags for plugin categorization
    Tags []string `yaml:"tags,omitempty"`
}
```

### Metric

Represents a single metric:

```go
type Metric struct {
    // Name of the metric
    Name string `json:"name"`
    
    // Value of the metric
    Value float64 `json:"value"`
    
    // Labels/tags associated with the metric
    Labels map[string]string `json:"labels"`
    
    // Timestamp when the metric was collected
    Timestamp time.Time `json:"timestamp"`
    
    // Description provides additional context
    Description string `json:"description,omitempty"`
}
```

### PluginStatus

Represents the current state of a plugin:

```go
type PluginStatus struct {
    // State indicates the plugin's operational state
    State PluginState `json:"state"`
    
    // LastError contains the most recent error
    LastError string `json:"last_error,omitempty"`
    
    // StartTime when the plugin was initialized
    StartTime time.Time `json:"start_time"`
    
    // Metrics contains plugin-specific metrics
    Metrics map[string]interface{} `json:"metrics,omitempty"`
}
```

### PluginState

Enum for plugin states:

```go
type PluginState string

const (
    StateInitializing PluginState = "initializing"
    StateRunning      PluginState = "running"
    StateStopped      PluginState = "stopped"
    StateError        PluginState = "error"
)
```

## Plugin Lifecycle

### Initialization

```go
func (p *YourPlugin) Initialize(ctx context.Context, config *Config) error {
    // 1. Validate configuration
    if err := p.validateConfig(config); err != nil {
        return fmt.Errorf("invalid configuration: %w", err)
    }
    
    // 2. Set up resources
    if err := p.setupResources(config); err != nil {
        return fmt.Errorf("failed to setup resources: %w", err)
    }
    
    // 3. Initialize metrics
    if err := p.initMetrics(); err != nil {
        return fmt.Errorf("failed to initialize metrics: %w", err)
    }
    
    return nil
}
```

### Start/Stop

```go
func (p *YourPlugin) Start(ctx context.Context) error {
    // 1. Start background workers
    if err := p.startWorkers(ctx); err != nil {
        return fmt.Errorf("failed to start workers: %w", err)
    }
    
    // 2. Begin processing/collecting/exporting
    go p.run(ctx)
    
    return nil
}

func (p *YourPlugin) Stop(ctx context.Context) error {
    // 1. Stop accepting new work
    p.stopWorkers()
    
    // 2. Flush pending operations
    if err := p.flush(ctx); err != nil {
        return fmt.Errorf("failed to flush: %w", err)
    }
    
    // 3. Clean up resources
    p.cleanup()
    
    return nil
}
```

## Error Handling

### Error Types

```go
var (
    // ErrInvalidConfig indicates invalid plugin configuration
    ErrInvalidConfig = errors.New("invalid configuration")
    
    // ErrNotInitialized indicates plugin not properly initialized
    ErrNotInitialized = errors.New("plugin not initialized")
    
    // ErrAlreadyRunning indicates plugin is already running
    ErrAlreadyRunning = errors.New("plugin already running")
)
```

### Error Handling Best Practices

```go
func (p *YourPlugin) Process(ctx context.Context, data []byte) ([]byte, error) {
    // 1. Validate input
    if len(data) == 0 {
        return nil, fmt.Errorf("%w: empty data", ErrInvalidInput)
    }
    
    // 2. Handle context cancellation
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
    }
    
    // 3. Process with timeout
    result, err := p.processWithTimeout(ctx, data)
    if err != nil {
        // Log error details
        p.logger.Error("Processing failed",
            zap.Error(err),
            zap.Int("data_size", len(data)))
        return nil, fmt.Errorf("processing failed: %w", err)
    }
    
    return result, nil
}
```

## Configuration

### YAML Configuration Example

```yaml
plugins:
  your-plugin:
    enabled: true
    type: processor
    version: "1.0.0"
    settings:
      interval: 10s
      batch_size: 100
      timeout: 5s
      retry:
        enabled: true
        max_attempts: 3
        initial_interval: 1s
    tags:
      - production
      - metrics
```

### Configuration Validation

```go
func (p *YourPlugin) validateConfig(config *Config) error {
    // 1. Check required fields
    if config.Name == "" {
        return fmt.Errorf("%w: name is required", ErrInvalidConfig)
    }
    
    // 2. Validate settings
    settings, ok := config.Settings.(map[string]interface{})
    if !ok {
        return fmt.Errorf("%w: invalid settings format", ErrInvalidConfig)
    }
    
    // 3. Type-specific validation
    if err := p.validateTypeSpecific(settings); err != nil {
        return fmt.Errorf("%w: %v", ErrInvalidConfig, err)
    }
    
    return nil
}
```

## Metrics and Telemetry

### Standard Metrics

```go
var (
    // ProcessingTime measures time taken to process data
    ProcessingTime = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "plugin_processing_time_seconds",
            Help: "Time taken to process data",
            Buckets: prometheus.DefBuckets,
        },
        []string{"plugin_name", "status"},
    )
    
    // ProcessedItems counts number of items processed
    ProcessedItems = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "plugin_processed_items_total",
            Help: "Total number of items processed",
        },
        []string{"plugin_name", "status"},
    )
)
```

### Telemetry Integration

```go
func (p *YourPlugin) recordMetrics(ctx context.Context, start time.Time, status string) {
    duration := time.Since(start).Seconds()
    
    // Record processing time
    ProcessingTime.WithLabelValues(p.config.Name, status).Observe(duration)
    
    // Increment processed items counter
    ProcessedItems.WithLabelValues(p.config.Name, status).Inc()
    
    // Log detailed metrics
    p.logger.Debug("Processing metrics",
        zap.String("plugin", p.config.Name),
        zap.Float64("duration_seconds", duration),
        zap.String("status", status))
}
```

## See Also

- [Plugin Development Guide](../plugins/development.md)
- [Plugin Examples](../plugins/examples/README.md)
- [Best Practices](best-practices.md)
- [Troubleshooting Guide](../development/troubleshooting.md)
