# SREDIAG TODO List

## SREDIAG — Master TODO & Execution Plan (sync 2025‑04‑27)

> **Scope** – Concrete engineering work. Strategy, risk and glossary live in [`docs/specification.md`](specification.md). Update with **minimal diffs**, append your name + date in commit messages.

---

## 0 · Snapshot Summary

* **Current tag**: `v0.0.1` (baseline collector + plugin loader)  — built from commit `HEAD@main ‹replace‑sha›`  
* **Upcoming tag**: `v0.1.0` (Phase 1 Perf Core, due 2025‑06‑14)

---

## 1 · Requirements Coverage (Spec §3)

| ID | Requirement | Phase | Status | Work‑Item |
| :-- | :-- | :-: | :-: | :-- |
| **F‑1** | OTLP ingestion (gRPC + HTTP) | 0 | ✅ | core bootstrap |
| **F‑2** | Hot‑reload plugins | 1 | 🔴 | IPC / Plugin Manager |
| **F‑3** | 30 s log dedup | 2 | 🟠 | `vectorhashprocessor` |
| **F‑4** | `%CMDB_HASH%` attribute | 3 | 🟢 | fingerprint library |
| **F‑5** | Tenant auth & rate limits | 4 | ⚪ | quotas middleware |
| **F‑6** | `/healthz` & `/metrics` | 0 | ✅ | core exporter |
| **F‑7** | Signed remote‑config | 3 | 🟢 | remote‑config engine |
| **F‑8** | Log auto‑format detect | 2 | 🟠 | dedup heuristics |
| **F‑9** | 10 GiB offline buffer | 2 | 🟠 | RocksDB spill cache |
| **F‑10** | Dynamic tail‑sampling | 4 | ⚪ | tail‑sampler rules |

_Gaps_: **F‑8** & **F‑9** require prototype → schedule sprint **2025‑07‑01**.

---

## 2 · Phase Roadmap (aligns Spec §6)

| Phase | Deliverables | Tag | ETA | Lead | Status |
| :-- | :-- | :-- | :-- | :-- | :-- |
| 0 Baseline | Static collector · plugin loader · basic CI | `v0.0.1` | 2025‑04‑15 | Core | ✅ |
| 1 Perf Core | shmipc IPC · `mem_guard` · hash‑index cache | `v0.1.0` | 2025‑06‑14 | Core | 🔴 |
| 2 Dedup + Compression | vectorhashprocessor · adaptive LZ4/ZSTD | `v0.2.0` | 2025‑07‑19 | Perf | 🟠 |
| 3 Governance | CMDB drift → OTLP · audit channel · RocksDB tier | `v0.3.0` | 2025‑08‑23 | Platform | 🟢 |
| 4 Multi‑Tenancy | SPIFFE mTLS · quotas · ClickHouse exporter | `v0.4.0` | 2025‑09‑27 | Security | ⚪ |
| 5 Auto‑Pilot | Helm chart · K8s Operator · side‑car inject | `v0.5.0‑beta1` | 2025‑10‑25 | DevOps | ⚪ |
| GA | FedRAMP draft · FIPS build · docs freeze | `v1.0.0` | 2025‑11‑29 | PMO | ⚪ |

Legend — ✅ done • 🔴 active • 🟠 queued • 🟢 design • ⚪ backlog

---

## 3 · Detailed Work‑Streams

### 3.1 Core Collector & Runtime

| ID | Task | Owner | ETA | Notes |
| :-- | :-- | :-- | :-- | :-- |
| C‑01 | Upgrade to OTel v0.124.0 / API v1.30.0 | Core | 2025‑05‑10 | update `go.mod`, `go mod tidy` |
| C‑02 | Pipeline builder (Go → YAML) | Core | 2025‑05‑24 | template pkg |
| C‑03 | Component registry with lazy load | Core | 2025‑06‑07 | `internal/core/registry.go` |
| C‑04 | Graceful shutdown (flush + RocksDB close) | Core | 2025‑06‑14 | tie into signals |

### 3.2 Plugin Framework

| P‑ID | Task | Phase | ETA |
| :-- | :-- | :-: | :-- |
| P‑01 | IPC fuzz/integration tests | 1 | 2025‑05‑28 |
| P‑02 | Manifest v1 JSON schema | 1 | 2025‑05‑31 |
| P‑03 | seccomp profile generator | 1 | 2025‑06‑10 |
| P‑04 | Heartbeat RPC + Prom metric | 1 | 2025‑06‑12 |

### 3.3 Dedup & Compression (Phase 2)

* RocksDB wrapper (WAL sync, CF per tenant)
* Content‑defined chunker (>4 KiB messages)
* Adaptive encoder benchmark (zstd‑1 vs lz4‑fast)
* Integration test: replay 10 M logs → ≥70 % egress reduction

### 3.4 Diagnostics Plugins

| D‑ID | Plugin / Feature | Phase | Status |
| :-- | :-- | :-: | :-- |
| D‑01 | **System diagnostics** (CPU, mem, IO, net) | 3 | 🟢 design |
| D‑02 | **Kubernetes diagnostics** (cluster, node, pod) | 3 | 🟢 design |
| D‑03 | **Cloud provider stubs** (AWS, Azure, GCP) | 5 | ⚪ |
| D‑04 | **IaC analyzers** (Terraform, K8s manifests, Helm) | 5 | ⚪ |

### 3.5 Observability & Dashboards

* Self‑metrics exporter, zPages endpoint, Grafana JSON dashboards (Phase 2‑3)

### 3.6 Docs & DX

* MkDocs build in CI, plugin SDK tutorial, contributing & style guides.

### 3.7 Testing & Quality

* Unit ≥ 85 % lines, integration in Kind, perf benchstat gate, chaos scripts.

### 3.8 Security & Compliance

* SBOM, SLSA L2, cosign signatures (Phase 2) • FIPS build (Phase 4) • FedRAMP draft (Phase 5).

---

## 4 · Deployment & Ops Rollout (Spec §9)

| Milestone | Env | Owner | Status | Artefact |
| :-- | :-- | :-- | :-- | :-- |
| M0 (05‑2025) PoC | Kind | DevOps | ✅ | `reports/m0‑poc.md` |
| M1 (06‑2025) Plugin GA | EKS / OKE | SRE | 🔴 | `helm/values‑msp.yaml` |
| M2 (07‑2025) Dedup lab | Perf lab | Perf | 🟠 | `bench/m2‑plan.md` |
| M3 (08‑2025) CMDB pilot | MSP pilot | Cust Success | 🟢 | `pilots/cmdb‑run.md` |
| M4 (09‑2025) Quotas & Billing | SaaS stage | FinOps | ⚪ | `docs/billing‑hooks.md` |
| M5 (10‑2025) FIPS build | Gov cloud | Security | ⚪ | `compliance/fedramp‑draft.md` |
| GA (11‑2025) SaaS launch | Prod | Product | ⚪ | `release‑checklist.md` |

Artefacts must be updated (link + status) at milestone closure.

---

## 5 · Compliance & Supply‑Chain Checklist

| Control | Tool / Evidence | Phase |
| :-- | :-- | :-: |
| SBOM (CycloneDX) | `make sbom` via Syft | 2 |
| SLSA L2 | `slsa-github-generator` | 2 |
| Cosign signatures | Release pipeline | 2 |
| FIPS‑140‑3 build | `CGO_ENABLED=1`, BoringCrypto toolchain | 4 |
| FedRAMP SSP | docs/compliance/fedramp | 5 |
