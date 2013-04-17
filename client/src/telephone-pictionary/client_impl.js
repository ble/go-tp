goog.provide('ble.telephone_pictionary.ClientImpl');

goog.require('ble.telephone_pictionary.Client');

goog.require('ble._2d.DrawPart');

goog.require('goog.result.Result');
goog.require('goog.labs.net.xhr');
goog.require('goog.Uri');

goog.scope(function() { 

var _ = ble.telephone_pictionary;

var SimpleResult = goog.result.SimpleResult;
var transform = goog.result.transform;
var JSON = window.JSON;
var xhr = goog.labs.net.xhr;

_.jsonHeader = {'headers': {'Content-Type': 'application/json'}};
/**
 * @constructor
 * @param {string} gameUrl
 * @implements _.Client
 */
_.ClientImpl = function(gameUrl) {
  this.url = gameUrl;
  var urlTemp = new goog.Uri(this.url);
  this.baseUrl = urlTemp.resolve(new goog.Uri('../../')).toString();
};

_.ClientImpl.prototype.getGameState = function() {
  return transform(xhr.get(this.url), JSON.parse);
};

_.ClientImpl.prototype.getStack = function(stackId) {
  return transform(xhr.get(this.baseUrl + 'stack/' + stackId), JSON.parse);
};

_.ClientImpl.prototype.getDrawing = function(drawingId) {
  return transform(xhr.get(this.baseUrl + 'drawing/' + drawingId), JSON.parse);
};

_.ClientImpl.prototype.appendToDrawing = function(drawingId, part) {
  //TODO: good place for ClientImpl to hook into model of current game state
  var action = {'actionType': 'draw', 'content': part};
  action = JSON.stringify(action);
  return xhr.post(this.baseUrl + 'drawing/' + drawingId, action, _.jsonHeader);
};

_.ClientImpl.prototype.passStack = function(stackId){
  //TODO: this would be a good place for ClientImpl to hook into the
  //model of the current game state and enforce the logical rule
  //that you can only pass the stack that you are currently working on,
  //etc.
  var action = {'actionType': 'passStack'};
  action = JSON.stringify(action);
  return xhr.post(this.url + 'pass', action, _.jsonHeader);
};

_.ClientImpl.prototype.chat = function(message){
  var action = {'actionType': 'chat', 'content': message};
  action = JSON.stringify(action);
  return xhr.post(this.url + 'chat', action, _.jsonHeader);
};


});
