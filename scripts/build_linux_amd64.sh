#!/bin/bash

BUILD_DIR="./build/linux/amd64"

echo "🔨 Building application"
mkdir -p $BUILD_DIR && \
env CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o $BUILD_DIR/WordToExcel ./cmd/app/main.go && \
echo "✅ Build finished successfully"