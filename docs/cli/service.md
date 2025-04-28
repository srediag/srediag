# `srediag service` — Operating the SREDIAG Agent

`service` manages the long-lived SREDIAG daemon:

* Core runtime (logging, plugin-mgr, HTTP endpoints)  
* Embedded **OpenTelemetry Collector** (when `collector.enabled: true`)

Two distinct **modes** are supported:

| Mode | Who runs it | Files & ports | Use-case |
| :--- | :---------- | :------------ | :------- |
| **System mode** | `root` (init system) | `/etc/srediag/…`, port 8080 | Fleet / prod nodes |
| **User mode**   | Non-root user | `$HOME/.config/srediag/…`, port 808x | Dev laptops, per-user tests |

A single host may run both; the control socket keeps namespaces
separate.

```bash

Usage:
  srediag service <sub-command> [flags]

```

---

## 0 · Command Matrix

| Category | Sub-cmd | System | User | Description |
| :------- | :------ | :----: | :--: | :---------- |
| **Lifecycle** | `start` (default) | ✓ | ✓ | Foreground unless `--detach` |
| | `stop` | ✓ | ✓ | Graceful shutdown via PID / control socket |
| | `restart` | ✓ | ✓ | Stop → Start |
| | `reload` | ✓ | ✓ | Hot-reload both YAML files |
| | `detach` | ✓ | ✓ | Fork to background (Unix) |
| **Introspection** | `status` | ✓ | ✓ | Health snapshot, resource usage |
| | `health` | ✓ | ✓ | Exit 0 if `/healthz` ready |
| | `profile` | ✓ | ✓ | Gather CPU+heap profile bundle |
| | `tail-logs` | ✓ | ✓ | Stream live service logs |
| **Validation** | `validate` | ✓ | ✓ | Dry-run parse YAML + plugin refs |
| **systemd Helper** | `install-unit` | ✓ | — | Create & enable system unit |
| | `uninstall-unit` | ✓ | — | Remove unit |
| **Maintenance** | `gc` | ✓ | ✓ | Purge stale PID/socket/logs |

---

## 1 · Global Flags

| Flag | Purpose | Default (sys) | Default (user) |
| :--- | :------ | :------------ | :------------- |
| `--config <file>` | Core YAML (`srediag.yaml`) | `/etc/srediag/srediag.yaml` | `$XDG_CONFIG_HOME/srediag/srediag.yaml` |
| `--service-yaml <file>` | Collector YAML | `/etc/srediag/srediag-service.yaml` | `$XDG_CONFIG_HOME/srediag/srediag-service.yaml` |
| `--mem-limit <MiB>` | Override RSS guard | value from YAML | value from YAML |
| `--timeout <dur>` | Wait time for stop/reload | `30s` | `30s` |
| `--detach` | Background daemonise | `false` | `false` |
| `--user` | **Force user mode**, even if root | `false` | auto |

*A command runs in **user mode** when `--user` is present **or**
effective UID ≠ 0.*

---

## 2 · Lifecycle Reference

### 2.1 `start`

System mode (foreground):

```bash
sudo srediag service start --config /etc/srediag/srediag.yaml
```

User mode (background):

```bash
srediag service start --user --detach
```

Artifacts created:

| Item | System path | User path |
| :--- | :---------- | :-------- |
| PID file | `/run/srediag.pid` | `$XDG_RUNTIME_DIR/srediag.pid` |
| Control socket | `/run/srediag.sock` | `$XDG_RUNTIME_DIR/srediag.sock` |
| HTTP port | `service.port` (8080) | next free 8080 + UID%100 |

### 2.2 `stop / restart`

`stop` sends **SIGTERM** via control socket, waits `--timeout`.

### 2.3 `reload`

Re-reads both YAMLs → validates → hot-swaps pipelines.  
Rollback on failure; exit **2** if invalid.

---

## 3 · Introspection

### 3.1 `status`

```bash
srediag service status --format yaml --user
```

YAML fields: `state`, `uptime`, `rss_mib`, `cpu_pct`,
`plugins.active`, `collector.pipelines`, etc.

### 3.2 `profile`

Collects **30 s** CPU + heap + goroutine to
`/tmp/srediag-profile-<ts>.zip` (root) or `$TMPDIR`.

---

## 4 · Validation

```bash
srediag service validate --config my.yaml --service-yaml pipeline.yaml
```

Checks syntax + plugin presence + alias uniqueness.

---

## 5 · systemd Integration

### 5.1 System unit

```bash
sudo srediag service install-unit \
     --config /etc/srediag/srediag.yaml \
     --service-yaml /etc/srediag/srediag-service.yaml
```

Creates `/etc/systemd/system/srediag.service`:

```ini
[Service]
Type=simple
User=srediag
ExecStart=/usr/local/bin/srediag service start --config %f
RuntimeDirectory=srediag
AmbientCapabilities=CAP_NET_BIND_SERVICE
ReadOnlyPaths=/
```

Enable & control via **systemctl**:

```bash
sudo systemctl enable --now srediag
sudo systemctl reload srediag        # ↔ srediag service reload
sudo systemctl status srediag
```

### 5.2 User unit (linger-enabled)

```bash
srediag service install-unit --user
systemctl --user enable --now srediag
```

---

## 6 · Garbage Collection (`gc`)

Purges:

* orphan PID files / sockets  
* logs > `--retention` (default 14 d)

```bash
sudo srediag service gc --retention 7d
```

---

## 7 · Default Collector Components

Activated in **system mode** unless disabled:

| Component ID | Plugin binary |
| :----------- | :------------ |
| `receiver/otlpreceiver` | `otlpreceiver` |
| `processor/batchprocessor` | `batchprocessor` |
| `processor/memorylimiterprocessor` | `memorylimiterprocessor` |
| `exporter/otlpexporter` | `otlpexporter` |
| `extension/healthcheckextension` | `healthcheckextension` |
| `extension/zpagesextension` | `zpagesextension` |

Disable permanently:

```yaml
plugins:
  enabled:
    - extension/zpagesextension@service=false   # boolean shorthand
```

or at runtime:

```bash
sudo srediag plugin disable --scope service zpagesextension
sudo srediag service reload
```

---

## 8 · Cheat-Sheet

| Task | System mode | User mode |
| :--- | :---------- | :-------- |
| Start | `systemctl start srediag` | `srediag service start --user --detach` |
| Hot-reload YAMLs | `srediag service reload` | same |
| Verify health (probe) | `srediag service health` | `srediag service health --user` |
| Grab profile | `sudo srediag service profile --output /tmp/a.zip` | `srediag service profile --user` |
| Remove daemon | `srediag service uninstall-unit && systemctl daemon-reload` | `systemctl --user disable --now srediag` |

---

## 9 · Exit Codes

| Code | Meaning |
| :--- | :------ |
| 0 | Success |
| 1 | Generic error |
| 2 | Validation failure (reload/validate) |
| 3 | Permission denied / root required |
| 4 | Daemon not running |
| 5 | Timeout waiting for action |

---

## 10 · Related Pages

* Collector & YAML — [configuration/service.md](../configuration/service.md)  
* Plugin lifecycle — [cli/plugin.md](plugin.md)  
* Build pipeline — [cli/build.md](build.md)

---

## Parameter Reference

| YAML Key                | Env Var                    | CLI Flag                |
|-------------------------|----------------------------|-------------------------|
| `service.port`          | `SREDIAG_SERVICE_PORT`     | `--service-port`        |
| `service.name`          | `SREDIAG_SERVICE_NAME`     | `--service-name`        |
| `collector.enabled`     | `SREDIAG_COLLECTOR_ENABLED`| `--collector-enabled`   |
| `collector.config_path` | `SREDIAG_COLLECTOR_CONFIG_PATH` | `--service-yaml`   |

> **Warning:** Do **not** use `--config` for service/collector YAML; this is reserved for the main SREDIAG config. Use `--service-yaml`/`SREDIAG_COLLECTOR_CONFIG_PATH` for collector pipeline configuration.

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
* Always check the effective config with `srediag service --print-config` (if available).
* Unknown YAML keys are logged at debug level and ignored for forward compatibility.
