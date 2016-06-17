// $("#play").click(load);
var ticktock
var songId = -1
var currDir = "/"

$( document ).ready( function () { tickMaster(); ticktock = setInterval(tickMaster, 1000);});
$( document ).ready( browseTo("/") );
$( document ).ready( function () { $("#music tr").click(musicDecend);});
$( document ).ready( function () { $("#home").click(musicDecend);});

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



function tickCurrentSong() {
  $.getJSON( "/currentsong", function( data ) {
    if ("Artist" in data && "Title" in data) {
      $("title").html("HttpMpc :: " + data.Artist + " - " + data.Title)
    }
    songId = -1
    if ("Id" in data) {
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

function browseTo(path) {
  $.getJSON( "/listinfo/" + path, function( data ) {
    $("#home").html("&nbsp;" + path);
    $("#music").html("");
    $.each( data, function( id, val ) {
      if ("directory" in val) {
        //
        $("#music").append("<tr><td><span class=\"ion-plus-circled\"></span>&nbsp;<span class=\"ion-folder\">&nbsp;" + val.directory + "</span></td></tr>\n");
      }
      if ("file" in val) {
        $("#music").append("<tr><td></span>&nbsp;<span class=\"ion-music-note\">&nbsp;" + val.file + "</span></td></tr>\n");
      }
    });
    $("#music span").click(musicDecend);
  });
}

function musicDecend(path) {
  //figure out what they clicked on:
  // - the + button
  //   either add whole directory, or single file
  if ($(this).attr('id') == "home") { 
    browseTo("/") 
    return
  }
  if ($(this).attr("class") == "ion-plus-circled"){
    var sibling = $(this).next()
    var child = $(sibling).children()
    console.dir(sibling)
    console.dir(child)
    var cousin = $(sibling).find("span")
    console.log("Adding")
    
    return
  }
  if ($(this).attr("class") == "ion-folder") {
    console.log("Decending")
    browseTo($(this).html().replace("&nbsp;",""))
    return
  }
  console.log("Adding file '" + $(this).html().replace("&nbsp;","") + "'")
  
  //- the file/dir name
  //  either decend into the folder, or add the file
  //get the folder clicked on
  // currDir += $(this).html()
  
  // console.log("Registering click in music table")
  // console.log($(this).index())
  // console.log()
  // console.dir($(this))
//   $(document).ready(function(){
//     $("#myTable td").click(function() {     
 
//         var column_num = parseInt( $(this).index() ) + 1;
//         var row_num = parseInt( $(this).parent().index() )+1;    
 
//         $("#result").html( "Row_num =" + row_num + "  ,  Rolumn_num ="+ column_num );   
//     });
// });
}