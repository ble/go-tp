goog.require('goog.events');

goog.require('goog.ui.Component');
goog.require('goog.ui.Button');

goog.require('goog.net.EventType');

goog.require('goog.labs.net.xhr');

goog.require('ble.scribble.EventType');

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
var Stack = ble.tpg.model.Stack;

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
  this.game = new Game('./');
  goog.events.listen(this.game, modelType.READ_STATE, this.setupWithState, false, this);
  this.game.requestState();
};


Client.prototype.bindToExistingDom = function() {
  this.chatContainerDiv    = document.getElementById('chat-container');
  this.drawingContainerDiv = document.getElementById('drawing-container');
  this.startButtonDiv = document.getElementById('start-game-button');
};

Client.prototype.setupWithState = function() {
  goog.events.listen(this.game, modelType.PASS, this.handlePass, false, this);

  //set up the chat ui
  this.chatContainer = new ble.tpg.ui.ChatContainer(this.game);
  ble.util.replaceElemWithComponent(this.chatContainerDiv, this.chatContainer);

  //set up the scribbler
  this.scribbler = new ble.tpg.ui.Scribbler(this.game);
  ble.util.replaceElemWithComponent(this.drawingContainerDiv, this.scribbler);
  this.scribbler.scribble.setEnabled(false);
  goog.events.listen(
      this.scribbler, 
      ble.scribble.EventType.DRAW_END,
      this.postDraw,
      false,
      this);

  if(!this.game.isStarted) {
  //set up the start-game button
    this.startButton = new goog.ui.Button("start game!"); 
    ble.util.replaceElemWithComponent(this.startButtonDiv, this.startButton);
    goog.events.listen(
      this.startButton,
      goog.ui.Component.EventType.ACTION,
      this.postStartGame,
      false,
      this);
  }

  //set up the comet loop
  goog.events.listen(this.cometLoop, cometType.COMET_DATA, this.game);
  this.cometLoop.runAt(this.game.lastTime); 
  this.setupGameState(); 
}

Client.prototype.postDraw = function(e) {
  //TODO: disable drawing while waiting for a draw part to post?
  var myStack = this.game.getMyStacks()[0];
  var myDrawing = myStack.drawings[0];
  console.log(myDrawing);
  var drawRequest = xhr.post(
      myDrawing.url,
      JSON.stringify({
        'actionType': 'draw',
        'content': e.drawn}),
      {'headers':
        {'Content-Type': 'application/json'}});

  console.log(e);
  drawRequest.wait((function(result) {
    console.log(result.getState());
    console.log(result);
  }).bind(this));

};

Client.prototype.postStartGame = function() {
  this.startButton.setEnabled(false);
  var startRequest = xhr.post(
      "./start",
      JSON.stringify({'actionType': 'startGame'}),
      {
        'headers': {
          'Content-Type': 'application/json'
        }  
      });
  startRequest.wait((function(result) {
    if(result.getState() == resultState.SUCCESS) {
      this.startButton.dispose(); 
    } else {
      this.startButton.setEnabled(true);
    }
  }).bind(this));

};

Client.prototype.handlePass = function(e) {
  if(e.to.isYou && e.stack === this.game.getMyStacks()[0]) {
    this.makeDrawingReady();
  }
};

Client.prototype.makeDrawingReady = function(e) {
  var myStack = this.game.getMyStacks()[0];
  var stackResponse = xhr.get(myStack.url);
  stackResponse.wait(this.handleStackResult.bind(this));
};

Client.prototype.handleStackResult = function(result) {
  if(result.getState() == resultState.SUCCESS) {
    var newStack = Stack.fromJSON(JSON.parse(result.getValue()));
    var oldStack = this.game.stacksById[newStack.id];
    oldStack.drawings = newStack.drawings; 
    this.setupGameState();
    this.scribbler.scribble.setEnabled(true);
    //TODO: get the relevant drawings here...
  } else {
    console.error("shoot");
  }
};

Client.prototype.setupGameState = function() {
  var game = this.game;
  var me = game.getMyPlayer();
  if(game.isStarted &&
     goog.isDefAndNotNull(me) ) {
    this.scribbler.scribble.setEnabled(true);
    console.log("i guess we should have a drawing or something");
  }
};

window.tpg_client = new ble.tpg.Client();
window.tpg_client.initialize();
});
