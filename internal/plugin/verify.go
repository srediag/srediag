// TODO(architecture/plugin.md §2, §3, §5, §6, §9): Implement validation for the plugin trust chain, including checking for manifest presence, verifying SHA-256 digest, validating cosign signature, ensuring ABI compatibility, and checking RBAC capabilities.
// TODO(architecture/plugin.md §2, §3, §5, §6, §9): Refuse to load any plugin if the trust chain validation fails at any step.
// TODO(architecture/plugin.md §2, §3, §5, §6, §9): Log an error and refuse to load the plugin if ABI compatibility checks fail (Go 1.24, OTel API v1.30.0).
// TODO(architecture/plugin.md §2, §3, §5, §6, §9): Mark the plugin as invalid if the manifest file is missing or incomplete.
// TODO(architecture/plugin.md §2, §3, §5, §6, §9): Enforce and track plugin result states: active, disabled, or invalid.
package plugin
