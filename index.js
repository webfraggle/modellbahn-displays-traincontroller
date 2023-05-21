const axios = require('axios');

// first read arguments
var argv = require('minimist')(process.argv.slice(2));
console.log(argv);

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
        path = "GleisA"
    } else {
        path = "GleisA"
    }
} else {
    path = "GleisA"
}

console.log("Endpoint:", endpoint, "path:", path);


// set Time
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
        console.log('SUCCESS',response);
    })
    .catch(function (error) {
        console.log('ERROR',error);
        // console.log('ERROR',error.response.status, error.response.data);
    });
}