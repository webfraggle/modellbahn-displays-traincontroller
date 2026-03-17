# mbd-cli v2.0.0 — Release Notes

**Kompletter Rewrite in Go** — die Node.js/pkg-Implementierung wurde vollständig ersetzt.

## Was ist neu

- **Native Binary** — kein eingebetteter Node.js-Runtime mehr. Deutlich weniger Antivirus-Fehlalarme (Windows) und keine Gatekeeper-Probleme (macOS).
- **Nicht-blockierend** — mbd-cli kehrt sofort zurück und blockiert TrainController nicht mehr. Der HTTP-Request läuft im Hintergrund.
- **Grafische Konfigurationsoberfläche** — Aufruf ohne Argumente öffnet ein Fenster zum Bearbeiten von Endpunkten, Anlegen/Löschen von Profilen und zum Testen der Verbindung.
- **CLI-Builder** — die Oberfläche enthält einen Befehlsgenerator, der fertige Kommandos für TrainController zusammenstellt und direkt ausführen kann.
- **`--setAllTrains`** — neues Flag zum gleichzeitigen Setzen aller drei Zugslots.
- **`--image`** — neues Flag zum Anzeigen eines Bildes auf dem Display.
- **UTF-8-Reparatur** — behebt doppelt kodierte Umlaute aus TrainController auf Windows.
- **`debug.log`** — wird automatisch neben der ausführbaren Datei geschrieben.

## Umbenennung

Die Binaries heißen jetzt **`mbd-cli`** (bisher `mbd-tc`):

| Neu | Alt | Plattform |
|---|---|---|
| `mbd-cli-arm64` | `mbd-tc-arm64` | macOS (Apple Silicon) |
| `mbd-cli-x64` | `mbd-tc-x64` | macOS (Intel) |
| `mbd-cli.exe` | `mbd-tc.exe` | Windows (x64) |

Wer die alten Namen beibehalten möchte, kann die Dateien einfach umbenennen — die CLI-Flags und das Verhalten sind identisch zur v0.x-Reihe. **Bestehende TrainController-Makros müssen nur dann angepasst werden, wenn der Dateiname im Makro steht.**

## Hinweis für bestehende Nutzer

Abgesehen vom neuen Dateinamen sind keine Änderungen an TrainController-Makros notwendig. Alle bisherigen Flags (`--next`, `--prev`, `--setTime`, `--setTrain1`–`3`, `--gleis`, `--conf`, `--timeout`) funktionieren wie bisher.
