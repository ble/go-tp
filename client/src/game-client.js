goog.require('goog.events');
goog.require('goog.net.EventType');
goog.require('goog.net.XhrIo');

goog.require('ble.net.CometLoop');

goog.provide('ble.tpg.game.setupClient');


goog.scope(function() {
var netType = goog.net.EventType;
var scope = ble.tpg.game;
var console = window.console;
var JSON = window.JSON;

scope.onResponse = function(event) {
  console.log(event);
};

scope.setupClient = function() {
  scope.cometLoop = new ble.tpg.game.QueryTimeComet("./events");
  var cometLoop = scope.cometLoop;
  var otherXhr = new goog.net.XhrIo();
  scope.otherXhr = otherXhr;
  goog.events.listenOnce(
      otherXhr, 
      [netType.SUCCESS],
      function(event) {
        console.log(this.getResponse());
        try {
          var jsonObj = JSON.parse(this.getResponse());
          var lastTime = jsonObj['lastTime'] || 0;
          cometLoop.runAt(lastTime);
        } catch(e) {
          console.log(e);
          console.log(event);
          console.log(this.getResponse());
          console.log("frickin' error");
        }
      });
  goog.events.listenOnce(
      otherXhr, 
      [netType.TIMEOUT, netType.ERROR],
      function(event) {
        console.log(event);
        console.log("frickin' error");
      });


  otherXhr.send("./")
  //cometLoop.run();
};

/**
 * @constructor
 * @param{string} uri
 * @extends{ble.net.CometLoop}
 */
ble.tpg.game.QueryTimeComet = function(uri) {
  ble.net.CometLoop.call(this, uri, 10000, 2500, 2500);
  this.lastQuery = null;
};
goog.inherits(ble.tpg.game.QueryTimeComet, ble.net.CometLoop);

var QTCp = ble.tpg.game.QueryTimeComet.prototype;

/**
 * @param{number} lastTime
 */
QTCp.runAt = function(lastTime) {
  this.lastQuery = lastTime;
  this.run();
}

QTCp.getUri = function() {
  var baseUri = ble.tpg.game.QueryTimeComet.superClass_.getUri.call(this);
  if(this.lastQuery == null)
    return baseUri;
  return baseUri + "?lastTime=" + this.lastQuery;
};

QTCp.handleEvent = function(event) {
  ble.tpg.game.QueryTimeComet.superClass_.handleEvent.call(this, event);
  if(event.type == netType.SUCCESS) {
    var jsonObj = JSON.parse(this.xhr.getResponseText());
    this.lastQuery = jsonObj['lastTime'];
    var events = jsonObj['events'] || [];
    for(var i = 0; i < events.length; i++) {
      console.log(events[i]);
    }
  }
};



scope.setupClient();
});
