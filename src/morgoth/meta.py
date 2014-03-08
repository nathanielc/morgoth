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

from morgoth.data.mongo_clients import MongoClients
from morgoth.config import Config
from morgoth.schedule import Schedule
from morgoth.utils import timedelta_from_str
from morgoth.utc import now
import morgoth.detectors
import morgoth.notifiers

import gevent
import re

import logging
logger = logging.getLogger(__name__)

__all__ = ['Meta']

class Meta(object):
    _db = MongoClients.Normal.morgoth
    _db_admin = MongoClients.Normal.admin
    _db_name = Config.mongo.database_name
    _use_sharding = Config.mongo.use_sharding
    _needs_updating = {}
    _refresh_interval = Config.get(['metric_meta', 'refresh_interval'], 60)
    _finishing = False
    _null_manager = None
    """ Dict containg the meta data for each metric """
    _meta = {}
    """ Dict of patterns to the MetricManager"""
    _managers = {}


    @classmethod
    def load(cls):
        """
        Loads the existing meta data
        """

        # Load managers from conf
        for pattern, conf in Config.metrics.items():
            cls._managers[pattern] = MetricManager(pattern, conf)

        # Load metrics from database
        for meta in cls._db.meta.find():
            metric = meta['_id']
            cls._meta[metric] = meta
            manager = cls._match_metric(metric)
            manager.add_metric(metric)
            manager.start()


    @classmethod
    def update(cls, metric, value):
        """
        Update a metrics meta data

        A metric's meta data will be only eventually consistent
        """
        if cls._finishing:
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
            #logger.debug("Created new meta %s" % str(meta))
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
        if not cls._finishing:
            cls._finishing = True
            for metric in cls._needs_updating:
                cls._update(metric)
            cls._needs_updating = {}


    @classmethod
    def notify_anomalous(cls, window):
        """
        Notify that a given metric is anomalous for the given window

        @param window: a window object representing the anomalous time frame
        """
        # Add the anomaly to the db
        cls._db.anomalies.insert({
            'metric' : window.metric,
            'start' : window.start,
            'stop' : window.stop,
        })

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
                #logger.debug("Got existing meta %s" % str(existing_meta))
                meta['version'] = existing_meta['version']
                meta['min'] = min(existing_meta['min'], meta['min'])
                meta['max'] = max(existing_meta['max'], meta['max'])
                meta['count'] = max(existing_meta['count'], meta['count'])
            #logger.debug("Saving meta %s for metric %s"% (str(meta), metric))
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
        cls._db.meta.insert(meta)

        manager = cls._match_metric(metric)
        manager.add_metric(metric)

        manager.start()

    @classmethod
    def _match_metric(cls, metric):
        """
        Determine which pattern matches the given metric

        @param metric: the name of the metric
        @return the MetricManager for the given metric
        """
        for pattern, manager in cls._managers.items():
            if re.match(pattern, metric):
                return manager

        # No config for the metric
        logger.warn("Metric '%s' has no matching configuration" % metric)
        return cls._null_manager




class MetricManager(object):
    """
    Manages all the activity around metrics
    An instance of this class will be created for each entry in the 'metrics' section of the config
    """
    def __init__(self, pattern, conf):
        """

        @param pattern: the regex string that matches the metrics
        @param conf: the conf object from the 'metric' section
        """
        self._pattern = pattern
        self._detectors = []
        self._notifiers = []
        self._metrics = set()
        self._started = False

        # Load Detectors
        d_loader = morgoth.detectors.get_loader()
        detectors = conf.get('detectors', {})
        if not detectors:
            logger.warn('No Detectors defined for metric pattern "%s"' % pattern)
        else:
            self._detectors = d_loader.load(detectors)

        # Load Notifiers
        n_loader = morgoth.notifiers.get_loader()
        notifiers = conf.get('notifiers', {})
        if not notifiers:
            logger.warn('No Notifiers defined for metric pattern "%s"' % pattern)
        else:
            self._notifiers = d_loader.load(notifiers)

        # Load schedule
        self._duration = timedelta_from_str(conf.schedule.duration)
        self._period = timedelta_from_str(conf.schedule.period)
        self._delay = timedelta_from_str(conf.schedule.delay)
        self._aligned = conf.get(['schedule', 'aligned'], True)

        self._schedule = Schedule(self._period, self._check_metrics, self._delay)



    def add_metric(self, metric):
        """
        Add new metric to the manager
        """
        self._metrics.add(metric)

    def start(self):
        """
        Start watching the metrics for anomanlies
        """
        if not self._started:
            logger.debug(
                    "Starting MetricManager %s with detectors %s and notifiers %s",
                    self._pattern,
                    self._detectors,
                    self._notifiers,
                )
            self._started = True
            if self._aligned:
                self._schedule.start_aligned()
            else:
                self._schedule.start()

    def _check_metrics(self):
       """
       Handle the next check
       """
       end = now() - self._delay
       start = end - self._duration
       logger.debug("Checking metrics for next window %s:%s", start, end)
       for metric in self._metrics:
           gevent.spawn(self._check_window, metric, start, end)

    def _check_window(self, metric, start, end):
        """
        Check an indivdual window

        @param metric: the name of the metric
        @param start: the start time
        @param end: the end time
        """
        votes = 1
        for detector in self._detectors:
            window = detector.is_anomalous(metric, start, end)
            if window.anomalous:
                votes += 1
            else:
                votes -= 1

        if votes > 0:
            Meta.notify_anomalous(window)
            for notifier in self._notifiers:
                notifier.notify(window)



class NullMetricManager(MetricManager):
    def __init__(self):
        pass
    def start(self):
        pass
    def add_metric(self, metric):
        pass

Meta._null_manager = NullMetricManager()

