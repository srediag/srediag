# SREDIAG — Plugins Configuration

Every Collector component — receiver, processor, exporter, extension, or
diagnostic helper — is delivered as an **individual ELF binary**.
The agent launches each binary in its own sandbox and keeps it alive
over the agent's control socket.

A component becomes usable only when its plugin is **enabled** for at
least one **scope**:

| Scope | Who uses it | Traffic path |
| :---- | :---------- | :----------- |
| **service** | Long-lived daemon (embedded Collector) | OTel pipeline |
| **cli** | Ad-hoc `srediag …` commands | stdout / terminal |

---

## 0 · Directory Layout & Scope Detection

| Scope | Executables (`plugins.exec_dir`) | Config (`plugins.d/`) |
| :---- | :------------------------------ | :-------------------- |
| **System mode** (unit runs as `root` or service user) | `/usr/libexec/srediag` | `/etc/srediag/plugins.d` |
| **User mode** (non-root, `--user`) | `$HOME/.local/libexec/srediag` | `$HOME/.config/srediag/plugins.d` |

Override paths with either:

```yaml
plugins:
  exec_dir: /opt/srediag/plugins
```

or environment:

```bash
export SREDIAG_PLUGINS_DIR=/opt/srediag/plugins
```

The agent auto-detects **scope**:

* *System* when EUID == 0 **and** `--user` flag not present.
* *User* otherwise.

---

## 1 · Default-Enabled Plugins (“Zero-config startup”)

### 1.1 System scope (Collector)

| Component ID | Binary | Purpose |
| :----------- | :----- | :------ |
| `receiver/otlpreceiver`           | `otlpreceiver` | OTLP ingest |
| `receiver/nopreceiver`            | `nopreceiver`  | Dummy receiver for boot |
| `processor/batchprocessor`        | `batchprocessor` | Overload smoothing |
| `processor/memorylimiterprocessor`| `memorylimiterprocessor` | RSS guard |
| `exporter/otlpexporter`           | `otlpexporter` | OTLP egress |
| `extension/healthcheckextension`  | `healthcheckextension` | `/healthz` |
| `extension/zpagesextension`       | `zpagesextension` | In-process trace UI |

### 1.2 CLI scope (user diagnostics)

| Command family | Binary | Size |
| :------------- | :----- | ---: |
| `system snapshot / monitor`  | `systemsnapshot`  | 2 MiB |
| `performance cpu-prof / mem-prof` | `perfprofiler` | 4 MiB |
| `security cis-bench` | `cisbaseline` | 6 MiB |

### 1.3 Optional (shipped but **OFF**)

| Binary | Typical flags | Why not default? |
| :----- | :------------ | :--------------- |
| `k8sclusterdiagnostics` | `diagnose kubernetes …` | Requires `kubectl`, heavy |
| `netlatencydiag` | `diagnose network latency …` | Raw sockets |
| `clickhouseexporter` | Pipeline exporter | External DB creds |

Enable on-demand:

```bash
sudo srediag plugin enable --scope service exporter clickhouseexporter
sudo srediag service reload
```

Disable default permanently:

```yaml
plugins:
  enabled:
    - extension/zpagesextension@service=false
```

---

## 2 · Plugin Lifecycle via CLI

| Task | Command | Notes |
| :--- | :------ | :---- |
| **Install** new binary | `srediag plugin install --file *.tar.gz` | Verifies SHA-256 & cosign |
| **Enable** in scope(s) | `srediag plugin enable --scope cli\|service \<type\> \<name\>` | Registers + hot-reloads |
| **Disable** temporarily | `srediag plugin disable …` | Persists only if YAML is edited |
| **Persist defaults** | `srediag plugin set <name> key value` | Writes to `plugins.d/` |
| **Bind** to pipeline | `srediag plugin bind attach … <alias>` | Creates alias if absent |
| **Reload** binary | `srediag plugin reload <name>` | Zero-drop hand-off |
| **Update** to new ver. | `srediag plugin update <name> --to 0.4.1` | `doctor` validates afterwards |
| **Uninstall** binary | `srediag plugin uninstall <name>` | Fails if active |

All lifecycle commands respect `--scope` and `--exec-dir`.

---

## 3 · `plugins.enabled:` YAML Grammar

```yaml
plugins:
  enabled:
    # Enabled for both scopes
    - processor/vectorhashprocessor
    # Service-only
    - receiver/journaldreceiver@service
    # CLI-only
    - diag/perfprofiler@cli
    # Explicit disable of default
    - extension/zpagesextension@service=false
```

* Suffixes: `@service`, `@cli`, `@both`, or `@service=false`.
* Wildcard prefix allowed: `receiver/*@service=false` (disable all receivers).

Precedence (high → low):

```bash
CLI flags   >  plugins.enabled (YAML)  >  Builder defaults
```

---

## 4 · `plugins.d/` — Per-Plugin Defaults

Merged order inside a plugin instance:

```bash
Collector YAML (highest)  >  plugins.d/<name>.yaml  >  plugin hard-coded defaults
```

Example **system scope** file:

```yaml
# /etc/srediag/plugins.d/vectorhashprocessor.yaml
window: 45s
max_cache_mib: 128
similarity: 92
```

User scope path:

```bash
$HOME/.config/srediag/plugins.d/vectorhashprocessor.yaml
```

---

## 5 · Multiple Instances (Aliases)

```yaml
processors:
  vectorhashprocessor/fast:
    window: 15s
  vectorhashprocessor/deep:
    window: 120s
    similarity: 92
```

CLI binding:

```bash
sudo srediag plugin bind attach processor vectorhashprocessor fast logs/default
```

Aliases may also be **disabled**:

```bash
sudo srediag plugin bind detach vectorhashprocessor deep logs/default
```

---

## 6 · Installation Patterns

### 6.1 Package manager (rpm / deb)

* Binaries → `/usr/libexec/srediag/`
* YAML stub → `/etc/srediag/plugins.d/`
* Checksums & sig = package manager responsibility.

### 6.2 Tarball (manual)

```bash
sudo tar -C /usr/libexec/srediag -xzf vectorhashprocessor-0.4.0.tar.gz
echo "enabled: true" | sudo tee /etc/srediag/plugins.d/vectorhashprocessor.yaml
sudo srediag plugin verify vectorhashprocessor
sudo srediag plugin enable --scope service processor vectorhashprocessor
sudo srediag service reload
```

### 6.3 Developer experiment (user scope)

```bash
mkdir -p ~/.local/libexec/srediag
cp ./mycpustatreceiver ~/.local/libexec/srediag/
srediag plugin enable --scope cli receiver mycpustatreceiver
```

---

## 7 · Security Hardening

| Layer | Enforcement |
| :---- | :---------- |
| **Checksum** | SHA-256 verified against build manifest or `*.sha256` file |
| **Signature** | `cosign verify-blob`; system scope **fails** on mismatch, user scope warns |
| **Sandbox** | UID 65534 (`nobody`), `seccomp-strict`, mount namespace, no raw sockets |
| **Cgroups** | Memory & CPU charged to agent slice; overridable via `plugins.resources` |

---

## 8 · Troubleshooting Cheats

| Symptom | Likely cause / remedy |
| :------ | :-------------------- |
| *"component not found"* during reload | Plugin not **enabled** in that scope or alias spelled wrong |
| *checksum mismatch* | Download correct artefact & verify cosign |
| Plugin stuck in `error` state | Crash-loop; inspect `journalctl -u srediag` or `~/.config/srediag/logs/` |
| `bind` fails | Alias already attached elsewhere; detach first |

Verbose logging:

```bash
SREDIAG_LOG_LEVEL=debug srediag plugin list --scope service
```

---

## 9 · Quick-reference of CLI Commands

| Command | Short-hand |
| :------ | :--------- |
| Inventory | `list`, `paths` |
| Inspection | `info`, `verify` |
| Lifecycle | `enable`, `disable`, `reload` |
| Config edit | `set`, `unset`, `bind attach/detach` |
| Distribution | `install`, `update`, `uninstall` |
| Maintenance | `doctor` |

See full syntax in [CLI manual](../cli/plugin.md).

---

## 10 · Further Reading

* Build & packaging — [build.md](build.md)
* Service pipelines — [service.md](service.md)
* Diagnose commands — [cli/diagnose.md](../cli/diagnose.md)

---

## Parameter Reference

| YAML Key           | Env Var                    | CLI Flag         |
|--------------------|---------------------------|------------------|
| `plugins.dir`      | `SREDIAG_PLUGINS_DIR`      | `--plugins-dir`  |
| `plugins.enabled`  | —                         | `plugin enable`  |
| `plugins.exec_dir` | `SREDIAG_PLUGINS_EXEC_DIR` | `--exec-dir`     |
| `srediag.config`   | `SREDIAG_CONFIG`           | `--config`       |

> **Warning:** Do **not** use `--config` for plugin-specific settings; this is reserved for the main SREDIAG config. Use the above flags/envs for plugin configuration.

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
* Always check the effective config with `srediag plugin --print-config` (if available).
* Unknown YAML keys are logged at debug level and ignored for forward compatibility.
