
from config import Config
from mongo_clients import MongoClients
from metric_meta import MetricMeta
from datetime import datetime
from gevent.queue import JoinableQueue
from gevent.event import Event
import gevent

import logging
logger = logging.getLogger(__name__)

class Collector(object):
    __time_fmt = "%Y%m%d%H"
    def __init__(self):
        # Write optimized MongoClient
        self._db = MongoClients.Normal.morgoth
        self._queue = JoinableQueue(maxsize=Config.get(['write_queue', 'max_size'], 1000))
        self._worker_count = Config.get(['write_queue', 'worker_count'], 2)
        self._running = Event()
        self._closing = False
        for i in xrange(self._worker_count):
            gevent.spawn(self._worker)


    def _worker(self):
        while True:
            self._running.wait()
            while not self._queue.empty():
                dt_utc, metric, value = self._queue.get()
                self._db[metric].insert({'_id' : dt_utc, 'value' : value})
                MetricMeta.update(metric, value)
                self._queue.task_done()
            self._running.clear()

    def insert(self, dt_utc, metric, value):
        if self._closing:
            logger.debug("Collector is closed")
            return
        self._queue.put((dt_utc, metric, value))
        self._running.set()

    def close(self):
        self._closing = True
        self._queue.join()
        logger.debug("Collector queue is empty")
        MetricMeta.finish()
