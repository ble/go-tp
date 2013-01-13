goog.require('goog.events.EventTarget');
goog.require('ble.scribble.Drawing');


goog.provide('ble.tpg.model.EventType');
goog.provide('ble.tpg.model.Player');
goog.provide('ble.tpg.model.Drawing');
goog.provide('ble.tpg.model.Stack');
goog.provide('ble.tpg.model.Game');

goog.scope(function() {
//scope-start

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

Game.prototype.handleEvent = function(event) {

};
//scope-end
});
