# mbd-cli unter Linux kompilieren

Diese Anleitung beschreibt Schritt für Schritt, wie du **mbd-cli** auf einem Linux-System aus dem Quellcode kompilierst.

---

## Voraussetzungen

### 1. Systempakete installieren

**mbd-cli** nutzt [Fyne](https://fyne.io/) als GUI-Framework, das CGO und OpenGL-Bibliotheken benötigt. Zusätzlich wird `git` zum Klonen des Repositories benötigt.

**Debian / Ubuntu / Raspberry Pi OS:**
```bash
sudo apt update
sudo apt install -y git gcc libgl1-mesa-dev xorg-dev
```

**Fedora / RHEL / CentOS:**
```bash
sudo dnf install -y git gcc mesa-libGL-devel libX11-devel
```

**Arch Linux:**
```bash
sudo pacman -S --needed git gcc mesa libx11 libxcursor libxrandr libxinerama libxi
```

---

### 2. Go 1.21 oder neuer installieren

Prüfe zuerst, ob Go bereits installiert ist:
```bash
go version
```

Falls nicht vorhanden oder die Version zu alt ist, **direkt von go.dev laden** (nicht über `apt`):

> Ubuntu/Debian liefern über `apt install golang-go` oft veraltete Versionen. `gccgo-go` ist ein alternativer Compiler und **nicht** der offizielle Go-Compiler — beide apt-Varianten werden nicht empfohlen.

```bash
# Aktuelle Version von go.dev herunterladen (AMD64 — Beispiel: 1.24.1)
wget https://go.dev/dl/go1.24.1.linux-amd64.tar.gz

# Alte Installation entfernen und neu entpacken
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.24.1.linux-amd64.tar.gz

# PATH dauerhaft setzen (in ~/.bashrc oder ~/.profile eintragen)
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

> **ARM-Systeme (z. B. Raspberry Pi):** `linux-arm64.tar.gz` (64-Bit) bzw. `linux-armv6l.tar.gz` (32-Bit) von [go.dev/dl](https://go.dev/dl/) verwenden.

Installationserfolg prüfen:
```bash
go version
# Ausgabe z. B.: go version go1.22.4 linux/amd64
```

---

## Repository klonen

```bash
# Beliebiges Zielverzeichnis wählen, z. B. im Home-Verzeichnis
mkdir -p ~/projekte
cd ~/projekte

# Repository klonen
git clone https://github.com/webfraggle/modellbahn-displays-traincontroller.git mbd-cli
cd mbd-cli
```

---

## Abhängigkeiten laden

```bash
go mod download
```

---

## Kompilieren

### Binärdatei für das aktuelle System erstellen

```bash
go build -o mbd-cli .
```

Die fertige Binärdatei `mbd-cli` liegt danach im aktuellen Verzeichnis.

### Alternativ: Mit Optimierungen (kleinere Dateigröße)

```bash
go build -ldflags "-s -w" -o mbd-cli .
```

---

## Konfigurationsverzeichnis anlegen

Die Binärdatei erwartet einen `config/`-Ordner **neben der ausführbaren Datei** mit mindestens einer JSON-Datei:

```bash
mkdir -p config

# Standardkonfiguration anlegen
cat > config/default.json <<EOF
{ "endpoint": "http://192.168.178.155" }
EOF
```

Passe die IP-Adresse auf die deines Zugzielanzeiger-Controllers an.

---

## Starten

### Konfigurations-Oberfläche öffnen (ohne Argumente):
```bash
./mbd-cli
```

### Direkt verwenden (Beispiele):
```bash
./mbd-cli --next
./mbd-cli --setTime "12:30"
./mbd-cli --setTrain1 "ICE123|12:30|Berlin|Hannover - Wolfsburg|0|"
```

---

## Verzeichnisstruktur nach dem Build

```
mbd-cli/           ← geklontes Repository
  mbd-cli          ← kompilierte Binärdatei
  config/
    default.json   ← deine Konfiguration
  debug.log        ← wird automatisch beim ersten Aufruf erstellt
```

---

## Häufige Fehler

| Fehler | Ursache | Lösung |
|---|---|---|
| `gcc: command not found` | C-Compiler fehlt | `sudo apt install gcc` |
| `GL/gl.h: No such file or directory` | OpenGL-Header fehlen | `sudo apt install libgl1-mesa-dev` |
| `X11/Xlib.h: No such file or directory` | X11-Header fehlen | `sudo apt install xorg-dev` |
| `go: command not found` | Go nicht im PATH | `source ~/.bashrc` oder Terminal neu starten |
| `config/default.json: no such file` | Konfigurationsordner fehlt | `mkdir -p config && echo '{"endpoint":"http://..."}' > config/default.json` |

---

## Hinweis: Headless-Systeme (ohne Grafikoberfläche)

Die GUI (`./mbd-cli` ohne Argumente) benötigt eine laufende X11- oder Wayland-Sitzung. Auf Servern oder im Headless-Betrieb (z. B. reiner SSH-Zugriff) kann die GUI nicht geöffnet werden.

Die **CLI-Befehle** (`--next`, `--setTrain1`, `--setTime` usw.) funktionieren jedoch vollständig ohne Grafikoberfläche.
