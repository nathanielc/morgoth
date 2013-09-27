
from datetime import datetime, timedelta
from morgoth.anomaly_detector import AnomalyDetector
import time
import sched
import calendar


import logging
logger = logging.getLogger(__name__)

class Watcher(object):
    def __init__(self,
            metric,
            watch_interval,
            training_windows,
            anomaly_handler,
            good_window_count=2):
        """
        Create a watcher on a given metric

        @param metric: the name of the metric to watch
        @param watch_inertval: an interval to watch specified as a timedelta object
            For example:
                    a 15 min timedelta will watch every 15 min window thoughout the day
                    a 1 hour timedelta will watch every 1 hour window throughout the day
            All intervals will be aligned with midnight UTC
            The interval must be evenly divisible into 24 hours and naturally less than 24 hours.
        @param anomaly_handler: an object to handle discovered anomalies.
            Must have a method `handle(window)`
        @param training_windows: list of training windows of type (timedelta, timedelta) tuples.
            The first element is the period to check (day, week, month etc)
            The second element is the duration of the window (15min, hour etc)
        @param good_window_count: the number of windows of a certain behavior that must be found
            in order to count that behavior as normal.
            For example:
                A good_window_count of 2 will require that at least two good windows exist in the training set.
                It will also ensure that 1 bad training window will not mask another bad window.
        """
        self._metric = metric
        self._watch_interval = watch_interval
        self._anomaly_handler = anomaly_handler
        assert self._watch_interval.days == 0
        assert (24*60*60) % self._watch_interval.seconds == 0
        self._training_windows = training_windows
        self._good_window_count = good_window_count
        self._ad = AnomalyDetector(self._metric, training_windows, count_threshold=self._good_window_count)

    def _check_window(self, start, end, window_size):
        if self._ad.is_anomaloous(start, end):
            print start, end


    def watch(self, now=datetime.utcnow()):
        """
         Watch the metric for anomalies in the given interval
        """

        self._schedule = sched.scheduler(time.time, time.sleep)

        start = datetime(now.year, now.month, now.day)
        while start < now:
            start += self._watch_interval

        now_timestamp = calendar.timegm(now.utctimetuple())

        next_check = now_timestamp + (start - now).seconds

        self._schedule.enterabs(next_check, 1, self._check_window, (start,))
        self._schedule.run()

    def _check_window(self, start):
        end = start + self._watch_interval
        self._schedule.enter(
            self._watch_interval.seconds,
            1,
            self._check_window,
            (end,)
        )
        logger.debug("Checking window %s %s" % (start, end))
        window = self._ad.is_anomalous(start, end)
        if window.anomalous:
            logger.debug("Window %s %s was anomalous" % (start, end))
            self._anomaly_handler.handle(window)



