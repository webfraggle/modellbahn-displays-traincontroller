# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**mbd-tc** is a CLI tool that bridges model railway software (e.g., TrainController) with "Zugzielanzeiger" (train destination display) hardware controllers via their REST API.

The active implementation is the **Go rewrite** on branch `go-rewrite`. The original Node.js implementation (`index.js`) remains on `main` for reference.

## Go Implementation (active, branch: go-rewrite)

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

Produces binaries in `dist/`. For Windows, see the commented section in `build.sh` (requires `brew install mingw-w64` for cross-compilation, or build natively on Windows).

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

**Fyne UI** (`internal/ui/ui.go`) — opens when binary is called with no arguments. Shows list of config profiles, allows editing endpoint URL, creating/deleting profiles, and testing the connection.

## Original Node.js Implementation (branch: main, kept for reference)

Entry point: `index.js`. Built with `ncc` + `pkg`. Known issues: antivirus false positives on Windows, Gatekeeper issues on macOS, blocking caller while HTTP request runs.
