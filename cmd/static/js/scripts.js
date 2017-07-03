window.chartColors = {
    red: 'rgb(255, 99, 132)',
    orange: 'rgb(255, 159, 64)',
    yellow: 'rgb(255, 205, 86)',
    green: 'rgb(75, 192, 192)',
    blue: 'rgb(54, 162, 235)',
    purple: 'rgb(153, 102, 255)',
    grey: 'rgb(201, 203, 207)'
};

var color = Chart.helpers.color;


(function() {
    "use strict";

    // custom scrollbar

    $("html").niceScroll({styler:"fb",cursorcolor:"#1b93e1", cursorwidth: '6', cursorborderradius: '10px', background: '#FFFFFF', spacebarenabled:false, cursorborder: '0',  zindex: '1000'});

    $(".scrollbar1").niceScroll({styler:"fb",cursorcolor:"#1b93e1", cursorwidth: '6', cursorborderradius: '0',autohidemode: 'false', background: '#FFFFFF', spacebarenabled:false, cursorborder: '0'});

	
	
    $(".scrollbar1").getNiceScroll();
    if ($('body').hasClass('scrollbar1-collapsed')) {
        $(".scrollbar1").getNiceScroll().hide();
    }

})(jQuery);


function apiRequest(req) {
    return $.getJSON( AJAXURI + req)
        .fail(function( data ) {
            console.log( "[API REQUEST FAILED]" , req, data );
            jsAlert("API request failed", "danger");
        });
}

function callRequest(device, f, args) {
    return $.post( AJAXURI + "/" + device + "/" + f, args)
        .success(function( data ) {
            data = JSON.parse(data);
            jsAlert(data.message, data.type);
            console.log( "[CALL REQUEST SUCCESS]" , device, f, args, data );
        })
        .fail(function( data ) {
            jsAlert("API request failed", "danger");
            console.log( "[CALL REQUEST FAILED]" , device, f, args, data );
        });
}

function jsAlert(message, type) {
    if (type == "error") {
        type = "danger";
    }
    if (type != "info" && type != "warning" && type != "danger") {
        type = "success";
    }
    var alertDiv = document.getElementById("js-alert");
    alertDiv.className = "alert alert-" + type;
    alertDiv.innerHTML = message;
}