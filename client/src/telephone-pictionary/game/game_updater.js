goog.provide('ble.telephone_pictionary.GameUpdater');

goog.scope(function() {
var _ = ble.telephone_pictionary;


/** @interface */
_.GameUpdater = function() {};

/** @param {string} playerId
 *  @param {string} playerName
 *  @param {boolean} isMe */
_.GameUpdater.prototype.joinGame = function(playerId, playerName, isMe) {};

/** @param {?string} from
 *  @param {?string} to
 *  @param {string} stackId
 *  @param {string} stackUrl */
_.GameUpdater.prototype.passStack = function(from, to, stackId, stackUrl) {};

/** @param {string} whoId */
_.GameUpdater.prototype.startGame = function(whoId) {};

/** @param {number} time */
_.GameUpdater.prototype.updateTime = function(time) {};
});
