window.root=window.root||{},$(document).ready(function(){$("#root-host").text(window.location.host),$("#root-protocol").text(window.location.protocol),$("#new-catcher").submit(function(o){o.preventDefault();var t=$("#subdomain").val();if(t){var n=window.location.protocol+"//"+t+"."+window.location.host+"/";window.location=n}}),$("#subdomain").focus()});
//# sourceMappingURL=root.1764eeb1.js.map
