#!/bin/bash
# Build script for mbd-tc
# Requires: Go 1.21+, Xcode Command Line Tools (for CGO/Fyne on macOS)
# For Windows cross-compilation: brew install mingw-w64

set -e

LDFLAGS="-s -w"
OUTDIR="dist"

mkdir -p "$OUTDIR/config"
cp config/*.json "$OUTDIR/config/"

echo "Building macOS ARM64 (Apple Silicon)..."
GOOS=darwin GOARCH=arm64 go build -ldflags="$LDFLAGS" -o "$OUTDIR/mbd-tc-arm64" .

echo "Building macOS AMD64 (Intel)..."
GOOS=darwin GOARCH=amd64 go build -ldflags="$LDFLAGS" -o "$OUTDIR/mbd-tc-x64" .

# Windows cross-compilation with Fyne requires mingw-w64 for CGO.
# Install with: brew install mingw-w64
# Then uncomment:
#
# echo "Building Windows AMD64..."
# CGO_ENABLED=1 \
#   CC=x86_64-w64-mingw32-gcc \
#   GOOS=windows GOARCH=amd64 \
#   go build -ldflags="$LDFLAGS -H windowsgui" -o "$OUTDIR/mbd-tc.exe" .
#
# Note: -H windowsgui suppresses the console window popup on Windows.
# For the CLI to be fully silent, this flag is required.
# For Windows builds, prefer building natively on a Windows machine
# or via GitHub Actions to avoid CGO cross-compilation complexity.

echo "Done. Binaries in $OUTDIR/"
