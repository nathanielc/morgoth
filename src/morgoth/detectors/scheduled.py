
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

