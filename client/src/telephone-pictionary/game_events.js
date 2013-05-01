goog.provide('ble.telephone_pictionary.EventType');

goog.provide('ble.telephone_pictionary.JoinEvent');
goog.provide('ble.telephone_pictionary.PassEvent');
goog.provide('ble.telephone_pictionary.StartEvent');

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
  JOIN: 'JOIN',
  PASS: 'PASS',
  START: 'START'
});

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

});
