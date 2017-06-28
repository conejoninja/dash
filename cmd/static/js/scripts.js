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


    $.getJSON( API_PROTO + "://" + API_URI + ":" + API_PORT +"/sensor/food01-t1")
        .success(function( data ) {


            var gdata = [];
            var gdata2 = [];
            var l = data.length;
            var i = 0;
            for(i=0;i<l;i++) {
                var dtime = Date.parse(data[i].time);
                gdata.push({
                    x: dtime,
                    y: data[i].value
                });
                gdata2.push({
                    x: dtime + 400,
                    y: data[l - i - 1].value
                });
            }


            var scatterChartData = {
                datasets: [{
                    label: "Current",
                    borderColor: window.chartColors.red,
                    backgroundColor: color(window.chartColors.red).alpha(0.2).rgbString(),
                    data: gdata
                }, {
                    label: "Past",
                    borderColor: window.chartColors.blue,
                    backgroundColor: color(window.chartColors.blue).alpha(0.2).rgbString(),
                    data: gdata2
                }]
            };

            var ctx = document.getElementById("line").getContext("2d");
            window.myScatter = Chart.Scatter(ctx, {
                data: scatterChartData,
                options: {
                    scales: {
                        xAxes: [{
                            type: 'time',
                        }]
                    },
                }
            });


        })
        .fail(function( data ) {
            console.log( "FAILE: " , data );
        });

    $.getJSON( API_PROTO + "://" + API_URI + ":" + API_PORT +"/sensor/food01-h1")
        .success(function( data ) {


            var gdata = [];
            var gdata2 = [];
            var l = data.length;
            var i = 0;
            for(i=0;i<l;i++) {
                var dtime = Date.parse(data[i].time);
                gdata.push({
                    x: dtime,
                    y: data[i].value
                });
                gdata2.push({
                    x: dtime + 400,
                    y: data[l - i - 1].value
                });
            }


            var scatterChartData = {
                datasets: [{
                    label: "Current",
                    borderColor: window.chartColors.red,
                    backgroundColor: color(window.chartColors.red).alpha(0.2).rgbString(),
                    data: gdata
                }, {
                    label: "Past",
                    borderColor: window.chartColors.blue,
                    backgroundColor: color(window.chartColors.blue).alpha(0.2).rgbString(),
                    data: gdata2
                }]
            };

            var ctx2 = document.getElementById("line2").getContext("2d");
            window.myScatter = Chart.Scatter(ctx2, {
                data: scatterChartData,
                options: {
                    scales: {
                        xAxes: [{
                            type: 'time',
                        }]
                    },
                }
            });


        })
        .fail(function( data ) {
            console.log( "FAILE: " , data );
        });


    $.getJSON( API_PROTO + "://" + API_URI + ":" + API_PORT +"/meta/food01-t1")
        .success(function( data ) {

            console.log(data);


        })
        .fail(function( data ) {
            console.log( "FAILE: " , data );
        });


})(jQuery);

                     
     
  