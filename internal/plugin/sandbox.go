// TODO(architecture/plugin.md §3, §4, §5, §9): Enforce seccomp, AppArmor, cgroup, memlock, and netcls sandboxing for all plugin processes to ensure runtime isolation and resource control.
// TODO(architecture/plugin.md §3, §4, §5, §9): Implement support for sandbox policy configuration via YAML and add OPA-based policy tests for plugin isolation.
// TODO(architecture/plugin.md §3, §4, §5, §9): Enforce read-only root filesystem and require all plugin processes to run as a non-root user (UID).
package plugin
