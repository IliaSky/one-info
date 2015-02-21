//======================================== VIEWS ========================================
SG.VIEWS = {
  NAVIGATION: '*a(href:#/@@)[@@]',

  FOLDER: "table > tr{ th[From] + th[Subject] } + *tr > @(headers) > td[@From] + td[@Subject] ",

  MODULE: ".module-outer > .module-inner > section.module > h1[@name] + @(content) > table > *tr > *td[@@]",
  MODULES: "* $.MODULE",

  // Authentication forms
  USERNAME: "label(for:username)[Username] + input#username(type:text,required)",
  PASSWORD: "label(for:password)[Password] + input#password(type:password,required)",
  PASSWORD_AGAIN: "label(for:password-again)[Password Again] + input#password-again(type:password,required)",
  REMEMBER: "label(for:remember)[Remember me] + input#remember(type:checkbox)",
  SUBMIT: "input(type:submit)",

  LOGIN :    "form#login-form    > $.USERNAME + $.PASSWORD + $.REMEMBER + $.SUBMIT",
  REGISTER : "form#register-form > $.USERNAME + $.PASSWORD + $.PASSWORD_AGAIN + $.REMEMBER + $.SUBMIT",
};
function $(id) {
  return document.getElementById(id);
}

function load(pattern, data) {
  $('content').innerHTML = SG('$.' + pattern.toUpperCase(), data);
}

function on(eventName, selector, callback){
  document.addEventListener(eventName, function(e){
  // document['on' + eventName] = function (e) {

    var elements = [].slice.call(document.querySelectorAll(selector));
    if (elements.indexOf(e.target) != -1)
      return callback.call(e.target, e);
  });
}
function ajax_get(url, success, failure){
  var xhr = window.XMLHttpRequest ? new XMLHttpRequest() : new ActiveXObject("Microsoft.XMLHTTP");
  xhr.onreadystatechange = function(){
    if (xhr.readyState == 4)
      xhr.status == 200 ? success(xhr.responseText) : failure(xhr.responseText);
  };
  xhr.open("GET", url, true);
  xhr.send();
}
function range(x) {
  var result = [];
  for (var i = 0; i < x; i++) {
    result[i] = [i]
  }
  return result;
}
// $('content').innerHTML = 
var content = range(5).map(function(e, i){
  return range(3).map(function(x){ return e*x; })
});
load('modules', ['Work', 'Go', 'JS', 'WWW', 'Uni'].map(function(e){
  return {name: e, content: content};
}));
$('sizes').innerHTML = ".module {width: 20%;height:400px;}";

function setSizes(obj) {
  $('sizes').innerHTML = '.module-outer {' + Object.keys(obj).map(function(key){
    return key + ':' + obj[key];
  }).join(';') + '}';
}

$('size').onchange = function (e){
  var grid = this.value.split('x'), x = grid[0]|0 || 1, y=grid[1]|0 || 1;
  // $('sizes').innerHTML = ".module {width: " + 90/x + "%;height: " + 90/y + "%;}";
  setSizes({
    width: 100/x + '%',
    height: 100/y + '%',
  });
};
$('size').onchange()
//======================================== MAIN ========================================

// window.onhashchange = function () {
//   var m, url = location.hash.slice(2);
//   console.log('Getting /api/folders/'+url);
//   ajax_get('/api/folders/' + url, function(data){
//     console.log(JSON.parse(data));
//     load('folder', JSON.parse(data));
//   }, function(err) {
//     alert(err);
//   })
// };
// window.onhashchange();

// ajax_get('/api/folders', function(data){
//   $('side-nav').innerHTML = SG('$.NAVIGATION', JSON.parse(data));
// }, function(){alert('Could not connect :(');})