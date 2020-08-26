/*
  This file contains custom functionality for the website
  */

var workerPassword = "000000000";
var count = 0;
var cond = false;

var lion = false;
var giraffe = false;
var snake = false;
var rhino = false;
var elephant = false;
var cow = false;
var camel = false;
var panther = false;
var monkey = false;

/* Puts the functionality to the GO fiel in the Toggle table for pagination*/
function changeGoLink(quantity, service, id) {
  var go = document.getElementById("go");
  var goField = document.getElementById("goField");

  var goToStep = goField.value;

  var goToStep = parseInt(goToStep) || -1;

  if (goToStep == -1) {
    alert("The step should be a number!");
    return false;
  } else if (goToStep < 0 || (goToStep > quantity && goToStep != "")) {
    alert("You selected a step that does not exist!");
    return false;
  }
  goToStep = goToStep - 1;
  window.location.href =
    "/videos/step-video-course/" + id + "/" + goToStep + "/" + service;
}

/* Sets the pagination under the table to fix the number of pages are for the table
quantity is the length of the array, numberNoURL is the /:id of the URL + 1 */
function setPagination(quantity, numberNoURL, service, id, cond) {
  var prev = document.getElementById("prev");
  var next = document.getElementById("next");
  var go = document.getElementById("go");
  var goField = document.getElementById("goField");

  //Disable color if we are on the limits
  if (numberNoURL - 1 == 0) {
    prev.setAttribute("style", "background-color: #5DBCD2");
  }
  if (quantity == 1) {
    go.setAttribute("style", "background-color: #5DBCD2");
    goField.setAttribute("disabled", true);
    goField.setAttribute("style", "background-color: #e9ecef");
    prev.setAttribute("href", "");
    next.setAttribute("href", "");
  }
  if (numberNoURL == quantity) {
    next.setAttribute("style", "background-color: #5DBCD2");
    next.setAttribute("href", "");
  }

  if (numberNoURL < quantity) {
    if (cond == true) {
      next.setAttribute(
        "href",
        "/videos/edit-step/" + id + "/" + numberNoURL + "/" + service
      );
    } else {
      next.setAttribute(
        "href",
        "/videos/step-video-course/" + id + "/" + numberNoURL + "/" + service
      );
    }
  }
  if (numberNoURL > 0 && numberNoURL > 1) {
    if (cond == true) {
      prev.setAttribute(
        "href",
        "/videos/edit-step/" + id + "/" + (numberNoURL - 2) + "/" + service
      );
    } else {
      prev.setAttribute(
        "href",
        "/videos/step-video-course/" +
          id +
          "/" +
          (numberNoURL - 2) +
          "/" +
          service
      );
    }
  }
}

/*This functions translates the int number into a symbols for the numerical alphabet */
function fromNumberToSymbol(number) {
  var dict = {
    star: 0,
    pentagon: 0,
    square: 0,
    triangle: 0,
    cross: 0,
    line: 0,
    circle: 0,
  };

  var count = 0;
  while (number >= 1000) {
    count++;
    number = number - 1000;
  }
  dict.star = count;

  count = 0;
  while (number >= 500) {
    count++;
    number = number - 500;
  }
  dict.pentagon = count;

  count = 0;
  while (number >= 100) {
    count++;
    number = number - 100;
  }
  dict.square = count;

  count = 0;
  while (number >= 50) {
    count++;
    number = number - 50;
  }
  dict.triangle = count;

  count = 0;
  while (number >= 10) {
    count++;
    number = number - 10;
  }
  dict.cross = count;

  count = 0;
  while (number >= 5) {
    count++;
    number = number - 5;
  }
  dict.line = count;

  count = 0;
  while (number > 0) {
    count++;
    number = number - 1;
  }
  dict.circle = count;

  return dict;
}

/*This functions writes the symbols for the money into the page*/
function putSymbols(price, id) {
  if (price < 0) {
    price = -price;
  }
  var d = fromNumberToSymbol(price);
  var first = document.getElementById("first" + id);
  var second = document.getElementById("second" + id);
  var third = document.getElementById("third" + id);

  for (var i = 0; i < d.star; i++) {
    var a = document.createElement("li");
    var b = document.createElement("img");
    b.setAttribute("src", "/static/svg/star.svg");
    a.appendChild(b);
    first.appendChild(a);
  }

  for (var i = 0; i < d.pentagon; i++) {
    var a = document.createElement("li");
    var b = document.createElement("img");
    b.setAttribute("src", "/static/svg/polygon.svg");
    a.appendChild(b);
    first.appendChild(a);
  }

  for (var i = 0; i < d.square; i++) {
    var a = document.createElement("li");
    var b = document.createElement("img");
    b.setAttribute("src", "/static/svg/square.svg");
    a.appendChild(b);
    second.appendChild(a);
  }
  for (var i = 0; i < d.triangle; i++) {
    var a = document.createElement("li");
    var b = document.createElement("img");
    b.setAttribute("src", "/static/svg/triangle.svg");
    a.appendChild(b);
    second.appendChild(a);
  }
  for (var i = 0; i < d.cross; i++) {
    var a = document.createElement("li");
    var b = document.createElement("img");
    b.setAttribute("src", "/static/svg/cross.svg");
    a.appendChild(b);
    second.appendChild(a);
  }
  for (var i = 0; i < d.line; i++) {
    var a = document.createElement("li");
    var b = document.createElement("img");
    b.setAttribute("src", "/static/svg/line.svg");
    a.appendChild(b);
    third.appendChild(a);
  }
  for (var i = 0; i < d.circle; i++) {
    var a = document.createElement("li");
    var b = document.createElement("img");
    b.setAttribute("src", "/static/svg/circle.svg");
    a.appendChild(b);
    third.appendChild(a);
  }
}

/*This functions writes the symbols for the money into the page with another id for the page w-wallet.html*/
function putSymbols2(price, id) {
  if (price < 0) {
    price = -price;
  }
  var d = fromNumberToSymbol(price);
  var first = document.getElementById("f" + id);
  var second = document.getElementById("s" + id);
  var third = document.getElementById("t" + id);

  for (var i = 0; i < d.star; i++) {
    var a = document.createElement("li");
    var b = document.createElement("img");
    b.setAttribute("src", "/static/svg/star.svg");
    a.appendChild(b);
    first.appendChild(a);
  }

  for (var i = 0; i < d.pentagon; i++) {
    var a = document.createElement("li");
    var b = document.createElement("img");
    b.setAttribute("src", "/static/svg/polygon.svg");
    a.appendChild(b);
    first.appendChild(a);
  }

  for (var i = 0; i < d.square; i++) {
    var a = document.createElement("li");
    var b = document.createElement("img");
    b.setAttribute("src", "/static/svg/square.svg");
    a.appendChild(b);
    second.appendChild(a);
  }
  for (var i = 0; i < d.triangle; i++) {
    var a = document.createElement("li");
    var b = document.createElement("img");
    b.setAttribute("src", "/static/svg/triangle.svg");
    a.appendChild(b);
    second.appendChild(a);
  }
  for (var i = 0; i < d.cross; i++) {
    var a = document.createElement("li");
    var b = document.createElement("img");
    b.setAttribute("src", "/static/svg/cross.svg");
    a.appendChild(b);
    second.appendChild(a);
  }
  for (var i = 0; i < d.line; i++) {
    var a = document.createElement("li");
    var b = document.createElement("img");
    b.setAttribute("src", "/static/svg/line.svg");
    a.appendChild(b);
    third.appendChild(a);
  }
  for (var i = 0; i < d.circle; i++) {
    var a = document.createElement("li");
    var b = document.createElement("img");
    b.setAttribute("src", "/static/svg/circle.svg");
    a.appendChild(b);
    third.appendChild(a);
  }
}

/*The function when you press the delete button in the worker password screen */
function deleteClick() {
  lion = false;
  giraffe = false;
  snake = false;
  rhino = false;
  elephant = false;
  cow = false;
  camel = false;
  panther = false;
  monkey = false;

  document.getElementById("interrogation1").src =
    "/static/svg/interrogation.svg";
  document.getElementById("interrogation2").src =
    "/static/svg/interrogation.svg";
  document.getElementById("interrogation3").src =
    "/static/svg/interrogation.svg";
  document.getElementById("interrogation4").src =
    "/static/svg/interrogation.svg";

  document.getElementById("lion").src =
    "/static/img/login/Animals_OFF/Lion_OFF.png";
  document.getElementById("giraffe").src =
    "/static/img/login/Animals_OFF/Giraffe_OFF.png";
  document.getElementById("snake").src =
    "/static/img/login/Animals_OFF/Snake_OFF.png";
  document.getElementById("monkey").src =
    "/static/img/login/Animals_OFF/Monkey_OFF.png";
  document.getElementById("panther").src =
    "/static/img/login/Animals_OFF/Panther_OFF.png";
  document.getElementById("camel").src =
    "/static/img/login/Animals_OFF/Camel_OFF.png";
  document.getElementById("cow").src =
    "/static/img/login/Animals_OFF/Cow_OFF.png";
  document.getElementById("elephant").src =
    "/static/img/login/Animals_OFF/Elephant_OFF.png";
  document.getElementById("rhino").src =
    "/static/img/login/Animals_OFF/Rhino_OFF.png";
  workerPassword = "000000000";
  count = 0;
  document.getElementById("password").value = workerPassword;
  var child = document.getElementById("passwordShow");
  child.innerHTML = "";
}

//Change Lion image
function lionClick() {
  if (checkSizePassword()) {
    return;
  }
  if (lion == false) {
    lion = true;
    document.getElementById("lion").src =
      "/static/img/login/Animals_ON/Lion_ON.png";
    count++;
    setCharAt(0, count);
    // Write Lion in the password
    var id = "interrogation" + count;
    document.getElementById(id).src =
      "/static/img/login/Animals_ON/Lion_ON.png";
  }
}

//Change Giraffe image
function giraffeClick() {
  if (checkSizePassword()) {
    return;
  }
  if (giraffe == false) {
    giraffe = true;
    document.getElementById("giraffe").src =
      "/static/img/login/Animals_ON/Giraffe_ON.png";
    count++;
    setCharAt(1, count);
    // Write Giraffe in the password
    var id = "interrogation" + count;
    document.getElementById(id).src =
      "/static/img/login/Animals_ON/Giraffe_ON.png";
  }
}

//Change Snake image
function snakeClick() {
  if (checkSizePassword()) {
    return;
  }
  if (snake == false) {
    snake = true;
    document.getElementById("snake").src =
      "/static/img/login/Animals_ON/Snake_ON.png";
    count++;
    setCharAt(2, count);
    // Write Snake in the password
    var id = "interrogation" + count;
    document.getElementById(id).src =
      "/static/img/login/Animals_ON/Snake_ON.png";
  }
}

//Change Rhino image
function rhinoClick() {
  if (checkSizePassword()) {
    return;
  }
  if (rhino == false) {
    rhino = true;
    document.getElementById("rhino").src =
      "/static/img/login/Animals_ON/Rhino_ON.png";
    count++;
    setCharAt(3, count);
    // Write Snake in the password
    var id = "interrogation" + count;
    document.getElementById(id).src =
      "/static/img/login/Animals_ON/Rhino_ON.png";
  }
}

//Change Elephant image
function elephantClick() {
  if (checkSizePassword()) {
    return;
  }
  if (elephant == false) {
    elephant = true;
    document.getElementById("elephant").src =
      "/static/img/login/Animals_ON/Elephant_ON.png";
    count++;
    setCharAt(4, count);
    // Write Elephant in the password
    var id = "interrogation" + count;
    document.getElementById(id).src =
      "/static/img/login/Animals_ON/Elephant_ON.png";
  }
}

//Change Cow image
function cowClick() {
  if (checkSizePassword()) {
    return;
  }
  if (cow == false) {
    cow = true;
    document.getElementById("cow").src =
      "/static/img/login/Animals_ON/Cow_ON.png";
    count++;
    setCharAt(5, count);
    // Write Cow in the password
    var id = "interrogation" + count;
    document.getElementById(id).src = "/static/img/login/Animals_ON/Cow_ON.png";
  }
}

//Change Camel image
function camelClick() {
  if (checkSizePassword()) {
    return;
  }
  if (camel == false) {
    camel = true;
    document.getElementById("camel").src =
      "/static/img/login/Animals_ON/Camel_ON.png";
    count++;
    setCharAt(6, count);
    // Write Camel in the password
    var id = "interrogation" + count;
    document.getElementById(id).src =
      "/static/img/login/Animals_ON/Camel_ON.png";
  }
}

//Change Panther image
function pantherClick() {
  if (checkSizePassword()) {
    return;
  }
  if (panther == false) {
    panther = true;
    document.getElementById("panther").src =
      "/static/img/login/Animals_ON/Panther_ON.png";
    count++;
    setCharAt(7, count);
    // Write Panther in the password
    var id = "interrogation" + count;
    document.getElementById(id).src =
      "/static/img/login/Animals_ON/Panther_ON.png";
  }
}

//Change Monkey image count, , cond,
function monkeyClick() {
  if (checkSizePassword()) {
    return;
  }
  if (monkey == false) {
    monkey = true;
    document.getElementById("monkey").src =
      "/static/img/login/Animals_ON/Monkey_ON.png";
    count++;
    setCharAt(8, count);
    // Write Monkey in the password
    var id = "interrogation" + count;
    document.getElementById(id).src =
      "/static/img/login/Animals_ON/Monkey_ON.png";
  }
}

//Add the div for write the workerPassword when you click the image
function divPassword() {
  b = document.getElementById("button");
  var c = document.createElement("div");
  c.setAttribute("class", "row");
  c.setAttribute("id", "passwordShow");
  c.setAttribute("style", "margin-top: 40px;");
  var d = document.getElementById("myForm");
  d.insertBefore(c, b);
}

//To change the number in the Worker Password
function setCharAt(index, chr) {
  if (index > workerPassword.length - 1) return workerPassword;
  workerPassword =
    workerPassword.substr(0, index) + chr + workerPassword.substr(index + 1);
  document.getElementById("password").value = workerPassword;
}

// To put the name of the photo to upload
function feedbackPhoto() {
  var place = document.getElementById("placeHolderPhoto");
  place.innerHTML = "Photo added!";
  // place.innerHTML = ic.value.replace(/^.*[\\\/]/, '');
  place.classList.remove("is-invalid");
}

//Everytime we write something in the Name or Surname field, we put the username according to the name and username
function fillUsername(event) {
  var name = document.getElementById("name").value.toLowerCase();
  var surname = document.getElementById("surname").value.toLowerCase();
  document.getElementById("username").value = name + surname;
}

/* This method put the icons for the worker in order to choose the password
 every time we change the value of the selector from Worker to Manager */
function passwordWorker() {
  var selection = document.getElementById("role").value;
  if (selection == "Worker") {
    isWorker();
  } else {
    isManager();
  }
}

/* Set the password and the icons for the worker */
function isWorker() {
  divPassword();
  var d = document.getElementById("password");
  d.value = "";
  d.setAttribute("readonly", "true");
  showIcons();
}

/* Set the password and removes the icons for the manager */
function isManager() {
  deleteClick();
  showIcons();
  workerPassword = "000000000";
  count = 0;
  var d = document.getElementById("password");
  d.value = "";
  d.removeAttribute("readonly");
  document.getElementById("password").value = "";
}

//To put the name of the file to upload
function feedbackFileAudio() {
  var place = document.getElementById("placeHolderAudio");
  var ic = document.getElementById("audio");
  place.innerHTML = ic.value.replace(/^.*[\\\/]/, "");
  place.classList.remove("is-invalid");
}

//To put the name of the file to upload
function feedbackFilePhoto() {
  var place = document.getElementById("placeHolderPhoto");
  // var ic = document.getElementById("photo");
  // place.innerHTML = ic.value.replace(/^.*[\\\/]/, '');
  place.innerHTML = "Photo Added!";
  place.classList.remove("is-invalid");
}

//To put the name of the file to upload
function feedbackFileVideo() {
  var place = document.getElementById("placeHolderVideo");
  var ic = document.getElementById("video");
  place.innerHTML = ic.value.replace(/^.*[\\\/]/, "");
  place.classList.remove("is-invalid");
}

function goBack() {
  window.history.go(-1);
  return false;
}

/* This method shows the icons for the worker in the new/edit User pages  */
var condNavbar = false;
function showIcons() {
  var imgClick = document.getElementById("workerPassword");
  if (!condNavbar) {
    imgClick.style.display = "block";
    condNavbar = true;
  } else {
    imgClick.style.display = "none";
    condNavbar = false;
  }
}

/* returns true if we reach the limit of 4 animals in the password */
function checkSizePassword() {
  if (count >= 4) {
    alert("Sorry, the maximun animals for the password is 4 animals!");
    return true;
  } else {
    return false;
  }
}

/* Change the selec image if the user clicks the image, this is more situable for the tablets and tactil interaction */
$(".image-workers").click(function (e) {
  if ($(this).hasClass("no-selected")) {
    var product = $(this).attr("id");
    console.log($("#workers-" + product));
    $(".image-workers").removeClass("selected");
    $(".image-workers").addClass("no-selected");
    $("#workers-" + product).prop("checked", true);
    $(this).addClass("selected");
    $(this).removeClass("no-selected");
    increaseWorker();
  } else if ($(this).hasClass("selected")) {
    var product = $(this).attr("id");

    $(this).addClass("no-selected");
    $("#workers-" + product).prop("checked", false);
    $(this).removeClass("selected");
  }
});

function playAudio() {
  var x = document.getElementById("myAudio");
  x.play();
}

function playAudioInformation() {
  var x = document.getElementById("audio-information");
  x.play();
}

/* This function is used in the form to Capitalice the first letter */
function capitalizeFirstLetter(string) {
  return string.charAt(0).toUpperCase() + string.slice(1);
}

/* THIS ARE ALL THE FUNCTIONS RELATED WITH NAVIGATION BAR NAV BAR NAVBAR */
//Highlight the icon Performance in the navigation Bar
function performanceON() {
  var normalIcon = document.getElementById("first-icon");
  normalIcon.setAttribute("src", "/static/svg/performance-hover.svg");
  var mobileIcon = document.getElementById("first-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/performance-hover.svg");
  firstON();
}

//Removes the highligh from the navigation bar
function performanceOFF() {
  var normalIcon = document.getElementById("first-icon");
  normalIcon.setAttribute("src", "/static/svg/performance.svg");
  var mobileIcon = document.getElementById("first-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/performance.svg");
  anyOFF("first");
}

//Highlight the icon Users in the navigation Bar
function backON() {
  var normalIcon = document.getElementById("back-icon");
  normalIcon.setAttribute("src", "/static/svg/back-hover.svg");
  var mobileIcon = document.getElementById("back-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/back-hover.svg");
  var text = document.getElementById("back-text");
  text.setAttribute("style", "color: #CF1F26;");
  var section = document.getElementById("back-section");
  section.setAttribute(
    "style",
    "color: #CF1F26; border-right: 4px solid #CF1F26;"
  );
  var textMobile = document.getElementById("back-text-mobile");
  textMobile.setAttribute("style", "color: #CF1F26;");
  var sectionMobile = document.getElementById("back-section-mobile");
  sectionMobile.setAttribute(
    "style",
    "color: #CF1F26; border-bottom: 4px solid #CF1F26;"
  );
}

//Removes the highligh from the navigation bar
function backOFF() {
  var normalIcon = document.getElementById("back-icon");
  normalIcon.setAttribute("src", "/static/svg/back.svg");
  var mobileIcon = document.getElementById("back-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/back.svg");
  var text = document.getElementById("back-text");
  text.removeAttribute("style");
  var section = document.getElementById("back-section");
  section.removeAttribute("style");
  var textMobile = document.getElementById("back-text-mobile");
  textMobile.removeAttribute("style");
  var sectionMobile = document.getElementById("back-section-mobile");
  sectionMobile.removeAttribute("style");
}

// Highlight the icon Sync in the navigation Bar
function syncON() {
  var normalIcon = document.getElementById("sync-icon");
  normalIcon.setAttribute("src", "/static/svg/sync-hover.svg");
  var mobileIcon = document.getElementById("sync-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/sync-hover.svg");
  secondON();
}

// Removes the highligh from the navigation bar
function syncOFF() {
  var normalIcon = document.getElementById("sync-icon");
  normalIcon.setAttribute("src", "/static/svg/sync.svg");
  var mobileIcon = document.getElementById("sync-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/sync.svg");
  anyOFF("second");
}
// Highlight the icon Deliveries in the navigation Bar
function deliveriesON() {
  var normalIcon = document.getElementById("first-icon");
  normalIcon.setAttribute("src", "/static/svg/deliveries-hover.svg");
  var mobileIcon = document.getElementById("first-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/deliveries-hover.svg");
  firstON();
}

// Removes the highligh from the navigation bar
function deliveriesOFF() {
  var normalIcon = document.getElementById("first-icon");
  normalIcon.setAttribute("src", "/static/svg/deliveries.svg");
  var mobileIcon = document.getElementById("first-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/deliveries.svg");
  anyOFF("first");
}

// Highlight the icon Deliveries in the navigation Bar
function testON() {
  var normalIcon = document.getElementById("second-icon");
  normalIcon.setAttribute("src", "/static/svg/test-hover.svg");
  var mobileIcon = document.getElementById("second-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/test-hover.svg");
  secondON();
}

// Removes the highligh from the navigation bar
function testOFF() {
  var normalIcon = document.getElementById("second-icon");
  normalIcon.setAttribute("src", "/static/svg/test.svg");
  var mobileIcon = document.getElementById("second-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/test.svg");
  anyOFF("second");
}

// Highlight the icon Products in the navigation Bar
function activityON() {
  var normalIcon = document.getElementById("activity-icon");
  normalIcon.setAttribute("src", "/static/svg/activity-hover.svg");
  var mobileIcon = document.getElementById("activity-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/activity-hover.svg");
  firstON();
}

//Removes the highligh from the navigation bar
function activityOFF() {
  var normalIcon = document.getElementById("activity-icon");
  normalIcon.setAttribute("src", "/static/svg/activity.svg");
  var mobileIcon = document.getElementById("activity-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/activity.svg");
  anyOFF("first");
}

//Highlight the icon Records in the navigation Bar
function generalON() {
  var normalIcon = document.getElementById("general-icon");
  normalIcon.setAttribute("src", "/static/svg/general-hover.svg");
  var mobileIcon = document.getElementById("general-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/general-hover.svg");
  thirdON();
}

//Removes the highligh from the navigation bar
function generalOFF() {
  var normalIcon = document.getElementById("general-icon");
  normalIcon.setAttribute("src", "/static/svg/general.svg");
  var mobileIcon = document.getElementById("general-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/general.svg");
  anyOFF("third");
}

// Removes the highligh from the navigation bar
function walletOFF() {
  var normalIcon = document.getElementById("first-icon");
  normalIcon.setAttribute("src", "/static/svg/wallet.svg");
  var mobileIcon = document.getElementById("first-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/wallet.svg");
  anyOFF("first");
}

// Highlight the icon Wallet in the navigation Bar
function walletON() {
  var normalIcon = document.getElementById("first-icon");
  normalIcon.setAttribute("src", "/static/svg/wallet-hover.svg");
  var mobileIcon = document.getElementById("first-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/wallet-hover.svg");
  firstON();
}

// Removes the highligh from the navigation bar
function timeOFF() {
  var normalIcon = document.getElementById("first-icon");
  normalIcon.setAttribute("src", "/static/svg/time.svg");
  var mobileIcon = document.getElementById("first-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/time.svg");
  anyOFF("first");
}

// Highlight the icon Wallet in the navigation Bar
function timeON() {
  var normalIcon = document.getElementById("first-icon");
  normalIcon.setAttribute("src", "/static/svg/time-hver.svg");
  var mobileIcon = document.getElementById("first-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/time-hver.svg");
  firstON();
}

// Removes the highligh from the navigation bar
function qrOFF() {
  var normalIcon = document.getElementById("first-icon");
  normalIcon.setAttribute("src", "/static/svg/qr.svg");
  var mobileIcon = document.getElementById("first-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/qr.svg");
  anyOFF("first");
}

// Highlight the icon QR in the navigation Bar
function qrON() {
  var normalIcon = document.getElementById("first-icon");
  normalIcon.setAttribute("src", "/static/svg/qr-hover.svg");
  var mobileIcon = document.getElementById("first-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/qr-hover.svg");
  firstON();
}

// Removes the highligh from the navigation bar
function paymentsOFF() {
  var normalIcon = document.getElementById("fifth-icon");
  normalIcon.setAttribute("src", "/static/svg/records-payments.svg");
  var mobileIcon = document.getElementById("fifth-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/records-payments.svg");
  anyOFF("fifth");
}

// Highlight the icon QR in the navigation Bar
function paymentsON() {
  var normalIcon = document.getElementById("fifth-icon");
  normalIcon.setAttribute("src", "/static/svg/records-payments-hover.svg");
  var mobileIcon = document.getElementById("fifth-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/records-payments-hover.svg");
  fifthON();
}

// Removes the highligh from the navigation bar
function materialsOFF() {
  var normalIcon = document.getElementById("fourth-icon");
  normalIcon.setAttribute("src", "/static/svg/records-materials.svg");
  var mobileIcon = document.getElementById("fourth-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/records-materials.svg");
  anyOFF("fourth");
}

// Highlight the icon QR in the navigation Bar
function materialsON() {
  var normalIcon = document.getElementById("fourth-icon");
  normalIcon.setAttribute("src", "/static/svg/records-materials-hover.svg");
  var mobileIcon = document.getElementById("fourth-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/records-materials-hover.svg");
  fourthON();
}

// Removes the highligh from the navigation bar
function conversationsOFF() {
  var normalIcon = document.getElementById("first-icon");
  normalIcon.setAttribute("src", "/static/svg/conversation.svg");
  var mobileIcon = document.getElementById("first-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/conversation.svg");
  anyOFF("first");
}

// Highlight the icon Wallet in the navigation Bar
function conversationsON() {
  var normalIcon = document.getElementById("first-icon");
  normalIcon.setAttribute("src", "/static/svg/conversation-hover.svg");
  var mobileIcon = document.getElementById("first-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/conversation-hover.svg");
  firstON();
}

// Highlight the icon Access in the navigation Bar
function videosON() {
  var normalIcon = document.getElementById("second-icon");
  normalIcon.setAttribute("src", "/static/svg/video.svg");
  var mobileIcon = document.getElementById("second-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/video.svg");
  secondON();
}

// Removes the highligh from the navigation bar
function videosOFF() {
  var normalIcon = document.getElementById("second-icon");
  normalIcon.setAttribute("src", "/static/svg/video-blue.svg");
  var mobileIcon = document.getElementById("second-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/video-blue.svg");
  anyOFF("second");
}

// Highlight the icon Access in the navigation Bar
function reportsON() {
  var normalIcon = document.getElementById("second-icon");
  normalIcon.setAttribute("src", "/static/svg/report-hover.svg");
  var mobileIcon = document.getElementById("second-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/report-hover.svg");
  secondON();
}

// Removes the highligh from the navigation bar
function reportsOFF() {
  var normalIcon = document.getElementById("second-icon");
  normalIcon.setAttribute("src", "/static/svg/report.svg");
  var mobileIcon = document.getElementById("second-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/report.svg");
  anyOFF("second");
}

// Highlight the icon Villages in the navigation Bar
function villagesON() {
  var normalIcon = document.getElementById("fourth-icon");
  normalIcon.setAttribute("src", "/static/svg/tent-hover.svg");
  var mobileIcon = document.getElementById("fourth-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/tent-hover.svg");
  fourthON();
}

// Removes the highligh from the navigation bar
function villagesOFF() {
  var normalIcon = document.getElementById("fourth-icon");
  normalIcon.setAttribute("src", "/static/svg/tent.svg");
  var mobileIcon = document.getElementById("fourth-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/tent.svg");
  anyOFF("fourth");
}

//Highlight the icon Records in the navigation Bar
function recordsON() {
  var normalIcon = document.getElementById("third-icon");
  normalIcon.setAttribute("src", "/static/svg/records-hover.svg");
  var mobileIcon = document.getElementById("third-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/records-hover.svg");
  thirdON();
}

//Removes the highligh from the navigation bar
function recordsOFF() {
  var normalIcon = document.getElementById("third-icon");
  normalIcon.setAttribute("src", "/static/svg/records.svg");
  var mobileIcon = document.getElementById("third-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/records.svg");
  anyOFF("third");
}

// Highlight the icon Access in the navigation Bar
function accessON() {
  var normalIcon = document.getElementById("second-icon");
  normalIcon.setAttribute("src", "/static/svg/access-hover.svg");
  var mobileIcon = document.getElementById("second-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/access-hover.svg");
  secondON();
}

// Removes the highligh from the navigation bar
function accessOFF() {
  var normalIcon = document.getElementById("second-icon");
  normalIcon.setAttribute("src", "/static/svg/access.svg");
  var mobileIcon = document.getElementById("second-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/access.svg");
  anyOFF("second");
}

// Highlight the icon Users in the navigation Bar
function usersON() {
  var normalIcon = document.getElementById("first-icon");
  normalIcon.setAttribute("src", "/static/svg/worker-hover.svg");
  var mobileIcon = document.getElementById("first-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/worker-hover.svg");
  firstON();
}

//Removes the highligh from the navigation bar
function usersOFF() {
  var normalIcon = document.getElementById("first-icon");
  normalIcon.setAttribute("src", "/static/svg/worker.svg");
  var mobileIcon = document.getElementById("first-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/worker.svg");
  anyOFF("first");
}

//Removes the highligh from the navigation bar
function recordsSalesOFF() {
  var normalIcon = document.getElementById("second-icon");
  normalIcon.setAttribute("src", "/static/svg/records-products.svg");
  var mobileIcon = document.getElementById("second-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/records-products.svg");
  anyOFF("second");
}

//Highlight the icon records-products in the navigation Bar
function recordsSalesON() {
  var normalIcon = document.getElementById("second-icon");
  normalIcon.setAttribute("src", "/static/svg/records-products-hover.svg");
  var mobileIcon = document.getElementById("second-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/records-products-hover.svg");
  secondON();
}

//Highlight the icon Records in the navigation Bar
function recordsOrdersON() {
  var normalIcon = document.getElementById("third-icon");
  normalIcon.setAttribute("src", "/static/svg/records-orders-hover.svg");
  var mobileIcon = document.getElementById("third-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/records-orders-hover.svg");
  thirdON();
}

//Removes the highligh from the navigation bar
function recordsOrdersOFF() {
  var normalIcon = document.getElementById("third-icon");
  normalIcon.setAttribute("src", "/static/svg/records-orders.svg");
  var mobileIcon = document.getElementById("third-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/records-orders.svg");
  anyOFF("third");
}

//Highlight the icon calendar in the navigation Bar
function categoriesON() {
  var normalIcon = document.getElementById("third-icon");
  normalIcon.setAttribute("src", "/static/svg/categories-hover.svg");
  var mobileIcon = document.getElementById("third-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/categories-hover.svg");
  thirdON();
}

//Removes the highligh from the navigation bar
function categoriesOFF() {
  var normalIcon = document.getElementById("third-icon");
  normalIcon.setAttribute("src", "/static/svg/categories.svg");
  var mobileIcon = document.getElementById("third-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/categories.svg");
  anyOFF("third");
}

//Highlight the icon calendar in the navigation Bar
function calendarON() {
  var normalIcon = document.getElementById("first-icon");
  normalIcon.setAttribute("src", "/static/svg/calendar-hover.svg");
  var mobileIcon = document.getElementById("first-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/calendar-hover.svg");
  firstON();
}

//Removes the highligh from the navigation bar
function calendarOFF() {
  var normalIcon = document.getElementById("first-icon");
  normalIcon.setAttribute("src", "/static/svg/calendar.svg");
  var mobileIcon = document.getElementById("first-icon-mobile");
  mobileIcon.setAttribute("src", "/static/svg/calendar.svg");
  anyOFF("first");
}

// Higlight the text and the line of the first icon of the navbar
function firstON() {
  var text = document.getElementById("first-text");
  text.setAttribute("style", "color: #63B190;");
  var section = document.getElementById("first-section");
  section.setAttribute(
    "style",
    "color: #63B190; border-right: 4px solid #63B190;"
  );
  var textMobile = document.getElementById("first-text-mobile");
  textMobile.setAttribute("style", "color: #63B190;");
  var sectionMobile = document.getElementById("first-section-mobile");
  sectionMobile.setAttribute(
    "style",
    "color: #63B190; border-bottom: 4px solid #63B190;"
  );
}

// Higlight the text and the line of the second icon of the navbar
function secondON() {
  var text = document.getElementById("second-text");
  text.setAttribute("style", "color: #ED037C;");
  var section = document.getElementById("second-section");
  section.setAttribute(
    "style",
    "color: #ED037C; border-right: 4px solid #ED037C;"
  );
  var textMobile = document.getElementById("second-text-mobile");
  textMobile.setAttribute("style", "color: #ED037C;");
  var sectionMobile = document.getElementById("second-section-mobile");
  sectionMobile.setAttribute(
    "style",
    "color: #ED037C; border-bottom: 4px solid #ED037C;"
  );
}

// Higlight the text and the line of the second icon of the navbar
function thirdON() {
  var text = document.getElementById("third-text");
  text.setAttribute("style", "color: #A97C50;");
  var section = document.getElementById("third-section");
  section.setAttribute(
    "style",
    "color: #A97C50; border-right: 4px solid #A97C50;"
  );
  var textMobile = document.getElementById("third-text-mobile");
  textMobile.setAttribute("style", "color: #A97C50;");
  var sectionMobile = document.getElementById("third-section-mobile");
  sectionMobile.setAttribute(
    "style",
    "color: #A97C50; border-bottom: 4px solid #A97C50;"
  );
}

//Highlight the fourth icon in the navigation Bar
function fourthON() {
  var text = document.getElementById("fourth-text");
  text.setAttribute("style", "color: #F7941D;");
  var section = document.getElementById("fourth-section");
  section.setAttribute(
    "style",
    "color: #F7941D; border-right: 4px solid #F7941D;"
  );
  var textMobile = document.getElementById("fourth-text-mobile");
  textMobile.setAttribute("style", "color: #F7941D;");
  var sectionMobile = document.getElementById("fourth-section-mobile");
  sectionMobile.setAttribute(
    "style",
    "color: #F7941D; border-bottom: 4px solid #F7941D;"
  );
}

//Highlight the fifth icon in the navigation Bar
function fifthON() {
  var text = document.getElementById("fifth-text");
  text.setAttribute("style", "color: #63B190;");
  var section = document.getElementById("fifth-section");
  section.setAttribute(
    "style",
    "color: #63B190; border-right: 4px solid #63B190;"
  );
  var textMobile = document.getElementById("fifth-text-mobile");
  textMobile.setAttribute("style", "color: #63B190;");
  var sectionMobile = document.getElementById("fifth-section-mobile");
  sectionMobile.setAttribute(
    "style",
    "color: #63B190; border-bottom: 4px solid #63B190;"
  );
}

// Removes the higlight the text and the line of the icon of the navbar given
function anyOFF(number) {
  var text = document.getElementById(number + "-text");
  text.removeAttribute("style");
  var section = document.getElementById(number + "-section");
  section.removeAttribute("style");
  var textMobile = document.getElementById(number + "-text-mobile");
  textMobile.removeAttribute("style");
  var sectionMobile = document.getElementById(number + "-section-mobile");
  sectionMobile.removeAttribute("style");
}

// WEBSITE

// Highlight the icon Performance in the navigation Bar
function aboutON() {
  var normalIcon = document.getElementById("first-icon");
  normalIcon.setAttribute("src", "/static/svg/about-hover.svg");
  var section = document.getElementById("first-section");
  section.setAttribute(
    "style",
    "color: #63B190; border-bottom: 2px solid #63B190; box-sizing: ;"
  );
  var textMobile = document.getElementById("first-text");
  textMobile.setAttribute("style", "color: #63B190;");
}

// Removes the highligh from the navigation bar
function aboutOFF() {
  var normalIcon = document.getElementById("first-icon");
  normalIcon.setAttribute("src", "/static/svg/about.svg");
  var text = document.getElementById("first-text");
  text.removeAttribute("style");
  var section = document.getElementById("first-section");
  section.removeAttribute("style");
}

// Highlight the icon Performance in the navigation Bar
function designON() {
  var normalIcon = document.getElementById("second-icon");
  normalIcon.setAttribute("src", "/static/svg/design-hover.svg");
  var section = document.getElementById("second-section");
  section.setAttribute(
    "style",
    "color: #ED037C; border-bottom: 2px solid #ED037C; box-sizing: ;"
  );
  var textMobile = document.getElementById("second-text");
  textMobile.setAttribute("style", "color: #ED037C;");
}

// Removes the highligh from the navigation bar
function designOFF() {
  var normalIcon = document.getElementById("second-icon");
  normalIcon.setAttribute("src", "/static/svg/design.svg");
  var text = document.getElementById("second-text");
  text.removeAttribute("style");
  var section = document.getElementById("second-section");
  section.removeAttribute("style");
}

// Highlight the icon Performance in the navigation Bar
function featuresON() {
  var normalIcon = document.getElementById("third-icon");
  normalIcon.setAttribute("src", "/static/svg/features-hover.svg");
  var section = document.getElementById("third-section");
  section.setAttribute(
    "style",
    "color: #F7941D; border-bottom: 2px solid #F7941D; box-sizing: ;"
  );
  var textMobile = document.getElementById("third-text");
  textMobile.setAttribute("style", "color: #F7941D;");
}

// Removes the highligh from the navigation bar
function featuresOFF() {
  var normalIcon = document.getElementById("third-icon");
  normalIcon.setAttribute("src", "/static/svg/features.svg");
  var text = document.getElementById("third-text");
  text.removeAttribute("style");
  var section = document.getElementById("third-section");
  section.removeAttribute("style");
}

// Highlight the icon Performance in the navigation Bar
function manualsON() {
  var normalIcon = document.getElementById("fourth-icon");
  normalIcon.setAttribute("src", "/static/svg/manuals-hover.svg");
  var section = document.getElementById("fourth-section");
  section.setAttribute(
    "style",
    "color: #A97C50; border-bottom: 2px solid #A97C50; box-sizing: ;"
  );
  var textMobile = document.getElementById("fourth-text");
  textMobile.setAttribute("style", "color: #A97C50;");
}

// Removes the highligh from the navigation bar
function manualsOFF() {
  var normalIcon = document.getElementById("fourth-icon");
  normalIcon.setAttribute("src", "/static/svg/manuals.svg");
  var text = document.getElementById("fourth-text");
  text.removeAttribute("style");
  var section = document.getElementById("fourth-section");
  section.removeAttribute("style");
}

// Highlight the icon Performance in the navigation Bar
function demoON() {
  var normalIcon = document.getElementById("back-icon");
  normalIcon.setAttribute("src", "/static/svg/demo-hover.svg");
  var section = document.getElementById("back-section");
  section.setAttribute(
    "style",
    "color: #CF1F26; border-bottom: 2px solid #CF1F26; box-sizing: ;"
  );
  var textMobile = document.getElementById("back-text");
  textMobile.setAttribute("style", "color: #CF1F26;");
}

// Removes the highligh from the navigation bar
function demoOFF() {
  var normalIcon = document.getElementById("back-icon");
  normalIcon.setAttribute("src", "/static/svg/demo.svg");
  var text = document.getElementById("back-text");
  text.removeAttribute("style");
  var section = document.getElementById("back-section");
  section.removeAttribute("style");
}

function wait(ms) {
  var start = new Date().getTime();
  var end = start;
  while (end < start + ms) {
    end = new Date().getTime();
  }
}
