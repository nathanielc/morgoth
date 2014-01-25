function (key, values) {
    var a = values[0];
    for ( var i = 1; i < values.length; i++) {
        var b = values[i];
        a.prob_dist[b.discrete] += b.prob_dist[b.discrete];
        a.count += b.count;
    }
    return a;
}
