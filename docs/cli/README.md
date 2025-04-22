# SREDIAG CLI

SREDIAG implements a comprehensive security model to protect system resources, data, and communications.

SREDIAG CLI provides powerful diagnostic and analysis tools built on top of OpenTelemetry
for system troubleshooting, monitoring, and management.

## Table of Contents

1. [Overview](#overview)
2. [Command Categories](#command-categories)
3. [Global Options](#global-options)
4. [Output Formats](#output-formats)
5. [Environment Variables](#environment-variables)
6. [Exit Codes](#exit-codes)
7. [OpenTelemetry Integration](#opentelemetry-integration)

## Overview

SREDIAG CLI is organized into several command categories, each focusing on specific
diagnostic and management tasks. All commands follow a consistent pattern:

```bash
srediag <category> <command> [options]
```

## Command Categories

### System Diagnostics

```bash
# Basic system diagnostics
srediag diagnose system
srediag diagnose system --resource cpu
srediag diagnose system --resource memory
srediag diagnose system --resource disk

# Real-time monitoring
srediag monitor system --interval 5s

# Process analysis
srediag analyze process --pid 1234
```

### Kubernetes Diagnostics

```bash
# Cluster health check
srediag diagnose kubernetes --cluster my-cluster

# Resource analysis
srediag analyze kubernetes resources --namespace default
srediag analyze kubernetes deployments --show-events

# Configuration audit
srediag audit kubernetes --namespace production
```

### Performance Analysis

```bash
# CPU profiling
srediag profile cpu --duration 30s --output cpu.prof

# Memory analysis
srediag analyze memory --threshold 90

# Bottleneck detection
srediag analyze bottlenecks --service my-service
```

### Security Diagnostics

```bash
# Vulnerability scanning
srediag scan vulnerabilities --severity high

# Compliance checking
srediag check compliance --standard pci-dss

# Configuration audit
srediag audit config --path /etc/srediag
```

## Global Options

All commands support these options:

| Option      | Description                    | Example                        |
|-------------|--------------------------------|--------------------------------|
| --config    | Configuration file path        | --config=/etc/srediag.yaml     |
| --format    | Output format (json/yaml/table)| --format=json                  |
| --verbose   | Enable verbose logging         | --verbose                      |
| --quiet     | Suppress non-error output      | --quiet                        |
| --output    | Output file path               | --output=results.json          |

## Output Formats

### JSON Format

```bash
srediag diagnose system --format json
```

```json
{
  "status": "healthy",
  "metrics": {
    "cpu_usage": 45.2,
    "memory_usage": 68.7
  }
}
```

### YAML Format

```bash
srediag analyze config --format yaml
```

```yaml
status: healthy
metrics:
  cpu_usage: 45.2
  memory_usage: 68.7
```

### Table Format

```bash
srediag monitor resources --format table
```

```text
+----------+--------+---------+-----------+
| Resource | Usage  | Status  | Threshold |
+----------+--------+---------+-----------+
| CPU      | 45.2%  | OK      | 80%       |
| Memory   | 68.7%  | Warning | 70%       |
+----------+--------+---------+-----------+
```

## Environment Variables

| Variable              | Description           | Example                        |
|-----------------------|-----------------------|--------------------------------|
| SREDIAG_CONFIG        | Configuration path    | /etc/srediag/config.yaml       |
| SREDIAG_OUTPUT_FORMAT | Default output format | json                           |
| SREDIAG_LOG_LEVEL     | Logging level         | debug                          |
| SREDIAG_API_KEY       | API key for cloud     | sk_live_123...                 |

## Exit Codes

| Code | Description        |
|------|--------------------|
| 0    | Success            |
| 1    | General error      |
| 2    | Config error       |
| 3    | Permission error   |
| 4    | Resource not found |
| 5    | Timeout error      |

## OpenTelemetry Integration

SREDIAG integrates seamlessly with OpenTelemetry:

```bash
# Export to OpenTelemetry Collector
srediag diagnose system --export otlp://localhost:4317

# Use OpenTelemetry context
srediag monitor --context otel-context.txt

# Configure sampling
srediag analyze --sampling-ratio 0.1
```

## See Also

- [Configuration Guide](../configuration/README.md)
- [Plugin Development](../plugins/development.md)
- [Troubleshooting](../reference/troubleshooting.md)
- [Security Guide](../security/README.md)
