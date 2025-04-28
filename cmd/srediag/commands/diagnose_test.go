package commands

import (
	"bytes"
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/srediag/srediag/internal/core"
)

// Mock implementations for diagnose package functions
func mockSystemDiagnostics(_ *core.AppContext, _ *cobra.Command, _ []string) error {
	return nil
}
func mockPerformanceDiagnostics(_ *core.AppContext, _ *cobra.Command, _ []string) error {
	return nil
}
func mockSecurityDiagnostics(_ *core.AppContext, _ *cobra.Command, _ []string) error {
	return nil
}

func mockSystemDiagnosticsErr(_ *core.AppContext, _ *cobra.Command, _ []string) error {
	return errors.New("system error")
}
func mockPerformanceDiagnosticsErr(_ *core.AppContext, _ *cobra.Command, _ []string) error {
	return errors.New("performance error")
}
func mockSecurityDiagnosticsErr(_ *core.AppContext, _ *cobra.Command, _ []string) error {
	return errors.New("security error")
}

// Remove withMockedDiagnose and use dependency injection or interfaces for mocking in your command constructors.
// For example, pass the diagnostic functions as parameters to newSystemDiagCmd, newPerformanceDiagCmd, and newSecurityDiagCmd.
// Update your test cases to inject the mock implementations directly when constructing the commands.

// --- Dependency-injected subcommand constructors for testing ---
func newSystemDiagCmdWithFunc(ctx *core.AppContext, fn func(*core.AppContext, *cobra.Command, []string) error) *cobra.Command {
	return &cobra.Command{
		Use:   "system",
		Short: "Run system diagnostics",
		Long:  `Check system health, resource usage, and configuration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := fn(ctx, cmd, args)
			if err != nil {
				cmd.PrintErrf("Error running system diagnostics: %v\n", err)
			}
			return err
		},
	}
}

func newPerformanceDiagCmdWithFunc(ctx *core.AppContext, fn func(*core.AppContext, *cobra.Command, []string) error) *cobra.Command {
	return &cobra.Command{
		Use:   "performance",
		Short: "Run performance diagnostics",
		Long:  `Analyze system and application performance metrics.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := fn(ctx, cmd, args)
			if err != nil {
				cmd.PrintErrf("Error running performance diagnostics: %v\n", err)
			}
			return err
		},
	}
}

func newSecurityDiagCmdWithFunc(ctx *core.AppContext, fn func(*core.AppContext, *cobra.Command, []string) error) *cobra.Command {
	return &cobra.Command{
		Use:   "security",
		Short: "Run security diagnostics",
		Long:  `Check security configurations and potential vulnerabilities.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := fn(ctx, cmd, args)
			if err != nil {
				cmd.PrintErrf("Error running security diagnostics: %v\n", err)
			}
			return err
		},
	}
}

func TestNewDiagnoseCmd_HasSubcommands(t *testing.T) {
	ctx := &core.AppContext{}
	cmd := newDiagnoseCmd(ctx)
	subCmds := cmd.Commands()
	var names []string
	for _, c := range subCmds {
		names = append(names, c.Name())
	}
	assert.Contains(t, names, "system")
	assert.Contains(t, names, "performance")
	assert.Contains(t, names, "security")
}

func TestRunDiagnose_NoArgs_ShowsHelp(t *testing.T) {
	ctx := &core.AppContext{}
	cmd := newDiagnoseCmd(ctx)
	b := &bytes.Buffer{}
	cmd.SetOut(b)
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	assert.NoError(t, err)
	out := b.String()
	assert.Contains(t, out, "Usage:")
	assert.Contains(t, out, "diagnose [type]")
}

func TestRunDiagnose_InvalidType_ShowsError(t *testing.T) {
	ctx := &core.AppContext{}
	cmd := newDiagnoseCmd(ctx)
	b := &bytes.Buffer{}
	cmd.SetOut(b)
	cmd.SetArgs([]string{"invalidtype"})
	err := cmd.Execute()
	assert.Error(t, err)
}

func TestRunDiagnose_ValidTypes(t *testing.T) {
	validTypes := []string{"system", "performance", "security"}
	for _, typ := range validTypes {
		ctx := &core.AppContext{}
		cmd := newDiagnoseCmd(ctx)
		b := &bytes.Buffer{}
		cmd.SetOut(b)
		cmd.SetArgs([]string{typ})
		err := cmd.Execute()
		assert.NoError(t, err)
	}
}

func TestSystemDiagCmd_DelegatesToDiagnose(t *testing.T) {
	ctx := &core.AppContext{}
	cmd := newSystemDiagCmdWithFunc(ctx, mockSystemDiagnostics)
	b := &bytes.Buffer{}
	cmd.SetOut(b)
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	assert.NoError(t, err)
}

func TestPerformanceDiagCmd_DelegatesToDiagnose(t *testing.T) {
	ctx := &core.AppContext{}
	cmd := newPerformanceDiagCmdWithFunc(ctx, mockPerformanceDiagnostics)
	b := &bytes.Buffer{}
	cmd.SetOut(b)
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	assert.NoError(t, err)
}

func TestSecurityDiagCmd_DelegatesToDiagnose(t *testing.T) {
	ctx := &core.AppContext{}
	cmd := newSecurityDiagCmdWithFunc(ctx, mockSecurityDiagnostics)
	b := &bytes.Buffer{}
	cmd.SetOut(b)
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	assert.NoError(t, err)
}

func TestSystemDiagCmd_ErrorOutput(t *testing.T) {
	ctx := &core.AppContext{}
	cmd := newSystemDiagCmdWithFunc(ctx, mockSystemDiagnosticsErr)
	b := &bytes.Buffer{}
	cmd.SetOut(b)
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	require.Error(t, err)
	out := b.String()
	assert.Contains(t, out, "Usage:")
}

func TestPerformanceDiagCmd_ErrorOutput(t *testing.T) {
	ctx := &core.AppContext{}
	cmd := newPerformanceDiagCmdWithFunc(ctx, mockPerformanceDiagnosticsErr)
	b := &bytes.Buffer{}
	cmd.SetOut(b)
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	require.Error(t, err)
	out := b.String()
	assert.Contains(t, out, "Usage:")
}

func TestSecurityDiagCmd_ErrorOutput(t *testing.T) {
	ctx := &core.AppContext{}
	cmd := newSecurityDiagCmdWithFunc(ctx, mockSecurityDiagnosticsErr)
	b := &bytes.Buffer{}
	cmd.SetOut(b)
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	require.Error(t, err)
	out := b.String()
	assert.Contains(t, out, "Usage:")
}
