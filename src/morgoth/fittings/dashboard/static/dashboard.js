
function Metrics() {
    this.series = [];
    this.graph = undefined;
    this.requests = [];
}

Metrics.prototype.draw = function () {

    var that = this;
    var series = this.series;

    var num_points = 250;
    var start = new Date();
    var start_date = $('#start_date').val();
    if (start_date) {
        start = new Date(start_date);
    }
    else {
        start.setDate(start.getDate() - 1);
    }

    var stop = new Date();
    var stop_date = $('#stop_date').val();
    if (stop_date) {
        stop = new Date(stop_date);
    }

    var step = ((stop.getTime() - start.getTime()) / 1000 ) / num_points;


    var metric_pattern = $('#metric_pattern').val();

    var palette = new Rickshaw.Color.Palette();

    that.json("http://localhost:8001/metrics?pattern=" + metric_pattern, function(data) {
        metrics = data.metrics
        metrics.forEach(function(metric, index) {
            var url = "http://localhost:8001/data/"
                    + metric
                    + '?start=' + start.toString()
                    + '&stop=' + stop.toString()
                    + '&step=' + step + 's';
            that.json(url, function(rows) {
                var data = [];
                if (rows.data.length > 0) {
                    rows.data.forEach(function(r) {
                        data.push({x: new Date(r[0]).getTime() / 1000, y: r[1]});
                    });
                    series[index] = {
                        data: data,
                        name: metric,
                        color: palette.color(),
                    };
                    that.update();
                } else {
                    console.log('no data received for ' + metric);
                }
            });
        });
    });
}

Metrics.prototype.redraw = function () {
    $('#legend').empty();
    $('#chart_container').html(
        '<div id="y_axis"></div><div id="chart"></div><div id="slider"></div>'
    );
    this.series = [];
    this.graph = undefined;
    this.abort_all();
    this.draw();
}

Metrics.prototype.update = function() {

    if (this.graph == undefined) {
        this._init_graph();
    }
    var graph = this.graph;
    graph.update();
    $('#legend').empty().css('height', '');
    var legend = new Rickshaw.Graph.Legend( {
        graph: graph,
        element: document.getElementById('legend')

    } );

    var shelving = new Rickshaw.Graph.Behavior.Series.Toggle( {
        graph: graph,
        legend: legend
    } );

    var order = new Rickshaw.Graph.Behavior.Series.Order( {
        graph: graph,
        legend: legend
    } );

    var highlighter = new Rickshaw.Graph.Behavior.Series.Highlight( {
        graph: graph,
        legend: legend
    } );
}

Metrics.prototype._init_graph = function() {

    var series = this.series;
    var graph = new Rickshaw.Graph( {
                element: document.querySelector("#chart"),
                width: 1000,
                height: 500,
                renderer: 'line',
                interpolation: 'linear',
                series: series,
    } );
    this.graph = graph;

    var x_axis = new Rickshaw.Graph.Axis.Time( {
        graph: graph,
        timeFixture: new Rickshaw.Fixtures.Time.Local(),
    } );

    var y_axis = new Rickshaw.Graph.Axis.Y( {
            graph: graph,
            orientation: 'left',
            tickFormat: Rickshaw.Fixtures.Number.formatKMBT,
            element: document.getElementById('y_axis'),
    } );


    graph.render();

    var smoother = new Rickshaw.Graph.Smoother( {
        graph: graph,
        element: $('#smoother')
    } );

    var preview = new Rickshaw.Graph.RangeSlider.Preview( {
        graph: graph,
        element: document.getElementById('slider'),
    } );

    var previewXAxis = new Rickshaw.Graph.Axis.Time({
        graph: preview.previews[0],
        timeFixture: new Rickshaw.Fixtures.Time.Local(),
    });

    previewXAxis.render();

    $('input.datepicker')
        .datepicker()
        .change(function () {
            window.Metrics.redraw();
        });

    $('#metric_pattern').change(function () { window.Metrics.redraw(); });
}

Metrics.prototype.json = function() {
    this.requests.push(d3.json.apply(null, arguments));
}

Metrics.prototype.abort_all = function () {
    this.requests.forEach(function (request) {
        request.abort();
    });
    this.requests = [];
}



function draw() {
    window.Metrics = new Metrics();
    window.Metrics.draw();
}
$(document).ready(draw);
