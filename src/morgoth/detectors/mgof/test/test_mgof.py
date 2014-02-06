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

import unittest
from morgoth.collector import Collector
from morgoth.anomaly_detector import AnomalyDetector
from datetime import datetime, timedelta


class AnomalyDetectorTestCase(unittest.TestCase):
    def create_metric_data(self, metric, start):
        self.delete_metric_data(metric)
        c = Collector()
        for w in range(6):
            for h in range(1):
                for m in range(60):
                    for s in range(60):
                        delta = timedelta(
                            weeks=w,
                            hours=h,
                            minutes=m,
                            seconds=s
                        )
                        value = m*60 + s
                        if w == 5:
                           value *= 4
                        c.insert(
                            start + delta,
                            metric,
                            value
                            )
    def delete_metric_data(self, metric):
        c = Collector()
        c.delete_metric(metric)
    def test_anomaly_detector_01(self):
        metric = 'test_anomaly_detector'
        start = datetime(2013, 9, 1)
        self.create_metric_data(metric, start)
        ad = AnomalyDetector(metric)

        a_start = start + timedelta(weeks=5)
        a_end = a_start + timedelta(hours=1)
        self.assertTrue(ad.is_anomalous(a_start, a_end))

        na_start = start + timedelta(weeks=4)
        na_end = na_start + timedelta(hours=1)
        self.assertFalse(ad.is_anomalous(na_start, na_end))

if __name__ == '__main__':
    unittest.main()

