package plugin

import (
	"context"
	"testing"
	"time"

	"github.com/srediag/srediag/internal/core"
)

// BenchmarkPluginManagerCreation tests plugin manager creation performance
func BenchmarkPluginManagerCreation(b *testing.B) {
	logger, err := core.NewLogger(&core.Logger{
		Level:  "info",
		Format: "console",
	})
	if err != nil {
		b.Fatalf("Failed to create logger: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewManager(logger, "/tmp/plugins")
	}
}

// BenchmarkPluginList tests plugin listing performance
func BenchmarkPluginList(b *testing.B) {
	logger, err := core.NewLogger(&core.Logger{
		Level:  "info",
		Format: "console",
	})
	if err != nil {
		b.Fatalf("Failed to create logger: %v", err)
	}

	manager := NewManager(logger, "/tmp/plugins")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = manager.List()
	}
}

// BenchmarkPluginHealthCheck tests health check performance
func BenchmarkPluginHealthCheck(b *testing.B) {
	logger, err := core.NewLogger(&core.Logger{
		Level:  "info",
		Format: "console",
	})
	if err != nil {
		b.Fatalf("Failed to create logger: %v", err)
	}

	manager := NewManager(logger, "/tmp/plugins")
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = manager.CheckHealth(ctx)
	}
}

// BenchmarkPluginMetadataProcessing tests metadata processing performance
func BenchmarkPluginMetadataProcessing(b *testing.B) {
	metadata := PluginMetadata{
		Name:         "test-plugin",
		Version:      "1.0.0",
		Type:         "receiver",
		Description:  "Test plugin for benchmarking",
		Capabilities: []string{"metrics", "traces", "logs"},
		SHA256:       "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// Simulate metadata validation and processing
		_ = len(metadata.Name) > 0
		_ = len(metadata.Version) > 0
		_ = len(metadata.SHA256) == 64
		_ = len(metadata.Capabilities) > 0
	}
}

// BenchmarkLoaderCreation tests loader creation performance
func BenchmarkLoaderCreation(b *testing.B) {
	logger, err := core.NewLogger(&core.Logger{
		Level:  "info",
		Format: "console",
	})
	if err != nil {
		b.Fatalf("Failed to create logger: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewLoader(logger, nil)
	}
}

// BenchmarkConcurrentPluginOperations tests concurrent plugin operations
func BenchmarkConcurrentPluginOperations(b *testing.B) {
	logger, err := core.NewLogger(&core.Logger{
		Level:  "info",
		Format: "console",
	})
	if err != nil {
		b.Fatalf("Failed to create logger: %v", err)
	}

	manager := NewManager(logger, "/tmp/plugins")
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Simulate concurrent operations
			_ = manager.List()
			_ = manager.CheckHealth(ctx)
			_, _ = manager.Get("non-existent-plugin")
		}
	})
}

// BenchmarkMemoryUsagePlugin tests memory allocation in plugin operations
func BenchmarkMemoryUsagePlugin(b *testing.B) {
	logger, err := core.NewLogger(&core.Logger{
		Level:  "info",
		Format: "console",
	})
	if err != nil {
		b.Fatalf("Failed to create logger: %v", err)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		manager := NewManager(logger, "/tmp/plugins")

		// Simulate typical operations
		_ = manager.List()
		_, _ = manager.Get("test-plugin")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		_ = manager.CheckHealth(ctx)
		cancel()
	}
}
