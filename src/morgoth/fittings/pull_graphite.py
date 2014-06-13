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


from dateutil.tz import tzoffset, gettz
from datetime import timedelta, datetime
from morgoth.date_utils import now, utc
from morgoth.fittings.fitting import Fitting
from morgoth.schedule import Schedule
from morgoth.utils import timedelta_from_str
from morgoth.date_utils import from_epoch, total_seconds
import urllib2
import json

import logging
logger = logging.getLogger(__name__)

class PullGraphite(Fitting):
    """
    Fitting to periodically pull data from a graphite
        install and store the data in morgoth
    """
    _date_format = "%H:%M_%Y%m%d"
    def __init__(self,
            app,
            metric_pattern,
            graphite_url,
            period,
            metric_format,
            lag=timedelta(),
            offset=0,
            tz=None,
            user=None,
            password=None):
        """
        Initialize a pull graphite fitting

        @param metric_pattern: a graphite pattern for the metrics to pull
        @param graphite_url: the url of the graphite instance http://host[:port]
        @param period: timedelta object for how often the data should be pulled
        @param metric_format: python format string to reformat the metric name. Must have exactly one '%s' option
        @param lag: timedelta object for how much to lag from the current
            time when requesting data
        @param offset: the UTC offset to use when requesting data from graphite. see param tz
        @param tz: tz string that can be passed to dateutil.tz.gettz to
            find a time zone of the graphite server.
            NOTE: the tz conf takes precedence over the offset param
        @param user: username for basic http authentication
        @param password: password for basic http authentication
        """
        super(PullGraphite, self).__init__()
        self._app = app
        self._metric_pattern = metric_pattern
        self._graphite_url = graphite_url
        self._period = period
        self._metric_format = metric_format
        self._lag = lag
        if tz is None:
            self._tz_offset = tzoffset(None, offset)
        else:
            self._tz_offset = gettz(tz)
        self._user = user
        self._password = password
        self._schedule = Schedule(self._period, self._pull)
        self._writer = self._app.engine.get_writer()

    @classmethod
    def from_conf(cls, conf, app):
        metric_pattern = conf.get('metric_pattern', '*')
        graphite_url = conf.get('graphite_url', 'http://localhost')
        period = timedelta_from_str(conf.get('period', '5m'))
        metric_format = conf.get('metric_format', '%s')
        lag = timedelta_from_str(conf.get('lag', '1m'))
        offset = timedelta_from_str(conf.get('offset', '0m'))
        tz = conf.get('tz', None)
        user = conf.get('user', None)
        password = conf.get('password', None)
        return PullGraphite(
                app,
                metric_pattern,
                graphite_url,
                period,
                metric_format,
                lag,
                total_seconds(offset),
                tz,
                user,
                password
            )

    def start(self):
        logger.info("Starting pull graphite fitting")
        self._schedule.start_aligned()

    def stop(self):
        self._schedule.stop()

    def _pull(self):
        """
        Pull the next chunk of data from the graphite instance
        """
        stop = now() - self._lag
        start = stop - self._period
        self._pull_data(self._metric_pattern, start, stop)

    def _pull_data(self, metric_target, start, stop):
        """
        Request the data from graphite and insert into morgoth

        @param metric_target: the metric name or a pattern
        @param start: the start time
        @param stop: the stop time
        """
        g_start = start.astimezone(self._tz_offset)
        g_stop = stop.astimezone(self._tz_offset)

        url = "%s/render?target=%s&format=json&from=%s&until=%s" % (
                self._graphite_url,
                metric_target,
                g_start.strftime(self._date_format),
                g_stop.strftime(self._date_format),
            )
        logger.debug("Pulling data for %s: %s", metric_target, url)
        request = urllib2.Request(url)
        self._add_auth(request)
        data = json.load(urllib2.urlopen(request))
        if not data:
            logger.warn("Could not find metrics for target %s", metric_target)
        if not data:
            return
        for dataset in data:
            metric_name = self._metric_format % dataset['target']
            logger.debug("Inserting datapoints for %s", metric_name)
            for value, timestamp in dataset['datapoints']:
                if value is not None:
                    dt_utc = from_epoch(timestamp)
                    self._writer.insert(dt_utc, metric_name, value)

    def _add_auth(self, request):
        """
        Add the authentication headers to a request

        @param request: the request object
        """
        header = "Basic " + (
            ('%s:%s' % (self._user, self._password)).encode('base64')
        )
        if self._user and self._password:
            request.add_header("Authorization", header)

