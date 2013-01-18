goog.require('goog.ui.Component');

goog.require('ble.scribble.UI');

goog.require('ble.tpg.model.Game');

goog.provide('ble.tpg.ui.Scribbler');

goog.scope(function() {


var Component = goog.ui.Component;
/**
 * @constructor
 * @param {ble.tpg.model.Game} game
 * @extends {Component}
 */
ble.tpg.ui.Scribbler = function(game) {
  Component.call(this);
  this.game = game;
  this.scribble = new ble.scribble.UI(this.game.width, this.game.height);
  this.addChild(this.scribble, true);
};
goog.inherits(ble.tpg.ui.Scribbler, Component);


});
