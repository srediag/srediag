package core

import (
	"testing"
	"time"

	"go.uber.org/zap"
)

// BenchmarkLogger tests logger performance
func BenchmarkLogger(b *testing.B) {
	logger, err := NewLogger(&Logger{
		Level:  "info",
		Format: "json",
	})
	if err != nil {
		b.Fatalf("Failed to create logger: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("Benchmark test message",
			zap.String("key1", "value1"),
			zap.Int("key2", i),
			zap.Duration("duration", time.Millisecond*100))
	}
}

// BenchmarkLoggerParallel tests logger performance under concurrent load
func BenchmarkLoggerParallel(b *testing.B) {
	logger, err := NewLogger(&Logger{
		Level:  "info",
		Format: "json",
	})
	if err != nil {
		b.Fatalf("Failed to create logger: %v", err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			logger.Info("Parallel benchmark test message",
				zap.String("worker", "worker"),
				zap.Int("iteration", i),
				zap.Duration("duration", time.Microsecond*50))
			i++
		}
	})
}

// BenchmarkConfigValidation tests config validation performance
func BenchmarkConfigValidation(b *testing.B) {
	config := &Config{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ValidateConfig(config)
	}
}

// BenchmarkComponentManager tests component manager performance
func BenchmarkComponentManager(b *testing.B) {
	logger, err := NewLogger(&Logger{
		Level:  "info",
		Format: "console",
	})
	if err != nil {
		b.Fatalf("Failed to create logger: %v", err)
	}

	manager := NewComponentManager(logger)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = manager
	}
}

// BenchmarkBuildInfoGeneration tests build info generation performance
func BenchmarkBuildInfoGeneration(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buildInfo := &BuildInfo{
			Version: "1.0.0",
			Commit:  "abc123def456",
			Date:    time.Now().Format(time.RFC3339),
		}
		_ = buildInfo.Version
	}
}

// BenchmarkMemoryUsage tests memory allocation patterns
func BenchmarkMemoryUsage(b *testing.B) {
	logger, err := NewLogger(&Logger{
		Level:  "info",
		Format: "json",
	})
	if err != nil {
		b.Fatalf("Failed to create logger: %v", err)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// Simulate typical usage patterns
		config := &Config{}

		_ = ValidateConfig(config)

		logger.Info("Processing request",
			zap.String("request_id", "req-12345"),
			zap.Int("request_number", i),
			zap.Duration("processing_time", time.Millisecond*10))
	}
}
