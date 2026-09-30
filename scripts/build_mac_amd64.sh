#!/bin/bash

PROJECT_DIR="$(pwd)"
BUILD_DIR="./build/macOS/amd64"
DMG_STAGE_DIR="$BUILD_DIR/dmg_stage"
NAME_DMG="WordToExcel_mac_amd64.dmg"

ICONS_DIR="$PROJECT_DIR/icons"
SRC_IMG="$ICONS_DIR/icon_mac.png"
if [ -z "$SRC_IMG" ]; then
  exit 1
fi
ICONSET_PATH="$ICONS_DIR/AppIcon.iconset"

echo "⌛ Preparing the icon"
mkdir -p "$ICONSET_PATH" && \
sips -z 16 16     "$SRC_IMG" --out "$ICONSET_PATH/icon_16x16.png" && \
sips -z 32 32     "$SRC_IMG" --out "$ICONSET_PATH/icon_16x16@2x.png" && \
sips -z 32 32     "$SRC_IMG" --out "$ICONSET_PATH/icon_32x32.png" && \
sips -z 64 64     "$SRC_IMG" --out "$ICONSET_PATH/icon_32x32@2x.png" && \
sips -z 128 128   "$SRC_IMG" --out "$ICONSET_PATH/icon_128x128.png" && \
sips -z 256 256   "$SRC_IMG" --out "$ICONSET_PATH/icon_128x128@2x.png" && \
sips -z 256 256   "$SRC_IMG" --out "$ICONSET_PATH/icon_256x256.png" && \
sips -z 512 512   "$SRC_IMG" --out "$ICONSET_PATH/icon_256x256@2x.png" && \
sips -z 512 512   "$SRC_IMG" --out "$ICONSET_PATH/icon_512x512.png" && \
sips -z 1024 1024 "$SRC_IMG" --out "$ICONSET_PATH/icon_512x512@2x.png" && \
echo "✅ The icon has been successfully prepared."

iconutil -c icns "$ICONSET_PATH" && \
rm -rf "$ICONSET_PATH" && \

echo "🔨 Building application" && \
mkdir -p $BUILD_DIR && \
env CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -x -o $BUILD_DIR/WordToExcel ./cmd/app/main.go && \
mkdir -p $BUILD_DIR/WordToExcel.app/Contents/MacOS $BUILD_DIR/WordToExcel.app/Contents/Resources && \
cp "$BUILD_DIR/WordToExcel" "$BUILD_DIR/WordToExcel.app/Contents/MacOS/" && \
cp "$ICONS_DIR/AppIcon.icns" "$BUILD_DIR/WordToExcel.app/Contents/Resources/" && \
cp "$PROJECT_DIR/Info.plist" "$BUILD_DIR/WordToExcel.app/Contents/" && \
echo "✅ Build finished successfully"

echo "🔨 Building .dmg file"
mkdir -p "$DMG_STAGE_DIR" && \
cp -R "$BUILD_DIR/WordToExcel.app" "$DMG_STAGE_DIR/" && \
ln -s /Applications "$DMG_STAGE_DIR/Applications" && \
rm -f "$BUILD_DIR/$NAME_DMG" && \
hdiutil create -volname "WordToExcel Installer" \
               -srcfolder "$DMG_STAGE_DIR" \
               -ov -format UDZO \
               "$BUILD_DIR/$NAME_DMG" && \
rm -rf "$DMG_STAGE_DIR" && \
echo "✅ Build .dmg file finished successfully"
