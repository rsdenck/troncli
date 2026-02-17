#!/bin/bash
# TRONCLI Installer Script
# Usage: curl -sS https://raw.githubusercontent.com/rsdenck/troncli/main/setup-tron.sh | bash

set -e

REPO="rsdenck/troncli"

# Fetch the latest release tag
LATEST_RELEASE=$(curl -sS https://api.github.com/repos/$REPO/releases/latest | grep "tag_name" | cut -d '"' -f 4)

# Fallback if no release found (e.g., first run)
if [ -z "$LATEST_RELEASE" ]; then
  echo "No release found, trying v0.1.1..."
  LATEST_RELEASE="v0.1.1"
fi

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
if [ "$OS" != "linux" ]; then
  echo "Error: Only Linux is supported."
  exit 1
fi

# Detect Architecture
ARCH=$(uname -m)
if [ "$ARCH" == "x86_64" ]; then
  ARCH="amd64"
elif [ "$ARCH" == "aarch64" ]; then
  ARCH="arm64"
else
  echo "Error: Unsupported architecture $ARCH."
  exit 1
fi

# Construct Download URL
# Version without 'v' prefix
VERSION_NO_V=${LATEST_RELEASE#v}
FILENAME="troncli_${VERSION_NO_V}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/$REPO/releases/download/$LATEST_RELEASE/$FILENAME"

echo "Installing TRONCLI $LATEST_RELEASE..."
echo "Downloading from $URL..."

# Create temp directory
TMP_DIR=$(mktemp -d)
trap "rm -rf $TMP_DIR" EXIT

if curl -sSL -o "$TMP_DIR/$FILENAME" "$URL"; then
    echo "Download complete."
else
    echo "Error: Failed to download release."
    exit 1
fi

echo "Extracting..."
tar -xzf "$TMP_DIR/$FILENAME" -C "$TMP_DIR"

echo "Installing to /usr/local/bin/troncli..."
if [ -w /usr/local/bin ]; then
    mv "$TMP_DIR/troncli" /usr/local/bin/troncli
else
    sudo mv "$TMP_DIR/troncli" /usr/local/bin/troncli
fi

chmod +x /usr/local/bin/troncli

echo "---------------------------------------------"
echo "TRONCLI installed successfully!"
echo "Run 'troncli --help' to get started."
echo "---------------------------------------------"
troncli --version
