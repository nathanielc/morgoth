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
import logging
import unittest
from datetime import datetime, timedelta

from morgoth.config import Config
from morgoth.date_utils import utc
from morgoth.detectors.mgof.mgof import MGOF
from morgoth.test.app_test_case import AppTestCase


logger = logging.getLogger(__name__)


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
    mongo_conf = """
data_engine:
    MongoEngine:
        use_sharding: false
metrics:
  .*:
    schedule:
      duration: 5m
      period: 5m
      delay: 1m
    """

    influx_conf = """
data_engine:
    InfluxEngine:
        user: root
        password: root
        database: test
metrics:
  .*:
    schedule:
      duration: 5m
      period: 5m
      delay: 1m
    """

    def create_metric_data(self, writer, metric, start):

        writer.delete_metric(metric)

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

    def test_mgof_01(self):

        for conf in [self.mongo_conf, self.influx_conf]:
            print('Starting test for conf %s' % conf)
            app, tdir, config_path = self.set_up_app(conf)
            writer = app.engine.get_writer()
            try:
                metric = 'test_mgof'
                start = datetime(2013, 9, 1, tzinfo=utc)
                self.create_metric_data(writer, metric, start)
                mgof = MGOF.from_conf(Config.loads(self.mgof_conf), app)

                a_start = start + timedelta(weeks=5)
                a_end = a_start + timedelta(hours=1)
                anomalous, window = mgof.is_anomalous(metric, a_start, a_end)
                self.assertTrue(anomalous)
                self.assertTrue(window.anomalous)

                na_start = start + timedelta(weeks=4)
                na_end = na_start + timedelta(hours=1)
                anomalous, window = mgof.is_anomalous(metric, na_start, na_end)
                self.assertFalse(anomalous)
                self.assertFalse(window.anomalous)

            finally:
                writer.delete_metric(metric)



if __name__ == '__main__':
    unittest.main()

