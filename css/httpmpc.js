// $("#play").click(load);
var ticktock

$( document ).ready( function () { tickMaster(); ticktock = setInterval(tickMaster, 1000);})

function tickMaster() {
  tickStatus()
  tickCurrentSong()
  tickPlaylistInfo()
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


var songId = -1
function tickCurrentSong() {
  $.getJSON( "/currentsong", function( data ) {
    if ("Artist" in data && "Title" in data) {
      $("title").html("HttpMpc :: " + data.Artist + " - " + data.Title)
    }
    songId = -1
    if ("Id" in data) {
      console.log(data.Id)
      songid = data.Id;
    }
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

function tickPlaylistInfo() {
  $.getJSON( "/playlistinfo", function( data ) {
    $("#queue").html("<thead><tr><th>ID</th><th>Track</th><th>Title</th><th>Artist</th><th>Album</th><th>Length</th></tr></thead>");
    $.each( data, function( key, val ) {
      if (val.Id == songid) {
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

var currDir = "/"
function browseTo(subdir) {

}