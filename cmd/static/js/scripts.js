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


    apiRequest("/sensor/food01-t1;food01-h1")
        .success(function( data ) {

            for(var sensor in data) {
                var current_data = [];
                var past_data = [];
                var l = data[sensor]["current"].length;
                var i = 0;
                for(i=0;i<l;i++) {
                    var dtime = new Date(Date.parse(data[sensor]["current"][i].time));
                    current_data.push({
                        x: dtime,
                        y: data[sensor]["current"][i].value
                    });
                }
                var l = data[sensor]["past"].length;
                var i = 0;
                for(i=0;i<l;i++) {
                    var dtime = new Date(Date.parse(data[sensor]["past"][i].time) + 24 * 3600000);
                    past_data.push({
                        x: dtime,
                        y: data[sensor]["past"][i].value
                    });
                }

                var scatterChartData = {
                    datasets: [{
                        label: "Current",
                        borderColor: window.chartColors.red,
                        backgroundColor: color(window.chartColors.red).alpha(0.2).rgbString(),
                        data: current_data
                    }, {
                        label: "Past",
                        borderColor: window.chartColors.blue,
                        backgroundColor: color(window.chartColors.blue).alpha(0.2).rgbString(),
                        data: past_data
                    }]
                };

                var ctx = document.getElementById(sensor).getContext("2d");
                window.myScatter = Chart.Scatter(ctx, {
                    data: scatterChartData,
                    options: {
                        scales: {
                            xAxes: [{
                                type: 'time',
                            }]
                        },
                        tooltips: {
                            callbacks: {
                                label: function(tooltipItem, data) {
                                    var t = Date.parse(data["datasets"][tooltipItem.datasetIndex]["data"][tooltipItem.index].x);
                                    if(tooltipItem.datasetIndex==1) {
                                        t = t - 24 * 3600000;
                                    }
                                    return "(" +
                                        new Date(t).toString().substr(4,17) +
                                        ", " +
                                        data["datasets"][tooltipItem.datasetIndex]["data"][tooltipItem.index].y +
                                        ")";
                                }
                            },
                        },
                    }
                });
            }
        });

    apiRequest("/meta/food01-t1;food01-h1")
        .success(function( data ) {
            if(data["food01-t1"]!=undefined) {
                if(data["food01-t1"]["max"]!=undefined) {
                    document.getElementById("food01-t1-meta-max").innerHTML = data["food01-t1"]["max"].toFixed(1);
                }
                if(data["food01-t1"]["min"]!=undefined) {
                    document.getElementById("food01-t1-meta-min").innerHTML = data["food01-t1"]["min"].toFixed(1);
                }
            }
            if(data["food01-h1"]!=undefined) {
                if(data["food01-h1"]["avg"]!=undefined) {
                    document.getElementById("food01-h1-meta-avg").innerHTML = data["food01-h1"]["avg"].toFixed(1);
                }
            }
          console.log("META", data);
        });

})(jQuery);


function apiRequest(req) {
    return $.getJSON( API_PROTO + "://" + API_URI + ":" + API_PORT + req)
        .fail(function( data ) {
        console.log( "[API REQUEST FAILED]" , req, data );
    });
}
     
  