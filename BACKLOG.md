# mbd-cli Backlog

Ideen und geplante Features. Keine bestimmte Reihenfolge — dient als Sammlung für Diskussion und Planung.

---

## Ideen

### Linux-Unterstützung
**Status:** Idee

Da es viele verschiedene Linux-Distributionen und Architekturen gibt, ist eine vorkompilierte Binary wenig praktikabel. Stattdessen: Anleitung zum Selbstbauen.

- Go und CGO/Fyne müssen auf dem Zielsystem vorhanden sein (`sudo apt install golang gcc`)
- `go build` auf dem Zielsystem reicht aus — kein Cross-Compile nötig
- Eventuell Headless-Modus prüfen: Fyne benötigt eine Display-Umgebung; auf Server-Systemen ohne GUI würde der Start ohne Argumente fehlschlagen
- Mögliche Lösung: wenn kein Display vorhanden (`$DISPLAY` leer / kein Wayland), direkt mit Fehlerhinweis abbrechen statt zu crashen
