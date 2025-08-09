#!/bin/bash

echo "Setting up development environment for Serial Plotter Android build..."

# Install Android NDK via Homebrew if not already installed
if ! brew list | grep -q android-ndk; then
    echo "Installing Android NDK..."
    brew install --cask android-ndk
else
    echo "Android NDK already installed"
fi

# Get Go bin directory
GO_BIN_DIR=$(go env GOPATH)/bin

# Set up environment variables
echo "Setting up environment variables..."

# Add to shell profile (zsh for macOS)
SHELL_PROFILE="$HOME/.zshrc"

# Android NDK environment
if ! grep -q "ANDROID_NDK_HOME" "$SHELL_PROFILE"; then
    echo 'export ANDROID_NDK_HOME="/opt/homebrew/share/android-ndk"' >> "$SHELL_PROFILE"
    echo "Added ANDROID_NDK_HOME to $SHELL_PROFILE"
else
    echo "ANDROID_NDK_HOME already configured"
fi

# Add Go bin to PATH if not already there
if ! grep -q "$GO_BIN_DIR" "$SHELL_PROFILE"; then
    echo "export PATH=\$PATH:$GO_BIN_DIR" >> "$SHELL_PROFILE"
    echo "Added Go bin directory to PATH in $SHELL_PROFILE"
else
    echo "Go bin directory already in PATH"
fi

# Install Fyne CLI tool
echo "Installing Fyne CLI tool..."
go install fyne.io/tools/cmd/fyne@latest

echo ""
echo "Setup complete! Please restart your terminal or run:"
echo "source $SHELL_PROFILE"
echo ""
echo "Then you can build for Android with:"
echo "./build-android.sh" 