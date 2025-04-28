// TODO: Implement SPIFFE-mTLS 1.3, cert rotation, and least-priv gRPC credentials (see architecture/security.md §0, §5)
// TODO: Enforce TLS 1.3 / AES-256-GCM-SHA384 for all data-plane connections (see architecture/security.md §5)
// TODO: Ensure all tls.* keys reside on tmpfs or are injected via KMS (see architecture/security.md §5)
package core
