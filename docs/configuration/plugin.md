# SREDIAG — Plugins (Binary Model)

All collector components (receivers, processors, exporters, extensions)
are delivered as **stand-alone binaries** launched by the agent.  
A component is usable only if its plugin is **enabled**.

---

## 1 · Paths & Scope

| Scope | Executables (`plugins.exec_dir`) | Config (`plugins.d/`) |
| :---- | :------------------------------ | :-------------------- |
| **System** (`root`, service unit) | `/usr/libexec/srediag` | `/etc/srediag/plugins.d` |
| **User** (`UID ≠ 0`) | `$HOME/.local/libexec/srediag` | `$HOME/.config/srediag/plugins.d` |

Override with `plugins.exec_dir` or `SREDIAG_PLUGINS_DIR`.

---

## 2 · Enabling Plugins

```yaml
plugins:
  enabled:
    - processor/vectorhashprocessor          # service + CLI
    - receiver/journaldreceiver@service      # service only
    - exporter/clickhouseexporter@cli        # CLI only
```

* `@service`, `@cli` suffix narrows scope.  
* Absence of suffix means **both** runtime contexts.

### CLI at runtime

```bash
# Session-only enable for CLI tools
srediag plugin enable --scope cli receiver journaldreceiver
```

---

## 3 · plugins.d/ — Default Settings

File: `plugins.d/<plugin>.yaml`

```yaml
# journaldreceiver.yaml
path: /var/log/journal
unit_filter: ["ssh.service", "nginx.service"]
```

*Merged before* Collector-YAML overrides.

---

## 4 · Multiple Instances

A single binary can host many *instances* differentiated by **alias**:

```yaml
processors:
  vectorhashprocessor/fast:      # aliases `fast` and `deep`
    window: 15s
  vectorhashprocessor/deep:
    window: 120s
    similarity: 92
```

Internally the plugin receives both configs and tags telemetry with
`processor.name`.

---

## 5 · Installation Patterns

### Package Manager

* Binaries → `/usr/libexec/srediag/`  
* YAML stub → `/etc/srediag/plugins.d/`  
* Checksums & signatures handled by the package.

### Tar.gz Manual

```bash
tar -C /usr/libexec/srediag -xzvf vectorhashprocessor-0.3.0-linux-amd64.tar.gz
echo "enabled: true" | sudo tee /etc/srediag/plugins.d/vectorhashprocessor.yaml
```

### User Experiment

```bash
mkdir -p ~/.local/libexec/srediag
cp ./mycpustatreceiver ~/.local/libexec/srediag/
srediag plugin enable --scope cli receiver mycpustatreceiver
```

---

## 6 · Security

| Layer | Policy |
| :---- | :----- |
| Checksum | SHA-256 verified against manifest or package DB |
| Signature | **cosign** bundles accepted for system scope; user scope warns if absent |
| Sandboxing | UID 65534, seccomp-strict (`clone3`, raw sockets blocked) |
| Resources | Sum RSS and CPU throttled under agent cgroup |

---

## 7 · Best Practices

1. **Immutable binaries** (`chmod 555`) in system dir.  
2. Keep `plugins.d/` under source-control next to collector YAML.  
3. Use user scope for dev; promote to system after review.  
4. Limit plugin exposure by attaching them only to required pipelines.

---

## 8 · Troubleshooting

| Symptom | Remedy |
| :------ | :----- |
| `binary not executable` | `chmod +x <binary>` |
| `checksum mismatch` | Re-download artefact; cross-check release manifest |
| Component “not found” during reload | Add plugin to `plugins.enabled` or correct alias spelling |

Enable debug logs:

```bash
SREDIAG_LOG_LEVEL=debug srediag service --config /etc/srediag/srediag.yaml
```

---

## 9 · References

* [Build Configuration](build.md)  
* [Service / Collector](service.md)  
* [CLI guide — plugin commands](../cli/README.md)
