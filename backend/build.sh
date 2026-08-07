#!/bin/bash
set -e
echo "Building inspirate-backend..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/server ./cmd/api
echo "Binary: backend/bin/server"
