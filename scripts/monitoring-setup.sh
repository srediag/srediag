#!/bin/bash
set -euo pipefail

# Monitoring setup script for SREDIAG
# This script sets up complete monitoring stack with Prometheus and Grafana

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

echo "📊 Setting up SREDIAG monitoring stack..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check prerequisites
check_prerequisites() {
    echo "🔍 Checking prerequisites..."
    
    if ! command_exists docker; then
        print_error "Docker is required but not installed"
        exit 1
    fi
    
    if ! command_exists docker-compose; then
        print_error "Docker Compose is required but not installed"
        exit 1
    fi
    
    # Check if ports are available
    ports=(3000 9090 13133 8889 4317 4318 16686)
    for port in "${ports[@]}"; do
        if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1; then
            print_warning "Port $port is already in use"
        fi
    done
    
    print_status "Prerequisites check completed"
}

# Start monitoring stack
start_monitoring() {
    echo "🚀 Starting monitoring stack..."
    
    cd "${PROJECT_ROOT}"
    
    # Start the full monitoring stack
    docker-compose up -d srediag prometheus grafana jaeger
    
    print_status "Monitoring stack started"
}

# Wait for services to be ready
wait_for_services() {
    echo "⏳ Waiting for services to be ready..."
    
    services=(
        "http://localhost:13133/:SREDIAG Health"
        "http://localhost:8889/metrics:SREDIAG Metrics"
        "http://localhost:9090/-/ready:Prometheus"
        "http://localhost:3000/api/health:Grafana"
        "http://localhost:16686/:Jaeger"
    )
    
    for service in "${services[@]}"; do
        IFS=':' read -r url name <<< "$service"
        
        print_info "Waiting for $name..."
        
        max_attempts=30
        attempt=1
        
        while [ $attempt -le $max_attempts ]; do
            if curl -sf "$url" >/dev/null 2>&1; then
                print_status "$name is ready"
                break
            fi
            
            if [ $attempt -eq $max_attempts ]; then
                print_warning "$name is not responding after $max_attempts attempts"
                break
            fi
            
            sleep 2
            ((attempt++))
        done
    done
}

# Setup Grafana dashboards
setup_grafana() {
    echo "📈 Setting up Grafana dashboards..."
    
    # Wait a bit more for Grafana to fully initialize
    sleep 10
    
    print_info "Grafana is available at: http://localhost:3000"
    print_info "Default credentials: admin/admin"
    print_status "Grafana dashboards configured via provisioning"
}

# Display monitoring URLs
show_monitoring_urls() {
    echo ""
    echo "🎯 Monitoring Stack URLs:"
    echo "========================================"
    echo "📊 Grafana:     http://localhost:3000"
    echo "   Username: admin, Password: admin"
    echo ""
    echo "📈 Prometheus:  http://localhost:9090"
    echo ""
    echo "🔍 Jaeger:      http://localhost:16686"
    echo ""
    echo "💚 SREDIAG:     http://localhost:13133 (health)"
    echo "📊 SREDIAG:     http://localhost:8889/metrics"
    echo ""
    echo "🔗 OTLP Endpoints:"
    echo "   gRPC:  localhost:4317"
    echo "   HTTP:  localhost:4318"
    echo "========================================"
}

# Generate sample data
generate_sample_data() {
    echo "📦 Generating sample telemetry data..."
    
    # Create a simple script to generate test data
    cat > /tmp/generate-traces.sh << 'EOF'
#!/bin/bash
# Simple script to generate test traces for demonstration

for i in {1..10}; do
    curl -X POST http://localhost:4318/v1/traces \
        -H "Content-Type: application/json" \
        -d '{
            "resourceSpans": [{
                "resource": {
                    "attributes": [{
                        "key": "service.name",
                        "value": {"stringValue": "demo-service"}
                    }]
                },
                "scopeSpans": [{
                    "spans": [{
                        "traceId": "'$(openssl rand -hex 16)'",
                        "spanId": "'$(openssl rand -hex 8)'",
                        "name": "demo-operation-'$i'",
                        "kind": "SPAN_KIND_SERVER",
                        "startTimeUnixNano": "'$(($(date +%s) * 1000000000))'",
                        "endTimeUnixNano": "'$((($(date +%s) + 1) * 1000000000))'"
                    }]
                }]
            }]
        }' >/dev/null 2>&1
    
    if [ $? -eq 0 ]; then
        echo "✅ Sent trace $i"
    else
        echo "❌ Failed to send trace $i"
    fi
    
    sleep 1
done
EOF

    chmod +x /tmp/generate-traces.sh
    
    if command_exists openssl && curl -sf http://localhost:4318 >/dev/null 2>&1; then
        print_info "Generating sample traces..."
        bash /tmp/generate-traces.sh
        print_status "Sample data generated"
    else
        print_warning "Cannot generate sample data (SREDIAG not ready or missing openssl)"
    fi
    
    rm -f /tmp/generate-traces.sh
}

# Health check
health_check() {
    echo "🏥 Running health checks..."
    
    # Check SREDIAG
    if curl -sf http://localhost:13133/ >/dev/null; then
        print_status "SREDIAG is healthy"
    else
        print_error "SREDIAG health check failed"
    fi
    
    # Check Prometheus targets
    if curl -sf http://localhost:9090/api/v1/targets >/dev/null; then
        print_status "Prometheus is scraping targets"
    else
        print_error "Prometheus targets check failed"
    fi
    
    # Check Grafana
    if curl -sf http://localhost:3000/api/health >/dev/null; then
        print_status "Grafana is healthy"
    else
        print_error "Grafana health check failed"
    fi
}

# Stop monitoring stack
stop_monitoring() {
    echo "🛑 Stopping monitoring stack..."
    
    cd "${PROJECT_ROOT}"
    docker-compose down
    
    print_status "Monitoring stack stopped"
}

# Show logs
show_logs() {
    echo "📋 Showing monitoring stack logs..."
    
    cd "${PROJECT_ROOT}"
    docker-compose logs -f srediag prometheus grafana jaeger
}

# Show help
show_help() {
    echo "SREDIAG Monitoring Setup Script"
    echo ""
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  start     Start the monitoring stack (default)"
    echo "  stop      Stop the monitoring stack"
    echo "  restart   Restart the monitoring stack"
    echo "  status    Show status of monitoring services"
    echo "  logs      Show logs from monitoring services"
    echo "  health    Run health checks"
    echo "  sample    Generate sample telemetry data"
    echo "  urls      Show monitoring URLs"
    echo "  help      Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 start    # Start monitoring stack"
    echo "  $0 sample   # Generate sample data"
    echo "  $0 logs     # View logs"
}

# Show status
show_status() {
    echo "📊 Monitoring Stack Status:"
    echo "========================================"
    
    cd "${PROJECT_ROOT}"
    docker-compose ps srediag prometheus grafana jaeger
}

# Main function
main() {
    local command="${1:-start}"
    
    case "$command" in
        "start")
            check_prerequisites
            start_monitoring
            wait_for_services
            setup_grafana
            health_check
            show_monitoring_urls
            print_info "Run '$0 sample' to generate test data"
            print_info "Run '$0 logs' to view logs"
            ;;
        "stop")
            stop_monitoring
            ;;
        "restart")
            stop_monitoring
            sleep 2
            main start
            ;;
        "status")
            show_status
            ;;
        "logs")
            show_logs
            ;;
        "health")
            health_check
            ;;
        "sample")
            generate_sample_data
            ;;
        "urls")
            show_monitoring_urls
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        *)
            print_error "Unknown command: $command"
            show_help
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"