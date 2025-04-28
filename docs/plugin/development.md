# SREDIAG Plugin Development Guide

## Table of Contents

1. [Getting Started](#getting-started)
2. [OpenTelemetry Integration](#opentelemetry-integration)
3. [Plugin Development](#plugin-development)
4. [Testing](#testing)
5. [Packaging and Distribution](#packaging-and-distribution)

## Getting Started

### Prerequisites

- Go 1.21 or later
- OpenTelemetry Collector SDK
- SREDIAG SDK
- Docker (for containerized development)

### Project Setup

```bash
# Create new plugin project
srediag plugin create my-plugin --type receiver|processor|exporter

# Initialize Go module
go mod init github.com/username/my-plugin

# Get dependencies
go get go.opentelemetry.io/collector
go get github.com/srediag/srediag/pkg/plugin
```

## OpenTelemetry Integration

### Component Factory Registration

```go
package myreceiver

import (
    "go.opentelemetry.io/collector/component"
    "github.com/srediag/srediag/pkg/plugin"
)

func NewFactory() component.ReceiverFactory {
    return plugin.NewReceiverFactory(
        "myreceiver",
        createDefaultConfig,
        createReceiver,
    )
}

func createDefaultConfig() component.Config {
    return &Config{
        DiagnosticSettings: plugin.DefaultDiagnosticSettings(),
    }
}

func createReceiver(
    ctx context.Context,
    params component.ReceiverCreateSettings,
    cfg component.Config,
    nextConsumer consumer.Metrics,
) (component.MetricsReceiver, error) {
    // Implementation
}
```

### Data Model Integration

```go
type DiagnosticReceiver struct {
    // OpenTelemetry fields
    config   *Config
    settings component.ReceiverCreateSettings
    consumer consumer.Metrics
    
    // SREDIAG extensions
    diagnostics plugin.DiagnosticCapabilities
    patterns    plugin.PatternDetector
}

func (r *DiagnosticReceiver) Start(ctx context.Context, host component.Host) error {
    // Initialize both OTel and SREDIAG components
    if err := r.initOpenTelemetry(ctx, host); err != nil {
        return err
    }
    return r.initDiagnostics(ctx)
}
```

## Plugin Development

### Directory Structure

```text
my-plugin/
├── cmd/
│   └── myreceiver/
│       └── main.go
├── config.go
├── config_test.go
├── doc.go
├── factory.go
├── go.mod
├── go.sum
├── metadata.yaml
├── receiver.go
├── receiver_test.go
└── testdata/
    └── config.yaml
```

### Configuration

```go
type Config struct {
    // Embed OpenTelemetry config
    *otelconfig.ReceiverSettings `mapstructure:",squash"`
    
    // SREDIAG diagnostic settings
    DiagnosticSettings plugin.DiagnosticSettings `mapstructure:"diagnostics"`
    
    // Plugin-specific settings
    CollectionInterval time.Duration     `mapstructure:"collection_interval"`
    Patterns          []PatternRule     `mapstructure:"patterns"`
    Resources         ResourceConfig    `mapstructure:"resources"`
}

func (c *Config) Validate() error {
    if err := c.ReceiverSettings.Validate(); err != nil {
        return err
    }
    return c.DiagnosticSettings.Validate()
}
```

Example configuration:

```yaml
receivers:
  myreceiver:
    collection_interval: 10s
    diagnostics:
      patterns:
        - name: error_spike
          threshold: 100
      correlation:
        enabled: true
        window: 5m
    resources:
      cpu_limit: 1
      memory_limit: 512Mi
```

### Metrics Pipeline

```go
func (r *DiagnosticReceiver) collectMetrics(ctx context.Context) error {
    // 1. Collect raw metrics
    metrics, err := r.gatherMetrics(ctx)
    if err != nil {
        return err
    }
    
    // 2. Apply diagnostic analysis
    diagnostics, err := r.analyzeDiagnostics(ctx, metrics)
    if err != nil {
        return err
    }
    
    // 3. Create OTel metric data
    md := r.convertToOTel(metrics, diagnostics)
    
    // 4. Push to next consumer
    return r.consumer.ConsumeMetrics(ctx, md)
}
```

## Testing

### Unit Tests

```go
func TestReceiver_Start(t *testing.T) {
    // Create test context
    ctx := context.Background()
    
    // Create mock consumer
    consumer := consumertest.NewNop()
    
    // Create receiver with test config
    factory := NewFactory()
    cfg := factory.CreateDefaultConfig()
    receiver, err := factory.CreateMetricsReceiver(
        ctx,
        component.ReceiverCreateSettings{},
        cfg,
        consumer,
    )
    require.NoError(t, err)
    
    // Test startup
    err = receiver.Start(ctx, componenttest.NewNopHost())
    require.NoError(t, err)
}
```

### Integration Tests

```go
func TestReceiver_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    
    // Create test pipeline
    factories, err := componenttest.NopFactories()
    require.NoError(t, err)
    
    factories.Receivers[typeStr] = NewFactory()
    
    cfg, err := configtest.LoadConfigAndValidate(path.Join("testdata", "config.yaml"), factories)
    require.NoError(t, err)
    
    // Run test pipeline
    pipeline, err := service.BuildAndStart(
        context.Background(),
        cfg,
        factories,
    )
    require.NoError(t, err)
    defer pipeline.Shutdown(context.Background())
    
    // Verify results
    // ...
}
```

## Packaging and Distribution

### Build Configuration

```makefile
.PHONY: build
build:
    go build -o bin/receiver ./cmd/myreceiver

.PHONY: docker
docker:
    docker build -t srediag/myreceiver:latest .

.PHONY: test
test:
    go test -v ./...
    go test -v ./... -tags=integration

.PHONY: lint
lint:
    golangci-lint run
```

### Docker Support

```dockerfile
FROM golang:1.21 as builder
WORKDIR /build
COPY . .
RUN go build -o receiver ./cmd/myreceiver

FROM gcr.io/distroless/base
COPY --from=builder /build/receiver /
ENTRYPOINT ["/receiver"]
```

### Distribution Manifest

```yaml
name: myreceiver
version: 1.0.0
type: receiver
compatibility:
  otel_version: "0.81.0"
  srediag_version: "1.0.0"
supported_signals:
  - metrics
  - diagnostics
capabilities:
  pattern_detection: true
  correlation: true
```

## See Also

- [OpenTelemetry Collector Development](https://opentelemetry.io/docs/collector/development/)
- [Plugin System Overview](README.md)
- [API Reference](../reference/api.md)
- [Examples](examples/README.md)
