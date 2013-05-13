
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
var isDefNotNull = goog.isDefAndNotNull;
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
  EventTarget.call(this);
  this.client = client;
  this.state = isDefNotNull(state) ? state : new _.GameImpl.GameState();
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
  var addedPlayer = state.addPlayer(playerId, playerName, isMe);
  if(isDefNotNull(addedPlayer)) {
    this.dispatchEvent(new _.JoinEvent(this, addedPlayer));
  }
};

_.GameImpl.prototype.passStack = function(from, to, stackId, stackUrl) {
  var state = this.state,
      stack;
  if(isDefNotNull(from)) {
    stack = state.takeStackFrom(stackId, from);
    if(!isDefNotNull(from)) {
      return;
    }
  } else {
    stack = state.createStack(stackId, stackUrl);
  }

  if(isDefNotNull(to)) {
    stack = state.giveStackTo(stackId, to);
    if(!isDefNotNull(stack)) {
      return;
    }
  }
  this.dispatchEvent(
    new _.PassEvent(
      this,
      state.player(from),
      state.player(to),
      stack));
};

_.GameImpl.prototype.startGame = function(whoId) {
  if(this.state.startGame())
    this.dispatchEvent(new _.StartEvent(this, this.state.playersById[whoId]));
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

_.GameImpl.GameState.prototype.player = function(id) {
  return this.playersById[id];
};

_.GameImpl.GameState.prototype.setFromJSON = function(obj) {
  var id           = obj['id'],
      isComplete   = obj['isComplete'],
      isStarted    = obj['isStarted'],
      lastTime     = obj['lastTime'],
      players      = obj['players'],
      stacks       = obj['stacks'],
      stacksInPlay = obj['stacksInPlay'],
      url          = obj['url'];
      lastTime     = obj['lastTime'];
  if(!goog.array.every(
      [id, isComplete, isStarted, lastTime, players, stacks, stacksInPlay, url],
      function(elem, index, array) { return isDefNotNull(elem); })) {
    console.error("missing field in JSON");
    return;
  }
  this.id = id;
  this.lastTime = lastTime;

  this.isComplete = Boolean(isComplete);
  this.isStarted = Boolean(isStarted);
  this.lastTime = Math.floor(Number(lastTime));
  this.url = String(url);

  this.players = [];
  this.playersById = {};
  this.stacks = [];
  this.stacksById = {};
  this.stacksInPlay = {}

  for(var i = 0; i < stacks.length; i++) {
    var sId = stacks[i]['id'];
    var sUrl = stacks[i]['url'];
    if(!isDefNotNull(this.createStack(sId, sUrl))) {
      return;
    }
  }

  for(var i = 0; i < players.length; i++) {
    var pId = players[i]['id'];
    var pName = players[i]['pseudonym'];
    var pIsMe = Boolean(players[i]['isYou']);
    if(!isDefNotNull(this.addPlayer(pId, pName, pIsMe))) {
      return;
    }
  }

  for(var playerId in stacksInPlay) {
    var held = stacksInPlay[playerId];
    for(var i = 0; i < held.length; i++) {
      if(!isDefNotNull(this.giveStackTo(held[i], playerId))) {
        return;
      }
    }
  }
};

_.GameImpl.GameState.prototype.addPlayer = function(id, name, isMe) {
  if(id in this.playersById) {
    console.error("duplicate player");
    return null;
  }
  if(isMe && this.playerMe != null) {
    console.error("duplicate identification of using player");
    return null;
  }
  var player = new _.GameImpl.Player(this, id, name, isMe);

  this.players.push(player);
  this.playersById[player.id()] = player;
  this.stacksInPlay[player.id()] = [];
  if(player.isMe())
    this.playerMe = player;
  return player;
};

_.GameImpl.GameState.prototype.createStack = function(id, url) {
  if(id in this.stacksById) {
    console.error("duplicate stack id");
    return null;
  }
  var stack = new _.GameImpl.Stack(id, url)

  this.stacks.push(stack);
  this.stacksById[stack.id()] = stack;
  return stack;
};

_.GameImpl.GameState.prototype.takeStackFrom = function(stackId, holderId) {
  if(!(stackId in this.stacksById)) {
    console.error("no such stack id");
    return null;
  }
  if(!(holderId in this.playersById)) {
    console.error("no such player id");
    return null;
  }
  var held = this.stacksInPlay[holderId],
      stack = this.stacksById[stackId],
      ix = held.indexOf(stack);
  if(ix < 0) {
    console.error("stack not held by player");
    return null;
  }
  held.splice(ix);
  return stack;
};

_.GameImpl.GameState.prototype.giveStackTo = function(stackId, receiverId) {
  if(!(stackId in this.stacksById)) {
    console.error("no such stack id");
    return null;
  }
  if(!(receiverId in this.playersById)) {
    console.error("no such player id");
    return null;
  }
  var held = this.stacksInPlay[receiverId],
      stack = this.stacksById[stackId],
      ix = held.indexOf(stack);
  if(ix >= 0) {
    console.error("stack already held by player");
    return null;
  }
  held.push(stack);
  return stack;
};

_.GameImpl.GameState.prototype.startGame = function() {
  if(this.isStarted) {
    console.error("game is already started");
    return false;
  }
  if(this.isComplete) {
    console.error("game is already complete");
    return false;
  }
 this.started = true;
  return true;
};

/** @constructor
 *  @implements {_.Player}
 *  @param {_.GameImpl.GameState} state
 *  @param {string} id
 *  @param {string} name
 *  @param {boolean} isMe
 *  */
_.GameImpl.Player = function(state, id, name, isMe) {
  this.state = state;
  this._id = id;
  this._name = name;
  this._isMe = isMe;
};

/** @return {Array.<_.Stack>} */
_.GameImpl.Player.prototype.stacksHeld = function() {
  return this.state.stacksInPlay[this.id()];
};

/** @return {string} */
_.GameImpl.Player.prototype.id = function() {
  return this._id;
};

/** @return {string} */
_.GameImpl.Player.prototype.name = function() {
  return this._name;
};

/** @return {string} */
_.GameImpl.Player.prototype.styleClass = function() {
  return "";
};

/** @return {boolean} */
_.GameImpl.Player.prototype.isMe = function() {
  return this._isMe;
};


/** @constructor
 *  @implements {_.Stack}
 *  @param {string} id
 *  @param {string} url */
_.GameImpl.Stack = function(id, url) {
  this._id = id;
  this._url = url;
};

/** @return {Result} */
_.GameImpl.Stack.prototype.fetchState = function() {
  throw "unimplemented";
};

/** @return {?Array.<_.Drawing>} */
_.GameImpl.Stack.prototype.drawings = function() {
  throw "unimplemented";
};

/** @return {string} */
_.GameImpl.Stack.prototype.id = function() {
  return this._id;
};

/** @constructor */
_.GameImpl.Drawing = function(id, player, stack) {
  throw "unimplemented";
};

/** @return {Result} */
_.GameImpl.Drawing.prototype.fetchState = function() {
  throw "unimplemented";
};

/** @return {?Array.<DrawPart>} */
_.GameImpl.Drawing.prototype.content = function() {
  throw "unimplemented";
};
/** @return {string} */
_.GameImpl.Drawing.prototype.id = function() {
  throw "unimplemented";
};

/** @return {_.Player} */
_.GameImpl.Drawing.prototype.player = function() {
  throw "unimplemented";
};
});
