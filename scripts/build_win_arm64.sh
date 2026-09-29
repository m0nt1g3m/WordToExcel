#!/bin/bash

PROJECT_DIR="$(pwd)"
BUILD_DIR="$PROJECT_DIR/build/win/arm64"
APP_DIR="$PROJECT_DIR/cmd/app"
MAGICK_CMD="magick"

if command -v magick &> /dev/null; then
    echo "Check magic commands"
    MAGICK_CMD="magick"
elif command -v magick convert &> /dev/null; then
    MAGICK_CMD="magick convert"
fi

echo "Building application"
mkdir -p "$BUILD_DIR" && \
rm $APP_DIR/*.syso && \
$MAGICK_CMD ./icons/icon.png -define icon:auto-resize=256,128,64,48,32,16 ./icons/icon.ico && \
rsrc -arch arm64 -ico "./icons/icon.ico" -o "$APP_DIR/rsrc_windows_arm64.syso" && \
cd $APP_DIR && \
env CGO_ENABLED=1  GOOS=windows GOARCH=arm64 go build -x -ldflags="-H windowsgui" -o "$BUILD_DIR/WordToExcel.exe" . && \
echo "✅ Build finished successfully"