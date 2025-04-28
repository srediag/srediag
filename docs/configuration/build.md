# SREDIAG — Build-time Configuration

| YAML Key           | Env Var                    | CLI Flag         |
|--------------------|---------------------------|------------------|
| `build.config`     | `SREDIAG_BUILD_CONFIG`     | `--build-config` |
| `build.output_dir` | `SREDIAG_BUILD_OUTPUT_DIR` | `--output-dir`   |

> **Warning:** Do **not** use `--config` for builder YAML; this is reserved for the main SREDIAG config. Always use `--build-config` for build operations.

SREDIAG compiles a *static* OpenTelemetry Collector plus any **first-party or
third-party plugins** using the upstream **otelcol-builder**.  
All build inputs live in a single YAML: **`srediag-build.yaml`**.

---

## 1 · File Schema

```yaml
dist:
  name: srediag            # Binary label
  description: SRE Diagnostics Agent
  version: 0.1.0
  output_path: ./bin/srediag

receivers:
  - gomod: go.opentelemetry.io/collector/receiver/otlpreceiver v0.124.0
  - gomod: github.com/srediag/receivers/journaldreceiver       v0.1.0

processors:
  - gomod: go.opentelemetry.io/collector/processor/batchprocessor v0.124.0
  - gomod: github.com/srediag/processors/vectorhashprocessor       v0.1.0

exporters:
  - gomod: go.opentelemetry.io/collector/exporter/otlpexporter v0.124.0
  - gomod: github.com/srediag/exporters/clickhouseexporter     v0.1.0

extensions:
  - gomod: go.opentelemetry.io/collector/extension/healthcheckextension v0.124.0
```

### Key rules

| Rule | Rationale |
| :--- | :--------- |
| **Exact Go module path** with **version tag** | Reproducibility |
| Components omitted here **cannot** be referenced at runtime | hard-fail on start |
| Use the `update-yaml-versions` sub-command to sync with `go.mod` | avoids drift |

---

## 2 · Building

```bash
# Compile agent + every declared component
srediag build all --config build/srediag-build.yaml

# Produce only a plugin (journaldreceiver)
srediag build plugin --type receiver --name journaldreceiver
```

Artifacts land in `./bin/` (agent) and `./plugins/` (plugins).

---

## 3 · CI Integration

* `mage ci` target always invokes `srediag build all` to guarantee
  artefacts match the YAML.  
* SBOM and cosign signatures are attached in the same step.

---

## 4 · Common Pitfalls

| Symptom | Likely Cause |
| :------ | :----------- |
| "component not found" during `srediag service` | present in collector YAML but missing from **build** YAML |
| ABI mismatch errors when loading `.so` plugins | agent built with Go < plugin Go version |

---

## 5 · References

* [OpenTelemetry Collector Builder](https://opentelemetry.io/docs/collector/build/)  
* [Technical Specification §6 Implementation Plan](../TECHNICAL_SPECIFICATION.md)

## Discovery Order & Precedence

1. CLI flags (highest)
2. Environment variables
3. YAML config file (see discovery order below)
4. Built-in defaults (lowest)

**Config file discovery order:**

1. `--build-config <file>` flag
2. `SREDIAG_BUILD_CONFIG` env var
3. `build/srediag-build.yaml` (default)

---
