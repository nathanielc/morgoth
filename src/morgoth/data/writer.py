

from gevent.queue import JoinableQueue
from gevent.event import Event
from morgoth.data.mongo_clients import MongoClients
from morgoth.config import Config
from morgoth.data import get_col_for_metric
from morgoth.meta import Meta
import gevent
import pymongo

import logging
logger = logging.getLogger(__name__)

class Writer(object):
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
                col = get_col_for_metric(self._db, metric)
                col.insert({
                    'time' : dt_utc,
                    'value' : value,
                    'metric' : metric}
                )
                Meta.update(metric, value)
                self._queue.task_done()
            self._running.clear()

    def insert(self, dt_utc, metric, value):
        if self._closing:
            logger.debug("Writer is closed")
            return
        self._queue.put((dt_utc, metric, value))
        self._running.set()

    def close(self):
        self._closing = True
        self._queue.join()
        logger.debug("Writer queue is empty")
        Meta.finish()


