
function Metrics() {
    this.series = [];
    this.graph = undefined;
}

Metrics.prototype.draw = function () {

    var that = this;
    var series = this.series;

    var num_points = 1000;
    var start = new Date('2014-03-29T00:00');
    var stop = new Date('2014-03-29T01:00');
    var step = ((stop.getTime() - start.getTime()) / 1000 ) / num_points;

    var metric_pattern = 'qualtrics.cp.*';

    var palette = new Rickshaw.Color.Palette();

    d3.json("http://localhost:8001/metrics?pattern=" + metric_pattern, function(data) {
        metrics = data.metrics
        metrics.forEach(function(metric, index) {
            var url = "http://localhost:8001/data/"
                    + metric
                    + '?start=' + start.toString()
                    + '&stop=' + stop.toString()
                    + '&step=' + step + 's';
            d3.json(url, function(rows) {
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

Metrics.prototype.update = function() {
    if (this.graph == undefined) {
        this._init_graph();
    }
    var graph = this.graph;
    graph.update();
    $('#legend').empty();
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
    //Rickshaw.Series.fill(series, null);
    var graph = new Rickshaw.Graph( {
                element: document.querySelector("#chart"),
                width: 1000,
                height: 50,
                renderer: 'horizon',
                series: series,
                interpolation: 'step-after',
    } );
    this.graph = graph;

    var x_axis = new Rickshaw.Graph.Axis.Time( { graph: graph } );

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
}


function draw() {
    var metrics = new Metrics();
    metrics.draw();
}
$(document).ready(draw);
