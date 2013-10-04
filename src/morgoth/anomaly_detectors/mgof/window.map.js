function () {
    var P = [];
    for ( var i = 0; i < %(n_bins)d; i++) {
        P[i] = 1;
    }
    var start_h = %(start_h)d;
    var start_m = %(start_m)d;
    var start_s = %(start_s)d;
    var end_h = %(end_h)d;
    var end_m = %(end_m)d;
    var end_s = %(end_s)d;
    if (this.hour != start_h) {
        start_m = 0;
        start_s = 0;
    }
    if (this.hour != end_h) {
        end_m = 59;
        end_s = 59;
    }
    for (var m = start_m; m <= end_m; m++) {
        for (var s = start_s; s <= end_s; s++) {
            if (!(m in this.data) || !(s in this.data[m])) continue;
            var value = this.data[m][s];
            var discrete = Math.floor(
                    (value - %(m_min)f) / %(step_size)f
                );
            P[discrete] += 10
            emit('%(id)s', {
                discrete : discrete,
                P : P,
                count: 1,
                version: %(version)d,
                metric : this.metric,
            });
            P[discrete] -= 10
         }
    }
}
