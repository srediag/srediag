#!/bin/bash
set -euo pipefail

# Security check script for SREDIAG
# This script runs various security scans and checks

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

echo "🔒 Running security checks for SREDIAG..."

# Change to project root
cd "${PROJECT_ROOT}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${GREEN}✅ $1${NC}"
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

# Install required tools if not present
install_security_tools() {
    echo "📦 Installing security tools..."
    
    # Install gosec
    if ! command_exists gosec; then
        echo "Installing gosec..."
        go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
    fi
    
    # Install nancy
    if ! command_exists nancy; then
        echo "Installing nancy..."
        go install github.com/sonatypecommunity/nancy@latest
    fi
    
    # Install govulncheck
    if ! command_exists govulncheck; then
        echo "Installing govulncheck..."
        go install golang.org/x/vuln/cmd/govulncheck@latest
    fi
    
    # Install syft for SBOM generation
    if ! command_exists syft; then
        echo "Installing syft..."
        curl -sSfL https://raw.githubusercontent.com/anchore/syft/main/install.sh | sh -s -- -b "${HOME}/.local/bin"
        export PATH="${HOME}/.local/bin:${PATH}"
    fi
    
    # Install grype for vulnerability scanning
    if ! command_exists grype; then
        echo "Installing grype..."
        curl -sSfL https://raw.githubusercontent.com/anchore/grype/main/install.sh | sh -s -- -b "${HOME}/.local/bin"
        export PATH="${HOME}/.local/bin:${PATH}"
    fi
    
    print_status "Security tools installed"
}

# Run gosec static analysis
run_gosec() {
    echo "🔍 Running gosec static analysis..."
    
    if gosec -fmt json -out gosec-report.json -severity medium ./...; then
        print_status "Gosec scan completed"
    else
        print_warning "Gosec found security issues (see gosec-report.json)"
    fi
    
    # Generate human-readable report
    gosec -fmt text ./... > gosec-report.txt || true
    print_status "Gosec reports generated: gosec-report.json, gosec-report.txt"
}

# Run dependency vulnerability check with nancy
run_nancy() {
    echo "🔍 Running Nancy dependency vulnerability check..."
    
    if go list -json -deps ./... | nancy sleuth; then
        print_status "Nancy scan completed - no vulnerabilities found"
    else
        print_warning "Nancy found vulnerable dependencies"
        return 1
    fi
}

# Run govulncheck
run_govulncheck() {
    echo "🔍 Running govulncheck..."
    
    if govulncheck ./...; then
        print_status "Govulncheck completed - no vulnerabilities found"
    else
        print_warning "Govulncheck found vulnerabilities"
        return 1
    fi
}

# Generate SBOM (Software Bill of Materials)
generate_sbom() {
    echo "📋 Generating SBOM..."
    
    # Generate SPDX format
    syft packages . -o spdx-json=sbom.spdx.json
    
    # Generate CycloneDX format
    syft packages . -o cyclonedx-json=sbom.cyclonedx.json
    
    # Generate human-readable format
    syft packages . -o table=sbom.txt
    
    print_status "SBOM generated: sbom.spdx.json, sbom.cyclonedx.json, sbom.txt"
}

# Run vulnerability scan with grype
run_grype() {
    echo "🔍 Running Grype vulnerability scan..."
    
    if grype . -o json > grype-results.json; then
        print_status "Grype scan completed"
    else
        print_warning "Grype found vulnerabilities (see grype-results.json)"
    fi
    
    # Generate human-readable report
    grype . -o table > grype-results.txt || true
    print_status "Grype reports generated: grype-results.json, grype-results.txt"
}

# Check file permissions
check_file_permissions() {
    echo "🔍 Checking file permissions..."
    
    # Check for world-writable files
    if find . -type f -perm /002 -not -path "./vendor/*" -not -path "./.git/*"; then
        print_warning "Found world-writable files"
    else
        print_status "No world-writable files found"
    fi
    
    # Check for executable scripts
    echo "Executable files:"
    find . -type f -perm /111 -not -path "./vendor/*" -not -path "./.git/*" | head -10
}

# Check for secrets in code
check_secrets() {
    echo "🔍 Checking for potential secrets..."
    
    # Simple regex patterns for common secrets
    secret_patterns=(
        "password\s*=\s*['\"][^'\"]*['\"]"
        "api[_-]?key\s*=\s*['\"][^'\"]*['\"]"
        "secret\s*=\s*['\"][^'\"]*['\"]"
        "token\s*=\s*['\"][^'\"]*['\"]"
        "private[_-]?key"
    )
    
    found_secrets=false
    for pattern in "${secret_patterns[@]}"; do
        if grep -rE -i "${pattern}" . --exclude-dir=vendor --exclude-dir=.git --exclude="*.md" 2>/dev/null; then
            found_secrets=true
        fi
    done
    
    if [ "$found_secrets" = true ]; then
        print_warning "Potential secrets found in code"
    else
        print_status "No obvious secrets found"
    fi
}

# Generate security report
generate_security_report() {
    echo "📊 Generating security report..."
    
    report_file="security-report.md"
    
    cat > "${report_file}" << EOF
# SREDIAG Security Report

Generated on: $(date)

## Summary

This report contains the results of various security scans and checks performed on the SREDIAG project.

## Gosec Static Analysis

\`\`\`
$(cat gosec-report.txt 2>/dev/null || echo "Gosec report not available")
\`\`\`

## Dependency Vulnerabilities

### Nancy Results
\`\`\`
Nancy scan completed $(date)
\`\`\`

### Govulncheck Results
\`\`\`
Govulncheck scan completed $(date)
\`\`\`

## SBOM (Software Bill of Materials)

SBOM files generated:
- sbom.spdx.json (SPDX format)
- sbom.cyclonedx.json (CycloneDX format)
- sbom.txt (Human-readable)

## Vulnerability Scan (Grype)

\`\`\`
$(cat grype-results.txt 2>/dev/null || echo "Grype report not available")
\`\`\`

## File Permissions Check

All file permissions have been checked for security issues.

## Secrets Scan

Code has been scanned for potential hardcoded secrets.

## Recommendations

1. Regularly update dependencies to patch known vulnerabilities
2. Run security scans in CI/CD pipeline
3. Use proper secret management for sensitive data
4. Follow security best practices for Go development
5. Keep SBOM updated for supply chain security

EOF

    print_status "Security report generated: ${report_file}"
}

# Cleanup function
cleanup() {
    echo "🧹 Cleaning up temporary files..."
    # Remove temporary files if needed
}

# Main execution
main() {
    echo "🚀 Starting SREDIAG security checks..."
    
    # Trap to ensure cleanup on exit
    trap cleanup EXIT
    
    # Install tools
    install_security_tools
    
    # Run security checks
    echo ""
    echo "=========================="
    echo "Running Security Checks"
    echo "=========================="
    
    run_gosec || true
    echo ""
    
    run_nancy || true
    echo ""
    
    run_govulncheck || true
    echo ""
    
    generate_sbom
    echo ""
    
    run_grype || true
    echo ""
    
    check_file_permissions
    echo ""
    
    check_secrets
    echo ""
    
    generate_security_report
    
    echo ""
    echo "=========================="
    print_status "Security checks completed!"
    echo "=========================="
    
    echo ""
    echo "📁 Generated files:"
    echo "   - gosec-report.json (detailed security issues)"
    echo "   - gosec-report.txt (human-readable security issues)"
    echo "   - sbom.spdx.json (SPDX format SBOM)"
    echo "   - sbom.cyclonedx.json (CycloneDX format SBOM)"
    echo "   - sbom.txt (human-readable SBOM)"
    echo "   - grype-results.json (detailed vulnerability scan)"
    echo "   - grype-results.txt (human-readable vulnerability scan)"
    echo "   - security-report.md (comprehensive security report)"
}

# Run main function
main "$@"