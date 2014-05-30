
from morgoth.data.reader import Reader
import calendar

import logging
logger = logging.getLogger(__name__)

DATE_FORMAT = '%Y-%m-%d %H:%M:%S'

class InfluxReader(Reader):

    def __init__(self, db):
        super(InfluxReader, self).__init__()
        self._db = db

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
        for _, total, bucket_start, count in result[0]['points']:
            i = int((bucket_start - m_min) / step_size)
            hist[i] = count * 10 / float(total * 10 + n_bins)

        return hist, total



