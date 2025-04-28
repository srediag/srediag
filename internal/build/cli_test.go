package build

import (
	"bytes"
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/srediag/srediag/internal/core"
)

// --- Tests ---

type mockBuildManager struct {
	buildAllErr       error
	buildPluginErr    error
	generateErr       error
	installPluginsErr error
}

func (m *mockBuildManager) BuildAll() error               { return m.buildAllErr }
func (m *mockBuildManager) BuildPlugin(t, n string) error { return m.buildPluginErr }
func (m *mockBuildManager) Generate(t, n string) error    { return m.generateErr }
func (m *mockBuildManager) InstallPlugins() error         { return m.installPluginsErr }

func TestCLI_BuildAll_Success(t *testing.T) {
	viper.Set("build.output_dir", "/tmp/test-build")
	defer viper.Reset()
	var buf bytes.Buffer
	logger := core.NewTestLogger(&buf)
	origNewBuildManagerFunc := newBuildManagerFunc
	newBuildManagerFunc = func(logger *core.Logger, outputDir string) BuildManagerInterface {
		return &mockBuildManager{}
	}
	defer func() { newBuildManagerFunc = origNewBuildManagerFunc }()
	ctx := &core.AppContext{Logger: logger}
	cmd := &cobra.Command{}
	err := CLI_BuildAll(ctx, cmd, nil)
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "BuildAll completed successfully")
}

func TestCLI_BuildAll_Error(t *testing.T) {
	viper.Set("build.output_dir", "")
	defer viper.Reset()
	var buf bytes.Buffer
	logger := core.NewTestLogger(&buf)
	origNewBuildManagerFunc := newBuildManagerFunc
	newBuildManagerFunc = func(logger *core.Logger, outputDir string) BuildManagerInterface {
		return &mockBuildManager{buildAllErr: errors.New("fail")}
	}
	defer func() { newBuildManagerFunc = origNewBuildManagerFunc }()
	ctx := &core.AppContext{Logger: logger}
	cmd := &cobra.Command{}
	err := CLI_BuildAll(ctx, cmd, nil)
	assert.Error(t, err)
	assert.Contains(t, buf.String(), "BuildAll failed")
}

func TestCLI_BuildPlugin_Success(t *testing.T) {
	viper.Set("build.output_dir", "/tmp/test-build")
	defer viper.Reset()
	var buf bytes.Buffer
	logger := core.NewTestLogger(&buf)
	origNewBuildManagerFunc := newBuildManagerFunc
	newBuildManagerFunc = func(logger *core.Logger, outputDir string) BuildManagerInterface {
		return &mockBuildManager{}
	}
	defer func() { newBuildManagerFunc = origNewBuildManagerFunc }()
	ctx := &core.AppContext{Logger: logger}
	cmd := &cobra.Command{}
	cmd.Flags().String("type", "foo", "")
	cmd.Flags().String("name", "bar", "")
	err := cmd.Flags().Set("type", "foo")
	require.NoError(t, err)
	err = cmd.Flags().Set("name", "bar")
	require.NoError(t, err)
	err = CLI_BuildPlugin(ctx, cmd, nil)
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "BuildPlugin completed successfully")
}

func TestCLI_BuildPlugin_MissingType(t *testing.T) {
	var buf bytes.Buffer
	logger := core.NewTestLogger(&buf)
	ctx := &core.AppContext{Logger: logger}
	cmd := &cobra.Command{}
	cmd.Flags().String("type", "", "")
	cmd.Flags().String("name", "bar", "")
	err := cmd.Flags().Set("name", "bar")
	require.NoError(t, err)
	err = CLI_BuildPlugin(ctx, cmd, nil)
	assert.ErrorContains(t, err, "--type flag is required")
}

func TestCLI_BuildPlugin_MissingName(t *testing.T) {
	var buf bytes.Buffer
	logger := core.NewTestLogger(&buf)
	ctx := &core.AppContext{Logger: logger}
	cmd := &cobra.Command{}
	cmd.Flags().String("type", "foo", "")
	cmd.Flags().String("name", "", "")
	err := cmd.Flags().Set("type", "foo")
	require.NoError(t, err)
	err = CLI_BuildPlugin(ctx, cmd, nil)
	assert.ErrorContains(t, err, "--name flag is required")
}

func TestCLI_BuildPlugin_Error(t *testing.T) {
	viper.Set("build.output_dir", "")
	defer viper.Reset()
	var buf bytes.Buffer
	logger := core.NewTestLogger(&buf)
	origNewBuildManagerFunc := newBuildManagerFunc
	newBuildManagerFunc = func(logger *core.Logger, outputDir string) BuildManagerInterface {
		return &mockBuildManager{buildPluginErr: errors.New("fail")}
	}
	defer func() { newBuildManagerFunc = origNewBuildManagerFunc }()
	ctx := &core.AppContext{Logger: logger}
	cmd := &cobra.Command{}
	cmd.Flags().String("type", "foo", "")
	cmd.Flags().String("name", "bar", "")
	err := cmd.Flags().Set("type", "foo")
	require.NoError(t, err)
	err = cmd.Flags().Set("name", "bar")
	require.NoError(t, err)
	err = CLI_BuildPlugin(ctx, cmd, nil)
	assert.Error(t, err)
	assert.Contains(t, buf.String(), "BuildPlugin failed")
}

func TestCLI_Generate_Success(t *testing.T) {
	viper.Set("build.output_dir", "/tmp/test-build")
	defer viper.Reset()
	var buf bytes.Buffer
	logger := core.NewTestLogger(&buf)
	origNewBuildManagerFunc := newBuildManagerFunc
	newBuildManagerFunc = func(logger *core.Logger, outputDir string) BuildManagerInterface {
		return &mockBuildManager{}
	}
	defer func() { newBuildManagerFunc = origNewBuildManagerFunc }()
	ctx := &core.AppContext{Logger: logger}
	cmd := &cobra.Command{}
	cmd.Flags().String("type", "foo", "")
	cmd.Flags().String("name", "bar", "")
	err := cmd.Flags().Set("type", "foo")
	require.NoError(t, err)
	err = cmd.Flags().Set("name", "bar")
	require.NoError(t, err)
	err = CLI_Generate(ctx, cmd, nil)
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "Generate completed successfully")
}

func TestCLI_Generate_Error(t *testing.T) {
	viper.Set("build.output_dir", "")
	defer viper.Reset()
	var buf bytes.Buffer
	logger := core.NewTestLogger(&buf)
	origNewBuildManagerFunc := newBuildManagerFunc
	newBuildManagerFunc = func(logger *core.Logger, outputDir string) BuildManagerInterface {
		return &mockBuildManager{generateErr: errors.New("fail")}
	}
	defer func() { newBuildManagerFunc = origNewBuildManagerFunc }()
	ctx := &core.AppContext{Logger: logger}
	cmd := &cobra.Command{}
	cmd.Flags().String("type", "foo", "")
	cmd.Flags().String("name", "bar", "")
	err := cmd.Flags().Set("type", "foo")
	require.NoError(t, err)
	err = cmd.Flags().Set("name", "bar")
	require.NoError(t, err)
	err = CLI_Generate(ctx, cmd, nil)
	assert.Error(t, err)
	assert.Contains(t, buf.String(), "Generate failed")
}

func TestCLI_InstallPlugins_Success(t *testing.T) {
	viper.Set("build.output_dir", "/tmp/test-build")
	defer viper.Reset()
	var buf bytes.Buffer
	logger := core.NewTestLogger(&buf)
	origNewBuildManagerFunc := newBuildManagerFunc
	newBuildManagerFunc = func(logger *core.Logger, outputDir string) BuildManagerInterface {
		return &mockBuildManager{}
	}
	defer func() { newBuildManagerFunc = origNewBuildManagerFunc }()
	ctx := &core.AppContext{Logger: logger}
	cmd := &cobra.Command{}
	err := CLI_InstallPlugins(ctx, cmd, nil)
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "InstallPlugins completed successfully")
}

func TestCLI_InstallPlugins_Error(t *testing.T) {
	viper.Set("build.output_dir", "")
	defer viper.Reset()
	var buf bytes.Buffer
	logger := core.NewTestLogger(&buf)
	origNewBuildManagerFunc := newBuildManagerFunc
	newBuildManagerFunc = func(logger *core.Logger, outputDir string) BuildManagerInterface {
		return &mockBuildManager{installPluginsErr: errors.New("fail")}
	}
	defer func() { newBuildManagerFunc = origNewBuildManagerFunc }()
	ctx := &core.AppContext{Logger: logger}
	cmd := &cobra.Command{}
	err := CLI_InstallPlugins(ctx, cmd, nil)
	assert.Error(t, err)
	assert.Contains(t, buf.String(), "InstallPlugins failed")
}

func TestCLI_UpdateBuilder_Success(t *testing.T) {
	viper.Set("build.output_dir", "/tmp/test-build")
	defer viper.Reset()
	var buf bytes.Buffer
	logger := core.NewTestLogger(&buf)
	ctx := &core.AppContext{Logger: logger}
	cmd := &cobra.Command{}
	cmd.Flags().String("yaml", "foo.yaml", "")
	cmd.Flags().String("gomod", "bar.mod", "")
	cmd.Flags().String("plugin-gen", "baz", "")
	err := cmd.Flags().Set("yaml", "foo.yaml")
	require.NoError(t, err)
	err = cmd.Flags().Set("gomod", "bar.mod")
	require.NoError(t, err)
	err = cmd.Flags().Set("plugin-gen", "baz")
	require.NoError(t, err)
	origUpdateBuilderFunc := updateBuilderFunc
	updateBuilderFunc = func(yaml, gomod, pluginGen string) error {
		require.Equal(t, "foo.yaml", yaml)
		require.Equal(t, "bar.mod", gomod)
		require.Equal(t, "baz", pluginGen)
		return nil
	}
	defer func() { updateBuilderFunc = origUpdateBuilderFunc }()
	err = CLI_UpdateBuilder(ctx, cmd, nil)
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "UpdateBuilder completed successfully")
}

func TestCLI_UpdateBuilder_MissingFlags(t *testing.T) {
	var buf bytes.Buffer
	logger := core.NewTestLogger(&buf)
	ctx := &core.AppContext{Logger: logger}

	origUpdateBuilderFunc := updateBuilderFunc
	updateBuilderFunc = func(yaml, gomod, pluginGen string) error { return nil }
	defer func() { updateBuilderFunc = origUpdateBuilderFunc }()

	// Missing plugin-gen
	cmd := &cobra.Command{}
	cmd.Flags().String("yaml", "foo.yaml", "")
	cmd.Flags().String("gomod", "bar.mod", "")
	cmd.Flags().String("plugin-gen", "", "")
	err := CLI_UpdateBuilder(ctx, cmd, nil)
	assert.ErrorContains(t, err, "--plugin-gen flag is required")

	// Missing gomod
	cmd = &cobra.Command{}
	cmd.Flags().String("yaml", "foo.yaml", "")
	cmd.Flags().String("gomod", "", "")
	cmd.Flags().String("plugin-gen", "baz", "")
	err = CLI_UpdateBuilder(ctx, cmd, nil)
	assert.ErrorContains(t, err, "--gomod flag is required")

	// Missing yaml
	cmd = &cobra.Command{}
	cmd.Flags().String("yaml", "", "")
	cmd.Flags().String("gomod", "bar.mod", "")
	cmd.Flags().String("plugin-gen", "baz", "")
	err = CLI_UpdateBuilder(ctx, cmd, nil)
	assert.ErrorContains(t, err, "--yaml flag is required")
}

func TestCLI_UpdateBuilder_Error(t *testing.T) {
	viper.Set("build.output_dir", "")
	defer viper.Reset()
	var buf bytes.Buffer
	logger := core.NewTestLogger(&buf)
	ctx := &core.AppContext{Logger: logger}
	cmd := &cobra.Command{}
	cmd.Flags().String("yaml", "foo.yaml", "")
	cmd.Flags().String("gomod", "bar.mod", "")
	cmd.Flags().String("plugin-gen", "baz", "")
	err := cmd.Flags().Set("yaml", "foo.yaml")
	require.NoError(t, err)
	err = cmd.Flags().Set("gomod", "bar.mod")
	require.NoError(t, err)
	err = cmd.Flags().Set("plugin-gen", "baz")
	require.NoError(t, err)
	origUpdateBuilderFunc := updateBuilderFunc
	updateBuilderFunc = func(yaml, gomod, pluginGen string) error {
		return errors.New("fail")
	}
	defer func() { updateBuilderFunc = origUpdateBuilderFunc }()
	err = CLI_UpdateBuilder(ctx, cmd, nil)
	assert.Error(t, err)
	assert.Contains(t, buf.String(), "UpdateBuilder failed")
}
