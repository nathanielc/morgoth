

from datetime import datetime
from gevent.queue import JoinableQueue
from gevent.event import Event
from mongo_clients import MongoClients
from morgoth.config import Config
from morgoth.data import get_col_for_metric
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
                MetricMeta.update(metric, value)
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
        MetricMeta.finish()


class MetricMeta(object):
    _db = MongoClients.Normal.morgoth
    _db_admin = MongoClients.Normal.admin
    _db_name = Config.mongo.database_name
    _use_sharding = Config.mongo.use_sharding
    _meta = {}
    _needs_updating = {}
    _refresh_interval = Config.get(['metric_meta', 'refresh_interval'], 60)
    _finishing = False

    @classmethod
    def update(cls, metric, value):
        """
        Update a metrics meta data

        A metric's meta data will be only eventually consistent
        """
        if cls._finishing:
            #logger.info("MetricMeta is finishing and not accepting any more updates")
            return
        if metric not in cls._meta:
            meta = {
                '_id' : metric,
                'version': 0,
                'min' : value,
                'max' : value,
                'count' : 1,
            }
            conf_meta = cls._get_meta_from_config(metric)
            meta.update(conf_meta)
            logger.debug("Created new meta %s" % str(meta))
            cls._meta[metric] = meta
        else:
            meta = cls._meta[metric]
            #logger.debug("Updating meta with new value: %s %f" % (str(meta), value))
            meta['min'] = min(meta['min'], value)
            meta['max'] = max(meta['max'], value)
            meta['count'] += 1

        if metric not in cls._needs_updating:
            cls._needs_updating[metric] = True
            gevent.spawn(cls._update_eventually, metric)
        #else:
        #    logger.debug("Metric already scheduled for update...")


    @classmethod
    def finish(cls):
        logger.debug("Finishing MetricMeta")
        if not cls._finishing:
            cls._finishing = True
            for metric in cls._needs_updating:
                cls._update(metric)
            cls._needs_updating = {}


    @classmethod
    def _update_eventually(cls, metric):
        gevent.sleep(cls._refresh_interval)
        if cls._finishing: return
        del cls._needs_updating[metric]
        cls._update(metric)

    @classmethod
    def _get_meta_from_config(cls, metric):
        return {}

    @classmethod
    def _update(cls, metric):
        """
        Perform the actual update of the meta data
        """
        meta = cls._meta[metric]
        # Update meta information
        success = False
        while not success:
            existing_meta = cls._db.meta.find_one({'_id': metric})
            if existing_meta is None:
                cls._db.meta.insert(meta)
            else:
                # Populate in memory meta with existing meta
                logger.debug("Got existing meta %s" % str(existing_meta))
                meta['version'] = existing_meta['version']
                meta['min'] = min(existing_meta['min'], meta['min'])
                meta['max'] = max(existing_meta['max'], meta['max'])
                meta['count'] = max(existing_meta['count'], meta['count'])
            logger.debug("Saving meta %s for metric %s"% (str(meta), metric))
            ret = cls._db.meta.update(
                {
                    '_id' : metric,
                    'version' : meta['version']
                }, {
                    '$set' : {
                        'min' : meta['min'],
                        'max' : meta['max'],
                        'count' : meta['count'],
                    },
                    '$inc' : {'version' : 1}
                })
            success = ret['updatedExisting']
            if not success:
                logger.error("The meta version changed. This should not have happend as this class controls all updates.")



