# Modellbahn Displays - TrainController Anbindung
## How to use
Go to config folder and edit default.json.
Enter your display's IP-Adresse

    {
    "endpoint":"http://192.168.178.155"
    }

### Option 1: Skip next and prev
Manage all your trains with the webinterface of your "Zugzielanzeiger" and the skip next or previous train via command-line / external program with this tool. 
    
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

## How to debug in TrainController
Use debug.bat. Change path to your folder and use the debug.bat instead of the mbd-tc.exe. 
This bat-file opens an extra notepad, this prevents the terminal to close and you can read the log messages.

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

To prevent showing up a command prompt window, use PE Tools to change Subsystem from 3 to 2.
Explanation here:
https://stackoverflow.com/questions/22653010/prevent-console-window-from-being-created-in-custom-node-js-build