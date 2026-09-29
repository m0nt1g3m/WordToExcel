#!/bin/bash

OS_TYPE="$(uname -s)"

install_imagemagick() {
    case "$OS_TYPE" in
        Darwin*)
            if ! command -v brew &> /dev/null; then
                echo "Ошибка: Homebrew не установлен. Установите brew или ImageMagick вручную."
                exit 1
            fi
            echo "Установка ImageMagick через Homebrew..."
            brew install imagemagick
            ;;
        Linux*)
            if command -v apt-get &> /dev/null; then
                sudo apt-get update && sudo apt-get install -y imagemagick
            elif command -v dnf &> /dev/null; then
                sudo dnf install -y imagemagick
            elif command -v pacman &> /dev/null; then
                sudo pacman -S --noconfirm imagemagick
            else
                echo "Ошибка: Не удалось определить пакетный менеджер. Установите ImageMagick вручную."
                exit 1
            fi
            ;;
    esac
}

if ! command -v magick &> /dev/null && ! command -v convert &> /dev/null; then
    echo "ImageMagick не найден."
    echo "Установка imagemagick..."
    install_imagemagick
else
    echo "ImageMagick уже установлен."
fi

if ! command -v rsrc &> /dev/null; then
    echo "rsrc не найдена." 
    echo "Установка через rsrc..."
    go install github.com/akavel/rsrc@latest
    export PATH="$PATH:$(go env GOPATH)/bin"
fi