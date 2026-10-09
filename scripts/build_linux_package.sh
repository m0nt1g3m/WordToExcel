#!/bin/bash
set -euo pipefail

APP_NAME="WordExcel"
APP_VERSION="${APP_VERSION:-0.1.0}"
APP_VERSION="${APP_VERSION#v}"
TARGET_OS="linux"
TARGET_ARCH="${1:-amd64}"

case "$TARGET_ARCH" in
  amd64|x86_64)
    GOARCH="amd64"
    PACKAGE_ARCH="amd64"
    ;;
  arm64|aarch64)
    GOARCH="arm64"
    PACKAGE_ARCH="arm64"
    ;;
  *)
    echo "Unsupported Linux architecture: $TARGET_ARCH" >&2
    exit 1
    ;;
esac

BUILD_DIR="./build/linux/${GOARCH}"
PACKAGE_ROOT="${BUILD_DIR}/package"
APP_BINARY_PATH="${BUILD_DIR}/${APP_NAME}"
OUTPUT_PACKAGE_PATH="${BUILD_DIR}/${APP_NAME}_${TARGET_OS}_${PACKAGE_ARCH}.deb"

mkdir -p "${PACKAGE_ROOT}/usr/local/bin" \
         "${PACKAGE_ROOT}/usr/share/applications" \
         "${PACKAGE_ROOT}/usr/share/icons/hicolor/256x256/apps" \
         "${PACKAGE_ROOT}/DEBIAN"

echo "🔨 Building ${APP_NAME} for Linux ${GOARCH}"
env CGO_ENABLED=1 GOOS=linux GOARCH="${GOARCH}" go build -o "${APP_BINARY_PATH}" ./cmd/app/main.go

cp "${APP_BINARY_PATH}" "${PACKAGE_ROOT}/usr/local/bin/wordtoexcel"
chmod 755 "${PACKAGE_ROOT}/usr/local/bin/wordtoexcel"

ICON_SOURCE=""
for candidate in \
  "./icons/icon.png" \
  "./internal/gui/assets/public/icon.png" \
  "./icons/icon_mac.png" \
  "./icons/icon_win.png"; do
  if [ -f "$candidate" ]; then
    ICON_SOURCE="$candidate"
    break
  fi
done

if [ -n "$ICON_SOURCE" ]; then
  cp "$ICON_SOURCE" "${PACKAGE_ROOT}/usr/share/icons/hicolor/256x256/apps/wordtoexcel.png"
fi

cat > "${PACKAGE_ROOT}/usr/share/applications/wordtoexcel.desktop" <<EOF
[Desktop Entry]
Name=WordToExcel
Comment=Convert Word tables to Excel
Exec=/usr/local/bin/wordtoexcel
Terminal=false
Type=Application
Icon=wordtoexcel
Categories=Office;
EOF

cat > "${PACKAGE_ROOT}/DEBIAN/control" <<EOF
Package: wordtoexcel
Version: ${APP_VERSION}
Section: utils
Priority: optional
Architecture: ${PACKAGE_ARCH}
Maintainer: WordToExcel <noreply@example.com>
Description: Convert Word tables to Excel
EOF

rm -f "${OUTPUT_PACKAGE_PATH}"
dpkg-deb --root-owner-group --build "${PACKAGE_ROOT}" "${OUTPUT_PACKAGE_PATH}"

echo "✅ Linux package created: ${OUTPUT_PACKAGE_PATH}"
