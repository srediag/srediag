# OpenTelemetry Integration

SREDIAG leverages OpenTelemetry for standardized observability and diagnostics capabilities.

## Architecture Overview

```ascii
+------------------------------------------------+
|                  SREDIAG Core                   |
|                                                 |
|    +---------------+      +-----------------+   |
|    |  OTel SDK    |<---->|  OTel Bridge     |   |
|    +---------------+      +-----------------+   |
|           ^                       ^             |
|           |                       |             |
|    +---------------+      +-----------------+   |
|    | OTel Pipeline |<---->| OTel Exporters  |   |
|    +---------------+      +-----------------+   |
+-------------------------------------------------+
                   |
         +---------+---------+
         |                   |
+----------------+    +-----------------+
| OTel Collector |    | Custom Backend  |
+----------------+    +-----------------+
```

## Pipeline Components

### 1. Data Collection

```ascii
+------------------+
| Data Collection  |
|                  |
| +-------------+  |    +--------------+
| | Metrics     |<----->| System Stats |
| +-------------+  |    +--------------+
|       ^          |
|       |          |    +--------------+
| +-------------+  |    | Application  |
| | Traces      |<----->| Flows        |
| +-------------+  |    +--------------+
|       ^          |
|       |          |    +--------------+
| +-------------+  |    | System       |
| | Logs        |<----->| Events       |
| +-------------+  |    +--------------+
+------------------+
```

### 2. Processing Pipeline

```ascii
+----------------------+
|  Processing Pipeline |
|                      |
| +-----------------+  |
| | Batch Processor |  |
| +-----------------+  |
|          ^           |
|          |           |
| +-----------------+  |
| | Filter          |  |
| +-----------------+  |
|          ^           |
|          |           |
| +-----------------+  |
| | Transform       |  |
| +-----------------+  |
+----------------------+
```

### 3. Export Pipeline

```ascii
+------------------------+
|    Export Pipeline     |
|                        |
| +------------------+   |
| | OTLP Exporter    |   |
| |                  |   |
| | - gRPC/HTTP      |   |
| | - TLS/mTLS       |   |
| | - Authentication |   |
| +------------------+   |
|          ^             |
|          |             |
| +------------------+   |
| | Custom Exporters |   |
| |                  |   |
| | - Prometheus     |   |
| | - Jaeger         |   |
| | - Zipkin         |   |
| +------------------+   |
+------------------------+
```

## Configuration Examples

### 1. Core Setup

```yaml
telemetry:
  opentelemetry:
    endpoint: "localhost:4317"
    protocol: "grpc"
    metrics:
      enabled: true
      interval: "10s"
    traces:
      enabled: true
      sampler:
        type: "parentbased_traceidratio"
        ratio: 0.1
    logs:
      enabled: true
      level: "info"
```

### 2. Exporters

```yaml
exporters:
  otlp:
    endpoint: "collector:4317"
    tls:
      enabled: true
      cert_file: "/etc/certs/client.crt"
      key_file: "/etc/certs/client.key"
  prometheus:
    endpoint: "0.0.0.0:9090"
    namespace: "srediag"
```

## Resource Detection

```ascii
+------------------------+
|   Resource Detection   |
|                        |
| +------------------+   |
| | Host Detection   |   |
| +------------------+   |
|          ^             |
|          |             |
| +------------------+   |
| | K8s Detection    |   |
| +------------------+   |
|          ^             |
|          |             |
| +------------------+   |
| | Cloud Detection  |   |
| +------------------+   |
+------------------------+
```

## Best Practices

### 1. Data Collection Best Practices

- Use appropriate sampling rates
- Configure batch sizes
- Set reasonable intervals
- Monitor resource usage

### 2. Security Best Practicess

- Enable TLS/mTLS
- Use authentication
- Implement authorization
- Protect sensitive data

### 3. Performance Best Practices

- Optimize batch processing
- Configure appropriate buffers
- Monitor pipeline health
- Handle backpressure

## Troubleshooting

```ascii
+------------------------+
|    Troubleshooting     |
|                        |
| 1. Connection Issues   |
|    - Check endpoints   |
|    - Verify TLS        |
|    - Test network      |
|                        |
| 2. Pipeline Issues     |
|    - Check configs     |
|    - Monitor queues    |
|    - Verify exports    |
|                        |
| 3. Performance Issues  |
|    - Check resources   |
|    - Monitor latency   |
|    - Adjust batching   |
+------------------------+
```

## Further Reading

- [OpenTelemetry Docs](https://opentelemetry.io/docs/)
- [Collector Guide](../configuration/collector.md)
- [Exporters Guide](../reference/exporters.md)
- [Monitoring Guide](../operations/monitoring.md)
