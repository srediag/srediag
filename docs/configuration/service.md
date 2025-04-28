# Service-mode & Collector Configuration

| YAML Key                | Env Var                        | CLI Flag         |
|-------------------------|--------------------------------|------------------|
| `srediag.config`        | `SREDIAG_CONFIG`               | `--config`       |
| `collector.config_path` | `SREDIAG_COLLECTOR_CONFIG_PATH`| `--service-yaml` |
| `service.port`          | `SREDIAG_SERVICE_PORT`         | `--service-port` |
| `service.name`          | `SREDIAG_SERVICE_NAME`         | `--service-name` |
| `collector.enabled`     | `SREDIAG_COLLECTOR_ENABLED`    | `--collector-enabled` |

> **Warning:** Do **not** use `--config` for service/collector YAML; this is reserved for the main SREDIAG config. Use `--service-yaml`/`SREDIAG_COLLECTOR_CONFIG_PATH` for collector pipeline configuration.

When you run `srediag service` the agent:

1. Boots **core services** (logging, healthz, plugin manager).  
2. Starts an embedded **OpenTelemetry Collector** (if
   `collector.enabled: true`).

> **Important** — A receiver / processor / exporter **MUST** exist as an
> enabled *plugin* before it can be referenced in `srediag-service.yaml`.

---

## 1 · Minimal Core YAML

```yaml
collector:
  enabled: true
  config_path: /etc/srediag/srediag-service.yaml
  memory_limit_mib: 1024
```

---

## 2 · Default Plugin Set (shipped with SREDIAG)

| Component | Type | Plugin Name | Enabled by default |
| :-------- | :--- | :---------- | :----------------- |
| OTLP gRPC / HTTP | receiver | `otlpreceiver` | ✓ |
| No-op receiver   | receiver | `nopreceiver`  | ✓ |
| Batch            | processor| `batchprocessor`| ✓ |
| MemoryLimiter    | processor| `memorylimiterprocessor` | ✓ |
| OTLP exporter    | exporter | `otlpexporter` | ✓ |
| HealthCheck      | extension| `healthcheckextension` | ✓ |
| zPages           | extension| `zpagesextension` | ✓ |

These binaries live in **`/usr/libexec/srediag/`** and are declared in
`plugins.enabled` out-of-the-box.

---

## 3 · Referencing Plugins in `srediag-service.yaml`

```yaml
receivers:
  otlp/default:            # first instance
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
  otlp/from_gateway:       # second instance (alias)
    protocols:
      http:
        endpoint: 0.0.0.0:4318

processors:
  batch/default:
  memory_limiter/main:
    limit_mib: 1024
    check_interval: 2s
  vectorhashprocessor/high_card: {}   # requires plugin enabled

exporters:
  otlp/gateway:
    endpoint: gw.observo.local:4317
```

### Multiple Instances

* Use **`/<alias>`** after the component name.  
* Each alias gets its own configuration block.  
* The plugin binary is started **once** and multiplexes instances.

---

## 4 · Enabling Extra Plugins

1. **System-wide & Persisted**

   ```yaml
   # srediag.yaml
   plugins:
     enabled:
       - processor/vectorhashprocessor        # service + CLI
       - receiver/journaldreceiver@service    # only service-mode
   ```

   *`@service`* or *`@cli`* suffix scopes the plugin. Omit for both.

2. **Ad-hoc (non-persistent) in a shell**

   ```bash
   srediag plugin enable --scope cli receiver journaldreceiver
   ```

   Scope options: `cli`, `service`, `both` (default).

---

## 5 · Hot-Reload & Validation

* Send **`SIGHUP`** → core YAML and collector YAML reload.  
* If the collector references a plugin not yet enabled, reload fails and
  the previous config stays active (logged at `error` level).  
* Use  
  `srediag validate collector --file srediag-service.yaml` *(planned)* to
  pre-check references.

---

## 6 · Troubleshooting

| Issue | Fix |
| :---- | :-- |
| `component not found: vectorhashprocessor` | Add `processor/vectorhashprocessor` to `plugins.enabled` |
| Reload fails after YAML edit | Check plugin scope; ensure alias names match `type/alias` syntax |
| High RSS | Lower `memory_limit_mib` or adjust MemoryLimiter settings |

---

## 7 · Related Docs

* [Plugins (binary) Guide](plugins.md)  
* [Build Configuration](build.md)  
* [CLI Reference](../cli/README.md)

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
