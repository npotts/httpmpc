(function(){
  $(window).scroll(function () {
      var top = $(document).scrollTop();
      $('.splash').css({
        'background-position': '0px -'+(top/3).toFixed(2)+'px'
      });
      if(top > 50)
        $('#home > .navbar').removeClass('navbar-transparent');
      else
        $('#home > .navbar').addClass('navbar-transparent');
  });

  $("a[href='#']").click(function(e) {
    e.preventDefault();
  });

  var $button = $("<div id='source-button' class='btn btn-primary btn-xs'>&lt; &gt;</div>").click(function(){
    var html = $(this).parent().html();
    html = cleanSource(html);
    $("#source-modal pre").text(html);
    $("#source-modal").modal();
  });

  $('.bs-component [data-toggle="popover"]').popover();
  $('.bs-component [data-toggle="tooltip"]').tooltip();

  $(".bs-component").hover(function(){
    $(this).append($button);
    $button.show();
  }, function(){
    $button.hide();
  });

  function cleanSource(html) {
    html = html.replace(/×/g, "&times;")
               .replace(/«/g, "&laquo;")
               .replace(/»/g, "&raquo;")
               .replace(/←/g, "&larr;")
               .replace(/→/g, "&rarr;");

    var lines = html.split(/\n/);

    lines.shift();
    lines.splice(-1, 1);

    var indentSize = lines[0].length - lines[0].trim().length,
        re = new RegExp(" {" + indentSize + "}");

    lines = lines.map(function(line){
      if (line.match(re)) {
        line = line.substring(indentSize);
      }

      return line;
    });

    lines = lines.join("\n");

    return lines;
  }

})();

$("#play").click(load);
// $( document ).ready( loadQueue )

function load() {
  // /playlistinfo
  $.getJSON( "/playlistinfo", function( data ) {
    var items = [];
    $("#queue").html("<thead><tr><th>ID</th><th>Track</th><th>Title</th><th>Artist</th><th>Album</th><th>Length</th></tr></thead>");
    $.each( data, function( key, val ) {
      $("#queue").append("<tr><td>" + key + "</td><td>" + val.Track + "</td><td>" + val.Title + "</td><td>" + val.Artist + "</td><td>" + val.Album + "</td><td>" + val.Time + "</td></tr>\n")
    });
  });

  $.getJSON( "/listplaylists", function( data ) {
    var items = [];
    $("#playlists").html("");
    $.each( data, function( key, val ) {
      $("#playlists").append("<tr><td>" + val.playlist + "</td></tr>\n");
    });
  });

}