#!/bin/bash
# Script to fix invalid Go version references in the project

set -e

echo "Fixing Go version references from 1.24 to 1.21..."

# Files to update
files=(
    "Dockerfile"
    "Dockerfile.dev"
    ".github/workflows/ci.yml"
    ".github/workflows/release.yml"
    ".github/workflows/security.yml"
)

# Backup and update each file
for file in "${files[@]}"; do
    if [ -f "$file" ]; then
        echo "Updating $file..."
        # Create backup
        cp "$file" "$file.bak"
        # Replace Go version
        sed -i 's/golang:1\.24/golang:1.21/g' "$file"
        sed -i 's/go-version: "1\.24"/go-version: "1.21"/g' "$file"
        sed -i 's/GO_VERSION: "1\.24"/GO_VERSION: "1.21"/g' "$file"
        sed -i 's/"1\.24"/"1.21"/g' "$file"
        echo "✓ Updated $file"
    else
        echo "⚠ File not found: $file"
    fi
done

echo ""
echo "Go version references have been updated to 1.21"
echo "Backup files created with .bak extension"
echo ""
echo "Please review the changes and commit them:"
echo "  git diff Dockerfile* .github/workflows/*.yml"
echo "  git add Dockerfile* .github/workflows/*.yml"
echo "  git commit -m 'fix: Update Go version from invalid 1.24 to 1.21'"