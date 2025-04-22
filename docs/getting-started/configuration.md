# SREDIAG Basic Configuration Guide

This guide covers the essential configuration settings to get started with SREDIAG.

## Table of Contents

1. [Configuration File Formats](#configuration-file-formats)
2. [Basic Configuration](#basic-configuration)
3. [Environment Variables](#environment-variables)
4. [Command Line Arguments](#command-line-arguments)
5. [Configuration Examples](#configuration-examples)

## Configuration File Formats

SREDIAG supports multiple configuration formats:

- YAML (recommended)
- JSON
- TOML

### File Locations

1. Default search paths:
   - `/etc/srediag/config.yaml`
   - `$HOME/.srediag/config.yaml`
   - `./config.yaml`

2. Custom location:

   ```bash
   srediag start --config /path/to/config.yaml
   ```

## Basic Configuration

### Minimal Configuration

```yaml
# config.yaml
version: "1.0"
service:
  name: myapp
  environment: development

telemetry:
  enabled: true
  endpoint: localhost:4317

logging:
  level: info
  format: json
  output: stdout
```

### Core Components

1. **Service Settings**

   ```yaml
   service:
     name: myapp
     environment: development
     port: 8080
     host: 0.0.0.0
   ```

2. **Telemetry Configuration**

   ```yaml
   telemetry:
     enabled: true
     endpoint: localhost:4317
     sampling:
       type: probabilistic
       ratio: 0.1
     metrics:
       enabled: true
       port: 9090
   ```

3. **Logging Settings**

   ```yaml
   logging:
     level: info
     format: json
     output: file
     file:
       path: /var/log/srediag/app.log
       max_size: 100
       max_age: 7
       max_backups: 5
   ```

## Environment Variables

Environment variables override file configurations:

```bash
# Core settings
export SREDIAG_SERVICE_NAME=myapp
export SREDIAG_SERVICE_ENV=production
export SREDIAG_PORT=9090

# Telemetry settings
export SREDIAG_TELEMETRY_ENABLED=true
export SREDIAG_TELEMETRY_ENDPOINT=collector:4317

# Logging settings
export SREDIAG_LOG_LEVEL=debug
export SREDIAG_LOG_FORMAT=json
```

## Command Line Arguments

Command line arguments take precedence:

```bash
# Basic usage
srediag start --config config.yaml

# Override settings
srediag start \
  --config config.yaml \
  --log-level debug \
  --port 9090

# Multiple config files
srediag start \
  --config base.yaml \
  --config override.yaml
```

## Configuration Examples

### Development Setup

```yaml
version: "1.0"
service:
  name: myapp
  environment: development
  port: 8080

telemetry:
  enabled: true
  endpoint: localhost:4317
  sampling:
    type: probabilistic
    ratio: 1.0
  metrics:
    enabled: true
    port: 9090

logging:
  level: debug
  format: console
  output: stdout
```

### Production Setup

```yaml
version: "1.0"
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
  metrics:
    enabled: true
    port: 9090
    push_interval: 10s

logging:
  level: info
  format: json
  output: file
  file:
    path: /var/log/srediag/app.log
    max_size: 100
    max_age: 7
    max_backups: 5
    compress: true

security:
  tls:
    enabled: true
    cert_file: /etc/srediag/certs/server.crt
    key_file: /etc/srediag/certs/server.key
  authentication:
    enabled: true
    type: oauth2
    provider: keycloak
```

### Kubernetes Setup

```yaml
version: "1.0"
service:
  name: myapp
  environment: production
  port: 8080

telemetry:
  enabled: true
  endpoint: otel-collector.monitoring:4317
  sampling:
    type: probabilistic
    ratio: 0.1
  metrics:
    enabled: true
    port: 9090

logging:
  level: info
  format: json
  output: stdout

kubernetes:
  enabled: true
  namespace: monitoring
  service_account: srediag
  pod_annotations:
    prometheus.io/scrape: "true"
    prometheus.io/port: "9090"
```

## Configuration Validation

Validate your configuration:

```bash
# Validate configuration file
srediag config validate --config config.yaml

# Test configuration
srediag config test --config config.yaml
```

## See Also

- [Advanced Configuration](../configuration/README.md)
- [Telemetry Configuration](../configuration/telemetry.md)
- [Security Configuration](../security/README.md)
- [Kubernetes Setup](../cloud/kubernetes.md)
