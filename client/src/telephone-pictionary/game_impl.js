
goog.provide('ble.telephone_pictionary.game_impl.Game');

goog.require('ble.telephone_pictionary.Game');

goog.require('goog.events.EventTarget');
goog.require('goog.result');
goog.require('goog.result.Result');

goog.require('ble.scribbleDeserializer');

goog.scope(function() {
var _ = ble.telephone_pictionary;
var __ = ble.telephone_pictionary.game_impl;
var Result = goog.result.Result;
var EventTarget = goog.events.EventTarget;

/**
 * @constructor
 * @extends {EventTarget}
 * @param {_.Client} client
 * @implements {_.Game}
 */
__.Game = function(client) {
  this.client = client;
  this.started = false;
  this.finished = false;
  this.aPlayers = [];
  this.oPlayersById = {};
  this.aStacks = [];
  this.oStacksHeldByPlayers = {};
  this.playerMe = null;
};
goog.inherits(_.game_impl.Game, EventTarget);

__.Game.prototype.fetchState = function() {
  var requested = this.client.getGameState();
  requested.wait(goog.bind(this.finishFetchState, this));
  return requested.transform(function(result) { return "succeeded."; });
};

__.Game.prototype.isStarted = function(){ return this.started; };

__.Game.prototype.isFinished = function(){ return this.finished; };

__.Game.prototype.players = function(){ return this.aPlayers.slice(); };

__.Game.prototype.playersById = function(){
  return goog.object.clone(this.oPlayersById);
};

__.Game.prototype.stacks = function(){ return this.aStacks.slice(); };

__.Game.prototype.stacksByHoldingPlayerId = function(){
  return goog.object.clone(this.oStacksHeldByPlayers);
};

/** @return {?_.Player} */
__.Game.prototype.myPlayer = function(){ return this.playerMe; };

__.Game.prototype.finishFetchState = function(result) {
  throw "UNIMPLEMENTED";
};

///** @return {Array.<_.Stack>} */
//_.Player.prototype.stacksHeld = function(){};
//
///** @return {string} */
//_.Player.prototype.id = function(){};
//
///** @return {string} */
//_.Player.prototype.name = function(){};
//
///** @return {string} */
//_.Player.prototype.styleClass = function(){};
//
///** @return {Result} */
//_.Stack.prototype.fetchState = function(){};
//
///** @return {?Array.<_.Drawing>} */
//_.Stack.prototype.drawings = function(){};
//
///** @return {string} */
//_.Stack.prototype.id = function(){};
//
///** @return {Result} */
//_.Drawing.prototype.fetchState = function(){};
//
///** @return {?Array.<DrawPart>} */
//_.Drawing.prototype.content = function(){};
//
///** @return {string} */
//_.Drawing.prototype.id = function(){};
//
///** @return {_.Player} */
//_.Drawing.prototype.player = function(){};
//
});
