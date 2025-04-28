# SREDIAG Documentation

## Overview

SREDIAG is the edge data-plane for the OBSERVO reliability platform. It extends the upstream OpenTelemetry Collector and incrementally adds MSP-grade features such as hot-swappable plugins, content-aware deduplication, CMDB drift detection, and tenant isolation.

This directory contains detailed guides for installation, configuration, CLI usage, architecture, and extensibility.

For the full technical specification and roadmap, see `docs/specification.md`.

## Table of Contents

### Getting Started

- [Quick Start](getting-started/README.md)  
- [Installation](getting-started/installation.md)  
- [Basic Configuration](getting-started/configuration.md)

### Architecture

- [Overview](architecture/README.md)  
- [OpenTelemetry Integration](architecture/opentelemetry.md)  
- [Security Architecture](architecture/security.md)

### Configuration

- [Overview](configuration/README.md)  
- [Plugin Configuration](configuration/plugins.md)  
- [Security Settings](configuration/security.md)  
- [Advanced Options](configuration/advanced.md)

### Plugin System

- [Overview](plugins/README.md)  
- [Development Guide](plugins/development.md)  
- [API Reference](plugins/reference/api.md)  
- [Best Practices](plugins/reference/best-practices.md)

### CLI Tools

- [Overview](cli/README.md)  
- [Available Commands](cli/README.md)
- [Diagnostic Commands](cli/diagnostic.md)
- [Analysis Commands](cli/analysis.md)
- [Management Commands](cli/management.md)

### Cloud Integration

- [Overview](cloud/README.md)  
- [AWS Integration](cloud/aws.md)  
- [OCI Integration](cloud/oci.md)  
- [Azure Integration](cloud/azure.md)  
- [GCP Integration](cloud/gcp.md)  
- [Kubernetes Deployment](cloud/kubernetes.md)

### Development

- [Contributing](development/CODE_OF_CONDUCT.md)  
- [Development Setup](development/setup.md)  
- [Code Style](development/style.md)  
- [Testing](development/testing.md)

### Reference

- [API Reference](reference/api.md)  
- [Performance Guidelines](reference/performance.md)  
- [Production Hardening](reference/production.md)  
- [Troubleshooting](reference/troubleshooting.md)

---

## Support

- [GitHub Issues](https://github.com/srediag/srediag/issues)  
- [GitHub Discussions](https://github.com/srediag/srediag/discussions)  
- [Security Policy](../SECURITY.md)  
