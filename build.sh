#!/bin/bash

# Create plugins directory structure
mkdir -p bin/plugins/{receivers,processors,exporters,extensions}

# Build nopreceiver plugin
cd plugins/receiver/nop
go build -buildmode=plugin -o ../../../bin/plugins/receivers/nopreceiver.so

# Build other plugins as needed 
