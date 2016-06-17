// $("#play").click(load);
var ticktock

$( document ).ready( function () { tickMaster(); ticktock = setInterval(tickMaster, 1000);});
$( document ).ready( browseTo("") );
$( document ).ready( function () { $("#music tr").click(musicDecend);});
$( document ).ready( function () { $("#home").click(musicDecend);});
$( document ).ready( function () { $("#updir").click(musicDecend);});
//button wiring
$( document ).ready( function () { $("#play").click(function() {$.post("/play/-1");})});
// $( document ).ready( function () { $("#pause").click(function() {$.post("/play/-1");})});
$( document ).ready( function () { $("#stop").click(function() {$.post("/stop");})});
$( document ).ready( function () { $("#previous").click(function() {$.post("/previous");})});
$( document ).ready( function () { $("#next").click(function() {$.post("/next");})});
$( document ).ready( function () { $("#clear").click(function() {$.post("/clear");})});
// $( document ).ready( function () { $("#consume").click(function() {$.post("/play/-1");})});
// $( document ).ready( function () { $("#repeat").click(function() {$.post("/play/-1");})});
// $( document ).ready( function () { $("#random").click(function() {$.post("/play/-1");})});
// $( document ).ready( function () { $("#single").click(function() {$.post("/play/-1");})});


function tickMaster() {
  tickStatus()
  tickCurrentSong()
  tickPlaylists()
}



function tickStatus() {
  $.getJSON( "/status", function( data ) {
    $.each( ["consume", "repeat", "random", "single"], function(i, val) {
      if (val in data) {
        var item = $(document.getElementById(val))
        item.removeClass("btn-warning");
        item.removeClass("btn-primary");
        item.removeClass("btn-info");
        item.addClass( (data[val] == 1) ? "btn-primary" : "btn-info");
      }
    });
    if ("state" in data && "time" in data && data.state == "play") {
      times = data.time.split(":")
      $("#progress").attr("style", "width:" + 100*parseInt(times[0])/parseInt(times[1]) + "%")
    }
  });
}



function tickCurrentSong() {
  $.getJSON( "/currentsong", function( data ) {
    if ("Artist" in data && "Title" in data) {
      $("title").html("HttpMpc :: " + data.Artist + " - " + data.Title)
    }
    var songid = -1
    if ("Id" in data) {
      songid = data.Id;
    }
    tickPlaylistInfo(songid)
  });
}

function humanizeTime(secs) {
  secs = parseInt(secs)
  var hours = Math.floor(secs/3600);
  secs %= 3600;
  var mins = Math.floor(secs / 60);
  secs %= 60;
  var rtn = ""
  if (hours > 0) {
    rtn += hours + "h" 
  }
  if (rtn.length >0 || mins > 0) {
    rtn += mins + "m"
  }
  if (rtn.length > 0 || secs > 0) {
    rtn += secs + "s"
  }
  return rtn
}

function tickPlaylistInfo(songid) {
  $.getJSON( "/playlistinfo", function( data ) {
    $("#queue").html("<thead><tr><th>ID</th><th>Track</th><th>Title</th><th>Artist</th><th>Album</th><th>Length</th></tr></thead>");
    $.each( data, function( key, val ) {
      if (songid != undefined && val.Id == songid) {
        $("#queue").append("<tr class='info'><td>" + val.Id + "</td><td>" + parseInt(val.Track) + "</td><td>" + val.Title + "</td><td>" + val.Artist + "</td><td>" + val.Album + "</td><td>" + humanizeTime(val.Time) + "</td></tr>\n")
      } else {
        $("#queue").append("<tr><td>" + val.Id + "</td><td>" + parseInt(val.Track) + "</td><td>" + val.Title + "</td><td>" + val.Artist + "</td><td>" + val.Album + "</td><td>" + humanizeTime(val.Time) + "</td></tr>\n")
      }      
    });
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
      if ("directory" in val) {
        //
        $("#music").append("<tr><td><span class=\"ion-plus-circled\"></span>&nbsp;<span class=\"ion-folder\"></span>&nbsp;<span>" + val.directory + "</span></td></tr>\n");
      }
      if ("file" in val) {
        $("#music").append("<tr><td><span class=\"ion-music-note\"></span>&nbsp;<span>" + val.file + "</span></td></tr>\n");
      }
    });
    $("#music span").click(musicDecend);
  });
}

function musicDecend(path) {
  //figure out what they clicked on:
  if ($(this).attr('id') == "home") { browseTo(""); return; }
  if ($(this).attr('id') == "updir") {if ($("#path").html() != "") {sp = $("#path").html().split("/"); sp.splice(sp.length-1, 1);browseTo(sp.join("/"));}return;}
  if ($(this).attr("class") == "ion-plus-circled"){
    var uri = $(this).next().next().html();
    $.post("/add/" + uri);
    console.log("Adding Dir: " + uri);
    return;
  } //adding a subdir
  if ($(this).attr("class") == "ion-folder") {browseTo($(this).next().html());return;} //clicks to folder icon
  if ($(this).attr("class") == undefined && $(this).prev().attr("class") == "ion-folder") { //clicks to folder
    console.log("Decending into: " + $(this).html())
    browseTo($(this).html())
    return
  }
  console.log("Adding file '" + $(this).html() + "'")
  $.post("/add/" + $(this).html());
}