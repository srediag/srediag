# SREDIAG — OBSERVO Diagnostics Agent

> Go-native diagnostics & observability agent • built on the OpenTelemetry Collector  
> Maintained by **Integra SH**   •   Contact: <marlon@integra.sh>

<div align="center">

| Build | Go | Docs | Licence |
| :---: | :-: | :--: | :-----: |
| ![status](https://img.shields.io/badge/status-alpha-blue) | ![go](https://img.shields.io/badge/go-1.24.x-blue) | [docs ↗](./docs/README.md) | MIT |

</div>

---

SREDIAG is the **edge data-plane** of the OBSERVO reliability platform.  
It wraps the upstream **OpenTelemetry Collector** and adds, step-by-step, the capabilities MSPs need at scale: hot-swappable plugins, content-aware deduplication, ITIL-aligned CMDB drift, and strict tenant isolation.

This repository is **work-in-progress**. The table below names what is **already wired**, what is **actively being coded**, and what is **road-mapped** for later milestones.

For the full product spec see **`docs/specification.md`** (generated from the canonical draft in `docs/architecture/`).

---

## 1 — Feature Matrix

| Area | Implemented (`main`) | In Development | Planned |
| :-- | :-- | :-- | :-- |
| Collector bootstrap | ✅ Static build (`build/otelcol-builder.yaml`) | — | — |
| Plugin framework | ✅ MVP hot-loader (`internal/plugin/…`) | 🔴 IPC hardening · seccomp | — |
| Diagnostics CLI | ✅ Skeleton commands (`cmd/`) | 🟠 System/Perf/Security modules (see `TODO` tags) | — |
| Receivers / Processors | 🚧 OTLP & Nop receivers | 🟠 Journald & System receivers · Batch/MemLimiter processors | ⏳ Vector-hash dedup processor |
| Exporters | 🚧 OTLP exporter | 🟠 ClickHouse exporter | ⏳ Cloud-native sinks (S3, GCS) |
| CMDB drift | — | — | ⏳ Phase 3 |
| Multi-tenancy (SPIFFE, quotas) | — | — | ⚪ Phase 4 |

Legend — **✅ shipped** • **🔴 active** • **🟠 queued** • **⏳ designed** • **⚪ backlog**

---

## 2 — Documentation Map

| Topic | Path |
| :-- | :-- |
| **Architecture overview** | `docs/architecture/README.md` |
| OpenTelemetry integration | `docs/architecture/opentelemetry.md` |
| Security design | `docs/architecture/security.md` |
| **Getting started** (install, config, first run) | `docs/getting-started/` |
| Configuration reference | `docs/configuration/` |
| Command-line help | `docs/cli/` |
| Full spec & roadmap | `docs/specification.md` |

All doc pages build locally with `mkdocs serve`.

---

## 3 — Quick Start (developer build)

```bash
# 1 Clone & build (needs Go ≥ 1.24)
git clone https://github.com/srediag/srediag
cd srediag
make build            # wraps 'mage build'

# 2 Run the collector with a demo config
./bin/srediag --config configs/otel-config.yaml

# 3 Hot-load a sample plugin (once compiled)
curl -X POST --unix-socket /var/run/srediag.sock \
     -F file=@plugins/example.so http://plugin.load
```

### Container

```bash
make docker   # builds multi-arch distroless image
```

For advanced setups (Helm, K8s Operator) see **`docs/getting-started/installation.md`**.

---

## 4 — Roadmap

| Milestone | Deliverable | ETA |
| :-- | :-- | :-- |
| **M1** | IPC + `mem_guard.go` | **Jun 2025** |
| **M2** | Dedup processor + tiered ZSTD/LZ4 | Jul 2025 |
| **M3** | CMDB drift events | Aug 2025 |
| **Beta 1** | Helm chart, first public image | Sep 2025 |

Track progress in **GitHub Projects › “SREDIAG Roadmap”**.

---

## 5 — Building & Testing

```bash
make ci   # lint → unit tests → SBOM → cosign
```

* `golangci-lint` gates every PR  
* Unit coverage ≥ 80 %  
* Integration tests run in Kind (`tests/integration`)

---

## 6 — Contributing

1. Fork → branch → `make ci`.  
2. Follow the commit style in `CONTRIBUTING.md`.  
3. Open a PR — we review within 48 h.

Plugins built with [`srediag-plugin-sdk`](docs/cli/README.md#plugin-sdk) are very welcome.

---

## 7 — License

MIT — see [`LICENSE`](LICENSE).

---

## 8 — Contact

*Mailing list*: <srediag-dev@integra.sh>  
*Maintainer*: **Marlon Costa** (<marlon@integra.sh>)
