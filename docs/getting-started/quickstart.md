<!-- markdownlint-configure MD046 style=fenced -->
# SREDIAG Quick Start Guide

Get up and running with SREDIAG in minutes. This guide covers the essential steps to start using SREDIAG for diagnostics and monitoring.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Basic Setup](#basic-setup)
3. [First Run](#first-run)
4. [Basic Usage](#basic-usage)
5. [Next Steps](#next-steps)

## Prerequisites

Ensure you have:

- SREDIAG installed ([Installation Guide](installation.md))
- Basic understanding of monitoring concepts
- Access to a terminal/command prompt

## Basic Setup

1. **Create Configuration**

Create a file named `config.yaml`:

    ```yaml
    service:
    name: myapp
    environment: development
    port: 8080

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

2. **Set Environment Variables**

    ```bash
    # Linux/macOS
    export SREDIAG_CONFIG=/path/to/config.yaml
    export SREDIAG_LOG_LEVEL=info

    # Windows PowerShell
    $env:SREDIAG_CONFIG="C:\path\to\config.yaml"
    $env:SREDIAG_LOG_LEVEL="info"
    ```

## First Run

    1. **Start SREDIAG**

    ```bash
    # Start with default configuration
    srediag start

    # Start with custom configuration
    srediag start --config config.yaml
    ```

    2. **Verify Service**

    ```bash
    # Check health
    srediag health check

    # View metrics
    curl http://localhost:9090/metrics

    # Check API
    curl http://localhost:8080/api/v1/health
    ```

## Basic Usage

### 1. Monitor System Resources

    ```bash
    # View system metrics
    srediag metrics system

    # View process metrics
    srediag metrics process

    # View network metrics
    srediag metrics network
    ```

### 2. Collect Diagnostics

    ```bash
    # Collect basic diagnostics
    srediag collect basic

    # Collect full diagnostics
    srediag collect full --output diagnostics.zip
    ```

### 3. View Logs

    ```bash
    # View service logs
    srediag logs

    # View specific component logs
    srediag logs --component collector

    # Follow logs in real-time
    srediag logs -f
    ```

### 4. Manage Service

    ```bash
    # Start service
    srediag service start

    # Stop service
    srediag service stop

    # Restart service
    srediag service restart

    # View status
    srediag service status
    ```

## Common Operations

### 1. Update Configuration

    ```bash
    # Edit configuration
    srediag config edit

    # Validate configuration
    srediag config validate

    # Apply configuration
    srediag config apply
    ```

### 2. Manage Plugins

    ```bash
    # List available plugins
    srediag plugin list

    # Install plugin
    srediag plugin install <plugin-name>

    # Enable plugin
    srediag plugin enable <plugin-name>
    ```

### 3. View Dashboard

1. Open your browser
2. Navigate to `http://localhost:8080`
3. Log in with default credentials:
   - Username: `admin`
   - Password: `admin`
4. Change default password

## Troubleshooting

### Common Issues

    1. **Service Won't Start**

    ```bash
    # Check logs
    srediag logs --level error

    # Verify ports
    netstat -tulpn | grep srediag
    ```

    2. **Configuration Issues**

    ```bash
    # Validate config
    srediag config validate

    # Check syntax
    srediag config lint
    ```

    3. **Permission Issues**

    ```bash
    # Fix permissions
    sudo chown -R srediag:srediag /etc/srediag
    sudo chmod 755 /etc/srediag
    ```

## Next Steps

After completing this quick start guide, consider:

1. **Advanced Configuration**
   - [Telemetry Setup](../configuration/telemetry.md)
   - [Security Configuration](../configuration/security.md)
   - [Plugin Configuration](../plugins/README.md)

2. **Integration**
   - [Kubernetes Integration](../cloud/kubernetes.md)
   - [Cloud Provider Setup](../cloud/README.md)
   - [External Systems](../configuration/integrations.md)

3. **Best Practices**
   - [Security Guidelines](../security/README.md)
   - [Performance Tuning](../reference/performance.md)
   - [Production Deployment](../reference/production.md)

4. **Advanced Topics**
   - [Custom Plugins](../plugins/development.md)
   - [API Reference](../reference/api.md)
   - [Architecture Overview](../architecture/README.md)

## See Also

- [Configuration Guide](../configuration/README.md)
- [CLI Reference](../cli/README.md)
- [Troubleshooting Guide](../reference/troubleshooting.md)
- [FAQ](../reference/faq.md)
