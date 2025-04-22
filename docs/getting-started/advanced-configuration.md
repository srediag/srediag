# SREDIAG Advanced Configuration Guide

This guide covers advanced configuration options and customization for SREDIAG.

## Table of Contents

1. [Advanced Settings](#advanced-settings)
2. [Component Configuration](#component-configuration)
3. [Custom Plugins](#custom-plugins)
4. [Integration Options](#integration-options)
5. [Performance Tuning](#performance-tuning)

## Advanced Settings

### Environment Variables

```bash
# Core Settings
SREDIAG_HOME=/etc/srediag
SREDIAG_CONFIG=/etc/srediag/config.yaml
SREDIAG_LOG_LEVEL=info
SREDIAG_LOG_FORMAT=json
SREDIAG_PORT=8080

# Telemetry Settings
SREDIAG_TELEMETRY_ENABLED=true
SREDIAG_TELEMETRY_ENDPOINT=localhost:4317
SREDIAG_METRICS_PORT=9090
SREDIAG_TRACE_SAMPLE_RATIO=0.1

# Security Settings
SREDIAG_TLS_ENABLED=true
SREDIAG_TLS_CERT=/etc/srediag/certs/server.crt
SREDIAG_TLS_KEY=/etc/srediag/certs/server.key
SREDIAG_AUTH_ENABLED=true
SREDIAG_JWT_SECRET=your-secret-here
```

### Advanced Configuration File

```yaml
service:
  name: myapp
  environment: production
  version: 1.0.0
  port: 8080
  graceful_shutdown: 30s
  max_concurrent_requests: 1000

telemetry:
  enabled: true
  endpoint: collector:4317
  sampling:
    type: probabilistic
    ratio: 0.1
  metrics:
    enabled: true
    port: 9090
    path: /metrics
    push_interval: 10s
    histogram_buckets: [0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10]
  traces:
    enabled: true
    propagation:
      - tracecontext
      - baggage
      - b3
    attributes:
      environment: production
      service.version: 1.0.0
      deployment.region: us-west-2

logging:
  level: info
  format: json
  output: file
  file:
    path: /var/log/srediag/srediag.log
    max_size: 100
    max_age: 7
    max_backups: 5
    compress: true
  fields:
    environment: production
    region: us-west-2
    datacenter: dc1

security:
  tls:
    enabled: true
    cert_file: /etc/srediag/certs/server.crt
    key_file: /etc/srediag/certs/server.key
    ca_file: /etc/srediag/certs/ca.crt
    min_version: TLS1.3
  auth:
    enabled: true
    type: jwt
    jwt:
      secret: ${JWT_SECRET}
      expiration: 24h
  rbac:
    enabled: true
    default_role: viewer
    roles:
      admin:
        - "*"
      operator:
        - "read:*"
        - "write:config"
      viewer:
        - "read:config"
        - "read:metrics"

plugins:
  directory: /etc/srediag/plugins
  enabled:
    - name: system-metrics
      version: 1.0.0
      config:
        collection_interval: 10s
        include_processes: true
    - name: log-analyzer
      version: 2.1.0
      config:
        patterns_file: /etc/srediag/patterns.yaml
        max_lines: 10000
    - name: network-monitor
      version: 1.2.0
      config:
        interfaces: ["eth0", "wlan0"]
        capture_packets: false

resources:
  cpu:
    limit: 2
    request: 200m
  memory:
    limit: 2Gi
    request: 512Mi
  disk:
    storage_path: /var/lib/srediag
    min_free_space: 10Gi

monitoring:
  health_check:
    enabled: true
    port: 8081
    path: /health
    interval: 30s
    timeout: 5s
  metrics:
    retention:
      time: 15d
      size: 50Gi
    scrape_interval: 15s
    evaluation_interval: 30s

clustering:
  enabled: true
  discovery:
    method: kubernetes
    kubernetes:
      namespace: monitoring
      label_selector: app=srediag
  consensus:
    protocol: raft
    election_timeout: 5s
    heartbeat_interval: 1s
  replication:
    factor: 3
    sync_timeout: 10s

api:
  cors:
    enabled: true
    allowed_origins: ["https://dashboard.example.com"]
    allowed_methods: ["GET", "POST", "PUT", "DELETE"]
    allowed_headers: ["Authorization", "Content-Type"]
    max_age: 86400
  rate_limit:
    enabled: true
    requests_per_second: 100
    burst: 200
```

## Component Configuration

### 1. Collector Configuration

```yaml
collector:
  receivers:
    otlp:
      protocols:
        grpc:
          endpoint: 0.0.0.0:4317
        http:
          endpoint: 0.0.0.0:4318
    
    hostmetrics:
      collection_interval: 10s
      scrapers:
        cpu:
          metrics:
            system.cpu.time:
              enabled: true
            system.cpu.utilization:
              enabled: true
        memory:
          metrics:
            system.memory.usage:
              enabled: true
            system.memory.utilization:
              enabled: true

  processors:
    batch:
      timeout: 5s
      send_batch_size: 2048
    
    memory_limiter:
      check_interval: 1s
      limit_mib: 2048
    
    resourcedetection:
      detectors: [env, system, docker, kubernetes]
      timeout: 2s
      override: true
    
    k8sattributes:
      auth_type: serviceAccount
      extract:
        metadata:
          - k8s.pod.name
          - k8s.pod.uid
          - k8s.deployment.name
          - k8s.namespace.name
          - k8s.node.name

  exporters:
    otlp:
      endpoint: collector:4317
      tls:
        ca_file: /etc/certs/ca.crt
        cert_file: /etc/certs/client.crt
        key_file: /etc/certs/client.key
    
    prometheus:
      endpoint: 0.0.0.0:9090
      namespace: srediag
      const_labels:
        environment: production
        deployment: kubernetes

  service:
    telemetry:
      logs:
        level: info
        development: false
        encoding: json
      metrics:
        level: detailed
        address: 0.0.0.0:8888
    pipelines:
      traces:
        receivers: [otlp]
        processors: [memory_limiter, k8sattributes, resourcedetection, batch]
        exporters: [otlp]
      metrics:
        receivers: [otlp, hostmetrics]
        processors: [memory_limiter, k8sattributes, resourcedetection, batch]
        exporters: [otlp, prometheus]
```

### 2. Plugin Configuration

```yaml
plugins:
  system_metrics:
    collection:
      interval: 10s
      timeout: 5s
    metrics:
      cpu:
        enabled: true
        include_per_cpu: true
      memory:
        enabled: true
        include_swap: true
      disk:
        enabled: true
        include_partitions: true
        ignore_fs_types: [tmpfs, devtmpfs]
      network:
        enabled: true
        include_per_interface: true
    alerting:
      rules:
        - name: high_cpu_usage
          condition: cpu_usage > 80
          duration: 5m
          severity: warning
        - name: memory_pressure
          condition: memory_usage > 90
          duration: 10m
          severity: critical

  log_analyzer:
    patterns:
      - name: error_pattern
        regex: "ERROR|FATAL|CRITICAL"
        severity: error
      - name: warning_pattern
        regex: "WARN|WARNING"
        severity: warning
    aggregation:
      interval: 5m
      group_by: [severity, component]
    retention:
      time: 7d
      max_size: 10Gi

  network_monitor:
    interfaces:
      - name: eth0
        metrics:
          - throughput
          - packets
          - errors
      - name: wlan0
        metrics:
          - signal_strength
          - quality
    capture:
      enabled: false
      max_packet_size: 128
      buffer_size: 1Mi
```

## Custom Plugins

### 1. Plugin Structure

```go
package main

import (
    "context"
    "github.com/yourusername/srediag/pkg/plugin"
)

type MyPlugin struct {
    config *Config
    logger *zap.Logger
}

type Config struct {
    Interval    time.Duration `yaml:"interval"`
    MetricName  string        `yaml:"metric_name"`
    Labels      []string      `yaml:"labels"`
}

func (p *MyPlugin) Start(ctx context.Context) error {
    // Plugin initialization code
    return nil
}

func (p *MyPlugin) Stop(ctx context.Context) error {
    // Cleanup code
    return nil
}

func (p *MyPlugin) Collect(ctx context.Context) ([]plugin.Metric, error) {
    // Metric collection code
    return metrics, nil
}
```

### 2. Plugin Configuration of Diagnostic Plugins

```yaml
plugins:
  my_plugin:
    enabled: true
    interval: 30s
    metric_name: custom_metric
    labels:
      - environment
      - component
```

## Integration Options

### 1. Kubernetes Integration

```yaml
kubernetes:
  enabled: true
  namespace: monitoring
  service_account: srediag
  resources:
    requests:
      cpu: 200m
      memory: 256Mi
    limits:
      cpu: 1000m
      memory: 1Gi
  node_selector:
    role: monitoring
  tolerations:
    - key: node-role.kubernetes.io/master
      operator: Exists
      effect: NoSchedule
```

### 2. Cloud Provider Integration

```yaml
cloud:
  provider: aws
  region: us-west-2
  credentials:
    access_key_id: ${AWS_ACCESS_KEY_ID}
    secret_access_key: ${AWS_SECRET_ACCESS_KEY}
  services:
    cloudwatch:
      enabled: true
      log_group: /srediag/logs
      metrics_namespace: SREDIAG/Metrics
    s3:
      enabled: true
      bucket: srediag-diagnostics
      prefix: logs/
```

## Performance Tuning

### 1. Resource Limits

```yaml
resources:
  cpu:
    limit: 2
    request: 200m
  memory:
    limit: 2Gi
    request: 512Mi
  storage:
    size: 50Gi
    class: standard
```

### 2. Buffer Configuration

```yaml
buffers:
  traces:
    size: 1000000
    batch_size: 1000
    flush_interval: 5s
  metrics:
    size: 500000
    batch_size: 5000
    flush_interval: 10s
  logs:
    size: 100000
    batch_size: 1000
    flush_interval: 1s
```

### 3. Network Tuning

```yaml
network:
  tcp_keep_alive: true
  max_connections: 10000
  read_buffer_size: 4096
  write_buffer_size: 4096
  dial_timeout: 5s
  tls_handshake_timeout: 10s
```

## See Also

- [Configuration Overview](../configuration/README.md)
- [Plugin Development](../plugins/development.md)
- [Performance Guide](../reference/performance.md)
- [Security Configuration](../configuration/security.md)
