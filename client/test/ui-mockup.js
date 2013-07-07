var toggleOne = function(el) {
  if(el.classList.contains("inactive")) {
    el.classList.remove("inactive");
    el.classList.add("active");
  } else if(el.classList.contains("active")) {
    el.classList.remove("active");
    el.classList.add("inactive");
  }
};

var canToggle = function(el) {
  return (el instanceof HTMLElement) &&
         (  el.classList.contains("inactive") ||
            el.classList.contains("active")      );

}

var toggler = function(event) {
  var target = event.target;
  while(target != null && target != event.currentTarget) {
    if(canToggle(target)) {
      console.log("toggling one");
      toggleOne(target);
      return;
    }
    target = target.parentElement;
  }
};

var cycleInstructions = function(element) {
  var cl = element.classList;
  if(!cl.contains("instructions"))
    return false;
  if(cl.contains("start")) {
    cl.remove("start");
    cl.add("draw");
  } else if(cl.contains("draw")) {
    cl.remove("draw");
    cl.add("describe");
  } else if(cl.contains("describe")) {
    cl.remove("describe");
    cl.add("start");
  }
  return true;
};

var cycleInstructionListener = function(event) {
  if(cycleInstructions(event.target))
    event.stopImmediatePropagation();
};

document.body.addEventListener("click", cycleInstructionListener);
document.body.addEventListener("click", toggler);

