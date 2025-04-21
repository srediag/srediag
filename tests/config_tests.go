// tests/config_test.go
package tests

import (
	"os"
	"testing"

	"github.com/srediag/srediag/internal/config"
)

func TestLoadDefaultConfig(t *testing.T) {
	os.Setenv("CONFIG_PATH", "configs/config.yaml") // fixture
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.Port != 8080 {
		t.Errorf("expected default port 8080, got %d", cfg.Port)
	}
}
