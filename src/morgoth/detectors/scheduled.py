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

from morgoth.detectors.detector import Detector
from morgoth.meta import Meta
from morgoth.schedule import Schedule
from morgoth.utc import utc, now
import gevent

import logging
logger = logging.getLogger(__name__)

class Scheduled(Detector):
    def __init__(self, period, duration):
        """
        Creat a Scheduled AD

        @param period: the timedelta period when the detection should be performed
        @param duration: the timedelta duration of the window to analyze
        """
        super(Scheduled, self).__init__()
        self._period = period
        self._duration = duration


    def start(self):
        """
        Start watching
        """
        self._sched = Schedule(self._period, self._check_metrics)
        self._sched.start_aligned()

    def _check_metrics(self):
        """
        Check if the metrics are anomalous
        """
        end = now()
        start = end - self._duration
        for metric in self._metrics:
            gevent.spawn(self._check_window, metric, start, end)

    def _check_window(self, metric, start, end):
        """
        Check an indivdual window

        @param metric: the name of the metric
        @param start: the start time
        @param end: the end time
        """
        window = self.is_anomalous(metric, start, end)
        if window.anomalous:
            Meta.notify_anomalous(window)

