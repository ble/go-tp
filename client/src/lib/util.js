goog.require('goog.dom.forms');

goog.provide('ble.util.formToJSMap');
goog.provide('ble.util.clearFormText');

ble.util.formToJSMap = function(form) {
  var map = goog.dom.forms.getFormDataMap(form);
  var result = {};
  var keys = map.getKeys();
  for(var i = 0; i < keys.length; i++) {
    result[keys[i]] = map.get(keys[i])[0];
  }
  return result;
};

ble.util.clearFormText = function(form) {
  var children = form.children;
  for(var i = 0; i < children.length; i++) {
    var child = children[i];
    if(child.tagName == 'INPUT' && child.type == 'text') {
      child.value = '';
    }
  }
};
