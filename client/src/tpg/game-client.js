goog.require('goog.events');
goog.require('goog.net.EventType');

goog.require('goog.labs.net.xhr');

goog.require('ble.net.QueryTimeComet');
goog.require('ble.tpg.model.Game');

goog.require('ble.tpg.ui.ChatContainer');
goog.require('ble.tpg.ui.Scribbler');



goog.provide('ble.tpg.game.setupClient');



goog.scope(function() {
var console = window.console;
var JSON = window.JSON;

var xhr = goog.labs.net.xhr;
var netType = goog.net.EventType;
var resultState = goog.result.Result.State;

var cometType = ble.net.EventType;
var modelType = ble.tpg.model.EventType; 

var Game = ble.tpg.model.Game;

/** @constructor */
ble.tpg.Client = function() {
  this.cometLoop = new ble.net.QueryTimeComet('./events');
  this.initialized = false;
};
var Client = ble.tpg.Client;

Client.prototype.initialize = function() {
  if(this.initialized)
    return;
  this.bindToExistingDom();
  this.requestInitialState();
};


Client.prototype.bindToExistingDom = function() {
  this.chatContainerDiv    = document.getElementById('chat-container');
  this.drawingContainerDiv = document.getElementById('drawing-container');
};

Client.prototype.requestInitialState = function() {
  var url = './'; 
  var stateRequest = xhr.get(url);
  stateRequest.wait(goog.bind(this.processStateResponse, this));
};

Client.prototype.processStateResponse = function(stateResponse) {
  if(stateResponse.getState() == resultState.SUCCESS) {
    console.log('got state response');
    console.log(stateResponse.getValue());
    //try {
      //set up the game model
      var jsonObj = JSON.parse(stateResponse.getValue());
      var lastTime = jsonObj['lastTime'] || 0;
      this.game = Game.fromJSON(jsonObj);

      //set up the chat ui
      this.chatContainer = new ble.tpg.ui.ChatContainer(this.game);
      ble.util.replaceElemWithComponent(this.chatContainerDiv, this.chatContainer);
    
      //set up the scribbler
      this.scribbler = new ble.tpg.ui.Scribbler(this.game);
      ble.util.replaceElemWithComponent(this.drawingContainerDiv, this.scribbler);
      this.scribbler.scribble.setEnabled(false);
      //set up the comet loop
      goog.events.listen(this.cometLoop, cometType.COMET_DATA, this.game);
      goog.events.listen(this.game, modelType.ALL, this.handleGameEvent, false, this);
      this.cometLoop.runAt(this.game.lastTime); 
      this.setupGameState();
 /*   } catch(e) {
      console.log('frickin\' error.');
      console.log(e);
    }*/
  } else {
    this.requestInitialState();
  }
};

Client.prototype.handleGameEvent = function(e) {

};

Client.prototype.setupGameState = function() {
  var game = this.game;
  var me = game.getMyPlayer();
  if(game.isStarted &&
     goog.isDefAndNotNull(me) &&
     game.inPlay[me.id] > 0) {
    console.log("i guess we should have a drawing or something");
  }
};

window.tpg_client = new ble.tpg.Client();
window.tpg_client.initialize();
});
