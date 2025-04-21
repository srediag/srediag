// tests/plugins_test.go
package tests

import (
	"testing"

	"github.com/srediag/srediag/internal/plugins"
)

func TestPluginRegistry(t *testing.T) {
	registry := plugins.NewPluginRegistry()
	err := registry.Register("dummy", func() plugins.Plugin { return nil })
	if err != nil {
		t.Fatalf("failed to register plugin: %v", err)
	}
	if !registry.Exists("dummy") {
		t.Errorf("plugin ‘dummy’ should exist")
	}
}
