goog.require('goog.events');
goog.require('goog.net.EventType');

goog.require('ble.net.CometLoop');

goog.provide('ble.tpg.game.setupClient');


goog.scope(function() {
var netType = goog.net.EventType;
var scope = ble.tpg.game;
scope.onResponse = function(event) {
  window.console.log(event);
};

scope.setupClient = function() {
  scope.cometLoop = new ble.tpg.game.QueryTimeComet("./events");
  var cometLoop = scope.cometLoop;
  cometLoop.run();
};

ble.tpg.game.QueryTimeComet = function(uri) {
  ble.net.CometLoop.call(this, uri, 10000, 100, 2500);
  this.lastQuery = null;
};
goog.inherits(ble.tpg.game.QueryTimeComet, ble.net.CometLoop);

var QTCp = ble.tpg.game.QueryTimeComet.prototype;

QTCp.getUri = function() {
  var baseUri = ble.tpg.game.QueryTimeComet.superClass_.getUri.call(this);
  if(this.lastQuery == null)
    return baseUri;
  return baseUri + "?lastQuery=" + this.lastQuery;
};

QTCp.handleEvent = function(event) {
  ble.tpg.game.QueryTimeComet.superClass_.handleEvent.call(this, event);
  if(event.type == netType.SUCCESS) {
    var jsonObj = JSON.parse(this.xhr.getResponseText());
    this.lastQuery = jsonObj['queryTime'];
  }
};



scope.setupClient();
});
