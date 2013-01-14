goog.require('goog.dom');
goog.require('goog.dom.forms');
goog.require('goog.events');

goog.require('goog.labs.net.xhr');

goog.require('goog.ui.Component');

goog.require('ble.util.formToJSMap');
goog.require('ble.util.clearFormText');
goog.require('ble.tpg.templates');
goog.require('ble.tpg.model.EventType');

goog.provide('ble.tpg.ui.ChatContainer');
goog.scope(function() {
//scope-start

var Component = goog.ui.Component;
var ModelType = ble.tpg.model.EventType;
var templates = ble.tpg.templates;
var xhr = goog.labs.net.xhr;
var resultState = goog.result.Result.State;
var forms = goog.dom.forms;

/**
 * @constructor
 * @extends{goog.ui.Component}
 * @param{ble.tpg.model.Game} game
 */
ble.tpg.ui.ChatContainer = function(game) {
  Component.call(this);
  this.game = game;
  this.chats = new ble.tpg.ui.Chats(game);
  this.chatInput = new ble.tpg.ui.ChatInput(game);
  this.addChild(this.chats, true);
  this.addChild(this.chatInput, true);
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
      this.displayChat(event.player, event.content);
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

cp.displayChat = function(player, content) {
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

/**
 * @constructor
 * @extends{goog.ui.Component}
 */
ble.tpg.ui.ChatInput = function(game) {
  Component.call(this);
  this.game = game;
};
goog.inherits(ble.tpg.ui.ChatInput, Component);

var cip = ble.tpg.ui.ChatInput.prototype;

cip.createDom = function() {
  goog.base(this, 'createDom');
  var elt = this.getElement();
  var dom = this.getDomHelper();
  var text = dom.createDom(
      'input', 
      {'type': 'text', 'name': 'content', 'class': 'chat-text'});
  var button = dom.createDom(
      'input',
      {'type': 'submit', 'value': 'chat!', 'class': 'chat-button'});
  var actionType = dom.createDom(
      'input',
      {'type': 'hidden', 'name': 'actionType', 'value': 'chat'});
  var form = dom.createDom('form', null, text, button, actionType);
  this.form = form;
  elt.appendChild(form);
};

cip.enterDocument = function() {
  goog.base(this, 'enterDocument');
  goog.events.listen(this.form, 'submit', this);
};

cip.handleEvent = function(event) {
  event.preventDefault();
  if(this.form.children[0].disabled)
    return;
  var form = this.form;
  var map = ble.util.formToJSMap(this.form);
  forms.setDisabled(form, true);
  var chatRequest = xhr.post(
      this.game.url + "chat",
      JSON.stringify(map),
      {
        'headers': {
          'Content-Type': 'application/json'}
      });
  chatRequest.wait(function(result) {
    forms.setDisabled(form, false);
    if(result.getState() == resultState.SUCCESS) {
      ble.util.clearFormText(form);
    }
  });
  console.log(map);
  console.log(event);
};

//scope-end
});
