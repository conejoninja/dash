(function() {
    "use strict";

    function log(msg) {
        document.getElementById('dash-log').textContent = msg + '\n' + document.getElementById('dash-log').textContent;
    }

    // setup websocket with callbacks
    var ws = new WebSocket('ws://<%= ws_server %>:<%= ws_port %>/ws');
    ws.onopen = function() {
        log('CONNECTED WEBSOCKET');
    };
    ws.onclose = function() {
        log('DISCONNECTED WEBSOCKET');
    };
    ws.onmessage = function(event) {
        console.log(event);
        log(event.data);
    };

})(jQuery);
