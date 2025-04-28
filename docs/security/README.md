# SREDIAG Supply-Chain Security Process

Measures that protect SREDIAG from build-time or distribution-time tampering.

---

## 1 · Signed Artifacts

| Artifact | Tool | Verification |
| :-- | :-- | :-- |
| OCI images | `cosign sign` | Checked at agent start & CI |
| Go plugins (`*.so`) | `cosign sign-blob` | Verified by Plugin Manager |

---

## 2 · Build Provenance (SLSA Level 2)

* GitHub Actions → `slsa-github-generator`.  
* Attestation includes commit SHA, builder image, dependency list.

---

## 3 · SBOM

* Generated via **Syft** (`make sbom`) in CycloneDX JSON.  
* Published alongside every OCI image.

---

## 4 · Dependency Management

* Weekly Dependabot PRs.  
* Daily `govulncheck` in CI.  
* Critical CVEs ⇒ Control-Plane sets **quarantine** flag, blocking plugins fleet-wide.

---

## 5 · Release Checklist

1. CI green → artefacts built.  
2. Generate SBOM, SLSA provenance.  
3. cosign sign images & plugins.  
4. Upload artefacts + attestation to GH Releases / OCI registry.  
5. Update `docs/changelog.md`.

---

## 6 · See Also

* [Security Architecture](../architecture/security.md)  
* [CI/CD Pipeline Guide](../development/testing.md)
