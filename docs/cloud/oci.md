# Oracle Cloud Infrastructure (OCI) Integration

SREDIAG provides comprehensive integration with Oracle Cloud Infrastructure for monitoring, diagnostics, and analysis.

## Features

### 1. Resource Monitoring

- Compute instances
- Block volumes
- Object storage
- Database systems
- Kubernetes clusters
- Load balancers
- Virtual networks

### 2. Metric Collection

- CPU utilization
- Memory usage
- Disk I/O
- Network throughput
- Database performance
- Container metrics
- Custom metrics

### 3. Log Analysis

- Audit logs
- System logs
- Application logs
- Security logs
- Database logs
- Container logs

## Configuration

### Basic Setup

```yaml
cloud:
  oci:
    enabled: true
    auth:
      type: "instance_principal"  # or "api_key"
      # For API key authentication
      config_file: "/etc/srediag/oci/config"
      profile: "DEFAULT"
    
    region: "us-ashburn-1"
    compartment_id: "ocid1.compartment.oc1..."

    # Resource collection
    resources:
      compute:
        enabled: true
        include_tags: ["environment", "application"]
      
      storage:
        enabled: true
        volume_metrics: true
      
      database:
        enabled: true
        performance_metrics: true
      
      kubernetes:
        enabled: true
        include_nodes: true
        include_pods: true
```

### Authentication Methods

#### Instance Principal

For applications running on OCI instances:

```yaml
auth:
  type: "instance_principal"
```

#### API Key

For applications running outside OCI:

```yaml
auth:
  type: "api_key"
  config_file: "/etc/srediag/oci/config"
  profile: "DEFAULT"
```

Config file format:

```ini
[DEFAULT]
user=ocid1.user.oc1..
fingerprint=20:3b:97:13:55:1c:5b:0d:d3:37:d8:50:4e:c5:3a:34
key_file=/etc/srediag/oci/oci_api_key.pem
tenancy=ocid1.tenancy.oc1..
region=us-ashburn-1
```

## Integration Examples

### 1. Compute Instance Monitoring

```yaml
resources:
  compute:
    enabled: true
    metrics:
      - name: "cpu_utilization"
        interval: "1m"
        aggregation: "average"
      
      - name: "memory_utilization"
        interval: "1m"
        aggregation: "average"
      
      - name: "disk_iops"
        interval: "5m"
        aggregation: "sum"
    
    alarms:
      - name: "high_cpu"
        metric: "cpu_utilization"
        threshold: 80
        duration: "5m"
```

### 2. Database Monitoring

```yaml
resources:
  database:
    enabled: true
    systems:
      - type: "autonomous"
        metrics:
          - name: "cpu_utilization"
          - name: "storage_utilization"
          - name: "sessions"
      
      - type: "dbsystem"
        metrics:
          - name: "read_iops"
          - name: "write_iops"
          - name: "wait_time"
```

### 3. Kubernetes Integration

```yaml
resources:
  kubernetes:
    enabled: true
    clusters:
      - ocid: "ocid1.cluster.oc1..."
        metrics:
          - name: "kube_pod_status"
          - name: "kube_node_status"
        logs:
          enabled: true
          types: ["system", "application"]
```

## OpenTelemetry Integration

SREDIAG converts OCI metrics to OpenTelemetry format:

```yaml
exporters:
  oci:
    metric_mapping:
      "oci.compute.cpu_utilization": "system.cpu.usage"
      "oci.memory.utilization": "system.memory.usage"
    
    resource_attributes:
      cloud.provider: "oci"
      cloud.region: "${region}"
      cloud.account.id: "${tenancy}"
```

## Diagnostic Features

### 1. Resource Health Checks

```yaml
diagnostics:
  oci:
    health_checks:
      - type: "compute"
        checks:
          - "instance_status"
          - "cpu_pressure"
          - "memory_pressure"
      
      - type: "database"
        checks:
          - "connection_count"
          - "long_running_queries"
          - "tablespace_usage"
```

### 2. Performance Analysis

```yaml
analysis:
  oci:
    performance:
      compute:
        - metric: "cpu_utilization"
          baseline_period: "7d"
          anomaly_detection: true
      
      database:
        - metric: "read_iops"
          correlation:
            - "cpu_utilization"
            - "wait_time"
```

### 3. Cost Analysis

```yaml
analysis:
  oci:
    cost:
      enabled: true
      resources:
        - type: "compute"
          metrics:
            - "ocpus"
            - "memory_gb"
        - type: "storage"
          metrics:
            - "volume_size_gb"
            - "backup_size_gb"
```

## Best Practices

1. **Authentication**
   - Use instance principal when possible
   - Rotate API keys regularly
   - Use minimal required permissions

2. **Resource Collection**
   - Filter resources by compartment
   - Use appropriate metric intervals
   - Enable only needed metrics

3. **Performance**
   - Use metric aggregation
   - Implement rate limiting
   - Cache resource metadata

4. **Security**
   - Encrypt sensitive configuration
   - Use secure key storage
   - Monitor audit logs

## Troubleshooting

### Common Issues

1. **Authentication Failures**

   ```bash
   srediag diagnose oci auth --config /etc/srediag/config.yaml
   ```

2. **Metric Collection Issues**

   ```bash
   srediag analyze oci metrics --resource compute
   ```

3. **Resource Access Problems**

   ```bash
   srediag check oci permissions --compartment-id ocid1.compartment.oc1...
   ```

## Further Reading

- [OCI Documentation](https://docs.oracle.com/en-us/iaas/Content/home.htm)
- [Configuration Guide](../configuration/README.md)
- [Plugin Development](../plugins/development.md)
- [Monitoring Guide](../tutorials/monitoring.md)
