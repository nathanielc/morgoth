
from mongo_client_factory import MongoClientFactory
from datetime import datetime, timedelta

class Reader(object):
    def __init__(self):
        self._db = MongoClientFactory.create().morgoth

    def get_data(self, metric, dt_start=None, dt_end=None):
        time_query = {}
        if dt_start:
            time_query['$gte'] = dt_start
        if dt_end:
            time_query['$lte'] = dt_end
        data = self._db.metrics.find({
                'time' : time_query,
                'metric' : metric
            })
        time_data = []
        for doc in data:
            date = doc['time']
            for m in range(60):
                m = str(m)
                if m not in doc['data']: continue
                minute = doc['data'][m]
                for s in range(60):
                    s = str(s)
                    if s not in minute: continue
                    delta = timedelta(minutes=int(m), seconds=int(s))
                    time_data.append(((date + delta).isoformat(), float(minute[s])))
        return time_data

