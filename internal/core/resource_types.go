package core

// ResourceThresholds represents resource monitoring thresholds
type ResourceThresholds struct {
	CPUThreshold    float64 `json:"cpu_threshold"`
	MemoryThreshold float64 `json:"memory_threshold"`
	DiskThreshold   float64 `json:"disk_threshold"`
}

// ResourceUsage represents current resource usage
type ResourceUsage struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
	DiskUsage   float64 `json:"disk_usage"`
}

// MetricType represents the type of a metric
type MetricType string

const (
	// MetricTypeGauge represents a gauge metric
	MetricTypeGauge MetricType = "gauge"
	// MetricTypeCounter represents a counter metric
	MetricTypeCounter MetricType = "counter"
	// MetricTypeHistogram represents a histogram metric
	MetricTypeHistogram MetricType = "histogram"
)

// Metric represents a resource metric
type Metric struct {
	Name        string     `json:"name"`
	Value       float64    `json:"value"`
	Type        MetricType `json:"type"`
	Labels      []string   `json:"labels,omitempty"`
	Description string     `json:"description,omitempty"`
}
