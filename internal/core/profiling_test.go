package core

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProfilerCreation(t *testing.T) {
	logger, err := NewLogger(&Logger{
		Level:  "info",
		Format: "console",
	})
	require.NoError(t, err)

	profiler := NewProfiler(nil, logger)
	assert.NotNil(t, profiler)
	assert.False(t, profiler.config.Enabled)
}

func TestProfilerWithConfig(t *testing.T) {
	logger, err := NewLogger(&Logger{
		Level:  "info",
		Format: "console",
	})
	require.NoError(t, err)

	config := &ProfilerConfig{
		Enabled:     true,
		Port:        6061,
		ProfilePath: "/tmp/test-profiles",
		Duration:    "5s",
	}

	profiler := NewProfiler(config, logger)
	assert.NotNil(t, profiler)
	assert.True(t, profiler.config.Enabled)
	assert.Equal(t, 6061, profiler.config.Port)
}

func TestMemoryStats(t *testing.T) {
	logger, err := NewLogger(&Logger{
		Level:  "info",
		Format: "console",
	})
	require.NoError(t, err)

	profiler := NewProfiler(nil, logger)
	stats := profiler.GetMemoryStats()

	assert.Contains(t, stats, "alloc_mb")
	assert.Contains(t, stats, "num_goroutines")
	assert.Contains(t, stats, "num_gc")

	// Memory stats should be reasonable
	assert.Greater(t, stats["alloc_mb"].(float64), 0.0)
	assert.Greater(t, stats["num_goroutines"].(int), 0)
}

func TestProfilerStartStop(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping profiler test in short mode")
	}

	logger, err := NewLogger(&Logger{
		Level:  "info",
		Format: "console",
	})
	require.NoError(t, err)

	// Create temporary directory for profiles
	tmpDir := "/tmp/test-profiler-" + time.Now().Format("20060102-150405")
	defer os.RemoveAll(tmpDir)

	config := &ProfilerConfig{
		Enabled:     true,
		Port:        6062, // Use different port to avoid conflicts
		ProfilePath: tmpDir,
		Duration:    "1s",
	}

	profiler := NewProfiler(config, logger)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Start profiler
	err = profiler.Start(ctx)
	assert.NoError(t, err)

	// Wait a bit for server to start
	time.Sleep(100 * time.Millisecond)

	// Stop profiler
	err = profiler.Stop(ctx)
	assert.NoError(t, err)
}

func TestCleanupOldProfiles(t *testing.T) {
	logger, err := NewLogger(&Logger{
		Level:  "info",
		Format: "console",
	})
	require.NoError(t, err)

	// Create temporary directory
	tmpDir := "/tmp/test-cleanup-" + time.Now().Format("20060102-150405")
	defer os.RemoveAll(tmpDir)

	err = os.MkdirAll(tmpDir, 0755)
	require.NoError(t, err)

	config := &ProfilerConfig{
		Enabled:     true,
		ProfilePath: tmpDir,
	}

	profiler := NewProfiler(config, logger)

	// Create some test files
	testFile1 := tmpDir + "/old-profile.prof"
	testFile2 := tmpDir + "/new-profile.prof"

	// Create old file
	f1, err := os.Create(testFile1)
	require.NoError(t, err)
	f1.Close()

	// Set old timestamp
	oldTime := time.Now().Add(-25 * time.Hour)
	err = os.Chtimes(testFile1, oldTime, oldTime)
	require.NoError(t, err)

	// Create new file
	f2, err := os.Create(testFile2)
	require.NoError(t, err)
	f2.Close()

	// Cleanup files older than 24 hours
	err = profiler.CleanupOldProfiles(24 * time.Hour)
	assert.NoError(t, err)

	// Old file should be gone, new file should remain
	_, err = os.Stat(testFile1)
	assert.True(t, os.IsNotExist(err))

	_, err = os.Stat(testFile2)
	assert.NoError(t, err)
}
