#!/bin/bash

APP_NAME="WordToExcel"
APP_VERSION="${APP_VERSION:-0.1.0}"
TARGET_OS="macos"
TARGET_ARCH="amd64"
PROJECT_DIR="$(pwd)"
BUILD_DIR="./build/macOS/amd64"
DMG_STAGE_DIR="$BUILD_DIR/dmg_stage"
OUTPUT_DMG_PATH="$BUILD_DIR/${APP_NAME}_${TARGET_OS}_${TARGET_ARCH}.dmg"

ICON_DIR="$PROJECT_DIR/icons"
ICON_SOURCE="$ICON_DIR/icon_mac.png"
if [ -z "$ICON_SOURCE" ]; then
  exit 1
fi
ICONSET_DIR="$ICON_DIR/AppIcon.iconset"

echo "⌛ Preparing the icon"
mkdir -p "$ICONSET_DIR" && \
sips -z 16 16     "$ICON_SOURCE" --out "$ICONSET_DIR/icon_16x16.png" && \
sips -z 32 32     "$ICON_SOURCE" --out "$ICONSET_DIR/icon_16x16@2x.png" && \
sips -z 32 32     "$ICON_SOURCE" --out "$ICONSET_DIR/icon_32x32.png" && \
sips -z 64 64     "$ICON_SOURCE" --out "$ICONSET_DIR/icon_32x32@2x.png" && \
sips -z 128 128   "$ICON_SOURCE" --out "$ICONSET_DIR/icon_128x128.png" && \
sips -z 256 256   "$ICON_SOURCE" --out "$ICONSET_DIR/icon_128x128@2x.png" && \
sips -z 256 256   "$ICON_SOURCE" --out "$ICONSET_DIR/icon_256x256.png" && \
sips -z 512 512   "$ICON_SOURCE" --out "$ICONSET_DIR/icon_256x256@2x.png" && \
sips -z 512 512   "$ICON_SOURCE" --out "$ICONSET_DIR/icon_512x512.png" && \
sips -z 1024 1024 "$ICON_SOURCE" --out "$ICONSET_DIR/icon_512x512@2x.png" && \
echo "✅ The icon has been successfully prepared."

iconutil -c icns "$ICONSET_DIR" && \
rm -rf "$ICONSET_DIR" && \

echo "🔨 Building application" && \
mkdir -p "$BUILD_DIR" && \
rm -rf "$BUILD_DIR/${APP_NAME}.app" "$DMG_STAGE_DIR" && \
env CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -x -o "$BUILD_DIR/$APP_NAME" ./cmd/app/main.go && \
mkdir -p "$BUILD_DIR/${APP_NAME}.app/Contents/MacOS" "$BUILD_DIR/${APP_NAME}.app/Contents/Resources" && \
cp "$BUILD_DIR/$APP_NAME" "$BUILD_DIR/${APP_NAME}.app/Contents/MacOS/" && \
cp "$ICON_DIR/AppIcon.icns" "$BUILD_DIR/${APP_NAME}.app/Contents/Resources/" && \
cp "$PROJECT_DIR/Info.plist" "$BUILD_DIR/${APP_NAME}.app/Contents/" && \
echo "✅ Build finished successfully"

echo "🔨 Building .dmg file"
mkdir -p "$DMG_STAGE_DIR" && \
cp -R "$BUILD_DIR/${APP_NAME}.app" "$DMG_STAGE_DIR/" && \
ln -s /Applications "$DMG_STAGE_DIR/Applications" && \
rm -f "$OUTPUT_DMG_PATH" && \
hdiutil create -volname "WordToExcel Installer" \
               -srcfolder "$DMG_STAGE_DIR" \
               -ov -format UDZO \
               "$OUTPUT_DMG_PATH" && \
rm -rf "$DMG_STAGE_DIR" && \
echo "✅ Build .dmg file finished successfully"