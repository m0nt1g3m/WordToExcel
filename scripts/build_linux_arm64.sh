#!/bin/bash

BUILD_DIR="./build/linux/arm64"

echo "Building application"
mkdir -p $BUILD_DIR && \
env CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -o $BUILD_DIR/WordToExcel ./cmd/app/main.go && \
echo "✅ Build finished successfully"