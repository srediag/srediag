package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/srediag/srediag/internal/core"
)

// mockServiceAppContext is a minimal stub for core.AppContext.
type mockServiceAppContext struct {
	core.AppContext
}

func TestNewServiceCmd_HasAllSubcommands(t *testing.T) {
	ctx := &mockServiceAppContext{}
	cmd := NewServiceCmd(&ctx.AppContext)

	expectedSubs := []string{
		"start", "stop", "restart", "reload", "detach", "status", "health",
		"profile", "tail-logs", "validate", "install-unit", "uninstall-unit", "gc",
	}
	found := map[string]bool{}
	for _, c := range cmd.Commands() {
		found[c.Name()] = true
	}
	for _, name := range expectedSubs {
		assert.Truef(t, found[name], "subcommand %q should be present", name)
	}
}

func TestServiceSubcommands_RunE_NoError(t *testing.T) {
	ctx := &mockServiceAppContext{}
	root := NewServiceCmd(&ctx.AppContext)
	subcommands := root.Commands()

	for _, sub := range subcommands {
		t.Run(sub.Name(), func(t *testing.T) {
			sub.SetArgs([]string{})
			err := sub.Execute()
			assert.NoError(t, err, "subcommand %q should not error", sub.Name())
		})
	}
}
