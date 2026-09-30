#!/bin/bash

ensure_brew() {
    if ! command -v brew &> /dev/null; then
        echo "⚠️ Homebrew not found."
        echo "📥 Installing Homebrew..."
        /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
        
        if [ -f "/opt/homebrew/bin/brew" ]; then
            eval "$(/opt/homebrew/bin/brew shellenv)"
        elif [ -f "/home/linuxbrew/.linuxbrew/bin/brew" ]; then
            eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"
        fi
    fi
}

install_imagemagick() {
    ensure_brew
    echo "📥 Installing ImageMagick via Homebrew..."
    brew install imagemagick
}

if ! command -v magick &> /dev/null && ! command -v convert &> /dev/null; then
    echo "ImageMagick not found."
    install_imagemagick
else
    echo "ImageMagick is already installed."
fi

if ! command -v rsrc &> /dev/null; then
    echo "rsrc not found."
    echo "📥 Installing rsrc..."
    go install github.com/akavel/rsrc@latest
    
    GOPATH_BIN="$(go env GOPATH)/bin"
    export PATH="$PATH:$GOPATH_BIN"
fi