goog.require('goog.object');
goog.require('goog.events.EventTarget');
goog.require('goog.events.Event');
goog.require('goog.labs.net.xhr');

goog.require('ble.scribble.Drawing');

goog.provide('ble.tpg.model.EventType');
goog.provide('ble.tpg.model.Player');
goog.provide('ble.tpg.model.Drawing');
goog.provide('ble.tpg.model.Stack');
goog.provide('ble.tpg.model.Game');

goog.scope(function() {
//scope-start

var xhr = goog.labs.net.xhr;

/**
 * @enum{string}
 */
ble.tpg.model.EventType = ({
  CHAT: 'tpg-chat',
  PASS: 'tpg-pass',
  START_GAME: 'tpg-start-game',
  COMPLETE_GAME: 'tpg-complete-game',
  JOIN_GAME: 'tpg-join-game'});
var EventType = ble.tpg.model.EventType;
EventType.ALL = [
  EventType.CHAT,
  EventType.PASS,
  EventType.START_GAME,
  EventType.COMPLETE_GAME,
  EventType.JOIN_GAME];

/**
 * @constructor
 * @param{string} id
 * @param{string} name
 * @param{boolean=} isYou
 */
ble.tpg.model.Player = function(id, name, isYou) {
  this.id = id;
  this.name = name;
  this.isYou = isYou ? true : false;
  this.styleName = '';
};
var Player = ble.tpg.model.Player;

Player.fromJSON = function(o) {
  return new Player(o['id'], o['pseudonym'], o['isYou']);
};

Player.arrayFromJSON = function(o) {
  var result = [];
  for(var i = 0; i < o.length; i++) {
    var player = Player.fromJSON(o[i]);
    player.styleName = 'player-' + i.toString();
    result.push(player);
  }
  return result;
};

Player.newForArray = function(array, id, name, isYou) {
  var player = new Player(id, name, isYou);
  player.styleName = 'player-' + array.length.toString();
  return player;
};

/**
 * @constructor
 * @param{string} id
 * @param{string?} content
 * @param{string} url
 */
ble.tpg.model.Drawing = function(id, url, content) {
  this.id = id;
  this.content = content;
  this.url = url;
}
var Drawing = ble.tpg.model.Drawing;

Drawing.fromJSON = function(o) {
  return new Drawing(o['id'], null, o['url']);
};

Drawing.arrayFromJSON = function(o) {
  var result = [];
  for(var i = 0; i < o.length; i++) {
    var drawing = Drawing.fromJSON(o[i]);
    result.push(drawing);
  }
  return result;
};

/**
 * @constructor
 * @param{string} id
 * @param{string} url
 * @param{Array.<Drawing>} drawings
 */
ble.tpg.model.Stack = function(id, url, drawings) {
  this.id = id;
  this.url = url;
  this.drawings = drawings;
};
var Stack = ble.tpg.model.Stack;

Stack.fromJSON = function(o) {
  return new Stack(o['id'], o['url'], Drawing.arrayFromJSON(o['drawings']));
};

Stack.arrayFromJSON = function(o) {
  var result = [];
  for(var i = 0; i < o.length; i++) {
    var stack = Stack.fromJSON(o[i]);
    result.push(stack);
  }
  return result;
}
/**
 * @constructor
 * @extends{goog.events.EventTarget}
 */
ble.tpg.model.Game = function(
    id, lastTime, players, stacks,
    inPlay, url, isStarted, isComplete) {
  this.id = id;
  this.lastTime = lastTime;
  this.url = url;

  this.isStarted = isStarted;
  this.isComplete = isComplete;

  this.playerMe = null;
  this.players = [];
  this.playersById = {};

  for(var i = 0; i < players.length; i++) {
    this.addPlayer(players[i]);
  }

  this.stacks = [];
  this.stacksById = {};
  this.stacksByHolderId = {};

  this.addStacks(stacks, inPlay);
};
goog.inherits(ble.tpg.model.Game, goog.events.EventTarget);
var Game = ble.tpg.model.Game;

Game.prototype.addStacks = function(stacks, inPlay) {
  for(var i = 0; i < stacks.length; i++) {
    var stack = stacks[i];
    this.stacks.push(stack);
    this.stacksById[stack.id] = stack;
  }
  goog.object.forEach(
      inPlay,
      function(stackIds, holderId) {
        var held = [];
        this.stacksByHolderId[holderId] = held;
        for(var i = 0; i < stackIds.length; i++) {
          held.push(this.stacksById[stackIds[i]]);
        }
      },
      this);
};

Game.prototype.width = 360;
Game.prototype.height = 270;

Game.fromJSON = function(o) {
  return new Game(
      o['id'],
      o['lastTime'],
      Player.arrayFromJSON(o['players']),
      Stack.arrayFromJSON(o['stacks']),
      o['stacksInPlay'],
      o['url'],
      o['isStarted'],
      o['isComplete']);
};

Game.prototype.getMyPlayer = function() {
  return this.playerMe;
};

Game.prototype.getMyStacks = function() {
  var me = this.playerMe; 
  return this.stacksByHolderId[me.id];
}

var cometType = ble.net.EventType;
var JSON = window.JSON;
Game.prototype.handleEvent = function(event) {
  if(event.type == cometType.COMET_DATA) {
    var json = JSON.parse(event.responseText);
    var receivedEvents = json['events'];
    if(goog.isDefAndNotNull(receivedEvents)) {
      for(var i = 0; i < receivedEvents.length; i++) {
        this.processJsonEvent(receivedEvents[i]);
      }
    }
  }
};

Game.prototype.addPlayer = function(player) {
  if(player.id in this.playersById)
    throw new Error('duplicate player id');
  var newPlayer = Player.newForArray(this.players, player.id, player.name, player.isYou);
  this.players.push(newPlayer);
  this.playersById[newPlayer.id] = newPlayer;
  if(newPlayer.isYou)
    this.playerMe = newPlayer;
  return newPlayer;
};

Game.prototype.passStack = function(pFrom, pTo, stackId, url) {
  var theStack;
  //if this stack is not being passed from someone,
  //it's just been created at the start of the game.
  if(!goog.isDefAndNotNull(pFrom)) {
    theStack = new Stack(stackId, url, null);
    this.stacks.push(theStack);
    this.stacksByHolderId[pTo.id] = [];
    this.stacksById[theStack.id] = theStack;
  } else {
    theStack = this.stacksById[stackId];
  }

  this.stacksByHolderId[pTo.id].push(theStack); 

  if(goog.isDefAndNotNull(pFrom)) {
    var heldByPasser = this.stacksByHolderId[pFrom];
    if(heldByPasser[0] !== theStack)
      throw new Error("expected first stack to be passed");
    heldByPasser.shift();
  }
};

var Event = goog.events.Event;
Game.prototype.processJsonEvent = function(o) {
  var event;
  switch(o['actionType']) {
    case 'joinGame':
      
      var newPlayer = this.addPlayer(new Player(o['who'], o['name'], o['isYou']));
      event = new Event(EventType.JOIN_GAME, this);
      event.player = newPlayer;
      this.dispatchEvent(event);
      break;
    case 'chat':
      var player = this.playersById[o['who']];
      event = new Event(EventType.CHAT, this);
      event.player = player;
      event.content = o['content'];
      this.dispatchEvent(event);
      break;
    case 'passStack':
      var playerFrom = this.playersById[o['who']];
      var playerTo = this.playersById[o['toWhom']];
      var stackId = o['stackId'];
      var url = o['url'];
      this.passStack(playerFrom, playerTo, stackId, url);
      event = new Event(EventType.PASS);
      event.from = playerFrom;
      event.to = playerTo;
      event.stack = this.stacksById[stackId];
      this.dispatchEvent(event);
      break;
  }
}
//scope-end
});
