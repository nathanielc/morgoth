import logging
import re

from morgoth.data.reader import Reader
from morgoth.date_utils import from_epoch, to_epoch, total_seconds


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
            metric = row['name']
            if metric == 'morgoth_anomalies':
                continue
            if pattern is None:
                metrics.append(metric)
            elif re.search(pattern, metric):
                metrics.append(metric)

        return metrics

    def get_data(self, metric, start=None, stop=None, step=None):
        super(InfluxReader, self).get_data(metric, start, stop, step)
        selection = 'time, value'
        group_by = None
        if step:
            selection = 'time, mean(value)'
            group_by = ' group by time(%ds)' % total_seconds(step)
        query = "select %s from %s " % (selection, metric)
        where = []
        if start:
            where.append("time > %ds" % (to_epoch(start)))
        if stop:
            where.append("time < %ds" % (to_epoch(stop)))

        if where:
            query += 'where ' + ' and '.join(where)

        if group_by:
            query += group_by


        result = self._db.query(query, time_precision='s')

        time_data = []
        if result:
            for row in result[0]['points']:
                if group_by:
                    epoch, value = row
                else:
                    epoch, _, value = row
                time_data.insert(0, (from_epoch(epoch).isoformat(), value))

        return time_data



    def get_anomalies(self, metric, start=None, stop=None):
        super(InfluxReader, self).get_anomalies(metric, start, stop)

        anomalies = []

        query = "select time, start, stop from morgoth_anomalies where metric = '%s'" % metric

        time_clause = None
        if start and stop:
            time_clause = 'stop >= %d and start <= %d' % (to_epoch(start), to_epoch(stop))
        elif start:
            time_clause = 'start >= %d' % to_epoch(start)
        elif stop:
            time_clause = 'stop <= %d' % to_epoch(stop)

        if time_clause:
            query += 'and %s' % time_clause

        result = self._db.query(query, time_precision='s')
        if result:
            for time, num, start, stop in result[0]['points']:
                anomalies.append({
                    'start' : from_epoch(start).isoformat(),
                    'stop' : from_epoch(stop).isoformat(),
                    'id' : str(time) + str(num)
                })
        return anomalies




    def get_histogram(self, metric, n_bins, start, stop):
        super(InfluxReader, self).get_histogram(metric, n_bins, start, stop)

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
        result = self._db.query(query, time_precision='s')
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


    def get_percentile(self, metric, percentile, start=None, stop=None):
        super(InfluxReader, self).get_percentile(metric, percentile, start, stop)
        query = "select percentile(value, %f) from %s " % (percentile, metric)
        where = []
        if start:
            where.append("time > %ds" % (to_epoch(start)))
        if stop:
            where.append("time < %ds" % (to_epoch(stop)))

        if where:
            query += 'where ' + ' and '.join(where)


        result = self._db.query(query, time_precision='s')

        percentile = None
        if result:
            percentile = result[0]['points'][0][1]

        return percentile

