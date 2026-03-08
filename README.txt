mbd-cli — Modellbahn Displays CLI
===================================

WICHTIG: macOS – Datei kann nicht geöffnet werden
---------------------------------------------------
Wenn macOS meldet, dass die Datei nicht geprüft werden kann,
einmalig im Terminal eingeben:

    xattr -d com.apple.quarantine mbd-cli-arm64

Oder: Rechtsklick auf die Datei → Öffnen → Trotzdem öffnen,
danach unter Systemeinstellungen → Datenschutz & Sicherheit bestätigen.


WAS IST DAS?
------------
mbd-cli verbindet Modellbahn-Software (z.B. TrainController) mit
Zugzielanzeiger-Displays über deren REST-Schnittstelle.
Das Programm kehrt sofort zurück und blockiert die aufrufende
Software nicht.


KONFIGURATION
-------------
Beim ersten Start ohne Argumente öffnet sich die Konfigurationsoberfläche:

    macOS Apple Silicon:  ./mbd-cli-arm64
    macOS Intel:          ./mbd-cli-x64
    Windows:              .\mbd-cli.exe

Dort Endpoint-URL des Displays eintragen (z.B. http://192.168.178.155),
Verbindung testen und speichern. Mehrere Profile für mehrere Displays
möglich. Der eingebaute Befehlsgenerator hilft beim Zusammenstellen
der korrekten Befehle.


BEFEHLE (Beispiele)
-------------------
Nächsten Zug anzeigen:
    mbd-cli.exe --next

Uhrzeit setzen:
    mbd-cli.exe --setTime "12:30"

Zuginfo setzen (Nr|Zeit|Ziel|Via|Verspätung|Info):
    mbd-cli.exe --setTrain1 "ICE123|12:30|Berlin|Hannover|0|"

Alle drei Zuginfos auf einmal:
    mbd-cli.exe --setTrain1 "..." --setTrain2 "..." --setTrain3 "..."

Bild anzeigen:
    mbd-cli.exe --image 00logo.png

Zweites Gleis (Doppelspuranzeige):
    mbd-cli.exe --gleis B --next

Anderes Display-Profil verwenden:
    mbd-cli.exe --conf gleis1 --setTime "12:30"

Timeout anpassen (Millisekunden):
    mbd-cli.exe --timeout 5000 --next


FEHLERBEHEBUNG
--------------
Das Programm schreibt eine debug.log neben der Programmdatei.
Dort sind alle aufgerufenen Argumente und Fehler protokolliert.


===================================
ENGLISH VERSION
===================================

IMPORTANT: macOS – File cannot be opened
-----------------------------------------
If macOS says it cannot verify the file, run this once in Terminal:

    xattr -d com.apple.quarantine mbd-cli-arm64

Or: Right-click the file → Open → Open Anyway,
then confirm in System Settings → Privacy & Security.


WHAT IS THIS?
-------------
mbd-cli connects model railway software (e.g. TrainController) with
Zugzielanzeiger train destination displays via their REST API.
The program returns immediately and does not block the calling software.


CONFIGURATION
-------------
When launched without arguments, the configuration UI opens:

    macOS Apple Silicon:  ./mbd-cli-arm64
    macOS Intel:          ./mbd-cli-x64
    Windows:              .\mbd-cli.exe

Enter the display's endpoint URL (e.g. http://192.168.178.155),
test the connection and save. Multiple profiles for multiple displays
are supported. The built-in command builder helps assemble the correct
commands to use in your model railway software.


COMMANDS (examples)
-------------------
Skip to next train:
    mbd-cli.exe --next

Set time:
    mbd-cli.exe --setTime "12:30"

Set train info (Nr|Time|Dest|Via|Delay|Info):
    mbd-cli.exe --setTrain1 "ICE123|12:30|Berlin|Hannover|0|"

Set all three train slots at once:
    mbd-cli.exe --setTrain1 "..." --setTrain2 "..." --setTrain3 "..."

Show image:
    mbd-cli.exe --image 00logo.png

Second track (double-track display):
    mbd-cli.exe --gleis B --next

Use a different display profile:
    mbd-cli.exe --conf gleis1 --setTime "12:30"

Adjust timeout (milliseconds):
    mbd-cli.exe --timeout 5000 --next


TROUBLESHOOTING
---------------
The program writes a debug.log file next to the executable.
It contains all arguments and any errors that occurred.
