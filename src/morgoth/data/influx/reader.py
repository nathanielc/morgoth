import logging

from morgoth.data.reader import Reader
from morgoth.utc import from_epoch


logger = logging.getLogger(__name__)

DATE_FORMAT = '%Y-%m-%d %H:%M:%S'

class InfluxReader(Reader):

    def __init__(self, db):
        super(InfluxReader, self).__init__()
        self._db = db

    def get_metrics(self, pattern=None):
        metrics = []
        result = self._db.query("list series")
        for row in result:
            metrics.append(row['name'])

        return metrics

    def get_data(self, metric, start=None, stop=None, step=None):
        query = "select time, value from %s " % metric
        where = []
        if start:
            where.append("time > '%s'" % start.strftime(DATE_FORMAT))
        if stop:
            where.append("time < '%s'" % stop.strftime(DATE_FORMAT))

        if where:
            query += 'where ' + ' and '.join(where)


        result = self._db.query(query, time_precision='m')

        time_data = []
        for epoch, _, value in result[0]['points']:
            time_data.insert(0, (from_epoch(epoch / 1000.0).isoformat(), value))

        return time_data




    def get_histogram(self, metric, n_bins, start, stop):

        result = self._db.query("select min(value), max(value) from %s" % metric)
        m_min = result[0]['points'][0][1]
        m_max = result[0]['points'][0][2]

        step_size = ((m_max * 1.01) - m_min) / float(n_bins)

        query = "select count(value), histogram(value, %f) from %s where time > '%s' and time < '%s'" % (
            step_size,
            metric,
            start.strftime(DATE_FORMAT),
            stop.strftime(DATE_FORMAT),
        )
        result = self._db.query(query)
        if not result:
            return [0] * n_bins, 0

        total = result[0]['points'][0][1]
        empty_value = 1.0 / float(total * 10 + n_bins)
        hist = [empty_value] * n_bins
        s = 0
        for _, total, bucket_start, count in result[0]['points']:
            i = int(round((bucket_start - m_min) / step_size))
            s += count
            hist[i] = count * 10 / float(total * 10 + n_bins)

        return hist, total



