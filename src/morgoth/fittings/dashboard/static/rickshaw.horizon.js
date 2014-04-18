Rickshaw.namespace('Rickshaw.Graph.Renderer.Horizon');

Rickshaw.Graph.Renderer.Horizon = Rickshaw.Class.create( Rickshaw.Graph.Renderer, {

    name: 'horizon',

    defaults: function($super) {

        return Rickshaw.extend( $super(), {
            unstack: true,
            fill: false,
            stroke: true,
            padding:{ top: 0.01, right: 0.01, bottom: 0.01, left: 0.01 },
            dotSize: 3,
            strokeWidth: 2
        } );
    },

    initialize: function($super, args) {
        $super(args);
    },

    seriesPathFactory: function() {

        var graph = this.graph;

        var factory = d3.svg.area()
            .x( function(d) { return graph.x(d.x) } )
            .y0( function(d) { return graph.y(d.y0) } )
            .y1( function(d) { return graph.y(d.y + d.y0) } )
            .interpolate(graph.interpolation).tension(this.tension);

        factory.defined && factory.defined( function(d) { return d.y !== null } );
        return factory;
    },

    _renderLines: function() {

        var graph = this.graph;

        var nodes = graph.vis.selectAll("path")
            .data(this.graph.stackedData)
            .enter().append("svg:path")
            .attr("d", this.seriesPathFactory());

        var i = 0;
        graph.series.forEach(function(series) {
            if (series.disabled) return;
            series.path = nodes[0][i++];
            this._styleSeries(series);
        }, this);
    },

    domain: function (data) {

        var stackedData = data || this.graph.stackedData || this.graph.stackData();

        var xMin = +Infinity;
        var xMax = -Infinity;

        var yMin = +Infinity;
        var yMax = -Infinity;

        stackedData.forEach(function (series) {
            var min = d3.min(series, function(d) { return d.y;});
            var max = d3.max(series, function(d) { return d.y;});

            var step = (max - min) / depth;
            if (step > yMax) yMax = step;

            if (!series.length) return;

            if (series[0].x < xMin) xMin = series[0].x;
            if (series[series.length - 1].x > xMax) xMax = series[series.length - 1].x;
        });
        return { x: [xMin, xMax], y: [yMin, yMax]};
    },

    _renderHorizon: function() {
        var graph = this.graph;
        var stackedData = graph.stackedData || graph.stackData();
        console.log(stackedData);

        var depth = 2;
        var horizon_data = [];
        var domain = this.domain();
        stackedData.forEach(function (data) {

            var bottom = 0;
            var top = step;
            for ( var i = 0; i < depth; i++) {
                hdata = [];
                data.forEach(function (d) {
                    var y = Math.min(d.y - bottom, top);
                    hdata.push({x: d.x, y: y, y0: d.y0});
                });
                top += step;
                bottom += step;
                horizon_data.push(hdata);
            }
        });

        console.log(horizon_data);
        var nodes = graph.vis.selectAll("path")
            .data(horizon_data)
            .enter().append("svg:g")

        nodes.append('svg:path')
            .attr("d", this.seriesPathFactory())
            .attr("class", "area")

        console.log(nodes);
        var i = 0;
        graph.series.forEach(function(series) {
            if (series.disabled) return;
            series.path = nodes[0].slice(i, i + depth);
            this._styleSeries(series);
            i += depth;
        }, this);
    },

    _styleSeries: function(series) {

        if (!series.path) return;

        series.path.forEach(function (path, i) {
            d3.select(path).select('.area')
                .attr('fill', d3.interpolateRgb(series.color, 'black')(0.25 * i));
        });


        if (series.className) {
            series.path.forEach(function (path, i) {
                path.setAttribute('class', series.className);
            });
        }
    },

    render: function() {

        var graph = this.graph;

        graph.vis.selectAll('*').remove();

        this._renderHorizon();
    }
} );
