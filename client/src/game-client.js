goog.require('goog.events');
goog.require('goog.net.EventType');
goog.require('goog.net.XhrIo');


goog.require('ble.net.QueryTimeComet');
goog.require('ble.tpg.model.Game');
goog.require('ble.tpg.ui.ChatContainer');
goog.provide('ble.tpg.game.setupClient');


goog.scope(function() {
var netType = goog.net.EventType;
var cometType = ble.net.EventType;
var modelType = ble.tpg.model.EventType;

var scope = ble.tpg.game;
var console = window.console;
var JSON = window.JSON;
var model = ble.tpg.model;

scope.onResponse = function(event) {
  console.log(event);
};

scope.setupClient = function() {
  scope.cometLoop = new ble.net.QueryTimeComet("./events");
  var cometLoop = scope.cometLoop;
  var otherXhr = new goog.net.XhrIo();
  scope.otherXhr = otherXhr;
  goog.events.listenOnce(
      otherXhr, 
      [netType.SUCCESS, netType.TIMEOUT, netType.ERROR],
      function(event) {
        if(event.type == netType.SUCCESS) {
          console.log(this.getResponse());
          try {
            var jsonObj = JSON.parse(this.getResponse());
            var lastTime = jsonObj['lastTime'] || 0;
            var game = model.Game.fromJSON(jsonObj);
            scope.chatContainer = new ble.tpg.ui.ChatContainer(game);
            goog.events.listen(cometLoop, cometType.COMET_DATA, game);
            goog.events.listen(game, modelType.JOIN_GAME, function(e) { console.log("soooeee"); console.log(e)});
            var chatContainer = document.getElementById('chat-container');
            scope.chatContainer.render(chatContainer);
            goog.dom.replaceNode(scope.chatContainer.getElement(), chatContainer);
            scope.chatContainer.getElement().id = chatContainer.id;
            scope.chatContainer.getElement().className += chatContainer.className;
            cometLoop.runAt(game.lastTime);
          } catch(e) {
            console.log(e);
            console.log(event);
            console.log(this.getResponse());
            console.log("frickin' error");
          }
        } else { 
          console.log(event);
          console.log("frickin' error");
        }
      });
  otherXhr.send("./")
};




scope.setupClient();
});
