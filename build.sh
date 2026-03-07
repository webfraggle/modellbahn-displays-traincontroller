#!/bin/bash
# Build script for mbd-tc
#
# Requirements:
#   macOS ARM64:  Go 1.21+, Xcode Command Line Tools
#   macOS Intel + Windows:  Docker + fyne-cross
#     go install github.com/fyne-io/fyne-cross@latest

LDFLAGS="-s -w"
OUTDIR="dist"

mkdir -p "$OUTDIR/config"
cp config/*.json "$OUTDIR/config/"

ok=0
skipped=0

# ── macOS ARM64 (native, no Docker needed) ───────────────────────────────────
echo "Building macOS ARM64 (Apple Silicon)..."
if go build -ldflags="$LDFLAGS" -o "$OUTDIR/mbd-tc-arm64" . 2>&1; then
    echo "  → $OUTDIR/mbd-tc-arm64"
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
    if fyne-cross darwin -arch amd64 -ldflags="$LDFLAGS" -output mbd-tc 2>&1; then
        cp fyne-cross/bin/darwin-amd64/mbd-tc "$OUTDIR/mbd-tc-x64"
        echo "  → $OUTDIR/mbd-tc-x64"
        ((ok++))
    else
        echo "  FAILED"
    fi
    rm -rf fyne-cross/

    # Windows: fyne-cross builds with GUI subsystem by default (no console popup)
    echo "Building Windows AMD64 via fyne-cross..."
    if fyne-cross windows -arch amd64 -ldflags="$LDFLAGS" -output mbd-tc 2>&1; then
        cp fyne-cross/bin/windows-amd64/mbd-tc.exe "$OUTDIR/mbd-tc.exe"
        echo "  → $OUTDIR/mbd-tc.exe"
        ((ok++))
    else
        echo "  FAILED"
    fi
    rm -rf fyne-cross/
fi

echo ""
echo "Done: $ok built, $skipped skipped. Binaries in $OUTDIR/"
