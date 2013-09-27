
import unittest
from morgoth.collector import Collector
from morgoth.watcher import Watcher
from datetime import timedelta, datetime
import morgoth.watcher
import calendar

import logging
logger = logging.getLogger('WatcherTestCase')
logging.basicConfig(level=logging.DEBUG)

class WatcherTestCase(unittest.TestCase):
    def create_metric_data(self, metric, start):
        logger.debug('Creating metric data')
        self.delete_metric_data(metric)
        c = Collector()
        for h in range(10):
            for m in range(60):
                for s in range(0, 60, 2):
                    delta = timedelta(
                        hours=h,
                        minutes=m,
                        seconds=s
                    )
                    value = m*60 + s
                    #Anomalous windows
                    if h in [3, 7]:
                       value *= 6
                    c.insert(
                        start + delta,
                        metric,
                        value
                    )

    def delete_metric_data(self, metric):
        c = Collector()
        c.delete_metric(metric)

    def test_watcher_01(self):


        metric = 'test_watcher_01'
        start = datetime(2013, 1, 1)
        now = start + timedelta(hours=5)
        end = start + timedelta(hours=9)
        self.create_metric_data(metric, start)

        morgoth.watcher.time = MockTime()
        morgoth.watcher.time._time = calendar.timegm(now.utctimetuple())
        morgoth.watcher.time._end_time = calendar.timegm(end.utctimetuple())

        training_windows = [
            (timedelta(hours=5), timedelta(minutes=60)),
            (timedelta(hours=4), timedelta(minutes=60)),
            (timedelta(hours=3), timedelta(minutes=60)),
            (timedelta(hours=2), timedelta(minutes=60)),
            (timedelta(hours=1), timedelta(minutes=60)),
        ]

        handler = MockAnomalyHandler()

        w = Watcher(
            metric,
            timedelta(minutes=60),
            training_windows,
            handler,
            5
        )

        try:
            w.watch(now)
        except StopIteration:
            pass

        self.assertEqual(1, len(handler.windows))
        self.assertEqual(
            (datetime(2013,1,1,7), datetime(2013,1,1,8)),
            handler.windows[0].range
        )



class MockAnomalyHandler:
    windows = []
    def handle(self, window):
        self.windows.append(window)

class MockTime:
    _time = 0
    _end_time = None
    def sleep(self, seconds):
        self._time += seconds
        logger.debug("Sleeping %d" % seconds)
    def time(self):
        if self._time > self._end_time:
            raise StopIteration()
        return self._time


if __name__ == '__main__':
    unittest.main()
