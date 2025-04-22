# SREDIAG OpenTelemetry Collector Configuration

## Overview

SREDIAG uses OpenTelemetry Collector as its core telemetry processing pipeline. This document
describes how to configure the collector for optimal performance and integration.

## Configuration Structure

```yaml
receivers:
  # Data ingestion points
  otlp:
    protocols:
      grpc:
        endpoint: ${OTEL_GRPC_ENDPOINT:-0.0.0.0:4317}
      http:
        endpoint: ${OTEL_HTTP_ENDPOINT:-0.0.0.0:4318}
  
  # Host metrics collection
  hostmetrics:
    collection_interval: 10s
    scrapers:
      cpu:
      disk:
      filesystem:
      memory:
      network:
      process:

processors:
  # Batch processing
  batch:
    timeout: ${BATCH_TIMEOUT:-1s}
    send_batch_size: ${BATCH_SIZE:-1024}
  
  # Memory management
  memory_limiter:
    check_interval: ${MEMORY_CHECK_INTERVAL:-1s}
    limit_mib: ${MEMORY_LIMIT_MIB:-1024}
  
  # Resource detection
  resourcedetection:
    detectors: [env, system, docker, kubernetes]
    timeout: 2s
  
  # Kubernetes attributes
  k8sattributes:
    auth_type: "serviceAccount"
    extract:
      metadata:
        - k8s.pod.name
        - k8s.pod.uid
        - k8s.deployment.name
        - k8s.namespace.name
        - k8s.node.name

exporters:
  # OTLP export
  otlp:
    endpoint: ${OTLP_ENDPOINT:-localhost:4317}
    tls:
      insecure: true
  
  # Prometheus export
  prometheus:
    endpoint: ${PROMETHEUS_ENDPOINT:-0.0.0.0:9090}

extensions:
  health_check:
  pprof:
  zpages:

service:
  extensions: [health_check, pprof, zpages]
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

## Components

### Receivers

1. **OTLP Receiver**
   - Supports gRPC and HTTP protocols
   - Configurable endpoints
   - TLS/mTLS support
   - Load balancing capabilities

2. **Host Metrics Receiver**
   - System resource monitoring
   - Configurable collection intervals
   - Multiple metric scrapers
   - Custom metric definitions

### Processors

1. **Batch Processor**
   - Optimizes data transmission
   - Configurable batch sizes
   - Timeout settings
   - Memory usage control

2. **Memory Limiter**
   - Prevents OOM conditions
   - Dynamic memory management
   - Configurable limits
   - Check intervals

3. **Resource Detection**
   - Automatic environment detection
   - Cloud provider integration
   - Container orchestration support
   - Custom attribute mapping

4. **Kubernetes Attributes**
   - Pod metadata enrichment
   - Service discovery
   - Namespace filtering
   - Label extraction

### Exporters

1. **OTLP Exporter**
   - Standard OpenTelemetry export
   - Multiple backend support
   - Compression options
   - Retry mechanisms

2. **Prometheus Exporter**
   - Metric exposition
   - Custom label mapping
   - Histogram configuration
   - Scrape endpoint setup

## Configuration Examples

### Basic Setup

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317

processors:
  batch:
    timeout: 1s
    send_batch_size: 1000

exporters:
  prometheus:
    endpoint: 0.0.0.0:9090

service:
  pipelines:
    metrics:
      receivers: [otlp]
      processors: [batch]
      exporters: [prometheus]
```

### Production Setup

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
        tls:
          cert_file: /etc/certs/server.crt
          key_file: /etc/certs/server.key
      http:
        endpoint: 0.0.0.0:4318
        tls:
          cert_file: /etc/certs/server.crt
          key_file: /etc/certs/server.key
  
  hostmetrics:
    collection_interval: 30s
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
    filter:
      node_from_env_var: KUBE_NODE_NAME

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

extensions:
  health_check:
    endpoint: 0.0.0.0:13133
  
  pprof:
    endpoint: 0.0.0.0:1777
  
  zpages:
    endpoint: 0.0.0.0:55679

service:
  extensions: [health_check, pprof, zpages]
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

## Best Practices

1. **Performance Optimization**
   - Configure appropriate batch sizes
   - Set reasonable timeouts
   - Monitor memory usage
   - Use compression when needed

2. **Security**
   - Enable TLS/mTLS
   - Use secure endpoints
   - Implement authentication
   - Follow least privilege

3. **Reliability**
   - Configure retries
   - Set up health checks
   - Monitor pipeline status
   - Implement fallbacks

4. **Scalability**
   - Use load balancing
   - Configure resource limits
   - Implement horizontal scaling
   - Monitor bottlenecks

## See Also

- [Configuration Overview](README.md)
- [Telemetry Configuration](telemetry.md)
- [Security Configuration](security.md)
- [Troubleshooting](../reference/troubleshooting.md)
