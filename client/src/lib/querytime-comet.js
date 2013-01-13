goog.require('ble.net.CometLoop');
goog.require('goog.net.EventType');

goog.provide('ble.net.QueryTimeComet');

goog.scope(function() {
var netType = goog.net.EventType;
var JSON = window.JSON;
var console = window.console;

/**
 * @constructor
 * @param{string} uri
 * @extends{ble.net.CometLoop}
 */
ble.net.QueryTimeComet = function(uri) {
  ble.net.CometLoop.call(this, uri, 10000, 2500, 2500);
  this.lastQuery = null;
};
goog.inherits(ble.net.QueryTimeComet, ble.net.CometLoop);

var QTCp = ble.net.QueryTimeComet.prototype;

/**
 * @param{number} lastTime
 */
QTCp.runAt = function(lastTime) {
  this.lastQuery = lastTime;
  this.run();
}

QTCp.getUri = function() {
  var baseUri = ble.net.QueryTimeComet.superClass_.getUri.call(this);
  if(this.lastQuery == null)
    return baseUri;
  return baseUri + "?lastTime=" + this.lastQuery;
};

QTCp.handleEvent = function(event) {
  if(event.type == netType.SUCCESS) {
    var jsonObj = JSON.parse(this.xhr.getResponseText());
    this.lastQuery = jsonObj['lastTime'];
  }
  ble.net.QueryTimeComet.superClass_.handleEvent.call(this, event);
};


});
