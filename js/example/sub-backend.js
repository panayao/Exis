var riffle = require('jsriffle');

riffle.setFabricLocal();
riffle.setLogLevelDebug();

var app = riffle.Domain("xs.demo.test");
var backend = app.subdomain("backend");
var client = app.subdomain("client");


function test(){
  console.log(arguments);
  console.log(arguments[0].fullname());

  //TODO This doesn't actually every return to the caller
  return "String";
}


backend.onJoin = function() {

    // Example Pub/Sub Basic - a very basic pub/sub example
    this.subscribe("basicSub", riffle.want(function(s) {
        console.log(s); // Expects a String, like "Hello"
    }, String));
    // End Example Pub/Sub Basic
    
    // Example Pub/Sub Basic Two - a basic pub/sub example
    this.subscribe("basicSubTwo", riffle.want(function(s, i) {
        console.log(s, i); // Expects a String, like "Hello 3"
    }, String, Number));
    // End Example Pub/Sub Basic Two

};


backend.join()


