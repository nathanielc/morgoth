
from mongo_client_factory import MongoClientFactory
from Queue import Queue
from datetime import datetime

class Collector(object):
    __time_fmt = "%Y%m%d%H"
    def __init__(self):
        self._db = MongoClientFactory.create().morgoth


    def delete_metric(self, metric):
        self._db.metrics.remove({'metric' : metric})
        self._db.meta.remove({'_id' : metric})
        self._db.windows.remove({'value.metric' : metric})

    def insert(self, dt_utc, metric, value):
        value = float(value)
        id = self._get_id(dt_utc, metric)
        time = datetime(dt_utc.year, dt_utc.month, dt_utc.day, dt_utc.hour)
        h = dt_utc.hour
        m = dt_utc.minute
        s = dt_utc.second
        query = {'_id' : id}
        update = { '$set' : {
                'hour' : h,
                'time' : time,
                'metric' : metric,
                'data.%d.%d' % (m, s) : value,
                }
            }
        self._db.metrics.update(query, update, upsert=True)

        # Update meta information
        success = False
        while not success:
            meta = self._db.meta.find_one({'_id': metric})
            if meta is None:
                meta = {
                    '_id' : metric,
                    'n_bins' : 10,
                    'version': 0,
                    'max' : value,
                    'min' : value,
                    'count' : 0,
                }
                self._db.meta.insert(meta)
            if value > meta['max'] or value < meta['min']:
                set_cmd = {}
                if value > meta['max']:
                    set_cmd = { 'max' : value }
                else:
                    set_cmd = { 'min' : value }
                ret = self._db.meta.update(
                    {
                        '_id' : metric,
                        'version' : meta['version']
                    }, {
                        '$set' : set_cmd,
                        '$inc' : { 'version' : 1}
                    })
                success = ret['updatedExisting']
            else:
                success = True

        self._db.meta.update({'_id' : metric}, { '$inc' : { 'count' : 1 } })


    def _get_id(self, dt_utc, metric):
        time = dt_utc.strftime(self.__time_fmt)
        return "%s:%s" % (time, metric)
