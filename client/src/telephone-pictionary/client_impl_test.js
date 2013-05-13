goog.require('ble.telephone_pictionary.ClientImpl');

goog.require('goog.testing.jsunit');
goog.require('goog.testing.PropertyReplacer');
goog.require('goog.testing.TestCase');
goog.require('goog.testing.TestRunner');

goog.require('goog.labs.mock');
goog.require('goog.labs.testing.Matcher');
goog.require('goog.labs.testing.HasEntriesMatcher');
goog.require('goog.labs.testing.AnythingMatcher');

/**
 * @constructor
 * @implements {goog.labs.testing.Matcher}
 */
goog.labs.testing.FromJsonMatcher = function(innerMatcher) {
  this.inner = innerMatcher;
};

goog.labs.testing.FromJsonMatcher.prototype.describe = function(v, opt_desc) {
  var o = window.JSON.parse(v);
  return this.inner.describe(o, opt_desc);
};

goog.labs.testing.FromJsonMatcher.prototype.matches = function(v) {
  return this.inner.matches(window.JSON.parse(v));
}

function fromJson(m) { return new goog.labs.testing.FromJsonMatcher(m); }
var myTests;

goog.scope(function() {
  var console = window.console;
  var Test = goog.testing.TestCase.Test;

  var mock = goog.labs.mock;
  var testing = goog.labs.testing;
  var xhr = goog.labs.net.xhr;
  var result = goog.result;

  var _ = ble.telephone_pictionary;
  var ClientImpl = _.ClientImpl;

  myTests = new goog.testing.TestCase('ClientImpl tests');
  myTests.setUp = function() {
    this.subber = new goog.testing.PropertyReplacer();
  };
  myTests.tearDown = function() {
    this.subber.reset();
    this.subber = null;
  }

  myTests.add(new Test(
      'ClientImpl makes correct method / URI XHRs',
      function() {
        var baseUrl = 'http://localhost:24769/shard0/';
        var gameUrl = 'http://localhost:24769/shard0/game/foobar/';

        var instance = new ClientImpl(gameUrl);
        this.subber.replace(xhr, 'get', mock.mockFunction(xhr.get));
        this.subber.replace(xhr, 'post', mock.mockFunction(xhr.post));
        var fake = new result.SimpleResult();

        mock.when(xhr.get)(anything()).thenReturn(fake);
        mock.when(xhr.post)(anything(), fromJson(anything()), _.jsonHeader).thenReturn(fake);

        instance.getGameState();
        mock.verify(xhr.get)(gameUrl);

        instance.getStack('piffler');
        mock.verify(xhr.get)(baseUrl+'stack/piffler');

        instance.getDrawing('barfoo');
        mock.verify(xhr.get)(baseUrl+'drawing/barfoo');

        instance.startGame();
        mock.verify(xhr.post)(
            gameUrl+'start',
            fromJson(hasEntries({'actionType':'startGame'})),
            _.jsonHeader);

        instance.appendToDrawing('barfoo', "hello");
        mock.verify(xhr.post)(
            baseUrl+'drawing/barfoo',
            fromJson(hasEntries({'actionType':'draw','content':'hello'})),
            _.jsonHeader);

        instance.passStack('piffler');
        mock.verify(xhr.post)(
            gameUrl+'pass',
            fromJson(hasEntries({'actionType':'passStack'})),
            _.jsonHeader);

        instance.chat('colorless green ideas sleep furiously');
        mock.verify(xhr.post)(
            gameUrl+'chat',
            fromJson(hasEntries({
              'actionType':'chat',
              'content':'colorless green ideas sleep furiously'})),
            _.jsonHeader);

      },
      myTests));
});

window.G_testRunner.initialize(myTests);

