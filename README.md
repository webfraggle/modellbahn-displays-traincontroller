# mbd-cli — Modellbahn Displays CLI

Kommandozeilentool zur Steuerung von Zugzielanzeiger-Displays über deren REST API.
Entwickelt für den Einsatz aus Modellbahn-Software wie TrainController, kehrt sofort zurück ohne den Aufrufer zu blockieren.

*(English version below)*

---

## Konfiguration

Beim ersten Start ohne Argumente öffnet sich die Konfigurations-Oberfläche:

```
./mbd-cli-arm64        # macOS Apple Silicon
./mbd-cli-x64          # macOS Intel
.\mbd-cli.exe          # Windows
```

Dort können Konfigurationsprofile angelegt, die Endpoint-URL eingetragen und die Verbindung getestet werden. Die Konfigurationen werden als JSON-Dateien im `config/`-Ordner neben der ausführbaren Datei gespeichert:

```json
{ "endpoint": "http://192.168.178.155" }
```

Der Befehlsgenerator in der rechten Hälfte der Oberfläche hilft beim Zusammenstellen der korrekten Befehle zum Kopieren oder direkten Ausführen.

---

## Optionen

### Nächsten / vorherigen Zug anzeigen

```
mbd-cli.exe --next
mbd-cli.exe --prev
```

### Uhrzeit setzen

```
mbd-cli.exe --setTime "12:30"
```

### Zuginformationen direkt setzen

Setzt bis zu drei Zugslots gleichzeitig. Der Wert ist eine durch `|` getrennte Zeichenfolge:
`TrainID|Zeit|Ziel|Via|Verspätung|Sonderinfo`

```
mbd-cli.exe --setTrain1 "ICE123|12:30|Berlin|Hannover - Wolfsburg|0|"
mbd-cli.exe --setTrain1 "ICE123|12:30|Berlin|Hannover|0|Info" --setTrain2 "RE50|21:12|Bebra|Hünfeld|+10|"
```

Leeres Display (Doppelspuranzeige):
```
mbd-cli.exe --gleis B --setTrain1 "|||||"
```

### Bild anzeigen

```
mbd-cli.exe --image 00Logo.png
```

Zeigt ein auf den Controller geladenes Bild (in den Modi Manuell, Intervall, Bilder). Das Bild bleibt bis zum nächsten Zug oder Bild aktiv.

### Mehrere Displays

Lege für jedes Display eine eigene Konfigurationsdatei an (z. B. `gleis1.json`) und lade sie mit `--conf`:

```
mbd-cli.exe --conf gleis1 --setTime "12:30"
```

### Doppelspuranzeige

```
mbd-cli.exe --gleis B --setTime "12:30"
```

Standard ist Gleis A.

### Weitere Optionen

| Option | Beschreibung |
|---|---|
| `--timeout <ms>` | HTTP-Timeout in Millisekunden (Standard: 30000) |
| `--conf <name>` | Lädt `config/<name>.json` statt `config/default.json` |

---

## macOS: Gatekeeper-Warnung beim Download

Wird die Binary aus dem Internet heruntergeladen und entpackt, blockiert macOS die Ausführung mit der Meldung „kann nicht überprüft werden". Das liegt am Quarantäne-Attribut, das beim Download gesetzt wird.

**Einmalig im Terminal nach dem Entpacken:**

```bash
xattr -d com.apple.quarantine mbd-cli-arm64
```

**Oder:** Im Finder Rechtsklick auf die Datei → **Öffnen** → **Trotzdem öffnen**. Danach auch in Systemeinstellungen → Datenschutz & Sicherheit sichtbar.

---

## Debugging

Das Tool schreibt automatisch eine `debug.log` Datei neben der ausführbaren Datei. Dort sind alle aufgerufenen Argumente und eventuelle Fehler protokolliert.

---

## Build

```
./build.sh
```

Erzeugt drei Binaries in `dist/`:

| Datei | Ziel |
|---|---|
| `mbd-cli-arm64` | macOS Apple Silicon (nativ) |
| `mbd-cli-x64` | macOS Intel (via fyne-cross + Docker) |
| `mbd-cli.exe` | Windows AMD64 (via fyne-cross + Docker) |

Voraussetzungen für macOS Intel- und Windows-Build:
- Docker (muss laufen)
- `go install github.com/fyne-io/fyne-cross@latest`
- `go install fyne.io/fyne/v2/cmd/fyne@latest` *(altes fyne CLI — Deprecation-Hinweis ignorieren, fyne-cross benötigt diese Version)*

---

---

# English Version

Command-line tool to control Zugzielanzeiger train destination displays via their REST API.
Designed for use from model railway software like TrainController — returns immediately without blocking the caller.

## Configuration

When launched without arguments, the configuration UI opens:

```
./mbd-cli-arm64        # macOS Apple Silicon
./mbd-cli-x64          # macOS Intel
.\mbd-cli.exe          # Windows
```

Create profiles, enter the display's endpoint URL, and test the connection. Configs are stored as JSON files in a `config/` folder next to the executable:

```json
{ "endpoint": "http://192.168.178.155" }
```

The command builder on the right side of the UI helps assemble the correct commands for copying or direct execution.

---

## Options

### Skip to next / previous train

```
mbd-cli.exe --next
mbd-cli.exe --prev
```

### Set time

```
mbd-cli.exe --setTime "12:30"
```

### Set train information directly

Sets up to three train slots at once. The value is a `|`-separated string:
`TrainID|Time|Destination|Via|Delay|SpecialInfo`

```
mbd-cli.exe --setTrain1 "ICE123|12:30|Berlin|Hannover - Wolfsburg|0|"
mbd-cli.exe --setTrain1 "ICE123|12:30|Berlin|Hannover|0|Info" --setTrain2 "RE50|21:12|Bebra|Hünfeld|+10|"
```

Clear a slot on a double-track display:
```
mbd-cli.exe --gleis B --setTrain1 "|||||"
```

### Show an image

```
mbd-cli.exe --image 00Logo.png
```

Displays an image previously uploaded to the controller (in Manual, Interval, or Images mode). The image stays on screen until the next train or image command.

### Multiple displays

Create a config file for each display (e.g. `gleis1.json`) and load it with `--conf`:

```
mbd-cli.exe --conf gleis1 --setTime "12:30"
```

### Double-track display

```
mbd-cli.exe --gleis B --setTime "12:30"
```

Default is track A.

### Further options

| Option | Description |
|---|---|
| `--timeout <ms>` | HTTP timeout in milliseconds (default: 30000) |
| `--conf <name>` | Loads `config/<name>.json` instead of `config/default.json` |

---

## macOS: Gatekeeper warning after download

When the binary is downloaded from the internet and unzipped, macOS blocks execution with a message saying it "cannot be verified". This is caused by the quarantine attribute set during download.

**Once in Terminal after unzipping:**

```bash
xattr -d com.apple.quarantine mbd-cli-arm64
```

**Or:** Right-click the file in Finder → **Open** → **Open Anyway**. Also visible afterwards in System Settings → Privacy & Security.

---

## Debugging

The tool automatically writes a `debug.log` file next to the executable, logging all arguments and any errors.

---

## Build

```
./build.sh
```

Produces three binaries in `dist/`:

| File | Target |
|---|---|
| `mbd-cli-arm64` | macOS Apple Silicon (native) |
| `mbd-cli-x64` | macOS Intel (via fyne-cross + Docker) |
| `mbd-cli.exe` | Windows AMD64 (via fyne-cross + Docker) |

Requirements for macOS Intel and Windows builds:
- Docker (must be running)
- `go install github.com/fyne-io/fyne-cross@latest`
- `go install fyne.io/fyne/v2/cmd/fyne@latest` *(old fyne CLI — ignore the deprecation notice, fyne-cross requires this version)*
