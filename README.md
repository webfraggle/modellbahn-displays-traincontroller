# Modellbahn Displays - TrainController Anbindung

(English version below)
## Konfiguration
Gehen Sie zum Konfigurationsordner und bearbeiten Sie default.json. Geben
Sie die IP-Adresse Ihres Displays ein:

    {
    "endpoint":"http://192.168.178.155"
    }

### Option 1: Zu nächsten und vorherigen Zugziel springen:
Verwalten Sie alle Ihre Züge über die Weboberfläche Ihres
Zugzielanzeigers und springen Sie dann per Kommandozeile / externem Programm zum nächsten oder vorherigen Zug.

    mbd-tc.exe --next
    mbd-tc.exe --prev
### Option 2: Aufruf der Anzeigen vom Webinterface über externes Programm (z.B.TC):
Verwalten Sie alle Ihre Züge über das Webinterface Ihres Zugzielanzeigers und stellen Sie die Uhrzeit per Kommandozeile / externem Programm ein.

    mbd-tc.exe --setTime "12:30"
### Option 3: Zuginformationen direkt in externem Programm (z.B. TC) einstellen:
Verwenden Sie die Optionen setTrain1 bis setTrain3, um alle Informationen direkt festzulegen. Sie können einen oder mehrere gleichzeitig verwenden.
Der Wert ist eine durch 5 Pipes getrennte Zeichenfolge mit folgendem Schema:

„TrainID|Time|Destination|Via|Delay|Special info“

    mbd-tc.exe --setTrain1 "ICE123|12:30|Berlin|Hannover - Wolfsburg|0|Kommt von der Commandline" --setTrain2 "RE50|21:12|Bebra|Hünfeld|+10|LOL" --setTrain3 "ICE3|09:45|Lübeck|Hamburg|0|"

### Einsatz mehrerer Displays:
Erstellen Sie eine Kopie des Standard-JSON. Benennen Sie es in den Namen Ihres Displays um, z. B. „gleis1.json“. Verwenden Sie die Befehlszeilenoption --conf, um diese neue Konfiguration zu laden:

    mbd-tc.exe --conf gleis1 --setTime "12:30"
    
### Wie man eine Doppelspuranzeige verwendet:
Verwenden Sie den Parameter „gleis“ mit „A“ oder „B“.

    mbd-tc.exe --gleis B --setTime "12:30"

### Ein leeres Display anzeigen (Doppelspuranzeige):
    mbd-tc.exe –-gleis B --setTrain1 "|||||"

### Bild anzeigen
    --image 00Logo.png
Damit kann man in den Modi Manuell, Interval und Bilder ein Bild, welches auf den Controller geladen wurde, anzeigen. Das Bild wird so lange angezeigt, bis ein neuer Zug oder ein neues Bild angezeigt wird.

### Spezialoptionen
    --timeout 1000
Das http-timeout in ms um bei fehlerhafter Verbindung das Tool schneller zu beenden.

### So verhindern Sie das Popup eines Befehlszeilenfensters:
Verwenden Sie die Datei **mbd-tc-hidden.exe** anstelle von mbd-tc.exe

### So debuggen Sie:
Verwenden Sie debug.bat. Ändern Sie den Pfad zu Ihrem Ordner und verwenden Sie debug.bat anstelle von mbd-tc.exe. Diese Bat-Datei öffnet einen zusätzlichen Texteditor, der verhindert, dass das Terminal geschlossen wird und Sie die Protokollmeldungen lesen können.

# English Version
## How to use
Go to config folder and edit default.json.
Enter your display's IP-Adresse

    {
    "endpoint":"http://192.168.178.155"
    }

### Option 1: Skip next and prev
Manage all your trains with the webinterface of your "Zugzielanzeiger" and then skip to next or previous train via command-line / external program with this tool. 
    
    mbd-tc.exe --next
    mbd-tc.exe --prev

### Option 2: setTime
Manage all your trains with the webinterface of your "Zugzielanzeiger" and the set the time via command-line / external program with this tool. 
    
    mbd-tc.exe --setTime "12:30"
### Option 3: set train infos directly
Use option setTrain1 to setTrain3 to set all the infos directly. You can use one ore more at a time. The value is a pipe separted string. With the following schema: "TrainID|Time|Destination|Via|Delay|Special info"
 
    
    mbd-tc.exe --setTrain1 "ICE123|12:30|Berlin|Hannover - Wolfsburg|0|Kommt von der Commandline" --setTrain2 "RE50|21:12|Bebra|Hünfeld|+10|LOL" --setTrain3 "ICE3|09:45|Lübeck|Hamburg|0|"

### Use of more displays
Create a copy of the default json. Rename to your display's name, e.g. "gleis1.json". Use the command-line option --conf to load this new config.

    mbd-tc.exe --conf gleis1 --setTime "12:30"

### how to use a double track display
Use the "gleis" parameter with "A" or "B".

    mbd-tc.exe --gleis B --setTime "12:30"

### Show an Image
    --image 00Logo.png
This allows you to display an image that has been loaded onto the controller in Manual, Interval and Images modes. The picture is displayed until a new move or a new picture is displayed.

### Special Option
    --timeout 1000
The http-timeout in ms for cancelling the tool earlier because of connection issues.

## How to prevent the popup of a command line window
Use the mbd-tc-hidden.exe file instead of mbd-tc.exe

## How to debug
Use debug.bat. Change path to your folder and use the debug.bat instead of the mbd-tc.exe. 
This bat-file opens an extra notepad, this prevents the terminal to close and you can read the log messages.

# How to develop.
Please use nvm to switch to node version 18.16
    
    nvm use 18.16
    npm install



## How to build executable.
I used "pkg" to build Windows & MacOS executables.

You need the following node modules
 
    npm install -g pkg
    npm i -g @vercel/ncc

Then you can build it with this command

    ncc build index.js -o build
    pkg ./build/index.js -t node18-win-x64,node18-macos-arm64,node18-macos-x64 -o ./dist/mbd-tc

    #### Mac-Only
    pkg ./build/index.js -t node18-macos-arm64,node18-macos-x64 -o ./dist/mbd-tc
or 
    
    npm run build

To prevent showing up a command prompt window, use PE Tools to change Subsystem from 3 to 2.
Explanation here:
https://stackoverflow.com/questions/22653010/prevent-console-window-from-being-created-in-custom-node-js-build
