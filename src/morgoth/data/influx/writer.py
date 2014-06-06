
from morgoth.data.writer import DefaultWriter
from morgoth.utc import to_epoch


import logging
logger = logging.getLogger(__name__)

class InfluxWriter(DefaultWriter):
    def __init__(self, db, app, max_size=None, worker_count=None):
        super(InfluxWriter, self).__init__(app, max_size, worker_count)
        self._db = db


    def _insert(self, dt_utc, metric, value):
        data = [{
            'name' : metric,
            'columns' : ['time', 'value'],
            'points': [
                [to_epoch(dt_utc) * 1000, value]
            ],
        }]
        self._db.write_points_with_precision(data, time_precision='m')
