
from datetime import datetime
from morgoth.anomaly_detector import AnomalyDetector
import time
import sched
import calendar

class Watcher(object):
    def __init__(self, metric, window_sizes, timefunc=time.time, delayfunc=time.sleep):
        self._metric = metric
        self._window_sizes = window_sizes
        self._timefunc = timefunc
        self._delayfunc = delayfunc
        self._ad = AnomalyDetector(self._metric)

    def _check_window(self, start, end, window_size):
        if self._ad.is_anomaloous(start, end):
            print start, end

    def watch(self):
        schedule = sched.scheduler(self._timefunc, self._delayfunc)
        start = datetime.now()
        for window_size in self._window_sizes:
            s = start + window_size
            t = calendar.timegm(s.utctimetuple())
            schedule.enterabs(t, 1, self._check_window, s, s + window_size, window_size)


