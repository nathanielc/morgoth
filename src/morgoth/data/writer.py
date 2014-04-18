#
# Copyright 2014 Nathaniel Cook
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


from gevent.queue import JoinableQueue
from gevent.event import Event
from morgoth.data.mongo_clients import MongoClients
from morgoth.data import get_col_for_metric
from morgoth.meta import Meta
import gevent
import pymongo

import logging
logger = logging.getLogger(__name__)

class Writer(object):
    """
    Class that provides write access to store metrics in morgoth
    """
    __time_fmt = "%Y%m%d%H"
    _max_size = 1000
    _worker_count = 2
    _flush_token = 'FLUSH'
    def __init__(self, max_size=None, worker_count=None):
        self._db = MongoClients.Normal.morgoth
        if max_size is None:
            max_size = Writer._max_size
        self._queue = JoinableQueue(maxsize=max_size)
        if worker_count is None:
            worker_count = Writer._worker_count
        self._worker_count = worker_count
        self._running = Event()
        self._flushed = Event()
        self._flushing = False
        self._closing = False
        for i in xrange(self._worker_count):
            gevent.spawn(self._worker)

    @classmethod
    def configure_defaults(cls, config):
        """
        Configure the default writer options

        @param config: the app configuration object
        @type config: morgoth.config.Config
        """
        cls._max_size = config.get(['writer', 'queue', 'max_size'], cls._max_size)
        cls._worker_count = config.get(['writer', 'queue', 'worker_count'], cls._worker_count)

    def _worker(self):
        while True:
            self._running.wait()
            while not self._queue.empty():
                dt_utc, metric, value = self._queue.get()
                if dt_utc == self._flush_token:
                    Meta.flush()
                    self._flushed.set()
                    continue
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
        """
        Insert a data point for a given metric

        @param dt_utc: the utc datetime of the data point
        @param metric: the name of the metric
        @param value: the value of the metric, must not be None
        @raise ValueError of `value` is None
        """
        if self._closing:
            logger.debug("Writer is closed")
            return
        if value is None:
            raise ValueError('value cannot be None')
        self._queue.put((dt_utc, metric, value))
        self._running.set()

    def delete_metric(self, metric):
        Meta.delete_metric(metric)

    def close(self):
        self._closing = True
        self._queue.join()
        logger.debug("Writer queue is empty")
        Meta.finish()

    def flush(self):
        if not self._flushing:
            self._flushing = True
            self._queue.put((self._flush_token, None, None))

        self._flushed.wait()
        self._flushing = False
        self._flushed.clear()


