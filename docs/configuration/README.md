# SREDIAG — Configuration Guide

SREDIAG is an extensible agent whose behaviour is defined by:

1. A **core YAML** (`srediag.yaml`) – controls services, logging,
   security, plugin directories.
2. An **optional collector YAML** (`srediag-service.yaml`) – standard
   OpenTelemetry pipelines for logs/metrics/traces.
3. **Environment variables + command-line flags** that override both
   at runtime.

| Layer | File / Source | Applies To | Can be absent? |
| :---- | :------------ | :--------- | :------------- |
| **1** | Command-line flags | CLI + Service | Yes |
| **2** | Environment variables | CLI + Service | Yes |
| **3** | `srediag.yaml` | CLI + Service | Yes (defaults kick in) |
| **4** | `srediag-service.yaml` | Service only | Yes (if `collector.enabled=false`) |

---

## 1 · File Locations & Discovery

```text
/etc/srediag/
├─ srediag.yaml             # Core (always parsed if present)
├─ srediag-service.yaml     # OTel pipeline (service mode)
└─ plugins/                 # .so artifacts (optional)
```

Discovery order for `srediag.yaml`:

1. `--config <file>` flag  
2. `SREDIAG_CONFIG` env var  
3. First existing path in  
   `/etc/srediag/srediag.yaml` → `$HOME/.srediag/config.yaml` →
   `./config/srediag.yaml` → `./srediag.yaml`

---

## 2 · Core YAML (`srediag.yaml`) Reference

```yaml
service:
  name: srediag
  port: 8080                       # HTTP UI / healthz
  environment: prod                # free-form tag

logging:
  level: info                      # debug|info|warn|error
  format: console                  # console|json

security:
  tls:
    enabled: true
    cert_file: /etc/srediag/cert.pem
    key_file:  /etc/srediag/key.pem

collector:                          # Parsed **only** in service mode
  enabled: true
  config_path: /etc/srediag/srediag-service.yaml
  memory_limit_mib: 1024

plugins:
  dir: /var/lib/srediag/plugins
  enabled:                         # pre-load on start
    - processor/vectorhashprocessor
```

Unknown keys are ignored (logged at `debug` level), allowing forward
compatibility.

---

## 3 · Collector YAML (`srediag-service.yaml`)

Follows the upstream **OpenTelemetry Collector v0.124.0** schema.
Typical snippet 👇

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317

processors:
  batch:
  memory_limiter:
    check_interval: 2s
    limit_mib: 1024

exporters:
  otlp:
    endpoint: gw.observo.local:4317
    tls:
      insecure: false

service:
  pipelines:
    traces:
      receivers:  [otlp]
      processors: [memory_limiter, batch]
      exporters:  [otlp]
```

Extra SREDIAG-only processors (e.g., `vectorhashprocessor`) become
available automatically once the corresponding plugin is **built and
loaded**.

Full field reference: `docs/configuration/service.md`.

---

## 4 · Build-time YAML (`srediag-build.yaml`)

Controls what **components are compiled in** (`otelcol-builder` spec).

```yaml
dist:
  name: srediag
  version: 0.1.0
receivers:
  - gomod: go.opentelemetry.io/collector/receiver/otlpreceiver  v0.124.0
processors:
  - gomod: github.com/srediag/processors/vectorhashprocessor     v0.1.0
exporters:
  - gomod: go.opentelemetry.io/collector/exporter/otlpexporter  v0.124.0
```

Regenerate the binary after edits:

```bash
srediag build all --config build/srediag-build.yaml
```

---

## 5 · Environment Variable Map

| YAML Key | Env Var | Notes |
| :------- | :------ | :---- |
| `logging.level` | `SREDIAG_LOG_LEVEL` | overrides YAML |
| `logging.format` | `SREDIAG_LOG_FORMAT` | `console` / `json` |
| `plugins.dir` | `SREDIAG_PLUGINS_DIR` | path to look for `.so` files |

---

## 6 · Precedence Rules

```text
Flags   >   Env vars   >   YAML   >   Built-ins
```

*Example* — given

```bash
export SREDIAG_LOG_LEVEL=debug
srediag diagnose system --log-level warn
```

the effective level is **warn** (flag beats env).

---

## 7 · Validation & Reload

* Core YAML is validated on startup; fatal errors abort execution.  
* Service mode watches both YAML files for `SIGHUP` — the collector
  reloads pipelines live; the core layer restarts plugin discovery.

---

## 8 · Minimal Working Examples

### 8.1 CLI-only

```yaml
service:
  name: srediag
logging:
  level: info
  format: console
```

```bash
srediag diagnose system        # uses defaults, no collector
```

### 8.2 Full Service Pipeline

```yaml
service:
  name: srediag
  environment: staging

collector:
  enabled: true
  config_path: /etc/srediag/srediag-service.yaml

plugins:
  enabled:
    - processor/vectorhashprocessor
```

---

## 9 · Best Practices

* Keep core YAML lean; heavy pipeline logic lives in
  `srediag-service.yaml`.  
* Store secrets outside YAML (env vars, K8s Secrets, Vault).  
* In Kubernetes mount both YAMLs as **ConfigMaps** and send `SIGHUP`
  for zero-downtime reloads.  
* Components **not compiled in** cannot be referenced at runtime —
  keep `srediag-build.yaml` and collector YAML in sync.

---

## 10 · Related Docs

* [CLI Guide](../cli/README.md)  
* [Service-mode Deep-Dive](service.md)  
* [Build System](../build.md)  
* [Plugin Architecture](../plugins/README.md)

---

## Parameter Reference

| YAML Key                | Env Var                        | CLI Flag                |
|-------------------------|--------------------------------|-------------------------|
| `srediag.config`        | `SREDIAG_CONFIG`               | `--config`              |
| `build.config_path`          | `SREDIAG_BUILD_CONFIG_PATH``         | `--build-config`        |
| `build.output_dir`      | `SREDIAG_BUILD_OUTPUT_DIR`     | `--output-dir`          |
| `plugins.dir`           | `SREDIAG_PLUGINS_DIR`          | `--plugins-dir`         |
| `plugins.exec_dir`      | `SREDIAG_PLUGINS_EXEC_DIR`     | `--exec-dir`            |
| `collector.config_path` | `SREDIAG_COLLECTOR_CONFIG_PATH`| `--service-yaml`        |
| `diagnostics.config_path` | `SREDIAG_DIAGNOSTICS_CONFIG_PATH`| `--diag-service-yaml`        |
| `diagnostics.defaults.output_format` | `SREDIAG_DIAG_OUTPUT_FORMAT` | `--diag-output-format` |
| `diagnostics.defaults.timeout`       | `SREDIAG_DIAG_TIMEOUT`       | `--diag-timeout`             |

> **Parameter Naming and Precedence:**
>
> * `--config`/`SREDIAG_CONFIG` is **only** for the main SREDIAG config. Subsystem configs (build, plugin, service, diagnostics) use their own unique flags and env vars as shown above.
> * **Precedence:** CLI flags > Env vars > YAML config > Built-in defaults

---

## Discovery Order & Precedence

1. CLI flags (highest)
2. Environment variables
3. YAML config file (see discovery order below)
4. Built-in defaults (lowest)

**Config file discovery order:**

1. `--config <file>` flag (main config)
2. `SREDIAG_CONFIG` env var (main config)
3. `/etc/srediag/srediag.yaml`
4. `$HOME/.srediag/config.yaml`
5. `./config/srediag.yaml`
6. `./srediag.yaml`

---

## Best Practices

* Use CLI flags or environment variables for automation and CI/CD.
* Use YAML for persistent, version-controlled configuration.
* Always check the effective config with `srediag --print-config` (if available).
* Unknown YAML keys are logged at debug level and ignored for forward compatibility.
