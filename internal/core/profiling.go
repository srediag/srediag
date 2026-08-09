package core

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof" // Import pprof for profiling endpoints
	"os"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"time"

	"go.uber.org/zap"
)

// ProfilerConfig holds configuration for performance profiling
type ProfilerConfig struct {
	Enabled     bool   `yaml:"enabled"`
	Port        int    `yaml:"port"`
	ProfilePath string `yaml:"profile_path"`
	Duration    string `yaml:"duration"`
}

// Profiler manages performance profiling for SREDIAG
type Profiler struct {
	config *ProfilerConfig
	logger *Logger
	server *http.Server
}

// NewProfiler creates a new performance profiler
func NewProfiler(config *ProfilerConfig, logger *Logger) *Profiler {
	if config == nil {
		config = &ProfilerConfig{
			Enabled:     false,
			Port:        6060,
			ProfilePath: "/tmp/srediag-profiles",
			Duration:    "30s",
		}
	}

	return &Profiler{
		config: config,
		logger: logger,
	}
}

// Start starts the profiling server and collection
func (p *Profiler) Start(ctx context.Context) error {
	if !p.config.Enabled {
		p.logger.Debug("Profiler disabled, skipping start")
		return nil
	}

	// Ensure profile directory exists
	if err := os.MkdirAll(p.config.ProfilePath, 0755); err != nil {
		return fmt.Errorf("failed to create profile directory: %w", err)
	}

	// Start pprof HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", http.DefaultServeMux.ServeHTTP)

	p.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", p.config.Port),
		Handler: mux,
	}

	go func() {
		p.logger.Info("Starting profiler server",
			zap.Int("port", p.config.Port),
			zap.String("profile_path", p.config.ProfilePath))

		if err := p.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			p.logger.Error("Profiler server error", zap.Error(err))
		}
	}()

	// Start automatic profile collection
	go p.collectProfiles(ctx)

	return nil
}

// Stop stops the profiling server
func (p *Profiler) Stop(ctx context.Context) error {
	if p.server == nil {
		return nil
	}

	p.logger.Info("Stopping profiler server")
	return p.server.Shutdown(ctx)
}

// collectProfiles automatically collects performance profiles
func (p *Profiler) collectProfiles(ctx context.Context) {
	duration, err := time.ParseDuration(p.config.Duration)
	if err != nil {
		p.logger.Error("Invalid profile duration", zap.Error(err))
		return
	}

	ticker := time.NewTicker(duration * 2) // Collect every 2x duration
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.collectCPUProfile(duration)
			p.collectMemProfile()
			p.collectTraceProfile(duration)
		}
	}
}

// createProfileFile creates a profile file with timestamp
func (p *Profiler) createProfileFile(profileType, extension string) (*os.File, string, error) {
	timestamp := time.Now().Format("20060102-150405")
	filename := fmt.Sprintf("%s/%s-profile-%s.%s", p.config.ProfilePath, profileType, timestamp, extension)

	f, err := os.Create(filename)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create %s profile file: %w", profileType, err)
	}

	return f, filename, nil
}

// collectCPUProfile collects CPU profile
func (p *Profiler) collectCPUProfile(duration time.Duration) {
	f, filename, err := p.createProfileFile("cpu", "prof")
	if err != nil {
		p.logger.Error("Failed to create CPU profile file", zap.Error(err))
		return
	}
	defer f.Close()

	p.logger.Debug("Starting CPU profile collection",
		zap.String("file", filename),
		zap.Duration("duration", duration))

	if err := pprof.StartCPUProfile(f); err != nil {
		p.logger.Error("Failed to start CPU profile", zap.Error(err))
		return
	}

	time.Sleep(duration)
	pprof.StopCPUProfile()

	p.logger.Info("CPU profile collected", zap.String("file", filename))
}

// collectMemProfile collects memory profile
func (p *Profiler) collectMemProfile() {
	f, filename, err := p.createProfileFile("mem", "prof")
	if err != nil {
		p.logger.Error("Failed to create memory profile file", zap.Error(err))
		return
	}
	defer f.Close()

	runtime.GC() // Force GC before collecting memory profile

	if err := pprof.WriteHeapProfile(f); err != nil {
		p.logger.Error("Failed to write memory profile", zap.Error(err))
		return
	}

	p.logger.Info("Memory profile collected", zap.String("file", filename))
}

// collectTraceProfile collects execution trace
func (p *Profiler) collectTraceProfile(duration time.Duration) {
	f, filename, err := p.createProfileFile("trace", "trace")
	if err != nil {
		p.logger.Error("Failed to create trace profile file", zap.Error(err))
		return
	}
	defer f.Close()

	p.logger.Debug("Starting trace collection",
		zap.String("file", filename),
		zap.Duration("duration", duration))

	if err := trace.Start(f); err != nil {
		p.logger.Error("Failed to start trace", zap.Error(err))
		return
	}

	time.Sleep(duration)
	trace.Stop()

	p.logger.Info("Trace profile collected", zap.String("file", filename))
}

// GetMemoryStats returns current memory statistics
func (p *Profiler) GetMemoryStats() map[string]interface{} {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return map[string]interface{}{
		"alloc_mb":        bToMb(m.Alloc),
		"total_alloc_mb":  bToMb(m.TotalAlloc),
		"sys_mb":          bToMb(m.Sys),
		"heap_alloc_mb":   bToMb(m.HeapAlloc),
		"heap_sys_mb":     bToMb(m.HeapSys),
		"heap_inuse_mb":   bToMb(m.HeapInuse),
		"stack_inuse_mb":  bToMb(m.StackInuse),
		"num_gc":          m.NumGC,
		"gc_cpu_fraction": m.GCCPUFraction,
		"num_goroutines":  runtime.NumGoroutine(),
	}
}

// LogMemoryStats logs current memory statistics
func (p *Profiler) LogMemoryStats() {
	stats := p.GetMemoryStats()
	p.logger.Info("Memory statistics",
		zap.Float64("alloc_mb", stats["alloc_mb"].(float64)),
		zap.Float64("heap_alloc_mb", stats["heap_alloc_mb"].(float64)),
		zap.Uint32("num_gc", stats["num_gc"].(uint32)),
		zap.Int("num_goroutines", stats["num_goroutines"].(int)))
}

// bToMb converts bytes to megabytes
func bToMb(b uint64) float64 {
	return float64(b) / 1024 / 1024
}

// CleanupOldProfiles removes profile files older than specified duration
func (p *Profiler) CleanupOldProfiles(maxAge time.Duration) error {
	if !p.config.Enabled {
		return nil
	}

	entries, err := os.ReadDir(p.config.ProfilePath)
	if err != nil {
		return fmt.Errorf("failed to read profile directory: %w", err)
	}

	now := time.Now()
	cleaned := 0

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		if now.Sub(info.ModTime()) > maxAge {
			filepath := fmt.Sprintf("%s/%s", p.config.ProfilePath, entry.Name())
			if err := os.Remove(filepath); err != nil {
				p.logger.Warn("Failed to remove old profile",
					zap.String("file", filepath),
					zap.Error(err))
			} else {
				cleaned++
			}
		}
	}

	if cleaned > 0 {
		p.logger.Info("Cleaned up old profiles", zap.Int("count", cleaned))
	}

	return nil
}
