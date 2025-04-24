#!/bin/bash

# Build script for example plugins

set -e

# Build the simple receiver plugin
echo "Building simple receiver plugin..."
cd simplereceiver
go build -buildmode=plugin -o ../../simplereceiver.so .
echo "Simple receiver plugin built successfully!"

# Add more plugin builds here as needed

echo "All example plugins built successfully!" 
