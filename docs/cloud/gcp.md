# Google Cloud Platform Integration

## Overview

SREDIAG provides comprehensive monitoring and diagnostic capabilities for Google Cloud Platform
services through deep integration with Cloud Monitoring, Cloud Trace, and OpenTelemetry.

## Supported Services

### Compute

- Compute Engine
- Google Kubernetes Engine (GKE)
- Cloud Run
- Cloud Functions
- App Engine

### Storage

- Cloud Storage
- Persistent Disk
- Filestore
- Cloud Storage for Firebase
- Transfer Service

### Database

- Cloud SQL
- Cloud Spanner
- Cloud Bigtable
- Cloud Firestore
- Cloud Memorystore

### Networking

- Virtual Private Cloud
- Cloud Load Balancing
- Cloud CDN
- Cloud Interconnect
- Cloud VPN

## Authentication

### Service Account

```bash
# Environment variable
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account.json

# GCloud CLI authentication
srediag configure gcp --login
```

### IAM Role Configuration

```yaml
gcp:
  role:
    name: SREDiagRole
    permissions:
      - monitoring.viewer
      - cloudtrace.user
      - logging.viewer
```

## Monitoring Features

### Compute Engine Monitoring

```bash
# Basic instance monitoring
srediag monitor gcp instance --project myproject --instance myinstance

# Resource utilization
srediag analyze gcp instance --project myproject --instance myinstance --metrics cpu,memory,disk

# Performance analysis
srediag diagnose gcp instance --project myproject --instance myinstance --performance
```

### Container Insights

```bash
# GKE cluster monitoring
srediag monitor gcp gke --cluster mycluster

# Pod analysis
srediag analyze gcp gke --cluster mycluster --namespace myapp

# Container metrics
srediag monitor gcp containers --cluster mycluster --metrics all
```

### Serverless Monitoring

```bash
# Cloud Function monitoring
srediag monitor gcp function --function myfunction

# Function performance analysis
srediag analyze gcp function --function myfunction --period 24h

# Cold start analysis
srediag analyze gcp function --function myfunction --cold-starts
```

## Integration with GCP Services

### Cloud Monitoring Integration

```bash
# Export metrics to Cloud Monitoring
srediag monitor gcp --export monitoring

# Import Cloud Monitoring metrics
srediag analyze gcp --import monitoring --project myproject

# Custom metrics
srediag monitor gcp --metrics-config custom-metrics.yaml
```

### Cloud Trace Integration

```bash
# Enable Cloud Trace
srediag trace gcp --enable

# Analyze traces
srediag analyze gcp traces --project myproject

# Trace visualization
srediag visualize gcp traces --output trace-map.html
```

### Cloud Logging Integration

```bash
# Query logs
srediag analyze gcp logs --project myproject --filter 'resource.type="gce_instance"'

# Export logs
srediag export gcp logs --project myproject --sink mysink

# Configure log routing
srediag configure gcp logs --project myproject --destination bigquery
```

## Configuration Examples

### Basic Configuration

```yaml
gcp:
  project_id: your-project-id
  zone: us-central1-a
  monitoring:
    interval: 60s
    metrics:
      - compute.googleapis.com/instance/cpu/utilization
      - compute.googleapis.com/instance/memory/usage
      - compute.googleapis.com/instance/disk/read_bytes_count
```

### Advanced Monitoring

```yaml
gcp:
  monitoring:
    enhanced: true
    custom_metrics:
      - name: application_errors
        type: custom.googleapis.com/application/errors
        labels:
          - service
          - environment
    alerts:
      - name: high_cpu_usage
        metric: compute.googleapis.com/instance/cpu/utilization
        threshold: 0.8
        duration: 300s
```

## OpenTelemetry Integration

### Metric Collection

```bash
# Export to OpenTelemetry Collector
srediag monitor gcp --export otlp://localhost:4317

# Configure sampling
srediag monitor gcp --sampling-ratio 0.1
```

### Trace Collection

```bash
# Enable OpenTelemetry tracing
srediag trace gcp --provider otel

# Configure trace export
srediag configure gcp traces --endpoint otel-collector:4317
```

## Best Practices

1. **Security**
   - Use service accounts with minimal permissions
   - Rotate service account keys regularly
   - Enable audit logging
   - Use Workload Identity when possible

2. **Resource Monitoring**
   - Set appropriate monitoring intervals
   - Configure meaningful alert policies
   - Use log-based metrics when needed
   - Enable enhanced monitoring for critical workloads

3. **Cost Management**
   - Monitor resource consumption
   - Set up budget alerts
   - Use labels for cost allocation
   - Configure quotas appropriately

4. **Performance Optimization**
   - Regular performance analysis
   - Resource right-sizing
   - Use autoscaling effectively
   - Monitor network latency

## Troubleshooting

### Common Issues

1. **Authentication Failures**

   ```bash
   # Verify credentials
   srediag check gcp auth

   # Test IAM permissions
   srediag verify gcp permissions
   ```

2. **Monitoring Issues**

   ```bash
   # Check metric ingestion
   srediag status gcp collector

   # Verify metric flow
   srediag verify gcp metrics
   ```

3. **Integration Problems**

   ```bash
   # Test connectivity
   srediag test gcp connectivity

   # Verify API endpoints
   srediag verify gcp endpoints
   ```

## See Also

- [Cloud Integration Overview](README.md)
- [AWS Integration](aws.md)
- [Azure Integration](azure.md)
- [OCI Integration](oci.md)
