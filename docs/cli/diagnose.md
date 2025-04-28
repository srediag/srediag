# `srediag diagnose` — On-Demand Host & Application Diagnostics

Every diagnostic area is delivered by a **binary plugin** with **CLI
scope** enabled.
If a command is missing, check that its plugin is loaded:

```bash
srediag plugin list --scope cli
```

You may enable a diagnostic plugin at runtime:

```bash
srediag plugin enable --scope cli diag k8sclusterdiagnostics
```

---

```bash
Usage:
  srediag diagnose <area> <command> [flags]
```

Global flags:

| Flag | Purpose | Default |
| :--- | :------ | :------ |
| `--output (json\|yaml\|table)` | Render style | `table` |
| `--quiet` | Suppress headings / timestamp | `false` |
| `--timeout <dur>` | Hard timeout per check | `30s` |
| `--format` | Alias of `--output` | — |

---

## 0 · Command Matrix & Plugin Mapping

| Area | Command | Backing plugin | Scope enabled by default* |
| :--- | :------ | :------------- | :------------------------ |
| **System** | `snapshot` | `systemsnapshot` | cli |
|  | `monitor` | `systemsnapshot` | cli |
| **Performance** | `cpu-prof` | `perfprofiler` | cli |
|  | `mem-prof` | `perfprofiler` | cli |
| **Security** | `cis-bench` | `cisbaseline` | cli |
| **Kubernetes** | `cluster` | `k8sclusterdiagnostics` | *opt-in* |
|  | `resources` | same | *opt-in* |
| **Network** | `latency` | `netlatencydiag` | *opt-in* |
| **Filesystem** | `inode-usage` | `fsmonitor` | *opt-in* |

\* *“opt-in” = plugin binary shipped but disabled by default to avoid heavy deps (kubectl, netperf, etc.).*

Activate with:

```bash
srediag plugin enable --scope cli diag fsmonitor
```

---

## 1 · System Diagnostics (`systemsnapshot`)

### 1.1 `snapshot`

One-shot health dump (CPU, mem, disk, net):

```bash
srediag diagnose system snapshot --output yaml
```

Flags:

| Flag | Description | Default |
| --- | --- | --- |
| `--resource <cpu\|mem\|disk\|net\|all>` | limit scope | all |
| `--detail <low\|med\|high>` | sample density | med |

### 1.2 `monitor`

Stream resource metrics every *interval*:

```bash
srediag diagnose system monitor --interval 2s
```

Stop with **Ctrl-C**.  Combine with `--output json` for piping.

---

## 2 · Performance Profiling (`perfprofiler`)

### 2.1 `cpu-prof`

```bash
srediag diagnose performance cpu-prof --duration 45s --output /tmp/cpu.pprof
```

Generates a pprof profile.

Flags: `--duration`, `--output`.

### 2.2 `mem-prof`

Same as above but heap profile.

---

## 3 · Security Checks (`cisbaseline`)

### 3.1 `cis-bench`

Runs a trimmed CIS 1.3 benchmark:

```bash
srediag diagnose security cis-bench --format table
```

Flags:

| Flag | Purpose | Default |
| :--- | :------ | :------ |
| `--level <1\|2>` | CIS level | 1 |
| `--include <regex>` | Only checks matching regex | — |
| `--exclude <regex>` | Skip checks | — |

Exit code **2** if any *warn / fail* findings.

---

## 4 · Kubernetes Diagnostics (`k8sclusterdiagnostics`)

> Requires `kubectl` on PATH and *kubeconfig* context.

Enable plugin:

```bash
srediag plugin enable --scope cli diag k8sclusterdiagnostics
```

### 4.1 `cluster`

```bash
srediag diagnose kubernetes cluster --context prod-us-east
```

Produces readiness of API server, etcd, nodes.

### 4.2 `resources`

```bash
srediag diagnose kubernetes resources --namespace prod --top 10
```

Shows top CPU / mem pods with YAML output.

---

## 5 · Network Diagnostics (`netlatencydiag`)

### 5.1 `latency`

```bash
srediag diagnose network latency --dest 1.1.1.1 --count 20
```

* Uses built-in ICMP (non-root via UDP fallback).
* Requires plugin enabled.

---

## 6 · Filesystem Diagnostics (`fsmonitor`)

### 6.1 `inode-usage`

```bash
srediag diagnose filesystem inode-usage --warning 85
```

Warns if any mount exceeds threshold.

---

## 7 · Output Examples

### Table (default)

```text
+-----------+--------------+---------+-------+
| RESOURCE  | USAGE (pct)  | STATUS  | NOTE  |
+-----------+--------------+---------+-------+
| CPU       | 43.1 %       | OK      |       |
| Memory    | 78.4 %       | WARN    | high  |
+-----------+--------------+---------+-------+
```

### JSON

```bash
srediag diagnose system snapshot --output json | jq .
```

```json
{
  "cpu": { "usage_pct": 43.1 },
  "memory": { "usage_pct": 78.4, "status": "warn" }
}
```

---

## 8 · Exit Codes (per-command)

| Code | Meaning |
| :--- | :------ |
| 0 | Success |
| 1 | Generic error |
| 2 | Findings above warn threshold |
| 3 | Missing dependency (kubectl, ping) |
| 4 | Plugin not enabled |
| 5 | Timeout |

---

## 9 · Tips & Best-Practice

* **Enable only what you need**—heavy plugins (K8s, network) stay off by default.
* **Automate** periodic snapshots with `cron` and `--output json`.
* Combine `diagnose … monitor` with `grep` or `jq` to feed alerts.
* Results feed into **OTel pipelines** when `diagotelprocessor`
  plugin is enabled (future roadmap).

---

## 10 · Related CLI Docs

* Plugin lifecycle — [cli/plugin.md](plugin.md)
* Service runtime — [cli/service.md](service.md)
