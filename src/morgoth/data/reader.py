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

from datetime import datetime, timedelta
from mongo_clients import MongoClients
from morgoth.data import get_col_for_metric
from morgoth.utc import utc
import re

import logging
logger = logging.getLogger(__name__)

class Reader(object):
    """
    Class that provides read access to the metric data and anomalies
    """
    def __init__(self):
        self._db = MongoClients.Normal.morgoth

    def get_metrics(self, pattern=None):
        metrics = []
        for metric in self._db.meta.find():
            name = metric['_id']
            if pattern and not re.match(pattern, name):
                continue
            metrics.append(name)
        return metrics

    def get_data(self, metric, start=None, stop=None, step=None):
        time_query = {}
        if start:
            time_query['$gte'] = start
        if stop:
            time_query['$lte'] = stop
        col = get_col_for_metric(self._db, metric)
        query = {'metric' : metric}
        if time_query:
            query['time'] = time_query
        data = col.find(query)
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
        query = {'metric' : metric}
        if start and stop:
            query['stop'] = { '$gte' : start }
            query['start'] = { '$lte' : stop }

        data = self._db.anomalies.find(query)
        anomalies = []
        for point in data:
            anomalies.append({
                'id' : str(point['_id']),
                'start' : point['start'].isoformat(),
                'stop' : point['stop'].isoformat()
            })
        return anomalies
