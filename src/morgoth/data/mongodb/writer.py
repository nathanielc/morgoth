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


from morgoth.data.writer import DefaultWriter
import gevent
import pymongo

import logging
logger = logging.getLogger(__name__)

class MongoWriter(DefaultWriter):
    """
    MongoDB implementation of the Writer class
    """
    _needs_updating = {}
    _refresh_interval = None
    _finishing = False
    """ Dict containg the meta data for each metric """
    _meta = {}

    def __init__(self, app, db, refresh_interval, max_size=None, worker_count=None):
        super(MongoWriter, self).__init__(app, max_size, worker_count)
        self._db = db
        self._refresh_interval = refresh_interval
        self._load()

    def _load(self):
        """
        Loads the existing meta data

        @param config: the app configuration object
        @type config: morgoth.config.Config
        """
        self._needs_updating = {}

        # Load metrics from database
        for meta in self._db.meta.find():
            metric = meta['_id']
            self._meta[metric] = meta
            #manager = self._match_metric(metric)
            #manager.add_metric(metric)
            #manager.start()


    def _insert(self, dt_utc, metric, value):
        """
        Perform actual insert into db backend
        """
        if self._finishing:
            return

        # Insert into metrics collection
        self._db.metrics.insert({
            'time' : dt_utc,
            'value' : value,
            'metric' : metric
        })


        if metric not in self._meta:
            meta = {
                '_id' : metric,
                'version': 0,
                'min' : value,
                'max' : value,
                'count' : 1,
            }
            #logger.debug("Created new meta %s" % str(meta))
            self._meta[metric] = meta
            self._update(metric)
        else:
            meta = self._meta[metric]
            #logger.debug("Updating meta with new value: %s %f" % (str(meta), value))
            meta['min'] = min(meta['min'], value)
            meta['max'] = max(meta['max'], value)
            meta['count'] += 1

            if metric not in self._needs_updating:
                self._needs_updating[metric] = True
                gevent.spawn(self._update_eventually, metric)
            #else:
            #    logger.debug("Metric already scheduled for update...")

    def _update(self, metric):
        """
        Update a metrics meta data

        A metric's meta data will be only eventually consistent
        """
        if metric in self._needs_updating:
            del self._needs_updating[metric]
        meta = self._meta[metric]
        # Update meta information
        success = False
        while not success:
            existing_meta = self._db.meta.find_one({'_id': metric})
            if existing_meta is None:
                self._init_new_metric(metric, meta)
            else:
                # Populate in memory meta with existing meta
                #logger.debug("Got existing meta %s" % str(existing_meta))
                meta['version'] = existing_meta['version']
                meta['min'] = min(existing_meta['min'], meta['min'])
                meta['max'] = max(existing_meta['max'], meta['max'])
                meta['count'] = max(existing_meta['count'], meta['count'])
            #logger.debug("Saving meta %s for metric %s"% (str(meta), metric))
            ret = self._db.meta.update(
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


    def _update_eventually(self, metric):
        gevent.sleep(self._refresh_interval)

    def _init_new_metric(self, metric, meta):
        """
        Initialize a new metric

        @param metric: the name of the metric
        @param meta: the metadata about the metric
        """
        self._db.meta.insert(meta)

        #manager = self._match_metric(metric)
        #manager.add_metric(metric)

        #manager.start()

    def _flush(self):
        """
        Perform the actual flush of all data in the queue
        """
        metrics = self._needs_updating.keys()
        for metric in metrics:
            self._update(metric)


    def record_anomalous(self, metric, start, stop):
       # Add the anomaly to the db
        self._db.anomalies.insert({
            'metric' : metric,
            'start' : start,
            'stop' : stop,
        })

    def delete_metric(self, metric):
        #logger.debug('Deleting metric %s', metric)
        self._db.metrics.remove({'metric' : metric})
        self._db.windows.remove({'value.metric' : metric})
        self._db.meta.remove({'_id' : metric})

