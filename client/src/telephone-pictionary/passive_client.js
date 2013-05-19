
goog.provide('ble.telephone_pictionary.PassiveClient');

goog.require('ble.telephone_pictionary.GameUpdater');
goog.require('ble.telephone_pictionary.ChatEvent');
goog.require('goog.events.EventTarget');
goog.require('goog.events.Event');

goog.scope(function() {

var _ = ble.telephone_pictionary;
var EventTarget = goog.events.EventTarget;
var Event = goog.events.Event;
var console = window.console;
var JSON = window.JSON;

/**
 @constructor
 @param {_.GameUpdater} game
 @extends {EventTarget}
**/
_.PassiveClient = function(game) {
  EventTarget.call(this);
  this.game = game;
};
goog.inherits(_.PassiveClient, EventTarget);

_.PassiveClient.prototype.handleEvent = function(event) {
  var text = event['responseText'];
  if(text === undefined) {
    console.error('missing response text');
    return;
  }
  var obj = JSON.parse(text);
  var events = obj['events'];
  if(events === undefined) {
    console.error('missing events');
    return;
  }
  for(var i = 0; i < events.length; i++) {
    var gameEvent = events[i];
    var type = gameEvent['actionType'];
    switch(type) {
    case 'joinGame':
      console.log(gameEvent['who']);
      console.log(gameEvent['name']);
      this.game.joinGame(gameEvent['who'], gameEvent['name'], false);
      break;
    case 'chat':
      this.dispatchEvent(new _.ChatEvent(this, gameEvent['content']));
      break;
    case 'startGame':
      this.game.startGame(gameEvent['who']);
      break;
    case 'passStack':
      this.game.passStack(
          gameEvent['who'],
          gameEvent['toWhom'],
          gameEvent['stackId'],
          gameEvent['url']);
      break;
    }
    console.log(events[i]);
  }
};

});
