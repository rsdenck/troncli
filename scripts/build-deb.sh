#!/bin/bash
# scripts/build-deb.sh
# Build Debian package for TronCLI

set -e

VERSION="0.2.21"
ARCH="amd64"
PKG_DIR="packaging/deb"
BUILD_DIR="dist/deb/troncli_${VERSION}_${ARCH}"

echo "Building Debian package for TronCLI v${VERSION}..."

# Clean previous build
rm -rf "dist/deb"
mkdir -p "$BUILD_DIR/usr/local/bin"
mkdir -p "$BUILD_DIR/usr/local/share/troncli/man"
mkdir -p "$BUILD_DIR/DEBIAN"

# Compile binary
echo "Compiling binary..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o "$BUILD_DIR/usr/local/bin/troncli" cmd/troncli/main.go

# Copy man page
echo "Copying documentation..."
cp docs/man/troncli.1 "$BUILD_DIR/usr/local/share/troncli/man/"
gzip "$BUILD_DIR/usr/local/share/troncli/man/troncli.1"

# Copy control files
echo "Copying control files..."
cp "$PKG_DIR/DEBIAN/control" "$BUILD_DIR/DEBIAN/"
cp "$PKG_DIR/DEBIAN/postinst" "$BUILD_DIR/DEBIAN/"
cp "$PKG_DIR/DEBIAN/prerm" "$BUILD_DIR/DEBIAN/"
chmod 755 "$BUILD_DIR/DEBIAN/postinst"
chmod 755 "$BUILD_DIR/DEBIAN/prerm"

# Build package
echo "Creating .deb package..."
dpkg-deb --build "$BUILD_DIR" "dist/deb/troncli_${VERSION}_${ARCH}.deb"

echo "Debian package built successfully at dist/deb/troncli_${VERSION}_${ARCH}.deb"
