# SREDIAG – OBSERVO Diagnostics Agent

SREDIAG is a powerful diagnostic and observability agent built on top of OpenTelemetry Collector.
It extends the OpenTelemetry Collector's capabilities with advanced diagnostic features, providing deep insights into systems, applications, and infrastructure.
Whether you're troubleshooting a Kubernetes cluster, analyzing cloud infrastructure, or monitoring application performance, SREDIAG combines OpenTelemetry's robust telemetry collection with specialized diagnostic tools.

## Key Features

- **OpenTelemetry Integration**
  - Built on OpenTelemetry Collector architecture
  - Full compatibility with OpenTelemetry protocols
  - Support for all OpenTelemetry core components
  - Extended collector functionality
  - Custom processors and exporters

- **Advanced Diagnostics**
  - Real‑time system analysis and troubleshooting
  - Kubernetes cluster health diagnostics
  - Infrastructure configuration analysis
  - Performance bottleneck detection
  - Root cause analysis automation

- **Management Capabilities**
  - Configuration management
  - Resource provisioning
  - Policy enforcement
  - Automated remediation
  - Change management
  - State reconciliation

- **Multi‑Platform Support**
  - Kubernetes environments
  - Oracle Cloud Infrastructure (OCI)
  - AWS environments
  - Azure platforms
  - Bare‑metal servers
  - Virtual machines

- **Configuration Analysis**
  - Infrastructure as Code validation
  - Security compliance checking
  - Best practices enforcement
  - Configuration drift detection
  - Cost optimization recommendations

- **Extended Plugin System**
  - Diagnostic plugins
  - Analysis plugins
  - Integration plugins
  - Management plugins
  - Custom plugin development
  - Plugin marketplace support

- **CLI Tools**
  - Interactive diagnostics
  - System analysis
  - Performance profiling
  - Configuration management
  - Resource optimization
  - State management

## Architecture

SREDIAG extends the OpenTelemetry Collector architecture with additional capabilities:

1. **OpenTelemetry Core**
   - Full OpenTelemetry Collector integration
   - Standard receivers, processors, and exporters
   - Native protocol support
   - Built-in scalability and reliability

2. **Diagnostic Engine**
   - Advanced analysis capabilities
   - Real-time troubleshooting
   - Custom diagnostic processors
   - Extended metrics collection

3. **Management Engine**
   - Configuration management
   - State reconciliation
   - Change application
   - Policy enforcement

4. **Plugin System**
   - Diagnostic plugins
   - Analysis plugins
   - Integration plugins
   - Management plugins

5. **Security Layer**
   - Authentication and authorization
   - Data encryption
   - Compliance monitoring
   - Access control

For detailed architecture information, see [Architecture Documentation](docs/architecture/overview.md).

## Getting Started

### Prerequisites

- Go 1.24 or higher
- Make
- Git
- kubectl (for Kubernetes features)
- Cloud provider CLI tools (optional)

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/srediag/srediag.git
   cd srediag
   ```

2. Build the agent:

   ```bash
   make build
   ```

3. Configuration:

   The configuration files are located in the `configs` directory:
   - `configs/srediag.yaml` - Main application configuration
   - `configs/otel-config.yaml` - OpenTelemetry Collector configuration

   SREDIAG will look for configuration files in the following order:
   1. Path specified in `SREDIAG_CONFIG` environment variable
   2. `configs/srediag.yaml` in the project directory
   3. `srediag.yaml` in the current directory
   4. `.srediag.yaml` in the current directory
   5. `~/.srediag/config/srediag.yaml` in the home directory
   6. `~/.srediag.yaml` in the home directory

   You can also override the plugins directory using the `SREDIAG_PLUGIN_DIR` environment variable.

### Basic Usage

1. Run system diagnostics:

   ```bash
   srediag diagnose system
   ```

2. Analyze Kubernetes cluster:

   ```bash
   srediag diagnose kubernetes --cluster my-cluster
   ```

3. Apply configuration changes:

   ```bash
   srediag apply config --path /path/to/config
   ```

4. Start monitoring:

   ```bash
   srediag monitor --target system|kubernetes|cloud
   ```

For detailed usage instructions, see [CLI Documentation](docs/cli/README.md).

## Documentation

- [Getting Started](docs/getting-started/README.md)
- [Architecture](docs/architecture/README.md)
- [Configuration](docs/configuration/README.md)
- [Plugin System](docs/plugins/README.md)
- [CLI Reference](docs/cli/README.md)
- [Cloud Integration](docs/cloud/README.md)
- [API Reference](docs/reference/api.md)
- [Security](docs/architecture/security.md)

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.

## Support

- [Documentation](docs/README.md)
- [Issue Tracker](https://github.com/srediag/srediag/issues)
- [Discussions](https://github.com/srediag/srediag/discussions)

## Roadmap

See our [TODO List](docs/TODO.md) for planned features and improvements.

## CLI Tools

SREDIAG provides powerful diagnostic CLI tools for interactive troubleshooting and analysis:

### System Diagnostics

```bash
# Basic system diagnostics
srediag diagnose system

# Specific resource analysis
srediag diagnose system --resource cpu
srediag diagnose system --resource memory
srediag diagnose system --resource disk

# Real-time monitoring
srediag monitor system --interval 5s

# Process analysis
srediag analyze process --pid 1234
```

### Kubernetes Diagnostics

```bash
# Cluster health check
srediag diagnose kubernetes --cluster my-cluster

# Resource analysis
srediag analyze kubernetes resources --namespace default
srediag analyze kubernetes deployments --show-events

# Configuration audit
srediag audit kubernetes --namespace production
```

### Performance Analysis

```bash
# CPU profiling
srediag profile cpu --duration 30s --output cpu.prof

# Memory analysis
srediag analyze memory --threshold 90

# Goroutine inspection
srediag debug goroutines --trace

# Bottleneck detection
srediag analyze bottlenecks --service my-service
```

### Security Diagnostics

```bash
# Vulnerability scanning
srediag scan vulnerabilities --severity high

# Compliance checking
srediag check compliance --standard pci-dss

# Configuration audit
srediag audit config --path /etc/srediag
```

### Analysis Features

- Real-time metric collection and analysis
- Automatic correlation detection
- Root cause analysis
- Anomaly detection
- Performance bottleneck identification
- Security vulnerability assessment
- Configuration validation
- Resource optimization recommendations

## Plugin Development

### Building Plugins

To build a plugin, use the following bash command:

```bash
go build -buildmode=plugin -o bin/plugins/.tmp/plugin_name.so .
```

This will create a shared object file that can be loaded by SREDIAG. Make sure to:

1. Use bash shell (not PowerShell) for building plugins
2. Place the built plugin in the `bin/plugins/.tmp` directory
3. Name the plugin file with a `.so` extension

For example, to build the simple receiver plugin:

```bash
cd plugins/examples/simplereceiver
go build -buildmode=plugin -o bin/plugins/.tmp/simplereceiver.so .
```

### Plugin Configuration

```
