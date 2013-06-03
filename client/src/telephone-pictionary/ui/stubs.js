goog.provide('ble.telephone_pictionary.RoomTitle');
goog.provide('ble.telephone_pictionary.Roster');
goog.provide('ble.telephone_pictionary.Chatroom');
goog.provide('ble.telephone_pictionary.StateOfPlay');
goog.provide('ble.telephone_pictionary.TaskDescription');
goog.provide('ble.telephone_pictionary.DrawingToInterpret');
goog.provide('ble.telephone_pictionary.DrawingInProgress');
goog.provide('ble.telephone_pictionary.StackToReview');
goog.provide('ble.telephone_pictionary.AllStacks');

goog.require('goog.ui.Component');

goog.scope(function() {


var _ = ble.telephone_pictionary;
var Component = goog.ui.Component;

/** @constructor
 *  @extends {Component} */
_.RoomTitle = function() { Component.call(this); }
goog.inherits(_.RoomTitle, Component);


/** @constructor
 *  @extends {Component} */
_.Roster = function() { Component.call(this); }
goog.inherits(_.Roster, Component);


/** @constructor
 *  @extends {Component} */
_.Chatroom = function() { Component.call(this); }
goog.inherits(_.Chatroom, Component);


/** @constructor
 *  @extends {Component} */
_.StateOfPlay = function() { Component.call(this); }
goog.inherits(_.StateOfPlay, Component);

/** @constructor
 *  @extends {Component} */
_.TaskDescription = function() { Component.call(this); }
goog.inherits(_.TaskDescription, Component);

/** @constructor
 *  @extends {Component} */
_.DrawingToInterpret = function() { Component.call(this); }
goog.inherits(_.DrawingToInterpret, Component);

/** @constructor
 *  @extends {Component} */
_.DrawingInProgress = function() { Component.call(this); }
goog.inherits(_.DrawingInProgress, Component);

/** @constructor
 *  @extends {Component} */
_.StackToReview = function() { Component.call(this); }
goog.inherits(_.StackToReview, Component);

/** @constructor
 *  @extends {Component} */
_.AllStacks = function() { Component.call(this); }
goog.inherits(_.AllStacks, Component);

/*
---->+ Pass drawing button
---->+ Start game button
*/


});
