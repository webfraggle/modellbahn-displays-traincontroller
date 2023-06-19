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


// first read arguments
var argv = require('minimist')(process.argv.slice(2));
// console.log(argv);

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
    outJson.vonnach = trainInfos[2];
    outJson.via = trainInfos[3];
    outJson.abw = trainInfos[4];
    outJson.hinweis = trainInfos[5];

    return {
        url:endpoint+'/zug'+nr,
        train: outJson
    }

}

console.log("Trains", trains);
if (trains.length)
{
    setTrains(trains);
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