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

from bson.code import Code
from bson.son import SON
from datetime import datetime, timedelta
from morgoth.data.reader import Reader
from morgoth.error import MorgothError
from morgoth.utc import utc
import re
import os

import logging
logger = logging.getLogger(__name__)

__dir__ = os.path.dirname(__file__)

class MongoReader(Reader):
    """
    Class that provides read access to the metric data and anomalies
    """
    map = None
    reduce_code = None
    finalize = None
    def __init__(self, db):
        """
        Create MongoReader

        @param db: database connection object pymongo.MongoClient.<database>
        """
        super(MongoReader, self).__init__()
        self._db = db

    def get_metrics(self, pattern=None):
        metrics = []
        for metric in self._db.meta.find():
            name = metric['_id']
            if pattern and not re.match(pattern, name):
                continue
            metrics.append(name)
        return metrics

    def get_data(self, metric, start=None, stop=None, step=None):
        super(MongoReader, self).get_data(metric, start, stop, step)
        time_query = {}
        if start:
            time_query['$gte'] = start
        if stop:
            time_query['$lte'] = stop
        query = {'metric' : metric}
        if time_query:
            query['time'] = time_query
        data = self._db.metrics.find(query)
        time_data = []

        count = 0
        total = 0.0
        boundary = None
        if start and step:
            boundary = (start + step).replace(tzinfo=utc)
        raw_data = []
        for point in data:
            raw_data.append((point['time'].isoformat(), point['value']))
            if boundary and step:
                if point['time'] > boundary:
                    if count > 0:
                        time_data.append((boundary.isoformat(), total / count))
                    boundary += step
                    while boundary + step <= point['time']:
                        boundary += step
                    count = 0
                    total = 0.0
                count += 1
                total += point['value']

            else:
                time_data.append((point['time'].isoformat(), point['value']))
        if count > 0:
            time_data.append((boundary.isoformat(), total / count))

        return time_data

    def get_anomalies(self, metric, start=None, stop=None):
        super(MongoReader, self).get_anomalies(metric, start, stop)
        query = {'metric' : metric}
        if start and stop:
            query['stop'] = {'$gte' : start}
            query['start'] = {'$lte' : stop}

        data = self._db.anomalies.find(query)
        anomalies = []
        for point in data:
            anomalies.append({
                'id' : str(point['_id']),
                'start' : point['start'].isoformat(),
                'stop' : point['stop'].isoformat()
            })
        return anomalies

    def get_histogram(self, metric, n_bins, start, stop):
        super(MongoReader, self).get_histogram(metric, n_bins, start, stop)

        meta = self._db.meta.find_one({'_id' : metric})
        if meta is None:
            raise MorgothError("No such metric '%s'" % metric)
        m_max = meta['max']
        m_min = meta['min']
        version = meta['version']


        step_size = ((m_max * 1.01) - m_min) / float(n_bins)

        map_values = {
            'step_size' : step_size,
            'm_min' : m_min,
            'n_bins' : n_bins,
        }

        finalize_values = {
            'start' : start.isoformat(),
            'stop' : stop.isoformat(),
            'version': version,
            'metric' : metric,
        }

        map_code = Code(self.map % map_values)
        finalize_code = Code(self.finalize % finalize_values)


        query = {
            'metric' : metric,
            'time' : {'$gte' : start, '$lt' : stop},
        }
        result = self._db.metrics.inline_map_reduce(map_code, self.reduce_code,
            query=query,
            finalize=finalize_code
        )
        if result:
            return result[0]['value']['prob_dist'], result[0]['value']['count']
        return [0] * n_bins, 0

# Initialize js code
if MongoReader.map is None:
    with open(os.path.join(__dir__, 'window.map.js')) as f:
        MongoReader.map = f.read()
    with open(os.path.join(__dir__, 'window.reduce.js')) as f:
        MongoReader.reduce_code = Code(f.read())
    with open(os.path.join(__dir__, 'window.finalize.js')) as f:
        MongoReader.finalize = f.read()
