#!/bin/bash
set -euo pipefail

APP_VERSION="${APP_VERSION:-0.1.0}"
export APP_VERSION

./scripts/build_linux_package.sh arm64
