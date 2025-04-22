# Azure Integration

## Overview

SREDIAG provides comprehensive monitoring and diagnostic capabilities for Azure services
through deep integration with Azure Monitor, Application Insights, and OpenTelemetry.

## Supported Services

### Compute

- Virtual Machines
- Virtual Machine Scale Sets
- Azure Kubernetes Service (AKS)
- Azure Functions
- App Services

### Storage

- Azure Storage Accounts
- Azure Files
- Azure Disks
- Azure NetApp Files
- Azure Backup

### Database

- Azure SQL Database
- Azure Cosmos DB
- Azure Database for MySQL
- Azure Database for PostgreSQL
- Azure Cache for Redis

### Networking

- Virtual Networks
- Load Balancers
- Application Gateway
- ExpressRoute
- Azure Firewall

## Authentication

### Service Principal

```bash
# Environment variables
export AZURE_SUBSCRIPTION_ID=your_subscription_id
export AZURE_TENANT_ID=your_tenant_id
export AZURE_CLIENT_ID=your_client_id
export AZURE_CLIENT_SECRET=your_client_secret

# Azure CLI login
srediag configure azure --login
```

### Role Configuration

```yaml
azure:
  role:
    name: SREDiagRole
    assignments:
      - Monitoring Reader
      - Metrics Reader
      - Log Analytics Reader
```

## Monitoring Features

### Virtual Machine Monitoring

```bash
# Basic VM monitoring
srediag monitor azure vm --resource-group mygroup --name myvm

# Resource utilization
srediag analyze azure vm --resource-group mygroup --name myvm --metrics cpu,memory,disk

# Performance analysis
srediag diagnose azure vm --resource-group mygroup --name myvm --performance
```

### Container Insights

```bash
# AKS cluster monitoring
srediag monitor azure aks --cluster mycluster

# Pod analysis
srediag analyze azure aks --cluster mycluster --namespace myapp

# Container metrics
srediag monitor azure containers --cluster mycluster --metrics all
```

### Serverless Monitoring

```bash
# Function App monitoring
srediag monitor azure function --app myapp --function myfunction

# Function performance analysis
srediag analyze azure function --app myapp --function myfunction --period 24h

# Cold start analysis
srediag analyze azure function --app myapp --function myfunction --cold-starts
```

## Integration with Azure Services

### Azure Monitor Integration

```bash
# Export metrics to Azure Monitor
srediag monitor azure --export monitor

# Import Azure Monitor metrics
srediag analyze azure --import monitor --namespace Microsoft.Compute

# Custom metrics
srediag monitor azure --metrics-config custom-metrics.yaml
```

### Application Insights Integration

```bash
# Enable Application Insights
srediag trace azure --enable appinsights

# Analyze traces
srediag analyze azure traces --app myapp

# Trace visualization
srediag visualize azure traces --output trace-map.html
```

### Log Analytics Integration

```bash
# Query logs
srediag analyze azure logs --workspace myworkspace --query mykusto

# Export logs
srediag export azure logs --workspace myworkspace --timespan 24h

# Configure log collection
srediag configure azure logs --workspace myworkspace --sources vm,aks
```

## Configuration Examples

### Basic Configuration

```yaml
azure:
  subscription_id: your_subscription_id
  resource_group: your_resource_group
  monitoring:
    interval: 60s
    metrics:
      - Percentage CPU
      - Available Memory Bytes
      - Disk Read Bytes
```

### Advanced Monitoring

```yaml
azure:
  monitoring:
    enhanced: true
    custom_metrics:
      - name: ApplicationErrors
        namespace: Custom/Application
        dimensions:
          - Service
          - Environment
    alerts:
      - name: HighCPUUsage
        metric: Percentage CPU
        threshold: 80
        window: 5m
```

## OpenTelemetry Integration

### Metric Collection

```bash
# Export to OpenTelemetry Collector
srediag monitor azure --export otlp://localhost:4317

# Configure sampling
srediag monitor azure --sampling-ratio 0.1
```

### Trace Collection

```bash
# Enable OpenTelemetry tracing
srediag trace azure --provider otel

# Configure trace export
srediag configure azure traces --endpoint otel-collector:4317
```

## Best Practices

1. **Security**
   - Use Managed Identities when possible
   - Implement least privilege access
   - Regular credential rotation
   - Enable Azure AD authentication

2. **Resource Monitoring**
   - Configure appropriate metric collection intervals
   - Set up meaningful alert rules
   - Use diagnostic settings effectively
   - Enable enhanced monitoring for critical resources

3. **Cost Management**
   - Monitor resource consumption
   - Set up cost allocation tags
   - Configure budget alerts
   - Use auto-scaling appropriately

4. **Performance Optimization**
   - Regular performance assessment
   - Resource right-sizing
   - Implement auto-scaling
   - Monitor application dependencies

## Troubleshooting

### Common Issues

1. **Authentication Failures**

   ```bash
   # Verify credentials
   srediag check azure auth

   # Test role assignments
   srediag verify azure permissions
   ```

2. **Monitoring Issues**

   ```bash
   # Check data collection
   srediag status azure collector

   # Verify metric flow
   srediag verify azure metrics
   ```

3. **Integration Problems**

   ```bash
   # Test connectivity
   srediag test azure connectivity

   # Verify service endpoints
   srediag verify azure endpoints
   ```

## See Also

- [Cloud Integration Overview](README.md)
- [AWS Integration](aws.md)
- [GCP Integration](gcp.md)
- [OCI Integration](oci.md)
