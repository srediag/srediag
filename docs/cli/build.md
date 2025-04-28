# `srediag build`

Wraps *otelcol-builder* and helper scripts to compile the agent or
individual plugins.

```bash

Usage:
  srediag build <command> [flags]

Persistent flags:
  --config <file>      Path to builder YAML  (default "build/srediag-build.yaml")
  --output-dir <path>  Where artefacts are stored  (default "/tmp/srediag-build")

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
srediag build all --config build/srediag-build.yaml

# Rebuild clickhouse exporter only
srediag build plugin --type exporter --name clickhouseexporter

# Generate skeleton code for new processor
srediag build generate --type processor --name myprocessor
```

Artifacts:

* Agent → `<output-dir>/srediag`
* Plugins → `<output-dir>/plugins/<name>/<name>`  
  (copy to `plugins.exec_dir` or use `install` command)
