# SREDIAG Configuration Guide

SREDIAG uses a flexible configuration system that supports multiple formats and sources,
with OpenTelemetry integration at its core.

## Table of Contents

1. [Overview](#overview)
2. [Configuration Files](#configuration-files)
3. [Environment Variables](#environment-variables)
4. [Command Line Arguments](#command-line-arguments)
5. [OpenTelemetry Configuration](#opentelemetry-configuration)
6. [Examples](#examples)
7. [Best Practices](#best-practices)

## Overview

SREDIAG's configuration system is built on the following principles:

- Hierarchical configuration with inheritance
- Multiple format support (YAML, JSON, TOML)
- Environment variable overrides
- Dynamic configuration reloading
- OpenTelemetry integration

## Configuration Files

### Main Configuration

The main configuration file (`srediag.yaml`) defines core settings:

```yaml
# SREDIAG Configuration
version: "v0.1.0"

# Service settings
service:
  name: "srediag"
  environment: "${SREDIAG_ENV:-production}"
  port: "${SREDIAG_PORT:-8080}"

# Collector settings
collector:
  enabled: true
  config_path: "configs/otel-config.yaml"
  memory_limit_mib: 1024
  cpu_limit_cores: 2

# Application logging
logging:
  level: "${LOG_LEVEL:-info}"
  format: "${LOG_FORMAT:-json}"
  file:
    enabled: "${LOG_FILE_ENABLED:-false}"
    path: "${LOG_FILE_PATH:-/var/log/srediag/srediag.log}"

```

### Component Configuration

Component-specific configurations are stored in separate files:

```yaml
# otel-config.yaml
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

## Environment Variables

Environment variables override file configurations:

```bash
# Core settings
export SREDIAG_SERVICE_NAME=myservice
export SREDIAG_SERVICE_ENV=staging
export SREDIAG_PORT=9090

# Telemetry settings
export SREDIAG_TELEMETRY_ENABLED=true
export SREDIAG_TELEMETRY_ENDPOINT=collector:4317

# Logging settings
export SREDIAG_LOG_LEVEL=debug
export SREDIAG_LOG_FORMAT=json
```

## Command Line Arguments

Command line arguments take precedence over other settings:

```bash
# Basic configuration
srediag --config /etc/srediag/config.yaml

# Override specific settings
srediag --log-level debug --port 9090

# Multiple config files
srediag --config base.yaml --config override.yaml
```

## OpenTelemetry Configuration

### Collector Configuration

```yaml
# otel-collector.yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: ${OTEL_GRPC_ENDPOINT:-0.0.0.0:4317}
      http:
        endpoint: ${OTEL_HTTP_ENDPOINT:-0.0.0.0:4318}

processors:
  batch:
    timeout: ${BATCH_TIMEOUT:-1s}
    send_batch_size: ${BATCH_SIZE:-1024}

  memory_limiter:
    check_interval: ${MEMORY_CHECK_INTERVAL:-1s}
    limit_mib: ${MEMORY_LIMIT_MIB:-1024}

exporters:
  otlp:
    endpoint: ${OTLP_ENDPOINT:-localhost:4317}
    tls:
      insecure: true

  prometheus:
    endpoint: ${PROMETHEUS_ENDPOINT:-0.0.0.0:9090}

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [memory_limiter, batch]
      exporters: [otlp]
    metrics:
      receivers: [otlp]
      processors: [memory_limiter, batch]
      exporters: [otlp, prometheus]
```

### Resource Detection

```yaml
# resource-detection.yaml
processors:
  resourcedetection:
    detectors: [env, system, docker, kubernetes]
    timeout: 2s
    override: true
    attributes:
      - key: service.name
        value: ${SERVICE_NAME}
      - key: deployment.environment
        value: ${DEPLOYMENT_ENV}
```

## Examples

### Basic Service Configuration

```yaml
# basic-config.yaml
service:
  name: myapp
  environment: production
  port: 8080

telemetry:
  enabled: true
  endpoint: collector:4317
  sampling:
    type: probabilistic
    ratio: 0.1

logging:
  level: info
  format: json
  file:
    enabled: true
    path: /var/log/srediag/service.log
```

### Advanced Configuration

```yaml
# advanced-config.yaml
service:
  name: myapp
  environment: production
  port: 8080
  features:
    - monitoring
    - diagnostics
    - profiling

telemetry:
  enabled: true
  endpoint: collector:4317
  sampling:
    type: probabilistic
    ratio: 0.1
  metrics:
    interval: 10s
    prefix: myapp
  traces:
    enabled: true
    propagation: b3

logging:
  level: info
  format: json
  output:
    - stdout
    - file
  file:
    enabled: true
    path: /var/log/srediag/service.log
    rotation:
      max_size: 100MB
      max_age: 7d
      max_backups: 5
```

## Best Practices

1. **Configuration Organization**
   - Use separate files for different components
   - Keep sensitive data in secure storage
   - Use environment-specific overrides
   - Version control configurations

2. **Security**
   - Never commit secrets to version control
   - Use environment variables for sensitive data
   - Implement proper access controls
   - Regularly rotate credentials

3. **Monitoring**
   - Configure appropriate logging levels
   - Set meaningful metric intervals
   - Use proper sampling rates
   - Monitor configuration changes

4. **Performance**
   - Optimize batch sizes
   - Configure appropriate intervals
   - Use resource limits
   - Monitor impact of configuration

## See Also

- [CLI Documentation](../cli/README.md)
- [Cloud Integration](../cloud/README.md)
- [Security Guide](../security/README.md)
- [Troubleshooting](../reference/troubleshooting.md)
