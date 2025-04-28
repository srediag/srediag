# SREDIAG — Diagnostics Configuration

This guide explains how to configure diagnostic commands, plugins, and runtime settings within SREDIAG. The diagnostic subsystem provides real-time insights and health-check capabilities, complementing the core Collector and plugin systems.

**Scope:**

- CLI-driven diagnostic commands
- Diagnostic plugin configurations
- Runtime diagnostic execution settings

**See also:**

- [Diagnostics Architecture](../architecture/diagnose.md)
- [Plugin Configuration](plugins.md)
- [CLI Diagnostics Usage](../cli/diagnose.md)

---

## 1 · Configuration Files & Directories

Diagnostics configurations are managed through:

| Scope        | Path                                          | Description                                             |
|--------------|-----------------------------------------------|---------------------------------------------------------|
| **System**   | `/etc/srediag/diagnostics.yaml`               | System-wide diagnostic default configurations           |
| **User**     | `$HOME/.config/srediag/diagnostics.yaml`      | User-specific overrides for diagnostic runs             |
| **Plugins**  | `/etc/srediag/plugins.d/<plugin>.yaml`        | Plugin-specific diagnostic settings (merged at runtime) |

Configurations follow the standard YAML hierarchy and support overrides via CLI flags or environment variables (`SREDIAG_DIAG_*`).

---

## 2 · Example YAML Configuration

A representative diagnostic configuration looks like this:

```yaml
diagnostics:
  defaults:
    output_format: json      # Default output format (json/yaml/table)
    timeout: 30s             # Default timeout per diagnostic command

  plugins:
    systemsnapshot:
      resources: [cpu, memory, disk]
      detail_level: high
    perfprofiler:
      default_duration: 60s
      max_duration: 300s
```

- **Global defaults** apply unless overridden by CLI flags.
- **Plugin-specific sections** configure individual plugin behavior.

---

## 3 · Configuration Reference

### 3.1 · Global Diagnostic Defaults

| Parameter      | Type     | Default   | Description                                      |
|----------------|----------|-----------|--------------------------------------------------|
| `output_format`| string   | `table`   | Default diagnostic output format (`json/yaml/table`)|
| `timeout`      | duration | `30s`     | Default timeout applied to all diagnostic runs   |
| `max_retries`  | integer  | `1`       | Retries per diagnostic command if initial run fails|

### 3.2 · `systemsnapshot` Plugin

| Parameter      | Type     | Default                 | Description                                          |
|----------------|----------|-------------------------|------------------------------------------------------|
| `resources`    | []string | `[cpu, memory, disk, net]` | Resources included by default in snapshots        |
| `detail_level` | string   | `medium`                | Snapshot detail: `low`, `medium`, or `high`          |

### 3.3 · `perfprofiler` Plugin

| Parameter          | Type     | Default  | Description                                      |
|--------------------|----------|----------|--------------------------------------------------|
| `default_duration` | duration | `30s`    | Default profiling duration                       |
| `max_duration`     | duration | `300s`   | Maximum allowed profiling duration               |
| `allowed_users`    | []string | `[]`     | Users permitted to run profiler (empty means no restrictions)|

### 3.4 · `cisbaseline` Plugin

| Parameter          | Type      | Default   | Description                                      |
|--------------------|-----------|-----------|--------------------------------------------------|
| `level`            | integer   | `1`       | CIS Benchmark Level (`1` or `2`)                 |
| `include`          | []string  | `[]`      | Specific checks to include (by IDs)              |
| `exclude`          | []string  | `[]`      | Checks to explicitly exclude                     |

---

## 4 · Runtime Overrides (CLI & ENV)

CLI flags always override YAML configurations, ensuring flexibility:

```bash
# Run diagnostic with explicit timeout and format
srediag diagnose system snapshot --timeout 60s --format json
```

Environment variables (`SREDIAG_DIAG_*`) are also supported:

```bash
export SREDIAG_DIAG_OUTPUT_FORMAT=json
export SREDIAG_DIAG_TIMEOUT=90s
```

---

## 5 · Diagnostic Plugins & Plugin Manager Integration

Diagnostic commands are implemented as plugins that integrate seamlessly into the SREDIAG Plugin Manager:

- Plugins register their diagnostic commands into the CLI (`srediag diagnose`).
- Each plugin includes default configurations and can be further configured in `plugins.d/<plugin>.yaml`.
- Plugins run in an isolated sandbox environment (seccomp/AppArmor policies enforced), leveraging the standard plugin runtime architecture.

**Example plugin YAML configuration:**

`/etc/srediag/plugins.d/systemsnapshot.yaml`

```yaml
resources: [cpu, memory]
detail_level: high
```

Plugin-level settings override global diagnostic defaults.

---

## 6 · Security & Permissions

Diagnostics execution respects security boundaries defined in the global configuration and security policy (`security.yaml`):

- **RBAC:** Diagnostic capabilities map to specific RBAC verbs (e.g., `diag:read`, `diag:execute`).
- **Sandboxing:** Diagnostics run under strict seccomp and AppArmor profiles, ensuring minimal privilege.
- **Permissions:** File-based configurations can specify allowed users or groups for sensitive diagnostics (`perfprofiler`).

---

## 7 · Observability & Metrics

Diagnostic plugins report metrics via the standard SREDIAG observability interfaces (OTel metrics):

| Metric Name                             | Type    | Description                                          |
|-----------------------------------------|---------|------------------------------------------------------|
| `srediag_diag_runs_total{plugin, status}`| counter | Counts diagnostic runs (labels: `plugin`, `status`) |
| `srediag_diag_duration_seconds{plugin}`  | histogram | Duration of each diagnostic run                     |
| `srediag_diag_errors_total{plugin}`      | counter | Diagnostic execution errors                         |

These metrics help track diagnostics health, performance, and failure rates.

---

## 8 · Troubleshooting

Common diagnostic runtime troubleshooting steps:

| Symptom                      | Potential Causes                        | Recommended Action                                     |
|------------------------------|-----------------------------------------|--------------------------------------------------------|
| Diagnostic timeout           | Insufficient timeout                    | Increase `timeout` via CLI or YAML                     |
| Unexpected output format     | Incorrect default settings              | Explicitly set `--format` on CLI or ENV                |
| Permission denied (profiler) | User restrictions                       | Verify `allowed_users` in `perfprofiler` config        |
| Plugin missing / not enabled | Plugin Manager configuration issue      | Check `plugins.enabled` in main `srediag.yaml`         |

---

## 9 · Cross-Reference & Further Reading

| Document                                | Details                                 |
|-----------------------------------------|-----------------------------------------|
| [Diagnostics Architecture](../architecture/diagnose.md)| Architecture and runtime details          |
| [Plugin Configuration](plugins.md)      | Plugin manager settings                 |
| [Security Configuration](security.md)   | RBAC, sandbox, and permissions          |
| [CLI Diagnostics Usage](../cli/diagnose.md)| CLI command syntax and examples         |
| [Collector Service Config](service.md)  | Service runtime configurations          |

---

## 10 · Governance & Revision Tracking

- Maintained under MIT License.
