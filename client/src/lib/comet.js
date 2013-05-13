goog.require('goog.net.XhrIo');
goog.require('goog.events.EventTarget');
goog.require('goog.events.Event');

goog.require('goog.Uri');
goog.require('goog.Uri.QueryData');

goog.require('goog.functions');

goog.provide('ble.net.CometLoop');
goog.provide('ble.net.EventType');

////////////////////////////////////////////////////////////////////////////////
                                                        goog.scope(function(){
////////////////////////////////////////////////////////////////////////////////

/**@enum{string}*/
ble.net.EventType = ({
  COMET_START: 'COMET_START',
  COMET_STOP: 'COMET_STOP',
  COMET_DATA: 'COMET_DATA'});

var Event = goog.events.Event;
var netType = goog.net.EventType;
var cType = ble.net.EventType;


//Basic comet loop:
//  hit a fixed URL with get requests, again and again.
/**
 * @constructor
 * @extends {goog.events.EventTarget}
 * @param {string} uri
 * @param {number} timeout milliseconds
 * @param {number} successWait milliseconds
 * @param {number} retryWait milliseconds
 */
ble.net.CometLoop =
 function(uri, timeout, successWait, retryWait) {
  goog.events.EventTarget.call(this);
  this.uri = uri
  this.timeout = timeout;
  this.successWait = successWait;
  /** @type {function(number):number}*/
  this.retryWait = goog.functions.constant(retryWait);

  this.previousFailures = 0;
  this.running = false;
  this.xhr = new goog.net.XhrIo();
  this.xhr.setTimeoutInterval(this.timeout);
  this.xhr.setParentEventTarget(this);
  goog.events.listen(
      this.xhr,
      [netType.SUCCESS, netType.TIMEOUT, netType.ERROR],
      this);
  this.pendingSend = null;
};
goog.inherits(ble.net.CometLoop, goog.events.EventTarget);

var bnCL = ble.net.CometLoop.prototype;

bnCL.run = function() {
  if(this.running)
    return;
  this.running = true;
  this.send_();
  this.dispatchEvent(new Event(cType.COMET_START));
};

bnCL.stop = function() {
  if(!this.running)
    return;
  this.xhr.abort();
  this.running = false;
  if(this.pendingSend !== null)
    window.clearTimeout(this.pendingSend);
  this.pendingSend = null;
  this.previousFailures = 0;
  this.dispatchEvent(new Event(cType.COMET_STOP));
};
/**
 * @protected
 */
bnCL.getUri = function() {
  return this.uri;
};

/**
 * @private
 */
bnCL.send_ = function() {
  this.pendingSend = null;
  var uri = this.getUri();
  this.preSend(uri);
  this.xhr.send(uri);
};

/**
 * @protected
 */
bnCL.preSend = function(uri) {};

/**
 * @param {goog.events.Event} event
 */
bnCL.handleEvent = function(event) {
  switch (event.type) {
    case netType.ERROR:
    case netType.TIMEOUT:
      var delay = this.retryWait(this.previousFailures);
      this.previousFailures++;
      this.pendingSend = window.setTimeout(goog.bind(this.send_, this), delay);
      this.dispatchEvent(event);
      break;
    case netType.ABORT:
      break;
    case netType.SUCCESS:
      this.processSuccess(event);
      var delay = this.successWait;
      var dataEvent = new Event(ble.net.EventType.COMET_DATA, this);
      dataEvent.responseText = this.xhr.getResponseText();
      this.dispatchEvent(dataEvent);
      this.pendingSend = window.setTimeout(goog.bind(this.send_, this), delay);
      break;
    default:
  }
};

bnCL.processSuccess = function(event) {
};

bnCL.disposeInternal = function() {
  this.xhr.dispose();
  ble.net.CometLoop.superClass_.disposeInternal.call(this);
};

////////////////////////////////////////////////////////////////////////////////
                                                                           });
////////////////////////////////////////////////////////////////////////////////
