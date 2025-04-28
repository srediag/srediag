# SREDIAG — Security Configuration Guide

> **Audience** – Platform / SRE engineers operating the SREDIAG agent in
> dev, staging, and production.
> **Out of scope** – Threat-model & crypto architecture (see
> `docs/architecture/security.md`) and supply-chain controls
> (`docs/security/supply-chain.md`).

SREDIAG enforces security along **four concentric rings**:

```ascii

┌──────────────────────────────────────────────┐
│  4  Supply-chain   (SBOM, cosign, SLSA)      │
│  3  Runtime guard  (seccomp, cgroup, RO FS)  │
│  2  Control plane  (TLS, AuthN, RBAC, quotas)│
│  1  Data plane     (mTLS, dedup, PII scrubs) │
└──────────────────────────────────────────────┘

```

This guide shows **how to configure ring 2-3** via YAML and runtime
flags.

---

## 0 · File & Flag Locations

| Scope | YAML key | Default path |
| :---- | :------- | :----------- |
| **Service daemon** | `security:` | `/etc/srediag/srediag.yaml` |
| **CLI runtime** | `security.cli:` (*subset*) | `$HOME/.config/srediag/srediag.yaml` |
| **Per-plugin tweaks** | `plugins.d/<name>.yaml` → `security:` | Same dir as plugin YAML |
| **Flags** | `--tls-*`, `--auth-*` | override YAML at startup |

> **Tip** – *System-mode* flags start with `srediag service …`;
> *user-mode* flags start with `srediag … --user`.

---

## 1 · Complete YAML Reference (system scope)

```yaml
security:
  tls:
    enabled: true
    cert_file: /etc/srediag/certs/server.crt
    key_file:  /etc/srediag/certs/server.key
    ca_file:   /etc/srediag/certs/ca.crt
    min_version: TLS1.3          # 1.2 allowed for legacy if explicit
    verify_client: true          # mTLS

  auth:               # Collector & CLI invoke same backend
    type: oauth2      # oauth2 | jwt | none
    jwt:
      secret:   "${JWT_SECRET}"
      lifetime: 24h
    oauth2:
      issuer:  https://login.example.com
      client_id: srediag
      client_secret: "${OIDC_SECRET}"
      scopes: [openid, profile, email]

  rbac:               # Shared by service HTTP UI & plugin IPC
    enabled: true
    default_role: viewer
    roles:
      admin:    ["*"]
      operator: ["read:*", "write:config", "write:telemetry"]
      viewer:   ["read:*"]

  quotas:             # per-tenant limits enforced by processors
    spans_per_second: 50_000
    logs_mib_per_min: 100

  rate_limit:         # HTTP & gRPC endpoints
    enabled: true
    rps: 200
    burst: 400

  runtime:            # Ring 3 – sandboxing
    seccomp_profile: runtime/default
    apparmor_profile: strict
    read_only_rootfs: true
    mem_guard_mib: 256          # hard RSS fence per plugin proc
    cpu_guard_pct: 80           # soft cap via cgroup v2 weight
```

Anything left unspecified inherits **plugin-hard-coded** safe defaults.

---

## 2 · Quick-Start Profiles (copy-paste blocks)

### 2.1 Development (HTTP, JWT, no sandbox)

```yaml
security:
  tls.enabled: false
  auth:
    type: jwt
    jwt.secret: dev-secret
  rbac.enabled: false
  runtime:
    seccomp_profile: unconfined
    read_only_rootfs: false
```

### 2.2 Production hardened (mTLS + OAuth2 + full sandbox)

```yaml
security:
  tls:
    enabled: true
    verify_client: true
    min_version: TLS1.3
  auth:
    type: oauth2
    oauth2:
      issuer: https://idp.corp.example
      client_id: srediag-agent
      client_secret: "${OIDC_CLIENT_SECRET}"
  rbac:
    enabled: true
    default_role: viewer
  runtime:
    read_only_rootfs: true
    seccomp_profile: runtime/default
    mem_guard_mib: 512
```

---

## 3 · Plugin-Specific Hardening

Each plugin inherits system defaults **unless** overridden under
`plugins.d/<name>.yaml`:

```yaml
# /etc/srediag/plugins.d/clickhouseexporter.yaml
security:
  rbac:
    roles:
      exporter: ["write:telemetry"]
  runtime:
    mem_guard_mib: 256
```

*You cannot grant a plugin more privilege than the parent agent.*
Attempting to lift `read_only_rootfs: false` when the agent runs with
RO FS is ignored.

---

## 4 · CLI Runtime vs Service Daemon

| Setting | Service obeys | CLI obeys | Notes |
| :------ | :------------ | :-------- | :---- |
| `tls.*` | ✓ | ✗ | CLI uses localhost socket |
| `auth.*` | ✓ | ✓ (for remote URLs) | |
| `rbac.*` | ✓ | ✓ | CLI commands map to RBAC verbs |
| `runtime.*` | ✓ | ✓ (per-plugin) | |

Pass `--user` to force *user-mode* and read the *user* YAML overlay:

```bash
srediag diagnose system snapshot --user
```

---

## 5 · Validation & Monitoring

* **Startup guard** – Agent refuses to boot if
  * TLS enabled but certificate missing
  * RBAC enabled and `default_role` undefined
* **Live verification**

```bash
srediag plugin verify otlpreceiver
srediag service status --format yaml | yq .security
```

* **Metrics exposed**

| Metric | Meaning |
| :----- | :------ |
| `srediag_auth_failures_total` | AuthN failures (label: `reason`) |
| `srediag_rbac_denies_total`   | Access denials (label: `role`) |
| `srediag_plugin_sandbox_violations_total` | seccomp hits |

Alert threshold: **> 0** for sandbox violations.

---

## 6 · Incident-Response Playbook (condensed)

1. **Contain** – `srediag service stop` (systemd) or
   `srediag plugin disable …` for targeted issue.
2. **Gather** – `srediag service profile` + relevant logs.
3. **Quarantine plugin** across fleet:

   ```bash
   srediag plugin update --scope service --to quarantine vectorhashprocessor
   ```

4. **Patch** – Build fixed plugin, `install` + `enable` on canary.

---

## 7 · Best-Practice Checklist

| ✔︎ | Recommendation |
| :- | :------------- |
| ✓ | Rotate TLS certs every 90 d (use cert-manager or Vault side-car) |
| ✓ | Avoid wild-card (`*`) actions outside `admin` |
| ✓ | Enable **read-only root FS** and `seccomp` everywhere |
| ✓ | Pin plugin SHA-256 hashes in GitOps manifests |
| ✓ | Audit `srediag_plugin_sandbox_violations_total` weekly |
| ✓ | Keep `plugins.d/` under version control alongside pipelines |

---

## 8 · Troubleshooting Quick-Hits

| Symptom | Fix |
| :------ | :-- |
| `x509: unknown issuer` on startup | Point `tls.ca_file` to correct CA bundle |
| HTTP 403 despite valid token | Check `rbac.roles` mapping to `auth.claims.role` |
| Plugin crashes at boot | Lower `mem_guard_mib` or inspect seccomp log in `dmesg` |
| Collector reload fails (`plugin not enabled`) | Add plugin to `plugins.enabled` or correct alias name |

Increase verbosity:

```bash
SREDIAG_LOG_LEVEL=debug srediag service reload
```

---

## 9 · Reference & Further Reading

* Threat model & crypto flow – `docs/architecture/security.md`
* Supply-chain (cosign, SBOM, provenance) – `docs/security/supply-chain.md`
* Plugin sandbox internals – `docs/architecture/plugin-manager.md`
* Service & Collector YAML – `configuration/service.md`
* Plugin lifecycle commands – `cli/plugin.md`

---

## Parameter Reference

| YAML Key                | Env Var                    | CLI Flag                |
|-------------------------|----------------------------|-------------------------|
| `security.tls.enabled`  | `SREDIAG_TLS_ENABLED`      | `--tls-enabled`         |
| `security.tls.cert_file`| `SREDIAG_TLS_CERT_FILE`    | `--tls-cert`            |
| `security.tls.key_file` | `SREDIAG_TLS_KEY_FILE`     | `--tls-key`             |
| `security.tls.ca_file`  | `SREDIAG_TLS_CA_FILE`      | `--tls-ca`              |
| `security.auth.type`    | `SREDIAG_AUTH_TYPE`        | `--auth-type`           |
| `security.rbac.enabled` | `SREDIAG_RBAC_ENABLED`     | `--rbac-enabled`        |
| `srediag.config`        | `SREDIAG_CONFIG`           | `--config`              |

> **Warning:** Use `--config`/`SREDIAG_CONFIG` only for the main SREDIAG config. Security settings must use the above hierarchical keys/flags.

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
* Always check the effective config with `srediag security --print-config` (if available).
* Unknown YAML keys are logged at debug level and ignored for forward compatibility.
