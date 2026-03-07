# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**mbd-cli** is a CLI tool that bridges model railway software (e.g., TrainController) with "Zugzielanzeiger" (train destination display) hardware controllers via their REST API.

The Go implementation is the only active codebase on `main`. The original Node.js implementation has been removed and is accessible via the git tag `last-nodejs-version`.

## Go Implementation

### Setup

Requires Go 1.21+ and Xcode Command Line Tools (for CGO/Fyne):

```
brew install go          # if not installed
go mod download          # fetch dependencies
```

### Running in development

```
go run . --next
go run . --setTime "12:30"
go run . --setTrain1 "ICE123|12:30|Berlin|Hannover - Wolfsburg|0|Info"
go run .                 # opens config UI (Fyne window)
```

### Building

```
./build.sh
```

Produces binaries in `dist/`: `mbd-cli-arm64` (native), `mbd-cli-x64` and `mbd-cli.exe` via fyne-cross + Docker. Docker must be running for the latter two. App-ID: `de.modellbahn-displays.mbd-cli`.

**fyne-cross prerequisite:** `fyne-cross` requires the **old** fyne CLI (`fyne.io/fyne/v2/cmd/fyne`), not the new `fyne.io/tools/cmd/fyne`. The new CLI renamed internal flags and is not yet compatible with fyne-cross.

```
go install github.com/fyne-io/fyne-cross@latest
go install fyne.io/fyne/v2/cmd/fyne@latest   # ← required, ignore the deprecation notice
```

### Architecture

```
main.go                      Entry point: no args → UI, args → spawn+exit
internal/
  spawn/spawn_windows.go     Detach() — Windows: CREATE_NO_WINDOW
  spawn/spawn_unix.go        Detach() — macOS/Linux: Setsid
  config/config.go           Load/Save/List/Delete JSON configs
  api/client.go              HTTP client for all REST endpoints
  api/utf8.go                repairUTF8() — Windows-1252 double-encode fix
  ui/ui.go                   Fyne config editor window
config/
  default.json               Default display endpoint
  gleis1.json                Example second display
build.sh                     Cross-compile script
```

**Non-blocking execution pattern:**
The foreground process (called by TrainController) spawns itself with `--bg` appended, then exits immediately. The background copy performs the actual HTTP call. This means the calling software never waits.

**CLI flags** (identical to original Node.js version):
- `--next` / `--prev` — GET `/skipNext` or `/skipPrev`
- `--setTime "HH:MM"` — POST `/setTime`
- `--setTrain1|2|3 "nr|time|dest|via|delay|info"` — POST `/zug1|2|3`
- `--image <filename>` — GET `/showImage`
- `--gleis A|B` — track selection for dual-track displays (default: `GleisA`)
- `--conf <name>` — loads `config/<name>.json` instead of `config/default.json`
- `--timeout <ms>` — HTTP timeout in ms (default: 30000)

**Config files** live next to the executable in a `config/` subdirectory. During `go run`, falls back to `./config/`. Only required key: `endpoint` (base URL of the display controller).

**REST API payload** for train slots (`/zug1`, `/zug2`, `/zug3`):
```json
{ "nr": "RE50", "zeit": "12:30", "vonnach": "Berlin", "via": "Hannover",
  "abw": "0", "hinweis": "", "fusszeile": "", "abschnitte": "", "reihung": "", "path": "GleisA" }
```
Fields `fusszeile`, `abschnitte`, `reihung` are sent but have no effect on the display.

**`repairUTF8()`** (`internal/api/utf8.go`) — fixes double-encoded UTF-8 from TrainController on Windows (Windows-1252 mis-interpreted as UTF-8).

**Fyne UI** (`internal/ui/ui.go`) — opens when binary is called with no arguments. Left panel: config profiles (create/delete/edit endpoint/test connection). Right panel: CLI command builder that assembles a ready-to-run command string based on selected command, track (Gleis A/B), and dynamic argument fields. The command string is editable and can be copied or executed directly via the "Ausführen" button (calls `os.Executable()` externally so the spawn mechanism and `debug.log` work identically to TrainController). The command prefix is platform-aware: `./mbd-cli-arm64` or `./mbd-cli-x64` on macOS (derived from the running binary name), `.\mbd-cli.exe` on Windows.

## Original Node.js Implementation

Removed from `main`. Accessible via git tag `last-nodejs-version`. Entry point was `index.js`, built with `ncc` + `pkg`. Known issues: antivirus false positives on Windows, Gatekeeper issues on macOS, blocking caller while HTTP request runs.
