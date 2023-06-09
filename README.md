# Modellbahn Displays - TrainController Anbindung
## How to use
Go to config folder and edit default.json.
Enter your display's IP-Adresse

    {
    "endpoint":"http://192.168.178.155"
    }
### Option 1: setTime
Manage all your trains with the webinterface of your "Zugzielanzeiger" and the set the time via command-line / external program with this tool. 
    
    mbd-tc-win.exe --setTime "12:30"
### Option 2: set train infos directly
 
    
    mbd-tc-win.exe --setTrain1 "ICE123|12:30|Berlin|Hannover - Wolfsburg|0|Kommt von der Commandline" --setTrain2 "RE50|21:12|Bebra|Hünfeld|+10|LOL" --setTrain3 "ICE3|09:45|Lübeck|Hamburg|0|"

### Use of more displays
Create a copy of the default json. Rename to your display's name, e.g. "gleis1.json". Use the command-line option --conf to load this new config.

    mbd-tc-win.exe --conf gleis1 --setTime "12:30"

### how to use a double track display
Use the path parameter with "GleisA" or "GleisB".

    mbd-tc-win.exe --path GleisB --setTime "12:30"

## How to develop.
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
    pkg ./build/index.js -t node18-win-arm64,node18-macos-arm64 -o ./dist/mbd-tc
or 
    
    npm run build