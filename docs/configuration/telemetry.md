# SREDIAG Telemetry Configuration

## Overview

SREDIAG uses OpenTelemetry as the foundation for all instrumentation and telemetry collection. This document describes how to configure telemetry for your environment.

## Configuration Structure

```yaml
telemetry:
  service_name: "srediag"
  service_version: "1.0.0"
  
  # Logs Configuration
  logs:
    level: ${LOG_LEVEL:-info}
    format: ${LOG_FORMAT:-json}
    output: ${LOG_OUTPUT:-stdout}
    file:
      path: ${LOG_FILE:-/var/log/srediag/srediag.log}
      max_size: ${LOG_MAX_SIZE:-100}  # MB
      max_age: ${LOG_MAX_AGE:-7}      # days
      max_backups: ${LOG_MAX_BACKUPS:-5}
      compress: true
  
  # Metrics Configuration
  metrics:
    enabled: ${METRICS_ENABLED:-true}
    host: ${METRICS_HOST:-0.0.0.0}
    port: ${METRICS_PORT:-9090}
    path: ${METRICS_PATH:-/metrics}
    push_interval: ${METRICS_PUSH_INTERVAL:-10s}
    histogram_buckets: [0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10]
  
  # Traces Configuration
  traces:
    enabled: ${TRACES_ENABLED:-true}
    sampler:
      type: ${TRACE_SAMPLER_TYPE:-parentbased_traceidratio}
      ratio: ${TRACE_SAMPLER_RATIO:-0.1}
    propagation:
      - tracecontext
      - baggage
      - b3
    attributes:
      environment: ${ENV:-production}
      deployment: ${DEPLOYMENT:-kubernetes}
  
  # Exporter Configuration
  exporter:
    type: ${EXPORTER_TYPE:-otlp}
    endpoint: ${EXPORTER_ENDPOINT:-localhost:4317}
    headers:
      Authorization: ${EXPORTER_AUTH_TOKEN:-}
    compression: ${EXPORTER_COMPRESSION:-gzip}
    timeout: ${EXPORTER_TIMEOUT:-5s}
    retry:
      enabled: true
      initial_interval: 5s
      max_interval: 30s
      max_elapsed_time: 300s
```

## Components

### Logs

1. **Log Levels**
   - `debug`: Detailed information for development
   - `info`: General operational information
   - `warn`: Unexpected but non-critical situations
   - `error`: Errors affecting functionality
   - `fatal`: Critical errors causing termination

2. **Log Formats**
   - `json`: Structured format for parsing
   - `text`: Human-readable format
   - `console`: Colored format for development

3. **Log Rotation**
   - Based on maximum size
   - Based on maximum age
   - Automatic compression
   - Configurable retention

### Metrics

1. **Metric Types**
   - Counters
   - Gauges
   - Histograms
   - Summary metrics

2. **Endpoints**
   - HTTP exposure
   - Prometheus format
   - Port/path configuration

3. **Histogram Buckets**
   - Custom configuration
   - Latency optimization
   - Logarithmic distribution

### Traces

1. **Sampling**
   - Trace ID based
   - Parent based
   - Configurable rate
   - Custom rules

2. **Propagation**
   - W3C TraceContext
   - W3C Baggage
   - B3 (single/multi)
   - Jaeger

3. **Attributes**
   - Environment
   - Deployment
   - Version
   - Customizable

## Configuration Examples

### Basic Configuration

```yaml
telemetry:
  service_name: "srediag"
  logs:
    level: info
    format: json
    output: stdout
  
  metrics:
    enabled: true
    port: 9090
  
  traces:
    enabled: true
    sampler:
      type: parentbased_traceidratio
      ratio: 0.1
  
  exporter:
    type: otlp
    endpoint: localhost:4317
```

### Production Configuration

```yaml
telemetry:
  service_name: "srediag"
  service_version: "1.2.3"
  
  logs:
    level: info
    format: json
    output: file
    file:
      path: /var/log/srediag/srediag.log
      max_size: 100
      max_age: 7
      max_backups: 5
      compress: true
  
  metrics:
    enabled: true
    host: 0.0.0.0
    port: 9090
    path: /metrics
    push_interval: 10s
    histogram_buckets: [0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10]
  
  traces:
    enabled: true
    sampler:
      type: parentbased_traceidratio
      ratio: 0.1
    propagation:
      - tracecontext
      - baggage
    attributes:
      environment: production
      deployment: kubernetes
      region: us-west-2
      datacenter: dc1
  
  exporter:
    type: otlp
    endpoint: collector:4317
    headers:
      Authorization: "Bearer ${OTEL_AUTH_TOKEN}"
    compression: gzip
    timeout: 5s
    retry:
      enabled: true
      initial_interval: 5s
      max_interval: 30s
      max_elapsed_time: 300s
```

## Best Practices

1. **Log Management**
   - Use appropriate levels
   - Configure rotation
   - Implement structuring
   - Monitor disk usage

2. **Metrics Optimization**
   - Define relevant buckets
   - Limit cardinality
   - Configure intervals
   - Monitor memory usage

3. **Trace Configuration**
   - Adjust sampling rate
   - Define useful attributes
   - Configure propagation
   - Implement context

4. **Exporter**
   - Configure retry
   - Use compression
   - Set timeouts
   - Monitor backpressure

## See Also

- [Configuration Overview](README.md)
- [Collector Configuration](collector.md)
- [Security Configuration](security.md)
- [Troubleshooting](../reference/troubleshooting.md)
