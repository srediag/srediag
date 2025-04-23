package diagnostic

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"

	"github.com/srediag/srediag/pkg/plugins"
)

const (
	pluginType    = "exporter/diagnostic"
	pluginName    = "diagnostic"
	pluginVersion = "v0.1.0"
)

// Config represents the exporter configuration
type Config struct {
	OutputDir     string        `mapstructure:"output_dir"`
	FlushInterval time.Duration `mapstructure:"flush_interval"`
	PrettyPrint   bool          `mapstructure:"pretty_print"`
}

// DiagnosticReport represents a diagnostic report
type DiagnosticReport struct {
	Timestamp time.Time         `json:"timestamp"`
	Resource  map[string]string `json:"resource"`
	Metrics   []MetricData      `json:"metrics"`
}

// MetricData represents metric data in the report
type MetricData struct {
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Unit        string                   `json:"unit"`
	Type        string                   `json:"type"`
	DataPoints  []map[string]interface{} `json:"dataPoints"`
}

// Exporter implements the diagnostic exporter
type Exporter struct {
	logger     *zap.Logger
	config     *Config
	host       component.Host
	reports    []*DiagnosticReport
	reportChan chan *DiagnosticReport
}

// NewFactory creates a new diagnostic exporter factory
func NewFactory() plugins.Factory {
	return &factory{}
}

type factory struct{}

func (f *factory) Type() string { return pluginType }

func (f *factory) CreatePlugin(cfg interface{}) (plugins.Plugin, error) {
	config, ok := cfg.(*Config)
	if !ok {
		return nil, fmt.Errorf("invalid configuration type")
	}

	return &Exporter{
		config:     config,
		reportChan: make(chan *DiagnosticReport, 100),
		reports:    make([]*DiagnosticReport, 0),
	}, nil
}

// Type returns the plugin type
func (e *Exporter) Type() string { return pluginType }

// Name returns the plugin name
func (e *Exporter) Name() string { return pluginName }

// Version returns the plugin version
func (e *Exporter) Version() string { return pluginVersion }

// Start initializes the exporter
func (e *Exporter) Start(ctx context.Context, host component.Host) error {
	e.host = host
	e.logger = zap.L()

	if err := os.MkdirAll(e.config.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	go e.flushReports(ctx)
	return nil
}

// Shutdown stops the exporter
func (e *Exporter) Shutdown(ctx context.Context) error {
	close(e.reportChan)
	return e.flush()
}

// ConsumeMetrics processes and exports the metrics data
func (e *Exporter) ConsumeMetrics(ctx context.Context, md pmetric.Metrics) error {
	report := e.createReport(md)
	select {
	case e.reportChan <- report:
	default:
		e.logger.Warn("report channel full, dropping report")
	}
	return nil
}

func (e *Exporter) createReport(md pmetric.Metrics) *DiagnosticReport {
	report := &DiagnosticReport{
		Timestamp: time.Now(),
		Resource:  make(map[string]string),
		Metrics:   make([]MetricData, 0),
	}

	rms := md.ResourceMetrics()
	for i := 0; i < rms.Len(); i++ {
		rm := rms.At(i)
		attrs := rm.Resource().Attributes()
		attrs.Range(func(k string, v pcommon.Value) bool {
			report.Resource[k] = v.AsString()
			return true
		})

		sms := rm.ScopeMetrics()
		for j := 0; j < sms.Len(); j++ {
			sm := sms.At(j)
			metrics := sm.Metrics()
			for k := 0; k < metrics.Len(); k++ {
				metric := metrics.At(k)
				metricData := MetricData{
					Name:        metric.Name(),
					Description: metric.Description(),
					Unit:        metric.Unit(),
					Type:        metric.Type().String(),
					DataPoints:  make([]map[string]interface{}, 0),
				}

				switch metric.Type() {
				case pmetric.MetricTypeGauge:
					dps := metric.Gauge().DataPoints()
					for l := 0; l < dps.Len(); l++ {
						dp := dps.At(l)
						metricData.DataPoints = append(metricData.DataPoints, map[string]interface{}{
							"timestamp": dp.Timestamp().AsTime(),
							"value":     dp.DoubleValue(),
						})
					}
				}

				report.Metrics = append(report.Metrics, metricData)
			}
		}
	}

	return report
}

func (e *Exporter) flushReports(ctx context.Context) {
	ticker := time.NewTicker(e.config.FlushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case report := <-e.reportChan:
			e.reports = append(e.reports, report)
		case <-ticker.C:
			if err := e.flush(); err != nil {
				e.logger.Error("failed to flush reports", zap.Error(err))
			}
		}
	}
}

func (e *Exporter) flush() error {
	if len(e.reports) == 0 {
		return nil
	}

	filename := filepath.Join(e.config.OutputDir,
		fmt.Sprintf("diagnostic-report-%d.json", time.Now().Unix()))

	var data []byte
	var err error
	if e.config.PrettyPrint {
		data, err = json.MarshalIndent(e.reports, "", "  ")
	} else {
		data, err = json.Marshal(e.reports)
	}
	if err != nil {
		return fmt.Errorf("failed to marshal reports: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write report file: %w", err)
	}

	e.reports = e.reports[:0]
	return nil
}
