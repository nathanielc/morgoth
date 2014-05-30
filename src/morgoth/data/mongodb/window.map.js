function () {
    var prob_dist = [];
    for ( var i = 0; i < %(n_bins)d; i++) {
        prob_dist[i] = 0;
    }
    var discrete = Math.floor(
            (this.value - %(m_min)f) / %(step_size)f
        );
    prob_dist[discrete] += 10
    emit('histogram', {
        discrete : discrete,
        prob_dist : prob_dist,
        count: 1,
    });
}
