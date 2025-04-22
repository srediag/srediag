# Plugin API Reference

This document provides detailed information about the SREDIAG plugin API.

## Core Interfaces

### IPlugin

The main interface that all plugins must implement.

```go
type IPlugin interface {
    // Name returns the unique identifier of the plugin
    Name() string

    // Description provides a brief overview of the plugin's functionality
    Description() string

    // Version returns the semantic version of the plugin
    Version() string

    // Start initializes the plugin with the given context
    Start(ctx context.Context) error

    // Shutdown performs cleanup operations
    Shutdown(ctx context.Context) error

    // Execute runs the plugin's main functionality
    Execute(ctx context.Context, input interface{}) (interface{}, error)
}
```

## Plugin Types

### Diagnostic Plugin

Interface for plugins that collect system metrics and perform diagnostics.

```go
type DiagnosticPlugin interface {
    IPlugin
    
    // CollectMetrics gathers system metrics
    CollectMetrics(ctx context.Context) ([]Metric, error)
    
    // PerformCheck executes diagnostic checks
    PerformCheck(ctx context.Context) (*CheckResult, error)
}
```

### Analysis Plugin

Interface for plugins that analyze collected data and detect patterns.

```go
type AnalysisPlugin interface {
    IPlugin
    
    // Analyze processes collected data
    Analyze(ctx context.Context, data []byte) (*Analysis, error)
    
    // GetInsights retrieves analysis results
    GetInsights(ctx context.Context) ([]Insight, error)
}
```

### Management Plugin

Interface for plugins that manage system resources and configurations.

```go
type ManagementPlugin interface {
    IPlugin
    
    // ApplyConfig applies configuration changes
    ApplyConfig(ctx context.Context, config []byte) error
    
    // GetStatus retrieves current management status
    GetStatus(ctx context.Context) (*Status, error)
}
```

## Common Types

### Metric

Represents a single metric data point.

```go
type Metric struct {
    Name      string
    Value     float64
    Labels    map[string]string
    Timestamp time.Time
}
```

### CheckResult

Represents the result of a diagnostic check.

```go
type CheckResult struct {
    Status    string
    Message   string
    Details   map[string]interface{}
    Timestamp time.Time
}
```

### Analysis

Represents the result of data analysis.

```go
type Analysis struct {
    ID        string
    Results   []AnalysisResult
    Metadata  map[string]interface{}
    Timestamp time.Time
}
```

### Status

Represents the current status of a management operation.

```go
type Status struct {
    State     string
    Message   string
    Details   map[string]interface{}
    Timestamp time.Time
}
```

## Error Handling

Plugins should use standard error types and follow these guidelines:

```go
// ErrNotImplemented indicates functionality not implemented
var ErrNotImplemented = errors.New("functionality not implemented")

// ErrInvalidConfig indicates invalid configuration
var ErrInvalidConfig = errors.New("invalid configuration")

// ErrOperationFailed indicates a failed operation
var ErrOperationFailed = errors.New("operation failed")
```

## Configuration

### Plugin Config

Standard configuration structure for plugins.

```go
type PluginConfig struct {
    // Enabled indicates if the plugin is active
    Enabled bool `yaml:"enabled"`

    // Interval between operations
    Interval time.Duration `yaml:"interval"`

    // Settings holds plugin-specific settings
    Settings map[string]interface{} `yaml:"settings"`
}
```

## Events

### Plugin Events

Standard events that plugins can emit.

```go
type PluginEvent struct {
    // Type of the event
    Type string

    // Source plugin that generated the event
    Source string

    // Payload contains event-specific data
    Payload interface{}

    // Timestamp when the event occurred
    Timestamp time.Time
}
```

## Best Practices

1. **Error Handling**
   - Use standard error types
   - Provide context with errors
   - Handle timeouts appropriately

2. **Configuration**
   - Validate all configurations
   - Use strong typing where possible
   - Provide defaults

3. **Context Usage**
   - Always respect context cancellation
   - Implement proper timeouts
   - Clean up resources

4. **Event Emission**
   - Use standard event types
   - Include relevant context
   - Follow naming conventions
