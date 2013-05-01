
goog.provide('ble.telephone_pictionary.GameImpl');

goog.require('ble.telephone_pictionary.Game');
goog.require('ble.telephone_pictionary.GameUpdater');
goog.require('ble.telephone_pictionary.JoinEvent');
goog.require('ble.telephone_pictionary.PassEvent');

goog.require('goog.events.EventTarget');
goog.require('goog.result');
goog.require('goog.result.Result');

goog.require('ble._2d.DrawPart');
goog.require('ble.scribbleDeserializer');

goog.scope(function() {
var _ = ble.telephone_pictionary;
var isDef = goog.isDefAndNotNull;
var Result = goog.result.Result;
var result = goog.result;
var EventTarget = goog.events.EventTarget;
var DrawPart = ble._2d.DrawPart;
var console = window.console;

/**
 * @constructor
 * @extends {EventTarget}
 * @param {_.Client} client
 * @param {_.GameImpl.GameState=} state
 * @implements {_.Game}
 * @implements {_.GameUpdater}
 */
_.GameImpl= function(client, state) {
  this.client = client;
  this.state = isDef(state) ? state : new _.GameImpl.GameState();
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

_.GameImpl.prototype.joinGame = function(playerId, playerName, isMe) {
  var state = this.state;
  var player = new _.GameImpl.Player(this, playerId, playerName);
  if(player.id() in state.playersById) {
    console.error("duplicate player");
    return;
  }
  if(isMe && state.playerMe != null) {
    console.error("duplicate identification of player");
    return;
  }
  state.players.push(player);
  state.playersById[player.id()] = player;
  this.dispatchEvent(new _.JoinEvent(this, player));
};

_.GameImpl.prototype.passStack = function(from, to, stackId, stackUrl) {
  var state = this.state, playerFrom, playerTo, stack;
  if(isDef(from)) {

    if(from in state.playersById) {
      playerFrom = state.playersById[from];
    } else {
      console.error("pass from missing player");
      return;
    }

  } else {
    playerFrom = null;
  }
  if(isDef(to)) {

    if(to in state.playersById) {
      playerTo = state.playersById[to];
    } else {
      console.error("pass from missing player");
      return;
    }

  } else {
    playerTo = null;
  }

  //when passed from no one, means a new stack
  if(playerFrom === null) {
    //stack id must be unique
    if(stackId in state.stacksById) {
      console.error("duplicate stack id");
      return;
    }
    //must be passed to somebody
    if(playerTo === null) {
      console.error("pass from nobody, to nobody")
    }
    stack = new _.GameImpl.Stack(stackId, stackUrl);
    state.stacks.push(stack);
    state.stacksById[stack.id()] = stack;
  } else {
    //it must already exist
    if(!(stackId in state.stacksById)) {
      console.error("missing stack id");
      return;
    }
    stack = state.stacksById[stackId];

    //when coming from someone, we need to remove it from their stacks held
    var stacksHeld = state.stacksInPlay[playerFrom.id()];
    var indexToRemove = stacksHeld.indexOf(stack);
    if(indexToRemove === -1) {
      console.error("could not remove stack from passing player's held stacks");
      return;
    }
    stacksHeld.splice(indexToRemove, 1);
  }

  if(isDef(playerTo)) {
    if(!(playerTo.id() in state.stacksInPlay))
      state.stacksInPlay[playerTo.id()] = [];
    state.stacksInPlay[playerTo.id()].push(stack);
  }
  this.dispatchEvent(new _.PassEvent(this, playerFrom, playerTo, stack));

};

_.GameImpl.prototype.startGame = function(whoId) {
  var state = this.state;
  if(state.isStarted) {
    console.error("game is already started");
    return;
  }
  if(state.isComplete) {
    console.error("game is already complete");
    return;
  }
  if(!(whoId in state.playersById)) {
    console.error("unrecognized player started game");
    return;
  }
  state.started = true;
  this.dispatchEvent(new _.StartEvent(this, state.playersById[whoId]));
};

_.GameImpl.prototype.updateTime = function(time) {
  var state = this.state;
  if(state.lastTime > time) {
    console.error("time value decreased");
    return;
  }
  state.lastTime = time;
};






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

  /**@type {Object.<string, Array.<_.Stack>>}*/
  this.stacksInPlay = {};

  /**@type {Object.<string, _.Stack>} */
  this.stacksById = {};
  this.url = "";

  /**@type {?_.Player}*/
  this.playerMe = null;
};

_.GameImpl.GameState.prototype.setFromJSON = function(obj) {};

/** @constructor
 *  @implements {_.Player}
 *  @param {_.Game} game
 *  @param {string} id
 *  @param {string} name
 *  */
_.GameImpl.Player = function(game, id, name) {};

/** @return {Array.<_.Stack>} */
_.GameImpl.Player.prototype.stacksHeld = function(){};
/** @return {string} */
_.GameImpl.Player.prototype.id = function(){};
/** @return {string} */
_.GameImpl.Player.prototype.name = function(){};
/** @return {string} */
_.GameImpl.Player.prototype.styleClass = function(){};


/** @constructor
 *  @implements {_.Stack}
 *  @param {string} id
 *  @param {string} url */
_.GameImpl.Stack = function(id, url) {};
/** @return {Result} */
_.GameImpl.Stack.prototype.fetchState = function(){};
/** @return {?Array.<_.Drawing>} */
_.GameImpl.Stack.prototype.drawings = function(){};
/** @return {string} */
_.GameImpl.Stack.prototype.id = function(){};

/** @constructor */
_.GameImpl.Drawing = function() {};
/** @return {Result} */
_.GameImpl.Drawing.prototype.fetchState = function(){};
/** @return {?Array.<DrawPart>} */
_.GameImpl.Drawing.prototype.content = function(){};
/** @return {string} */
_.GameImpl.Drawing.prototype.id = function(){};
/** @return {_.Player} */
_.GameImpl.Drawing.prototype.player = function(){};
});
