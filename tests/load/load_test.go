//go:build load
// +build load

package load

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// LoadTestConfig holds configuration for load testing
type LoadTestConfig struct {
	Endpoint    string
	Duration    time.Duration
	Concurrency int
	RequestRate int
}

// TestOTLPLoad tests OTLP endpoint under load
func TestOTLPLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load tests in short mode")
	}

	config := LoadTestConfig{
		Endpoint:    "http://localhost:4318/v1/traces",
		Duration:    30 * time.Second,
		Concurrency: 10,
		RequestRate: 100, // requests per second
	}

	// Generate test trace data
	traceData := generateTestTrace()

	ctx, cancel := context.WithTimeout(context.Background(), config.Duration+10*time.Second)
	defer cancel()

	// Statistics
	var totalRequests int64
	var successfulRequests int64
	var failedRequests int64
	var totalLatency time.Duration
	var mu sync.Mutex

	// Worker function
	worker := func(workerID int) {
		client := &http.Client{Timeout: 5 * time.Second}

		ticker := time.NewTicker(time.Second / time.Duration(config.RequestRate/config.Concurrency))
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				start := time.Now()

				// Send request
				resp, err := client.Post(config.Endpoint, "application/json", bytes.NewReader(traceData))

				latency := time.Since(start)

				mu.Lock()
				totalRequests++
				totalLatency += latency

				if err != nil || resp.StatusCode >= 400 {
					failedRequests++
					if err != nil {
						t.Logf("Worker %d: Request failed: %v", workerID, err)
					} else {
						t.Logf("Worker %d: Request failed with status: %d", workerID, resp.StatusCode)
						resp.Body.Close()
					}
				} else {
					successfulRequests++
					resp.Body.Close()
				}
				mu.Unlock()
			}
		}
	}

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < config.Concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			worker(id)
		}(i)
	}

	// Wait for test duration
	time.Sleep(config.Duration)
	cancel()
	wg.Wait()

	// Calculate statistics
	successRate := float64(successfulRequests) / float64(totalRequests) * 100
	avgLatency := totalLatency / time.Duration(totalRequests)
	requestsPerSecond := float64(totalRequests) / config.Duration.Seconds()

	// Log results
	t.Logf("Load Test Results:")
	t.Logf("  Duration: %v", config.Duration)
	t.Logf("  Concurrency: %d", config.Concurrency)
	t.Logf("  Total Requests: %d", totalRequests)
	t.Logf("  Successful: %d", successfulRequests)
	t.Logf("  Failed: %d", failedRequests)
	t.Logf("  Success Rate: %.2f%%", successRate)
	t.Logf("  Average Latency: %v", avgLatency)
	t.Logf("  Requests/Second: %.2f", requestsPerSecond)

	// Assertions
	assert.Greater(t, successRate, 95.0, "Success rate should be > 95%")
	assert.Less(t, avgLatency, 100*time.Millisecond, "Average latency should be < 100ms")
	assert.Greater(t, requestsPerSecond, float64(config.RequestRate)*0.8, "Request rate should be >= 80% of target")
}

// TestMemoryUsageUnderLoad tests memory usage during high load
func TestMemoryUsageUnderLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load tests in short mode")
	}

	// This test would require access to SREDIAG's memory metrics
	// For now, we'll test that the metrics endpoint is responsive under load

	config := LoadTestConfig{
		Endpoint:    "http://localhost:8889/metrics",
		Duration:    10 * time.Second,
		Concurrency: 5,
		RequestRate: 20,
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.Duration+5*time.Second)
	defer cancel()

	var successfulRequests int64
	var mu sync.Mutex

	worker := func() {
		client := &http.Client{Timeout: 2 * time.Second}

		for {
			select {
			case <-ctx.Done():
				return
			default:
				resp, err := client.Get(config.Endpoint)
				if err == nil && resp.StatusCode == 200 {
					mu.Lock()
					successfulRequests++
					mu.Unlock()
					resp.Body.Close()
				}
				time.Sleep(time.Second / time.Duration(config.RequestRate/config.Concurrency))
			}
		}
	}

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < config.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker()
		}()
	}

	time.Sleep(config.Duration)
	cancel()
	wg.Wait()

	t.Logf("Metrics endpoint requests: %d", successfulRequests)
	assert.Greater(t, successfulRequests, int64(10), "Metrics endpoint should be responsive")
}

// TestConcurrentConnections tests handling of many concurrent connections
func TestConcurrentConnections(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load tests in short mode")
	}

	const numConnections = 100
	const requestsPerConnection = 10

	traceData := generateTestTrace()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var successCount int64
	var mu sync.Mutex

	connectionWorker := func(connID int) {
		client := &http.Client{Timeout: 5 * time.Second}

		for i := 0; i < requestsPerConnection; i++ {
			select {
			case <-ctx.Done():
				return
			default:
				resp, err := client.Post("http://localhost:4318/v1/traces",
					"application/json",
					bytes.NewReader(traceData))

				if err == nil && resp.StatusCode < 400 {
					mu.Lock()
					successCount++
					mu.Unlock()
					resp.Body.Close()
				}

				// Small delay between requests
				time.Sleep(10 * time.Millisecond)
			}
		}
	}

	// Start all connections
	var wg sync.WaitGroup
	start := time.Now()

	for i := 0; i < numConnections; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			connectionWorker(id)
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)

	expectedRequests := numConnections * requestsPerConnection
	successRate := float64(successCount) / float64(expectedRequests) * 100

	t.Logf("Concurrent Connections Test:")
	t.Logf("  Connections: %d", numConnections)
	t.Logf("  Requests per connection: %d", requestsPerConnection)
	t.Logf("  Total expected requests: %d", expectedRequests)
	t.Logf("  Successful requests: %d", successCount)
	t.Logf("  Success rate: %.2f%%", successRate)
	t.Logf("  Duration: %v", duration)

	assert.Greater(t, successRate, 90.0, "Success rate should be > 90%")
}

// generateTestTrace creates a realistic test trace for load testing
func generateTestTrace() []byte {
	trace := map[string]interface{}{
		"resourceSpans": []map[string]interface{}{
			{
				"resource": map[string]interface{}{
					"attributes": []map[string]interface{}{
						{
							"key": "service.name",
							"value": map[string]interface{}{
								"stringValue": "load-test-service",
							},
						},
						{
							"key": "service.version",
							"value": map[string]interface{}{
								"stringValue": "1.0.0",
							},
						},
					},
				},
				"scopeSpans": []map[string]interface{}{
					{
						"spans": []map[string]interface{}{
							{
								"traceId":           fmt.Sprintf("%032d", time.Now().UnixNano()),
								"spanId":            fmt.Sprintf("%016d", time.Now().UnixNano()),
								"name":              "load-test-operation",
								"kind":              "SPAN_KIND_SERVER",
								"startTimeUnixNano": fmt.Sprintf("%d", time.Now().UnixNano()),
								"endTimeUnixNano":   fmt.Sprintf("%d", time.Now().Add(100*time.Millisecond).UnixNano()),
								"attributes": []map[string]interface{}{
									{
										"key": "http.method",
										"value": map[string]interface{}{
											"stringValue": "POST",
										},
									},
									{
										"key": "http.status_code",
										"value": map[string]interface{}{
											"intValue": "200",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	data, _ := json.Marshal(trace)
	return data
}

// BenchmarkOTLPThroughput benchmarks OTLP endpoint throughput
func BenchmarkOTLPThroughput(b *testing.B) {
	traceData := generateTestTrace()
	client := &http.Client{Timeout: 5 * time.Second}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp, err := client.Post("http://localhost:4318/v1/traces",
				"application/json",
				bytes.NewReader(traceData))

			if err != nil {
				b.Errorf("Request failed: %v", err)
				continue
			}

			if resp.StatusCode >= 400 {
				b.Errorf("Request failed with status: %d", resp.StatusCode)
			}

			resp.Body.Close()
		}
	})
}
