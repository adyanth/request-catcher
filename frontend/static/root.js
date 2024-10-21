window.root = window.root || {};

$(document).ready(function() {
  $('#root-host').text(window.location.host);
  $('#root-protocol').text(window.location.protocol);
  $('#new-catcher').submit(function(e) {
    e.preventDefault();

    var subdomain = $('#subdomain').val();
    if (subdomain) {
      var url = window.location.protocol + '//' + subdomain + '.' + window.location.host + '/';
      window.location = url;
    }
  });

  // See https://requestcatcher.com/assets/its-free-software.gif

  $('#subdomain').focus();
});
