# SREDIAG Examples

This document provides practical examples for using SREDIAG in various scenarios.

## Table of Contents

- [Basic Usage](#basic-usage)
- [Configuration Examples](#configuration-examples)
- [Plugin Development](#plugin-development)
- [Integration Examples](#integration-examples)
- [Monitoring Setup](#monitoring-setup)

## Basic Usage

### Starting SREDIAG Service

```bash
# Start with default configuration
srediag service start

# Start with custom config
srediag service start --config /path/to/config.yaml

# Start with debug logging
srediag service start --log-level debug --log-format json
```

### Running Diagnostics

```bash
# Run system diagnostics
srediag diagnose system

# Run performance diagnostics
srediag diagnose performance --output json

# Run security diagnostics
srediag diagnose security --output-file security-report.json
```

### Plugin Management

```bash
# List available plugins
srediag plugin list

# Load a plugin
srediag plugin load --type receiver --name custom-receiver

# Check plugin health
srediag plugin health --name custom-receiver

# Unload a plugin
srediag plugin unload --name custom-receiver
```

## Configuration Examples

### Basic Service Configuration

```yaml
# srediag.yaml
service:
  log_level: info
  log_format: console
  
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318

processors:
  batch:
    timeout: 1s
    send_batch_size: 1024
  
  memory_limiter:
    limit_mib: 512

exporters:
  otlp:
    endpoint: http://jaeger:14250
    tls:
      insecure: true
  
  prometheus:
    endpoint: "0.0.0.0:8889"

extensions:
  health_check:
    endpoint: 0.0.0.0:13133
  
  pprof:
    endpoint: 0.0.0.0:1777

service:
  extensions: [health_check, pprof]
  pipelines:
    traces:
      receivers: [otlp]
      processors: [memory_limiter, batch]
      exporters: [otlp]
    
    metrics:
      receivers: [otlp]
      processors: [memory_limiter, batch]
      exporters: [prometheus]
```

### Plugin Configuration

```yaml
# plugins.yaml
plugins:
  directory: /opt/srediag/plugins
  
  configs:
    custom-receiver:
      type: receiver
      config:
        endpoint: 0.0.0.0:9999
        protocol: http
        
    custom-processor:
      type: processor
      config:
        buffer_size: 1000
        timeout: 5s
```

### Diagnostics Configuration

```yaml
# diagnostics.yaml
diagnostics:
  system:
    enabled: true
    interval: 30s
    checks:
      - cpu_usage
      - memory_usage
      - disk_usage
      - network_stats
  
  performance:
    enabled: true
    profiling:
      cpu: true
      memory: true
      goroutine: true
  
  security:
    enabled: true
    scans:
      - file_permissions
      - open_ports
      - process_list
```

## Plugin Development

### Basic Plugin Structure

```go
// plugin/custom-receiver/main.go
package main

import (
    "context"
    "fmt"
    
    "go.opentelemetry.io/collector/component"
    "go.opentelemetry.io/collector/receiver"
    
    "github.com/srediag/srediag/internal/plugin"
)

type Config struct {
    Endpoint string `mapstructure:"endpoint"`
    Protocol string `mapstructure:"protocol"`
}

type customReceiver struct {
    config *Config
    logger component.Logger
}

func (r *customReceiver) Start(ctx context.Context, host component.Host) error {
    r.logger.Info("Starting custom receiver", 
        "endpoint", r.config.Endpoint,
        "protocol", r.config.Protocol)
    return nil
}

func (r *customReceiver) Shutdown(ctx context.Context) error {
    r.logger.Info("Shutting down custom receiver")
    return nil
}

func createReceiver(ctx context.Context, set receiver.CreateSettings, cfg component.Config, consumer receiver.Consumer) (receiver.Traces, error) {
    config := cfg.(*Config)
    
    return &customReceiver{
        config: config,
        logger: set.Logger,
    }, nil
}

func createFactory() receiver.Factory {
    return receiver.NewFactory(
        "custom",
        func() component.Config { return &Config{} },
        receiver.WithTraces(createReceiver, component.StabilityLevelAlpha),
    )
}

func main() {
    factory := createFactory()
    plugin.ServeReceiver(factory)
}
```

### Plugin Manifest

```yaml
# plugin/custom-receiver/manifest.yaml
name: custom-receiver
version: 1.0.0
type: receiver
description: Custom receiver for special data sources
capabilities:
  - traces
  - metrics
author: Your Name <email@example.com>
license: Apache-2.0

dependencies:
  go_version: "1.24"
  otel_version: "1.30.0"

config_schema:
  type: object
  properties:
    endpoint:
      type: string
      description: Endpoint to receive data
    protocol:
      type: string
      enum: [http, grpc]
      description: Protocol to use
```

## Integration Examples

### Docker Compose Setup

```yaml
# docker-compose.yml
version: '3.8'

services:
  srediag:
    image: ghcr.io/srediag/srediag:latest
    container_name: srediag
    ports:
      - "4317:4317"   # OTLP gRPC
      - "4318:4318"   # OTLP HTTP
      - "8889:8889"   # Prometheus metrics
      - "13133:13133" # Health check
    volumes:
      - ./configs:/etc/srediag
      - ./plugins:/opt/srediag/plugins
    command: ["service", "start", "--config", "/etc/srediag/srediag.yaml"]
    environment:
      - SREDIAG_LOG_LEVEL=info
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:13133/"]
      interval: 30s
      timeout: 10s
      retries: 3

  jaeger:
    image: jaegertracing/all-in-one:latest
    ports:
      - "16686:16686"
      - "14250:14250"
    environment:
      - COLLECTOR_OTLP_ENABLED=true

  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
```

### Kubernetes Deployment

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: srediag
  namespace: observability
spec:
  replicas: 3
  selector:
    matchLabels:
      app: srediag
  template:
    metadata:
      labels:
        app: srediag
    spec:
      containers:
      - name: srediag
        image: ghcr.io/srediag/srediag:latest
        ports:
        - containerPort: 4317
          name: otlp-grpc
        - containerPort: 4318
          name: otlp-http
        - containerPort: 8889
          name: metrics
        env:
        - name: SREDIAG_LOG_LEVEL
          value: "info"
        volumeMounts:
        - name: config
          mountPath: /etc/srediag
        resources:
          requests:
            memory: "256Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /
            port: 13133
          initialDelaySeconds: 30
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /
            port: 13133
          initialDelaySeconds: 5
          periodSeconds: 10
      volumes:
      - name: config
        configMap:
          name: srediag-config

---
apiVersion: v1
kind: Service
metadata:
  name: srediag
  namespace: observability
spec:
  selector:
    app: srediag
  ports:
  - name: otlp-grpc
    port: 4317
    targetPort: 4317
  - name: otlp-http
    port: 4318
    targetPort: 4318
  - name: metrics
    port: 8889
    targetPort: 8889
```

## Monitoring Setup

### Prometheus Configuration

```yaml
# prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'srediag'
    static_configs:
      - targets: ['srediag:8889']
    scrape_interval: 10s
    metrics_path: /metrics
```

### Grafana Dashboard

```json
{
  "dashboard": {
    "title": "SREDIAG Monitoring",
    "panels": [
      {
        "title": "Received Spans per Second",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(otelcol_receiver_accepted_spans_total[5m])",
            "legendFormat": "Spans/sec"
          }
        ]
      },
      {
        "title": "Memory Usage",
        "type": "graph",
        "targets": [
          {
            "expr": "process_resident_memory_bytes",
            "legendFormat": "Memory Usage"
          }
        ]
      }
    ]
  }
}
```

### Health Check Script

```bash
#!/bin/bash
# health-check.sh

SREDIAG_ENDPOINT="http://localhost:13133"
TIMEOUT=5

check_health() {
    response=$(curl -s -w "%{http_code}" --max-time $TIMEOUT "$SREDIAG_ENDPOINT/" -o /dev/null)
    
    if [ "$response" = "200" ]; then
        echo "✅ SREDIAG is healthy"
        return 0
    else
        echo "❌ SREDIAG health check failed (HTTP $response)"
        return 1
    fi
}

check_metrics() {
    response=$(curl -s --max-time $TIMEOUT "$SREDIAG_ENDPOINT/metrics" | grep -c "otelcol_")
    
    if [ "$response" -gt 0 ]; then
        echo "✅ Metrics endpoint is working ($response metrics found)"
        return 0
    else
        echo "❌ Metrics endpoint failed"
        return 1
    fi
}

echo "Checking SREDIAG health..."
check_health && check_metrics
```

## Advanced Usage

### Custom Build Configuration

```yaml
# build-config.yaml
build:
  output_dir: "./dist"
  targets:
    - os: linux
      arch: amd64
    - os: linux
      arch: arm64
    - os: darwin
      arch: amd64
    - os: windows
      arch: amd64

components:
  receivers:
    - otlp
    - prometheus
    - jaeger
  
  processors:
    - batch
    - memory_limiter
    - k8sattributes
  
  exporters:
    - otlp
    - prometheus
    - jaeger
  
  extensions:
    - health_check
    - pprof
    - zpages
```

### Security Configuration

```yaml
# security.yaml
security:
  tls:
    cert_file: /etc/ssl/certs/srediag.crt
    key_file: /etc/ssl/private/srediag.key
    ca_file: /etc/ssl/certs/ca.crt
  
  auth:
    type: bearer_token
    token_file: /etc/srediag/auth-token
  
  rbac:
    policies_file: /etc/srediag/rbac-policies.yaml
```

This examples document provides comprehensive guidance for using SREDIAG across different scenarios and environments.