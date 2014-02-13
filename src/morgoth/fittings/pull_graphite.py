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



from datetime import timedelta, datetime
from morgoth.config import Config
from morgoth.utc import now, utc
from morgoth.data.writer import Writer
from morgoth.fittings.fitting import Fitting
from morgoth.schedule import Schedule
from morgoth.utils import timedelta_from_str
import urllib2
import json

import logging
logger = logging.getLogger(__name__)

class PullGraphite(Fitting):
    _date_format = "%H:%M_%Y%m%d"
    def __init__(self, metric_pattern, graphite_url, period, user=None, password=None):
        """
        Initialize a pull graphite fitting

        @param metric_pattern: a graphite pattern for the metrics to pull
        @param graphite_url: the url of the graphite instance http://host[:port]
        @param period: timedelta object for how often the data should be pulled
        @param user: username for basic http authentication
        @param password: password for basic http authentication
        """
        super(PullGraphite, self).__init__()
        self._metric_pattern = metric_pattern
        self._graphite_url = graphite_url
        self._period = period
        self._user = user
        self._password = password
        self._schedule = Schedule(self._period, self._pull)
        self._writer = Writer()

    @classmethod
    def from_conf(cls, conf):
        metric_pattern = conf.get('metric_pattern', '*')
        graphite_url = conf.get('graphite_url', 'http://localhost')
        period = timedelta_from_str(conf.get('period', '5m'))
        user = conf.get('user', None)
        password = conf.get('password', None)
        return PullGraphite(
                metric_pattern,
                graphite_url,
                period,
                user,
                password
            )

    def start(self):
        self._schedule.start_aligned()

    def stop(self):
        self._schedule.stop()

    def _pull(self):
        """
        Pull the next chunk of data from the graphite instance
        """
        stop = now() - timedelta(hours=7)
        start = stop - self._period
        self._pull_metrics(self._metric_pattern, start, stop)

    def _pull_metrics(self, metric, start, stop):
        logger.debug('Pulling metrics for %s', metric)
        request = urllib2.Request("%s/metrics/find?query=%s" % (self._graphite_url, metric))
        self._add_auth(request)
        children = json.load(urllib2.urlopen(request))
        for child in children:
            child_metric = child['id']
            if child['leaf'] == 0:
                self._pull_metrics("%s.*" % child_metric, start, stop)
            else:
                self._pull_data(child_metric, start, stop)
        if not children:
            self._pull_data(metric, start, stop)

    def _pull_data(self, metric, start, stop):
        url = "%s/render?target=%s&format=json&from=%s&until=%s" % (
                self._graphite_url,
                metric,
                start.strftime(self._date_format),
                stop.strftime(self._date_format),
            )
        logger.debug("Pulling data for %s:\n%s", metric, url)
        request = urllib2.Request(url)
        self._add_auth(request)
        data = json.load(urllib2.urlopen(request))
        if not data:
            logger.warn("Could not find metric %s", metric)
        if not data:
            return
        for value, timestamp in data[0]['datapoints']:
            dt_utc = datetime.fromtimestamp(timestamp + 60*60, utc)
            self._writer.insert(dt_utc, metric, value)






    def _add_auth(self, request):
        if self._user and self._password:
            request.add_header("Authorization", "Basic " + (('%s:%s' % (self._user, self._password)).encode('base64')))

