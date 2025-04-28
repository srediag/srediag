package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultPluginDirAndExecDir(t *testing.T) {
	// Should return system path if running as root or /usr/bin
	os.Setenv("HOME", "/tmp/testhome")
	defer os.Unsetenv("HOME")
	// Not root, not /usr/bin
	os.Args[0] = "/tmp/srediag"
	assert.True(t, strings.HasSuffix(DefaultPluginDir(), filepath.Join(".local", "libexec", "srediag")))
	assert.True(t, strings.HasSuffix(DefaultPluginExecDir(), filepath.Join(".local", "libexec", "srediag")))
}

func TestDefaultBuildOutputDir(t *testing.T) {
	os.Setenv("HOME", "/tmp/testhome")
	defer os.Unsetenv("HOME")
	os.Args[0] = "/tmp/srediag"
	assert.True(t, strings.HasSuffix(DefaultBuildOutputDir(), filepath.Join(".srediag", "build")))
}

func TestFindConfigFile(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "srediag.yaml")
	err := os.WriteFile(file, []byte("service:\n  name: test"), 0644)
	require.NoError(t, err)
	oldPaths := configSearchPaths
	defer func() { configSearchPaths = oldPaths }()
	configSearchPaths = []string{file}
	found, ext := findConfigFile()
	assert.Equal(t, file, found)
	assert.Equal(t, ".yaml", ext)
}

func TestStrictYAMLUnmarshal(t *testing.T) {
	type S struct {
		A string `yaml:"a"`
	}
	var s S
	data := []byte("a: foo\nb: bar\n")
	err := StrictYAMLUnmarshal(data, &s)
	assert.Error(t, err, "should error on unknown field b")
}

func TestLoadConfigWithOverlay_Defaults(t *testing.T) {
	var cfg Config
	cliFlags := map[string]string{}
	err := LoadConfigWithOverlay(&cfg, cliFlags)
	assert.NoError(t, err)
	assert.Equal(t, "info", cfg.Logging.Level)
	assert.Equal(t, "console", cfg.Logging.Format)
	assert.Equal(t, 8080, cfg.Service.Port)
	assert.Equal(t, "srediag", cfg.Service.Name)
}

func TestLoadConfigWithOverlay_CLIOverride(t *testing.T) {
	var cfg Config
	cliFlags := map[string]string{"logging.level": "debug", "service.port": "9090"}
	err := LoadConfigWithOverlay(&cfg, cliFlags)
	assert.NoError(t, err)
	assert.Equal(t, "debug", cfg.Logging.Level)
	assert.Equal(t, 9090, cfg.Service.Port)
}

func TestValidateConfig(t *testing.T) {
	cfg := &Config{
		Plugins: PluginsConfig{Dir: "/tmp/plugins"},
		Logging: LoggingConfig{Level: "info"},
		Service: ServiceConfig{Port: 1234},
	}
	assert.NoError(t, ValidateConfig(cfg))

	cfg.Plugins.Dir = ""
	assert.Error(t, ValidateConfig(cfg))
	cfg.Plugins.Dir = "/tmp/plugins"
	cfg.Logging.Level = ""
	assert.Error(t, ValidateConfig(cfg))
	cfg.Logging.Level = "info"
	cfg.Service.Port = 0
	assert.Error(t, ValidateConfig(cfg))
}

func TestPrintEffectiveConfig(t *testing.T) {
	cfg := &Config{
		Service: ServiceConfig{Name: "test", Port: 1234},
	}
	err := PrintEffectiveConfig(cfg)
	assert.NoError(t, err)
}

func TestWithConfigPathAndEnvPrefixAndSuffix(t *testing.T) {
	opt := WithConfigPath("foo.yaml")
	var o configSpecOpts
	opt(&o)
	assert.Equal(t, "foo.yaml", o.Path)

	opt2 := WithEnvPrefix("FOO_")
	opt2(&o)
	assert.Equal(t, "FOO_", o.EnvPrefix)

	opt3 := WithConfigPathSuffix("bar")
	opt3(&o)
	assert.Equal(t, "bar", o.PathSuffix)
}

func TestNewConfig(t *testing.T) {
	cfg := NewConfig()
	assert.NotNil(t, cfg)
}
