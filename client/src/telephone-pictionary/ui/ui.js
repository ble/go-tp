
goog.provide('ble.telephone_pictionary.UI');

goog.require('ble.telephone_pictionary.Roster');
goog.require('ble.telephone_pictionary.Chatroom');
goog.require('ble.telephone_pictionary.StateOfPlay');
goog.require('ble.telephone_pictionary.TaskDescription');
goog.require('ble.telephone_pictionary.DrawingToInterpret');
goog.require('ble.telephone_pictionary.DrawingInProgress');
goog.require('ble.telephone_pictionary.StackToReview');
goog.require('ble.telephone_pictionary.AllStacks');



goog.require('goog.ui.Component');


goog.scope(function() {


var _ = ble.telephone_pictionary;
var Component = goog.ui.Component;

/**
 * @constructor
 * @extends {Component}
 */
_.UI = function() {
  Component.call(this);
  this.tpgElements_ = {};
  for(var className in this.childrenCssClasses) {
    var child = new this.childrenCssClasses[className];
    this.tpgElements_[className] = child;
    this.addChild(child, true); 
    //really stupid placeholder
    child.getElement().appendChild(document.createTextNode(className));
  };
};
goog.inherits(_.UI, Component);

_.UI.prototype.decorateInternal = function(element) {
  goog.base(this, 'decorateInternal', element);
  for(var className in this.childrenCssClasses) {
    var child = this.tpgElements_[className];
    var withName = this.getElement().getElementsByClassName(className);
    if(withName.length != 1) {
      continue;
      throw "oh no";
    }
    var node = withName[0];
    var childElement = child.getElement();
    node.parentNode.replaceChild(childElement, node);
    for(var childContent in node.childNodes) {
      childElement.appendChild(node.childNodes[childContent]);
    }
    child.getElement().classList.add(className);
  } 
}
_.UI.prototype.enterDocument = function() {
  goog.base(this, 'enterDocument');
};

_.UI.prototype.canDecorate = function(element) {
  var allPresent = true;
  for(var className in this.childrenCssClasses) {
    var withName = element.getElementsByClassName(className);
    if(withName.length == 0) { 
      window.console.error("missing " + className);
      continue;
      allPresent = false;
    }
    if(withName.length > 1) {
      allPresent = false;
      window.console.error("multiple " + className);
    }
  }
  return allPresent;
};

_.UI.prototype.childrenCssClasses = {
  'room-title': _.RoomTitle,
  'roster': _.Roster,
  'chatroom': _.Chatroom,
  'state-of-play': _.StateOfPlay,
  'task-description': _.TaskDescription,
  'drawing-to-interpret': _.DrawingToInterpret,
  'drawing-in-progress': _.DrawingInProgress,
  'stack-to-review': _.StackToReview,
  'all-stacks': _.AllStacks
};


});
