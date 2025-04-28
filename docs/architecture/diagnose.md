# SREDIAG — Diagnostics Runtime Architecture

## 0 · Layered Model

| Layer | Location | Responsibility |
| :---- | :-------- | :------------- |
| **L0 CLI Shell** | `cmd/srediag/commands/diagnostic.go` | Top-level `srediag diagnose …` command factory |
| **L1 Diagnostic Core** | `internal/diagnostic` | Shared helpers (procfs, psutil adapter, perf events) |
| **L2 Diag Plugin Loader** | `internal/plugin` (scope **cli**) | Verify → sandbox → expose Cobra sub-cmds |
| **L3 Diag Plugins** | `$SREDIAG_PLUGINS_DIR/*`  `~/…` | Real logic (systemsnapshot, perfprofiler, cisbaseline, …) |

Data never flows upward; plugins call helpers but not the CLI.

---

## 1 · Plugin Contract (SDK §5.1 superset)

```go
type DiagPlugin interface {
    diagnostics.Plugin                // Info(), Init(), Shutdown(), Health()
    Register(root *cobra.Command)     // inject sub-commands
    Capabilities() []string           // e.g. ["diag/system", "diag/perf/cpu"]
}
```

* `Register` receives a **private** Cobra branch; plugins must not touch
  global flags.
* `Capabilities` are exported to **RBAC** (viewer / operator / admin).

---

## 2 · Execution Flow

```mermaid
flowchart TD
   C[Cobra root\n"srediag"] --> D[diagnose handler]
   D --> M(PluginMgr.load(cli))
   M -->|inject cmds| P1(P systemsnapshot)
   M --> P2(P perfprofiler)
   click P1 callback "Run(ctx)" "runs inside host\nseccomp(default)"
```

* Plugins launch **inside the main process** (no fork) to keep “one
  binary” UX.
* Heavy collectors (perf events, flamegraphs) *may* spawn helpers under
  `internal/diagnostic/cmdhelper` with their own seccomp profile.

---

## 3 · Streaming & Telemetry

| Mode | Transport | Use-case |
| :--- | :-------- | :------- |
| **STDOUT** | table / JSON / YAML | Default human output |
| **OTLP**   | `diag_exporter` internal processor | Fleet automation (Phase 3) |
| **File**   | `--output file.csv` | Air-gapped export |

`diag_exporter` is a **service-scope** plugin that converts diagnostic
events into OTLP Logs and pushes them to the CP for CMDB correlation.

---

## 4 · Sandbox & Resource Caps

| Guard | Default | Notes |
| :---- | :------ | :---- |
| **seccomp** | `runtime/default` | perfprofiler adds `perf_event_open` |
| **mem** | 128 Mi soft | counts towards agent `mem_limit_mib` |
| **cpu** | 1 core weight | sampled profilers may raise to 2 × target |
| **file IO** | RO FS + cwd | systemsnapshot needs `/proc`, `/sys` (RO) |

Violations increment `srediag_diag_sandbox_violations_total`.

---

## 5 · Built-in Plugin Matrix (2025.05.0)

| ID | Capability | Sources Read | Output |
| :-- | :--------- | :----------- | :----- |
| `systemsnapshot` | `diag/system` | procfs, sysfs, cgroup v2 | JSON / OTLP |
| `perfprofiler` | `diag/perf/*` | `perf_event_open`, BPF counters | pprof file |
| `cisbaseline` | `diag/security` | `/usr/bin/lynis` adapter | table / OTLP |

All three ship in the default tar.gz and are enabled **cli** scope.

---

## 6 · Control-Plane Feedback Loop (Phase 3)

1. Operator pushes `profile.cpu duration=30s` command via signed
   remote-config.
2. Agent spawns **perfprofiler** with arguments.
3. Resulting `cpu.pb.gz` uploaded through
   `diag_exporter` → Control-Plane storage bucket.
4. CP annotates CMDB asset with `diag.result=OK|WARN|FAIL`.

The mechanism re-uses remote-config digest logic described in the
Collector architecture (§2).

---

## 7 · Error & Exit-Code Semantics

| Scope | Code | Meaning |
| :---- | ---: | :------ |
| Plugin Run() | non-nil error | propagates to Cobra, exit 1 |
| Manager Load | 4 | plugin not found |
| Verification | 2 | SHA-256 / cosign failed |
| Sandbox | 3 | seccomp / rlimit triggered |
| Timeout | 5 | `--timeout` exceeded |

---

## 8 · Metrics Contract

| Metric | Type | Labels | Meaning |
| :----- | :--- | :----- | :------ |
| `srediag_diag_runs_total` | counter | `plugin`, `result` | finished jobs |
| `srediag_diag_duration_seconds` | histogram | `plugin` | wall-clock |
| `srediag_diag_sandbox_violations_total` | counter | `plugin`, `syscall` | seccomp hit |

Plugins **must** emit these via the SDK helper
`diagmetrics.NewRun(ctx, plugin)`.

---

## 9 · Cross-Document Index

| Topic | Go to |
| :---- | :---- |
| CLI user guide & examples | `docs/cli/diagnose.md` |
| Plugin manifest & exec paths | `configuration/plugins.md` |
| Collector & plugin hot-swap | `architecture/service-collector.md` |
| Build pipeline | `build/architecture.md` |
| Security hardening | `configuration/security.md` |
