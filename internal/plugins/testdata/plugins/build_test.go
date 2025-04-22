package plugins_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	// Compila o plugin mock antes de executar os testes
	if err := buildMockPlugin(); err != nil {
		panic(err)
	}

	// Executa os testes
	code := m.Run()

	// Limpa os arquivos gerados
	cleanupMockPlugin()

	os.Exit(code)
}

func buildMockPlugin() error {
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", "mock_plugin.so", "mock_plugin.go")
	cmd.Dir = "testdata/plugins"
	return cmd.Run()
}

func cleanupMockPlugin() {
	os.Remove(filepath.Join("testdata", "plugins", "mock_plugin.so"))
}
