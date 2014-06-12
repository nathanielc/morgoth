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
import gevent
import pymongo
from morgoth.utc import utc

import logging
logger = logging.getLogger(__name__)

class Writer(object):
    """
    Class that provides write access to store metrics in morgoth
    """
    def __init__(self, app):
        self._app = app

    @property
    def _metrics_manager(self):
        """
        The metrics manager instance
        """
        return self._app.metrics_manager

    def insert(self, dt_utc, metric, value):
        """
        Insert a data point for a given metric

        @param dt_utc: the utc datetime of the data point
        @param metric: the name of the metric
        @param value: the value of the metric, must not be None
        @raise ValueError of `value` is None
        """
        assert dt_utc.tzinfo == utc
        self._metrics_manager.new_metric(metric)

    def record_anomalous(self, metric, start, stop):
        """
        Record that a given metric is anomalous for the given window

        @param metric: the name of the metric
        @param start: the start time of the anomalous window
        @param stop: the stop time of the anomalous window
        """
        assert start.tzinfo == utc
        assert stop.tzinfo == utc



    def delete_metric(self, metric):
        """
        Delete all data for the metric
        """
        pass

    def close(self):
        """
        Close writer
        """
        pass

    def flush(self):
        """
        Flush all pending writes and block until complete
        """
        pass


class DefaultWriter(Writer):
    """
    Default implementation of the Writer class
    """
    __time_fmt = "%Y%m%d%H"
    _max_size = 1000
    _worker_count = 2
    _flush_token = 'FLUSH'
    def __init__(self, app, max_size=None, worker_count=None):
        super(DefaultWriter, self).__init__(app)
        if max_size is None:
            max_size = DefaultWriter._max_size
        self._queue = JoinableQueue(maxsize=max_size)
        if worker_count is None:
            worker_count = DefaultWriter._worker_count
        self._worker_count = worker_count
        self._running = Event()
        self._flushed = Event()
        self._flushing = False
        self._closing = False
        for i in xrange(self._worker_count):
            gevent.spawn(self._worker)

    @classmethod
    def get_options(cls, config):
        """
        Read writer conf and return dict of options for creating DefaultWriter
        The options dict can later be passed as key word args to the DefaultWriter constructor
        """
        options = {}
        try:
            options['max_size'] = config.writer.get(['writer', 'queue', 'max_size'], cls._max_size)
            options['worker_count'] = config.writer.get(['writer', 'queue', 'worker_count'], cls._worker_count)
        except KeyError:
            pass
        return options


    def _worker(self):
        """
        Worker method that pulls pending writes off a queue
        and performs them in batches
        """
        while True:
            self._running.wait()
            while not self._queue.empty():
                dt_utc, metric, value = self._queue.get()
                if dt_utc == self._flush_token:
                    self._flush()
                    self._flushed.set()
                    continue
                self._insert(dt_utc, metric, value)
                self._queue.task_done()
            self._running.clear()

    def insert(self, dt_utc, metric, value):
        if self._closing:
            logger.debug("Writer is closed")
            return
        if value is None:
            raise ValueError('value cannot be None')
        self._queue.put((dt_utc, metric, value))
        self._running.set()
        super(DefaultWriter, self).insert(dt_utc, metric, value)

    def _insert(self, dt_utc, metric, values):
        """
        Perform actual insert into db backend
        """
        pass


    def close(self):
        self._closing = True
        self._queue.join()
        logger.debug("Writer queue is empty")

    def flush(self):
        if not self._flushing:
            self._flushing = True
            self._queue.put((self._flush_token, None, None))

        self._flushed.wait()
        self._flushing = False
        self._flushed.clear()

    def _flush(self):
        """
        Perform the actual flush of all data in the queue
        """
        pass


