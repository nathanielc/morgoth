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

from morgoth.config import Config
from morgoth.test.app_test_case import AppTestCase
from morgoth.data.writer import Writer
from morgoth.detectors.mgof.mgof import MGOF
from datetime import datetime, timedelta
import os
import tempfile


class MGOFTest(AppTestCase):
    """
    NOTE: this is an an end to end test
    """
    mgof_conf = """
n_bins: 15
normal_count: 1
chi2_percentage: 0.95
windows:
    - {offset: 1w, duration: 1h}
    - {offset: 2w, duration: 1h}
    - {offset: 3w, duration: 1h}
    - {offset: 4w, duration: 1h}
    - {offset: 5w, duration: 1h}
    - {offset: 6w, duration: 1h}
"""
    conf = """
mongo:
  use_sharding: false
metrics:
  .*:
    detectors:
      MGOF:
        n_bins: 15
        normal_count: 1
        chi2_percentage: 0.95
        windows:
            - {offset: 1w, duration: 1h}
            - {offset: 2w, duration: 1h}
            - {offset: 3w, duration: 1h}
            - {offset: 4w, duration: 1h}
            - {offset: 5w, duration: 1h}
            - {offset: 6w, duration: 1h}
    schedule:
      duration: 5m
      period: 5m
      delay: 1m
    """

    def create_metric_data(self, metric, start):
        self.delete_metric_data(metric)
        writer = Writer()
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
                        writer.insert(
                            start + delta,
                            metric,
                            value
                            )
        writer.flush()
    def delete_metric_data(self, metric):
        writer = Writer()
        writer.delete_metric(metric)
    def test_mgof_01(self):

        try:
            metric = 'test_mgof'
            start = datetime(2013, 9, 1)
            self.create_metric_data(metric, start)
            mgof = MGOF.from_conf(Config.loads(self.mgof_conf))

            a_start = start + timedelta(weeks=5)
            a_end = a_start + timedelta(hours=1)
            self.assertTrue(mgof.is_anomalous(metric, a_start, a_end).anomalous)

            na_start = start + timedelta(weeks=4)
            na_end = na_start + timedelta(hours=1)
            self.assertFalse(mgof.is_anomalous(metric, na_start, na_end).anomalous)
        finally:
            self.delete_metric_data(metric)

if __name__ == '__main__':
    unittest.main()

