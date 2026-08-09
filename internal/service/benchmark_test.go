package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.opentelemetry.io/collector/component"

	"github.com/srediag/srediag/internal/core"
)

// BenchmarkServiceCreation tests service creation performance
func BenchmarkServiceCreation(b *testing.B) {
	logger, err := core.NewLogger(&core.Logger{
		Level:  "info",
		Format: "console",
	})
	if err != nil {
		b.Fatalf("Failed to create logger: %v", err)
	}

	receivers := make(map[component.Type]component.Factory)
	processors := make(map[component.Type]component.Factory)
	exporters := make(map[component.Type]component.Factory)
	extensions := make(map[component.Type]component.Factory)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewService(logger, receivers, processors, exporters, extensions)
	}
}

// BenchmarkServiceStartStop tests service start/stop performance
func BenchmarkServiceStartStop(b *testing.B) {
	logger, err := core.NewLogger(&core.Logger{
		Level:  "info",
		Format: "console",
	})
	if err != nil {
		b.Fatalf("Failed to create logger: %v", err)
	}

	receivers := make(map[component.Type]component.Factory)
	processors := make(map[component.Type]component.Factory)
	exporters := make(map[component.Type]component.Factory)
	extensions := make(map[component.Type]component.Factory)

	service := NewService(logger, receivers, processors, exporters, extensions)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.Start(ctx)
		_ = service.Stop(ctx)
	}
}

// BenchmarkConcurrentServiceOperations tests concurrent service operations
func BenchmarkConcurrentServiceOperations(b *testing.B) {
	logger, err := core.NewLogger(&core.Logger{
		Level:  "info",
		Format: "console",
	})
	if err != nil {
		b.Fatalf("Failed to create logger: %v", err)
	}

	receivers := make(map[component.Type]component.Factory)
	processors := make(map[component.Type]component.Factory)
	exporters := make(map[component.Type]component.Factory)
	extensions := make(map[component.Type]component.Factory)

	service := NewService(logger, receivers, processors, exporters, extensions)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		ctx := context.Background()
		for pb.Next() {
			// Simulate concurrent service operations
			_ = service.Start(ctx)
			_ = service.Stop(ctx)
		}
	})
}

// BenchmarkServiceWithComponents tests service performance with many components
func BenchmarkServiceWithComponents(b *testing.B) {
	logger, err := core.NewLogger(&core.Logger{
		Level:  "info",
		Format: "console",
	})
	if err != nil {
		b.Fatalf("Failed to create logger: %v", err)
	}

	// Create larger maps to simulate real-world usage
	receivers := make(map[component.Type]component.Factory, 10)
	processors := make(map[component.Type]component.Factory, 20)
	exporters := make(map[component.Type]component.Factory, 15)
	extensions := make(map[component.Type]component.Factory, 5)

	// Populate with dummy entries
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("receiver%d", i)
		typ, _ := component.NewType(key)
		receivers[typ] = nil
	}
	for i := 0; i < 20; i++ {
		key := fmt.Sprintf("processor%d", i)
		typ, _ := component.NewType(key)
		processors[typ] = nil
	}
	for i := 0; i < 15; i++ {
		key := fmt.Sprintf("exporter%d", i)
		typ, _ := component.NewType(key)
		exporters[typ] = nil
	}
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("extension%d", i)
		typ, _ := component.NewType(key)
		extensions[typ] = nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service := NewService(logger, receivers, processors, exporters, extensions)
		ctx := context.Background()
		_ = service.Start(ctx)
		_ = service.Stop(ctx)
	}
}

// BenchmarkMemoryUsageService tests memory allocation in service operations
func BenchmarkMemoryUsageService(b *testing.B) {
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
		receivers := make(map[component.Type]component.Factory)
		processors := make(map[component.Type]component.Factory)
		exporters := make(map[component.Type]component.Factory)
		extensions := make(map[component.Type]component.Factory)

		service := NewService(logger, receivers, processors, exporters, extensions)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
		_ = service.Start(ctx)
		_ = service.Stop(ctx)
		cancel()
	}
}
