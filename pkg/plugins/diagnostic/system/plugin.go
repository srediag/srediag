package system

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"

	"github.com/srediag/srediag/pkg/plugins"
)

const (
	pluginType    = "diagnostic/system"
	pluginName    = "system"
	pluginVersion = "v0.1.0"
)

// Config represents the plugin configuration
type Config struct {
	plugins.PluginConfig `mapstructure:",squash"`
	CollectInterval      time.Duration `mapstructure:"collect_interval"`
	CPUThreshold         float64       `mapstructure:"cpu_threshold"`
	MemThreshold         float64       `mapstructure:"mem_threshold"`
}

// Plugin implements the system diagnostic plugin
type Plugin struct {
	logger *zap.Logger
	config *Config
	host   component.Host
}

// NewFactory creates a new system plugin factory
func NewFactory() plugins.Factory {
	return &factory{}
}

type factory struct{}

// Type implements component.Factory
func (f *factory) Type() component.Type {
	return component.MustNewType("system")
}

// CreateDefaultConfig implements component.Factory
func (f *factory) CreateDefaultConfig() component.Config {
	return &Config{
		PluginConfig: plugins.PluginConfig{
			Enabled:  true,
			Settings: make(map[string]string),
		},
		CollectInterval: 30 * time.Second,
		CPUThreshold:    80.0,
		MemThreshold:    80.0,
	}
}

// CreatePlugin implements plugins.Factory
func (f *factory) CreatePlugin(cfg interface{}) (plugins.BasePlugin, error) {
	config, ok := cfg.(*Config)
	if !ok {
		return nil, fmt.Errorf("invalid configuration type")
	}

	return &Plugin{
		config: config,
	}, nil
}

// Type implements plugins.BasePlugin
func (p *Plugin) Type() component.Type {
	return component.MustNewType("system")
}

// Name implements plugins.BasePlugin
func (p *Plugin) Name() string { return pluginName }

// Version implements plugins.BasePlugin
func (p *Plugin) Version() string { return pluginVersion }

// Start implements component.Component
func (p *Plugin) Start(ctx context.Context, host component.Host) error {
	p.host = host
	p.logger = zap.L().Named("system-plugin")
	return nil
}

// Shutdown implements component.Component
func (p *Plugin) Shutdown(ctx context.Context) error {
	return nil
}

// Diagnose implements plugins.DiagnosticPlugin
func (p *Plugin) Diagnose(ctx context.Context, target string, options map[string]interface{}) (plugins.DiagnosticResult, error) {
	result := plugins.DiagnosticResult{
		Status:   "success",
		Message:  "System diagnostics completed",
		Details:  make(map[string]interface{}),
		Severity: "info",
	}

	// Collect CPU metrics
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return result, fmt.Errorf("failed to get CPU usage: %w", err)
	}

	// Collect memory metrics
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return result, fmt.Errorf("failed to get memory stats: %w", err)
	}

	// Add system information
	result.Details["os"] = runtime.GOOS
	result.Details["arch"] = runtime.GOARCH
	result.Details["cpu_count"] = runtime.NumCPU()
	result.Details["cpu_usage"] = cpuPercent[0]
	result.Details["memory_total"] = vmStat.Total
	result.Details["memory_used"] = vmStat.Used
	result.Details["memory_free"] = vmStat.Free
	result.Details["memory_usage"] = vmStat.UsedPercent

	// Check thresholds
	if cpuPercent[0] > p.config.CPUThreshold {
		result.Status = "warning"
		result.Severity = "warning"
		result.Message = fmt.Sprintf("CPU usage (%.2f%%) exceeds threshold (%.2f%%)",
			cpuPercent[0], p.config.CPUThreshold)
	}

	if vmStat.UsedPercent > p.config.MemThreshold {
		result.Status = "warning"
		result.Severity = "warning"
		result.Message = fmt.Sprintf("Memory usage (%.2f%%) exceeds threshold (%.2f%%)",
			vmStat.UsedPercent, p.config.MemThreshold)
	}

	return result, nil
}
