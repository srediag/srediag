# SREDIAG Architecture

This document provides the foundational architectural concepts underpinning **SREDIAG**, detailing the system structure, subsystem interactions, and especially the central role of the `internal/core` package, critical for consistent component management and operation across the project.

---

## 1 · High-Level System Overview

```ascii
┌───────────────────────────────────────────────────────────────────────┐
│                          SREDIAG Agent                                │
│                                                                       │
│ ┌───────────┐   ┌───────────────┐   ┌───────────────┐   ┌───────────┐ │
│ │ Receivers │──▶│  Processors   │──▶│   Exporters   │──▶│ ObservO CP│ │
│ └───────────┘   └───────────────┘   └───────────────┘   └───────────┘ │
│      ▲                 ▲                   ▲                          │
│      │                 │                   │                          │
│ ┌─────────────┐  ┌───────────────┐  ┌───────────────┐                 │
│ │ Plugin      │  │ Deduplication │  │ CMDB & Drift  │                 │
│ │ Manager     │  │ Compression   │  │ Detection     │                 │
│ │(hot-swap,   │  │(xxHash,LZ4,ZSTD│ │(fingerprint,  │                 │
│ │ IPC,seccomp)│  │ RocksDB cache) │ │delta tracking)│                 │
│ └─────────────┘  └───────────────┘  └───────────────┘                 │
│      ▼                 ▼                   ▼                          │
│ ┌───────────────────────────────────────────────────────────┐         │
│ │ Resource & Governance Controller                          │         │
│ │(Cgroups, Memory Guard, Quotas, Seccomp/AppArmor sandbox)  │         │
│ └───────────────────────────────────────────────────────────┘         │
│                      │                                                │
│            Self-Metrics & Observability                               │
└───────────────────────────────────────────────────────────────────────┘
```

**Key Architectural Patterns:**

- **Collector Pipeline**: Receivers → Processors → Exporters → ObservO Control Plane.
- **Plugin System**: Dynamic extension points for Receivers, Processors, Exporters, and CLI diagnostics.
- **Governance Layer**: Strict resource isolation, quota enforcement, and security hardening through Cgroups and seccomp.
- **Self-Observability**: Metrics, health endpoints, and tracing via OpenTelemetry zPages.

---

## 2 · Core Package Architecture (`internal/core`)

The **Core Package** is the foundational library shared across all SREDIAG commands and plugins. It abstracts essential functionality like configuration management, logging, component registration, lifecycle, and operational context management.

### 2.1 · Configuration Loader Architecture

The Configuration Loader provides a unified, consistent approach to configuration management for both CLI and daemon modes.

- **Configuration Sources** (in order of priority):
  1. **CLI Flags**
  2. **Environment Variables**
  3. **YAML Configuration** (`srediag.yaml` or specified via flag)
  4. **Built-in Defaults**

- **Validation & Reloading**:
  - Performs rigorous validation of configuration at startup and on dynamic reload (via `SIGHUP`).
  - Supports hierarchical configuration structure, clearly separating service, plugins, logging, security, and collector settings.

- **Dynamic Hot-Reloading**:
  - On `SIGHUP`, the runtime seamlessly applies new configuration values without downtime.

---

### 2.2 · Logging System

Unified logging via Zap across CLI and Service components:

- **Bootstrap Stage**:
  - Initial minimal logger setup (controlled by ENV variables: `SREDIAG_LOG_LEVEL`, `SREDIAG_LOG_FORMAT`).
- **Runtime Stage**:
  - Reconfigures logging as per loaded YAML config.
- **Consistency**:
  - All components share the same enriched logger (version, build metadata).

---

### 2.3 · Component Registry & Manager

Provides a robust, type-safe component registration system supporting built-in Collector core components and plugins.

- **Factory Registration**:
  - Collectors register components at build-time using the registry pattern, ensuring uniqueness and proper type adherence.
- **Component Manager**:
  - Dynamically retrieves component factories (receivers, processors, exporters, extensions) to build pipelines or invoke CLI plugins.
  - Facilitates dynamic plugin loading and lifecycle management.

---

### 2.4 · Build Information & Versioning

Captures detailed metadata about the running binary:

- Embeds CalVer, Git commit SHA, Go runtime version, OTel Collector version.
- Accessible via `srediag version` CLI and embedded within observability metrics and logs.

---

### 2.5 · Application Context (`AppContext`)

Encapsulates core services required across the application:

```go
type AppContext struct {
    Logger           *zap.Logger
    Config           *Config
    ComponentManager *ComponentManager
    BuildInfo        BuildInfo
}
```

**Usage Pattern**: Passed explicitly to all command handlers, plugins, and subsystems ensuring consistent runtime behavior and simplified dependency management.

---

### 2.6 · Runtime Lifecycle & CLI Integration

The Core package facilitates lifecycle handling for both CLI and service modes:

- **CLI Execution Flow**:

  ```bash
  srediag build|plugin|security|diagnose [...]
  ```

  - Parses configuration and initializes `AppContext`.
  - Dispatches to relevant handlers and subcommands.

- **Service Mode**:

  ```bash
  srediag service start [...]
  ```

  - Initializes runtime components and OTel pipelines.
  - Listens for Unix signals (`SIGHUP` for reload, `SIGTERM` for graceful shutdown, `SIGUSR2` for plugin reload).

---

## 3 · Subsystem Architectural Documents

| Document                              | Responsibility                                                   |
|---------------------------------------|------------------------------------------------------------------|
| [Build Architecture](build.md)        | CI/CD, binary & plugin compilation, versioning, SBOM, cosign signing |
| [Diagnostics Architecture](diagnose.md)| CLI-driven diagnostics and plugin-based health and performance checks |
| [Plugin Architecture](plugin.md)      | Hot-swap framework, sandboxing, IPC (shmipc-go), SDK contracts    |
| [Security Architecture](security.md)  | TLS/mTLS, certificate management, RBAC, runtime sandboxing, compliance |
| [Service Architecture](service.md)    | Long-running service lifecycle, pipeline orchestration, self-observability |

Each subsystem is described in detail within its dedicated architectural documentation.

---

## 4 · Cross-Document Reference

| Component                           | Architecture Document               |
|-------------------------------------|-------------------------------------|
| Plugin loading, sandbox, validation | [Plugin Architecture](plugin.md)    |
| Certificate handling, sandboxing    | [Security Architecture](security.md)|
| Service lifecycle and management    | [Service Architecture](service.md)  |
| Build, CI/CD, artifact generation   | [Build Architecture](build.md)      |
| Diagnostic command & plugin flow    | [Diagnostics Architecture](diagnose.md) |

---

## 5 · Further Reading & Resources

- **Operational Configuration**: [`configuration/README.md`](../configuration/README.md)  
- **CLI User Guide**: [`cli/README.md`](../cli/README.md)  

---

## 6 · Document Governance

- Documented under the MIT License.  
