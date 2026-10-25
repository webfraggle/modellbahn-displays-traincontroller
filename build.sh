#!/bin/bash
# Build script for mbd-cli
#
# Requirements:
#   macOS ARM64:  Go 1.26+, Xcode Command Line Tools
#   macOS Intel + Windows:  Docker + fyne-cross
#     go install github.com/fyne-io/fyne-cross@latest
#     go install fyne.io/fyne/v2/cmd/fyne@latest   ← fyne-cross requires the OLD fyne CLI
#     (fyne.io/tools/cmd/fyne is incompatible with fyne-cross until fyne-cross is updated)

# Ensure ~/go/bin is on PATH (needed when script is called directly, not from an interactive shell)
export PATH="$HOME/go/bin:$PATH"

LDFLAGS="-s -w"
# Note: do not wrap $LDFLAGS in extra quotes when passing to fyne-cross
# fyne-cross images may ship an older Go than go.mod requires (Windows image: 1.24).
# GOTOOLCHAIN=auto lets Go fetch the required toolchain inside the container.
# On the very first run the download notice confuses fyne's metadata step
# ("Failed to inject metadata") — the toolchain is cached afterwards, so just run again.
FYNE_CROSS_ENV="-env GOTOOLCHAIN=auto"
OUTDIR="dist"
# Version is maintained in VERSION (also embedded into the binary via go:embed)
VERSION="v$(tr -d '[:space:]' < VERSION)"
echo "mbd-cli $VERSION"

mkdir -p "$OUTDIR/config"
cp config/*.json "$OUTDIR/config/"
cp README.txt "$OUTDIR/README.txt"

ok=0
skipped=0

# ── macOS ARM64 (native, no Docker needed) ───────────────────────────────────
echo "Building macOS ARM64 (Apple Silicon)..."
if go build -ldflags "$LDFLAGS" -o "$OUTDIR/mbd-cli-arm64" . 2>&1; then
    echo "  → $OUTDIR/mbd-cli-arm64"
    ((ok++))
else
    echo "  FAILED"
fi

# ── macOS Intel + Windows via fyne-cross (requires Docker) ───────────────────
if ! command -v fyne-cross &>/dev/null; then
    echo ""
    echo "fyne-cross not found — skipping macOS Intel and Windows builds."
    echo "To enable:"
    echo "  1. Install Docker and start it"
    echo "  2. go install github.com/fyne-io/fyne-cross@latest"
    ((skipped+=2))
else
    echo "Building macOS Intel (AMD64) via fyne-cross..."
    if fyne-cross darwin $FYNE_CROSS_ENV -arch amd64 -app-id de.modellbahn-displays.mbd-cli -ldflags "$LDFLAGS" -output mbd-cli 2>&1; then
        cp fyne-cross/bin/darwin-amd64/mbd-cli "$OUTDIR/mbd-cli-x64"
        echo "  → $OUTDIR/mbd-cli-x64"
        ((ok++))
    else
        echo "  FAILED"
    fi
    rm -rf fyne-cross/

    # Windows: fyne-cross builds with GUI subsystem by default (no console popup)
    echo "Building Windows AMD64 via fyne-cross..."
    if fyne-cross windows $FYNE_CROSS_ENV -arch amd64 -app-id de.modellbahn-displays.mbd-cli -output mbd-cli 2>&1; then
        cp fyne-cross/bin/windows-amd64/mbd-cli.exe "$OUTDIR/mbd-cli.exe"
        echo "  → $OUTDIR/mbd-cli.exe"
        ((ok++))
    else
        echo "  FAILED"
    fi
    rm -rf fyne-cross/
fi

echo ""
echo "Done: $ok built, $skipped skipped. Binaries in $OUTDIR/"

# ── Release ZIP (only when all three binaries were built) ────────────────────
if [ "$ok" -eq 3 ]; then
    RELNAME="mbd-cli-$VERSION"
    rm -rf "$OUTDIR/$RELNAME" "$OUTDIR/$RELNAME.zip"
    mkdir -p "$OUTDIR/$RELNAME/config"
    cp "$OUTDIR/mbd-cli-arm64" "$OUTDIR/mbd-cli-x64" "$OUTDIR/mbd-cli.exe" "$OUTDIR/README.txt" "$OUTDIR/$RELNAME/"
    cp config/*.json "$OUTDIR/$RELNAME/config/"
    (cd "$OUTDIR" && zip -qrX "$RELNAME.zip" "$RELNAME" -x '*.DS_Store')
    echo "Release ZIP: $OUTDIR/$RELNAME.zip"
else
    echo "Release ZIP skipped (not all binaries built)."
fi
