# Getting Started with SREDIAG

Welcome to SREDIAG! This guide will help you get up and running quickly with our diagnostic and monitoring solution.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Installation](#installation)
3. [Quick Start](#quick-start)
4. [Basic Configuration](#basic-configuration)
5. [First Steps](#first-steps)
6. [Next Steps](#next-steps)

## Prerequisites

Before installing SREDIAG, ensure you have:

- Go 1.21 or later
- Docker 20.10.0 or later (for containerized deployment)
- Kubernetes 1.24+ (for Kubernetes deployment)
- OpenTelemetry Collector (optional, for advanced telemetry)

## Installation

### Using Go

```bash
go install github.com/yourusername/srediag/cmd/srediag@latest
```

### Using Docker

```bash
docker pull srediag/srediag:latest
```

### Using Helm (Kubernetes)

```bash
helm repo add srediag https://charts.srediag.io
helm repo update
helm install srediag srediag/srediag
```

## Quick Start

1. **Start SREDIAG**

    ```bash
    srediag start --config srediag.yaml
    ```

2. **Verify Installation**

    ```bash
    srediag health check
    ```

3. **View Basic Metrics**

    ```bash
    curl http://localhost:9090/metrics
    ```

## Basic Configuration

Create a basic configuration file `config.yaml`:

```yaml
service:
  name: myapp
  environment: development

telemetry:
  enabled: true
  endpoint: localhost:4317
  sampling:
    type: probabilistic
    ratio: 0.1

logging:
  level: info
  format: json
  output: stdout
```

## First Steps

1. **Configure Data Collection**
   - Set up basic metrics collection
   - Configure log aggregation
   - Enable trace sampling

2. **Explore the UI**
   - Access the dashboard at `http://localhost:8080`
   - Review system metrics
   - Explore trace visualization

3. **Set Up Alerts**
   - Configure basic alerting rules
   - Set up notification channels
   - Test alert delivery

## Next Steps

After completing the basic setup, consider exploring:

1. **Advanced Features**
   - [Custom Metrics](../configuration/telemetry.md)
   - [Advanced Tracing](../configuration/collector.md)
   - [Cloud Integration](../cloud/README.md)

2. **Integration Options**
   - [Kubernetes Integration](../cloud/kubernetes.md)
   - [Cloud Provider Setup](../cloud/README.md)
   - [External Systems](../configuration/integrations.md)

3. **Best Practices**
   - [Security Guidelines](../security/README.md)
   - [Performance Tuning](../reference/performance.md)
   - [Troubleshooting](../reference/troubleshooting.md)

## See Also

- [Configuration Guide](../configuration/README.md)
- [Architecture Overview](../architecture/README.md)
- [API Reference](../reference/api.md)
- [CLI Documentation](../cli/README.md)
