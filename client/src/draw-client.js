goog.require('goog.events');
goog.require('goog.events.Event');

goog.require('goog.net.XhrIo');
goog.require('goog.net.EventType');

goog.require('ble.scribble.UI');
goog.require('ble.scribble.MutableDrawing');
goog.require('ble.scribbleDeserializer');
goog.require('ble.scribble.EventType');

goog.provide('ble.tpg.drawing.setupClient');

goog.scope(function() {
var JSON = window.JSON;

var isSomeVal = goog.isDefAndNotNull;
var XhrIo = goog.net.XhrIo;
var nEventType = goog.net.EventType;

var Event = goog.events.Event;
var listen = goog.events.listen;
var listenOnce = goog.events.listenOnce;

var UI = ble.scribble.UI;
var Drawing = ble.scribble.MutableDrawing;
var deserializer = ble.scribbleDeserializer;
var sEventType = ble.scribble.EventType;

var scope = ble.tpg.drawing;

scope.width = 320;
scope.height = 240;

scope.setupClient = function() {
  var x = new XhrIo();
  listenOnce(x, nEventType.SUCCESS, scope.createUI, false, x);
  x.send("./", "GET");
};

/** @type function(this:XhrIo, Event) */
scope.createUI = function(event) {
  var xhr = /**@type {XhrIo}*/ this;
  var drawing = scope.deserializeDrawing(JSON.parse(xhr.getResponseText()));
  if(!isSomeVal) {
    scope.showError("Failed to deserialize drawing.");
    return;
  }
  scope.ui = new ble.scribble.UI(scope.width, scope.height);
  var ui = scope.ui;
  ui.render(document.body);
  ui.canvas.drawing = drawing;
  ui.canvas.withContext(ui.canvas.repaintComplete);
  scope.wireUpUI(xhr);
};

scope.showError = function(msgOrSomething) {
  window.console.log(msgOrSomething);
};

scope.wireUpUI = function(xhr) {
    listen(scope.ui, sEventType.DRAW_END, scope.onDrawEnd, false, xhr);
};

scope.deserializeDrawing = function(json) {
  var parts = [];
  var time = Infinity;
  var d = deserializer;
  for(var i = 0; i < json.length; i++) {
    var part = d.deserialize(json[i]); 
    if(part) {
      time = Math.min(time, part.start());
      parts.push(part);
    }
  }
  var drawing = new ble.scribble.MutableDrawing(time, parts);
  return drawing;
};

/** @type function(this:XhrIo, Event) */
scope.onDrawEnd = function(event) {
  var xhr = /**@type {XhrIo}*/ this;
  var part = event.drawn;
  xhr.send("./", "POST", JSON.stringify(part), {'Content-Type': 'application/json'}); 
}
scope.setupClient();
});
