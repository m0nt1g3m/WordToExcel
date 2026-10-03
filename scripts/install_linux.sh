#!/bin/bash
set -euo pipefail

if [ "$#" -ne 1 ]; then
  echo "Usage: $0 path/to/WordToExcel_<version>_<arch>.deb"
  exit 1
fi

DEB_PATH="$1"

if [ ! -f "$DEB_PATH" ]; then
  echo "Package not found: $DEB_PATH" >&2
  exit 1
fi

sudo dpkg -i "$DEB_PATH"

echo "✅ WordToExcel installed successfully."
echo "Run it with: wordtoexcel"
