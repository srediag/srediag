// Package plugin provides plugin sandboxing and security enforcement for SREDIAG plugins.
//
// This file defines placeholder structures for advanced security features to be implemented
// in future phases including seccomp, AppArmor, cgroup, and filesystem isolation.
//
// Note: Advanced security features to be implemented in Phase 2:
// - Seccomp, AppArmor, cgroup, memlock, and netcls sandboxing enforcement
// - Sandbox policy configuration via YAML with OPA-based policy validation
// - Read-only root filesystem and non-root user enforcement
//
// Usage:
//   - These structures serve as placeholders for future security implementation
//   - Current plugin system provides basic process isolation
//   - Advanced sandboxing will be implemented based on architecture/plugin.md requirements
package plugin

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"syscall"
)

// SandboxConfig represents the configuration for plugin sandboxing
type SandboxConfig struct {
	// EnableSeccomp enables seccomp filtering for plugin processes
	EnableSeccomp bool `yaml:"enable_seccomp"`
	// EnableAppArmor enables AppArmor profiles for plugin processes
	EnableAppArmor bool `yaml:"enable_apparmor"`
	// EnableCgroups enables cgroup resource limits for plugin processes
	EnableCgroups bool `yaml:"enable_cgroups"`
	// ReadOnlyRootFS enforces read-only root filesystem
	ReadOnlyRootFS bool `yaml:"readonly_rootfs"`
	// NonRootUser enforces non-root user execution
	NonRootUser bool `yaml:"non_root_user"`
	// MemoryLimit sets memory limit in bytes
	MemoryLimit int64 `yaml:"memory_limit"`
	// CPULimit sets CPU limit as percentage
	CPULimit int `yaml:"cpu_limit"`
}

// Sandbox manages plugin process sandboxing and security enforcement
type Sandbox struct {
	config *SandboxConfig
}

// NewSandbox creates a new sandbox manager with the given configuration
func NewSandbox(config *SandboxConfig) *Sandbox {
	if config == nil {
		config = &SandboxConfig{
			EnableSeccomp:  false, // Disabled until Phase 2
			EnableAppArmor: false, // Disabled until Phase 2
			EnableCgroups:  false, // Disabled until Phase 2
			ReadOnlyRootFS: false, // Disabled until Phase 2
			NonRootUser:    true,  // Basic non-root enforcement enabled
			MemoryLimit:    0,     // No limit until Phase 2
			CPULimit:       0,     // No limit until Phase 2
		}
	}
	return &Sandbox{config: config}
}

// ApplySandbox applies security restrictions to a plugin process
func (s *Sandbox) ApplySandbox(ctx context.Context, cmd *os.Process) error {
	// Basic security: ensure non-root user
	if s.config.NonRootUser {
		if err := s.enforceNonRootUser(); err != nil {
			return fmt.Errorf("failed to enforce non-root user: %w", err)
		}
	}

	// Future Phase 2 implementations will add:
	// - s.applySeccompFilter()
	// - s.applyAppArmorProfile()
	// - s.applyCgroupLimits()
	// - s.applyReadOnlyRootFS()

	return nil
}

// enforceNonRootUser checks that the current process is not running as root
func (s *Sandbox) enforceNonRootUser() error {
	currentUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}

	if currentUser.Uid == "0" {
		return fmt.Errorf("plugin cannot run as root user for security reasons")
	}

	return nil
}

// ValidateSecurityPolicy validates the sandbox configuration
func (s *Sandbox) ValidateSecurityPolicy() error {
	if s.config == nil {
		return fmt.Errorf("sandbox configuration is nil")
	}

	// Basic validation for enabled features
	if s.config.NonRootUser {
		if err := s.enforceNonRootUser(); err != nil {
			return fmt.Errorf("non-root user validation failed: %w", err)
		}
	}

	if s.config.MemoryLimit < 0 {
		return fmt.Errorf("memory limit cannot be negative")
	}

	if s.config.CPULimit < 0 || s.config.CPULimit > 100 {
		return fmt.Errorf("CPU limit must be between 0 and 100")
	}

	return nil
}

// GetSecurityStatus returns the current security enforcement status
func (s *Sandbox) GetSecurityStatus() map[string]bool {
	return map[string]bool{
		"seccomp_enabled":  s.config.EnableSeccomp,
		"apparmor_enabled": s.config.EnableAppArmor,
		"cgroups_enabled":  s.config.EnableCgroups,
		"readonly_rootfs":  s.config.ReadOnlyRootFS,
		"non_root_user":    s.config.NonRootUser,
		"memory_limited":   s.config.MemoryLimit > 0,
		"cpu_limited":      s.config.CPULimit > 0,
	}
}

// SetProcessCredentials sets the process credentials for plugin execution
func (s *Sandbox) SetProcessCredentials(uid, gid int) error {
	if uid == 0 && s.config.NonRootUser {
		return fmt.Errorf("cannot set root UID when non-root user is enforced")
	}

	// Set process credentials
	if err := syscall.Setgid(gid); err != nil {
		return fmt.Errorf("failed to set GID: %w", err)
	}

	if err := syscall.Setuid(uid); err != nil {
		return fmt.Errorf("failed to set UID: %w", err)
	}

	return nil
}
