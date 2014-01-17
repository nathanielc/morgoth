
def get_col_for_metric(db, metric):
    """
    Return the collection connection for a given metric


    NOTE: the current implementation uses just one collection for all metric data.
    I plan to change this use multiple collections later

    @param db: the db connection to use
    @param metric: the name of the metric
    """
    return db.metrics
