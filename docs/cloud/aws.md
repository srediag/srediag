# AWS Integration

## Overview

SREDIAG provides comprehensive monitoring and diagnostic capabilities for AWS services
through deep integration with AWS APIs and OpenTelemetry.

## Supported Services

### Compute

- EC2 instances
- ECS containers
- EKS clusters
- Lambda functions
- Auto Scaling groups

### Storage

- S3 buckets
- EBS volumes
- EFS file systems
- FSx file systems

### Database

- RDS instances
- DynamoDB tables
- ElastiCache clusters
- Aurora clusters

### Networking

- VPCs
- Load Balancers
- Transit Gateways
- Direct Connect
- VPN connections

## Authentication

### IAM Credentials

```bash
# Environment variables
export AWS_ACCESS_KEY_ID=your_access_key
export AWS_SECRET_ACCESS_KEY=your_secret_key
export AWS_REGION=us-west-2

# AWS CLI profile
srediag configure aws --profile myprofile
```

### IAM Role Configuration

```yaml
aws:
  role:
    name: SREDiagRole
    policies:
      - AWSCloudWatchReadOnlyAccess
      - AWSXRayReadOnlyAccess
      - AWSSystemsManagerReadOnlyAccess
```

## Monitoring Features

### EC2 Monitoring

```bash
# Basic instance monitoring
srediag monitor aws ec2 --instance-id i-1234567890abcdef0

# Resource utilization
srediag analyze aws ec2 --instance-id i-1234567890abcdef0 --metrics cpu,memory,disk

# Performance analysis
srediag diagnose aws ec2 --instance-id i-1234567890abcdef0 --performance
```

### Container Insights

```bash
# ECS cluster monitoring
srediag monitor aws ecs --cluster mycluster

# EKS pod analysis
srediag analyze aws eks --cluster mycluster --namespace myapp

# Container metrics
srediag monitor aws containers --cluster mycluster --metrics all
```

### Serverless Monitoring

```bash
# Lambda function monitoring
srediag monitor aws lambda --function myfunction

# Lambda performance analysis
srediag analyze aws lambda --function myfunction --period 24h

# Cold start analysis
srediag analyze aws lambda --function myfunction --cold-starts
```

## Integration with AWS Services

### CloudWatch Integration

```bash
# Export metrics to CloudWatch
srediag monitor aws --export cloudwatch

# Import CloudWatch metrics
srediag analyze aws --import cloudwatch --namespace AWS/EC2

# Custom metrics
srediag monitor aws --metrics-config custom-metrics.yaml
```

### X-Ray Integration

```bash
# Enable X-Ray tracing
srediag trace aws --enable

# Analyze traces
srediag analyze aws traces --service myapp

# Trace visualization
srediag visualize aws traces --output trace-map.html
```

### Systems Manager Integration

```bash
# Run diagnostics using SSM
srediag diagnose aws ssm --instance-id i-1234567890abcdef0

# Collect system information
srediag collect aws ssm --instance-id i-1234567890abcdef0 --type system-info

# Execute automation
srediag automate aws ssm --document AWS-RunPatchBaseline
```

## Configuration Examples

### Basic Configuration

```yaml
aws:
  region: us-west-2
  credentials:
    profile: default
  monitoring:
    interval: 60s
    metrics:
      - CPUUtilization
      - MemoryUtilization
      - DiskUsage
```

### Advanced Monitoring

```yaml
aws:
  monitoring:
    enhanced: true
    custom_metrics:
      - name: ApplicationErrors
        namespace: Custom/Application
        dimensions:
          - Service
          - Environment
    alarms:
      - name: HighCPUUsage
        metric: CPUUtilization
        threshold: 80
        period: 300
```

## OpenTelemetry Integration

### Metric Collection

```bash
# Export to OpenTelemetry Collector
srediag monitor aws --export otlp://localhost:4317

# Configure sampling
srediag monitor aws --sampling-ratio 0.1
```

### Trace Collection

```bash
# Enable OpenTelemetry tracing
srediag trace aws --provider otel

# Configure trace export
srediag configure aws traces --endpoint otel-collector:4317
```

## Best Practices

1. **IAM Security**
   - Use least privilege access
   - Rotate credentials regularly
   - Implement MFA for API access

2. **Resource Monitoring**
   - Set appropriate monitoring intervals
   - Configure meaningful alarms
   - Use enhanced monitoring for critical resources

3. **Cost Management**
   - Monitor resource utilization
   - Set up cost allocation tags
   - Configure budget alerts

4. **Performance Optimization**
   - Regular performance analysis
   - Resource right-sizing
   - Automated scaling configuration

## Troubleshooting

### Common Issues

1. **Authentication Failures**

   ```bash
   # Verify credentials
   srediag check aws auth

   # Test IAM permissions
   srediag verify aws permissions
   ```

2. **Monitoring Issues**

   ```bash
   # Check collector status
   srediag status aws collector

   # Verify metric flow
   srediag verify aws metrics
   ```

3. **Integration Problems**

   ```bash
   # Test connectivity
   srediag test aws connectivity

   # Verify service endpoints
   srediag verify aws endpoints
   ```

## See Also

- [Cloud Integration Overview](README.md)
- [Azure Integration](azure.md)
- [GCP Integration](gcp.md)
- [OCI Integration](oci.md)
