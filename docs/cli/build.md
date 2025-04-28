# `srediag build`

Wraps *otelcol-builder* and helper scripts to compile the agent or
individual plugins.

| YAML Key                | Env Var                    | CLI Flag                |
|-------------------------|----------------------------|-------------------------|
| `build.config`          | `SREDIAG_BUILD_CONFIG`     | `--build-config`        |
| `build.output_dir`      | `SREDIAG_BUILD_OUTPUT_DIR` | `--output-dir`          |

> **Warning:** Do **not** use `--config` for builder YAML; this is reserved for the main SREDIAG config. Always use `--build-config` for build operations.

```bash

Usage:
  srediag build <command> [flags]

Persistent flags:
  --build-config <file>  Path to builder YAML  (default "build/srediag-build.yaml")
  --output-dir <path>    Where artefacts are stored  (default "/tmp/srediag-build")

```

---

## 1 · Commands

| Command | Purpose |
| :------ | :------ |
| `all` | Build agent **and** every plugin declared in the YAML |
| `plugin` | Build a single plugin (`--type`, `--name` required) |
| `generate` | Produce plugin scaffold code (no compile) |
| `install` | Copy pre-built plugins into `plugins.exec_dir` |
| `update-yaml-versions` | Sync builder YAML with `go.mod` |

---

## 2 · Examples

```bash
# Build everything
srediag build all --build-config build/srediag-build.yaml

# Rebuild clickhouse exporter only
srediag build plugin --type exporter --name clickhouseexporter

# Generate skeleton code for new processor
srediag build generate --type processor --name myprocessor

# Sync builder compoments with go.mod (uses UpdateBuilderYAMLVersions)
srediag build update --build-config configs/srediag-builder.yaml --gomod go.mod --plugin-gen plugin/generated
```

Artifacts:

* Agent → `<output-dir>/srediag`
* Plugins → `<output-dir>/plugins/<name>/<name>`  
  (copy to `plugins.exec_dir` or use `install` command)

---

**Note for maintainers:**

* The YAML/go.mod sync logic is implemented in the function `UpdateBuilderYAMLVersions` in `internal/build/update.go`.

---

## Parameter Reference

| YAML Key                | Env Var                    | CLI Flag                |
|-------------------------|----------------------------|-------------------------|
| `build.output_dir`      | `SREDIAG_BUILD_OUTPUT_DIR` | `--output-dir`          |
| `build.config`          | `SREDIAG_BUILD_CONFIG`     | `--build-config`        |

---

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

## Best Practices

* Use CLI flags or environment variables for automation and CI/CD.
* Use YAML for persistent, version-controlled configuration.
* Always check the effective config with `srediag build --print-config` (if available).
* Unknown YAML keys are logged at debug level and ignored for forward compatibility.
