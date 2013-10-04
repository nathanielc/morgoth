function (key, values) {
    var a = values[0];
    for ( var i = 1; i < values.length; i++) {
        var b = values[i];
        a.P[b.discrete] += 10;
        a.count += b.count;
    }
    return a;
}
