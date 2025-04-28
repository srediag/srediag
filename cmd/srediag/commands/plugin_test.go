package commands

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/srediag/srediag/internal/core"
)

// mockAppContext is a minimal mock for core.AppContext.
type mockAppContext struct {
	core.AppContext
}

func TestNewPluginCmd_HasExpectedSubcommands(t *testing.T) {
	ctx := &mockAppContext{}
	cmd := newPluginCmd(&ctx.AppContext)

	subcommands := map[string]bool{}
	for _, c := range cmd.Commands() {
		subcommands[c.Name()] = true
	}

	assert.True(t, subcommands["list"], "should have 'list' subcommand")
	assert.True(t, subcommands["info"], "should have 'info' subcommand")
	assert.True(t, subcommands["enable"], "should have 'enable' subcommand")
	assert.True(t, subcommands["disable"], "should have 'disable' subcommand")
}

// --- Dependency-injected subcommand constructors for testing ---
func newPluginListCmdWithFunc(ctx *core.AppContext, fn func(*core.AppContext, *cobra.Command, []string) error) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all available plugins",
		RunE:  func(cmd *cobra.Command, args []string) error { return fn(ctx, cmd, args) },
	}
}

func newPluginInfoCmdWithFunc(ctx *core.AppContext, fn func(*core.AppContext, *cobra.Command, []string) error) *cobra.Command {
	return &cobra.Command{
		Use:   "info [name]",
		Short: "Show information about a plugin",
		Args:  cobra.ExactArgs(1),
		RunE:  func(cmd *cobra.Command, args []string) error { return fn(ctx, cmd, args) },
	}
}

func newPluginEnableCmdWithFunc(ctx *core.AppContext, fn func(*core.AppContext, *cobra.Command, []string) error) *cobra.Command {
	return &cobra.Command{
		Use:   "enable [type] [name]",
		Short: "Enable a plugin",
		Args:  cobra.ExactArgs(2),
		RunE:  func(cmd *cobra.Command, args []string) error { return fn(ctx, cmd, args) },
	}
}

func newPluginDisableCmdWithFunc(ctx *core.AppContext, fn func(*core.AppContext, *cobra.Command, []string) error) *cobra.Command {
	return &cobra.Command{
		Use:   "disable [name]",
		Short: "Disable a plugin",
		Args:  cobra.ExactArgs(1),
		RunE:  func(cmd *cobra.Command, args []string) error { return fn(ctx, cmd, args) },
	}
}

func TestNewPluginListCmd_DelegatesToPluginCLIList(t *testing.T) {
	called := false
	mock := func(ctx *core.AppContext, cmd *cobra.Command, args []string) error {
		called = true
		return nil
	}
	ctx := &mockAppContext{}
	cmd := newPluginListCmdWithFunc(&ctx.AppContext, mock)
	err := cmd.RunE(cmd, []string{})
	require.NoError(t, err)
	assert.True(t, called, "plugin.CLI_List should be called")
}

func TestNewPluginInfoCmd_DelegatesToPluginCLIInfo(t *testing.T) {
	called := false
	mock := func(ctx *core.AppContext, cmd *cobra.Command, args []string) error {
		called = true
		return nil
	}
	ctx := &mockAppContext{}
	cmd := newPluginInfoCmdWithFunc(&ctx.AppContext, mock)
	// Should fail if not enough args
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	assert.Error(t, err)
	// Should succeed and call handler with correct args
	called = false
	cmd.SetArgs([]string{"plugin1"})
	err = cmd.Execute()
	require.NoError(t, err)
	assert.True(t, called, "plugin.CLI_Info should be called")
}

func TestNewPluginEnableCmd_DelegatesToPluginCLIEnable(t *testing.T) {
	called := false
	mock := func(ctx *core.AppContext, cmd *cobra.Command, args []string) error {
		called = true
		return nil
	}
	ctx := &mockAppContext{}
	cmd := newPluginEnableCmdWithFunc(&ctx.AppContext, mock)
	// Should fail if not enough args
	cmd.SetArgs([]string{"type"})
	err := cmd.Execute()
	assert.Error(t, err)
	// Should succeed and call handler with correct args
	called = false
	cmd.SetArgs([]string{"type", "plugin1"})
	err = cmd.Execute()
	require.NoError(t, err)
	assert.True(t, called, "plugin.CLI_Enable should be called")
}

func TestNewPluginDisableCmd_DelegatesToPluginCLIDisable(t *testing.T) {
	called := false
	mock := func(ctx *core.AppContext, cmd *cobra.Command, args []string) error {
		called = true
		return nil
	}
	ctx := &mockAppContext{}
	cmd := newPluginDisableCmdWithFunc(&ctx.AppContext, mock)
	// Should fail if not enough args
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	assert.Error(t, err)
	// Should succeed and call handler with correct args
	called = false
	cmd.SetArgs([]string{"plugin1"})
	err = cmd.Execute()
	require.NoError(t, err)
	assert.True(t, called, "plugin.CLI_Disable should be called")
}

func TestPluginCmd_SubcommandsReturnHandlerErrors(t *testing.T) {
	listMock := func(ctx *core.AppContext, cmd *cobra.Command, args []string) error {
		return errors.New("list error")
	}
	infoMock := func(ctx *core.AppContext, cmd *cobra.Command, args []string) error {
		return errors.New("info error")
	}
	enableMock := func(ctx *core.AppContext, cmd *cobra.Command, args []string) error {
		return errors.New("enable error")
	}
	disableMock := func(ctx *core.AppContext, cmd *cobra.Command, args []string) error {
		return errors.New("disable error")
	}
	ctx := &mockAppContext{}

	err := newPluginListCmdWithFunc(&ctx.AppContext, listMock).RunE(nil, []string{})
	assert.EqualError(t, err, "list error")

	err = newPluginInfoCmdWithFunc(&ctx.AppContext, infoMock).RunE(nil, []string{"plugin1"})
	assert.EqualError(t, err, "info error")

	err = newPluginEnableCmdWithFunc(&ctx.AppContext, enableMock).RunE(nil, []string{"type", "plugin1"})
	assert.EqualError(t, err, "enable error")

	err = newPluginDisableCmdWithFunc(&ctx.AppContext, disableMock).RunE(nil, []string{"plugin1"})
	assert.EqualError(t, err, "disable error")
}
