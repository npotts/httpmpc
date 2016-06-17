// $("#play").click(load);
var ticktock

$( document ).ready( function () { tickMaster(); ticktock = setInterval(tickMaster, 1000);});
$( document ).ready( browseTo("") );
$( document ).ready( function () { $("#music tr").click(musicDecend);});
$( document ).ready( function () { $("#home").click(musicDecend);});
$( document ).ready( function () { $("#updir").click(musicDecend);});
//button wiring
$( document ).ready( function () { $("#play").click(function() {$.post("/play/-1");})});
$( document ).ready( function () { $("#pause").click(ctrlBool) });
$( document ).ready( function () { $("#stop").click(function() {$.post("/stop");})});
$( document ).ready( function () { $("#previous").click(function() {$.post("/previous");})});
$( document ).ready( function () { $("#next").click(function() {$.post("/next");})});
$( document ).ready( function () { $("#clear").click(function() {$.post("/clear");})});
$( document ).ready( function () { $("#consume").click(ctrlBool) });
$( document ).ready( function () { $("#repeat").click(ctrlBool) });
$( document ).ready( function () { $("#random").click(ctrlBool) });
$( document ).ready( function () { $("#single").click(ctrlBool) });

function humanizeTime(secs) {
  secs = parseInt(secs)
  var rtn = ""
  var hours = Math.floor(secs/3600); secs %= 3600;
  var mins = Math.floor(secs / 60); secs %= 60;
  if (hours > 0) { rtn += hours + "h" }
  if (rtn.length >0 || mins > 0) { rtn += mins + "m" }
  if (rtn.length > 0 || secs > 0) { rtn += secs + "s" }
  return rtn
}

function ctrlBool() {
  handle = $(this).attr("id")
  $.ajax("/" + handle, {method: ($("#" + handle).hasClass("btn-primary")) ? "PUT": "DELETE"});
}

function tickMaster() {
  tickStatus()
  tickCurrentSong()
  tickPlaylists()
}

function tickStatus() {
  $.getJSON( "/status", function( data ) {
    $.each( ["consume", "repeat", "random", "single"], function(i, val) {
      if (val in data) {
        $(document.getElementById(val)).removeClass("btn-warning").removeClass("btn-primary").removeClass("btn-info").addClass( (data[val] == 1) ? "btn-info" : "btn-primary");} });
    if ("state" in data) {
      if ("time" in data && data.state == "play") { times = data.time.split(":"); $("#progress").attr("style", "width:" + 100*parseInt(times[0])/parseInt(times[1]) + "%") }; //progress bar
      $("#pause").removeClass("btn-info").removeClass("btn-primary").addClass((data.state == "pause") ? "btn-info": "btn-primary");
    }
  });
}

function tickCurrentSong() {
  $.getJSON( "/currentsong", function( data ) {
    if ("Artist" in data && "Title" in data) { $("title").html("HttpMpc :: " + data.Artist + " - " + data.Title) }
    if ("Id" in data) { tickPlaylistInfo(data.Id) } else { tickPlaylistInfo(-1) }
  });
}

function tickPlaylistInfo(songid) {
  $.getJSON( "/playlistinfo", function( data ) {
    $("#queue").html("<thead><tr><th>Title</th><th>Artist</th><th>Album</th><th>Length</th></tr></thead>");
    $.each( data, function( key, val ) {
      if (songid != undefined && val.Id == songid) {
        $("#queue").append("<tr class='info' sid='" + val.Id + "'><td>" + val.Title + "</td><td>" + val.Artist + "</td><td>" + val.Album + "</td><td>" + humanizeTime(val.Time) + "</td></tr>\n")
      } else {
        $("#queue").append("<tr sid='" + val.Id + "'><td>" + val.Title + "</td><td>" + val.Artist + "</td><td>" + val.Album + "</td><td>" + humanizeTime(val.Time) + "</td></tr>\n")
      }
    });
    $("#queue tr").click(function() {$.post("/playid/" + $(this).attr("sid"));}); //want to play a song in the queue
  });
}

function tickPlaylists() {
  $.getJSON( "/listplaylists", function( data ) {
    $("#playlists").html("");
    $.each( data, function( key, val ) {
      $("#playlists").append("<tr><td>" + val.playlist + "</td></tr>\n");
    });
  });
}

function browseTo(path) {
  $.getJSON( "/listinfo/" + path, function( data ) {
    $("#path").html(path);
    $("#music").html("");
    $.each( data, function( id, val ) {
      if ("directory" in val) { $("#music").append("<tr><td><span class=\"ion-plus-circled\"></span>&nbsp;<span class=\"ion-folder\"></span>&nbsp;<span>" + val.directory + "</span></td></tr>\n"); }
      if ("file" in val) { $("#music").append("<tr><td><span class=\"ion-music-note\"></span>&nbsp;<span>" + val.file + "</span></td></tr>\n"); }
    });
    $("#music span").click(musicDecend); //reregister evt handler
  });
}

function musicDecend() { //figure out what they clicked on
  if ($(this).attr('id') == "home") { browseTo(""); return; }
  if ($(this).attr('id') == "updir") {if ($("#path").html() != "") {sp = $("#path").html().split("/"); sp.splice(sp.length-1, 1);browseTo(sp.join("/"));}return;}
  if ($(this).attr("class") == "ion-plus-circled"){ var uri = $(this).next().next().html(); $.post("/add/" + uri); return; } //adding a subdir
  if ($(this).attr("class") == "ion-folder") {browseTo($(this).next().html());return;} //clicks to folder icon
  if ($(this).attr("class") == undefined && $(this).prev().attr("class") == "ion-folder") { browseTo($(this).html());return;}//clicks to folder
  $.post("/add/" + $(this).html()); //otherwise it was jsut a file
}