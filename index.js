const axios = require('axios');

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
    axios.postForm(endpoint+'/setTime', {
        path: path,
        time: timeString
    })
    .then(function (response) {
        console.log('SUCCESS',response.data);
    })
    .catch(function (error) {
        console.error('ERROR',error);
        process.exit(1);
        return;
        // console.log('ERROR',error.response.status, error.response.data);
    });
}

// node index.js --setTrain1 "ICE123|12:30|Berlin|Hannover - Wolfsburg|0|Kommt von der Commandline" --setTrain2 "RE50|21:12|Bebra|Hünfeld|+10|LOL" --setTrain3 "ICE3|09:45|Lübeck|Hamburg|0|"

// --setTrain1 "ICE123|12:30|Berlin|Hannover - Wolfsburg|0|Kommt von der Commandline"
if (argv['setTrain1'])
{
    trainString = argv['setTrain1'];
    setTrain(1, trainString);
}

// --setTrain2 "RE50|21:12|Bebra|Hünfeld|+10|LOL"
if (argv['setTrain2'])
{
    trainString = argv['setTrain2'];
    setTrain(2, trainString);
}


// --setTrain3 "ICE3|09:45|Lübeck|Hamburg|0|"
if (argv['setTrain3'])
{
    trainString = argv['setTrain3'];
    setTrain(3, trainString);
}

function setTrain(nr, trainString)
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

    // console.log(outJson);

    axios.post(endpoint+'/zug'+nr, outJson)
    .then(function (response) {
        console.log('SUCCESS',response.data);
        process.exit(0);
        return;
    })
    .catch(function (error) {
        console.error('ERROR',error);
        process.exit(1);
        return;
        // console.log('ERROR',error.response.status, error.response.data);
    });
}