# SREDIAG CLI Overview

| YAML Key                | Env Var                        | CLI Flag                |
|-------------------------|--------------------------------|-------------------------|
| `srediag.config`        | `SREDIAG_CONFIG`               | `--config`              |
| `build.config`          | `SREDIAG_BUILD_CONFIG`         | `--build-config`        |
| `build.output_dir`      | `SREDIAG_BUILD_OUTPUT_DIR`     | `--output-dir`          |
| `plugins.dir`           | `SREDIAG_PLUGINS_DIR`          | `--plugins-dir`         |
| `plugins.exec_dir`      | `SREDIAG_PLUGINS_EXEC_DIR`     | `--exec-dir`            |
| `collector.config_path` | `SREDIAG_COLLECTOR_CONFIG_PATH`| `--service-yaml`        |
| `diagnostics.defaults.output_format` | `SREDIAG_DIAG_OUTPUT_FORMAT` | `--output` / `--format` |
| `diagnostics.defaults.timeout`       | `SREDIAG_DIAG_TIMEOUT`       | `--timeout`             |

> **Parameter Naming and Precedence:**
>
> - `--config`/`SREDIAG_CONFIG` is **only** for the main SREDIAG config. Subsystem configs (build, plugin, service, diagnostics) use their own unique flags and env vars as shown above.
> - **Precedence:** CLI flags > Env vars > YAML config > Built-in defaults

SREDIAG’s command-line interface unifies all agent operations—building artifacts, running the long-lived service, managing plugins, performing security checks, and on-demand diagnostics—under a single `srediag` binary.

Commands are organized into five groups, each with its own detailed guide:

| Group           | Purpose                                                       | Documentation                  |
| :-------------- | :------------------------------------------------------------ | :----------------------------- |
| **Service**     | Start and manage the embedded OpenTelemetry Collector service | [Service Commands](service.md) |
| **Build**       | Compile the collector, plugins, images, SBOMs, and signatures | [Build Commands](build.md)     |
| **Plugin**      | Install, enable, configure, bind, verify and curate plugins   | [Plugin Commands](plugin.md)   |
| **Security**    | Rotate TLS certs, verify binaries/plugins, inspect sandbox    | [Security Commands](security.md) |
| **Diagnostics** | Execute on-demand system/performance/security checks          | [Diagnostics Commands](diagnose.md) |

---

## Global Options

All `srediag` commands support these common flags:

| Flag                 | Description                                   |
| :------------------- | :-------------------------------------------- |
| `-c, --config <path>` | Path to `srediag.yaml` (env: `SREDIAG_CONFIG`) |
| `--format <fmt>`     | Output format: `table` _(default)_, `json`, `yaml` |
| `--quiet`            | Suppress informational output; show only errors |
| `-h, --help`         | Display usage for any command or sub-command   |

Use `srediag help` or `srediag <group> --help` to explore available sub-commands and flags.

---

## Version

```bash
srediag version
```

Prints the agent’s CalVer, Git commit, Go and OTel versions.

---

## Getting Started

1. **Build** your collector and plugins:

   ```bash
   srediag build all
   ```

2. **Start** the service:

   ```bash
   srediag service start --config /etc/srediag/srediag.yaml
   ```

3. **Install & enable** a plugin:

   ```bash
   srediag plugin install --file vectorhashprocessor-0.3.0.tar.gz
   srediag plugin enable processor vectorhashprocessor
   ```

4. **Rotate** your TLS certificates:

   ```bash
   srediag security cert rotate --cert new.crt --key new.key --ca ca.pem
   ```

5. **Run** a system health check:

   ```bash
   srediag diagnose system
   ```

---

## Further Reading

- **Configuration Reference** – `docs/configuration/README.md`
- **Architecture Overviews** – `docs/architecture/README.md`
