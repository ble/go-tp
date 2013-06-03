goog.provide('ble.telephone_pictionary.Game');
goog.provide('ble.telephone_pictionary.Player');
goog.provide('ble.telephone_pictionary.Stack');
goog.provide('ble.telephone_pictionary.Drawing');


goog.require('goog.result.Result');

goog.require('ble._2d.DrawPart');

goog.scope(function() {
var _ = ble.telephone_pictionary;
var Result = goog.result.Result;
var result = goog.result;

var DrawPart = ble._2d.DrawPart;

/** @interface */
_.Game = function() {};

/** @interface */
_.Player = function() {};

/** @interface */
_.Stack = function() {};

/** @interface */
_.Drawing = function() {};

/** @return {Result} */
_.Game.prototype.fetchState = function(){};

/** @return {boolean} */
_.Game.prototype.isStarted = function(){};

/** @return {boolean} */
_.Game.prototype.isFinished = function(){};

/** @return {Array.<_.Player>} */
_.Game.prototype.players = function(){};

/** @return {Object.<string, _.Player>} */
_.Game.prototype.playersById = function(){};

/** @return {Array.<_.Stack>} */
_.Game.prototype.stacks = function(){};

/** @return {Object.<string, _.Stack>} */
_.Game.prototype.stacksByHoldingPlayerId = function(){};

/** @return {?_.Player} */
_.Game.prototype.myPlayer = function(){};

/** @return {Array.<_.Stack>} */
_.Player.prototype.stacksHeld = function(){};

/** @return {string} */
_.Player.prototype.id = function(){};

/** @return {string} */
_.Player.prototype.name = function(){};

/** @return {boolean} */
_.Player.prototype.isMe = function(){};

/** @return {string} */
_.Player.prototype.styleClass = function(){};

/** @return {Result} */
_.Stack.prototype.fetchState = function(){};

/** @return {?Array.<_.Drawing>} */
_.Stack.prototype.drawings = function(){};

/** @return {string} */
_.Stack.prototype.id = function(){};

/** @return {Result} */
_.Drawing.prototype.fetchState = function(){};

/** @return {?Array.<DrawPart>} */
_.Drawing.prototype.content = function(){};

/** @return {string} */
_.Drawing.prototype.id = function(){};

/** @return {_.Player} */
_.Drawing.prototype.player = function(){};
});
