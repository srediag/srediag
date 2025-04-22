# SREDIAG Troubleshooting Guide

This guide helps you diagnose and resolve common issues with SREDIAG.

## Table of Contents

1. [Common Issues](#common-issues)
2. [Diagnostic Tools](#diagnostic-tools)
3. [Log Analysis](#log-analysis)
4. [Performance Issues](#performance-issues)
5. [Known Issues](#known-issues)

## Common Issues

### Service Won't Start

1. **Check Configuration**

   ```bash
   # Validate configuration
   srediag config validate --config config.yaml
   
   # Test configuration
   srediag config test --config config.yaml
   ```

2. **Check Permissions**

   ```bash
   # Fix directory permissions
   sudo chown -R $(id -u):$(id -g) /etc/srediag
   sudo chmod 755 /etc/srediag
   ```

3. **Check Port Availability**

   ```bash
   # Check if ports are in use
   netstat -tulpn | grep -E '8080|9090'
   
   # Kill process using port
   sudo fuser -k 8080/tcp
   ```

### Telemetry Issues

1. **Collector Connection**

   ```bash
   # Check collector status
   srediag collector status
   
   # Test collector connection
   srediag collector test-connection
   ```

2. **Missing Metrics**

   ```bash
   # Check metric endpoints
   curl http://localhost:9090/metrics
   
   # View metric status
   srediag metrics status
   ```

3. **Trace Issues**

   ```bash
   # Verify trace sampling
   srediag traces check-sampling
   
   # Test trace generation
   srediag traces generate-test
   ```

### Authentication Problems

1. **Reset Admin Password**

   ```bash
   srediag user reset-password --username admin
   ```

2. **Check Auth Config**

   ```bash
   srediag auth status
   srediag auth test-config
   ```

## Diagnostic Tools

### System Checks

```bash
# Full system check
srediag diagnostics run

# Component health check
srediag health check --all

# Resource usage
srediag system stats
```

### Configuration Validation

```bash
# Validate all configs
srediag config validate --all

# Check specific component
srediag config validate --component telemetry

# Show effective config
srediag config show --effective
```

### Connectivity Tests

```bash
# Test all connections
srediag connectivity test

# Test specific endpoint
srediag connectivity test --endpoint collector:4317

# Show network status
srediag network status
```

## Log Analysis

### Viewing Logs

1. **Service Logs**

   ```bash
   # View service logs
   srediag logs show
   
   # Follow logs
   srediag logs follow
   
   # Filter logs
   srediag logs show --level error
   ```

2. **Component Logs**

   ```bash
   # View collector logs
   srediag logs show --component collector
   
   # View processor logs
   srediag logs show --component processor
   ```

### Log Levels

Adjust log levels for debugging:

```yaml
logging:
  level: debug
  components:
    collector: debug
    processor: info
    exporter: warn
```

### Log Analysis Tools

```bash
# Search logs
srediag logs search "error"

# Analyze log patterns
srediag logs analyze

# Export logs
srediag logs export --format json
```

## Performance Issues

### High CPU Usage

1. **Check Resource Usage**

   ```bash
   srediag metrics show --metric cpu_usage
   ```

2. **Profile Service**

   ```bash
   srediag profile cpu --duration 30s
   ```

3. **Optimize Configuration**

   ```yaml
   telemetry:
     batch_size: 1000
     flush_interval: 10s
   ```

### Memory Leaks

1. **Monitor Memory**

   ```bash
   srediag metrics show --metric memory_usage
   ```

2. **Generate Heap Dump**

   ```bash
   srediag debug heap-dump
   ```

3. **Analyze Memory**

   ```bash
   srediag analyze memory-usage
   ```

### Slow Performance

1. **Check Latency**

   ```bash
   srediag metrics show --metric request_latency
   ```

2. **Trace Operations**

   ```bash
   srediag trace --operation data_collection
   ```

3. **Optimize Settings**

   ```yaml
   performance:
     buffer_size: 5000
     workers: 4
   ```

## Known Issues

### Version 1.x

1. **Issue #123: High Memory Usage**
   - Symptom: Memory grows over time
   - Workaround: Restart service daily
   - Fixed in: v1.2.3

2. **Issue #456: Metric Loss**
   - Symptom: Intermittent metric drops
   - Workaround: Increase batch timeout
   - Fixed in: v1.3.0

### Version 2.x

1. **Issue #789: Slow Startup**
   - Symptom: Service takes >30s to start
   - Workaround: Disable auto-discovery
   - Status: Under investigation

## Recovery Procedures

### Service Recovery

```bash
# Safe restart
srediag service restart --graceful

# Emergency restart
srediag service restart --force

# Recover from crash
srediag service recover
```

### Data Recovery

```bash
# Backup data
srediag backup create

# Restore from backup
srediag backup restore --file backup.tar.gz

# Verify data integrity
srediag data verify
```

## See Also

- [Installation Guide](installation.md)
- [Configuration Guide](configuration.md)
- [Advanced Troubleshooting](../reference/troubleshooting.md)
- [Support Resources](../reference/support.md)
