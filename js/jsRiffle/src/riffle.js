
require('./go.js');
var want = require('./want.js');
var ws = require('./transport/websocket.js');
var pjson = require('../package.json');
exports.version = pjson.version;

// Used to counteract uint generation on seemlingly 32 bit platforms
global.NewID = function() {
   return Math.floor(Math.random() * 9007199254740992);
};

// Dont need any of this-- just return the conn
var Ws = function () {
    this.open = function(url) {
        // Methods available on the conn: console.log, protocol, send, close, onmessage, onopen, onclose, info
        var factory = new ws.Factory({'type': 'websocket', 'url': url});
        this.conn = factory.create();
        this.conn.onmessage = this.onmessage;
        this.conn.onopen = this.onopen;
        this.conn.onclose = this.onclose;
    }
}; 

global.Renamer = function(domain) {
	for (var func in domain) {
		domain[func.substr(0, 1).toLowerCase() + func.substr(1)] = domain[func];
		delete domain[func];
	}
}

// Intercepts .then and sends down cumin args to the core. 
// Should only be used by Calls, and internally at that 
global.PromiseInterceptor = function(trueHandler, domain, cb, other) {
    // console.log("Interceptor receiving: " + arguments);

    return function(args, trueErrback) {
        console.log("args and b: ", trueErrback)
        domain.callExpects(cb, args.types);
        trueHandler(args.fp, trueErrback)
    }
}

global.WsWrapper = Ws;
exports.Domain = global.Domain.New;

exports.SetLogLevelOff = global.Config.SetLogLevelOff;
exports.SetLogLevelApp = global.Config.SetLogLevelApp;
exports.SetLogLevelErr = global.Config.SetLogLevelErr;
exports.SetLogLevelWarn = global.Config.SetLogLevelWarn;
exports.SetLogLevelInfo = global.Config.SetLogLevelInfo;
exports.SetLogLevelDebug = global.Config.SetLogLevelDebug;

exports.SetFabricDev = global.Config.SetFabricDev;
exports.SetFabricSandbox = global.Config.SetFabricSandbox;
exports.SetFabricProduction = global.Config.SetFabricProduction;
exports.SetFabricLocal = global.Config.SetFabricLocal;
exports.SetFabric = global.Config.SetFabric;

exports.Application = global.Config.Application;
exports.Debug = global.Config.Debug;
exports.Info = global.Config.Info;
exports.Warn = global.Config.Warn;
exports.Error = global.Config.Error;

//want.js exports
exports.want = want.want;
exports.wait = want.wait;
exports.ModelObject = want.ModelObject;
