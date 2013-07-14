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

var cycleClasses = function(classCycle) {

  var N = classCycle.length;
  var cc = classCycle.slice();

  return function(element) {
    var lst = element.classList;
    for(var i = 0; i < classCycle.length; i++) {
      var cl1 = cc[i], cl2 = cc[(i + 1) % N];
      if(lst.contains(cl1)) {
        lst.remove(cl1);
        lst.add(cl2);
        return;
      }
    }

    console.warn("didn't find any of the classes\n" + "\n\t".join(cc));
    lst.add(cc[0]);
  };
};

var cycleInstructions_internal = cycleClasses(["start","draw","describe"]);
var cycleInstructions = function(element) {
  var cl = element.classList;
  if(!cl.contains("instructions"))
    return false;
  cycleInstructions_internal(element);
  return true;
};

var cycleInstructionListener = function(event) {
  if(cycleInstructions(event.target))
    event.stopImmediatePropagation();
};

document.body.addEventListener("click", cycleInstructionListener);
document.body.addEventListener("click", toggler);

var headline = document.getElementsByClassName("status-headline")[0];
var statusCycler = cycleClasses(["before-game", "at-work", "waiting"]);
headline.addEventListener("click", function(event) {
  statusCycler(headline);
});

