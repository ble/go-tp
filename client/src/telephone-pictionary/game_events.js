goog.provide('ble.telephone_pictionary.EventType');

goog.provide('ble.telephone_pictionary.LoadedEvent');
goog.provide('ble.telephone_pictionary.JoinEvent');
goog.provide('ble.telephone_pictionary.PassEvent');
goog.provide('ble.telephone_pictionary.StartEvent');
goog.provide('ble.telephone_pictionary.ChatEvent');

goog.require('goog.events.Event');
goog.require('goog.events.Event');
goog.require('ble.telephone_pictionary.Game');
goog.require('ble.telephone_pictionary.Player');
goog.require('ble.telephone_pictionary.Stack');

goog.scope(function() {
var _ = ble.telephone_pictionary;
var Event = goog.events.Event;
var EventTarget = goog.events.EventTarget;

/** @enum {string} */
_.EventType = ({
  LOADED: 'LOADED',
  JOIN: 'JOIN',
  PASS: 'PASS',
  START: 'START',
  CHAT: 'CHAT'
});

/** @constructor
 *  @param {EventTarget} target
 *  @param {_.Game} game
 *  @extends {Event} */
_.LoadedEvent = function(target, game) {
  Event.call(this, _.EventType.LOADED, target);
  this.game = game;
};
goog.inherits(_.LoadedEvent, Event);

/** @constructor
 *  @param {EventTarget} game
 *  @param {_.Player} player
 *  @extends {Event} */
_.JoinEvent = function(game, player) {
  Event.call(this, _.EventType.JOIN, game);
  this.player = player;
};
goog.inherits(_.JoinEvent, Event);

/** @constructor
 *  @param {EventTarget} game
 *  @param {_.Player?} from
 *  @param {_.Player?} to
 *  @param {_.Stack} stack
 *  @extends {Event} */
_.PassEvent = function(game, from, to, stack) {
  Event.call(this, _.EventType.PASS, game);
  this.from = from;
  this.to = to;
  this.stack = stack;
};
goog.inherits(_.PassEvent, Event);

/** @constructor
 *  @param {EventTarget} game
 *  @param {_.Player} who
 *  @extends {Event} */
_.StartEvent = function(game, who) {
  Event.call(this, _.EventType.START, game);
  this.who = who;
};
goog.inherits(_.StartEvent, Event);


/** @constructor
 *  @param {EventTarget} target
 *  @param {string} chatContent
 *  @extends {Event} */
_.ChatEvent = function(target, chatContent) {
  Event.call(this, _.EventType.CHAT, target);
  this.content = chatContent;
};
goog.inherits(_.ChatEvent, Event);

});
