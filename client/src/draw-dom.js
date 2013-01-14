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
  this.chats = new ble.tpg.ui.Chats(game);
//  this.chatInput = new ble.tpg.ui.ChatInput();
  this.addChild(this.chats, true);
//  this.addChild(chatInput);
};
goog.inherits(ble.tpg.ui.ChatContainer, Component);

var ccp = ble.tpg.ui.ChatContainer.prototype;

/**
 * @param{Element} element
 * @return{boolean}
 */
ccp.canDecorate = function(element) {
  return element.tagName.toLowerCase() == 'div';
};

ccp.enterDocument = function() {
  goog.base(this, 'enterDocument');
};

ccp.exitDocument = function() {
  goog.base(this, 'exitDocument');
  goog.events.unlisten(
      this.game,
      ModelType.ALL,
      this);
};

/**
 * @constructor
 * @param{ble.tpg.model.Game} game
 * @extends{goog.ui.Component}
 */
ble.tpg.ui.Chats = function(game) {
  Component.call(this);
  this.game = game
};
goog.inherits(ble.tpg.ui.Chats, Component);

var cp = ble.tpg.ui.Chats.prototype;

cp.createDom_ = function() {
  goog.base(this, 'createDom_');
  this.getElement().className = 'chats';
};

cp.enterDocument = function() {
  goog.base(this, 'enterDocument');
  goog.events.listen(
      this.game,
      ModelType.ALL,
      this); 
};

cp.exitDocument = function() {
  goog.base(this, 'exitDocument');
  goog.events.unlisten(
      this.game,
      ModelType.ALL,
      this);
}

/**
 * @param{goog.events.Event} event
 */
cp.handleEvent = function(event) {
  switch(event.type) {
    case ModelType.CHAT:
      break;
    case ModelType.PASS:
      break;
    case ModelType.START_GAME:
      break;
    case ModelType.COMPLETE_GAME:
      break;
    case ModelType.JOIN_GAME:
      this.displayJoin(event.player);
      break;
  }
};

/**
 * @param{ble.tpg.model.Player} player
 */
cp.displayJoin = function(player) {
  var dom = this.dom_;
  var o = ({
    'name': player.name,
    'styleName': player.styleName});
  var line = dom.htmlToDocumentFragment(
      templates.joinLine(o));
  this.getElement().appendChild(line);
};

cp.displayChat = function(playerId, content) {
  var dom = this.dom_;
  var player = this.game.players[playerId];
  var o = ({
    'name': player.name,
    'styleName': player.styleName,
    'content': content});
  var line = dom.htmlToDocumentFragment(
      templates.chatLine(o));
  this.getElement().appendChild(line);
};

cp.displayPass = function(playerId, stackId, toWhom) {
  var dom = this.dom_;
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

cp.displayStart = function(playerId) {
  var dom = this.dom_;
  var player = this.game.players[playerId];
  var o = ({
    'name': player.name,
    'styleName': player.styleName});
  var line = dom.htmlToDocumentFragment(
      templates.startLine(o));
  this.getElement().appendChild(line);
};

//ble.tpg.ui.
//scope-end
});
