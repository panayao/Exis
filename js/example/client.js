
// "Real" client.js testing code
var riffle = require('jsriffle');

riffle.SetFabricLocal();
riffle.SetLogLevelDebug();

var app = riffle.Domain("xs.damouse");
var receiver = app.subdomain("alpha");
var me = app.subdomain("beta");


me.onJoin = function() {
    console.log("Sender Joined");

    receiver.call("iGiveInts", "Hi").then(riffle.wait(function(a) {
        console.log("Result: ", a);
    }, [Number]),
    function (err) {
        console.log("ERROR: ", err); 
    });

};

me.join()