// tests/plugins_test.go
package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/srediag/srediag/internal/plugins"
)

// DummyPlugin implementa a interface Plugin para testes
type DummyPlugin struct{}

func (d *DummyPlugin) Start(ctx context.Context) error {
	return nil
}

func (d *DummyPlugin) Stop(ctx context.Context) error {
	return nil
}

func TestPluginRegistry(t *testing.T) {
	// Criar um novo registro de plugins
	registry := plugins.NewPluginRegistry()

	// Registrar o plugin dummy
	registry.Register("dummy", &DummyPlugin{})

	// Verificar se o plugin foi registrado corretamente
	assert.True(t, registry.Exists("dummy"), "plugin 'dummy' deveria existir")

	// Obter o plugin registrado
	plugin := registry.Get("dummy")
	assert.NotNil(t, plugin, "plugin 'dummy' não deveria ser nil")

	// Verificar se o plugin está na lista de plugins
	plugins := registry.List()
	assert.Contains(t, plugins, "dummy", "plugin 'dummy' deveria estar na lista de plugins")
}
