var toggleAll = function() {
  var inactive = document.getElementsByClassName("inactive");
  inactive = [].slice.call(inactive);
  var active = document.getElementsByClassName("active");
  active = [].slice.call(active);
  for(var i = 0; i < inactive.length; i++) {
    inactive[i].classList.remove("inactive");
    inactive[i].classList.add("active");
  }
  for(var i = 0; i < active.length; i++) {
    active[i].classList.remove("active");
    active[i].classList.add("inactive");
  }
};

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
  if(canToggle(event.target)) {
    console.log("toggling one");
    toggleOne(event.target);
  } else {
    console.log("toggling all");
    toggleAll();
  }
};
document.body.addEventListener("click", toggler);

