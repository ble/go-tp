goog.provide('ble.telephone_pictionary.Client');

goog.require('goog.result');
goog.require('goog.result.Result');
goog.require('ble._2d.DrawPart');

goog.scope(function() {
var _ = ble.telephone_pictionary;
var Result = goog.result.Result;
var result = goog.result;

var DrawPart = ble._2d.DrawPart;
/**
 * @interface
 */
_.Client = function(){};

/** @return Result.<object> */
_.Client.prototype.getGameState = function(){};

/**
 * @param {string} stackId
 * @return Result.<object>
 */
_.Client.prototype.getStack = function(stackId){};

/**
 * @param {string} drawingId
 * @return Result.<object>
 */
_.Client.prototype.getDrawing = function(drawingId){};

/**
 * @param {DrawPart} part
 * @return Result.<*>
 */ 
_.Client.prototype.appendToDrawing = function(part){};

/**
 * @param {string} stackId
 * @return Result.<*>
 */
_.Client.prototype.passStack = function(stackId){};

/**
 * @param {string} message
 * @return Result.<*>
 */
_.Client.prototype.chat = function(message){};

});
