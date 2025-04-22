# SREDIAG – OBSERVO Diagnostics Agent

SREDIAG is a powerful diagnostic and observability agent designed for modern SRE practices.
It combines robust telemetry collection with advanced diagnostic capabilities, providing deep insights into systems, applications, and infrastructure.
Whether you’re troubleshooting a Kubernetes cluster, analyzing cloud infrastructure, or monitoring application performance, SREDIAG provides the tools you need for effective system reliability engineering.

## Key Features

- **Advanced Diagnostics**
  - Real‑time system analysis and troubleshooting
  - Kubernetes cluster health diagnostics
  - Infrastructure configuration analysis
  - Performance bottleneck detection
  - Root cause analysis

- **SRE Tooling**
  - SLO/SLI monitoring and tracking
  - Error budget management
  - Capacity planning
  - Performance analysis
  - Automated remediation suggestions

- **Multi‑Platform Support**
  - Kubernetes environments
  - Cloud providers (AWS, Azure, GCP)
  - Bare‑metal servers
  - Containerized applications
  - Virtual machines

- **Configuration Analysis**
  - Infrastructure as Code validation
  - Security compliance checking
  - Best practices enforcement
  - Configuration drift detection
  - Cost optimization recommendations

- **Plugin System**
  - Extensible architecture
  - Custom diagnostic capabilities
  - Integration with existing tools
  - Community plugin ecosystem
  - Plugin marketplace

- **CLI Tools**
  - Interactive diagnostics
  - Health checking
  - Performance analysis
  - Configuration validation
  - Resource optimization

## Architecture

SREDIAG is built on a modular architecture designed for reliability and extensibility:

1. **Diagnostic Engine**
   Advanced analysis and troubleshooting capabilities

2. **Plugin System**
   Extensible architecture for custom diagnostics

3. **CLI Interface**
   Powerful command‑line diagnostic tools

4. **Telemetry Pipeline**
   Efficient data collection and analysis

5. **Integration Layer**
   Connects with external systems and tools

6. **Security Layer**
   Ensures secure operation and data handling

## Getting Started

### Prerequisites

- Go 1.24 or higher
- Make
- Git
- kubectl (for Kubernetes features)
- Cloud provider CLI tools (optional)

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/observo/srediag.git
   cd srediag
   ```

2. Build the agent:

   ```bash
   make build
   ```

3. Create a configuration file:

   ```bash
   cp configs/config.yaml.example configs/config.yaml
   ```

4. Edit the configuration file according to your needs.

### Basic Usage

1. Run system diagnostics:

   ```bash
   srediag diagnose system
   ```

2. Analyze Kubernetes cluster:

   ```bash
   srediag diagnose kubernetes --cluster my-cluster
   ```

3. Check configuration:

   ```bash
   srediag analyze config --path /path/to/config
   ```

4. Start monitoring:

   ```bash
   srediag monitor --target system|kubernetes|cloud
   ```

## Configuration

SREDIAG uses a YAML configuration file with the following structure:

```yaml
debug: false
log_level: "info"

telemetry:
  enabled: true
  service_name: "srediag"
  endpoint: "http://localhost:4317"

diagnostics:
  system:
    enabled: true
    interval: "1m"
  kubernetes:
    enabled: true
    context: "my-cluster"
  cloud:
    enabled: true
    providers:
      - aws
      - gcp

plugins:
  directory: "plugins"
  enabled:
    - "system-diagnostics"
    - "k8s-analyzer"
    - "cloud-diagnostics"
  settings:
    system-diagnostics:
      collect_interval: "30s"
    k8s-analyzer:
      include_events: true
    cloud-diagnostics:
      regions:
        - "us-east-1"
        - "eu-west-1"
```

## Plugin Development

SREDIAG supports various types of plugins:

1. **Diagnostic Plugins**: Create custom diagnostic capabilities
2. **Collector Plugins**: Gather specific metrics, logs, or traces
3. **Analysis Plugins**: Implement custom analysis algorithms
4. **Integration Plugins**: Connect with external systems

For detailed information about plugin development, see [Plugin Development Guide](docs/plugin-development.md).

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.

## Support

- [Documentation](docs/)
- [Issue Tracker](https://github.com/observo/srediag/issues)
- [Discussions](https://github.com/observo/srediag/discussions)

## Roadmap

See our [TODO](TODO.md) file for planned features and improvements.
