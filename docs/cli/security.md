# `srediag security` — TLS & Integrity Utility

> **Purpose** – certificate lifecycle and integrity-check helper for the
> SREDIAG binary and its plugins – nothing more.
> For runtime knobs (YAML) see `configuration/security.md`.
> For architectural rationale see `architecture/security.md`.

---

## 0 · Usage Synopsis

```bash
srediag security <sub-command> [flags]
```

Global flags:

| Flag | Purpose | Default |
| :--- | :------ | :------ |
| `--format <json\|yaml\|table>` | Output renderer | `table` |
| `--quiet` | Suppress non-error output | `false` |
| `--scope <sys\|user>` | Operate on system (root) or user store | auto |

---

## 1 · Command Matrix

| Category | Sub-command | Scope(s) | Description |
| :------- | :---------- | :------- | :---------- |
| **TLS / SPIFFE** | `cert show` | both | Print active leaf & CA details |
|                 | `cert rotate` | sys | Atomically swap cert/key/CA bundle |
|                 | `spiffe id`   | both | Show computed SPIFFE ID |
| **Integrity**   | `verify binary`  | both | SHA-256 + cosign on collector |
|                 | `verify plugin`  | both | Ditto for one plugin bundle |
| **Sandbox**     | `sandbox stats`  | both | Show seccomp / rlimit counters |
|                 | `sandbox test`   | both | Execute a harmless syscall probe |
| **Audit Helper**| `doctor`         | both | Quick health scan (expiry, sigs) |

---

## 2 · TLS & SPIFFE Commands

### 2.1 `cert show`

```bash
srediag security cert show --format yaml
```

```yaml
common_name: edge-agent.observo.local
spiffe_id:   spiffe://observo/agent/node-01
not_before:  2025-04-26T12:00:00Z
not_after:   2025-07-25T12:00:00Z
signature:   SHA256-ECDSA
```

`--chain` prints full PEM.

### 2.2 `cert rotate`

```bash
sudo srediag security cert rotate \
      --cert /etc/ssl/new.crt \
      --key  /etc/ssl/new.key \
      --ca   /etc/ssl/ca.pem
```

* Writes files to the directory configured by `security.tls.*`.
* Sends **SIGHUP** to running service; zero-downtime.

Exit codes: `0 ok`, `4 file-not-found`, `5 reload-failed`.

### 2.3 `spiffe id`

Returns deterministic SPIFFE ID derived from cert SAN or fallback
`service.name=node-id`.

---

## 3 · Binary & Plugin Verification

### 3.1 Collector binary

```bash
srediag security verify binary /usr/local/bin/srediag
```

Output (table default):

| Field | Value |
| :---- | :---- |
| SHA-256 | `f2c3…` |
| cosign | **OK** (keyless) |
| ABI | Go 1.24, OTel API v1.30.0 |
| SBOM digest | `sha256:47b9…` |
| Status | **valid** |

Exit `2` if any check fails.

### 3.2 Plugin bundle

```bash
srediag security verify plugin /usr/libexec/srediag/vectorhashprocessor.tar.gz
```

Checks manifest → SHA-256 → cosign (system scope only).

---

## 4 · Sandbox Diagnostics

### 4.1 `sandbox stats`

Prints counters kept by the Plugin Manager.

| Metric | Value |
| :----- | ----: |
| seccomp_blocks_total | 0 |
| rss_softlimit_hits   | 3 |
| cpu_throttle_events  | 1 |

### 4.2 `sandbox test`

Runs a short probe that tries blocked syscalls (`ptrace`) and reports
if the profile is active.

Exit `3` when the sandbox is disabled or mis-configured.

---

## 5 · `doctor` — Quick Health Scan

```bash
srediag security doctor --format json
```

Checks performed:

| Check | Severity on fail |
| :---- | :--------------- |
| TLS expiry < 14 d | **warn** |
| Binary unsigned / mismatching SHA | **error** |
| Plugin bundle unsigned | **error** |
| seccomp profile inactive | **error** |
| CA bundle absence | **warn** |

Exit codes:

| Code | Interpretation |
| ---: | :------------- |
| 0 | all good |
| 1 | only warnings |
| 2 | at least one error |

---

## 6 · Examples

```bash
# Nightly certificate rotation (cron)
sudo srediag security cert rotate \
     --cert /etc/letsencrypt/live/new.crt \
     --key  /etc/letsencrypt/live/new.key \
     --ca   /etc/letsencrypt/live/ca.pem

# CI step – verify bundle before publishing
srediag security verify plugin ./build/vectorhashprocessor.tar.gz

# Check sandbox health on a running node
srediag security sandbox stats
```

---

## 7 · Cross-Reference Index

| Topic | Document |
| :---- | :------- |
| Architecture & threat model | `architecture/security.md` |
| YAML security knobs | `configuration/security.md` |
| Plugin verification flow | `architecture/service-collector.md` |
| Build signing pipeline | `build/architecture.md` |
