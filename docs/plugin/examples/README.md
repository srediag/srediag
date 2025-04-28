# SREDIAG Plugin Examples

## Table of Contents

1. [Overview](#overview)
2. [OpenTelemetry Integration Examples](#opentelemetry-integration-examples)
3. [Diagnostic Receivers](#diagnostic-receivers)
4. [Diagnostic Processors](#diagnostic-processors)
5. [Diagnostic Exporters](#diagnostic-exporters)

## Overview

This directory contains example plugins demonstrating SREDIAG's integration with OpenTelemetry Collector and custom diagnostic capabilities. Each example includes complete source code, configuration, and documentation.

## OpenTelemetry Integration Examples

### 1. OpenTelemetry Receiver Extension

Extends an OpenTelemetry receiver with diagnostic capabilities:

```go
package mysysreceiver

import (
    "context"
    "go.opentelemetry.io/collector/component"
    "go.opentelemetry.io/collector/consumer"
    "go.opentelemetry.io/collector/receiver/hostmetricsreceiver"
    "github.com/srediag/srediag/pkg/plugin"
)

type SysReceiver struct {
    *hostmetricsreceiver.Receiver
    diagnostics plugin.DiagnosticCapabilities
}

func NewFactory() component.ReceiverFactory {
    return plugin.NewReceiverFactory(
        "mysysreceiver",
        createDefaultConfig,
        createReceiver,
    )
}

func (r *SysReceiver) Start(ctx context.Context, host component.Host) error {
    // Start OpenTelemetry receiver
    if err := r.Receiver.Start(ctx, host); err != nil {
        return err
    }
    
    // Initialize diagnostic capabilities
    return r.diagnostics.Start(ctx)
}
```

Configuration:

```yaml
receivers:
  mysysreceiver:
    collection_interval: 10s
    scrapers:
      cpu:
        metrics:
          system.cpu.utilization:
            enabled: true
      memory:
        metrics:
          system.memory.usage:
            enabled: true
    diagnostics:
      pattern_detection:
        enabled: true
        rules:
          - name: cpu_spike
            metric: system.cpu.utilization
            threshold: 90
            duration: 5m
```

## Diagnostic Receivers

### 1. System Diagnostic Receiver

Collects system-level diagnostic information:

```go
package sysdiag

import (
    "context"
    "go.opentelemetry.io/collector/pdata/pmetric"
    "github.com/srediag/srediag/pkg/plugin"
)

type SystemDiagnostics struct {
    config *Config
    consumer consumer.Metrics
    patterns plugin.PatternDetector
}

func (r *SystemDiagnostics) Collect(ctx context.Context) error {
    // Collect system metrics
    metrics := pmetric.NewMetrics()
    rMetrics := metrics.ResourceMetrics().AppendEmpty()
    
    // Add resource attributes
    rMetrics.Resource().Attributes().PutStr("service.name", "system")
    
    // Collect CPU diagnostics
    if err := r.collectCPUDiagnostics(ctx, rMetrics); err != nil {
        return err
    }
    
    // Collect Memory diagnostics
    if err := r.collectMemoryDiagnostics(ctx, rMetrics); err != nil {
        return err
    }
    
    // Analyze patterns
    if err := r.patterns.Analyze(ctx, metrics); err != nil {
        return err
    }
    
    return r.consumer.ConsumeMetrics(ctx, metrics)
}
```

Configuration:

```yaml
receivers:
  sysdiag:
    collection_interval: 30s
    diagnostics:
      cpu:
        enabled: true
        include_processes: true
        threshold_cpu_user: 80
        threshold_cpu_system: 20
      memory:
        enabled: true
        include_swap: true
        threshold_memory_used: 90
        threshold_swap_used: 50
    correlation:
      enabled: true
      window: 10m
```

## Diagnostic Processors

### 1. Pattern Analysis Processor

Analyzes metrics for diagnostic patterns:

```go
package patternprocessor

import (
    "context"
    "go.opentelemetry.io/collector/pdata/pmetric"
    "github.com/srediag/srediag/pkg/plugin"
)

type PatternProcessor struct {
    config    *Config
    patterns  []Pattern
    detector  plugin.PatternDetector
    next      consumer.Metrics
}

func (p *PatternProcessor) ProcessMetrics(ctx context.Context, md pmetric.Metrics) (pmetric.Metrics, error) {
    // Analyze patterns
    results, err := p.detector.AnalyzeMetrics(ctx, md)
    if err != nil {
        return md, err
    }
    
    // Add pattern detection results
    rMetrics := md.ResourceMetrics()
    for i := 0; i < rMetrics.Len(); i++ {
        rm := rMetrics.At(i)
        if patterns := results.PatternsForResource(rm.Resource()); len(patterns) > 0 {
            // Add pattern information as metrics
            sm := rm.ScopeMetrics().AppendEmpty()
            sm.Scope().SetName("pattern_processor")
            
            for _, pattern := range patterns {
                m := sm.Metrics().AppendEmpty()
                m.SetName("pattern.detected")
                dp := m.SetEmptyGauge().DataPoints().AppendEmpty()
                dp.Attributes().PutStr("pattern.name", pattern.Name)
                dp.Attributes().PutStr("pattern.severity", pattern.Severity)
                dp.SetDoubleValue(1.0)
            }
        }
    }
    
    return md, nil
}
```

Configuration:

```yaml
processors:
  pattern_processor:
    patterns:
      - name: resource_exhaustion
        conditions:
          - metric: system.cpu.utilization
            threshold: "> 90"
            duration: "5m"
          - metric: system.memory.usage
            threshold: "> 85"
            duration: "5m"
        severity: high
      - name: io_bottleneck
        conditions:
          - metric: system.disk.io_time
            threshold: "> 80"
            duration: "3m"
        severity: medium
    correlation:
      enabled: true
      window: 15m
      max_patterns: 10
```

## Diagnostic Exporters

### 1. Diagnostic Report Exporter

Exports diagnostic results to structured reports:

```go
package reportexporter

import (
    "context"
    "go.opentelemetry.io/collector/pdata/pmetric"
    "github.com/srediag/srediag/pkg/plugin"
)

type ReportExporter struct {
    config     *Config
    templates  map[string]*template.Template
    formatter  plugin.DiagnosticFormatter
}

func (e *ReportExporter) ConsumeMetrics(ctx context.Context, md pmetric.Metrics) error {
    // Extract diagnostic patterns
    patterns := e.extractPatterns(md)
    
    // Generate diagnostic report
    report, err := e.formatter.FormatDiagnostics(ctx, patterns)
    if err != nil {
        return err
    }
    
    // Export report based on configuration
    switch e.config.Format {
    case "pdf":
        return e.exportPDF(ctx, report)
    case "html":
        return e.exportHTML(ctx, report)
    default:
        return e.exportJSON(ctx, report)
    }
}
```

Configuration:

```yaml
exporters:
  diagnostic_report:
    format: pdf
    output_dir: /var/log/srediag/reports
    template: detailed
    include_metrics: true
    include_traces: true
    retention:
      enabled: true
      max_age: 30d
      max_size: 1GB
```

## See Also

- [OpenTelemetry Collector Examples](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/examples)
- [Plugin Development Guide](../development.md)
- [API Reference](../../reference/api.md)
- [Architecture Overview](../../architecture/overview.md)
