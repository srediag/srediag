#!/bin/bash
set -euo pipefail

# Test data generator script for SREDIAG
# Generates realistic telemetry data for testing and demonstration

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

# Configuration
OTLP_ENDPOINT="${OTLP_ENDPOINT:-http://localhost:4318}"
DURATION="${DURATION:-60}"
RATE="${RATE:-10}"
SERVICES="${SERVICES:-web-api,payment-service,user-service,notification-service}"

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
    if ! command_exists curl; then
        print_error "curl is required but not installed"
        exit 1
    fi
    
    if ! command_exists jq; then
        print_warning "jq is not installed, JSON formatting will be disabled"
    fi
    
    # Test OTLP endpoint
    if ! curl -sf "${OTLP_ENDPOINT}" >/dev/null 2>&1; then
        print_warning "OTLP endpoint ${OTLP_ENDPOINT} is not responding"
        print_info "Make sure SREDIAG is running with OTLP HTTP receiver enabled"
    fi
}

# Generate trace ID
generate_trace_id() {
    openssl rand -hex 16
}

# Generate span ID
generate_span_id() {
    openssl rand -hex 8
}

# Generate timestamp in nanoseconds
generate_timestamp() {
    echo $(($(date +%s) * 1000000000))
}

# Generate realistic service names
get_service_names() {
    echo "${SERVICES}" | tr ',' '\n'
}

# Generate realistic operation names
get_operations() {
    local service="$1"
    case "$service" in
        "web-api")
            echo "GET /users,POST /users,PUT /users/{id},DELETE /users/{id},GET /health"
            ;;
        "payment-service")
            echo "process_payment,validate_card,refund_payment,get_transaction_history"
            ;;
        "user-service")
            echo "authenticate_user,create_user,update_profile,get_user_data,delete_user"
            ;;
        "notification-service")
            echo "send_email,send_sms,send_push_notification,get_notification_preferences"
            ;;
        *)
            echo "operation_a,operation_b,operation_c,health_check"
            ;;
    esac
}

# Generate random HTTP status code
generate_status_code() {
    local codes=(200 200 200 200 201 202 400 401 403 404 500 502 503)
    echo "${codes[$((RANDOM % ${#codes[@]}))]}"
}

# Generate random duration in nanoseconds (1ms to 2s)
generate_duration() {
    echo $(((RANDOM % 2000 + 1) * 1000000))
}

# Generate a realistic trace
generate_trace() {
    local service="$1"
    local operation="$2"
    local trace_id="$3"
    local parent_span_id="$4"
    
    local span_id=$(generate_span_id)
    local start_time=$(generate_timestamp)
    local duration=$(generate_duration)
    local end_time=$((start_time + duration))
    local status_code=$(generate_status_code)
    
    # Determine span kind and status
    local span_kind="SPAN_KIND_SERVER"
    local status_code_value="STATUS_CODE_OK"
    
    if [[ "$status_code" -ge 400 ]]; then
        status_code_value="STATUS_CODE_ERROR"
    fi
    
    # Generate trace JSON
    cat << EOF
{
  "resourceSpans": [{
    "resource": {
      "attributes": [
        {"key": "service.name", "value": {"stringValue": "$service"}},
        {"key": "service.version", "value": {"stringValue": "1.0.0"}},
        {"key": "deployment.environment", "value": {"stringValue": "development"}},
        {"key": "host.name", "value": {"stringValue": "$(hostname)"}},
        {"key": "process.pid", "value": {"intValue": "$$"}}
      ]
    },
    "scopeSpans": [{
      "scope": {
        "name": "$service",
        "version": "1.0.0"
      },
      "spans": [{
        "traceId": "$trace_id",
        "spanId": "$span_id",
        $(if [ -n "$parent_span_id" ]; then echo "\"parentSpanId\": \"$parent_span_id\","; fi)
        "name": "$operation",
        "kind": "$span_kind",
        "startTimeUnixNano": "$start_time",
        "endTimeUnixNano": "$end_time",
        "attributes": [
          {"key": "http.method", "value": {"stringValue": "$(echo "$operation" | cut -d' ' -f1 | head -c 10)"}},
          {"key": "http.route", "value": {"stringValue": "$(echo "$operation" | cut -d' ' -f2- | head -c 50)"}},
          {"key": "http.status_code", "value": {"intValue": "$status_code"}},
          {"key": "http.response_size", "value": {"intValue": "$((RANDOM % 10000 + 100))"}},
          {"key": "user.id", "value": {"stringValue": "user-$((RANDOM % 1000))"}}
        ],
        "status": {
          "code": "$status_code_value"
        }
      }]
    }]
  }]
}
EOF
}

# Generate metrics data
generate_metrics() {
    local service="$1"
    local timestamp=$(generate_timestamp)
    
    cat << EOF
{
  "resourceMetrics": [{
    "resource": {
      "attributes": [
        {"key": "service.name", "value": {"stringValue": "$service"}},
        {"key": "service.version", "value": {"stringValue": "1.0.0"}}
      ]
    },
    "scopeMetrics": [{
      "scope": {
        "name": "$service-metrics",
        "version": "1.0.0"
      },
      "metrics": [
        {
          "name": "http_requests_total",
          "description": "Total number of HTTP requests",
          "unit": "1",
          "sum": {
            "dataPoints": [{
              "attributes": [
                {"key": "method", "value": {"stringValue": "GET"}},
                {"key": "status", "value": {"stringValue": "200"}}
              ],
              "timeUnixNano": "$timestamp",
              "asInt": "$((RANDOM % 1000 + 100))"
            }],
            "aggregationTemporality": "AGGREGATION_TEMPORALITY_CUMULATIVE",
            "isMonotonic": true
          }
        },
        {
          "name": "http_request_duration_seconds",
          "description": "HTTP request duration in seconds",
          "unit": "s",
          "histogram": {
            "dataPoints": [{
              "timeUnixNano": "$timestamp",
              "count": "$((RANDOM % 100 + 10))",
              "sum": "$((RANDOM % 1000)).$(printf "%03d" $((RANDOM % 1000)))",
              "bucketCounts": [$(for i in {1..10}; do echo -n "$((RANDOM % 10)), "; done | sed 's/, $//')],
              "explicitBounds": [0.1, 0.2, 0.5, 1.0, 2.0, 5.0, 10.0, 30.0, 60.0]
            }],
            "aggregationTemporality": "AGGREGATION_TEMPORALITY_CUMULATIVE"
          }
        }
      ]
    }]
  }]
}
EOF
}

# Generate logs data
generate_logs() {
    local service="$1"
    local timestamp=$(generate_timestamp)
    local log_levels=("INFO" "WARN" "ERROR" "DEBUG")
    local level="${log_levels[$((RANDOM % ${#log_levels[@]}))]}"
    
    cat << EOF
{
  "resourceLogs": [{
    "resource": {
      "attributes": [
        {"key": "service.name", "value": {"stringValue": "$service"}},
        {"key": "service.version", "value": {"stringValue": "1.0.0"}}
      ]
    },
    "scopeLogs": [{
      "scope": {
        "name": "$service-logs",
        "version": "1.0.0"
      },
      "logRecords": [{
        "timeUnixNano": "$timestamp",
        "severityText": "$level",
        "severityNumber": "$((RANDOM % 4 + 1))",
        "body": {"stringValue": "[$level] Processing request for user-$((RANDOM % 1000)) - operation completed in $((RANDOM % 1000))ms"},
        "attributes": [
          {"key": "component", "value": {"stringValue": "$service"}},
          {"key": "user.id", "value": {"stringValue": "user-$((RANDOM % 1000))"}}
        ]
      }]
    }]
  }]
}
EOF
}

# Send data to OTLP endpoint
send_data() {
    local data="$1"
    local endpoint="$2"
    
    if command_exists jq; then
        # Pretty print if jq is available
        data=$(echo "$data" | jq -c .)
    fi
    
    response=$(curl -s -w "%{http_code}" -X POST "$endpoint" \
        -H "Content-Type: application/json" \
        -d "$data" \
        -o /dev/null)
    
    if [[ "$response" == "200" || "$response" == "202" ]]; then
        return 0
    else
        return 1
    fi
}

# Generate and send test data
generate_test_data() {
    local duration="$1"
    local rate="$2"
    
    print_info "Generating test data for ${duration} seconds at ${rate} events/second"
    print_info "Target endpoint: ${OTLP_ENDPOINT}"
    
    local services_array=($(get_service_names))
    local total_events=0
    local successful_events=0
    local failed_events=0
    
    local end_time=$(($(date +%s) + duration))
    
    while [[ $(date +%s) -lt $end_time ]]; do
        for ((i=0; i<rate; i++)); do
            # Select random service
            local service="${services_array[$((RANDOM % ${#services_array[@]}))]}"
            local operations_str=$(get_operations "$service")
            local operations_array=($(echo "$operations_str" | tr ',' '\n'))
            local operation="${operations_array[$((RANDOM % ${#operations_array[@]}))]}"
            
            # Generate trace ID for this transaction
            local trace_id=$(generate_trace_id)
            
            # Generate and send trace
            local trace_data=$(generate_trace "$service" "$operation" "$trace_id" "")
            if send_data "$trace_data" "${OTLP_ENDPOINT}/v1/traces"; then
                ((successful_events++))
            else
                ((failed_events++))
            fi
            
            # Generate and send metrics (every 5th event)
            if [[ $((total_events % 5)) -eq 0 ]]; then
                local metrics_data=$(generate_metrics "$service")
                send_data "$metrics_data" "${OTLP_ENDPOINT}/v1/metrics" || true
            fi
            
            # Generate and send logs (every 3rd event)
            if [[ $((total_events % 3)) -eq 0 ]]; then
                local logs_data=$(generate_logs "$service")
                send_data "$logs_data" "${OTLP_ENDPOINT}/v1/logs" || true
            fi
            
            ((total_events++))
            
            # Progress indicator
            if [[ $((total_events % 50)) -eq 0 ]]; then
                echo -n "."
            fi
        done
        
        sleep 1
    done
    
    echo ""
    print_status "Data generation completed"
    print_info "Total events: $total_events"
    print_info "Successful: $successful_events"
    print_info "Failed: $failed_events"
    
    if [[ $failed_events -gt 0 ]]; then
        print_warning "Some events failed to send. Check SREDIAG connectivity."
    fi
}

# Show help
show_help() {
    echo "SREDIAG Test Data Generator"
    echo ""
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -e, --endpoint URL    OTLP endpoint URL (default: http://localhost:4318)"
    echo "  -d, --duration SEC    Duration to generate data in seconds (default: 60)"
    echo "  -r, --rate NUM        Events per second (default: 10)"
    echo "  -s, --services LIST   Comma-separated service names (default: web-api,payment-service,user-service,notification-service)"
    echo "  -h, --help            Show this help message"
    echo ""
    echo "Environment Variables:"
    echo "  OTLP_ENDPOINT         OTLP endpoint URL"
    echo "  DURATION              Duration in seconds"
    echo "  RATE                  Events per second"
    echo "  SERVICES              Comma-separated service names"
    echo ""
    echo "Examples:"
    echo "  $0                                          # Generate data with defaults"
    echo "  $0 -d 300 -r 20                           # Generate for 5 minutes at 20 events/sec"
    echo "  $0 -e http://srediag:4318 -s web,api,db   # Custom endpoint and services"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -e|--endpoint)
            OTLP_ENDPOINT="$2"
            shift 2
            ;;
        -d|--duration)
            DURATION="$2"
            shift 2
            ;;
        -r|--rate)
            RATE="$2"
            shift 2
            ;;
        -s|--services)
            SERVICES="$2"
            shift 2
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# Main execution
main() {
    echo "🎯 SREDIAG Test Data Generator"
    echo "=============================="
    
    check_prerequisites
    
    print_info "Configuration:"
    print_info "  Endpoint: $OTLP_ENDPOINT"
    print_info "  Duration: $DURATION seconds"
    print_info "  Rate: $RATE events/second"
    print_info "  Services: $SERVICES"
    
    if ! command_exists openssl; then
        print_error "openssl is required for generating trace and span IDs"
        exit 1
    fi
    
    generate_test_data "$DURATION" "$RATE"
    
    print_status "Test data generation completed successfully!"
}

# Run main function
main "$@"