

goog.provide('ble.telephone_pictionary.ServerModel');

goog.require('goog.events.EventTarget');

goog.require('ble.net.QueryTimeComet');
goog.require('ble.net.EventType');

goog.require('ble.telephone_pictionary.Client');
goog.require('ble.telephone_pictionary.ClientImpl');
goog.require('ble.telephone_pictionary.GameImpl');
goog.require('ble.telephone_pictionary.PassiveClient');

goog.scope(function() {

var _ = ble.telephone_pictionary;
var QueryTimeComet = ble.net.QueryTimeComet;
var EventTarget = goog.events.EventTarget;
var NetEventType = ble.net.EventType;
var console = window.console;

/**
 * @constructor
 * @param {_.Client} client
 * @param {_.PassiveClient} passive
 * @param {_.GameImpl} game
 * @param {QueryTimeComet} comet
 * @extends {EventTarget}
 */
_.ServerModel = function(client, passive, game, comet) {
  EventTarget.call(this);
  this.client = client;
  this.passive = passive;
  this.game = game;
  this.comet = comet;
  this.wireUp_();
};
goog.inherits(_.ServerModel, EventTarget);

_.ServerModel.getUrlFromWindow = function(window) {
  var url = window.location.toString();
  url = url.replace(/\/[^\/]*$/, '/');
  return url;
};

_.ServerModel.initialize = function(window) {
  var url = _.ServerModel.getUrlFromWindow(window);
  var client = new _.ClientImpl(url);
  var game = new _.GameImpl(client);
  var passive = new _.PassiveClient(game);
  var comet = new QueryTimeComet(url + 'events');
  return new _.ServerModel(client, passive, game, comet);
};

/**
 * @private
 */
_.ServerModel.prototype.wireUp_ = function() {
  goog.events.listen(this.comet, NetEventType.COMET_DATA, this.passive);
  this.game.setParentEventTarget(this);
  this.passive.setParentEventTarget(this);
};

//`run` can be called more than once to re-sync to the game state
_.ServerModel.prototype.run = function() {
  this.comet.stop();
  var fetch = this.game.fetchState();
  fetch.wait(goog.bind(this.runFetch_, this));
};

/**
 * @private
 */
_.ServerModel.prototype.runFetch_ = function(fetchResult) {
  if(fetchResult.getError() !== undefined) {
    console.error(fetchResult.getError());
    return;
  }
  var value = fetchResult.getValue();
  var time = value['lastTime'];
  if(time === undefined) {
    console.error('didn\'t get a time');
    return
  }
  this.comet.runAt(time);
};

});
