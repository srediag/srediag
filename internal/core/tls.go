// TODO: Implement SPIFFE-mTLS 1.3, certificate rotation, and least-privilege gRPC credentials for secure communication (see docs/architecture/security.md §0, §5)
// TODO: Enforce TLS 1.3 with AES-256-GCM-SHA384 for all data-plane connections (see docs/architecture/security.md §5)
// TODO: Ensure all TLS keys are stored on tmpfs or injected via KMS for security (see docs/architecture/security.md §5)
package core
