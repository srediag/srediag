# Cloud Integration

SREDIAG Cloud provides seamless integration with major cloud providers for comprehensive
cloud infrastructure diagnostics, monitoring, and management.

## Table of Contents

1. [Overview](#overview)
2. [Supported Providers](#supported-providers)
3. [Authentication](#authentication)
4. [Features](#features)
5. [Configuration](#configuration)
6. [Usage Examples](#usage-examples)
7. [OpenTelemetry Integration](#opentelemetry-integration)

## Overview

SREDIAG Cloud integration enables diagnostic and monitoring capabilities across various
cloud platforms, providing unified insights into your cloud infrastructure.

## Supported Providers

### Amazon Web Services (AWS)

- EC2 instance monitoring
- ECS container insights
- EKS cluster diagnostics
- CloudWatch metrics integration
- AWS Lambda monitoring

### Microsoft Azure

- Azure VM diagnostics
- AKS cluster monitoring
- Azure Container Instances
- Azure Monitor integration
- Azure Functions monitoring

### Google Cloud Platform (GCP)

- GCE instance monitoring
- GKE cluster diagnostics
- Cloud Run container insights
- Cloud Monitoring integration
- Cloud Functions monitoring

### Oracle Cloud Infrastructure (OCI)

- Compute instance monitoring
- Container Engine (OKE) diagnostics
- Cloud Native monitoring
- OCI monitoring integration

## Authentication

Each cloud provider requires specific authentication methods:

```bash
# AWS authentication
export AWS_ACCESS_KEY_ID=your_access_key
export AWS_SECRET_ACCESS_KEY=your_secret_key
export AWS_REGION=us-west-2

# Azure authentication
export AZURE_SUBSCRIPTION_ID=your_subscription_id
export AZURE_TENANT_ID=your_tenant_id
export AZURE_CLIENT_ID=your_client_id
export AZURE_CLIENT_SECRET=your_client_secret

# GCP authentication
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/credentials.json

# OCI authentication
export OCI_CONFIG_FILE=/path/to/config
export OCI_PROFILE=DEFAULT
```

## Features

### Infrastructure Monitoring

```bash
# AWS EC2 monitoring
srediag monitor aws ec2 --instance-id i-1234567890abcdef0

# Azure VM monitoring
srediag monitor azure vm --resource-group mygroup --name myvm

# GCP instance monitoring
srediag monitor gcp instance --project myproject --instance myinstance

# OCI instance monitoring
srediag monitor oci instance --compartment mycompartment --instance myinstance
```

### Container Orchestration

```bash
# EKS cluster diagnostics
srediag diagnose aws eks --cluster mycluster

# AKS cluster monitoring
srediag monitor azure aks --cluster mycluster

# GKE cluster analysis
srediag analyze gcp gke --cluster mycluster

# OKE cluster health
srediag check oci oke --cluster mycluster
```

### Serverless Monitoring

```bash
# AWS Lambda monitoring
srediag monitor aws lambda --function myfunction

# Azure Functions analysis
srediag analyze azure function --app myapp --function myfunction

# GCP Cloud Functions
srediag monitor gcp function --function myfunction

# OCI Functions
srediag monitor oci function --app myapp --function myfunction
```

## Configuration

Example configuration in `srediag.yaml`:

```yaml
cloud:
  aws:
    region: us-west-2
    credentials:
      profile: default
  azure:
    subscription_id: your_subscription_id
    resource_group: your_resource_group
  gcp:
    project_id: your_project_id
    zone: us-central1-a
  oci:
    compartment_id: your_compartment_id
    region: us-phoenix-1
```

## Usage Examples

### Resource Health Check

```bash
# Check AWS resources
srediag check aws --region us-west-2

# Monitor Azure resources
srediag monitor azure --resource-group mygroup

# Analyze GCP resources
srediag analyze gcp --project myproject

# Check OCI resources
srediag check oci --compartment mycompartment
```

### Cost Analysis

```bash
# AWS cost analysis
srediag analyze aws costs --period 30d

# Azure cost management
srediag analyze azure costs --subscription mysub

# GCP billing analysis
srediag analyze gcp billing --project myproject

# OCI cost tracking
srediag analyze oci costs --compartment mycompartment
```

## OpenTelemetry Integration

SREDIAG integrates cloud monitoring with OpenTelemetry:

```bash
# Export AWS metrics to OpenTelemetry
srediag monitor aws --export otlp://localhost:4317

# Azure metrics with OpenTelemetry context
srediag monitor azure --context otel-context.txt

# GCP monitoring with sampling
srediag monitor gcp --sampling-ratio 0.1

# OCI metrics export
srediag monitor oci --export otlp://collector:4317
```

## See Also

- [CLI Documentation](../cli/README.md)
- [Configuration Guide](../configuration/README.md)
- [Security Guide](../security/README.md)
- [Troubleshooting](../reference/troubleshooting.md)
