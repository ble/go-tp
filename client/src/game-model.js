goog.require('goog.events.EventTarget');
goog.require('ble.scribble.Drawing');
goog.require('goog.net.EventType');
goog.require('goog.events.Event');

goog.provide('ble.tpg.model.EventType');
goog.provide('ble.tpg.model.Player');
goog.provide('ble.tpg.model.Drawing');
goog.provide('ble.tpg.model.Stack');
goog.provide('ble.tpg.model.Game');

goog.scope(function() {
//scope-start

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
 */
ble.tpg.model.Player = function(id, name) {
  this.id = id;
  this.name = name;
  this.styleName = '';
};
var Player = ble.tpg.model.Player;

Player.fromJSON = function(o) {
  return new Player(o['id'], o['pseudonym']);
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

Player.newForArray = function(array, id, name) {
  var player = new Player(id, name);
  player.styleName = 'player-' + array.length.toString();
  return player;
};

/**
 * @constructor
 * @param{string} id
 * @param{string?} content
 * @param{string} url
 */
ble.tpg.model.Drawing = function(id, content, url) {
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
ble.tpg.model.Game = function(id, lastTime, players, stacks, inPlay, url) {
  this.id = id;
  this.lastTime = lastTime;
  this.players = players;
  this.stacks = stacks;
  this.inPlay = inPlay;
  this.url = url;
};
goog.inherits(ble.tpg.model.Game, goog.events.EventTarget);
var Game = ble.tpg.model.Game;

Game.fromJSON = function(o) {
  return new Game(
      o['id'],
      o['lastTime'],
      Player.arrayFromJSON(o['players']),
      Stack.arrayFromJSON(o['stacks']),
      o['stacksInPlay'],
      o['url']);
};

var cometType = ble.net.EventType;
var console = window.console;
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
    console.log(json);
  }
};

var Event = goog.events.Event;
Game.prototype.processJsonEvent = function(o) {
  var event;
  switch(o['actionType']) {
    case 'joinGame':
      var newPlayer = Player.newForArray(this.players, o['who'], o['name']);
      this.players.push(newPlayer);
      event = new Event(EventType.JOIN_GAME, this);
      event.player = newPlayer; 
      this.dispatchEvent(event);
  }
}
//scope-end
});
