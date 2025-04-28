# SREDIAG Architecture

SREDIAG is the edge data-plane for the OBSERVO reliability platform. It embeds the OpenTelemetry Collector core (v0.124.0) and layers on MSP-grade features—hot-swappable plugins, advanced deduplication, CMDB drift, diagnostics, and hardened service runtime—powered by our reusable **internal/core** library.

---

## 1 · Core Library (`/internal/core`)

Every SREDIAG component—CLI, Collector pipelines, service, plugins & diagnostics—boots and wires up via **internal/core**:

- **Configuration**
  Loads defaults, environment overrides, and `--config` YAML into a single `Config` struct.
- **Logging**
  `core.NewLogger(&core.Logger{…})` produces a structured, Zap-based logger used across all modules.
- **Component Manager**
  `core.NewComponentManager(logger)` tracks and instantiates Collector factories (receivers/processors/exporters).
- **AppContext**
  Bundles `Logger`, `Config`, `ComponentManager`, and `BuildInfo` for passing into CLI commands and service.
- **Version & Flags**
  Exposes `core.NewVersionCmd()` and standard `--log-level`, `--log-format`, `--quiet`, etc., on every command.

This shared foundation ensures consistency in error handling, observability, and configuration across the five architectural domains below.

---

## 2 · High-Level View

```ascii
                                /─ diagnostics CLI ─┐
                               ▼                     │
   ┌────────────────────────────────────────────────────────────────┐
   │                        SREDIAG Agent (Go)                      │
   │                                                                │
   │  ┌─────┐   ┌───────┐   ┌───────┐   ┌────────┐   ┌──────────────┐ │
   │  │Core │──▶│Plugin │──▶│Pipeline│──▶│Service │──▶│Self-Metrics  │ │
   │  │Lib  │   │Manager│   │Components│ │Runtime │   │& zPages     │ │
   │  └─────┘   └───────┘   └───────┘   └────────┘   └──────────────┘ │
   │      ▲                             ▲                           │
   │      │                             │                           │
   │   `/internal/core`           `/internal/service`             │
   └────────────────────────────────────────────────────────────────┘
```

- **Core Library (`/internal/core`)**
  Bootstraps config, logging, CLI flags, and component registry.

- **Plugin Manager (`/internal/plugin`)**
  Hot-loads `.so` plugins via shmipc, cosign-verifies bundles, sandboxes via seccomp/AppArmor.
  See [plugin architecture][plugin] for details.

- **Pipeline Components (`/internal/…`)**
  Native OTel receivers/processors/exporters + custom modules (dedup, CMDB, compression).

- **Service Runtime (`/internal/service`)**
  Orchestrates Collector pipelines as a long-running service with graceful shutdown, metrics endpoints, and quotas.
  See [service architecture][service].

- **Diagnostics CLI (`/cmd/srediag/commands/diagnose`)**
  On-demand system, performance, and security checks implemented via plugin-style modules.
  See [diagnose architecture][diagnose].

---

## 3 · Sub-documents

Each of the five detailed architecture guides lives under `docs/architecture/`:

| Guide                                    | Description                                             |
|:-----------------------------------------|:--------------------------------------------------------|
| **[build][build]**                       | How Collector core & plugins are compiled via `otelcol-builder`, CI integration, and version sync. |
| **[diagnose][diagnose]**                 | Design of the `srediag diagnose …` commands, plugin hooks, and output integration. |
| **[plugin][plugin]**                     | Hot-swap plugin framework: IPC contracts, security checks, lifecycle, and SDK. |
| **[security][security]**                 | Runtime hardening: TLS/mTLS, RBAC, seccomp/AppArmor, rate-limits, and supply-chain controls. |
| **[service][service]**                   | Long-running service mode: pipeline orchestration, graceful shutdown, health/metrics endpoints. |

---

## 4 · How It All Fits

1. **Startup**
   `main.go` (in `/cmd/srediag`) calls `core.LoadConfig()` → `core.NewLogger()` → `core.NewComponentManager()` → builds an `AppContext` .

2. **CLI Dispatch**
   `Execute(ctx)` registers root flags and subcommands: `start` (service), `build`, `plugin`, `diagnose`, etc., reusing `core` defaults.

3. **Service Mode** (`srediag start`)
   - Uses `ComponentManager.GetFactories()` to assemble OTel pipelines.
   - Wraps them in `service.NewService(...)` which invokes `/internal/service` logic.
   - Exposes `/metrics`, `/healthz`, zPages, and responds to signals (SIGHUP, SIGTERM, SIGUSR2).

4. **Plugin Loading**
   - On startup, built-in plugins from `otelcol-builder.yaml` are registered.
   - `srediag plugin load/unload` posts to UNIX socket → `PluginManager.Load()` or `.Unload()`.
   - Bundles are validated (SHA256 + cosign), mapped by component type, and sandboxed.

5. **Diagnostics**
   - `srediag diagnose <system|performance|security>` spins up lightweight plugin modules in `/internal/diagnostic`, runs checks, and prints results.

---

## 5 · Next Steps

- Review each sub-guide for implementation details and code snippets.
- Ensure `/internal/core` changes propagate cleanly into CLI and service.
- Align build pipelines to emit both binary and plugin artifacts per [build guide][build].
- Keep architecture in sync with `docs/TECHNICAL_SPECIFICATION.md` for traceability.

---

[build]: build.md
[diagnose]: diagnose.md
[plugin]: plugin.md
[security]: security.md
[service]: service.md
