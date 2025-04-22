# SREDIAG First Steps

This guide will walk you through your first steps using SREDIAG after installation.

## Table of Contents

1. [Initial Setup](#initial-setup)
2. [Basic Operations](#basic-operations)
3. [Data Collection](#data-collection)
4. [Visualization](#visualization)
5. [Next Steps](#next-steps)

## Initial Setup

### 1. Create Working Directory

```bash
mkdir myapp-monitoring
cd myapp-monitoring
```

### 2. Create Basic Configuration

```bash
cat > config.yaml << EOF
version: "1.0"
service:
  name: myapp
  environment: development

telemetry:
  enabled: true
  endpoint: localhost:4317
  sampling:
    type: probabilistic
    ratio: 1.0

logging:
  level: info
  format: json
  output: stdout
EOF
```

### 3. Start SREDIAG

```bash
srediag start --config config.yaml
```

## Basic Operations

### 1. Check Service Status

```bash
# Check health
srediag health check

# View service status
srediag status

# List active components
srediag components list
```

### 2. View Metrics

```bash
# View basic metrics
curl http://localhost:9090/metrics

# Check specific metric
curl http://localhost:9090/metrics | grep srediag_
```

### 3. Access Dashboard

1. Open web browser
2. Navigate to `http://localhost:8080`
3. Log in with default credentials:
   - Username: `admin`
   - Password: `admin`

## Data Collection

### 1. Enable Basic Metrics

```yaml
metrics:
  enabled: true
  collectors:
    - system
    - process
    - network
  interval: 10s
```

### 2. Configure Logging

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

### 3. Set Up Tracing

```yaml
traces:
  enabled: true
  sampler:
    type: probabilistic
    ratio: 0.1
  propagation:
    - tracecontext
    - baggage
```

## Visualization

### 1. Metrics Dashboard

1. Access the metrics dashboard:

   ```bash
   open http://localhost:8080/metrics
   ```

2. Available views:
   - System Overview
   - Resource Usage
   - Application Metrics
   - Custom Dashboards

### 2. Log Viewer

1. Access the log viewer:

   ```bash
   open http://localhost:8080/logs
   ```

2. Features:
   - Real-time log streaming
   - Log filtering
   - Pattern matching
   - Export capabilities

### 3. Trace Explorer

1. Access the trace explorer:

   ```bash
   open http://localhost:8080/traces
   ```

2. Features:
   - Trace visualization
   - Service dependency mapping
   - Performance analysis
   - Error tracking

## Next Steps

### 1. Advanced Configuration

Explore advanced configuration options:

- [Advanced Configuration Guide](../configuration/README.md)
- [Telemetry Configuration](../configuration/telemetry.md)
- [Security Setup](../security/README.md)

### 2. Integration

Set up integrations with:

- [Kubernetes](../cloud/kubernetes.md)
- [Cloud Providers](../cloud/README.md)
- [External Systems](../configuration/integrations.md)

### 3. Custom Development

Learn about extending SREDIAG:

- [API Reference](../reference/api.md)
- [Plugin Development](../plugins/README.md)
- [Custom Metrics](../configuration/metrics.md)

### 4. Best Practices

Review best practices for:

- [Performance Tuning](../reference/performance.md)
- [Security Hardening](../security/hardening.md)
- [High Availability](../reference/ha.md)

## Common Tasks

### Adding Custom Metrics

1. Define metric in configuration:

    ```yaml
    metrics:
    custom:
        - name: app_requests_total
        type: counter
        help: "Total number of HTTP requests"
        labels:
            - method
            - path
    ```

2. Use the API to record metrics:

    ```bash
    curl -X POST http://localhost:8080/api/v1/metrics \
    -H "Content-Type: application/json" \
    -d '{
        "name": "app_requests_total",
        "value": 1,
        "labels": {
        "method": "GET",
        "path": "/api/v1/users"
        }
    }'
    ```

### Setting Up Alerts

1. Create alert rule:

    ```yaml
    alerts:
    rules:
        - name: high_cpu_usage
        condition: srediag_cpu_usage > 80
        duration: 5m
        severity: warning
        annotations:
            summary: High CPU usage detected
    ```

2. Configure notification:

    ```yaml
    notifications:
    channels:
        - type: slack
        name: ops-alerts
        webhook_url: ${SLACK_WEBHOOK_URL}
    ```

### Troubleshooting

1. Enable debug logging:

    ```bash
    srediag start --config config.yaml --log-level debug
    ```

2. Check component status:

    ```bash
    srediag components status
    ```

3. View detailed diagnostics:

    ```bash
    srediag diagnostics run
    ```

## See Also

- [Configuration Guide](configuration.md)
- [Installation Guide](installation.md)
- [CLI Reference](../cli/README.md)
- [Troubleshooting Guide](../reference/troubleshooting.md)
