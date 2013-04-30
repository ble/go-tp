
goog.provide('ble.telephone_pictionary.GameImpl');

goog.require('ble.telephone_pictionary.Game');

goog.require('goog.events.EventTarget');
goog.require('goog.result');
goog.require('goog.result.Result');

goog.require('ble.scribbleDeserializer');

goog.scope(function() {
var _ = ble.telephone_pictionary;
var Result = goog.result.Result;
var result = goog.result;
var EventTarget = goog.events.EventTarget;

/**
 * @constructor
 * @extends {EventTarget}
 * @param {_.Client} client
 * @param {_.GameState=} state
 * @implements {_.Game}
 */
_.GameImpl= function(client, state) {
  this.client = client; 
  this.state = goog.isDefAndNotNull(state) ? state : new _.GameState();
};
goog.inherits(_.GameImpl, EventTarget);

_.GameImpl.prototype.fetchState = function() {
  var requested = this.client.getGameState();
  requested.wait(goog.bind(this.finishFetchState, this));
  return result.transform(requested, function(result) { return result; });
};

_.GameImpl.prototype.isStarted = function(){ return this.state.isStarted; };

_.GameImpl.prototype.isFinished = function(){ return this.state.isComplete; };

_.GameImpl.prototype.players = function(){ return this.state.players.slice(); };

_.GameImpl.prototype.playersById = function(){
  return goog.object.clone(this.state.playersById);
};

_.GameImpl.prototype.stacks = function(){ return this.state.stacks.slice(); };

_.GameImpl.prototype.stacksByHoldingPlayerId = function(){
  return goog.object.clone(this.state.stacksInPlay);
};

/** @return {?_.Player} */
_.GameImpl.prototype.myPlayer = function(){ return this.state.playerMe; };

_.GameImpl.prototype.finishFetchState = function(jsonResult) {
  switch(jsonResult.getState()) {
    case Result.State.PENDING:
      throw "invalid state for call to finishFetchState";
    case Result.State.ERROR:
      throw jsonResult.getError();
    case Result.State.SUCCESS:
      return this.state.setFromJSON(jsonResult.getValue());
  }
};

/**
 * @constructor
 */ 
_.GameImpl.GameState = function() {
  this.id = "";
  this.isComplete = false;
  this.isStarted = false;
  this.lastTime = 0;
  /**@type {Array.<_.Player>}*/
  this.players = [];

  /**@type {Object.<string, _.Player>}*/
  this.playersById = {};

  /**@type {Array.<_.Stack>}*/
  this.stacks = [];

  /**@type {Object.<string, _.Stack>}*/
  this.stacksInPlay = {};
  this.url = "";

  /**@type {?_.Player}*/
  this.playerMe = null;
};

_.GameImpl.GameState.prototype.setFromJSON = function(obj) {}

/** @constructor */
_.GameImpl.Player = function() {}
/** @return {Array.<_.Stack>} */
_.GameImpl.Player.prototype.stacksHeld = function(){};
/** @return {string} */
_.GameImpl.Player.prototype.id = function(){};
/** @return {string} */
_.GameImpl.Player.prototype.name = function(){};
/** @return {string} */
_.GameImpl.Player.prototype.styleClass = function(){};


/** @constructor */
_.GameImpl.Stack = function() {}
/** @return {Result} */
_.GameImpl.Stack.prototype.fetchState = function(){};
/** @return {?Array.<_.Drawing>} */
_.GameImpl.Stack.prototype.drawings = function(){};
/** @return {string} */
_.GameImpl.Stack.prototype.id = function(){};
/** @return {Result} */
_.GameImpl.Drawing.prototype.fetchState = function(){};
/** @return {?Array.<DrawPart>} */
_.GameImpl.Drawing.prototype.content = function(){};
/** @return {string} */
_.GameImpl.Drawing.prototype.id = function(){};
/** @return {_.Player} */
_.GameImpl.Drawing.prototype.player = function(){};
});
