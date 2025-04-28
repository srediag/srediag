# SREDIAG Documentation

## Overview

SREDIAG is the edge data-plane for the OBSERVO reliability platform. It extends the upstream OpenTelemetry Collector and incrementally adds MSP-grade features such as hot-swappable plugins, content-aware deduplication, CMDB drift detection, and tenant isolation.

This directory contains detailed guides for installation, configuration, CLI usage, architecture, and extensibility.

For the full technical specification and roadmap, see `docs/specification.md`.

## Table of Contents

### General

- [README](README.md)
- [TODO](TODO.md)
- [Code of Conduct](CODE_OF_CONDUCT.md)
- [Contributing](CONTRIBUTING.md)

### Getting Started

- [Overview](getting-started/README.md)
- [Quick Start](getting-started/quickstart.md)
- [Installation](getting-started/installation.md)
- [First Steps](getting-started/first-steps.md)
- [Basic Configuration](getting-started/configuration.md)
- [Advanced Configuration](getting-started/advanced-configuration.md)
- [Troubleshooting](getting-started/troubleshooting.md)

### Architecture

- [Overview](architecture/README.md)
- [Plugin Architecture](architecture/plugin.md)
- [Service Architecture](architecture/service.md)
- [Build Architecture](architecture/build.md)
- [Diagnose Architecture](architecture/diagnose.md)
- [Security Architecture](architecture/security.md)

### Configuration

- [Overview](configuration/README.md)
- [Plugin Configuration](configuration/plugin.md)
- [Service Configuration](configuration/service.md)
- [Build Configuration](configuration/build.md)
- [Diagnose Configuration](configuration/diagnose.md)
- [Security Settings](configuration/security.md)

### Plugin System

- [Overview](plugin/README.md)
- [Development Guide](plugin/development.md)
- [API Reference](plugin/types.md)
- [Plugin Architecture](plugin/plugin-architecture.md)
- [Best Practices](plugin/best-practices.md)
  - [Examples](plugin/examples/README.md)

### CLI Tools

- [Overview](cli/README.md)
- [Plugin CLI](cli/plugin.md)
- [Service CLI](cli/service.md)
- [Build CLI](cli/build.md)
- [Diagnose CLI](cli/diagnose.md)
- [Security CLI](cli/security.md)

### Cloud Integration

- [Overview](cloud/README.md)
- [AWS Integration](cloud/aws.md)
- [OCI Integration](cloud/oci.md)
- [Azure Integration](cloud/azure.md)
- [GCP Integration](cloud/gcp.md)
- [Kubernetes Deployment](cloud/kubernetes.md)

### Reference

- [Overview](reference/README.md)
- [API Reference](reference/api.md)
- [Performance Guidelines](reference/performance.md)
- [Production Hardening](reference/production.md)
- [Troubleshooting](reference/troubleshooting.md)

### Security

- [Overview](security/README.md)

---

## Support

- [GitHub Issues](https://github.com/srediag/srediag/issues)  
- [GitHub Discussions](https://github.com/srediag/srediag/discussions)  
- [Security Policy](../SECURITY.md)  
