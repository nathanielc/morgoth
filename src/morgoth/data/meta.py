
from morgoth.data.mongo_clients import MongoClients
from morgoth.config import Config

import gevent
import re

import logging
logger = logging.getLogger(__name__)

class Meta(object):
    _db = MongoClients.Normal.morgoth
    _db_admin = MongoClients.Normal.admin
    _db_name = Config.mongo.database_name
    _use_sharding = Config.mongo.use_sharding
    _meta = {}
    _needs_updating = {}
    _refresh_interval = Config.get(['metric_meta', 'refresh_interval'], 60)
    _finishing = False


    @classmethod
    def load(cls):
        """
        Loads the existing meta data
        """
        for meta in cls._db.meta.find():
            metric = meta['_id']
            cls._meta[metric] = meta
            cls._start_ads(metric)


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
                cls._init_new_metric(metric, meta)
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

    @classmethod
    def _init_new_metric(cls, metric, meta):
        """
        Initialize a new metric

        @param metric: the name of the metric
        @param meta: the metadata about the metric
        """
        # Insert the meta data for the metric
        cls._db.meta.insert(meta)

        # Start anomaly dectectors
        cls._start_ads(metric)


    @classmethod
    def _start_ads(cls, metric):
        """
        Start the appropriate Anomaly Detectors for the given metric

        @param metric: the name of the metric
        """
        for pattern, conf in Config.metrics.items():
            if re.match(pattern, metric):
                logger.debug("Metric %s matched %s" % (metric, pattern))
                break



