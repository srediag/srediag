package commands

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/srediag/srediag/internal/core"
)

// --- Mocks for internal/build functions ---

var (
	calledBuildAll       bool
	calledBuildPlugin    bool
	calledGenerate       bool
	calledInstallPlugins bool
	calledUpdateBuilder  bool
	lastCmd              *cobra.Command
	lastArgs             []string
	buildAllErr          error
	buildPluginErr       error
	generateErr          error
	installPluginsErr    error
	updateBuilderErr     error
)

func resetMocks() {
	calledBuildAll = false
	calledBuildPlugin = false
	calledGenerate = false
	calledInstallPlugins = false
	calledUpdateBuilder = false
	lastCmd = nil
	lastArgs = nil
	buildAllErr = nil
	buildPluginErr = nil
	generateErr = nil
	installPluginsErr = nil
	updateBuilderErr = nil
}

// Patch build package functions for testing
func defaultMockBuildFuncs() func() {
	origBuildAll := buildCLI_BuildAll
	origBuildPlugin := buildCLI_BuildPlugin
	origGenerate := buildCLI_Generate
	origInstall := buildCLI_InstallPlugins
	origUpdate := buildCLI_UpdateBuilder

	buildCLI_BuildAll = func(ctx *core.AppContext, cmd *cobra.Command, args []string) error {
		calledBuildAll = true
		lastCmd = cmd
		lastArgs = args
		return buildAllErr
	}
	buildCLI_BuildPlugin = func(ctx *core.AppContext, cmd *cobra.Command, args []string) error {
		calledBuildPlugin = true
		lastCmd = cmd
		lastArgs = args
		return buildPluginErr
	}
	buildCLI_Generate = func(ctx *core.AppContext, cmd *cobra.Command, args []string) error {
		calledGenerate = true
		lastCmd = cmd
		lastArgs = args
		return generateErr
	}
	buildCLI_InstallPlugins = func(ctx *core.AppContext, cmd *cobra.Command, args []string) error {
		calledInstallPlugins = true
		lastCmd = cmd
		lastArgs = args
		return installPluginsErr
	}
	buildCLI_UpdateBuilder = func(ctx *core.AppContext, cmd *cobra.Command, args []string) error {
		calledUpdateBuilder = true
		lastCmd = cmd
		lastArgs = args
		return updateBuilderErr
	}
	return func() {
		buildCLI_BuildAll = origBuildAll
		buildCLI_BuildPlugin = origBuildPlugin
		buildCLI_Generate = origGenerate
		buildCLI_InstallPlugins = origInstall
		buildCLI_UpdateBuilder = origUpdate
	}
}

// --- Test helpers ---

// Dummy AppContext for tests
func testAppContext() *core.AppContext {
	return &core.AppContext{}
}

// --- Tests ---

func TestNewBuildCmd_HasSubcommandsAndFlags(t *testing.T) {
	ctx := testAppContext()
	cmd, err := NewBuildCmd(ctx)
	if err != nil {
		t.Fatalf("NewBuildCmd returned error: %v", err)
	}

	// Check persistent flags
	if cmd.PersistentFlags().Lookup("build-config") == nil {
		t.Error("Expected persistent flag 'build-config'")
	}
	if cmd.PersistentFlags().Lookup("output-dir") == nil {
		t.Error("Expected persistent flag 'output-dir'")
	}

	// Check subcommands
	wantSubs := []string{"all", "plugin", "generate", "install", "update"}
	for _, sub := range wantSubs {
		found := false
		for _, c := range cmd.Commands() {
			if c.Name() == sub {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected subcommand %q", sub)
		}
	}
}

func TestNewBuildAllCmd_DelegatesToBuildAll(t *testing.T) {
	resetMocks()
	restore := defaultMockBuildFuncs()
	defer restore()

	ctx := testAppContext()
	cmd := newBuildAllCmd(ctx)
	if cmd.Use != "all" {
		t.Errorf("Expected Use 'all', got %q", cmd.Use)
	}
	err := cmd.RunE(cmd, []string{"foo"})
	if err != nil {
		t.Errorf("RunE returned error: %v", err)
	}
	if !calledBuildAll {
		t.Error("Expected build.CLI_BuildAll to be called")
	}
	if len(lastArgs) > 0 && lastArgs[0] != "foo" {
		t.Errorf("Expected args to be passed, got %v", lastArgs)
	}
}

func TestNewBuildAllCmd_DelegatesToBuildAll_Error(t *testing.T) {
	resetMocks()
	restore := defaultMockBuildFuncs()
	defer restore()
	buildAllErr = errors.New("fail all")

	ctx := testAppContext()
	cmd := newBuildAllCmd(ctx)
	err := cmd.RunE(cmd, []string{"foo"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "fail all")
	assert.True(t, calledBuildAll)
}

func TestNewBuildPluginCmd_DelegatesToBuildPlugin(t *testing.T) {
	resetMocks()
	restore := defaultMockBuildFuncs()
	defer restore()

	ctx := testAppContext()
	cmd := newBuildPluginCmd(ctx)
	if cmd.Use != "plugin" {
		t.Errorf("Expected Use 'plugin', got %q", cmd.Use)
	}
	// Check flags
	if cmd.Flags().Lookup("type") == nil {
		t.Error("Expected flag 'type'")
	}
	if cmd.Flags().Lookup("name") == nil {
		t.Error("Expected flag 'name'")
	}
	// Pass required flags for handler to be called
	errSet := cmd.Flags().Set("type", "exporter")
	require.NoError(t, errSet)
	errSet = cmd.Flags().Set("name", "foo")
	require.NoError(t, errSet)
	err := cmd.RunE(cmd, []string{})
	if err != nil {
		assert.Contains(t, err.Error(), "build plugin failed")
	}
	if !calledBuildPlugin {
		t.Error("Expected build.CLI_BuildPlugin to be called")
	}
}

func TestNewBuildGenerateCmd_DelegatesToGenerate(t *testing.T) {
	resetMocks()
	restore := defaultMockBuildFuncs()
	defer restore()

	ctx := testAppContext()
	cmd := newBuildGenerateCmd(ctx)
	if cmd.Use != "generate" {
		t.Errorf("Expected Use 'generate', got %q", cmd.Use)
	}
	// Check flags
	if cmd.Flags().Lookup("type") == nil {
		t.Error("Expected flag 'type'")
	}
	if cmd.Flags().Lookup("name") == nil {
		t.Error("Expected flag 'name'")
	}
	// Pass required flags for handler to be called
	errSet := cmd.Flags().Set("type", "processor")
	require.NoError(t, errSet)
	errSet = cmd.Flags().Set("name", "bar")
	require.NoError(t, errSet)
	err := cmd.RunE(cmd, []string{})
	if err != nil {
		assert.Contains(t, err.Error(), "generate failed")
	}
	if !calledGenerate {
		t.Error("Expected build.CLI_Generate to be called")
	}
}

func TestNewBuildInstallCmd_DelegatesToInstallPlugins(t *testing.T) {
	resetMocks()
	restore := defaultMockBuildFuncs()
	defer restore()

	ctx := testAppContext()
	cmd := newBuildInstallCmd(ctx)
	if cmd.Use != "install" {
		t.Errorf("Expected Use 'install', got %q", cmd.Use)
	}
	// Pass a dummy arg to ensure handler is called
	err := cmd.RunE(cmd, []string{"dummy"})
	if err != nil {
		assert.Contains(t, err.Error(), "install plugins failed")
	}
	if !calledInstallPlugins {
		t.Error("Expected build.CLI_InstallPlugins to be called")
	}
}

func TestNewBuildUpdateCmd_DelegatesToUpdateBuilder(t *testing.T) {
	resetMocks()
	restore := defaultMockBuildFuncs()
	defer restore()

	ctx := testAppContext()
	cmd := newBuildUpdateCmd(ctx)
	if cmd.Use != "update" {
		t.Errorf("Expected Use 'update', got %q", cmd.Use)
	}
	// Check flags
	if cmd.Flags().Lookup("yaml") == nil {
		t.Error("Expected flag 'yaml'")
	}
	if cmd.Flags().Lookup("gomod") == nil {
		t.Error("Expected flag 'gomod'")
	}
	if cmd.Flags().Lookup("plugin-gen") == nil {
		t.Error("Expected flag 'plugin-gen'")
	}
	// Pass required flags for handler to be called
	errSet := cmd.Flags().Set("yaml", "test.yaml")
	require.NoError(t, errSet)
	errSet = cmd.Flags().Set("gomod", "go.mod")
	require.NoError(t, errSet)
	errSet = cmd.Flags().Set("plugin-gen", "plugin/generated")
	require.NoError(t, errSet)
	err := cmd.RunE(cmd, []string{"sync"})
	if err != nil {
		assert.Contains(t, err.Error(), "update builder failed")
	}
	if !calledUpdateBuilder {
		t.Error("Expected build.CLI_UpdateBuilder to be called")
	}
}

// --- Patch points for build package (to allow test compilation) ---

// These variables allow us to patch the build package functions in tests.
// In production, they are assigned to the real build.CLI_* functions.
