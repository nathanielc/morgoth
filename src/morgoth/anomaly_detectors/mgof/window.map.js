function () {
    var P = [];
    for ( var i = 0; i < %(n_bins)d; i++) {
        P[i] = 0;
    }
    var discrete = Math.floor(
            (this.value - %(m_min)f) / %(step_size)f
        );
    P[discrete] += 10
    emit('%(id)s', {
        discrete : discrete,
        P : P,
        count: 1,
    });
}
