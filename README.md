# SREDIAG — OBSERVO Diagnostics Agent

> Go-native diagnostics & observability agent • built on the OpenTelemetry Collector  
> Maintained by **Integra SH** • Contact: <marlon@integra.sh>

| Status | Go | Licence |
| :----: | :-: | :-----: |
| ![build](https://img.shields.io/badge/build-alpha-blue) | ![go](https://img.shields.io/badge/go-1.24.x-blue) | MIT |

---

## 1 — Project Scope

SREDIAG is the *edge data-plane* for the OBSERVO reliability platform.  
It starts as a **thin wrapper** around the upstream OpenTelemetry Collector and will grow—incrementally—into a fully-featured MSP-grade agent with hot-swappable plugins, volume-aware deduplication, CMDB drift reporting, and multi-tenant controls.

The repository is intentionally **early-stage**; most advanced features are *design-complete* but **not implemented yet**.  
Use it **only for experimentation** until the first beta tag (`2025.07.0`).

---

## 2 — Feature Matrix

| Area | Implemented (`main`) | In Development | Planned / Backlog |
| :-- | :-- | :-- | :-- |
| Collector core bootstrap | ✅ Static build via `otelcol-builder.yaml` | — | — |
| **Plugin framework** | ✅ MVP hot-loader (`internal/plugin/…`) | 🔴 IPC stabilisation · seccomp profiles | — |
| Diagnostics CLI | ✅ Skeleton commands in `cmd/` | 🟠 System / Perf / Security modules (TODOs) | — |
| Receivers / Processors / Exporters | 🚧 Upstream OTLP + Nop receivers | 🟠 Journald & System receivers · Batch/MemLimiter processors | ⏳ Vector-hash dedup · ClickHouse exporter |
| Dedup & Compression engine | — | — | ⏳ Phase 2 roadmap |
| CMDB drift & ITIL mapping | — | — | ⏳ Phase 3 roadmap |
| Multi-tenancy (quotas, SPIFFE) | — | — | ⚪ Phase 4 roadmap |
| K8s Operator & side-car injector | — | — | ⚪ Phase 5 roadmap |

Legend: **✅ shipped** • **🔴 active** • **🟠 queued** • **⏳ designed** • **⚪ backlog**

See `docs/specification.md` for the full technical roadmap.

---

## 3 — Repository Layout

```ascii

srediag/
├ cmd/                 # CLI entry points (Cobra)
├ internal/
│ ├ core/              # Version, config & logging helpers
│ ├ plugin/            # Hot-swap manager, shmipc stubs
│ └ diagnostic/        # TODO: system / perf / security commands
├ configs/             # Example collector configs + builder spec
├ docs/                # Product spec & appendices
├ magefiles/           # CI helper targets
└ Dockerfile           # Distroless build

```

---

## 4 — Quick Start

```bash
# 1. Clone & build (requires Go ≥ 1.24)
git clone https://github.com/srediag/srediag
cd srediag
make build            # wraps 'mage build'

# 2. Run a minimal collector
./bin/srediag --config configs/otel-config.yaml

# 3. Hot-load a sample plugin (once built)
curl -X POST --unix-socket /var/run/srediag.sock \
     -F file=@plugins/example.so http://plugin.load
```

### Container image

```bash
# Build a local multi-arch image
make docker           # docker buildx bakes distroless image
```

---

## 5 — Development Status & Roadmap

| Milestone | Target | ETA | Issue Board |
| :-- | :-- | :-- | :-- |
| **M0** — Baseline skeleton | Static OTel collector + plugin loader | **Done (2025-04)** | — |
| **M1** — IPC & Mem Guard | Stabilise shmipc-go, introduce `mem_guard.go` | **Jun 2025** | [link] |
| **M2** — Dedup & Compression | Vector-hash processor + tiered ZSTD/LZ4 | **Jul 2025** | [link] |
| **M3** — CMDB & Drift | Fingerprint lib, OTLP drift events | **Aug 2025** | [link] |
| **Beta 1** | Feature-freeze, Helm chart | **Sep 2025** | [link] |

Roadmap is tracked in **GitHub Projects › Milestones**.

---

## 6 — Building & Testing

```bash
# Lint, unit tests, SBOM & signed artifacts
make ci          # mage lint test sbom sign
```

* `golangci-lint` gates every PR  
* Unit coverage must stay **≥ 80 %**  
* Integration tests run in Kind (`tests/integration`)

---

## 7 — Contributing

1. **Fork** the repo & create a branch.  
2. Follow the **“Style & Commit”** guide in `CONTRIBUTING.md`.  
3. Run `make ci`; ensure all checks pass.  
4. Open a PR — A maintainer will review within 48 h.

We gladly accept **issue reports**, **feature proposals**, and **code PRs**—especially new plugins built with `srediag-plugin-sdk`.

---

## 8 — License

MIT.  See [`LICENSE`](LICENSE).

---

## 9 — Contact

*Mailing list*: <srediag-dev@integra.sh> •  
*Maintainer*: **Marlon Costa** (<marlon@integra.sh>)
