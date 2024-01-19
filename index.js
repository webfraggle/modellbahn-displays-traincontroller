const axios = require('axios');

// reroute console.log to file
var fs = require('fs');
var util = require('util');
var log_file = fs.createWriteStream('./debug.log', {flags : 'w'});
var log_stdout = process.stdout;

console.log = function() { //
    for (var i=0, numArgs = arguments.length; i<numArgs; i++){
        log_file.write(util.format(arguments[i]) + ' ');
        log_stdout.write(util.format(arguments[i]) + ' ');
    }
    log_file.write('\n');
    log_stdout.write('\n');
};
console.error = console.log;

// end reroute log to debug file


// first read arguments
var argv = require('minimist')(process.argv.slice(2));
console.log(argv);
console.log('execPath: ', process.execPath);
// console.log('report: ', process.report);
// console.log('report.filename: ', process.report.filename);
// console.log('report.directory: ', process.report.directory);

if (process.execPath.endsWith("tc.exe") || process.execPath.endsWith("tc-hidden.exe"))
{
    newDir = process.execPath.substring(0, process.execPath.lastIndexOf('\\'));
    console.log('Starting directory: ' + process.cwd());
    try {
        process.chdir(newDir);
        console.log('New directory: ' + process.cwd());
      }
      catch (err) {
        console.log('chdir: ' + err);
      }

}

if (process.execPath.endsWith("-macos") || process.execPath.endsWith("-arm64") || process.execPath.endsWith("-x64"))
{
    newDir = process.execPath.substring(0, process.execPath.lastIndexOf('/'));
    console.log('Starting directory: ' + process.cwd());
    try {
        process.chdir(newDir);
        console.log('New directory: ' + process.cwd());
      }
      catch (err) {
        console.log('chdir: ' + err);
      }

}




// second read config

if (argv['conf'])
{
    process.env.NODE_ENV = argv['conf'];
}
const config = require('config');
const endpoint = config.get('endpoint');

// for double track displays
if (argv['gleis'])
{
    if (argv['gleis'].slice(-1).toLowerCase() == "b")
    {
        path = "GleisB"
    } else {
        path = "GleisA"
    }
} else {
    path = "GleisA"
}

console.log("Endpoint:", endpoint, "path:", path);

if (argv['next'])
{
    url = endpoint+'/skipNext?path='+path;
    console.log('skip next, url:',url);
    axios
    .get(url)
    .then(function (response) {
        console.log('SUCCESS');
        process.exit(0);
    })
    .catch(function (error) {
        console.error('ERROR',error.toString());
        process.exit(1);
        return;
        // console.log('ERROR',error.response.status, error.response.data);
    });
}
if (argv['prev'])
{
    url = endpoint+'/skipPrev?path='+path;
    console.log('skip prev, url:',url);
    axios
    .get(url)
    .then(function (response) {
        console.log('SUCCESS');
        process.exit(0);
    })
    .catch(function (error) {
        console.error('ERROR',error.toString());
        process.exit(1);
        return;
        // console.log('ERROR',error.response.status, error.response.data);
    });
}

// node index.js --setTime "12:30"
if (argv['setTime'])
{
    timeString = argv['setTime'];
    const testRE = new RegExp(/^(\d{1,2}):(\d{1,2})$/);
    if (!testRE.test(timeString))
    {
        console.error('Wrong time format should be 12:23');
        process.exit(1);
        return;
    }
    var obj = {
        path: path,
        time: timeString
    };
    console.log("setTime obj:", obj);
    axios.postForm(endpoint+'/setTime', obj)
    .then(function (response) {
        console.log('SUCCESS',response.data);
        process.exit(0);
    })
    .catch(function (error) {
        console.error('ERROR',error.toString());
        process.exit(1);
        return;
        // console.log('ERROR',error.response.status, error.response.data);
    });
}

// node index.js --setTrain1 "ICE123|12:30|Berlin|Hannover - Wolfsburg|0|Kommt von der Commandline" --setTrain2 "RE50|21:12|Bebra|Hünfeld|+10|LOL" --setTrain3 "ICE3|09:45|Lübeck|Hamburg|0|"

var trains = []; 

// --setTrain1 "ICE123|12:30|Berlin|Hannover - Wolfsburg|0|Kommt von der Commandline"
if (argv['setTrain1'])
{
    trainString = argv['setTrain1'];
    trains.push(getTrainObject(1, trainString));
}

// --setTrain2 "RE50|21:12|Bebra|Hünfeld|+10|LOL"
if (argv['setTrain2'])
{
    trainString = argv['setTrain2'];
    trains.push(getTrainObject(2, trainString));
}


// --setTrain3 "ICE3|09:45|Lübeck|Hamburg|0|"
if (argv['setTrain3'])
{
    trainString = argv['setTrain3'];
    trains.push(getTrainObject(3, trainString));
}

function getTrainObject(nr, trainString)
{
    outJson = {
        "vonnach": "",
        "nr": "",
        "zeit": "",
        "via": "",
        "abw": 0,
        "hinweis": "",
        "fusszeile": "",
        "abschnitte": "",
        "reihung": "",
        "path": path
    }
    trainInfos = trainString.split("|");
    outJson.nr = trainInfos[0];
    outJson.zeit = trainInfos[1];
    outJson.vonnach = repairUTF8(trainInfos[2]);
    outJson.via = repairUTF8(trainInfos[3]);
    outJson.abw = trainInfos[4];
    outJson.hinweis = repairUTF8(trainInfos[5]);

    return {
        url:endpoint+'/zug'+nr,
        train: outJson
    }

}

if (trains.length)
{
    console.log("Sending Trains: ", trains);
    setTrains(trains);
} else {
    console.log("No Trains, Not sending");

}
function setTrains(trains)
{
    for (let i = 0; i < trains.length; i++) {
        const train = trains[i];
        axios.post(train.url, train.train)
        .then(function (response) {
            console.log('SUCCESS',response.data);
            // process.exit(0);
            // return;
        })
        .catch(function (error) {
            console.error('ERROR',error);
            // process.exit(1);
            // return;
            // console.log('ERROR',error.response.status, error.response.data);
        });
        
    }
}


// Functions

function repairUTF8(input) {
    const REPLACEMENTS = {
        "â‚¬": "€", "â€š": "‚", "â€ž": "„", "â€¦": "…", "Ë†": "ˆ",
        "â€¹": "‹", "â€˜": "‘", "â€™": "’", "â€œ": "“", "â€": "”",
        "â€¢": "•", "â€“": "–", "â€”": "—", "Ëœ": "˜", "â„¢": "™",
        "â€º": "›", "Å“": "œ", "Å’": "Œ", "Å¾": "ž", "Å¸": "Ÿ",
        "Å¡": "š", "Å½": "Ž", "Â¡": "¡", "Â¢": "¢", "Â£": "£",
        "Â¤": "¤", "Â¥": "¥", "Â¦": "¦", "Â§": "§", "Â¨": "¨",
        "Â©": "©", "Âª": "ª", "Â«": "«", "Â¬": "¬", "Â®": "®",
        "Â¯": "¯", "Â°": "°", "Â±": "±", "Â²": "²", "Â³": "³",
        "Â´": "´", "Âµ": "µ", "Â¶": "¶", "Â·": "·", "Â¸": "¸",
        "Â¹": "¹", "Âº": "º", "Â»": "»", "Â¼": "¼", "Â½": "½",
        "Â¾": "¾", "Â¿": "¿", "Ã€": "À", "Ã‚": "Â", "Ãƒ": "Ã",
        "Ã„": "Ä", "Ã…": "Å", "Ã†": "Æ", "Ã‡": "Ç", "Ãˆ": "È",
        "Ã‰": "É", "ÃŠ": "Ê", "Ã‹": "Ë", "ÃŒ": "Ì", "ÃŽ": "Î",
        "Ã‘": "Ñ", "Ã’": "Ò", "Ã“": "Ó", "Ã”": "Ô", "Ã•": "Õ",
        "Ã–": "Ö", "Ã—": "×", "Ã˜": "Ø", "Ã™": "Ù", "Ãš": "Ú",
        "Ã›": "Û", "Ãœ": "Ü", "Ãž": "Þ", "ÃŸ": "ß", "Ã¡": "á",
        "Ã¢": "â", "Ã£": "ã", "Ã¤": "ä", "Ã¥": "å", "Ã¦": "æ",
        "Ã§": "ç", "Ã¨": "è", "Ã©": "é", "Ãª": "ê", "Ã«": "ë",
        "Ã¬": "ì", "Ã­": "í", "Ã®": "î", "Ã¯": "ï", "Ã°": "ð",
        "Ã±": "ñ", "Ã²": "ò", "Ã³": "ó", "Ã´": "ô", "Ãµ": "õ",
        "Ã¶": "ö", "Ã·": "÷", "Ã¸": "ø", "Ã¹": "ù", "Ãº": "ú",
        "Ã»": "û", "Ã¼": "ü", "Ã½": "ý", "Ã¾": "þ", "Ã¿": "ÿ"
    }
    Object.entries(REPLACEMENTS).forEach(entry => {
        const [key, value] = entry;
        // console.log(key, value);
        input = input.replace(key, value);
    });
    return input;
}

