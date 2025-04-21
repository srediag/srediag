# srediag

**srediag** is a lightweight, plugin‑driven infrastructure diagnostics agent written in Go.  
It exports metrics, logs and traces via OpenTelemetry (OTLP) and stores them in ClickHouse.  
The agent core is extensible via HashiCorp Go‑Plugin.

---

## 🚀 Quickstart

```bash
# Clone and enter
git clone https://github.com/srediag/srediag.git
cd srediag

# Build the agent
make build

# Run with default config
./bin/srediag --config ./configs/config.yaml
```

---

## 📁 Repository Layout

```text
.
├── LICENSE
├── Makefile
├── Dockerfile
├── README.md
├── configs/
│   └── config.yaml
├── cmd/
│   └── srediag/
│       └── main.go
├── internal/
│   └── plugins/
│       └── plugins.go
├── metrics/
├── logs/
├── traces/
├── .golangci.yml
├── .gitignore
└── .github/
    └── workflows/
        └── ci.yml
```

---

## 🛠️ Features

- **Plugin architecture** powered by [hashicorp/go-plugin].  
- **OpenTelemetry** support (metrics, logs, traces) via OTLP gRPC.  
- **ClickHouse** backend for high‑performance analytics.  
- **Config** with [spf13/viper] + Cobra CLI.  
- **Structured logging** with [uber/zap].

---

## ⚙️ Configuration

Edit `configs/config.yaml`:

```yaml
otlp:
  endpoint: "localhost:4317"

plugins:
  - metrics
  - logs
  - traces
```

---

## 📦 Releases & Docker

Build and publish a Docker image:

```bash
make docker
```

Then run:

```bash
docker run --rm -v "$PWD/configs":/app/configs \
  srediag/srediag:latest --config /app/configs/config.yaml
```

---

## ✔️ Contributing

1. Fork the repo.  
2. Create `feature/xyz` branch.  
3. Implement and test.  
4. Send a Pull Request.  

Please read [`CODE_OF_CONDUCT.md`] and [`CONTRIBUTING.md`] for guidelines.

---

## 📜 License

This project is licensed under Apache 2.0. See [LICENSE] for details.

[hashicorp/go-plugin]: https://github.com/hashicorp/go-plugin
[spf13/viper]:    https://github.com/spf13/viper
[uber/zap]:       https://github.com/uber-go/zap
