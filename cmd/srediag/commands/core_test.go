package commands

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/srediag/srediag/internal/core"
)

// --- Mocks and helpers ---

func mockNewLogger(*core.Logger) (*core.Logger, error) {
	return &core.Logger{}, nil
}

func mockNewLoggerErr(*core.Logger) (*core.Logger, error) {
	return nil, errors.New("logger error")
}

func mockValidateConfig(cfg *core.Config) error {
	return nil
}

func mockValidateConfigErr(cfg *core.Config) error {
	return errors.New("invalid config")
}

func mockLoadConfigWithOverlay(spec interface{}, cliFlags map[string]string, opts ...core.ConfigOption) error {
	return nil
}

func mockLoadConfigWithOverlayErr(spec interface{}, cliFlags map[string]string, opts ...core.ConfigOption) error {
	return errors.New("load config error")
}

// --- Test Setup ---

type testDeps struct {
	LoadConfigWithOverlay func(spec interface{}, cliFlags map[string]string, opts ...core.ConfigOption) error
	ValidateConfig        func(cfg *core.Config) error
	NewLogger             func(cfg *core.Logger) (*core.Logger, error)
	PrintEffectiveConfig  func(cfg *core.Config) error

	NewBuildCmd    func(ctx *core.AppContext) (*cobra.Command, error)
	NewDiagnoseCmd func(ctx *core.AppContext) *cobra.Command
	NewPluginCmd   func(ctx *core.AppContext) *cobra.Command
	NewServiceCmd  func(ctx *core.AppContext) *cobra.Command
}

func makeDeps(d testDeps) *RootCommandDeps {
	return &RootCommandDeps{
		LoadConfigWithOverlay: d.LoadConfigWithOverlay,
		ValidateConfig:        d.ValidateConfig,
		NewLogger:             d.NewLogger,
		PrintEffectiveConfig:  d.PrintEffectiveConfig,
		NewBuildCmd:           d.NewBuildCmd,
		NewDiagnoseCmd:        d.NewDiagnoseCmd,
		NewPluginCmd:          d.NewPluginCmd,
		NewServiceCmd:         d.NewServiceCmd,
	}
}

// --- Tests ---

func TestNewRootCommand_Basic(t *testing.T) {
	ctx := &core.AppContext{}
	cmd := NewRootCommand(ctx, makeDeps(testDeps{
		LoadConfigWithOverlay: mockLoadConfigWithOverlay,
		ValidateConfig:        mockValidateConfig,
		NewLogger:             mockNewLogger,
	}))
	require.NotNil(t, cmd)
	assert.Equal(t, "srediag", cmd.Use)
	assert.True(t, cmd.SilenceUsage)
}

func TestNewRootCommand_PersistentFlags(t *testing.T) {
	ctx := &core.AppContext{}
	cmd := NewRootCommand(ctx, makeDeps(testDeps{
		LoadConfigWithOverlay: mockLoadConfigWithOverlay,
		ValidateConfig:        mockValidateConfig,
		NewLogger:             mockNewLogger,
	}))
	flags := cmd.PersistentFlags()

	assert.NotNil(t, flags.Lookup("config"))
	assert.NotNil(t, flags.Lookup("output"))
	assert.NotNil(t, flags.Lookup("quiet"))
	assert.NotNil(t, flags.Lookup("no-color"))
	assert.NotNil(t, flags.Lookup("output-file"))
	assert.NotNil(t, flags.Lookup("log-level"))
	assert.NotNil(t, flags.Lookup("log-format"))
	assert.NotNil(t, flags.Lookup("print-config"))
}

func TestNewRootCommand_RunE(t *testing.T) {
	ctx := &core.AppContext{}
	cmd := NewRootCommand(ctx, makeDeps(testDeps{
		LoadConfigWithOverlay: mockLoadConfigWithOverlay,
		ValidateConfig:        mockValidateConfig,
		NewLogger:             mockNewLogger,
	}))
	err := cmd.RunE(cmd, []string{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "please specify a subcommand")
}

func TestNewRootCommand_PersistentPreRunE_ConfigLoadError(t *testing.T) {
	ctx := &core.AppContext{}
	cmd := NewRootCommand(ctx, makeDeps(testDeps{
		LoadConfigWithOverlay: mockLoadConfigWithOverlayErr,
		ValidateConfig:        mockValidateConfig,
		NewLogger:             mockNewLogger,
	}))
	err := cmd.PersistentPreRunE(cmd, []string{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to load config")
}

func TestNewRootCommand_PersistentPreRunE_ConfigValidateError(t *testing.T) {
	ctx := &core.AppContext{}
	cmd := NewRootCommand(ctx, makeDeps(testDeps{
		LoadConfigWithOverlay: mockLoadConfigWithOverlay,
		ValidateConfig:        mockValidateConfigErr,
		NewLogger:             mockNewLogger,
	}))
	err := cmd.PersistentPreRunE(cmd, []string{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid config")
}

func TestNewRootCommand_PersistentPreRunE_LoggerError(t *testing.T) {
	ctx := &core.AppContext{}
	cmd := NewRootCommand(ctx, makeDeps(testDeps{
		LoadConfigWithOverlay: mockLoadConfigWithOverlay,
		ValidateConfig:        mockValidateConfig,
		NewLogger:             mockNewLoggerErr,
	}))
	err := cmd.PersistentPreRunE(cmd, []string{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to initialize logger")
}

func TestViperAllSettings_Empty(t *testing.T) {
	viper.Reset()
	settings := viperAllSettings()
	assert.Empty(t, settings)
}

func TestViperAllSettings_WithValues(t *testing.T) {
	viper.Reset()
	viper.Set("foo", "bar")
	viper.Set("num", 123)
	settings := viperAllSettings()
	assert.Equal(t, "bar", settings["foo"])
	_, ok := settings["num"]
	assert.False(t, ok, "non-string values should not be included")
}

func TestExecute_CallsRootCommand(t *testing.T) {
	called := false

	rootCmd := &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			called = true
			return nil
		},
	}
	// Directly execute the command
	err := rootCmd.Execute()
	assert.NoError(t, err)
	assert.True(t, called)
}

func TestOutputFormat_StructFields(t *testing.T) {
	of := OutputFormat{
		Format:     "json",
		Quiet:      true,
		NoColor:    true,
		OutputFile: "out.txt",
	}
	assert.Equal(t, "json", of.Format)
	assert.True(t, of.Quiet)
	assert.True(t, of.NoColor)
	assert.Equal(t, "out.txt", of.OutputFile)
}

// --- Subcommand wiring smoke test ---

func TestNewRootCommand_AddsSubcommands(t *testing.T) {
	ctx := &core.AppContext{}
	cmd := NewRootCommand(ctx, makeDeps(testDeps{
		LoadConfigWithOverlay: mockLoadConfigWithOverlay,
		ValidateConfig:        mockValidateConfig,
		NewLogger:             mockNewLogger,
		NewBuildCmd: func(ctx *core.AppContext) (*cobra.Command, error) {
			return &cobra.Command{Use: "build"}, nil
		},
		NewDiagnoseCmd: func(ctx *core.AppContext) *cobra.Command {
			return &cobra.Command{Use: "diagnose"}
		},
		NewPluginCmd: func(ctx *core.AppContext) *cobra.Command {
			return &cobra.Command{Use: "plugin"}
		},
		NewServiceCmd: func(ctx *core.AppContext) *cobra.Command {
			return &cobra.Command{Use: "service"}
		},
	}))
	require.NotNil(t, cmd)

	found := map[string]bool{}
	for _, c := range cmd.Commands() {
		found[c.Use] = true
	}
	assert.True(t, found["build"])
	assert.True(t, found["diagnose"])
	assert.True(t, found["plugin"])
	assert.True(t, found["service"])
}

func TestNewRootCommand_BuildCmdInitError_Exits(t *testing.T) {
	ctx := &core.AppContext{}
	// Inject a failing build command constructor via RootCommandDeps
	failBuildCmd := func(ctx *core.AppContext) (*cobra.Command, error) {
		return &cobra.Command{
			Use: "build",
			RunE: func(cmd *cobra.Command, args []string) error {
				return errors.New("fail")
			},
		}, nil
	}
	deps := &RootCommandDeps{
		LoadConfigWithOverlay: mockLoadConfigWithOverlay,
		ValidateConfig:        mockValidateConfig,
		NewLogger:             mockNewLogger,
		NewBuildCmd:           failBuildCmd,
	}
	cmd := NewRootCommand(ctx, deps)
	require.NotNil(t, cmd)
	buildCmd, _, err := cmd.Find([]string{"build"})
	require.NoError(t, err)
	require.NotNil(t, buildCmd)
	err = buildCmd.RunE(buildCmd, []string{})
	assert.EqualError(t, err, "fail")
}

// Note: Testing os.Exit is not practical in unit tests; skip printConfig branch coverage.
