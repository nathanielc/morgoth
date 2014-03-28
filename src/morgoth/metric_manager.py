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

from morgoth.schedule import Schedule
from morgoth.utils import timedelta_from_str
from morgoth.utc import now
import morgoth.detectors
import morgoth.notifiers

import gevent

import logging
logger = logging.getLogger(__name__)

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
        schedule = conf.get('schedule', None)
        if schedule:
            self._duration = timedelta_from_str(schedule.duration)
            self._period = timedelta_from_str(schedule.period)
            self._delay = timedelta_from_str(schedule.delay)
            self._aligned = schedule.get('aligned', True)

            self._schedule = Schedule(self._period, self._check_metrics, self._delay)
        else:
            self._schedule = None

        # Load consensus
        self._consensus = conf.get('consensus', 0.5)



    def add_metric(self, metric):
        """
        Add new metric to the manager
        """
        self._metrics.add(metric)

    def start(self):
        """
        Start watching the metrics for anomanlies
        """
        if not self._started and self._schedule:
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
       stop = now() - self._delay
       start = stop - self._duration
       logger.debug("Checking metrics for next window %s:%s", start, stop)
       for metric in self._metrics:
           gevent.spawn(self._check_window, metric, start, stop)

    def _check_window(self, metric, start, stop):
        """
        Check an indivdual window

        @param metric: the name of the metric
        @param start: the start time
        @param stop: the stop time
        """
        votes = 0.0
        windows = []
        for detector in self._detectors:
            window = detector.is_anomalous(metric, start, stop)
            if window.anomalous:
                votes += 1
            windows.appstop(window)

        if votes / len(self._detectors) > self._consensus:
            Meta.notify_anomalous(metric, start, stop)
            for notifier in self._notifiers:
                notifier.notify(windows)



class NullMetricManager(MetricManager):
    """
    A noop implementation of the MetricManager
    """
    def __init__(self):
        pass
    def start(self):
        pass
    def add_metric(self, metric):
        pass


