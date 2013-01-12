goog.require('goog.dom');
goog.require('goog.events');
goog.require('goog.ui.Component');
goog.require('ble.tpg.templates');
goog.require('ble.tpg.model.EventType');

goog.provide('ble.tpg.ui.ChatContainer');
goog.scope(function() {
//scope-start

var Component = goog.ui.Component;
var ModelType = ble.tpg.model.EventType;
var templates = ble.tpg.templates;

/**
 * @constructor
 * @extends{goog.ui.Component}
 * @param{ble.tpg.model.Game} game
 */
ble.tpg.ui.ChatContainer = function(game) {
  Component.call(this);
  this.game = game;
  this.chats = new ble.tpg.ui.Chats();
  this.chatInput = new ble.tpg.ui.ChatInput();
  this.addChild(chats);
  this.addChild(chatInput);
};
goog.inherits(ble.tpg.ui.ChatContainer, Component);

var ccp = ble.tpg.ui.ChatContainer.prototype;

/**
 * @param{Element} element
 * @return{boolean}
 */
ccp.canDecorate = function(element) {
  return element.tagName == 'div';
};

ccp.enterDocument = function() {
  goog.dom.removeChildren(this.getElement());
  Component.enterDocument.call(this);
  goog.events.listen(
      this.game,
      ModelType.ALL,
      this);
};

ccp.exitDocument = function() {
  Component.exitDocument.call(this);
  goog.events.unlisten(
      this.game,
      ModelType.ALL,
      this);
};

/**
 * @param{goog.events.Event} event
 */
ccp.handleEvent = function(event) {
  switch(event.type) {
    case ModelType.CHAT;
    case ModelType.PASS;
    case ModelType.START_GAME;
    case ModelType.COMPLETE_GAME;
    case ModelType.JOIN_GAME;
      this.displayEvent(event.fromServer);
    break;
  }
};

ccp.displayEvent = function(event) {
  var who = event['who'];
  var toWhom = event['toWhom'];
  var stackId = event['stackId'];
  var content = event['content'];

  switch(event['actionType']) {
    case "joinGame":
      displayJoin(who);
      break;
    case "chat":
      displayChat(who, content);
      break;
    case "passStack":
      displayPass(who, stackId, toWhom);
      break;
    case "startGame":
      displayStart(who);
      break;
    case "completeGame":
      break;
  }
};

ccp.displayJoin = function(playerId) {
  var dom = this.dom_;
  var player = this.game.players[playerId];
  var o = ({
    'name': player.name,
    'styleName': player.styleName});
  var line = dom.htmlToDocumentFragment(
      templates.joinLine(o));
  this.getElement().appendChild(line);
};

ccp.displayChat = function(playerId, content) {
  var player = this.game.players[playerId];
  var o = ({
    'name': player.name,
    'styleName': player.styleName,
    'content': content});
  var line = dom.htmlToDocumentFragment(
      templates.chatLine(o));
  this.getElement().appendChild(line);
};

ccp.displayPass = function(playerId, stackId, toWhom) {
  var player = this.game.players[playerId];
  var o = ({
    'name': player.name,
    'styleName': player.styleName});
  if(toWhom) {
    var recipient = this.game.players[toWhom];
    o['nameRecipient'] = recipient.name;
    o['styleRecipient'] = recipient.styleName;
  }
  var line = dom.htmlToDocumentFragment(
      templates.passLine(o));
  this.getElement().appendChild(line);
};

ccp.displayStart = function(playerId) {
  var player = this.game.players[playerId];
  var o = ({
    'name': player.name,
    'styleName': player.styleName});
  var line = dom.htmlToDocumentFragment(
      templates.startLine(o));
  this.getElement().appendChild(line);
};

/**
 * @constructor
 * @extends{goog.ui.Component}
 */
ble.tpg.ui.Chats = function() {
  Component.call(this);
};
goog.inherits(ble.tpg.ui.Chats, Component);

//var cp = 
//ble.tpg.ui.
//scope-end
});
