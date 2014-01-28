
from datetime import datetime, timedelta
from mongo_clients import MongoClients
from morgoth.data import get_col_for_metric
import re

class Reader(object):
    def __init__(self):
        self._db = MongoClients.Normal.morgoth

    def get_metrics(self, pattern=None):
        metrics = []
        for metric in self._db.meta.find():
            name = metric['_id']
            if pattern and not re.match(pattern, name):
                continue
            metrics.append(name)
        return metrics

    def get_data(self, metric, dt_start=None, dt_end=None):
        time_query = {}
        if dt_start:
            time_query['$gte'] = dt_start
        if dt_end:
            time_query['$lte'] = dt_end
        col = get_col_for_metric(self._db, metric)
        query = {'metric' : metric}
        if time_query:
            query['time'] = time_query
        data = col.find(query)
        time_data = []
        for point in data:
            time_data.append((point['time'].isoformat(), point['value']))
        return time_data

