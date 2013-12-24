function (key, value) {
    var len = value.count * 10 + value.P.length;
    for (var i = 0; i < value.P.length; i++) {
        value.P[i] = (value.P[i] + 1) / len;
    }
    delete value['discrete'];
    value.version = %(version)d;
    value.ad = 'mgof';
    value.start =  new Date('%(start)s');
    value.end = new Date('%(end)s');
    value.metric = '%(metric)s';
    return value;
}
