function (key, value) {
    var len = value.count * 10 + value.prob_dist.length;
    for (var i = 0; i < value.prob_dist.length; i++) {
        value.prob_dist[i] = (value.prob_dist[i] + 1) / len;
    }
    delete value['discrete'];
    value.version = %(version)d;
    value.ad = 'mgof';
    value.start =  new Date('%(start)s');
    value.end = new Date('%(end)s');
    value.metric = '%(metric)s';
    return value;
}
