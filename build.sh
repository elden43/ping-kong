#!/bin/bash

set -e

echo "Building Ping-Kong for multiple platforms..."

mkdir -p bin

# Linux
GOOS=linux GOARCH=amd64 go build -o bin/ping-kong-linux-amd64

# Windows
GOOS=windows GOARCH=amd64 go build -o bin/ping-kong-windows-amd64.exe

# macOS
GOOS=darwin GOARCH=amd64 go build -o bin/ping-kong-darwin-amd64

echo "✅ Build completed! Check the /bin directory for outputs."